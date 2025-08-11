package utils

import (
	"github.com/poteto-go/go-alchemy-sdk/constant"
	"github.com/poteto-go/go-alchemy-sdk/types"
)

func TransformBlock(rawTx types.BlockResponse) (types.Block, error) {
	blockNumber, err := FromHex(rawTx.Number)
	if err != nil {
		return types.Block{}, constant.ErrFailedToTransformBlockNumber
	}

	timestamp, err := FromHex(rawTx.Timestamp)
	if err != nil {
		return types.Block{}, constant.ErrFailedToTransformBlockNumber
	}

	difficulty, err := FromHex(rawTx.Difficulty)
	if err != nil {
		return types.Block{}, constant.ErrFailedToTransformDifficulty
	}

	gasLimit, err := FromBigHex(rawTx.GasLimit)
	if err != nil {
		return types.Block{}, constant.ErrFailedToTransformGasLimit
	}

	gasUsed, err := FromBigHex(rawTx.GasUsed)
	if err != nil {
		return types.Block{}, constant.ErrFailedToTransformGasLimit
	}

	return types.Block{
		Hash:         rawTx.Hash,
		ParentHash:   rawTx.ParentHash,
		Number:       blockNumber,
		Timestamp:    timestamp,
		Nonce:        rawTx.Nonce,
		Difficulty:   difficulty,
		GasLimit:     gasLimit,
		GasUsed:      gasUsed,
		Miner:        rawTx.Miner,
		Transactions: rawTx.Transactions,
	}, nil
}
