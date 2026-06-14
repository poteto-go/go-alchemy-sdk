package types

import "math/big"

// Erc1155 (ERC1155 multi-token) read interface for wallet.
// This is only defined for UX; it delegates to the Erc1155 namespace.
type WalletErc1155 interface {
	/*
		get the amount of tokens of the given tokenId owned by account
	*/
	BalanceOf(contractAddress, account string, tokenId *big.Int) (*big.Int, error)

	/*
		get the balances of multiple (account, tokenId) pairs in a single call.
		accounts and tokenIds must have the same length.
	*/
	BalanceOfBatch(contractAddress string, accounts []string, tokenIds []*big.Int) ([]*big.Int, error)

	/*
		get the metadata URI for the given tokenId
	*/
	Uri(contractAddress string, tokenId *big.Int) (string, error)
}
