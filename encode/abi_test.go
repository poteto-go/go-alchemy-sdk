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

func TestABIUint256Array(t *testing.T) {
	t.Run("encodes length word + one word per element", func(t *testing.T) {
		out := encode.ABIUint256Array([]*big.Int{big.NewInt(1), big.NewInt(256)})

		// length word + two element words.
		assert.Equal(t, 3*constant.ABIWordSize, len(out))
		assert.Equal(t, byte(0x02), out[constant.ABIWordSize-1])
		assert.Equal(t, byte(0x01), out[2*constant.ABIWordSize-1])
		assert.Equal(t, byte(0x01), out[3*constant.ABIWordSize-2])
		assert.Equal(t, byte(0x00), out[3*constant.ABIWordSize-1])
	})

	t.Run("empty slice -> length word only", func(t *testing.T) {
		out := encode.ABIUint256Array(nil)

		assert.Equal(t, constant.ABIWordSize, len(out))
		assert.Equal(t, make([]byte, constant.ABIWordSize), out)
	})

	t.Run("roundtrip with decode.Uint256Array", func(t *testing.T) {
		vs := []*big.Int{big.NewInt(0), big.NewInt(7), big.NewInt(123456789)}
		// Prepend the standard single-array return offset (0x20).
		out := append(encode.ABIUint256(big.NewInt(constant.ABIWordSize)), encode.ABIUint256Array(vs)...)

		decoded, err := decode.Uint256Array(out)
		assert.NoError(t, err)
		assert.Equal(t, len(vs), len(decoded))
		for i := range vs {
			assert.Equal(t, 0, vs[i].Cmp(decoded[i]))
		}
	})
}

func TestABIAddressArray(t *testing.T) {
	t.Run("encodes length word + one left-padded word per address", func(t *testing.T) {
		addr := "0xabc0000000000000000000000000000000000abc"
		out := encode.ABIAddressArray([]string{addr})

		// length word + one address word.
		assert.Equal(t, 2*constant.ABIWordSize, len(out))
		assert.Equal(t, byte(0x01), out[constant.ABIWordSize-1])
		assert.Equal(t, make([]byte, 12), out[constant.ABIWordSize:constant.ABIWordSize+12])
		assert.Equal(t, common.HexToAddress(addr).Bytes(), out[constant.ABIWordSize+12:])
	})

	t.Run("empty slice -> length word only", func(t *testing.T) {
		out := encode.ABIAddressArray(nil)

		assert.Equal(t, constant.ABIWordSize, len(out))
		assert.Equal(t, make([]byte, constant.ABIWordSize), out)
	})
}

func TestABIDynamicArgs(t *testing.T) {
	t.Run("emits one offset word per arg then the tails", func(t *testing.T) {
		accounts := encode.ABIAddressArray([]string{"0xabc0000000000000000000000000000000000abc"})
		ids := encode.ABIUint256Array([]*big.Int{big.NewInt(1)})

		out := encode.ABIDynamicArgs(accounts, ids)

		// head: two offset words, then both tails.
		assert.Equal(t, 2*constant.ABIWordSize+len(accounts)+len(ids), len(out))
		// first offset points past the head (2 words = 0x40).
		assert.Equal(t, byte(0x40), out[constant.ABIWordSize-1])
		// second offset points past the head + accounts tail.
		expectedSecond := big.NewInt(int64(2*constant.ABIWordSize + len(accounts)))
		assert.Equal(t, encode.ABIUint256(expectedSecond), out[constant.ABIWordSize:2*constant.ABIWordSize])
		// tails follow the head in order.
		assert.Equal(t, accounts, out[2*constant.ABIWordSize:2*constant.ABIWordSize+len(accounts)])
		assert.Equal(t, ids, out[2*constant.ABIWordSize+len(accounts):])
	})

	t.Run("no args -> empty", func(t *testing.T) {
		assert.Equal(t, 0, len(encode.ABIDynamicArgs()))
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
