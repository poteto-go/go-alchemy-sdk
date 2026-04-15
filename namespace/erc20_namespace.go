package namespace

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/poteto-go/go-alchemy-sdk/types"
)

type IERC20 interface {
	/*
		BalanceOf returns the balance of the specified address.
	*/
	BalanceOf(
		contract types.ERC20ContractInstance,
		contractAddress,
		walletAddress string,
	) (*big.Int, error)
}

type ERC20 struct {
	ether types.EtherApi
}

func NewERC20Namespace(ether types.EtherApi) IERC20 {
	return &ERC20{
		ether: ether,
	}
}

func (e *ERC20) BalanceOf(
	contract types.ERC20ContractInstance,
	contractAddress,
	walletAddress string,
) (*big.Int, error) {
	callData := contract.PackBalanceOf(common.HexToAddress(walletAddress))

	unpack := func(data []byte) (any, error) {
		return contract.UnpackBalanceOf(data)
	}

	output, err := e.ether.ContractCall(
		contract,
		common.HexToAddress(contractAddress),
		&bind.CallOpts{},
		callData,
		unpack,
	)
	if err != nil {
		return nil, err
	}
	return output.(*big.Int), nil
}
