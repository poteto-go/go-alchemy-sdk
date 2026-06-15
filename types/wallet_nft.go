package types

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
)

// walletNftApproval is the write-side approval interface shared by ERC-721 and ERC-1155.
type walletNftApproval interface {
	/*
		grant or revoke an operator's approval to manage all of the caller's tokens
			- wait for mined
			- gas limit is 300000 for default
			- stops waiting when ctx is canceled
	*/
	SetApprovalForAll(ctx context.Context, contractAddress, operator string, approved bool, gasLimit *uint64) (*gethTypes.Receipt, error)

	/*
		grant or revoke an operator's approval to manage all of the caller's tokens
			- gas limit is 300000 for default
	*/
	SetApprovalForAllNoWait(contractAddress, operator string, approved bool, gasLimit *uint64) (common.Hash, error)

	/*
		get whether the operator is approved to manage all of the owner's tokens
	*/
	IsApprovedForAll(contractAddress, owner, operator string) (bool, error)
}

// Nft (ERC721) interface for wallet.
// This is only defined for UX.
type WalletNft interface {
	walletNftApproval
	/*
		transfer the NFT with the given tokenId from one address to another
			- wait for mined
			- gas limit is 300000 for default
			- stops waiting when ctx is canceled
	*/
	TransferFrom(ctx context.Context, contractAddress, fromAddress, toAddress string, tokenId *big.Int, gasLimit *uint64) (*gethTypes.Receipt, error)

	/*
		transfer the NFT with the given tokenId from one address to another
			- gas limit is 300000 for default
	*/
	TransferFromNoWait(contractAddress, fromAddress, toAddress string, tokenId *big.Int, gasLimit *uint64) (common.Hash, error)

	/*
		safely transfer the NFT with the given tokenId from one address to another
			- wait for mined
			- gas limit is 300000 for default
			- stops waiting when ctx is canceled
	*/
	SafeTransferFrom(ctx context.Context, contractAddress, fromAddress, toAddress string, tokenId *big.Int, gasLimit *uint64) (*gethTypes.Receipt, error)

	/*
		safely transfer the NFT with the given tokenId from one address to another
			- gas limit is 300000 for default
	*/
	SafeTransferFromNoWait(contractAddress, fromAddress, toAddress string, tokenId *big.Int, gasLimit *uint64) (common.Hash, error)

	/*
		safely transfer the NFT with the given tokenId, passing additional data
		to the recipient's onERC721Received hook
			- wait for mined
			- gas limit is 300000 for default
			- stops waiting when ctx is canceled
	*/
	SafeTransferFromWithData(ctx context.Context, contractAddress, fromAddress, toAddress string, tokenId *big.Int, data []byte, gasLimit *uint64) (*gethTypes.Receipt, error)

	/*
		safely transfer the NFT with the given tokenId, passing additional data
		to the recipient's onERC721Received hook
			- gas limit is 300000 for default
	*/
	SafeTransferFromWithDataNoWait(contractAddress, fromAddress, toAddress string, tokenId *big.Int, data []byte, gasLimit *uint64) (common.Hash, error)

	/*
		approve another address to transfer the NFT with the given tokenId
			- wait for mined
			- gas limit is 300000 for default
			- stops waiting when ctx is canceled
	*/
	Approve(ctx context.Context, contractAddress, toAddress string, tokenId *big.Int, gasLimit *uint64) (*gethTypes.Receipt, error)

	/*
		approve another address to transfer the NFT with the given tokenId
			- gas limit is 300000 for default
	*/
	ApproveNoWait(contractAddress, toAddress string, tokenId *big.Int, gasLimit *uint64) (common.Hash, error)

	/*
		get the number of NFTs owned by the given address
	*/
	BalanceOf(contractAddress, owner string) (*big.Int, error)

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
}
