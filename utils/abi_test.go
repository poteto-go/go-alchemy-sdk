package utils_test

import (
	"testing"

	"github.com/poteto-go/go-alchemy-sdk/utils"
	"github.com/stretchr/testify/assert"
)

func TestEncodeABIString(t *testing.T) {
	t.Run("encode and decode roundtrip", func(t *testing.T) {
		str := "TestToken"
		encoded := utils.EncodeABIString(str)
		res, err := utils.DecodeABIString(encoded)
		assert.NoError(t, err)
		assert.Equal(t, str, res)
	})

	t.Run("empty string", func(t *testing.T) {
		encoded := utils.EncodeABIString("")
		res, err := utils.DecodeABIString(encoded)
		assert.NoError(t, err)
		assert.Equal(t, "", res)
	})
}

func TestDecodeABIString(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		str := "TestToken"
		encoded := utils.EncodeABIString(str)

		res, err := utils.DecodeABIString(encoded)
		assert.NoError(t, err)
		assert.Equal(t, str, res)
	})

	t.Run("too short", func(t *testing.T) {
		_, err := utils.DecodeABIString([]byte("too-short"))
		assert.Error(t, err)
	})

	t.Run("length mismatch", func(t *testing.T) {
		encoded := make([]byte, utils.ABIStringHeaderSize)
		encoded[utils.ABIStringHeaderSize-1] = 100 // Declared length 100, no data follows
		_, err := utils.DecodeABIString(encoded)
		assert.Error(t, err)
	})
}
