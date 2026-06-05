package utils_test

import (
	"strings"
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

		result, err := utils.DecodeABIAddress(input)

		assert.NoError(t, err)
		assert.Equal(t, addr, result)
	})

	t.Run("returns error for short input", func(t *testing.T) {
		_, err := utils.DecodeABIAddress([]byte{0x01})

		assert.Error(t, err)
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

	t.Run("roundtrip with string longer than 255 bytes", func(t *testing.T) {
		str := strings.Repeat("a", 300)
		encoded := utils.EncodeABIString(str)
		res, err := utils.DecodeABIString(encoded)
		assert.NoError(t, err)
		assert.Equal(t, str, res)
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

	t.Run("returns error for malicious length (no panic)", func(t *testing.T) {
		// length word's lower 8 bytes set to 0xFF...FF => Int64() == -1,
		// which previously slipped past the bounds check and panicked.
		encoded := make([]byte, constant.ABIStringHeaderSize)
		for i := constant.ABIWordSize; i < constant.ABIStringHeaderSize; i++ {
			encoded[i] = 0xFF
		}
		assert.NotPanics(t, func() {
			_, err := utils.DecodeABIString(encoded)
			assert.Error(t, err)
		})
	})
}
