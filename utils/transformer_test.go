package utils_test

import (
	"math"
	"math/big"
	"strconv"
	"testing"

	"github.com/poteto-go/go-alchemy-sdk/types"
	"github.com/poteto-go/go-alchemy-sdk/utils"
	"github.com/stretchr/testify/assert"
)

func TestTransformAlchemyReceiptToGeth(t *testing.T) {
	validLogs := []types.LogResponse{
		{
			Address: "0x1",
			Topics: []string{
				"0x2",
				"0x3",
			},
			Data:             "0x1",
			BlockNumber:      "0x1",
			TransactionHash:  "0x1",
			BlockHash:        "0x1",
			TransactionIndex: "0x1",
			LogIndex:         "0x1",
			Removed:          false,
		},
	}
	t.Run("transform alchemy.TransactionReceipt -> geth.Receipt", func(t *testing.T) {
		// Arrange
		receipt := types.TransactionReceipt{
			TransactionHash:   "0x504ce587a65bdbdb6414a0c6c16d86a04dd79bfcc4f2950eec9634b30ce5370f",
			TransactionIndex:  "0x0",
			BlockHash:         "0xe7212a92cfb9b06addc80dec2a0dfae9ea94fd344efeb157c41e12994fcad60a",
			BlockNumber:       "0x50",
			From:              "0x627306090abab3a6e1400e9345bc60c78a8bef57",
			To:                "0xf17f52151ebef6c7334fad080c5704d77216b732",
			ContractAddress:   "0xf17f52151EbEF6C7334FAD080c5704D77216b732",
			CumulativeGasUsed: "0x1458",
			GasUsed:           "0x1458",
			BlobGasUsed:       "0x1458",
			Logs:              validLogs,
			LogsBloom:         "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
			EffectiveGasPrice: "0x1",
			Type:              "0x0",
			Status:            "0x1",
		}

		// Act
		gethReceipt, err := utils.TransformAlchemyReceiptToGeth(receipt)

		// Assert
		assert.Nil(t, err)
		assert.Equal(t, gethReceipt.Type, uint8(0))
		assert.Equal(t, gethReceipt.Status, uint64(1))
		assert.Equal(t, gethReceipt.CumulativeGasUsed, uint64(5208))
		assert.Equal(t, gethReceipt.GasUsed, uint64(5208))
		assert.Equal(t, gethReceipt.EffectiveGasPrice.Cmp(big.NewInt(1)), 0)
		assert.Equal(t, gethReceipt.BlobGasUsed, uint64(5208))
		assert.Equal(t, gethReceipt.BlockNumber.Cmp(big.NewInt(80)), 0)
		assert.Equal(t, gethReceipt.TransactionIndex, uint(0))
		assert.Equal(t, gethReceipt.BlockHash.Hex(), receipt.BlockHash)
		assert.Equal(t, gethReceipt.TxHash.Hex(), receipt.TransactionHash)
		assert.Equal(t, gethReceipt.ContractAddress.Hex(), receipt.ContractAddress)
		assert.Equal(t, len(gethReceipt.Logs), 1)
	})

	t.Run("error case", func(t *testing.T) {
		t.Run("if logs is invalid", func(t *testing.T) {
			// Arrange
			receipt := types.TransactionReceipt{
				Logs: []types.LogResponse{
					{
						LogIndex: "hello",
					},
				},
			}

			// Act
			_, err := utils.TransformAlchemyReceiptToGeth(receipt)

			// Assert
			assert.Error(t, err)
		})

		t.Run("if type is invalid", func(t *testing.T) {
			// Arrange
			receipt := types.TransactionReceipt{
				Type: "hello",
				Logs: validLogs,
			}

			// Act
			_, err := utils.TransformAlchemyReceiptToGeth(receipt)

			// Assert
			assert.Error(t, err)
		})

		t.Run("if type is overflow", func(t *testing.T) {
			// Arrange
			var overFlowed uint64 = math.MaxUint
			receipt := types.TransactionReceipt{
				Type: "0x" + strconv.FormatUint(overFlowed, 16),
				Logs: validLogs,
			}

			// Act
			_, err := utils.TransformAlchemyReceiptToGeth(receipt)

			// Assert
			assert.Error(t, err)
		})

		t.Run("if status is invalid", func(t *testing.T) {
			// Arrange
			receipt := types.TransactionReceipt{
				Status: "hello",
				Logs:   validLogs,
				Type:   "0x1",
			}

			// Act
			_, err := utils.TransformAlchemyReceiptToGeth(receipt)

			// Assert
			assert.Error(t, err)
		})

		t.Run("if invalid cumulativeGasUsed", func(t *testing.T) {
			// Arrange
			receipt := types.TransactionReceipt{
				Status:            "0x1",
				Logs:              validLogs,
				Type:              "0x1",
				CumulativeGasUsed: "hello",
			}

			// Act
			_, err := utils.TransformAlchemyReceiptToGeth(receipt)

			// Assert
			assert.Error(t, err)
		})

		t.Run("if invalid GasUsed", func(t *testing.T) {
			// Arrange
			receipt := types.TransactionReceipt{
				Status:            "0x1",
				Logs:              validLogs,
				Type:              "0x1",
				CumulativeGasUsed: "0x1",
				GasUsed:           "hello",
			}

			// Act
			_, err := utils.TransformAlchemyReceiptToGeth(receipt)

			// Assert
			assert.Error(t, err)
		})

		t.Run("if invalid EffectiveGasUsed", func(t *testing.T) {
			// Arrange
			receipt := types.TransactionReceipt{
				Status:            "0x1",
				Logs:              validLogs,
				Type:              "0x1",
				CumulativeGasUsed: "0x1",
				GasUsed:           "0x1",
				EffectiveGasPrice: "hello",
			}

			// Act
			_, err := utils.TransformAlchemyReceiptToGeth(receipt)

			// Assert
			assert.Error(t, err)
		})

		t.Run("if invalid BlobGasUsed", func(t *testing.T) {
			// Arrange
			receipt := types.TransactionReceipt{
				Status:            "0x1",
				Logs:              validLogs,
				Type:              "0x1",
				CumulativeGasUsed: "0x1",
				GasUsed:           "0x1",
				EffectiveGasPrice: "0x1",
				BlobGasUsed:       "hello",
			}

			// Act
			_, err := utils.TransformAlchemyReceiptToGeth(receipt)

			// Assert
			assert.Error(t, err)
		})

		t.Run("if invalid BlockNumber", func(t *testing.T) {
			// Arrange
			receipt := types.TransactionReceipt{
				Status:            "0x1",
				Logs:              validLogs,
				Type:              "0x1",
				CumulativeGasUsed: "0x1",
				GasUsed:           "0x1",
				EffectiveGasPrice: "0x1",
				BlobGasUsed:       "0x1",
				BlockNumber:       "hello",
			}

			// Act
			_, err := utils.TransformAlchemyReceiptToGeth(receipt)

			// Assert
			assert.Error(t, err)
		})

		t.Run("if invalid TransactionIndex", func(t *testing.T) {
			// Arrange
			receipt := types.TransactionReceipt{
				Status:            "0x1",
				Logs:              validLogs,
				Type:              "0x1",
				CumulativeGasUsed: "0x1",
				GasUsed:           "0x1",
				EffectiveGasPrice: "0x1",
				BlobGasUsed:       "0x1",
				BlockNumber:       "0x1",
				TransactionIndex:  "hello",
			}

			// Act
			_, err := utils.TransformAlchemyReceiptToGeth(receipt)

			// Assert
			assert.Error(t, err)
		})

		t.Run("if TransactionIndex overflow", func(t *testing.T) {
			// Arrange
			var overFlowed uint64 = math.MaxUint
			receipt := types.TransactionReceipt{
				Status:            "0x1",
				Logs:              validLogs,
				Type:              "0x1",
				CumulativeGasUsed: "0x1",
				GasUsed:           "0x1",
				EffectiveGasPrice: "0x1",
				BlobGasUsed:       "0x1",
				BlockNumber:       "0x1",
				TransactionIndex:  "0x" + strconv.FormatUint(overFlowed, 16),
			}

			// Act
			_, err := utils.TransformAlchemyReceiptToGeth(receipt)

			// Assert
			assert.Error(t, err)
		})
	})
}

func TestTransformAlchemyLogToGeth(t *testing.T) {
	t.Run("transform alchemy.LogResponse -> geth.Log", func(t *testing.T) {
		// Arrange
		address := "0x0000000000000001111111111111111111111111"
		topic1 := "0x0000000000000000000000000000000000000000000000000000000000000001"
		topic2 := "0x0000000000000000000000000000000000000000000000000000000000000002"
		data := "0x0000000000000000000000000000000000000000000000000000000000000003"
		blockNumber := "0x1"
		blockHash := "0x0000000000000000000000000000000000000000000000000000000000000001"
		transactionHash := "0x0000000000000000000000000000000000000000000000000000000000000002"
		transactionIndex := "0x0000000000000000000000000000000000000000000000000000000000000002"
		logIndex := "0x0000000000000000000000000000000000000000000000000000000000000003"

		LogResponse := types.LogResponse{
			Address: address,
			Topics: []string{
				topic1,
				topic2,
			},
			Data:             data,
			BlockNumber:      blockNumber,
			BlockHash:        blockHash,
			TransactionHash:  transactionHash,
			TransactionIndex: transactionIndex,
			LogIndex:         logIndex,
			Removed:          false,
		}

		// Act
		log, err := utils.TransformAlchemyLogToGeth(LogResponse)

		// Assert
		assert.Nil(t, err)
		assert.Equal(t, log.Address.Hex(), address)
		assert.Equal(t, log.Topics[0].Hex(), topic1)
		assert.Equal(t, log.Topics[1].Hex(), topic2)
		assert.Equal(t, string(log.Data), data)
		assert.Equal(t, log.BlockNumber, uint64(1))
		assert.Equal(t, log.TxHash.Hex(), transactionHash)
		assert.Equal(t, log.TxIndex, uint(2))
		assert.Equal(t, log.BlockHash.Hex(), blockHash)
		assert.Equal(t, log.Index, uint(3))
		assert.Equal(t, log.Removed, false)
	})

	t.Run("error case", func(t *testing.T) {
		t.Run("if blockNumber is invalid", func(t *testing.T) {
			// Arrange
			blockNumber := "hello"
			logResponse := types.LogResponse{
				BlockNumber: blockNumber,
			}

			// Act
			_, err := utils.TransformAlchemyLogToGeth(logResponse)

			// Assert
			assert.Error(t, err)
		})

		t.Run("if txIndex is invalid", func(t *testing.T) {
			// Arrange
			blockNumber := "0x123"
			txIndex := "hello"
			logResponse := types.LogResponse{
				BlockNumber:      blockNumber,
				TransactionIndex: txIndex,
			}

			// Act
			_, err := utils.TransformAlchemyLogToGeth(logResponse)

			// Assert
			assert.Error(t, err)
		})

		t.Run("if txIndex overflow", func(t *testing.T) {
			// Arrange
			blockNumber := "0x123"
			var overFlowed uint64 = math.MaxUint
			txIndex := "0x" + strconv.FormatUint(overFlowed, 16)

			logResponse := types.LogResponse{
				BlockNumber:      blockNumber,
				TransactionIndex: txIndex,
			}

			// Act
			_, err := utils.TransformAlchemyLogToGeth(logResponse)

			// Assert
			assert.Error(t, err)
		})

		t.Run("if logIndex is invalid", func(t *testing.T) {
			// Arrange
			blockNumber := "0x123"
			txIndex := "0x123"
			logIndex := "hello"

			// Act
			_, err := utils.TransformAlchemyLogToGeth(types.LogResponse{
				BlockNumber:      blockNumber,
				TransactionIndex: txIndex,
				LogIndex:         logIndex,
			})

			// Assert
			assert.Error(t, err)
		})

		t.Run("if logIndex overflow", func(t *testing.T) {
			// Arrange
			blockNumber := "0x123"
			txIndex := "0x123"
			var overFlowed uint64 = math.MaxUint
			logIndex := "0x" + strconv.FormatUint(overFlowed, 16)

			// Act
			_, err := utils.TransformAlchemyLogToGeth(types.LogResponse{
				BlockNumber:      blockNumber,
				TransactionIndex: txIndex,
				LogIndex:         logIndex,
			})

			// Assert
			assert.Error(t, err)
		})
	})
}
