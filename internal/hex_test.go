package internal_test

import (
	"testing"

	"github.com/poteto-go/go-alchemy-sdk/internal"
	"github.com/stretchr/testify/assert"
)

func TestDecodeHex(t *testing.T) {
	t.Run("can decode hex begin w/ 0x", func(t *testing.T) {
		// Arrange
		sampleHex := "0xbcdf20249abf0ed6d944c0288fad489e33f66b3960d9e6229c1cd214ed3bbe31"

		// Act
		decoded, err := internal.DecodeHex(sampleHex)

		// Assert
		assert.Nil(t, err)
		assert.Equal(t, len(decoded), 32)
	})

	t.Run("can decode hex begin w/o 0x", func(t *testing.T) {
		// Arrange
		sampleHex := "bcdf20249abf0ed6d944c0288fad489e33f66b3960d9e6229c1cd214ed3bbe31"

		// Act
		decoded, err := internal.DecodeHex(sampleHex)

		// Assert
		assert.Nil(t, err)
		assert.Equal(t, len(decoded), 32)
	})

	t.Run("error on invalid hex", func(t *testing.T) {
		// Arrange
		sampleHex := "invalid"

		// Act
		_, err := internal.DecodeHex(sampleHex)

		// Assert
		assert.Error(t, err)
	})
}
