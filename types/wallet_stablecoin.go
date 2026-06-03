package types

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
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

type StableCoinEIP2612 interface {
	PermitNoWait(contractAddress, ownerAddress, spenderAddress string, value, deadline *big.Int, gasLimit *uint64) (common.Hash, error)
	Permit(ctx context.Context, contractAddress, ownerAddress, spenderAddress string, value, deadline *big.Int, gasLimit *uint64) (*gethTypes.Receipt, error)
}

type WalletStableCoin interface {
	WalletERC20
	StableCoinMinting
	StableCoinPausing
	StableCoinBlacklisting
	StableCoinMinterAdmin
	StableCoinRoleAdmin
	StableCoinInfo
	StableCoinEIP2612
}
