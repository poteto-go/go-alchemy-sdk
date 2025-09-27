package utils

import (
	"github.com/ethereum/go-ethereum/common"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/poteto-go/go-alchemy-sdk/types"
)

// from alchemy.Receipt -> geth.Receipt
func TransformAlchemyReceiptToGeth(receipt types.TransactionReceipt) (*gethTypes.Receipt, error) {
	logs := make([]*gethTypes.Log, len(receipt.Logs))
	for i, log := range receipt.Logs {
		gethLog, err := TransformAlchemyLogToGeth(log)
		if err != nil {
			return nil, err
		}

		logs[i] = gethLog
	}

	typeU64, err := FromHexU64(receipt.Type)
	if err != nil {
		return nil, err
	}

	status, err := FromHexU64(receipt.Status)
	if err != nil {
		return nil, err
	}

	cGasUsed, err := FromHexU64(receipt.CumulativeGasUsed)
	if err != nil {
		return nil, err
	}

	gasUsed, err := FromHexU64(receipt.GasUsed)
	if err != nil {
		return nil, err
	}

	eGasPrice, err := FromBigHex(receipt.EffectiveGasPrice)
	if err != nil {
		return nil, err
	}

	bGasUsed, err := FromHexU64(receipt.BlobGasUsed)
	if err != nil {
		return nil, err
	}

	blockNumber, err := FromBigHex(receipt.BlockNumber)
	if err != nil {
		return nil, err
	}

	txIndex, err := FromHexU64(receipt.TransactionIndex)
	if err != nil {
		return nil, err
	}

	return &gethTypes.Receipt{
		// nolint:gosec
		Type:              uint8(typeU64),
		PostState:         []byte(receipt.Root),
		Status:            status,
		CumulativeGasUsed: cGasUsed,
		Bloom:             gethTypes.Bloom([]byte(receipt.LogsBloom)),
		Logs:              logs,
		TxHash:            common.HexToHash(receipt.TransactionHash),
		ContractAddress:   common.HexToAddress(receipt.ContractAddress),
		GasUsed:           gasUsed,
		EffectiveGasPrice: eGasPrice,
		BlobGasUsed:       bGasUsed,
		BlockHash:         common.HexToHash(receipt.BlockHash),
		BlockNumber:       blockNumber,
		TransactionIndex:  uint(txIndex),
	}, nil
}

// transform alchemy.LogResponse -> geth.Log
func TransformAlchemyLogToGeth(log types.LogResponse) (*gethTypes.Log, error) {
	topics := make([]common.Hash, len(log.Topics))
	for i, topic := range log.Topics {
		topics[i] = common.HexToHash(topic)
	}

	blockNumber, err := FromHexU64(log.BlockNumber)
	if err != nil {
		return nil, err
	}

	txIndex, err := FromHexU64(log.TransactionIndex)
	if err != nil {
		return nil, err
	}

	logIndex, err := FromHexU64(log.LogIndex)
	if err != nil {
		return nil, err
	}

	return &gethTypes.Log{
		Address:     common.HexToAddress(log.Address),
		Topics:      topics,
		Data:        []byte(log.Data),
		BlockNumber: blockNumber,
		TxHash:      common.HexToHash(log.TransactionHash),
		TxIndex:     uint(txIndex),
		BlockHash:   common.HexToHash(log.BlockHash),
		Index:       uint(logIndex),
		Removed:     log.Removed,
	}, nil
}
