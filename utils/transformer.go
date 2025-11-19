package utils

import (
	"math"

	"github.com/ethereum/go-ethereum/common"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/poteto-go/go-alchemy-sdk/constant"
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
	if typeU64 > math.MaxUint8 {
		return nil, constant.ErrOverFlow
	}
	typeU8 := uint8(typeU64)

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

	txIndexU64, err := FromHexU64(receipt.TransactionIndex)
	if err != nil {
		return nil, err
	}
	if txIndexU64 > math.MaxUint32 {
		return nil, constant.ErrOverFlow
	}
	txIndex := uint(txIndexU64)

	return &gethTypes.Receipt{
		Type:              typeU8,
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
		TransactionIndex:  txIndex,
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

	txIndexU64, err := FromHexU64(log.TransactionIndex)
	if err != nil {
		return nil, err
	}
	if txIndexU64 > math.MaxUint32 {
		return nil, constant.ErrOverFlow
	}
	txIndex := uint(txIndexU64)

	logIndexU64, err := FromHexU64(log.LogIndex)
	if err != nil {
		return nil, err
	}
	if logIndexU64 > math.MaxUint32 {
		return nil, constant.ErrOverFlow
	}
	logIndex := uint(logIndexU64)

	return &gethTypes.Log{
		Address:     common.HexToAddress(log.Address),
		Topics:      topics,
		Data:        []byte(log.Data),
		BlockNumber: blockNumber,
		TxHash:      common.HexToHash(log.TransactionHash),
		TxIndex:     txIndex,
		BlockHash:   common.HexToHash(log.BlockHash),
		Index:       logIndex,
		Removed:     log.Removed,
	}, nil
}

func TransformTxRequestToGethTxData(txRequest types.TransactionRequest) (gethTypes.TxData, error) {
	toAddress := common.HexToAddress(txRequest.To)
	value, err := FromBigHex(txRequest.Value)
	if err != nil {
		return nil, err
	}

	data := common.FromHex(txRequest.Data)

	// EIP-1559 (DynamicFeeTx)
	if txRequest.MaxFeePerGas != nil || txRequest.MaxPriorityFeePerGas != nil {
		return &gethTypes.DynamicFeeTx{
			ChainID:    txRequest.ChainID,
			Nonce:      txRequest.Nonce,
			GasTipCap:  txRequest.MaxPriorityFeePerGas,
			GasFeeCap:  txRequest.MaxFeePerGas,
			Gas:        txRequest.GasLimit,
			To:         &toAddress,
			Value:      value,
			Data:       data,
			AccessList: gethTypes.AccessList{},
		}, nil
	}

	// EIP-2930 (AccessListTx)
	if txRequest.AccessList != nil {
		return &gethTypes.AccessListTx{
			ChainID:    txRequest.ChainID,
			Nonce:      txRequest.Nonce,
			GasPrice:   txRequest.GasPrice,
			Gas:        txRequest.GasLimit,
			To:         &toAddress,
			Value:      value,
			Data:       data,
			AccessList: gethTypes.AccessList{},
		}, nil
	}

	// LegacyTx (Type 0)
	return &gethTypes.LegacyTx{
		Nonce:    txRequest.Nonce,
		GasPrice: txRequest.GasPrice,
		Gas:      txRequest.GasLimit,
		To:       &toAddress,
		Value:    value,
		Data:     data,
	}, nil
}
