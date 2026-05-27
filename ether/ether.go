package ether

import (
	"context"
	"errors"
	"io"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"golang.org/x/crypto/sha3"

	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/go-viper/mapstructure/v2"

	"github.com/poteto-go/go-alchemy-sdk/constant"
	"github.com/poteto-go/go-alchemy-sdk/internal"
	"github.com/poteto-go/go-alchemy-sdk/types"
	"github.com/poteto-go/go-alchemy-sdk/utils"
)

type Ether struct {
	provider        types.IAlchemyProvider
	config          EtherApiConfig
	connCount       int
	client          *ethclient.Client
	clientCreatedAt int64
	mu              *sync.Mutex
}

func NewEtherApi(provider types.IAlchemyProvider, config EtherApiConfig) types.EtherApi {
	return &Ether{
		provider:  provider,
		config:    config,
		connCount: 0,
		client:    nil,
		mu:        &sync.Mutex{},
	}
}

func (ether *Ether) SetEthClient() error {
	ether.connCount += 1
	if ether.isClientJwsAlive() {
		return nil
	}

	ether.kill()

	ether.mu.Lock()
	defer ether.mu.Unlock()

	rpcClient, err := ether.createRpcClient()
	if err != nil {
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
	if ether.client == nil {
		return
	}

	ether.client.Close()
	ether.client = nil
}

// limitedReadCloser wraps a limited reader while preserving the original closer.
type limitedReadCloser struct {
	io.Reader
	io.Closer
}

// limitedTransport wraps http.RoundTripper to cap response bodies at maxBytes.
type limitedTransport struct {
	underlying http.RoundTripper
	maxBytes   int64
}

// RoundTrip intercepts every HTTP response from the geth rpc.Client call chain:
//
//	geth method call (e.g. BlockNumber)
//	  └─ ethclient.Client → rpc.Client
//	       └─ httpConn.doRequest()           [rpc/http.go:229]
//	            └─ hc.client.Do(req)         ← net/http http.Client.Do
//	                 └─ Transport.RoundTrip(req)   ← called automatically by Go's http.Client
//	                      └─ limitedTransport.RoundTrip()  [ether/ether.go]
//	                           ├─ t.underlying.RoundTrip(req)  ← actual HTTP communication
//	                           └─ resp.Body = LimitReader(resp.Body, maxBytes)  ← wrapped here
func (t *limitedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := t.underlying.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	resp.Body = limitedReadCloser{
		Reader: io.LimitReader(resp.Body, t.maxBytes),
		Closer: resp.Body,
	}
	return resp, nil
}

func (ether *Ether) createRpcClient() (*rpc.Client, error) {
	maxBytes := ether.config.maxResponseBytes
	if maxBytes == 0 {
		maxBytes = types.DefaultMaxResponseBytes
	}

	httpClient := &http.Client{
		Transport: &limitedTransport{
			underlying: http.DefaultTransport,
			maxBytes:   maxBytes,
		},
	}

	rpcClient, err := rpc.DialOptions(
		context.Background(),
		ether.config.url,
		rpc.WithHTTPClient(httpClient),
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
	if ether.client == nil {
		return
	}

	ether.connCount -= 1
	if ether.connCount > 0 {
		return
	}

	ether.kill()
}

func (ether *Ether) Client() *ethclient.Client {
	return ether.client
}

func (ether *Ether) BlockNumber() (uint64, error) {
	err := ether.SetEthClient()
	if err != nil {
		return uint64(0), err
	}
	defer ether.Close()

	res, err := internal.GethRequestWithBackOff(
		ether.config.backoffConfig,
		ether.config.requestTimeout,
		ether.client.BlockNumber,
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

	res, err := internal.GethRequestWithBackOff(
		ether.config.backoffConfig,
		ether.config.requestTimeout,
		ether.client.SuggestGasPrice,
	)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (ether *Ether) GetBalance(address string, blockTag string) (*big.Int, error) {
	if err := utils.ValidateBlockTag(blockTag); err != nil {
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

	code, err := internal.GethRequestTwoArgWithBackOff(
		ether.config.backoffConfig,
		ether.config.requestTimeout,
		ether.client.CodeAt,
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

	code, err := internal.GethRequestTwoArgWithBackOff(
		ether.config.backoffConfig,
		ether.config.requestTimeout,
		ether.client.CodeAtHash,
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

	tx, isPending, err := internal.GethRequestArgWithBackOffTuple(
		ether.config.backoffConfig,
		ether.config.requestTimeout,
		ether.client.TransactionByHash,
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

	res, err := internal.GethRequestThreeArgWithBackOff(
		ether.config.backoffConfig,
		ether.config.requestTimeout,
		ether.client.StorageAt,
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
		balanceSlice, ok := balances.([]map[string]any)
		if !ok {
			return types.TokenBalanceResponse{}, constant.ErrUnexpectedResponseType
		}
		for _, balance := range balanceSlice {
			if errStr, ok := balance["error"]; ok && errStr != nil {
				balance["error"] = errors.New(errStr.(string))
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

	res, err := internal.GethRequestArgWithBackOff(
		ether.config.backoffConfig,
		ether.config.requestTimeout,
		ether.client.EstimateGas,
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

	res, err := internal.GethRequestWithBackOff(
		ether.config.backoffConfig,
		ether.config.requestTimeout,
		ether.client.SuggestGasPrice,
	)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (ether *Ether) Call(tx types.TransactionRequest, blockTag string) (string, error) {
	if err := utils.ValidateBlockTag(blockTag); err != nil {
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

	output, err := internal.GethRequestTwoArgWithBackOff(
		ether.config.backoffConfig,
		ether.config.requestTimeout,
		ether.client.CallContract,
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

	txReceipt, err := internal.GethRequestArgWithBackOff(
		ether.config.backoffConfig,
		ether.config.requestTimeout,
		ether.client.TransactionReceipt,
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

	res, err := internal.GethRequestArgWithBackOff(
		ether.config.backoffConfig,
		ether.config.requestTimeout,
		ether.client.BlockByNumber,
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

	res, err := internal.GethRequestArgWithBackOff(
		ether.config.backoffConfig,
		ether.config.requestTimeout,
		ether.client.BlockByHash,
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

	nonce, err := internal.GethRequestArgWithBackOff(
		ether.config.backoffConfig,
		ether.config.requestTimeout,
		ether.client.PendingNonceAt,
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

	err = internal.GethRequestSingleErrorWithBackOff(
		ether.config.backoffConfig,
		ether.config.requestTimeout,
		ether.client.SendTransaction,
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

	res, err := internal.GethRequestWithBackOff(
		ether.config.backoffConfig,
		ether.config.requestTimeout,
		ether.client.ChainID,
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

	res, err := internal.GethRequestWithBackOff(
		ether.config.backoffConfig,
		ether.config.requestTimeout,
		ether.client.PeerCount,
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
	deployer := bind.DefaultDeployer(auth, ether.client)

	// create and submit the contract deployment
	deployRes, err := bind.LinkAndDeploy(&deployParams, deployer)
	if err != nil {
		return nil, err
	}

	return deployRes, err
}

// TODO: backoff
func (ether *Ether) ContractTransact(auth *bind.TransactOpts, contract types.ContractInstance, contractAddress string, data []byte) (*gethTypes.Transaction, error) {
	if contract == nil {
		return nil, constant.ErrContractInstanceIsNil
	}

	err := ether.SetEthClient()
	if err != nil {
		return nil, err
	}
	defer ether.Close()

	instance := contract.Instance(ether.client, common.HexToAddress(contractAddress))

	tx, err := bind.Transact(
		instance, auth, data,
	)
	if err != nil {
		return nil, err
	}

	return tx, nil

	// return ether.WaitMined(tx.Hash())
}

// TODO: support backoff
func (ether *Ether) WaitMined(txHash common.Hash) (*gethTypes.Receipt, error) {
	err := ether.SetEthClient()
	if err != nil {
		return nil, err
	}
	defer ether.Close()

	tx, err := bind.WaitMined(context.Background(), ether.client, txHash)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (ether *Ether) WaitDeployed(txHash common.Hash) (common.Address, error) {
	err := ether.SetEthClient()
	if err != nil {
		return common.Address{}, err
	}
	defer ether.Close()

	address, err := bind.WaitDeployed(context.Background(), ether.client, txHash)
	if err != nil {
		return common.Address{}, err
	}

	return address, nil
}

func (
	ether *Ether,
) ContractCall(
	contract types.ContractInstance,
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

	instance := contract.Instance(ether.client, contractAddress)

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
	hash := sha3.NewLegacyKeccak256()
	if _, err := hash.Write(method); err != nil {
		return nil, err
	}
	methodID := hash.Sum(nil)[:4]

	data := make([]byte, 0, 4+len(args)*32)
	data = append(data, methodID...)
	for _, arg := range args {
		data = append(data, arg...)
	}

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
