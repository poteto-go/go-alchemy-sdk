package ether

import (
	"encoding/json"
	"math/big"
	"strings"

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

	blockNumber, err := utils.FromHex(blockNumberHex)
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

	price, err := utils.FromHex(priceHex)
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

	balance, err := utils.FromBigHex(balanceHex)
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

	return code, nil
}

func (ether *Ether) GetTransaction(hash string) (types.TransactionResponse, error) {
	result, err := ether.provider.Send(core.Eth_GetTransactionByHash, hash)
	if err != nil {
		return types.TransactionResponse{}, err
	}

	var txRaw types.TransactionRawResponse
	if err := json.Unmarshal([]byte(result), &txRaw); err != nil {
		return types.TransactionResponse{}, core.ErrFailedToUnmarshalTransaction
	}

	tx, err := utils.TransformTransaction(txRaw)
	if err != nil {
		return types.TransactionResponse{}, err
	}

	return tx, nil
}
