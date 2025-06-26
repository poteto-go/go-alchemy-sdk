package utils_test

import (
	"math/big"
	"testing"

	"github.com/poteto-go/go-alchemy-sdk/core"
	"github.com/poteto-go/go-alchemy-sdk/utils"
	"github.com/stretchr/testify/assert"
)

func TestFromHex(t *testing.T) {
	t.Run("normal case:", func(t *testing.T) {
		// Arrange
		target := "0x123456"

		// Act
		result, err := utils.FromHex(target)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 1193046, result)
	})

	t.Run("error case:", func(t *testing.T) {
		t.Run("empty string => core.ErrInvalidHexString", func(t *testing.T) {
			// Act
			_, err := utils.FromHex("")

			// Assert
			assert.ErrorIs(t, err, core.ErrInvalidHexString)
		})

		t.Run("invalid hex string => core.ErrInvalidHexString", func(t *testing.T) {
			// Act
			_, err := utils.FromHex("unexpected")

			// Assert
			assert.ErrorIs(t, err, core.ErrInvalidHexString)
		})

		t.Run("not number => core.ErrInvalidHexString", func(t *testing.T) {
			// Act
			_, err := utils.FromHex("0xhello")

			// Assert
			assert.ErrorIs(t, err, core.ErrInvalidHexString)
		})
	})
}

func TestFromBigHex(t *testing.T) {
	t.Run("normal case:", func(t *testing.T) {
		t.Run("valid hex string", func(t *testing.T) {
			// Arrange
			target := "0x1234567890abcdef"
			expected := new(big.Int)
			expected.SetString("1311768467294899695", 10)

			// Act
			result, err := utils.FromBigHex(target)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, expected, result)
		})

		t.Run("zero value", func(t *testing.T) {
			// Arrange
			target := "0x0"
			expected := big.NewInt(0)

			// Act
			result, err := utils.FromBigHex(target)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, expected, result)
		})

		t.Run("large hex string", func(t *testing.T) {
			// Arrange
			target := "0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"
			expected := new(big.Int)
			expected.SetString("115792089237316195423570985008687907853269984665640564039457584007913129639935", 10) // 2^256 - 1

			// Act
			result, err := utils.FromBigHex(target)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, expected, result)
		})
	})

	t.Run("error case:", func(t *testing.T) {
		t.Run("empty string => core.ErrInvalidHexString", func(t *testing.T) {
			// Act
			_, err := utils.FromBigHex("")

			// Assert
			assert.ErrorIs(t, err, core.ErrInvalidHexString)
		})

		t.Run("invalid hex string => core.ErrInvalidHexString", func(t *testing.T) {
			// Act
			_, err := utils.FromBigHex("unexpected")

			// Assert
			assert.ErrorIs(t, err, core.ErrInvalidHexString)
		})

		t.Run("not a hex number => core.ErrInvalidHexString", func(t *testing.T) {
			// Act
			_, err := utils.FromBigHex("0xGHIJ")

			// Assert
			assert.ErrorIs(t, err, core.ErrInvalidHexString)
		})

		t.Run("missing 0x prefix => core.ErrInvalidHexString", func(t *testing.T) {
			// Act
			_, err := utils.FromBigHex("12345")

			// Assert
			assert.ErrorIs(t, err, core.ErrInvalidHexString)
		})
	})
}
