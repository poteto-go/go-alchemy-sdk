package namespace

import (
	"math/big"

	"github.com/poteto-go/go-alchemy-sdk/constant"
	"github.com/poteto-go/go-alchemy-sdk/decode"
	"github.com/poteto-go/go-alchemy-sdk/encode"
	"github.com/poteto-go/go-alchemy-sdk/types"
	"github.com/poteto-go/go-alchemy-sdk/validate"
)

type IErc1155 interface {
	iApprovedForAll

	// BalanceOfToken returns the amount of the given tokenId owned by account.
	BalanceOfToken(contractAddress, account string, tokenId *big.Int) (*big.Int, error)

	// BalanceOfBatch returns the balances of multiple (account, tokenId) pairs
	// in a single call. accounts and tokenIds must have the same length.
	BalanceOfBatch(contractAddress string, accounts []string, tokenIds []*big.Int) ([]*big.Int, error)

	// Uri returns the metadata URI for the given tokenId.
	Uri(contractAddress string, tokenId *big.Int) (string, error)
}

// Erc1155 embeds *Nft to reuse its IsApprovedForAll implementation.
type Erc1155 struct {
	*Nft
}

func NewErc1155Namespace(ether types.EtherApi) IErc1155 {
	return &Erc1155{
		Nft: &Nft{ether: ether},
	}
}

func (e *Erc1155) BalanceOfToken(contractAddress, account string, tokenId *big.Int) (*big.Int, error) {
	if err := validate.Addresses(contractAddress, account); err != nil {
		return nil, err
	}
	if err := validate.Uint256(tokenId); err != nil {
		return nil, err
	}
	output, err := e.ether.CallReadMethod(
		constant.BalanceOfTokenFnSignature,
		contractAddress,
		encode.ABIAddress(account),
		encode.ABIUint256(tokenId),
	)
	if err != nil {
		return nil, err
	}

	return decode.Uint256(output)
}

func (e *Erc1155) BalanceOfBatch(contractAddress string, accounts []string, tokenIds []*big.Int) ([]*big.Int, error) {
	if err := validate.Address(contractAddress); err != nil {
		return nil, err
	}
	if len(accounts) != len(tokenIds) {
		return nil, constant.ErrMismatchedArrayLength
	}
	if err := validate.Addresses(accounts...); err != nil {
		return nil, err
	}
	for _, tokenId := range tokenIds {
		if err := validate.Uint256(tokenId); err != nil {
			return nil, err
		}
	}

	// balanceOfBatch(address[],uint256[]) takes two dynamic arrays; ABIDynamicArgs
	// builds the offset head and appends each array tail.
	args := encode.ABIDynamicArgs(
		encode.ABIAddressArray(accounts),
		encode.ABIUint256Array(tokenIds),
	)
	output, err := e.ether.CallReadMethod(
		constant.BalanceOfBatchFnSignature,
		contractAddress,
		args,
	)
	if err != nil {
		return nil, err
	}

	balances, err := decode.Uint256Array(output)
	if err != nil {
		return nil, err
	}
	if len(balances) != len(accounts) {
		return nil, constant.ErrUnexpectedBalanceCount
	}
	return balances, nil
}

func (e *Erc1155) Uri(contractAddress string, tokenId *big.Int) (string, error) {
	if err := validate.Address(contractAddress); err != nil {
		return "", err
	}
	if err := validate.Uint256(tokenId); err != nil {
		return "", err
	}
	output, err := e.ether.CallReadMethod(
		constant.UriFnSignature,
		contractAddress,
		encode.ABIUint256(tokenId),
	)
	if err != nil {
		return "", err
	}

	return decode.ABIString(output)
}
