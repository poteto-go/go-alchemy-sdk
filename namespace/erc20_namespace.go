package namespace

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/poteto-go/go-alchemy-sdk/constant"
	"github.com/poteto-go/go-alchemy-sdk/types"
	"github.com/poteto-go/go-alchemy-sdk/utils"
	"golang.org/x/crypto/sha3"
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
	output, err := e.callReadMethod(
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
	output, err := e.callReadMethod(
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
	output, err := e.callReadMethod(
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
	output, err := e.callReadMethod(
		constant.NameFnSignature,
		contractAddress,
	)
	if err != nil {
		return "", err
	}

	return utils.DecodeABIString(output)
}

func (e *ERC20) Symbol(contractAddress string) (string, error) {
	output, err := e.callReadMethod(
		constant.SymbolFnSignature,
		contractAddress,
	)
	if err != nil {
		return "", err
	}

	return utils.DecodeABIString(output)
}

func (e *ERC20) Decimals(contractAddress string) (uint8, error) {
	output, err := e.callReadMethod(
		constant.DecimalsFnSignature,
		contractAddress,
	)
	if err != nil {
		return 0, err
	}

	decimals := new(big.Int).SetBytes(output).Uint64()
	if decimals > 255 {
		return 0, fmt.Errorf("decimals overflow: %d", decimals)
	}
	return uint8(decimals), nil
}

func (e *ERC20) callReadMethod(
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

	output, err := e.ether.CallContract(msg, "latest")
	if err != nil {
		return []byte{}, err
	}
	return output, nil
}
