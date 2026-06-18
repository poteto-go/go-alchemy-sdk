package namespace

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/poteto-go/go-alchemy-sdk/constant"
	"github.com/poteto-go/go-alchemy-sdk/types"
)

var bigOne = big.NewInt(1)

type SimulatedDebug struct {
	ether types.EtherApi

	// Debug compatible
	snapShotCount    *big.Int
	snapShotRegistry map[*big.Int]common.Hash
}

func NewSimulatedDebugNamespace(ether types.EtherApi) IDebug {
	return &SimulatedDebug{
		ether:            ether,
		snapShotCount:    big.NewInt(0),
		snapShotRegistry: make(map[*big.Int]common.Hash),
	}
}

func (sd *SimulatedDebug) Snapshot() (*big.Int, error) {
	block, err := sd.ether.GetBlockByNumber("latest")
	if err != nil {
		return nil, err
	}

	snapShotId := sd.snapShotCount
	sd.snapShotRegistry[snapShotId] = block.Hash()

	newCount := new(big.Int)
	newCount.Add(sd.snapShotCount, bigOne)
	sd.snapShotCount = newCount
	return snapShotId, nil
}

func (sd *SimulatedDebug) RevertTo(snapShotId *big.Int) (bool, error) {
	snapShotHash, ok := sd.snapShotRegistry[snapShotId]
	if !ok {
		return false, constant.ErrUnexpectedSnapshotId
	}

	if err := sd.ether.Fork(snapShotHash); err != nil {
		return false, err
	}

	newCount := new(big.Int)
	newCount.Sub(sd.snapShotCount, bigOne)
	sd.snapShotCount = newCount
	return true, nil
}
