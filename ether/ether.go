package ether

import (
	"context"
	"errors"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient/simulated"

	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/go-viper/mapstructure/v2"

	"github.com/poteto-go/go-alchemy-sdk/constant"
	"github.com/poteto-go/go-alchemy-sdk/encode"
	"github.com/poteto-go/go-alchemy-sdk/internal"
	"github.com/poteto-go/go-alchemy-sdk/types"
	"github.com/poteto-go/go-alchemy-sdk/utils"
	"github.com/poteto-go/go-alchemy-sdk/validate"
)

type Ether struct {
	provider        types.IAlchemyProvider
	config          EtherApiConfig
	connCount       int
	client          *ethclient.Client
	clientCreatedAt int64
	mu              *sync.Mutex
	httpClient      *http.Client // shared across all rpc.Client creations

	// simulated backend
	simBackend *simulated.Backend
	simClient  *simulated.Client
}

func NewEtherApi(provider types.IAlchemyProvider, config EtherApiConfig) types.EtherApi {
	return &Ether{
		provider:   provider,
		config:     config,
		connCount:  0,
		client:     nil,
		mu:         &sync.Mutex{},
		httpClient: utils.NewSharedHTTPClient(config.maxResponseBytes, config.requestTimeout, config.transport),
	}
}

func NewSimulatedApi(backend *simulated.Backend) types.EtherApi {
	client := backend.Client()
	return &Ether{
		// The simulated backend is in-process, so only the request timeout and
		// backoff config (used by the geth-request dispatcher) matter here.
		// A zero timeout would make every call deadline-exceed immediately.
		config:     NewEtherApiConfig("", 0, 10*time.Second, &types.DefaultBackoffConfig, nil, nil, 0, nil),
		simBackend: backend,
		connCount:  0,
		client:     nil,
		mu:         &sync.Mutex{},
		simClient:  &client,
	}
}

// it will return nil on simulated alchemy
func (ether *Ether) HttpClient() *http.Client {
	return ether.httpClient
}

func (ether *Ether) SetEthClient() error {
	ether.mu.Lock()
	defer ether.mu.Unlock()

	// simulated client must not killed or re-create
	// user should manage backend.Close w/o sdk
	if ether.simBackend != nil {
		return nil
	}

	ether.connCount++
	if ether.isClientJwsAlive() {
		return nil
	}

	ether.kill()

	rpcClient, err := ether.createRpcClient()
	if err != nil {
		ether.connCount--
		return err
	}

	ether.client = ethclient.NewClient(rpcClient)
	ether.clientCreatedAt = time.Now().Unix()
	return nil
}

// geth accepts a tight ~60s iat window; recreate before the boundary so
// clock skew or in-flight latency cannot push a request past 60s.
func (ether *Ether) isClientJwsAlive() bool {
	if ether.client == nil {
		return false
	}

	if len(ether.config.jwtSecret) == 0 {
		return ether.client != nil
	}

	now := time.Now().Unix()
	return now-ether.clientCreatedAt < constant.JwsAliveWindowSec
}

// kill all client
func (ether *Ether) kill() {
	// simulated client must not killed or re-create
	// user should manage backend.Close w/o sdk
	if ether.simBackend != nil {
		return
	}

	if ether.client == nil {
		return
	}

	ether.client.Close()
	ether.client = nil
}

func (ether *Ether) createRpcClient() (*rpc.Client, error) {
	rpcClient, err := rpc.DialOptions(
		context.Background(),
		ether.config.url,
		rpc.WithHTTPClient(ether.httpClient),
	)
	if err != nil {
		return nil, err
	}

	rpcClient.SetHeader("Content-Type", "application/json")
	rpcClient.SetHeader("Alchemy-Ethers-Sdk-Method", "send")
	for _, header := range ether.provider.CustomHeaders() {
		for key, values := range header {
			for _, value := range values {
				rpcClient.SetHeader(key, value)
			}
		}
	}

	if err := ether.generateAndSetAuthorization(rpcClient); err != nil {
		return nil, err
	}
	return rpcClient, nil
}

func (ether *Ether) generateAndSetAuthorization(rpcClient *rpc.Client) error {
	if len(ether.config.JwtSecret()) == 0 {
		return nil
	}

	jws, err := internal.GenerateJws(ether.config.jwtSecret)
	if err != nil {
		return err
	}

	rpcClient.SetHeader("Authorization", "Bearer "+jws)
	return nil
}

func (ether *Ether) Close() {
	ether.mu.Lock()
	defer ether.mu.Unlock()

	if ether.client == nil {
		return
	}

	ether.connCount--
	if ether.connCount > 0 {
		return
	}

	ether.kill()
}

func (ether *Ether) Client() types.EthClient {
	ether.mu.Lock()
	defer ether.mu.Unlock()
	if ether.simBackend != nil {
		return *ether.simClient
	}
	return ether.client
}

func (ether *Ether) Network() types.Network {
	if ether.provider == nil {
		return ""
	}
	return ether.provider.Network()
}

/*
BatchCall sends multiple JSON-RPC requests in a single HTTP round-trip using
geth's underlying rpc.Client.

Each element's Result/Error is populated in place (geth semantics): a per-request
RPC error is stored on the element's Error field, while the returned error is only
set for I/O level failures. Backoff retry therefore applies to I/O failures only.
*/
func (ether *Ether) BatchCall(elems []rpc.BatchElem) error {
	if err := ether.SetEthClient(); err != nil {
		return err
	}
	defer ether.Close()

	c, ok := ether.Client().(*ethclient.Client)
	if !ok {
		return constant.ErrUnSupportSimulatedMethod
	}

	return internal.GethRequestSingleErrorWithBackOff(
		ether.config.backoffConfig,
		ether.config.requestTimeout,
		c.Client().BatchCallContext,
		elems,
	)
}

func (ether *Ether) BlockNumber() (uint64, error) {
	err := ether.SetEthClient()
	if err != nil {
		return uint64(0), err
	}
	defer ether.Close()

	c := ether.Client()
	res, err := internal.GethRequestWithBackOff(
		ether.config.backoffConfig,
		ether.config.requestTimeout,
		c.BlockNumber,
	)
	if err != nil {
		return uint64(0), err
	}

	return res, nil
}

func (ether *Ether) GasPrice() (*big.Int, error) {
	err := ether.SetEthClient()
	if err != nil {
		return nil, err
	}
	defer ether.Close()

	c := ether.Client()
	res, err := internal.GethRequestWithBackOff(
		ether.config.backoffConfig,
		ether.config.requestTimeout,
		c.SuggestGasPrice,
	)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (ether *Ether) GetBalance(address string, blockTag string) (*big.Int, error) {
	if err := validate.BlockTag(blockTag); err != nil {
		return nil, err
	}

	balanceHex, err := ether.provider.Send(
		constant.Eth_GetBalance,
		types.RequestArgs{
			strings.ToLower(address),
			blockTag,
		},
	)
	if err != nil {
		return nil, err
	}

	balanceStr, ok := balanceHex.(string)
	if !ok {
		return nil, constant.ErrUnexpectedResponseType
	}
	balance, err := utils.FromBigHex(balanceStr)
	if err != nil {
		return nil, err
	}
	return balance, nil
}

func (ether *Ether) CodeAt(address string, blockTag string) (string, error) {
	err := ether.SetEthClient()
	if err != nil {
		return "", err
	}
	defer ether.Close()

	blockNumber, err := utils.ToBlockNumber(blockTag)
	if err != nil {
		return "", err
	}

	c := ether.Client()
	code, err := internal.GethRequestTwoArgWithBackOff(
		ether.config.backoffConfig,
		ether.config.requestTimeout,
		c.CodeAt,
		common.HexToAddress(address),
		blockNumber,
	)
	if err != nil {
		return "", err
	}

	return hexutil.Encode(code), nil
}

func (ether *Ether) CodeAtHash(address string, blockHash string) (string, error) {
	err := ether.SetEthClient()
	if err != nil {
		return "", err
	}
	defer ether.Close()

	c, ok := ether.Client().(*ethclient.Client)
	if !ok {
		return "", constant.ErrUnSupportSimulatedMethod
	}

	code, err := internal.GethRequestTwoArgWithBackOff(
		ether.config.backoffConfig,
		ether.config.requestTimeout,
		c.CodeAtHash,
		common.HexToAddress(address),
		common.HexToHash(blockHash),
	)
	if err != nil {
		return "", err
	}

	return hexutil.Encode(code), nil
}

func (ether *Ether) GetTransaction(hash string) (*gethTypes.Transaction, bool, error) {
	err := ether.SetEthClient()
	if err != nil {
		return nil, false, err
	}
	defer ether.Close()

	c := ether.Client()
	tx, isPending, err := internal.GethRequestArgWithBackOffTuple(
		ether.config.backoffConfig,
		ether.config.requestTimeout,
		c.TransactionByHash,
		common.HexToHash(hash),
	)
	if err != nil {
		return nil, isPending, err
	}

	return tx, isPending, nil
}

func (ether *Ether) StorageAt(address, position, blockTag string) (string, error) {
	err := ether.SetEthClient()
	if err != nil {
		return "", err
	}
	defer ether.Close()

	account := common.HexToAddress(address)
	key := common.HexToHash(position)
	blockNumber, err := utils.ToBlockNumber(blockTag)
	if err != nil {
		return "", err
	}

	c := ether.Client()
	res, err := internal.GethRequestThreeArgWithBackOff(
		ether.config.backoffConfig,
		ether.config.requestTimeout,
		c.StorageAt,
		account,
		key,
		blockNumber,
	)
	if err != nil {
		return "", err
	}

	return common.Bytes2Hex(res), nil
}

func (ether *Ether) GetTokenBalances(address string, params ...string) (types.TokenBalanceResponse, error) {
	paramsAny := make([]any, len(params)+1)
	paramsAny[0] = strings.ToLower(address)
	for i, param := range params {
		paramsAny[i+1] = param
	}

	result, err := ether.provider.Send(
		constant.Alchemy_GetTokenBalances,
		paramsAny,
	)
	if err != nil {
		return types.TokenBalanceResponse{}, err
	}

	resultMap, ok := result.(map[string]any)
	if !ok {
		return types.TokenBalanceResponse{}, constant.ErrUnexpectedResponseType
	}
	if balances, exists := resultMap["tokenBalances"]; exists {
		balancesAny, ok := balances.([]any)
		if !ok {
			return types.TokenBalanceResponse{}, constant.ErrUnexpectedResponseType
		}
		for _, b := range balancesAny {
			bm, ok := b.(map[string]any)
			if !ok {
				return types.TokenBalanceResponse{}, constant.ErrUnexpectedResponseType
			}
			if errStr, ok := bm["error"]; ok && errStr != nil {
				s, ok := errStr.(string)
				if !ok {
					return types.TokenBalanceResponse{}, constant.ErrUnexpectedResponseType
				}
				bm["error"] = errors.New(s)
			}
		}
	}

	var tokenBalanceResponse types.TokenBalanceResponse
	if err := mapstructure.Decode(resultMap, &tokenBalanceResponse); err != nil {
		return types.TokenBalanceResponse{}, constant.ErrFailedToMapTokenResponse
	}

	return tokenBalanceResponse, nil
}

func (ether *Ether) GetTokenMetadata(address string) (types.TokenMetadataResponse, error) {
	result, err := ether.provider.Send(
		constant.Alchemy_GetTokenMetadata,
		types.RequestArgs{
			strings.ToLower(address),
		},
	)
	if err != nil {
		return types.TokenMetadataResponse{}, err
	}

	resultMap, ok := result.(map[string]any)
	if !ok {
		return types.TokenMetadataResponse{}, constant.ErrUnexpectedResponseType
	}
	var tokenMetadata types.TokenMetadataResponse
	if err := mapstructure.Decode(resultMap, &tokenMetadata); err != nil {
		return types.TokenMetadataResponse{}, constant.ErrFailedToMapTokenResponse
	}

	return tokenMetadata, nil
}

func (ether *Ether) GetLogs(filter types.Filter) ([]types.LogResponse, error) {
	result, err := ether.provider.Send(
		constant.Eth_GetLogs,
		types.RequestArgs{
			filter,
		},
	)
	if err != nil {
		return []types.LogResponse{}, err
	}

	resultArr, ok := result.([]any)
	if !ok {
		return []types.LogResponse{}, constant.ErrUnexpectedResponseType
	}
	logs := make([]types.LogResponse, len(resultArr))
	for i, res := range resultArr {
		resultMap, ok := res.(map[string]any)
		if !ok {
			return []types.LogResponse{}, constant.ErrUnexpectedResponseType
		}
		var log types.LogResponse
		if err := mapstructure.Decode(resultMap, &log); err != nil {
			return []types.LogResponse{}, constant.ErrFailedToMapTokenResponse
		}

		logs[i] = log
	}
	return logs, nil
}

func (ether *Ether) EstimateGas(tx types.TransactionRequest) (*big.Int, error) {
	err := ether.SetEthClient()
	if err != nil {
		return nil, err
	}
	defer ether.Close()

	toAddress := common.HexToAddress(tx.To)
	value, err := utils.FromBigHex(tx.Value)
	if err != nil {
		return nil, err
	}

	c := ether.Client()
	res, err := internal.GethRequestArgWithBackOff(
		ether.config.backoffConfig,
		ether.config.requestTimeout,
		c.EstimateGas,
		ethereum.CallMsg{
			From:  common.HexToAddress(tx.From),
			To:    (&toAddress),
			Value: value,
			Data:  tx.Data,
		},
	)
	if err != nil {
		return nil, err
	}

	// NOTE: this is false positive
	// nolint:gosec
	return big.NewInt(int64(res)), nil
}

func (ether *Ether) SuggestGasPrice() (*big.Int, error) {
	err := ether.SetEthClient()
	if err != nil {
		return nil, err
	}
	defer ether.Close()

	c := ether.Client()
	res, err := internal.GethRequestWithBackOff(
		ether.config.backoffConfig,
		ether.config.requestTimeout,
		c.SuggestGasPrice,
	)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (ether *Ether) SuggestGasTipCap() (*big.Int, error) {
	if err := ether.SetEthClient(); err != nil {
		return nil, err
	}
	defer ether.Close()

	c := ether.Client()
	tip, err := internal.GethRequestWithBackOff(
		ether.config.backoffConfig,
		ether.config.requestTimeout,
		c.SuggestGasTipCap,
	)
	if err != nil {
		return nil, err
	}
	return tip, nil
}

func (ether *Ether) SuggestEIP1559Fees() (*big.Int, *big.Int, error) {
	if err := ether.SetEthClient(); err != nil {
		return nil, nil, err
	}
	defer ether.Close()

	c := ether.Client()

	tip, err := internal.GethRequestWithBackOff(
		ether.config.backoffConfig,
		ether.config.requestTimeout,
		c.SuggestGasTipCap,
	)
	if err != nil {
		return nil, nil, err
	}

	header, err := internal.GethRequestArgWithBackOff(
		ether.config.backoffConfig,
		ether.config.requestTimeout,
		c.HeaderByNumber,
		(*big.Int)(nil),
	)
	if err != nil {
		return nil, nil, err
	}
	if header.BaseFee == nil {
		return nil, nil, constant.ErrChainNotSupportEIP1559
	}

	maxFee := new(big.Int).Add(new(big.Int).Lsh(header.BaseFee, 1), tip)
	return tip, maxFee, nil
}

func (ether *Ether) Call(tx types.TransactionRequest, blockTag string) (string, error) {
	if err := validate.BlockTag(blockTag); err != nil {
		return "", err
	}

	result, err := ether.provider.Send(constant.Eth_Call, types.RequestArgs{
		tx,
		blockTag,
	})
	if err != nil {
		return "", err
	}

	resultStr, ok := result.(string)
	if !ok {
		return "", constant.ErrUnexpectedResponseType
	}
	return resultStr, nil
}

func (ether *Ether) CallContract(
	msg ethereum.CallMsg,
	blockTag string,
) ([]byte, error) {
	err := ether.SetEthClient()
	if err != nil {
		return nil, err
	}
	defer ether.Close()

	blockNumber, err := utils.ToBlockNumber(blockTag)
	if err != nil {
		return nil, err
	}

	c := ether.Client()
	output, err := internal.GethRequestTwoArgWithBackOff(
		ether.config.backoffConfig,
		ether.config.requestTimeout,
		c.CallContract,
		msg,
		blockNumber,
	)
	if err != nil {
		return nil, err
	}

	return output, nil
}

func (ether *Ether) GetTransactionReceipt(hash string) (*gethTypes.Receipt, error) {
	err := ether.SetEthClient()
	if err != nil {
		return nil, err
	}
	defer ether.Close()

	c := ether.Client()
	txReceipt, err := internal.GethRequestArgWithBackOff(
		ether.config.backoffConfig,
		ether.config.requestTimeout,
		c.TransactionReceipt,
		common.HexToHash(hash),
	)
	if err != nil {
		return nil, err
	}

	return txReceipt, nil
}

func (ether *Ether) GetTransactionReceipts(arg types.BlockNumberOrHash) ([]*gethTypes.Receipt, error) {
	result, err := ether.provider.Send(constant.Alchemy_TransactionReceipts, types.RequestArgs{
		arg,
	})
	if err != nil {
		return []*gethTypes.Receipt{}, err
	}

	resultMap, ok := result.(map[string]any)
	if !ok {
		return []*gethTypes.Receipt{}, constant.ErrUnexpectedResponseType
	}
	var txReceiptsRes types.TransactionReceiptsResponse
	if err := mapstructure.Decode(resultMap, &txReceiptsRes); err != nil {
		return []*gethTypes.Receipt{}, constant.ErrFailedToMapTransactionReceipt
	}

	txReceipts := make([]*gethTypes.Receipt, len(txReceiptsRes.Receipts))
	for i, receipt := range txReceiptsRes.Receipts {
		gethReceipt, err := utils.TransformAlchemyReceiptToGeth(receipt)
		if err != nil {
			return []*gethTypes.Receipt{}, err
		}
		txReceipts[i] = gethReceipt
	}
	return txReceipts, nil
}

func (ether *Ether) GetBlockByNumber(blockNumber string) (*gethTypes.Block, error) {
	err := ether.SetEthClient()
	if err != nil {
		return nil, err
	}
	defer ether.Close()

	bigBlockNumber, err := utils.ToBlockNumber(blockNumber)
	if err != nil {
		return nil, err
	}

	c := ether.Client()
	res, err := internal.GethRequestArgWithBackOff(
		ether.config.backoffConfig,
		ether.config.requestTimeout,
		c.BlockByNumber,
		bigBlockNumber,
	)
	if err != nil {
		return nil, err
	}

	if res == nil {
		return nil, constant.ErrResultIsNil
	}

	return res, nil
}

func (ether *Ether) GetBlockByHash(blockHash string) (*gethTypes.Block, error) {
	err := ether.SetEthClient()
	if err != nil {
		return nil, err
	}
	defer ether.Close()

	c := ether.Client()
	res, err := internal.GethRequestArgWithBackOff(
		ether.config.backoffConfig,
		ether.config.requestTimeout,
		c.BlockByHash,
		common.HexToHash(blockHash),
	)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, constant.ErrResultIsNil
	}

	return res, nil
}

/*
PendingNonceAt returns the account nonce of the given account in the pending state.
This is the nonce that should be used for the next transaction.

internal call geth
*/
func (ether *Ether) PendingNonceAt(address string) (uint64, error) {
	err := ether.SetEthClient()
	if err != nil {
		return uint64(0), err
	}
	defer ether.Close()

	c := ether.Client()
	nonce, err := internal.GethRequestArgWithBackOff(
		ether.config.backoffConfig,
		ether.config.requestTimeout,
		c.PendingNonceAt,
		common.HexToAddress(address),
	)
	if err != nil {
		return uint64(0), err
	}

	return nonce, nil
}

// send signed tx into the pending pool for execution w/ geth
func (ether *Ether) SendRawTransaction(signedTx *gethTypes.Transaction) error {
	err := ether.SetEthClient()
	if err != nil {
		return err
	}
	defer ether.Close()

	c := ether.Client()
	err = internal.GethRequestSingleErrorWithBackOff(
		ether.config.backoffConfig,
		ether.config.requestTimeout,
		c.SendTransaction,
		signedTx,
	)
	if err != nil {
		return err
	}

	return nil
}

func (ether *Ether) ChainID() (*big.Int, error) {
	err := ether.SetEthClient()
	if err != nil {
		return nil, err
	}
	defer ether.Close()

	c := ether.Client()
	res, err := internal.GethRequestWithBackOff(
		ether.config.backoffConfig,
		ether.config.requestTimeout,
		c.ChainID,
	)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (ether *Ether) PeerCount() (uint64, error) {
	err := ether.SetEthClient()
	if err != nil {
		return 0, err
	}
	defer ether.Close()

	c, ok := ether.Client().(*ethclient.Client)
	if !ok {
		return 0, constant.ErrUnSupportSimulatedMethod
	}

	res, err := internal.GethRequestWithBackOff(
		ether.config.backoffConfig,
		ether.config.requestTimeout,
		c.PeerCount,
	)
	if err != nil {
		return 0, err
	}

	return res, nil
}

// TODO: backoff
func (ether *Ether) DeployContract(
	auth *bind.TransactOpts,
	metaData *bind.MetaData,
) (*bind.DeploymentResult, error) {
	err := ether.SetEthClient()
	if err != nil {
		return nil, err
	}
	defer ether.Close()

	// set up params to deploy an instance of the metadata
	deployParams := bind.DeploymentParams{
		Contracts: []*bind.MetaData{metaData},
	}
	c := ether.Client()
	deployer := bind.DefaultDeployer(auth, c)

	// create and submit the contract deployment
	deployRes, err := bind.LinkAndDeploy(&deployParams, deployer)
	if err != nil {
		return nil, err
	}

	return deployRes, err
}

// rawBoundContract binds the address to the backend with an empty ABI.
// Callers pass pre-encoded data and decode results themselves, so bind
// only needs the address and backend.
func rawBoundContract(addr common.Address, backend bind.ContractBackend) *bind.BoundContract {
	return bind.NewBoundContract(addr, abi.ABI{}, backend, backend, backend)
}

// TODO: backoff
func (ether *Ether) ContractTransact(auth *bind.TransactOpts, contractAddress string, data []byte) (*gethTypes.Transaction, error) {
	err := ether.SetEthClient()
	if err != nil {
		return nil, err
	}
	defer ether.Close()

	instance := rawBoundContract(common.HexToAddress(contractAddress), ether.Client())

	tx, err := bind.Transact(
		instance, auth, data,
	)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (ether *Ether) Commit() (common.Hash, error) {
	if ether.simBackend == nil {
		return common.Hash{}, constant.ErrUnexpectedNilSimulatedBackend
	}

	return ether.simBackend.Commit(), nil
}

func (ether *Ether) Fork(snapShotHash common.Hash) error {
	if ether.simBackend == nil {
		return constant.ErrUnexpectedNilSimulatedBackend
	}

	return ether.simBackend.Fork(snapShotHash)
}

func (ether *Ether) simulatedMined(txHash common.Hash) (*gethTypes.Receipt, error) {
	if ether.simBackend == nil {
		return nil, constant.ErrUnexpectedNilSimulatedBackend
	}

	ether.simBackend.Commit()
	return ether.GetTransactionReceipt(txHash.Hex())
}

// TODO: support backoff
func (ether *Ether) WaitMined(ctx context.Context, txHash common.Hash) (*gethTypes.Receipt, error) {
	if ether.simBackend != nil {
		return ether.simulatedMined(txHash)
	}

	err := ether.SetEthClient()
	if err != nil {
		return nil, err
	}
	defer ether.Close()

	c := ether.Client()
	tx, err := bind.WaitMined(ctx, c, txHash)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (ether *Ether) simulatedDeployed(txHash common.Hash) (common.Address, error) {
	txReceipt, err := ether.simulatedMined(txHash)
	if err != nil {
		return common.Address{}, err
	}

	if txReceipt.ContractAddress == (common.Address{}) {
		return common.Address{}, constant.ErrUnexpectedNoContractAddress
	}

	return txReceipt.ContractAddress, nil
}

func (ether *Ether) WaitDeployed(ctx context.Context, txHash common.Hash) (common.Address, error) {
	if ether.simBackend != nil {
		return ether.simulatedDeployed(txHash)
	}

	err := ether.SetEthClient()
	if err != nil {
		return common.Address{}, err
	}
	defer ether.Close()

	c := ether.Client()
	address, err := bind.WaitDeployed(ctx, c, txHash)
	if err != nil {
		return common.Address{}, err
	}

	return address, nil
}

func (
	ether *Ether,
) ContractCall(
	contractAddress common.Address,
	opts *bind.CallOpts,
	callData []byte,
	unpack func([]byte) (any, error),
) (any, error) {
	err := ether.SetEthClient()
	if err != nil {
		return nil, err
	}
	defer ether.Close()

	instance := rawBoundContract(contractAddress, ether.Client())

	val, err := bind.Call(instance, opts, callData, unpack)
	if err != nil {
		return nil, err
	}

	return val, nil
}

func (ether *Ether) CallReadMethod(
	method []byte,
	contractAddress string,
	args ...[]byte,
) ([]byte, error) {
	data := encode.ReadCalldata(method, args...)

	contractAddr := common.HexToAddress(contractAddress)
	msg := ethereum.CallMsg{
		To:   &contractAddr,
		Data: data,
	}

	output, err := ether.CallContract(msg, "latest")
	if err != nil {
		return []byte{}, err
	}
	return output, nil
}

// assetTransfersReq is the JSON-serializable form sent to alchemy_getAssetTransfers.
type assetTransfersReq struct {
	FromBlock         string   `json:"fromBlock,omitempty"`
	ToBlock           string   `json:"toBlock,omitempty"`
	FromAddress       string   `json:"fromAddress,omitempty"`
	ToAddress         string   `json:"toAddress,omitempty"`
	ContractAddresses []string `json:"contractAddresses,omitempty"`
	Category          []string `json:"category"`
	WithMetadata      bool     `json:"withMetadata"`
	ExcludeZeroValue  bool     `json:"excludeZeroValue"`
	MaxCount          string   `json:"maxCount,omitempty"`
	PageKey           string   `json:"pageKey,omitempty"`
}

func (ether *Ether) GetAssetTransfers(params types.AssetTransfersParams) (types.AssetTransfersResponse, error) {
	req := assetTransfersReq{
		FromBlock:         params.FromBlock,
		ToBlock:           params.ToBlock,
		FromAddress:       strings.ToLower(params.FromAddress),
		ToAddress:         strings.ToLower(params.ToAddress),
		ContractAddresses: params.ContractAddresses,
		Category:          params.Category,
		WithMetadata:      params.WithMetadata,
		ExcludeZeroValue:  params.ExcludeZeroValue,
		PageKey:           params.PageKey,
	}
	if params.MaxCount > 0 {
		req.MaxCount = hexutil.EncodeUint64(uint64(params.MaxCount))
	}

	result, err := ether.provider.Send(
		constant.Alchemy_GetAssetTransfers,
		types.RequestArgs{req},
	)
	if err != nil {
		return types.AssetTransfersResponse{}, err
	}

	resultMap, ok := result.(map[string]any)
	if !ok {
		return types.AssetTransfersResponse{}, constant.ErrUnexpectedResponseType
	}

	var response types.AssetTransfersResponse
	// WeakDecode is required: the Alchemy API returns null for some string fields
	// (e.g. erc721TokenId, rawContract.address) which strict Decode cannot coerce to "".
	if err := mapstructure.WeakDecode(resultMap, &response); err != nil {
		return types.AssetTransfersResponse{}, constant.ErrFailedToMapAssetTransfers
	}
	return response, nil
}
