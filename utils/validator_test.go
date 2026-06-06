package utils

import (
	"testing"

	"github.com/poteto-go/go-alchemy-sdk/constant"
	"github.com/stretchr/testify/assert"
)

func TestValidateBlockTag(t *testing.T) {
	t.Run("normal case:", func(t *testing.T) {
		t.Run("latest", func(t *testing.T) {
			// Arrange
			blockTag := "latest"

			// Act
			err := ValidateBlockTag(blockTag)

			// Assert
			assert.NoError(t, err)
		})

		t.Run("earliest", func(t *testing.T) {
			// Arrange
			blockTag := "earliest"

			// Act
			err := ValidateBlockTag(blockTag)

			// Assert
			assert.NoError(t, err)
		})

		t.Run("pending", func(t *testing.T) {
			// Arrange
			blockTag := "pending"

			// Act
			err := ValidateBlockTag(blockTag)

			// Assert
			assert.NoError(t, err)
		})

		t.Run("safe", func(t *testing.T) {
			// Arrange
			blockTag := "safe"

			// Act
			err := ValidateBlockTag(blockTag)

			// Assert
			assert.NoError(t, err)
		})

		t.Run("finalized", func(t *testing.T) {
			// Arrange
			blockTag := "finalized"

			// Act
			err := ValidateBlockTag(blockTag)

			// Assert
			assert.NoError(t, err)
		})

		t.Run("hex block number", func(t *testing.T) {
			// Arrange
			blockTag := "0x1234"

			// Act
			err := ValidateBlockTag(blockTag)

			// Assert
			assert.NoError(t, err)
		})
	})

	t.Run("error case:", func(t *testing.T) {
		cases := []string{
			"unexpected",
			"",
			"0x",
			"0xgg",
		}
		for _, blockTag := range cases {
			t.Run(blockTag, func(t *testing.T) {
				// Act
				err := ValidateBlockTag(blockTag)

				// Assert
				assert.ErrorIs(t, err, constant.ErrInvalidBlockTag)
			})
		}
	})
}

func TestValidateABIString(t *testing.T) {
	t.Run("normal case: returns declared length", func(t *testing.T) {
		// Arrange: header declares length 3 with 3 bytes of data following.
		output := make([]byte, constant.ABIStringHeaderSize+constant.ABIWordSize)
		output[constant.ABIStringHeaderSize-1] = 3

		// Act
		length, err := ValidateABIString(output)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, int64(3), length)
	})

	t.Run("error case:", func(t *testing.T) {
		t.Run("output too short", func(t *testing.T) {
			// Act
			_, err := ValidateABIString([]byte("too-short"))

			// Assert
			assert.ErrorIs(t, err, constant.ErrInvalidABIString)
		})

		t.Run("declared length exceeds output", func(t *testing.T) {
			// Arrange
			output := make([]byte, constant.ABIStringHeaderSize)
			output[constant.ABIStringHeaderSize-1] = 100

			// Act
			_, err := ValidateABIString(output)

			// Assert
			assert.ErrorIs(t, err, constant.ErrInvalidABIString)
		})

		t.Run("malicious length does not panic", func(t *testing.T) {
			// Arrange: lower 8 bytes of the length word set to 0xFF...FF,
			// which would make Int64() negative and slip past a naive check.
			output := make([]byte, constant.ABIStringHeaderSize)
			for i := constant.ABIWordSize; i < constant.ABIStringHeaderSize; i++ {
				output[i] = 0xFF
			}

			// Act & Assert
			assert.NotPanics(t, func() {
				_, err := ValidateABIString(output)
				assert.ErrorIs(t, err, constant.ErrInvalidABIString)
			})
		})
	})
}
