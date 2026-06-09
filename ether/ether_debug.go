package ether

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/poteto-go/go-alchemy-sdk/constant"
	"github.com/poteto-go/go-alchemy-sdk/types"
	"github.com/poteto-go/go-alchemy-sdk/utils"
)

func (ether *Ether) Snapshot() (*big.Int, error) {
	result, err := ether.provider.Send(constant.Evm_Snapshot, types.RequestArgs{})
	if err != nil {
		return nil, err
	}

	resultStr, ok := result.(string)
	if !ok {
		return nil, constant.ErrUnexpectedResponseType
	}

	snapshotId, err := utils.FromBigHex(resultStr)
	if err != nil {
		return nil, err
	}

	return snapshotId, nil
}

func (ether *Ether) RevertTo(snapshotId *big.Int) (bool, error) {
	if snapshotId == nil {
		return false, constant.ErrNilSnapshotId
	}

	result, err := ether.provider.Send(constant.Evm_Revert, types.RequestArgs{
		hexutil.EncodeBig(snapshotId),
	})
	if err != nil {
		return false, err
	}

	reverted, ok := result.(bool)
	if !ok {
		return false, constant.ErrUnexpectedResponseType
	}

	return reverted, nil
}
