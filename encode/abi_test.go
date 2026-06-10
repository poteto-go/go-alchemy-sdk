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
