package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToBlockNumber(t *testing.T) {
	t.Run("if latest return nil, nil", func(t *testing.T) {
		// Act
		res, err := ToBlockNumber("latest")

		// Assert
		assert.Nil(t, res)
		assert.Nil(t, err)
	})

	t.Run("if from big hex, return bigNumber", func(t *testing.T) {
		// Act
		res, err := ToBlockNumber("0x123")

		// Assert
		assert.NotNil(t, res)
		assert.Nil(t, err)
	})

	t.Run("if failed to from bigHex, return err", func(t *testing.T) {
		// Act
		_, err := ToBlockNumber("unxpected")

		// Assert
		assert.Error(t, err)
	})
}
