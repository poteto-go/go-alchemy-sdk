package encode_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/poteto-go/go-alchemy-sdk/constant"
	"github.com/poteto-go/go-alchemy-sdk/decode"
	"github.com/poteto-go/go-alchemy-sdk/encode"
	"github.com/stretchr/testify/assert"
)

func TestABIString_roundtrip(t *testing.T) {
	t.Run("encode and decode roundtrip", func(t *testing.T) {
		str := "TestToken"
		encoded := encode.ABIString(str)
		res, err := decode.ABIString(encoded)
		assert.NoError(t, err)
		assert.Equal(t, str, res)
	})

	t.Run("empty string", func(t *testing.T) {
		encoded := encode.ABIString("")
		res, err := decode.ABIString(encoded)
		assert.NoError(t, err)
		assert.Equal(t, "", res)
	})
}

func TestABIBytes(t *testing.T) {
	t.Run("encodes length word + right-padded data", func(t *testing.T) {
		data := []byte{0xde, 0xad, 0xbe, 0xef}
		out := encode.ABIBytes(data)

		// length word (4) + one padded data word.
		assert.Equal(t, 2*constant.ABIWordSize, len(out))
		assert.Equal(t, byte(0x04), out[constant.ABIWordSize-1])
		assert.Equal(t, data, out[constant.ABIWordSize:constant.ABIWordSize+4])
		// trailing bytes of the data word are zero-padded.
		assert.Equal(t, make([]byte, constant.ABIWordSize-4), out[constant.ABIWordSize+4:])
	})

	t.Run("empty data -> length word only", func(t *testing.T) {
		out := encode.ABIBytes([]byte{})

		assert.Equal(t, constant.ABIWordSize, len(out))
		assert.Equal(t, make([]byte, constant.ABIWordSize), out)
	})

	t.Run("data spanning multiple words is padded to a word boundary", func(t *testing.T) {
		data := make([]byte, 33)
		out := encode.ABIBytes(data)

		// length word + two data words (33 bytes rounds up to 64).
		assert.Equal(t, 3*constant.ABIWordSize, len(out))
		assert.Equal(t, byte(0x21), out[constant.ABIWordSize-1])
	})
}

func TestReadCalldata(t *testing.T) {
	t.Run("prepends the 4-byte selector and appends args", func(t *testing.T) {
		arg := common.LeftPadBytes(common.HexToAddress("0x1").Bytes(), constant.ABIWordSize)
		data := encode.ReadCalldata([]byte("balanceOf(address)"), arg)

		// balanceOf(address) selector.
		assert.Equal(t, "70a08231", common.Bytes2Hex(data[:4]))
		assert.Equal(t, 4+constant.ABIWordSize, len(data))
		assert.Equal(t, arg, data[4:])
	})

	t.Run("no args -> selector only", func(t *testing.T) {
		data := encode.ReadCalldata([]byte("totalSupply()"))
		assert.Equal(t, 4, len(data))
	})
}

func TestABIAddress(t *testing.T) {
	out := encode.ABIAddress("0xabc0000000000000000000000000000000000abc")

	assert.Equal(t, constant.ABIWordSize, len(out))
	assert.Equal(t, make([]byte, 12), out[:12])
	assert.Equal(t, common.HexToAddress("0xabc0000000000000000000000000000000000abc").Bytes(), out[12:])
}

func TestABIUint256(t *testing.T) {
	t.Run("encodes to 32-byte word", func(t *testing.T) {
		v := big.NewInt(1)
		out := encode.ABIUint256(v)

		assert.Equal(t, constant.ABIWordSize, len(out))
		assert.Equal(t, make([]byte, 31), out[:31])
		assert.Equal(t, byte(0x01), out[31])
	})

	t.Run("encodes zero", func(t *testing.T) {
		out := encode.ABIUint256(big.NewInt(0))

		assert.Equal(t, constant.ABIWordSize, len(out))
		assert.Equal(t, make([]byte, 32), out)
	})

	t.Run("roundtrip with decode.Uint256", func(t *testing.T) {
		v := big.NewInt(123456789)
		out := encode.ABIUint256(v)

		decoded, err := decode.Uint256(out)
		assert.NoError(t, err)
		assert.Equal(t, 0, v.Cmp(decoded))
	})
}

func TestABIBool(t *testing.T) {
	t.Run("encodes true to 0x00..01", func(t *testing.T) {
		out := encode.ABIBool(true)

		assert.Equal(t, constant.ABIWordSize, len(out))
		assert.Equal(t, make([]byte, 31), out[:31])
		assert.Equal(t, byte(0x01), out[31])
	})

	t.Run("encodes false to all zeros", func(t *testing.T) {
		out := encode.ABIBool(false)

		assert.Equal(t, constant.ABIWordSize, len(out))
		assert.Equal(t, make([]byte, 32), out)
	})
}
