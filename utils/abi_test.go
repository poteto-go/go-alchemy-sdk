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

func TestEncodeReadCalldata(t *testing.T) {
	t.Run("prepends the 4-byte selector and appends args", func(t *testing.T) {
		arg := common.LeftPadBytes(common.HexToAddress("0x1").Bytes(), constant.ABIWordSize)
		data := utils.EncodeReadCalldata([]byte("balanceOf(address)"), arg)

		// balanceOf(address) selector.
		assert.Equal(t, "70a08231", common.Bytes2Hex(data[:4]))
		assert.Equal(t, 4+constant.ABIWordSize, len(data))
		assert.Equal(t, arg, data[4:])
	})

	t.Run("no args -> selector only", func(t *testing.T) {
		data := utils.EncodeReadCalldata([]byte("totalSupply()"))
		assert.Equal(t, 4, len(data))
	})
}

func TestEncodeABIAddress(t *testing.T) {
	out := utils.EncodeABIAddress("0xabc0000000000000000000000000000000000abc")

	assert.Equal(t, constant.ABIWordSize, len(out))
	// left-padded: high 12 bytes zero, low 20 bytes the address.
	assert.Equal(t, make([]byte, 12), out[:12])
	assert.Equal(t, common.HexToAddress("0xabc0000000000000000000000000000000000abc").Bytes(), out[12:])
}

func TestDecodeUint256(t *testing.T) {
	v, err := utils.DecodeUint256([]byte{0x01, 0x00})
	assert.NoError(t, err)
	assert.Equal(t, "256", v.String())

	zero, err := utils.DecodeUint256(nil)
	assert.NoError(t, err)
	assert.Equal(t, "0", zero.String())
}

func TestDecodeBool(t *testing.T) {
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
			got, err := utils.DecodeBool(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expect, got)
		})
	}
}

func TestDecodeUint8(t *testing.T) {
	t.Run("decodes value", func(t *testing.T) {
		v, err := utils.DecodeUint8(append(make([]byte, 31), 0x12))
		assert.NoError(t, err)
		assert.Equal(t, uint8(18), v)
	})

	t.Run("empty -> 0", func(t *testing.T) {
		v, err := utils.DecodeUint8(nil)
		assert.NoError(t, err)
		assert.Equal(t, uint8(0), v)
	})

	t.Run("overflow -> error", func(t *testing.T) {
		// 0x0100 == 256, exceeds uint8.
		_, err := utils.DecodeUint8([]byte{0x01, 0x00})
		assert.Error(t, err)
	})
}

func TestDecodeBytes32(t *testing.T) {
	t.Run("decodes the first 32-byte word", func(t *testing.T) {
		in := make([]byte, 32)
		in[0] = 0xde
		in[31] = 0xad
		out, err := utils.DecodeBytes32(in)
		assert.NoError(t, err)
		assert.Equal(t, byte(0xde), out[0])
		assert.Equal(t, byte(0xad), out[31])
	})

	t.Run("too short -> error", func(t *testing.T) {
		_, err := utils.DecodeBytes32(make([]byte, 31))
		assert.Error(t, err)
	})
}
