package types

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
)

// WalletERC1155 (ERC1155 multi-token) interface for wallet.
type WalletERC1155 interface {
	walletNftApproval

	/*
		get the amount of tokens of the given tokenId owned by account
	*/
	BalanceOfToken(contractAddress, account string, tokenId *big.Int) (*big.Int, error)

	/*
		get the balances of multiple (account, tokenId) pairs in a single call.
		accounts and tokenIds must have the same length.
	*/
	BalanceOfBatch(contractAddress string, accounts []string, tokenIds []*big.Int) ([]*big.Int, error)

	/*
		get the metadata URI for the given tokenId
	*/
	Uri(contractAddress string, tokenId *big.Int) (string, error)

	/*
		ERC-1155 safeTransferFrom(from,to,id,amount,data) — transfers `amount` units of
		token `id` from `fromAddress` to `toAddress`. data is passed to onERC1155Received.
			- wait for mined
			- gas limit is 300000 for default
			- stops waiting when ctx is canceled
	*/
	SafeTransferFrom(ctx context.Context, contractAddress, fromAddress, toAddress string, id, amount *big.Int, data []byte, gasLimit *uint64) (*gethTypes.Receipt, error)

	/*
		ERC-1155 safeTransferFrom(from,to,id,amount,data) — no-wait variant.
			- gas limit is 300000 for default
	*/
	SafeTransferFromNoWait(contractAddress, fromAddress, toAddress string, id, amount *big.Int, data []byte, gasLimit *uint64) (common.Hash, error)

	/*
		ERC-1155 safeBatchTransferFrom(from,to,ids,amounts,data) — transfers multiple
		token types in a single call. ids and amounts must have the same length.
			- wait for mined
			- gas limit is 300000 for default
			- stops waiting when ctx is canceled
	*/
	SafeBatchTransferFrom(ctx context.Context, contractAddress, fromAddress, toAddress string, ids, amounts []*big.Int, data []byte, gasLimit *uint64) (*gethTypes.Receipt, error)

	/*
		ERC-1155 safeBatchTransferFrom(from,to,ids,amounts,data) — no-wait variant.
			- gas limit is 300000 for default
	*/
	SafeBatchTransferFromNoWait(contractAddress, fromAddress, toAddress string, ids, amounts []*big.Int, data []byte, gasLimit *uint64) (common.Hash, error)
}
