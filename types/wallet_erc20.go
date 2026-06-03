package types

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
)

// ERC20 interface for wallet.
// This is only defined for UX.
type WalletERC20 interface {
	/*
		transfer erc20 token by provided wallet
			- wait for mined
			- gas limit is 300000 for default
			- stops waiting when ctx is canceled
	*/
	Transfer(ctx context.Context, contractAddress, toAddress string, amount *big.Int, gasLimit *uint64) (*gethTypes.Receipt, error)

	/*
		transfer erc20 token by provided wallet
			- gas limit is 300000 for default
	*/
	TransferNoWait(contractAddress, toAddress string, amount *big.Int, gasLimit *uint64) (common.Hash, error)

	/*
		transfer erc20 token from another address (requires prior approval)
			- wait for mined
			- gas limit is 300000 for default
			- stops waiting when ctx is canceled
	*/
	TransferFrom(ctx context.Context, contractAddress, fromAddress, toAddress string, amount *big.Int, gasLimit *uint64) (*gethTypes.Receipt, error)

	/*
		transfer erc20 token from another address (requires prior approval)
			- gas limit is 300000 for default
	*/
	TransferFromNoWait(contractAddress, fromAddress, toAddress string, amount *big.Int, gasLimit *uint64) (common.Hash, error)

	/*
		approve spender to spend erc20 token
			- wait for mined
			- gas limit is 300000 for default
			- stops waiting when ctx is canceled
	*/
	Approve(ctx context.Context, contractAddress, spenderAddress string, amount *big.Int, gasLimit *uint64) (*gethTypes.Receipt, error)

	/*
		approve spender to spend erc20 token
			- gas limit is 300000 for default
	*/
	ApproveNoWait(contractAddress, spenderAddress string, amount *big.Int, gasLimit *uint64) (common.Hash, error)

	/*
		get balance of provided wallet & erc20 token
	*/
	BalanceOf(contractAddress string) (*big.Int, error)

	/*
		get total supply of erc20 token
	*/
	TotalSupply(contractAddress string) (*big.Int, error)

	/*
		get allowance of erc20 token
	*/
	Allowance(contractAddress, owner, spender string) (*big.Int, error)

	/*
		get name of erc20 token
	*/
	Name(contractAddress string) (string, error)

	/*
		get symbol of erc20 token
	*/
	Symbol(contractAddress string) (string, error)

	/*
		get decimals of erc20 token
	*/
	Decimals(contractAddress string) (uint8, error)
}
