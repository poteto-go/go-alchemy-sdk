package utils

import (
	"math/big"

	"github.com/poteto-go/go-alchemy-sdk/constant"
	"github.com/poteto-go/go-alchemy-sdk/types"
)

func TransformTransaction(rawTx types.TransactionRawResponse) (types.TransactionResponse, error) {
	blockNumber, err := FromHex(rawTx.BlockNumber)
	if err != nil {
		return types.TransactionResponse{}, constant.ErrFailedToTransformBlockNumber
	}

	typeInt, err := FromHex(rawTx.Type)
	if err != nil {
		return types.TransactionResponse{}, constant.ErrFailedToTransformType
	}
	nonce, err := FromHex(rawTx.Nonce)
	if err != nil {
		return types.TransactionResponse{}, constant.ErrFailedToTransformNonce
	}
	gasPrice, err := FromBigHex(rawTx.GasPrice)
	if err != nil {
		return types.TransactionResponse{}, constant.ErrFailedToTransformGasPrice
	}
	gasLimit, err := FromBigHex(rawTx.Gas)
	if err != nil {
		return types.TransactionResponse{}, constant.ErrFailedToTransformGasLimit
	}
	valueInt, err := FromBigHex(rawTx.Value)
	if err != nil {
		return types.TransactionResponse{}, constant.ErrFailedToTransformValue
	}
	chainId, err := FromHex(rawTx.ChainId)
	if err != nil {
		return types.TransactionResponse{}, constant.ErrFailedToTransformChainId
	}
	intV, err := FromBigHex(rawTx.V)
	if err != nil {
		return types.TransactionResponse{}, constant.ErrFailedToTransformV
	}

	return types.TransactionResponse{
		BlockNumber:          blockNumber,
		BlockHash:            rawTx.BlockHash,
		Index:                1, // TODO: 仮
		Hash:                 rawTx.Hash,
		Type:                 typeInt,
		To:                   rawTx.To,
		From:                 rawTx.From,
		Nonce:                nonce,
		GasPrice:             gasPrice,
		GasLimit:             gasLimit,      // TODO: 仮
		MaxPriorityFeePerGas: big.NewInt(0), // TODO: 仮
		MaxFeePerGas:         big.NewInt(0), // TODO: 仮
		Data:                 rawTx.Input,   // TODO: 仮
		Value:                valueInt,
		ChainID:              chainId,
		Signature: types.Signature{
			R: rawTx.R,
			S: rawTx.S,
			V: intV,
		},
		AccessList:          []string{}, // TODO: 仮
		BlobVersionedHashes: []string{}, // TODO: 仮
		AuthorizationList:   []string{}, // TODO: 仮
	}, nil
}
