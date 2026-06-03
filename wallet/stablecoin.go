package wallet

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/poteto-go/go-alchemy-sdk/constant"
)

type WalletStableCoin interface {
	WalletERC20

	/*
		pause all token transfers (requires pauser role)
			- wait for mined
			- gas limit is 300000 for default
			- stops waiting when ctx is canceled
	*/
	Pause(ctx context.Context, contractAddress string, gasLimit *uint64) (*gethTypes.Receipt, error)

	/*
		pause all token transfers (requires pauser role)
			- gas limit is 300000 for default
	*/
	PauseNoWait(contractAddress string, gasLimit *uint64) (common.Hash, error)

	/*
		resume token transfers (requires pauser role)
			- wait for mined
			- gas limit is 300000 for default
			- stops waiting when ctx is canceled
	*/
	Unpause(ctx context.Context, contractAddress string, gasLimit *uint64) (*gethTypes.Receipt, error)

	/*
		resume token transfers (requires pauser role)
			- gas limit is 300000 for default
	*/
	UnpauseNoWait(contractAddress string, gasLimit *uint64) (common.Hash, error)

	/*
		check if the contract is currently paused
	*/
	Paused(contractAddress string) (bool, error)

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

	/*
		blacklist an address (requires blacklister role)
			- wait for mined
			- gas limit is 300000 for default
			- stops waiting when ctx is canceled
	*/
	Blacklist(ctx context.Context, contractAddress, address string, gasLimit *uint64) (*gethTypes.Receipt, error)

	/*
		blacklist an address (requires blacklister role)
			- gas limit is 300000 for default
	*/
	BlacklistNoWait(contractAddress, address string, gasLimit *uint64) (common.Hash, error)

	/*
		remove an address from the blacklist (requires blacklister role)
			- wait for mined
			- gas limit is 300000 for default
			- stops waiting when ctx is canceled
	*/
	UnBlacklist(ctx context.Context, contractAddress, address string, gasLimit *uint64) (*gethTypes.Receipt, error)

	/*
		remove an address from the blacklist (requires blacklister role)
			- gas limit is 300000 for default
	*/
	UnBlacklistNoWait(contractAddress, address string, gasLimit *uint64) (common.Hash, error)

	/*
		check if an address is blacklisted
	*/
	IsBlacklisted(contractAddress, address string) (bool, error)

	/*
		return the master minter address of the contract
	*/
	MasterMinter(contractAddress string) (common.Address, error)

	/*
		return the pauser address of the contract
	*/
	Pauser(contractAddress string) (common.Address, error)

	/*
		return the blacklister address of the contract
	*/
	Blacklister(contractAddress string) (common.Address, error)

	/*
		return the current owner address of the contract
	*/
	Owner(contractAddress string) (common.Address, error)

	/*
		transfer contract ownership to a new address (requires owner role)
			- wait for mined
			- gas limit is 300000 for default
			- stops waiting when ctx is canceled
	*/
	TransferOwnership(ctx context.Context, contractAddress, newOwner string, gasLimit *uint64) (*gethTypes.Receipt, error)

	/*
		transfer contract ownership to a new address (requires owner role)
			- gas limit is 300000 for default
	*/
	TransferOwnershipNoWait(contractAddress, newOwner string, gasLimit *uint64) (common.Hash, error)

	/*
		get the currency identifier (e.g. "USD")
	*/
	Currency(contractAddress string) (string, error)

	/*
		get the contract version string
	*/
	Version(contractAddress string) (string, error)

	/*
		configure a minter with an allowance (requires masterMinter role)
			- wait for mined
			- gas limit is 300000 for default
			- stops waiting when ctx is canceled
	*/
	ConfigureMinter(ctx context.Context, contractAddress, minter string, allowance *big.Int, gasLimit *uint64) (*gethTypes.Receipt, error)

	/*
		configure a minter with an allowance (requires masterMinter role)
			- gas limit is 300000 for default
	*/
	ConfigureMinterNoWait(contractAddress, minter string, allowance *big.Int, gasLimit *uint64) (common.Hash, error)

	/*
		remove a minter (requires masterMinter role)
			- wait for mined
			- gas limit is 300000 for default
			- stops waiting when ctx is canceled
	*/
	RemoveMinter(ctx context.Context, contractAddress, minter string, gasLimit *uint64) (*gethTypes.Receipt, error)

	/*
		remove a minter (requires masterMinter role)
			- gas limit is 300000 for default
	*/
	RemoveMinterNoWait(contractAddress, minter string, gasLimit *uint64) (common.Hash, error)

	/*
		check if an address is a configured minter
	*/
	IsMinter(contractAddress, address string) (bool, error)

	/*
		get the remaining mint allowance for a minter
	*/
	MinterAllowance(contractAddress, address string) (*big.Int, error)
}

type walletStableCoin struct {
	walletERC20
}

func (api *walletStableCoin) MintNoWait(contractAddress, toAddress string, amount *big.Int, gasLimit *uint64) (common.Hash, error) {
	return api.sendERC20Tx(contractAddress, gasLimit, constant.MintFnSignature,
		common.LeftPadBytes(common.HexToAddress(toAddress).Bytes(), constant.ABIWordSize),
		common.LeftPadBytes(amount.Bytes(), constant.ABIWordSize),
	)
}

func (api *walletStableCoin) Mint(ctx context.Context, contractAddress, toAddress string, amount *big.Int, gasLimit *uint64) (*gethTypes.Receipt, error) {
	return api.waitMined(ctx, func() (common.Hash, error) {
		return api.MintNoWait(contractAddress, toAddress, amount, gasLimit)
	})
}

func (api *walletStableCoin) BurnNoWait(contractAddress string, amount *big.Int, gasLimit *uint64) (common.Hash, error) {
	return api.sendERC20Tx(contractAddress, gasLimit, constant.BurnFnSignature,
		common.LeftPadBytes(amount.Bytes(), constant.ABIWordSize),
	)
}

func (api *walletStableCoin) Burn(ctx context.Context, contractAddress string, amount *big.Int, gasLimit *uint64) (*gethTypes.Receipt, error) {
	return api.waitMined(ctx, func() (common.Hash, error) {
		return api.BurnNoWait(contractAddress, amount, gasLimit)
	})
}

func (api *walletStableCoin) BlacklistNoWait(contractAddress, address string, gasLimit *uint64) (common.Hash, error) {
	return api.sendERC20Tx(contractAddress, gasLimit, constant.BlacklistFnSignature,
		common.LeftPadBytes(common.HexToAddress(address).Bytes(), constant.ABIWordSize),
	)
}

func (api *walletStableCoin) Blacklist(ctx context.Context, contractAddress, address string, gasLimit *uint64) (*gethTypes.Receipt, error) {
	return api.waitMined(ctx, func() (common.Hash, error) {
		return api.BlacklistNoWait(contractAddress, address, gasLimit)
	})
}

func (api *walletStableCoin) UnBlacklistNoWait(contractAddress, address string, gasLimit *uint64) (common.Hash, error) {
	return api.sendERC20Tx(contractAddress, gasLimit, constant.UnBlacklistFnSignature,
		common.LeftPadBytes(common.HexToAddress(address).Bytes(), constant.ABIWordSize),
	)
}

func (api *walletStableCoin) UnBlacklist(ctx context.Context, contractAddress, address string, gasLimit *uint64) (*gethTypes.Receipt, error) {
	return api.waitMined(ctx, func() (common.Hash, error) {
		return api.UnBlacklistNoWait(contractAddress, address, gasLimit)
	})
}

func (api *walletStableCoin) IsBlacklisted(contractAddress, address string) (bool, error) {
	sc := api.w.snapshotStableCoin()
	if sc == nil {
		return false, constant.ErrWalletIsNotConnected
	}
	return sc.IsBlacklisted(contractAddress, address)
}

func (api *walletStableCoin) MasterMinter(contractAddress string) (common.Address, error) {
	sc := api.w.snapshotStableCoin()
	if sc == nil {
		return common.Address{}, constant.ErrWalletIsNotConnected
	}
	return sc.MasterMinter(contractAddress)
}

func (api *walletStableCoin) Pauser(contractAddress string) (common.Address, error) {
	sc := api.w.snapshotStableCoin()
	if sc == nil {
		return common.Address{}, constant.ErrWalletIsNotConnected
	}
	return sc.Pauser(contractAddress)
}

func (api *walletStableCoin) Blacklister(contractAddress string) (common.Address, error) {
	sc := api.w.snapshotStableCoin()
	if sc == nil {
		return common.Address{}, constant.ErrWalletIsNotConnected
	}
	return sc.Blacklister(contractAddress)
}

func (api *walletStableCoin) Owner(contractAddress string) (common.Address, error) {
	sc := api.w.snapshotStableCoin()
	if sc == nil {
		return common.Address{}, constant.ErrWalletIsNotConnected
	}
	return sc.Owner(contractAddress)
}

func (api *walletStableCoin) PauseNoWait(contractAddress string, gasLimit *uint64) (common.Hash, error) {
	return api.sendERC20Tx(contractAddress, gasLimit, constant.PauseFnSignature)
}

func (api *walletStableCoin) Pause(ctx context.Context, contractAddress string, gasLimit *uint64) (*gethTypes.Receipt, error) {
	return api.waitMined(ctx, func() (common.Hash, error) {
		return api.PauseNoWait(contractAddress, gasLimit)
	})
}

func (api *walletStableCoin) UnpauseNoWait(contractAddress string, gasLimit *uint64) (common.Hash, error) {
	return api.sendERC20Tx(contractAddress, gasLimit, constant.UnpauseFnSignature)
}

func (api *walletStableCoin) Unpause(ctx context.Context, contractAddress string, gasLimit *uint64) (*gethTypes.Receipt, error) {
	return api.waitMined(ctx, func() (common.Hash, error) {
		return api.UnpauseNoWait(contractAddress, gasLimit)
	})
}

func (api *walletStableCoin) Paused(contractAddress string) (bool, error) {
	sc := api.w.snapshotStableCoin()
	if sc == nil {
		return false, constant.ErrWalletIsNotConnected
	}
	return sc.Paused(contractAddress)
}

func (api *walletStableCoin) TransferOwnershipNoWait(contractAddress, newOwner string, gasLimit *uint64) (common.Hash, error) {
	return api.sendERC20Tx(contractAddress, gasLimit, constant.TransferOwnershipFnSignature,
		common.LeftPadBytes(common.HexToAddress(newOwner).Bytes(), constant.ABIWordSize),
	)
}

func (api *walletStableCoin) TransferOwnership(ctx context.Context, contractAddress, newOwner string, gasLimit *uint64) (*gethTypes.Receipt, error) {
	return api.waitMined(ctx, func() (common.Hash, error) {
		return api.TransferOwnershipNoWait(contractAddress, newOwner, gasLimit)
	})
}

func (api *walletStableCoin) Currency(contractAddress string) (string, error) {
	sc := api.w.snapshotStableCoin()
	if sc == nil {
		return "", constant.ErrWalletIsNotConnected
	}
	return sc.Currency(contractAddress)
}

func (api *walletStableCoin) Version(contractAddress string) (string, error) {
	sc := api.w.snapshotStableCoin()
	if sc == nil {
		return "", constant.ErrWalletIsNotConnected
	}
	return sc.Version(contractAddress)
}

func (api *walletStableCoin) ConfigureMinterNoWait(contractAddress, minter string, allowance *big.Int, gasLimit *uint64) (common.Hash, error) {
	return api.sendERC20Tx(contractAddress, gasLimit, constant.ConfigureMinterFnSignature,
		common.LeftPadBytes(common.HexToAddress(minter).Bytes(), constant.ABIWordSize),
		common.LeftPadBytes(allowance.Bytes(), constant.ABIWordSize),
	)
}

func (api *walletStableCoin) ConfigureMinter(ctx context.Context, contractAddress, minter string, allowance *big.Int, gasLimit *uint64) (*gethTypes.Receipt, error) {
	return api.waitMined(ctx, func() (common.Hash, error) {
		return api.ConfigureMinterNoWait(contractAddress, minter, allowance, gasLimit)
	})
}

func (api *walletStableCoin) RemoveMinterNoWait(contractAddress, minter string, gasLimit *uint64) (common.Hash, error) {
	return api.sendERC20Tx(contractAddress, gasLimit, constant.RemoveMinterFnSignature,
		common.LeftPadBytes(common.HexToAddress(minter).Bytes(), constant.ABIWordSize),
	)
}

func (api *walletStableCoin) RemoveMinter(ctx context.Context, contractAddress, minter string, gasLimit *uint64) (*gethTypes.Receipt, error) {
	return api.waitMined(ctx, func() (common.Hash, error) {
		return api.RemoveMinterNoWait(contractAddress, minter, gasLimit)
	})
}

func (api *walletStableCoin) IsMinter(contractAddress, address string) (bool, error) {
	sc := api.w.snapshotStableCoin()
	if sc == nil {
		return false, constant.ErrWalletIsNotConnected
	}
	return sc.IsMinter(contractAddress, address)
}

func (api *walletStableCoin) MinterAllowance(contractAddress, address string) (*big.Int, error) {
	sc := api.w.snapshotStableCoin()
	if sc == nil {
		return nil, constant.ErrWalletIsNotConnected
	}
	return sc.MinterAllowance(contractAddress, address)
}
