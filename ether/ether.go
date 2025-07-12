package ether

import (
	"errors"
	"math/big"
	"strings"

	"github.com/go-viper/mapstructure/v2"

	"github.com/poteto-go/go-alchemy-sdk/core"
	"github.com/poteto-go/go-alchemy-sdk/types"
	"github.com/poteto-go/go-alchemy-sdk/utils"
)

type EtherApi interface {
	/* get  the number of the most recent block. */
	GetBlockNumber() (int, error)

	/* Returns the best guess of the current gas price to use in a transaction. */
	GetGasPrice() (int, error)

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

		NOTE: This is an alias for {@link TransactNamespace.getTransaction}.
	*/
	GetTransaction(hash string) (types.TransactionResponse, error)

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
		Returns an estimate of the amount of gas that would be required to submit transaction to the network.

		An estimate may not be accurate since there could be another transaction on the network that was not accounted for,
		but after being mined affects the relevant state.
		This is an alias for {@link TransactNamespace.estimateGas}.
	*/
	EstimateGas(transaction types.TransactionRequest) (*big.Int, error)
}

type Ether struct {
	provider types.IAlchemyProvider
}

func NewEtherApi(provider types.IAlchemyProvider) EtherApi {
	return &Ether{
		provider: provider,
	}
}

func (ether *Ether) GetBlockNumber() (int, error) {
	blockNumberHex, err := ether.provider.Send(core.Eth_BlockNumber)
	if err != nil {
		return 0, err
	}

	blockNumber, err := utils.FromHex(blockNumberHex.(string))
	if err != nil {
		return 0, err
	}
	return blockNumber, nil
}

func (ether *Ether) GetGasPrice() (int, error) {
	priceHex, err := ether.provider.Send(core.Eth_GasPrice)
	if err != nil {
		return 0, err
	}

	price, err := utils.FromHex(priceHex.(string))
	if err != nil {
		return 0, err
	}
	return price, nil
}

func (ether *Ether) GetBalance(address string, blockTag string) (*big.Int, error) {
	if err := utils.ValidateBlockTag(blockTag); err != nil {
		return big.NewInt(0), err
	}

	balanceHex, err := ether.provider.Send(
		core.Eth_GetBalance,
		strings.ToLower(address),
		blockTag,
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
		core.Eth_GetCode,
		strings.ToLower(address),
		blockTag,
	)
	if err != nil {
		return "", err
	}

	return code.(string), nil
}

func (ether *Ether) GetTransaction(hash string) (types.TransactionResponse, error) {
	result, err := ether.provider.Send(core.Eth_GetTransactionByHash, hash)
	if err != nil {
		return types.TransactionResponse{}, err
	}

	var txRaw types.TransactionRawResponse
	if err := mapstructure.Decode(result, &txRaw); err != nil {
		return types.TransactionResponse{}, core.ErrFailedToMapTransaction
	}

	tx, err := utils.TransformTransaction(txRaw)
	if err != nil {
		return types.TransactionResponse{}, err
	}

	return tx, nil
}

func (ether *Ether) GetStorageAt(address, position, blockTag string) (string, error) {
	if err := utils.ValidateBlockTag(blockTag); err != nil {
		return "", err
	}

	result, err := ether.provider.Send(
		core.Eth_GetStorageAt,
		strings.ToLower(address),
		position,
		blockTag,
	)
	if err != nil {
		return "", err
	}

	return result.(string), nil
}

func (ether *Ether) GetTokenBalances(address string, params ...string) (types.TokenBalanceResponse, error) {
	params = append([]string{address}, params...)

	result, err := ether.provider.Send(
		core.Alchemy_GetTokenBalances,
		params...,
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
		return types.TokenBalanceResponse{}, core.ErrFailedToMapTokenResponse
	}

	return tokenBalanceResponse, nil
}

func (ether *Ether) GetTokenMetadata(address string) (types.TokenMetadataResponse, error) {
	result, err := ether.provider.Send(
		core.Alchemy_GetTokenMetadata,
		strings.ToLower(address),
	)
	if err != nil {
		return types.TokenMetadataResponse{}, err
	}

	resultMap := result.(map[string]any)
	var tokenMetadata types.TokenMetadataResponse
	if err := mapstructure.Decode(resultMap, &tokenMetadata); err != nil {
		return types.TokenMetadataResponse{}, core.ErrFailedToMapTokenResponse
	}

	return tokenMetadata, nil
}

func (ether *Ether) EstimateGas(transaction types.TransactionRequest) (*big.Int, error) {
	result, err := ether.provider.SendTransaction(core.Eth_EstimateGas, transaction)
	if err != nil {
		return big.NewInt(0), err
	}

	estimatedGas, err := utils.FromBigHex(result.(string))
	if err != nil {
		return big.NewInt(0), err
	}

	return estimatedGas, nil
}
