package ether

import (
	"errors"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
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

type EtherApi interface {
	GetEthClient() (*ethclient.Client, error)

	/* get  the number of the most recent block. */
	BlockNumber() (uint64, error)

	/* Returns the best guess of the current gas price to use in a transaction. */
	GasPrice() (*big.Int, error)

	/* Returns the balance of a given address as of the provided block. */
	GetBalance(address string, blockTag string) (*big.Int, error)

	/*
		Returns the contract code of the provided address at the block.
		If there is no contract deployed, the result is 0x.
	*/
	GetCode(address, blockTag string) (string, error)

	/*
		Returns the transaction with hash or null if the transaction is unknown.

		If a transaction has not been mined, this method will search the
		transaction pool. Various backends may have more restrictive transaction
		pool access (e.g. if the gas price is too low or the transaction was only
		recently sent and not yet indexed) in which case this method may also return null.

		internal call geth.TransactionByHash

		NOTE: This is an alias for {@link TransactNamespace.getTransaction}.
	*/
	GetTransaction(hash string) (tx *gethTypes.Transaction, isPending bool, err error)

	/*
		Return the value of the provided position at the provided address, at the provided block in `Bytes32` format.
		For inspecting solidity code.
	*/
	GetStorageAt(address, position, blockTag string) (string, error)

	/*
		Returns the ERC-20 token balances for a specific owner address w, w/o params
	*/
	GetTokenBalances(address string, params ...string) (types.TokenBalanceResponse, error)

	/* Returns metadata for a given token contract address. */
	GetTokenMetadata(address string) (types.TokenMetadataResponse, error)

	/*
		Returns an array of logs that match the provided filter.
	*/
	GetLogs(filter types.Filter) ([]types.LogResponse, error)

	/*
		Returns an estimate of the amount of gas that would be required to submit transaction to the network.

		An estimate may not be accurate since there could be another transaction on the network that was not accounted for,
		but after being mined affects the relevant state.
		This is an alias for {@link TransactNamespace.estimateGas}.
	*/
	EstimateGas(tx types.TransactionRequest) (*big.Int, error)

	/*
		Returns the result of executing the transaction, using call.
		A call does not require any ether, but cannot change any state.
		This is useful for calling getters on Contracts.
	*/
	Call(tx types.TransactionRequest, blockTag string) (string, error)

	/*
		Null if the tx has not been mined.
		Returns the transaction receipt for hash.
		To stall until the transaction has been mined, consider the waitForTransaction method below.
	*/
	GetTransactionReceipt(hash string) (*gethTypes.Receipt, error)

	/*
		An enhanced API that gets all transaction receipts for a given block by number or block hash.
		Returns geth's Receipt.
	*/
	GetTransactionReceipts(arg types.TransactionReceiptsArg) ([]*gethTypes.Receipt, error)

	/*
		Simple wrapper around eth_getBlockByNumber.
		This returns the complete block information for the provided block number.
	*/
	GetBlockByNumber(blockNumber string) (*gethTypes.Block, error)

	/*
		Simple wrapper around eth_getBlockByHash.
		This returns the complete block information for the provided block hash.
	*/
	GetBlockByHash(blockHash string) (*gethTypes.Block, error)
}

type Ether struct {
	provider types.IAlchemyProvider
	config   EtherApiConfig
}

func NewEtherApi(provider types.IAlchemyProvider, config EtherApiConfig) EtherApi {
	return &Ether{
		provider: provider,
		config:   config,
	}
}

func (ether *Ether) GetEthClient() (*ethclient.Client, error) {
	rpcClient, err := rpc.Dial(ether.config.url)
	if err != nil {
		return nil, err
	}

	return ethclient.NewClient(rpcClient), nil
}

func (ether *Ether) BlockNumber() (uint64, error) {
	client, err := ether.GetEthClient()
	if err != nil {
		return uint64(0), err
	}
	defer client.Close()

	res, err := internal.GethRequestWithBackOff(
		ether.config.backoffConfig,
		ether.config.requestTimeout,
		client.BlockNumber,
	)
	if err != nil {
		return uint64(0), err
	}

	return res, nil
}

func (ether *Ether) GasPrice() (*big.Int, error) {
	client, err := ether.GetEthClient()
	if err != nil {
		return big.NewInt(0), err
	}
	defer client.Close()

	res, err := internal.GethRequestWithBackOff(
		ether.config.backoffConfig,
		ether.config.requestTimeout,
		client.SuggestGasPrice,
	)
	if err != nil {
		return big.NewInt(0), err
	}

	return res, nil
}

func (ether *Ether) GetBalance(address string, blockTag string) (*big.Int, error) {
	if err := utils.ValidateBlockTag(blockTag); err != nil {
		return big.NewInt(0), err
	}

	balanceHex, err := ether.provider.Send(
		constant.Eth_GetBalance,
		types.RequestArgs{
			strings.ToLower(address),
			blockTag,
		},
	)
	if err != nil {
		return big.NewInt(0), err
	}

	balance, err := utils.FromBigHex(balanceHex.(string))
	if err != nil {
		return big.NewInt(0), err
	}
	return balance, nil
}

func (ether *Ether) GetCode(address, blockTag string) (string, error) {
	if err := utils.ValidateBlockTag(blockTag); err != nil {
		return "", err
	}

	code, err := ether.provider.Send(
		constant.Eth_GetCode,
		types.RequestArgs{
			strings.ToLower(address),
			blockTag,
		},
	)
	if err != nil {
		return "", err
	}

	return code.(string), nil
}

func (ether *Ether) GetTransaction(hash string) (*gethTypes.Transaction, bool, error) {
	client, err := ether.GetEthClient()
	if err != nil {
		return nil, false, err
	}
	defer client.Close()

	tx, isPending, err := internal.GethRequestArgWithBackOffTuple(
		ether.config.backoffConfig,
		ether.config.requestTimeout,
		client.TransactionByHash,
		common.HexToHash(hash),
	)
	if err != nil {
		return nil, isPending, err
	}

	return tx, isPending, nil
}

func (ether *Ether) GetStorageAt(address, position, blockTag string) (string, error) {
	if err := utils.ValidateBlockTag(blockTag); err != nil {
		return "", err
	}

	result, err := ether.provider.Send(
		constant.Eth_GetStorageAt,
		types.RequestArgs{
			strings.ToLower(address),
			position,
			blockTag,
		},
	)
	if err != nil {
		return "", err
	}

	return result.(string), nil
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
	client, err := ether.GetEthClient()
	if err != nil {
		return big.NewInt(0), err
	}
	defer client.Close()

	toAddress := common.HexToAddress(tx.To)
	value, err := utils.FromBigHex(tx.Value)
	if err != nil {
		return big.NewInt(0), err
	}

	res, err := internal.GethRequestArgWithBackOff(
		ether.config.backoffConfig,
		ether.config.requestTimeout,
		client.EstimateGas,
		ethereum.CallMsg{
			From:  common.HexToAddress(tx.From),
			To:    (&toAddress),
			Value: value,
		},
	)
	if err != nil {
		return big.NewInt(0), err
	}

	// NOTE: this is false positive
	// nolint:gosec
	return big.NewInt(int64(res)), nil
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
	client, err := ether.GetEthClient()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	txReceipt, err := internal.GethRequestArgWithBackOff(
		ether.config.backoffConfig,
		ether.config.requestTimeout,
		client.TransactionReceipt,
		common.HexToHash(hash),
	)
	if err != nil {
		return nil, err
	}

	return txReceipt, nil
}

func (ether *Ether) GetTransactionReceipts(arg types.TransactionReceiptsArg) ([]*gethTypes.Receipt, error) {
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
	client, err := ether.GetEthClient()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	bigBlockNumber, err := utils.FromBigHex(blockNumber)
	if err != nil {
		return nil, err
	}

	res, err := internal.GethRequestArgWithBackOff(
		ether.config.backoffConfig,
		ether.config.requestTimeout,
		client.BlockByNumber,
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
	client, err := ether.GetEthClient()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	res, err := internal.GethRequestArgWithBackOff(
		ether.config.backoffConfig,
		ether.config.requestTimeout,
		client.BlockByHash,
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
