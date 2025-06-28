package ether

import (
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
