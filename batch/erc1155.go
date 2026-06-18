package batch

import (
	"math/big"

	"github.com/poteto-go/go-alchemy-sdk/constant"
	"github.com/poteto-go/go-alchemy-sdk/decode"
	"github.com/poteto-go/go-alchemy-sdk/encode"
	"github.com/poteto-go/go-alchemy-sdk/validate"
)

// ERC1155Batch queues ERC-1155 contract reads onto its Batcher. Mirrors the
// read-only methods of the Erc1155 namespace and validates the same inputs.
type ERC1155Batch struct {
	b *Batcher
}

func (e *ERC1155Batch) BalanceOfToken(contractAddress, account string, tokenId *big.Int) *Result[*big.Int] {
	if err := validate.Addresses(contractAddress, account); err != nil {
		return failed[*big.Int](err)
	}
	if err := validate.Uint256(tokenId); err != nil {
		return failed[*big.Int](err)
	}
	return AddCall(e.b, contractAddress, constant.BalanceOfTokenFnSignature, decode.Uint256,
		encode.ABIAddress(account), encode.ABIUint256(tokenId))
}

func (e *ERC1155Batch) Uri(contractAddress string, tokenId *big.Int) *Result[string] {
	if err := validate.Address(contractAddress); err != nil {
		return failed[string](err)
	}
	if err := validate.Uint256(tokenId); err != nil {
		return failed[string](err)
	}
	return AddCall(e.b, contractAddress, constant.UriFnSignature, decode.ABIString, encode.ABIUint256(tokenId))
}
