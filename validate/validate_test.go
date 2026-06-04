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
