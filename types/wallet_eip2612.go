package types

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
)

type WalletEIP2612 interface {
	// PermitNoWait submits a permit transaction using a pre-signed signature.
	PermitNoWait(contractAddress, ownerAddress, spenderAddress string, value, deadline *big.Int, v uint8, r, s [32]byte, gasLimit *uint64) (common.Hash, error)

	// Permit submits a permit transaction and waits for it to be mined.
	Permit(ctx context.Context, contractAddress, ownerAddress, spenderAddress string, value, deadline *big.Int, v uint8, r, s [32]byte, gasLimit *uint64) (*gethTypes.Receipt, error)

	// Nonces returns the current nonce for the given owner address.
	Nonces(contractAddress, ownerAddress string) (*big.Int, error)

	// DomainSeparator returns the EIP-712 domain separator for the contract.
	DomainSeparator(contractAddress string) ([32]byte, error)
}
