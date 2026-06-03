package wallet

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/poteto-go/go-alchemy-sdk/constant"
)

type StableCoinMinting interface {
	Mint(ctx context.Context, contractAddress, toAddress string, amount *big.Int, gasLimit *uint64) (*gethTypes.Receipt, error)
	MintNoWait(contractAddress, toAddress string, amount *big.Int, gasLimit *uint64) (common.Hash, error)
	Burn(ctx context.Context, contractAddress string, amount *big.Int, gasLimit *uint64) (*gethTypes.Receipt, error)
	BurnNoWait(contractAddress string, amount *big.Int, gasLimit *uint64) (common.Hash, error)
}

type StableCoinPausing interface {
	Pause(ctx context.Context, contractAddress string, gasLimit *uint64) (*gethTypes.Receipt, error)
	PauseNoWait(contractAddress string, gasLimit *uint64) (common.Hash, error)
	Unpause(ctx context.Context, contractAddress string, gasLimit *uint64) (*gethTypes.Receipt, error)
	UnpauseNoWait(contractAddress string, gasLimit *uint64) (common.Hash, error)
	Paused(contractAddress string) (bool, error)
}

type StableCoinBlacklisting interface {
	Blacklist(ctx context.Context, contractAddress, address string, gasLimit *uint64) (*gethTypes.Receipt, error)
	BlacklistNoWait(contractAddress, address string, gasLimit *uint64) (common.Hash, error)
	UnBlacklist(ctx context.Context, contractAddress, address string, gasLimit *uint64) (*gethTypes.Receipt, error)
	UnBlacklistNoWait(contractAddress, address string, gasLimit *uint64) (common.Hash, error)
	IsBlacklisted(contractAddress, address string) (bool, error)
}

type StableCoinMinterAdmin interface {
	ConfigureMinter(ctx context.Context, contractAddress, minter string, allowance *big.Int, gasLimit *uint64) (*gethTypes.Receipt, error)
	ConfigureMinterNoWait(contractAddress, minter string, allowance *big.Int, gasLimit *uint64) (common.Hash, error)
	RemoveMinter(ctx context.Context, contractAddress, minter string, gasLimit *uint64) (*gethTypes.Receipt, error)
	RemoveMinterNoWait(contractAddress, minter string, gasLimit *uint64) (common.Hash, error)
	IsMinter(contractAddress, address string) (bool, error)
	MinterAllowance(contractAddress, address string) (*big.Int, error)
}

type StableCoinRoleAdmin interface {
	UpdateMasterMinter(ctx context.Context, contractAddress, newMasterMinter string, gasLimit *uint64) (*gethTypes.Receipt, error)
	UpdateMasterMinterNoWait(contractAddress, newMasterMinter string, gasLimit *uint64) (common.Hash, error)
	UpdateBlacklister(ctx context.Context, contractAddress, newBlacklister string, gasLimit *uint64) (*gethTypes.Receipt, error)
	UpdateBlacklisterNoWait(contractAddress, newBlacklister string, gasLimit *uint64) (common.Hash, error)
	UpdatePauser(ctx context.Context, contractAddress, newPauser string, gasLimit *uint64) (*gethTypes.Receipt, error)
	UpdatePauserNoWait(contractAddress, newPauser string, gasLimit *uint64) (common.Hash, error)
	TransferOwnership(ctx context.Context, contractAddress, newOwner string, gasLimit *uint64) (*gethTypes.Receipt, error)
	TransferOwnershipNoWait(contractAddress, newOwner string, gasLimit *uint64) (common.Hash, error)
	Owner(contractAddress string) (common.Address, error)
	MasterMinter(contractAddress string) (common.Address, error)
	Pauser(contractAddress string) (common.Address, error)
	Blacklister(contractAddress string) (common.Address, error)
}

type StableCoinInfo interface {
	Currency(contractAddress string) (string, error)
	Version(contractAddress string) (string, error)
}

type WalletStableCoin interface {
	WalletERC20
	StableCoinMinting
	StableCoinPausing
	StableCoinBlacklisting
	StableCoinMinterAdmin
	StableCoinRoleAdmin
	StableCoinInfo
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

func (api *walletStableCoin) UpdateMasterMinterNoWait(contractAddress, newMasterMinter string, gasLimit *uint64) (common.Hash, error) {
	return api.sendERC20Tx(contractAddress, gasLimit, constant.UpdateMasterMinterFnSignature,
		common.LeftPadBytes(common.HexToAddress(newMasterMinter).Bytes(), constant.ABIWordSize),
	)
}

func (api *walletStableCoin) UpdateMasterMinter(ctx context.Context, contractAddress, newMasterMinter string, gasLimit *uint64) (*gethTypes.Receipt, error) {
	return api.waitMined(ctx, func() (common.Hash, error) {
		return api.UpdateMasterMinterNoWait(contractAddress, newMasterMinter, gasLimit)
	})
}

func (api *walletStableCoin) UpdateBlacklisterNoWait(contractAddress, newBlacklister string, gasLimit *uint64) (common.Hash, error) {
	return api.sendERC20Tx(contractAddress, gasLimit, constant.UpdateBlacklisterFnSignature,
		common.LeftPadBytes(common.HexToAddress(newBlacklister).Bytes(), constant.ABIWordSize),
	)
}

func (api *walletStableCoin) UpdateBlacklister(ctx context.Context, contractAddress, newBlacklister string, gasLimit *uint64) (*gethTypes.Receipt, error) {
	return api.waitMined(ctx, func() (common.Hash, error) {
		return api.UpdateBlacklisterNoWait(contractAddress, newBlacklister, gasLimit)
	})
}

func (api *walletStableCoin) UpdatePauserNoWait(contractAddress, newPauser string, gasLimit *uint64) (common.Hash, error) {
	return api.sendERC20Tx(contractAddress, gasLimit, constant.UpdatePauserFnSignature,
		common.LeftPadBytes(common.HexToAddress(newPauser).Bytes(), constant.ABIWordSize),
	)
}

func (api *walletStableCoin) UpdatePauser(ctx context.Context, contractAddress, newPauser string, gasLimit *uint64) (*gethTypes.Receipt, error) {
	return api.waitMined(ctx, func() (common.Hash, error) {
		return api.UpdatePauserNoWait(contractAddress, newPauser, gasLimit)
	})
}
