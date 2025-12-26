package ether

import (
	"context"
	"errors"
	"math/big"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"

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
	provider  types.IAlchemyProvider
	config    EtherApiConfig
	connCount int
	client    *ethclient.Client
	mu        *sync.Mutex
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
	if ether.client != nil {
		return nil
	}

	ether.mu.Lock()
	defer ether.mu.Unlock()

	rpcClient, err := rpc.Dial(ether.config.url)
	if err != nil {
		return err
	}

	ether.client = ethclient.NewClient(rpcClient)
	return nil
}

func (ether *Ether) GetEthClient() (*ethclient.Client, error) {
	err := ether.SetEthClient()
	if err != nil {
		return nil, err
	}

	return ether.client, nil
}

func (ether *Ether) Close() {
	if ether.client == nil {
		return
	}

	ether.connCount -= 1
	if ether.connCount > 0 {
		return
	}

	ether.client.Close()
	ether.client = nil
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

	balance, err := utils.FromBigHex(balanceHex.(string))
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

	return common.Bytes2Hex(code), nil
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

	return common.Bytes2Hex(code), nil
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
	paramsAny := []any{address}
	for _, param := range params {
		paramsAny = append(paramsAny, param)
	}

	result, err := ether.provider.Send(
		constant.Alchemy_GetTokenBalances,
		paramsAny,
	)
	if err != nil {
		return types.TokenBalanceResponse{}, err
	}

	resultMap := result.(map[string]any)
	if balances, ok := resultMap["tokenBalances"]; ok {
		for _, balance := range balances.([]map[string]any) {
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

	resultMap := result.(map[string]any)
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

	resultArr := result.([]any)
	logs := make([]types.LogResponse, len(resultArr))
	for i, res := range resultArr {
		resultMap := res.(map[string]any)
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

	return result.(string), nil
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

	resultMap := result.(map[string]any)
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
