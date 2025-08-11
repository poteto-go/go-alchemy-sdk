package utils_test

import (
	"math/big"
	"testing"

	"github.com/poteto-go/go-alchemy-sdk/constant"
	"github.com/poteto-go/go-alchemy-sdk/types"
	"github.com/poteto-go/go-alchemy-sdk/utils"
	"github.com/stretchr/testify/assert"
)

func newRawBlock() types.BlockResponse {
	return types.BlockResponse{
		Hash:         "0x1",
		ParentHash:   "0x2",
		Number:       "0x3",
		Timestamp:    "0x4",
		Nonce:        "0x5",
		Difficulty:   "0x6",
		GasLimit:     "0x7",
		GasUsed:      "0x8",
		Miner:        "0x9",
		Transactions: []string{"0xa", "0xb"},
	}
}

func TestTransformBlock(t *testing.T) {
	t.Run("normal case", func(t *testing.T) {
		t.Run("can transform from rawBlock to block", func(t *testing.T) {
			// Arrange
			rawBlock := newRawBlock()
			expected := types.Block{
				Hash:         "0x1",
				ParentHash:   "0x2",
				Number:       3,
				Timestamp:    4,
				Nonce:        "0x5",
				Difficulty:   6,
				GasLimit:     big.NewInt(7),
				GasUsed:      big.NewInt(8),
				Miner:        "0x9",
				Transactions: []string{"0xa", "0xb"},
			}

			// Act
			actual, err := utils.TransformBlock(rawBlock)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, expected, actual)
		})
	})

	t.Run("error case", func(t *testing.T) {
		t.Run("if error on blockNumber transform, return constant.ErrFailedToTransformBlockNumber", func(t *testing.T) {
			// Arrange
			rawBlock := newRawBlock()
			rawBlock.Number = "hoge"

			// Act
			_, err := utils.TransformBlock(rawBlock)

			// Assert
			assert.ErrorIs(t, err, constant.ErrFailedToTransformBlockNumber)
		})

		t.Run("if error on timestamp transform, return constant.ErrFailedToTransformBlockNumber", func(t *testing.T) {
			// Arrange
			rawBlock := newRawBlock()
			rawBlock.Timestamp = "hoge"

			// Act
			_, err := utils.TransformBlock(rawBlock)

			// Assert
			assert.ErrorIs(t, err, constant.ErrFailedToTransformBlockNumber)
		})

		t.Run("if error on difficulty transform, return constant.ErrFailedToTransformDifficulty", func(t *testing.T) {
			// Arrange
			rawBlock := newRawBlock()
			rawBlock.Difficulty = "hoge"

			// Act
			_, err := utils.TransformBlock(rawBlock)

			// Assert
			assert.ErrorIs(t, err, constant.ErrFailedToTransformDifficulty)
		})

		t.Run("if error on gasLimit transform, return constant.ErrFailedToTransformGasLimit", func(t *testing.T) {
			// Arrange
			rawBlock := newRawBlock()
			rawBlock.GasLimit = "hoge"

			// Act
			_, err := utils.TransformBlock(rawBlock)

			// Assert
			assert.ErrorIs(t, err, constant.ErrFailedToTransformGasLimit)
		})

		t.Run("if error on gasUsed transform, return constant.ErrFailedToTransformGasLimit", func(t *testing.T) {
			// Arrange
			rawBlock := newRawBlock()
			rawBlock.GasUsed = "hoge"

			// Act
			_, err := utils.TransformBlock(rawBlock)

			// Assert
			assert.ErrorIs(t, err, constant.ErrFailedToTransformGasLimit)
		})
	})
}
