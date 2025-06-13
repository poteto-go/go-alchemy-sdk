package utils_test

import (
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
