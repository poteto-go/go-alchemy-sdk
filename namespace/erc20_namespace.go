package namespace

import (
	"math/big"

	"github.com/poteto-go/go-alchemy-sdk/constant"
	"github.com/poteto-go/go-alchemy-sdk/types"
	"github.com/poteto-go/go-alchemy-sdk/utils"
	"github.com/poteto-go/go-alchemy-sdk/validate"
)

type IERC20 interface {
	// BalanceOf returns the balance of the specified address.
	BalanceOf(
		contractAddress,
		walletAddress string,
	) (*big.Int, error)

	// TotalSupply returns the total supply of the token.
	TotalSupply(contractAddress string) (*big.Int, error)

	// Allowance returns the amount of tokens the spender is allowed to spend on behalf of the owner.
	Allowance(contractAddress, owner, spender string) (*big.Int, error)

	// Name returns the name of the token.
	Name(contractAddress string) (string, error)

	// Symbol returns the symbol of the token.
	Symbol(contractAddress string) (string, error)

	// Decimals returns the number of decimals the token uses.
	Decimals(contractAddress string) (uint8, error)
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
	contractAddress,
	walletAddress string,
) (*big.Int, error) {
	if err := validate.Addresses(contractAddress, walletAddress); err != nil {
		return nil, err
	}
	output, err := e.ether.CallReadMethod(
		constant.BalanceOfFnSignature,
		contractAddress,
		utils.EncodeABIAddress(walletAddress),
	)
	if err != nil {
		return nil, err
	}

	return utils.DecodeUint256(output)
}

func (e *ERC20) TotalSupply(contractAddress string) (*big.Int, error) {
	if err := validate.Address(contractAddress); err != nil {
		return nil, err
	}
	output, err := e.ether.CallReadMethod(
		constant.TotalSupplyFnSignature,
		contractAddress,
	)
	if err != nil {
		return nil, err
	}

	return utils.DecodeUint256(output)
}

func (e *ERC20) Allowance(contractAddress, owner, spender string) (*big.Int, error) {
	if err := validate.Addresses(contractAddress, owner, spender); err != nil {
		return nil, err
	}
	output, err := e.ether.CallReadMethod(
		constant.AllowanceFnSignature,
		contractAddress,
		utils.EncodeABIAddress(owner),
		utils.EncodeABIAddress(spender),
	)
	if err != nil {
		return nil, err
	}

	return utils.DecodeUint256(output)
}

func (e *ERC20) Name(contractAddress string) (string, error) {
	if err := validate.Address(contractAddress); err != nil {
		return "", err
	}
	output, err := e.ether.CallReadMethod(
		constant.NameFnSignature,
		contractAddress,
	)
	if err != nil {
		return "", err
	}

	return utils.DecodeABIString(output)
}

func (e *ERC20) Symbol(contractAddress string) (string, error) {
	if err := validate.Address(contractAddress); err != nil {
		return "", err
	}
	output, err := e.ether.CallReadMethod(
		constant.SymbolFnSignature,
		contractAddress,
	)
	if err != nil {
		return "", err
	}

	return utils.DecodeABIString(output)
}

func (e *ERC20) Decimals(contractAddress string) (uint8, error) {
	if err := validate.Address(contractAddress); err != nil {
		return 0, err
	}
	output, err := e.ether.CallReadMethod(
		constant.DecimalsFnSignature,
		contractAddress,
	)
	if err != nil {
		return 0, err
	}

	return utils.DecodeUint8(output)
}
