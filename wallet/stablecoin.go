package wallet

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/poteto-go/go-alchemy-sdk/constant"
	"github.com/poteto-go/go-alchemy-sdk/types"
)

type WalletStableCoin interface {
	WalletERC20

	/*
		mint stablecoin tokens to an address (requires minter role)
			- wait for mined
			- gas limit is 300000 for default
			- stops waiting when ctx is canceled
	*/
	Mint(ctx context.Context, contractAddress, toAddress string, amount *big.Int, gasLimit *uint64) (*gethTypes.Receipt, error)

	/*
		mint stablecoin tokens to an address (requires minter role)
			- gas limit is 300000 for default
	*/
	MintNoWait(contractAddress, toAddress string, amount *big.Int, gasLimit *uint64) (common.Hash, error)

	/*
		burn stablecoin tokens from the caller's balance (requires minter role)
			- wait for mined
			- gas limit is 300000 for default
			- stops waiting when ctx is canceled
	*/
	Burn(ctx context.Context, contractAddress string, amount *big.Int, gasLimit *uint64) (*gethTypes.Receipt, error)

	/*
		burn stablecoin tokens from the caller's balance (requires minter role)
			- gas limit is 300000 for default
	*/
	BurnNoWait(contractAddress string, amount *big.Int, gasLimit *uint64) (common.Hash, error)
}

type walletStableCoin struct {
	walletERC20
}

func (api *walletStableCoin) MintNoWait(contractAddress, toAddress string, amount *big.Int, gasLimit *uint64) (common.Hash, error) {
	provider := api.w.snapshot()
	if provider == nil {
		return common.Hash{}, constant.ErrWalletIsNotConnected
	}

	data, err := buildERC20TxData(
		constant.MintFnSignature,
		common.LeftPadBytes(common.HexToAddress(toAddress).Bytes(), 32),
		common.LeftPadBytes(amount.Bytes(), 32),
	)
	if err != nil {
		return common.Hash{}, err
	}

	txRequest := types.TransactionRequest{
		From:     api.w.GetAddress(),
		To:       contractAddress,
		Value:    "0x0",
		GasLimit: resolveGasLimit(gasLimit),
		Data:     data,
	}

	return api.w.SendTransaction(txRequest)
}

func (api *walletStableCoin) Mint(ctx context.Context, contractAddress, toAddress string, amount *big.Int, gasLimit *uint64) (*gethTypes.Receipt, error) {
	return api.waitMined(ctx, func() (common.Hash, error) {
		return api.MintNoWait(contractAddress, toAddress, amount, gasLimit)
	})
}

func (api *walletStableCoin) BurnNoWait(contractAddress string, amount *big.Int, gasLimit *uint64) (common.Hash, error) {
	provider := api.w.snapshot()
	if provider == nil {
		return common.Hash{}, constant.ErrWalletIsNotConnected
	}

	data, err := buildERC20TxData(
		constant.BurnFnSignature,
		common.LeftPadBytes(amount.Bytes(), 32),
	)
	if err != nil {
		return common.Hash{}, err
	}

	txRequest := types.TransactionRequest{
		From:     api.w.GetAddress(),
		To:       contractAddress,
		Value:    "0x0",
		GasLimit: resolveGasLimit(gasLimit),
		Data:     data,
	}

	return api.w.SendTransaction(txRequest)
}

func (api *walletStableCoin) Burn(ctx context.Context, contractAddress string, amount *big.Int, gasLimit *uint64) (*gethTypes.Receipt, error) {
	return api.waitMined(ctx, func() (common.Hash, error) {
		return api.BurnNoWait(contractAddress, amount, gasLimit)
	})
}
