package validate_test

import (
	"math/big"
	"testing"

	"github.com/poteto-go/go-alchemy-sdk/constant"
	"github.com/poteto-go/go-alchemy-sdk/validate"
	"github.com/stretchr/testify/assert"
)

func TestUint256(t *testing.T) {
	tests := []struct {
		name    string
		v       *big.Int
		wantErr error
	}{
		{"nil", nil, constant.ErrNilAmount},
		{"negative", big.NewInt(-1), constant.ErrNegativeAmount},
		{"zero", big.NewInt(0), nil},
		{"positive", big.NewInt(100), nil},
		{
			"uint256 max",
			new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(1)),
			nil,
		},
		{
			"over uint256 max",
			new(big.Int).Lsh(big.NewInt(1), 256),
			constant.ErrAmountExceedsUint256,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.ErrorIs(t, validate.Uint256(tt.v), tt.wantErr)
		})
	}
}

func TestAddress(t *testing.T) {
	tests := []struct {
		name    string
		addr    string
		wantErr error
	}{
		{"valid address", "0xE25583099BA105D9ec0A67f5Ae86D90e50036425", nil},
		{"valid lowercase", "0xe25583099ba105d9ec0a67f5ae86d90e50036425", nil},
		{"empty string", "", constant.ErrInvalidAddress},
		{"too short", "0x1234", constant.ErrInvalidAddress},
		{"not hex", "not-an-address", constant.ErrInvalidAddress},
		{"no 0x prefix", "E25583099BA105D9ec0A67f5Ae86D90e50036425", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.ErrorIs(t, validate.Address(tt.addr), tt.wantErr)
		})
	}
}

func TestAddresses(t *testing.T) {
	valid := "0xE25583099BA105D9ec0A67f5Ae86D90e50036425"
	invalid := "invalid"

	tests := []struct {
		name    string
		addrs   []string
		wantErr error
	}{
		{"all valid", []string{valid, valid}, nil},
		{"first invalid", []string{invalid, valid}, constant.ErrInvalidAddress},
		{"second invalid", []string{valid, invalid}, constant.ErrInvalidAddress},
		{"empty slice", []string{}, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.ErrorIs(t, validate.Addresses(tt.addrs...), tt.wantErr)
		})
	}
}

func TestBlockTag(t *testing.T) {
	tests := []struct {
		name     string
		blockTag string
		wantErr  error
	}{
		{"latest", "latest", nil},
		{"earliest", "earliest", nil},
		{"pending", "pending", nil},
		{"safe", "safe", nil},
		{"finalized", "finalized", nil},
		{"hex block number", "0x1234", nil},
		{"unexpected", "unexpected", constant.ErrInvalidBlockTag},
		{"empty", "", constant.ErrInvalidBlockTag},
		{"only prefix", "0x", constant.ErrInvalidBlockTag},
		{"not hex", "0xgg", constant.ErrInvalidBlockTag},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.ErrorIs(t, validate.BlockTag(tt.blockTag), tt.wantErr)
		})
	}
}

func TestUrl(t *testing.T) {
	tests := []struct {
		name    string
		rawUrl  string
		wantErr error
	}{
		{"empty is allowed", "", nil},
		{"valid http", "http://localhost:8545", nil},
		{"valid https", "https://my-rpc.example.com", nil},
		{"valid ws", "ws://localhost:8546", nil},
		{"valid wss", "wss://my-rpc.example.com/v2/key", nil},
		{"invalid scheme", "ftp://bad-scheme.com", constant.ErrInvalidPrivateNetworkUrl},
		{"missing scheme", "localhost:8545", constant.ErrInvalidPrivateNetworkUrl},
		{"empty host", "http://", constant.ErrInvalidPrivateNetworkUrl},
		{"empty hostname with port", "http://:8545", constant.ErrInvalidPrivateNetworkUrl},
		{"empty ws host", "wss://", constant.ErrInvalidPrivateNetworkUrl},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.ErrorIs(t, validate.Url(tt.rawUrl), tt.wantErr)
		})
	}
}

func TestABIOffsetIsStandard(t *testing.T) {
	t.Run("standard 0x20 offset", func(t *testing.T) {
		out := make([]byte, constant.ABIWordSize)
		out[constant.ABIWordSize-1] = 0x20
		assert.True(t, validate.ABIOffsetIsStandard(out))
	})

	t.Run("non-standard offset value", func(t *testing.T) {
		out := make([]byte, constant.ABIWordSize)
		out[constant.ABIWordSize-1] = 0x40
		assert.False(t, validate.ABIOffsetIsStandard(out))
	})

	t.Run("non-zero high byte", func(t *testing.T) {
		// a huge offset whose low byte happens to be 0x20 must not pass.
		out := make([]byte, constant.ABIWordSize)
		out[0] = 0x01
		out[constant.ABIWordSize-1] = 0x20
		assert.False(t, validate.ABIOffsetIsStandard(out))
	})
}

func TestABIUint256Array(t *testing.T) {
	t.Run("valid: empty array", func(t *testing.T) {
		out := make([]byte, constant.ABIWordSize*2)
		out[constant.ABIWordSize-1] = 0x20 // standard offset
		assert.NoError(t, validate.ABIUint256Array(out))
	})

	t.Run("valid: length fits within output", func(t *testing.T) {
		// offset(32) + length=2(32) + 2 item words(64) = 128 bytes
		out := make([]byte, constant.ABIWordSize*4)
		out[constant.ABIWordSize-1] = 0x20 // standard offset
		out[constant.ABIWordSize*2-1] = 2
		assert.NoError(t, validate.ABIUint256Array(out))
	})

	t.Run("too short: under 64 bytes", func(t *testing.T) {
		assert.ErrorIs(t, validate.ABIUint256Array(make([]byte, 31)), constant.ErrInvalidABIArray)
	})

	t.Run("non-standard offset word", func(t *testing.T) {
		// offset=0x40 instead of the standard 0x20.
		out := make([]byte, constant.ABIWordSize*2)
		out[constant.ABIWordSize-1] = 0x40
		assert.ErrorIs(t, validate.ABIUint256Array(out), constant.ErrInvalidABIArray)
	})

	t.Run("declared length exceeds output", func(t *testing.T) {
		// offset=0x20 + length=5 but no element words follow
		out := append(append(make([]byte, constant.ABIWordSize-1), 0x20), append(make([]byte, constant.ABIWordSize-1), 0x05)...)
		assert.ErrorIs(t, validate.ABIUint256Array(out), constant.ErrInvalidABIArray)
	})

	t.Run("malicious length does not panic", func(t *testing.T) {
		out := make([]byte, constant.ABIWordSize*2)
		out[constant.ABIWordSize-1] = 0x20 // standard offset
		for i := constant.ABIWordSize; i < constant.ABIWordSize*2; i++ {
			out[i] = 0xFF
		}
		assert.NotPanics(t, func() {
			assert.ErrorIs(t, validate.ABIUint256Array(out), constant.ErrInvalidABIArray)
		})
	})
}

func TestABIString(t *testing.T) {
	t.Run("normal case: valid header and length", func(t *testing.T) {
		// header declares length 3 with a full data word following.
		output := make([]byte, constant.ABIStringHeaderSize+constant.ABIWordSize)
		output[constant.ABIWordSize-1] = 0x20 // standard offset
		output[constant.ABIStringHeaderSize-1] = 3

		assert.NoError(t, validate.ABIString(output))
	})

	t.Run("error case:", func(t *testing.T) {
		t.Run("output too short", func(t *testing.T) {
			assert.ErrorIs(t, validate.ABIString([]byte("too-short")), constant.ErrInvalidABIString)
		})

		t.Run("non-standard offset word", func(t *testing.T) {
			// offset=0x40 instead of the standard 0x20.
			output := make([]byte, constant.ABIStringHeaderSize)
			output[constant.ABIWordSize-1] = 0x40

			assert.ErrorIs(t, validate.ABIString(output), constant.ErrInvalidABIString)
		})

		t.Run("declared length exceeds output", func(t *testing.T) {
			output := make([]byte, constant.ABIStringHeaderSize)
			output[constant.ABIWordSize-1] = 0x20 // standard offset
			output[constant.ABIStringHeaderSize-1] = 100

			assert.ErrorIs(t, validate.ABIString(output), constant.ErrInvalidABIString)
		})

		t.Run("malicious length does not panic", func(t *testing.T) {
			// lower 8 bytes of the length word set to 0xFF...FF, which would
			// make Int64() negative and slip past a naive bounds check.
			output := make([]byte, constant.ABIStringHeaderSize)
			output[constant.ABIWordSize-1] = 0x20 // standard offset
			for i := constant.ABIWordSize; i < constant.ABIStringHeaderSize; i++ {
				output[i] = 0xFF
			}

			assert.NotPanics(t, func() {
				assert.ErrorIs(t, validate.ABIString(output), constant.ErrInvalidABIString)
			})
		})
	})
}
