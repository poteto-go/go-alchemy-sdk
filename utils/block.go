package utils

import (
	"github.com/poteto-go/go-alchemy-sdk/core"
	"github.com/poteto-go/go-alchemy-sdk/types"
)

func TransformBlock(rawTx types.BlockResponse) (types.Block, error) {
	blockNumber, err := FromHex(rawTx.Number)
	if err != nil {
		return types.Block{}, core.ErrFailedToTransformBlockNumber
	}

	timestamp, err := FromHex(rawTx.Timestamp)
	if err != nil {
		return types.Block{}, core.ErrFailedToTransformBlockNumber
	}

	difficulty, err := FromHex(rawTx.Difficulty)
	if err != nil {
		return types.Block{}, core.ErrFailedToTransformDifficulty
	}

	gasLimit, err := FromBigHex(rawTx.GasLimit)
	if err != nil {
		return types.Block{}, core.ErrFailedToTransformGasLimit
	}

	gasUsed, err := FromBigHex(rawTx.GasUsed)
	if err != nil {
		return types.Block{}, core.ErrFailedToTransformGasLimit
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
