package namespace

import (
	"math/big"

	"github.com/poteto-go/go-alchemy-sdk/types"
)

type IDebug interface {
	/*
		Snapshot takes a snapshot of the current blockchain state with evm_snapshot
		and returns the snapshot id.

		Only supported on simulated backend or development chains (hardhat, anvil, ganache, ...).
	*/
	Snapshot() (*big.Int, error)

	/*
		RevertTo reverts the blockchain state to the provided snapshot id with evm_revert.
		It returns true if the state was reverted.

		Only supported on simulated backend or development chains (hardhat, anvil, ganache, ...).
	*/
	RevertTo(snapshotId *big.Int) (bool, error)
}

type Debug struct {
	ether types.EtherApi
}

func NewDebugNamespace(ether types.EtherApi) IDebug {
	return &Debug{
		ether: ether,
	}
}

func (d *Debug) Snapshot() (*big.Int, error) {
	snapshotId, err := d.ether.Snapshot()
	if err != nil {
		return nil, err
	}

	return snapshotId, nil
}

func (d *Debug) RevertTo(snapshotId *big.Int) (bool, error) {
	reverted, err := d.ether.RevertTo(snapshotId)
	if err != nil {
		return false, err
	}

	return reverted, nil
}
