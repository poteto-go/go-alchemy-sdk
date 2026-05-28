package namespace

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/poteto-go/go-alchemy-sdk/constant"
	"github.com/poteto-go/go-alchemy-sdk/types"
	"github.com/poteto-go/go-alchemy-sdk/utils"
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
	output, err := e.ether.CallReadMethod(
		constant.BalanceOfFnSignature,
		contractAddress,
		common.LeftPadBytes(common.HexToAddress(walletAddress).Bytes(), 32),
	)
	if err != nil {
		return nil, err
	}

	outputInt := new(big.Int)
	return outputInt.SetBytes(output), nil
}

func (e *ERC20) TotalSupply(contractAddress string) (*big.Int, error) {
	output, err := e.ether.CallReadMethod(
		constant.TotalSupplyFnSignature,
		contractAddress,
	)
	if err != nil {
		return nil, err
	}

	outputInt := new(big.Int)
	return outputInt.SetBytes(output), nil
}

func (e *ERC20) Allowance(contractAddress, owner, spender string) (*big.Int, error) {
	output, err := e.ether.CallReadMethod(
		constant.AllowanceFnSignature,
		contractAddress,
		common.LeftPadBytes(common.HexToAddress(owner).Bytes(), 32),
		common.LeftPadBytes(common.HexToAddress(spender).Bytes(), 32),
	)
	if err != nil {
		return nil, err
	}

	outputInt := new(big.Int)
	return outputInt.SetBytes(output), nil
}

func (e *ERC20) Name(contractAddress string) (string, error) {
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
	output, err := e.ether.CallReadMethod(
		constant.DecimalsFnSignature,
		contractAddress,
	)
	if err != nil {
		return 0, err
	}

	out := new(big.Int).SetBytes(output)
	if out.BitLen() > 8 {
		return 0, fmt.Errorf("decimals overflow: %s", out.String())
	}
	b := out.Bytes()
	if len(b) == 0 {
		return 0, nil
	}
	return b[len(b)-1], nil
}
