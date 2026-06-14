package types

import "math/big"

// WalletERC1155 (ERC1155 multi-token) interface for wallet.
// Embeds WalletNft to inherit ERC-721 compatible read/write methods.
type WalletERC1155 interface {
	WalletNft

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
}
