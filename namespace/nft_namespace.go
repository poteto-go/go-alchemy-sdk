package namespace

import (
	"math/big"
	"strings"

	"github.com/poteto-go/go-alchemy-sdk/constant"
	"github.com/poteto-go/go-alchemy-sdk/decode"
	"github.com/poteto-go/go-alchemy-sdk/encode"
	"github.com/poteto-go/go-alchemy-sdk/types"
	"github.com/poteto-go/go-alchemy-sdk/validate"
)

type INft interface {
	// OwnerOf returns the owner of the NFT with the given tokenId.
	OwnerOf(contractAddress string, tokenId *big.Int) (string, error)

	// TokenURI returns the URI of the NFT metadata for the given tokenId.
	TokenURI(contractAddress string, tokenId *big.Int) (string, error)

	// Name returns the name of the NFT collection.
	Name(contractAddress string) (string, error)

	// Symbol returns the symbol of the NFT collection.
	Symbol(contractAddress string) (string, error)

	// GetApproved returns the approved address for the given tokenId.
	GetApproved(contractAddress string, tokenId *big.Int) (string, error)

	// IsApprovedForAll returns whether the operator is approved to manage all
	// of the owner's NFTs.
	IsApprovedForAll(contractAddress, owner, operator string) (bool, error)
}

type Nft struct {
	ether types.EtherApi
}

func NewNftNamespace(ether types.EtherApi) INft {
	return &Nft{
		ether: ether,
	}
}

func (n *Nft) OwnerOf(contractAddress string, tokenId *big.Int) (string, error) {
	if err := validate.Address(contractAddress); err != nil {
		return "", err
	}
	if err := validate.Uint256(tokenId); err != nil {
		return "", err
	}
	output, err := n.ether.CallReadMethod(
		constant.OwnerOfFnSignature,
		contractAddress,
		encode.ABIUint256(tokenId),
	)
	if err != nil {
		return "", err
	}

	addr, err := decode.ABIAddress(output)
	if err != nil {
		return "", err
	}
	return strings.ToLower(addr.Hex()), nil
}

func (n *Nft) TokenURI(contractAddress string, tokenId *big.Int) (string, error) {
	if err := validate.Address(contractAddress); err != nil {
		return "", err
	}
	if err := validate.Uint256(tokenId); err != nil {
		return "", err
	}
	output, err := n.ether.CallReadMethod(
		constant.TokenURIFnSignature,
		contractAddress,
		encode.ABIUint256(tokenId),
	)
	if err != nil {
		return "", err
	}

	return decode.ABIString(output)
}

func (n *Nft) Name(contractAddress string) (string, error) {
	if err := validate.Address(contractAddress); err != nil {
		return "", err
	}
	output, err := n.ether.CallReadMethod(
		constant.NameFnSignature,
		contractAddress,
	)
	if err != nil {
		return "", err
	}

	return decode.ABIString(output)
}

func (n *Nft) Symbol(contractAddress string) (string, error) {
	if err := validate.Address(contractAddress); err != nil {
		return "", err
	}
	output, err := n.ether.CallReadMethod(
		constant.SymbolFnSignature,
		contractAddress,
	)
	if err != nil {
		return "", err
	}

	return decode.ABIString(output)
}

func (n *Nft) GetApproved(contractAddress string, tokenId *big.Int) (string, error) {
	if err := validate.Address(contractAddress); err != nil {
		return "", err
	}
	if err := validate.Uint256(tokenId); err != nil {
		return "", err
	}
	output, err := n.ether.CallReadMethod(
		constant.GetApprovedFnSignature,
		contractAddress,
		encode.ABIUint256(tokenId),
	)
	if err != nil {
		return "", err
	}

	addr, err := decode.ABIAddress(output)
	if err != nil {
		return "", err
	}
	return strings.ToLower(addr.Hex()), nil
}

func (n *Nft) IsApprovedForAll(contractAddress, owner, operator string) (bool, error) {
	if err := validate.Addresses(contractAddress, owner, operator); err != nil {
		return false, err
	}
	output, err := n.ether.CallReadMethod(
		constant.IsApprovedForAllFnSignature,
		contractAddress,
		encode.ABIAddress(owner),
		encode.ABIAddress(operator),
	)
	if err != nil {
		return false, err
	}

	return decode.Bool(output)
}
