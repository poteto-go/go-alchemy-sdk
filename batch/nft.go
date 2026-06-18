package batch

import (
	"fmt"
	"math/big"

	"github.com/poteto-go/go-alchemy-sdk/constant"
	"github.com/poteto-go/go-alchemy-sdk/decode"
	"github.com/poteto-go/go-alchemy-sdk/encode"
	"github.com/poteto-go/go-alchemy-sdk/validate"
)

// NftBatch queues ERC-721 contract reads onto its Batcher. Mirrors the
// read-only methods of the Nft namespace and validates the same inputs.
type NftBatch struct {
	b *Batcher
}

// decodeAddressHex decodes an ABI address and returns it as a lowercase hex string.
func decodeAddressHex(output []byte) (string, error) {
	addr, err := decode.ABIAddress(output)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("0x%x", addr[:]), nil
}

// addAddressTokenCall validates address + tokenId, then queues an eth_call
// that decodes its result as a lowercase hex address. Used by OwnerOf and GetApproved.
func addAddressTokenCall(b *Batcher, fnSig []byte, contractAddress string, tokenId *big.Int) *Result[string] {
	if err := validate.Address(contractAddress); err != nil {
		return failed[string](err)
	}
	if err := validate.Uint256(tokenId); err != nil {
		return failed[string](err)
	}
	return AddCall(b, contractAddress, fnSig, decodeAddressHex, encode.ABIUint256(tokenId))
}

func (n *NftBatch) BalanceOf(contractAddress, owner string) *Result[*big.Int] {
	if err := validate.Addresses(contractAddress, owner); err != nil {
		return failed[*big.Int](err)
	}
	return AddCall(n.b, contractAddress, constant.BalanceOfFnSignature, decode.Uint256, encode.ABIAddress(owner))
}

func (n *NftBatch) OwnerOf(contractAddress string, tokenId *big.Int) *Result[string] {
	return addAddressTokenCall(n.b, constant.OwnerOfFnSignature, contractAddress, tokenId)
}

func (n *NftBatch) TokenURI(contractAddress string, tokenId *big.Int) *Result[string] {
	if err := validate.Address(contractAddress); err != nil {
		return failed[string](err)
	}
	if err := validate.Uint256(tokenId); err != nil {
		return failed[string](err)
	}
	return AddCall(n.b, contractAddress, constant.TokenURIFnSignature, decode.ABIString, encode.ABIUint256(tokenId))
}

func (n *NftBatch) Name(contractAddress string) *Result[string] {
	return addStringCall(n.b, contractAddress, constant.NameFnSignature)
}

func (n *NftBatch) Symbol(contractAddress string) *Result[string] {
	return addStringCall(n.b, contractAddress, constant.SymbolFnSignature)
}

func (n *NftBatch) GetApproved(contractAddress string, tokenId *big.Int) *Result[string] {
	return addAddressTokenCall(n.b, constant.GetApprovedFnSignature, contractAddress, tokenId)
}
