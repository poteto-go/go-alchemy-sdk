package utils_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/poteto-go/go-alchemy-sdk/constant"
	"github.com/poteto-go/go-alchemy-sdk/utils"
	"github.com/stretchr/testify/assert"
)

func TestDecodeABIAddress(t *testing.T) {
	addr := common.HexToAddress("0xabcdef1234567890abcdef1234567890abcdef12")

	t.Run("decodes left-padded 32-byte word", func(t *testing.T) {
		input := make([]byte, 32)
		copy(input[12:], addr.Bytes())

		result := utils.DecodeABIAddress(input)

		assert.Equal(t, addr, result)
	})

	t.Run("returns zero address for short input", func(t *testing.T) {
		result := utils.DecodeABIAddress([]byte{0x01})

		assert.Equal(t, common.Address{}, result)
	})
}

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
		encoded := make([]byte, constant.ABIStringHeaderSize)
		encoded[constant.ABIStringHeaderSize-1] = 100 // Declared length 100, no data follows
		_, err := utils.DecodeABIString(encoded)
		assert.Error(t, err)
	})
}
