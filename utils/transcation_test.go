package utils_test

import (
	"math/big"
	"testing"

	"github.com/poteto-go/go-alchemy-sdk/constant"
	"github.com/poteto-go/go-alchemy-sdk/types"
	"github.com/poteto-go/go-alchemy-sdk/utils"
	"github.com/stretchr/testify/assert"
)

func newRawTransaction() types.TransactionRawResponse {
	return types.TransactionRawResponse{
		BlockHash:   "0x1",
		BlockNumber: "0x2",
		From:        "0x3",
		To:          "0x4",
		Gas:         "0x5",
		GasPrice:    "0x6",
		Hash:        "0x7",
		Input:       "0x8",
		Nonce:       "0x9",
		Value:       "0xa",
		ChainId:     "0xb",
		V:           "0xc",
		R:           "0xd",
		S:           "0xe",
		Type:        "0x11",
	}
}

func TestTransformTransaction(t *testing.T) {
	t.Run("normal case", func(t *testing.T) {
		t.Run("can transform from rawTransaction to transaction", func(t *testing.T) {
			// Arrange
			rawTransaction := newRawTransaction()
			expected := types.TransactionResponse{
				BlockHash:   "0x1",
				Index:       1,
				BlockNumber: 2,
				From:        "0x3",
				To:          "0x4",
				GasLimit:    big.NewInt(5),
				GasPrice:    big.NewInt(6),
				Hash:        "0x7",
				Data:        "0x8",
				Nonce:       9,
				Type:        17,
				Value:       big.NewInt(10),
				ChainID:     11,
				Signature: types.Signature{
					R: "0xd",
					S: "0xe",
					V: big.NewInt(12),
				},
				MaxPriorityFeePerGas: big.NewInt(0),
				MaxFeePerGas:         big.NewInt(0),
				AccessList:           []string{},
				BlobVersionedHashes:  []string{},
				AuthorizationList:    []string{},
			}

			// Act
			actual, err := utils.TransformTransaction(rawTransaction)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, expected, actual)
		})
	})

	t.Run("error case", func(t *testing.T) {
		t.Run("if error on blockNumber transform, return constant.ErrFailedToTransformBlockNumber", func(t *testing.T) {
			// Arrange
			rawTransaction := newRawTransaction()
			rawTransaction.BlockNumber = "hoge"

			// Act
			_, err := utils.TransformTransaction(rawTransaction)

			// Assert
			assert.ErrorIs(t, err, constant.ErrFailedToTransformBlockNumber)
		})

		t.Run("if error on type transform, return constant.ErrFailedToTransformType", func(t *testing.T) {
			// Arrange
			rawTransaction := newRawTransaction()
			rawTransaction.Type = "hoge"

			// Act
			_, err := utils.TransformTransaction(rawTransaction)

			// Assert
			assert.ErrorIs(t, err, constant.ErrFailedToTransformType)
		})

		t.Run("if error on nonce transform, return constant.ErrFailedToTransformNonce", func(t *testing.T) {
			// Arrange
			rawTransaction := newRawTransaction()
			rawTransaction.Nonce = "hoge"

			// Act
			_, err := utils.TransformTransaction(rawTransaction)

			// Assert
			assert.ErrorIs(t, err, constant.ErrFailedToTransformNonce)
		})

		t.Run("if error on gasPrice transform, return constant.ErrFailedToTransformGasPrice", func(t *testing.T) {
			// Arrange
			rawTransaction := newRawTransaction()
			rawTransaction.GasPrice = "hoge"

			// Act
			_, err := utils.TransformTransaction(rawTransaction)

			// Assert
			assert.ErrorIs(t, err, constant.ErrFailedToTransformGasPrice)
		})

		t.Run("if error on gasLimit transform, return constant.ErrFailedToTransformGasLimit", func(t *testing.T) {
			// Arrange
			rawTransaction := newRawTransaction()
			rawTransaction.Gas = "hoge"

			// Act
			_, err := utils.TransformTransaction(rawTransaction)

			// Assert
			assert.ErrorIs(t, err, constant.ErrFailedToTransformGasLimit)
		})

		t.Run("if error on value transform, return constant.ErrFailedToTransformValue", func(t *testing.T) {
			// Arrange
			rawTransaction := newRawTransaction()
			rawTransaction.Value = "hoge"

			// Act
			_, err := utils.TransformTransaction(rawTransaction)

			// Assert
			assert.ErrorIs(t, err, constant.ErrFailedToTransformValue)
		})

		t.Run("if error on chainId transform, return constant.ErrFailedToTransformChainId", func(t *testing.T) {
			// Arrange
			rawTransaction := newRawTransaction()
			rawTransaction.ChainId = "hoge"

			// Act
			_, err := utils.TransformTransaction(rawTransaction)

			// Assert
			assert.ErrorIs(t, err, constant.ErrFailedToTransformChainId)
		})

		t.Run("if error on v transform, return constant.ErrFailedToTransformV", func(t *testing.T) {
			// Arrange
			rawTransaction := newRawTransaction()
			rawTransaction.V = "hoge"

			// Act
			_, err := utils.TransformTransaction(rawTransaction)

			// Assert
			assert.ErrorIs(t, err, constant.ErrFailedToTransformV)
		})
	})
}
