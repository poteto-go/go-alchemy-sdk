package types

import (
	"math/big"
)

// Nft (ERC721) interface for wallet.
// This is only defined for UX.
type WalletNft interface {
	/*
		get owner of the NFT with the given tokenId
	*/
	OwnerOf(contractAddress string, tokenId *big.Int) (string, error)

	/*
		get URI of the NFT metadata for the given tokenId
	*/
	TokenURI(contractAddress string, tokenId *big.Int) (string, error)

	/*
		get name of the NFT collection
	*/
	Name(contractAddress string) (string, error)

	/*
		get symbol of the NFT collection
	*/
	Symbol(contractAddress string) (string, error)

	/*
		get approved address for the given tokenId
	*/
	GetApproved(contractAddress string, tokenId *big.Int) (string, error)

	/*
		get whether the operator is approved to manage all of the owner's NFTs
	*/
	IsApprovedForAll(contractAddress, owner, operator string) (bool, error)
}
