package utils_test

import (
	"testing"

	"github.com/poteto-go/go-alchemy-sdk/utils"
	"github.com/stretchr/testify/assert"
)

func TestDecodeABIString(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		str := "TestToken"
		encoded := make([]byte, 96)
		encoded[31] = 0x20           // Offset
		encoded[63] = byte(len(str)) // Length
		copy(encoded[64:], []byte(str))

		res, err := utils.DecodeABIString(encoded)
		assert.NoError(t, err)
		assert.Equal(t, str, res)
	})

	t.Run("too short", func(t *testing.T) {
		_, err := utils.DecodeABIString([]byte("too-short"))
		assert.Error(t, err)
	})

	t.Run("length mismatch", func(t *testing.T) {
		encoded := make([]byte, 64)
		encoded[63] = 100 // Declared length 100
		_, err := utils.DecodeABIString(encoded)
		assert.Error(t, err)
	})
}
