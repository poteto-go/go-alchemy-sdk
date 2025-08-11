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
	})

	t.Run("error case:", func(t *testing.T) {
		// Arrange
		blockTag := "unexpected"

		// Act
		err := ValidateBlockTag(blockTag)

		// Assert
		assert.ErrorIs(t, err, constant.ErrInvalidBlockTag)
	})
}
