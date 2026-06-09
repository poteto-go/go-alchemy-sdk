package batch

import (
	"math/big"

	"github.com/poteto-go/go-alchemy-sdk/constant"
	"github.com/poteto-go/go-alchemy-sdk/decode"
	"github.com/poteto-go/go-alchemy-sdk/encode"
	"github.com/poteto-go/go-alchemy-sdk/validate"
)

// ERC20Batch queues ERC-20 contract reads onto its Batcher. Mirrors the
// read-only methods of the ERC20 namespace and validates the same inputs.
type ERC20Batch struct {
	b *Batcher
}

func (e *ERC20Batch) BalanceOf(contractAddress, walletAddress string) *Result[*big.Int] {
	if err := validate.Addresses(contractAddress, walletAddress); err != nil {
		return failed[*big.Int](err)
	}
	return AddCall(e.b, contractAddress, constant.BalanceOfFnSignature, decode.Uint256, encode.ABIAddress(walletAddress))
}

func (e *ERC20Batch) TotalSupply(contractAddress string) *Result[*big.Int] {
	if err := validate.Address(contractAddress); err != nil {
		return failed[*big.Int](err)
	}
	return AddCall(e.b, contractAddress, constant.TotalSupplyFnSignature, decode.Uint256)
}

func (e *ERC20Batch) Allowance(contractAddress, owner, spender string) *Result[*big.Int] {
	if err := validate.Addresses(contractAddress, owner, spender); err != nil {
		return failed[*big.Int](err)
	}
	return AddCall(e.b, contractAddress, constant.AllowanceFnSignature, decode.Uint256,
		encode.ABIAddress(owner), encode.ABIAddress(spender))
}

func (e *ERC20Batch) Name(contractAddress string) *Result[string] {
	if err := validate.Address(contractAddress); err != nil {
		return failed[string](err)
	}
	return AddCall(e.b, contractAddress, constant.NameFnSignature, decode.ABIString)
}

func (e *ERC20Batch) Symbol(contractAddress string) *Result[string] {
	if err := validate.Address(contractAddress); err != nil {
		return failed[string](err)
	}
	return AddCall(e.b, contractAddress, constant.SymbolFnSignature, decode.ABIString)
}

func (e *ERC20Batch) Decimals(contractAddress string) *Result[uint8] {
	if err := validate.Address(contractAddress); err != nil {
		return failed[uint8](err)
	}
	return AddCall(e.b, contractAddress, constant.DecimalsFnSignature, decode.Uint8)
}
