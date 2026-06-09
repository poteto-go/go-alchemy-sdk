package decode_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/poteto-go/go-alchemy-sdk/constant"
	"github.com/poteto-go/go-alchemy-sdk/decode"
	"github.com/stretchr/testify/assert"
)

// abiStringBytes returns a hand-crafted ABI-encoded string: 32-byte offset word,
// 32-byte length word, then the string data padded to a 32-byte boundary.
func abiStringBytes(s string) []byte {
	n := len(s)
	padded := ((n + constant.ABIWordSize - 1) / constant.ABIWordSize) * constant.ABIWordSize
	b := make([]byte, constant.ABIStringHeaderSize+padded)
	b[constant.ABIWordSize-1] = byte(constant.ABIWordSize)
	b[constant.ABIStringHeaderSize-1] = byte(n)
	copy(b[constant.ABIStringHeaderSize:], s)
	return b
}

func TestABIAddress(t *testing.T) {
	addr := common.HexToAddress("0xabcdef1234567890abcdef1234567890abcdef12")

	t.Run("decodes left-padded 32-byte word", func(t *testing.T) {
		input := make([]byte, 32)
		copy(input[12:], addr.Bytes())

		result, err := decode.ABIAddress(input)

		assert.NoError(t, err)
		assert.Equal(t, addr, result)
	})

	t.Run("returns error for short input", func(t *testing.T) {
		_, err := decode.ABIAddress([]byte{0x01})

		assert.Error(t, err)
	})
}

func TestABIString(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		res, err := decode.ABIString(abiStringBytes("TestToken"))
		assert.NoError(t, err)
		assert.Equal(t, "TestToken", res)
	})

	t.Run("too short", func(t *testing.T) {
		_, err := decode.ABIString([]byte("too-short"))
		assert.Error(t, err)
	})

	t.Run("length mismatch", func(t *testing.T) {
		encoded := make([]byte, constant.ABIStringHeaderSize)
		encoded[constant.ABIStringHeaderSize-1] = 100
		_, err := decode.ABIString(encoded)
		assert.Error(t, err)
	})

	t.Run("returns error for malicious length (no panic)", func(t *testing.T) {
		encoded := make([]byte, constant.ABIStringHeaderSize)
		for i := constant.ABIWordSize; i < constant.ABIStringHeaderSize; i++ {
			encoded[i] = 0xFF
		}
		assert.NotPanics(t, func() {
			_, err := decode.ABIString(encoded)
			assert.Error(t, err)
		})
	})
}

func TestUint256(t *testing.T) {
	v, err := decode.Uint256([]byte{0x01, 0x00})
	assert.NoError(t, err)
	assert.Equal(t, "256", v.String())

	zero, err := decode.Uint256(nil)
	assert.NoError(t, err)
	assert.Equal(t, "0", zero.String())
}

func TestBool(t *testing.T) {
	tests := []struct {
		name   string
		input  []byte
		expect bool
	}{
		{"last byte 1 -> true", append(make([]byte, 31), 0x01), true},
		{"last byte 0 -> false", make([]byte, 32), false},
		{"empty -> false", nil, false},
		{"last byte 2 -> false", append(make([]byte, 31), 0x02), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decode.Bool(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expect, got)
		})
	}
}

func TestUint8(t *testing.T) {
	t.Run("decodes value", func(t *testing.T) {
		v, err := decode.Uint8(append(make([]byte, 31), 0x12))
		assert.NoError(t, err)
		assert.Equal(t, uint8(18), v)
	})

	t.Run("empty -> 0", func(t *testing.T) {
		v, err := decode.Uint8(nil)
		assert.NoError(t, err)
		assert.Equal(t, uint8(0), v)
	})

	t.Run("overflow -> error", func(t *testing.T) {
		_, err := decode.Uint8([]byte{0x01, 0x00})
		assert.Error(t, err)
	})
}

func TestBytes32(t *testing.T) {
	t.Run("decodes the first 32-byte word", func(t *testing.T) {
		in := make([]byte, 32)
		in[0] = 0xde
		in[31] = 0xad
		out, err := decode.Bytes32(in)
		assert.NoError(t, err)
		assert.Equal(t, byte(0xde), out[0])
		assert.Equal(t, byte(0xad), out[31])
	})

	t.Run("too short -> error", func(t *testing.T) {
		_, err := decode.Bytes32(make([]byte, 31))
		assert.Error(t, err)
	})
}
