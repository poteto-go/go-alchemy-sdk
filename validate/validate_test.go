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
		{"invalid scheme", "ftp://bad-scheme.com", constant.ErrInvalidPrivateNetworkUrl},
		{"missing scheme", "localhost:8545", constant.ErrInvalidPrivateNetworkUrl},
		{"empty host", "http://", constant.ErrInvalidPrivateNetworkUrl},
		{"empty hostname with port", "http://:8545", constant.ErrInvalidPrivateNetworkUrl},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.ErrorIs(t, validate.Url(tt.rawUrl), tt.wantErr)
		})
	}
}

func TestABIString(t *testing.T) {
	t.Run("normal case: valid header and length", func(t *testing.T) {
		// header declares length 3 with a full data word following.
		output := make([]byte, constant.ABIStringHeaderSize+constant.ABIWordSize)
		output[constant.ABIStringHeaderSize-1] = 3

		assert.NoError(t, validate.ABIString(output))
	})

	t.Run("error case:", func(t *testing.T) {
		t.Run("output too short", func(t *testing.T) {
			assert.ErrorIs(t, validate.ABIString([]byte("too-short")), constant.ErrInvalidABIString)
		})

		t.Run("declared length exceeds output", func(t *testing.T) {
			output := make([]byte, constant.ABIStringHeaderSize)
			output[constant.ABIStringHeaderSize-1] = 100

			assert.ErrorIs(t, validate.ABIString(output), constant.ErrInvalidABIString)
		})

		t.Run("malicious length does not panic", func(t *testing.T) {
			// lower 8 bytes of the length word set to 0xFF...FF, which would
			// make Int64() negative and slip past a naive bounds check.
			output := make([]byte, constant.ABIStringHeaderSize)
			for i := constant.ABIWordSize; i < constant.ABIStringHeaderSize; i++ {
				output[i] = 0xFF
			}

			assert.NotPanics(t, func() {
				assert.ErrorIs(t, validate.ABIString(output), constant.ErrInvalidABIString)
			})
		})
	})
}
