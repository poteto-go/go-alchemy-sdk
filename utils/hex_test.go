package utils_test

import (
	"math/big"
	"testing"

	"github.com/poteto-go/go-alchemy-sdk/constant"
	"github.com/poteto-go/go-alchemy-sdk/utils"
	"github.com/stretchr/testify/assert"
)

func TestFromHex(t *testing.T) {
	t.Run("normal case:", func(t *testing.T) {
		t.Run("can transform from hex", func(t *testing.T) {
			// Arrange
			target := "0x123456"

			// Act
			result, err := utils.FromHex(target)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, 1193046, result)
		})

		t.Run("empty string => 0", func(t *testing.T) {
			// Act
			res, _ := utils.FromHex("")

			// Assert
			assert.Equal(t, 0, res)
		})
	})

	t.Run("error case:", func(t *testing.T) {
		t.Run("invalid one length string => constant.ErrInvalidHexString", func(t *testing.T) {
			// Act
			_, err := utils.FromHex("0")

			// Assert
			assert.ErrorIs(t, err, constant.ErrInvalidHexString)
		})
		t.Run("invalid hex string => constant.ErrInvalidHexString", func(t *testing.T) {
			// Act
			_, err := utils.FromHex("unexpected")

			// Assert
			assert.ErrorIs(t, err, constant.ErrInvalidHexString)
		})

		t.Run("not number => constant.ErrInvalidHexString", func(t *testing.T) {
			// Act
			_, err := utils.FromHex("0xhello")

			// Assert
			assert.ErrorIs(t, err, constant.ErrInvalidHexString)
		})
	})
}

func TestFromHexU64(t *testing.T) {
	t.Run("normal case:", func(t *testing.T) {
		t.Run("can transform from hex", func(t *testing.T) {
			// Arrange
			target := "0x123456"

			// Act
			result, err := utils.FromHexU64(target)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, uint64(1193046), result)
		})

		t.Run("empty string => 0", func(t *testing.T) {
			// Act
			res, _ := utils.FromHexU64("")

			// Assert
			assert.Equal(t, uint64(0), res)
		})
	})

	t.Run("error case:", func(t *testing.T) {
		t.Run("invalid one length string => constant.ErrInvalidHexString", func(t *testing.T) {
			// Act
			_, err := utils.FromHexU64("0")

			// Assert
			assert.ErrorIs(t, err, constant.ErrInvalidHexString)
		})
		t.Run("invalid hex string => constant.ErrInvalidHexString", func(t *testing.T) {
			// Act
			_, err := utils.FromHexU64("unexpected")

			// Assert
			assert.ErrorIs(t, err, constant.ErrInvalidHexString)
		})

		t.Run("not number => constant.ErrInvalidHexString", func(t *testing.T) {
			// Act
			_, err := utils.FromHexU64("0xhello")

			// Assert
			assert.ErrorIs(t, err, constant.ErrInvalidHexString)
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

		t.Run("empty string => 0", func(t *testing.T) {
			// Act
			res, err := utils.FromBigHex("")

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, res.Cmp(big.NewInt(0)), 0)
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
		t.Run("one letter => constant.ErrInvalidHexString", func(t *testing.T) {
			// Act
			_, err := utils.FromBigHex("0")

			// Assert
			assert.ErrorIs(t, err, constant.ErrInvalidHexString)
		})

		t.Run("invalid hex string => constant.ErrInvalidHexString", func(t *testing.T) {
			// Act
			_, err := utils.FromBigHex("unexpected")

			// Assert
			assert.ErrorIs(t, err, constant.ErrInvalidHexString)
		})

		t.Run("not a hex number => constant.ErrInvalidHexString", func(t *testing.T) {
			// Act
			_, err := utils.FromBigHex("0xGHIJ")

			// Assert
			assert.ErrorIs(t, err, constant.ErrInvalidHexString)
		})

		t.Run("missing 0x prefix => constant.ErrInvalidHexString", func(t *testing.T) {
			// Act
			_, err := utils.FromBigHex("12345")

			// Assert
			assert.ErrorIs(t, err, constant.ErrInvalidHexString)
		})
	})
}
