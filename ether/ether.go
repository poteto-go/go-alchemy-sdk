package ether

import (
	"math/big"
	"strings"

	"github.com/poteto-go/go-alchemy-sdk/core"
	"github.com/poteto-go/go-alchemy-sdk/types"
	"github.com/poteto-go/go-alchemy-sdk/utils"
)

func GetBlockNumber(provider types.IAlchemyProvider) (int, error) {
	blockNumberHex, err := provider.Send(core.Eth_BlockNumber)
	if err != nil {
		return 0, err
	}

	blockNumber, err := utils.FromHex(blockNumberHex)
	if err != nil {
		return 0, err
	}
	return blockNumber, nil
}

func GetGasPrice(provider types.IAlchemyProvider) (int, error) {
	priceHex, err := provider.Send(core.Eth_GasPrice)
	if err != nil {
		return 0, err
	}

	price, err := utils.FromHex(priceHex)
	if err != nil {
		return 0, err
	}
	return price, nil
}

func GetBalance(provider types.IAlchemyProvider, address string, blockTag string) (*big.Int, error) {
	if err := utils.ValidateBlockTag(blockTag); err != nil {
		return big.NewInt(0), err
	}

	balanceHex, err := provider.Send(
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
