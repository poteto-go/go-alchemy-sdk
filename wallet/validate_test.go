package wallet

import (
	"math/big"
	"testing"

	"github.com/poteto-go/go-alchemy-sdk/constant"
	"github.com/stretchr/testify/assert"
)

func TestValidateUint256(t *testing.T) {
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
			assert.ErrorIs(t, validateUint256(tt.v), tt.wantErr)
		})
	}
}

func TestValidateAddress(t *testing.T) {
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
			assert.ErrorIs(t, validateAddress(tt.addr), tt.wantErr)
		})
	}
}

func TestWallet_ERC20TransferNoWait_Validation(t *testing.T) {
	contractAddress := "0x1234567890123456789012345678901234567890"
	validAddr := "0xE25583099BA105D9ec0A67f5Ae86D90e50036425"

	t.Run("nil amount returns ErrNilAmount", func(t *testing.T) {
		w, _ := New(testPrivHex)
		_, err := w.ERC20().TransferNoWait(contractAddress, validAddr, nil, nil)
		assert.ErrorIs(t, err, constant.ErrNilAmount)
	})

	t.Run("negative amount returns ErrNegativeAmount", func(t *testing.T) {
		w, _ := New(testPrivHex)
		_, err := w.ERC20().TransferNoWait(contractAddress, validAddr, big.NewInt(-1), nil)
		assert.ErrorIs(t, err, constant.ErrNegativeAmount)
	})

	t.Run("overflow amount returns ErrAmountExceedsUint256", func(t *testing.T) {
		w, _ := New(testPrivHex)
		over := new(big.Int).Lsh(big.NewInt(1), 256)
		_, err := w.ERC20().TransferNoWait(contractAddress, validAddr, over, nil)
		assert.ErrorIs(t, err, constant.ErrAmountExceedsUint256)
	})

	t.Run("invalid address returns ErrInvalidAddress", func(t *testing.T) {
		w, _ := New(testPrivHex)
		_, err := w.ERC20().TransferNoWait(contractAddress, "invalid", big.NewInt(1), nil)
		assert.ErrorIs(t, err, constant.ErrInvalidAddress)
	})
}

func TestWallet_StableCoin_MintNoWait_Validation(t *testing.T) {
	contractAddress := "0x1234567890123456789012345678901234567890"
	validAddr := "0xE25583099BA105D9ec0A67f5Ae86D90e50036425"

	t.Run("nil amount returns ErrNilAmount", func(t *testing.T) {
		w, _ := New(testPrivHex)
		_, err := w.StableCoin().MintNoWait(contractAddress, validAddr, nil, nil)
		assert.ErrorIs(t, err, constant.ErrNilAmount)
	})

	t.Run("negative amount returns ErrNegativeAmount", func(t *testing.T) {
		w, _ := New(testPrivHex)
		_, err := w.StableCoin().MintNoWait(contractAddress, validAddr, big.NewInt(-1), nil)
		assert.ErrorIs(t, err, constant.ErrNegativeAmount)
	})

	t.Run("overflow amount returns ErrAmountExceedsUint256", func(t *testing.T) {
		w, _ := New(testPrivHex)
		over := new(big.Int).Lsh(big.NewInt(1), 256)
		_, err := w.StableCoin().MintNoWait(contractAddress, validAddr, over, nil)
		assert.ErrorIs(t, err, constant.ErrAmountExceedsUint256)
	})

	t.Run("invalid to-address returns ErrInvalidAddress", func(t *testing.T) {
		w, _ := New(testPrivHex)
		_, err := w.StableCoin().MintNoWait(contractAddress, "invalid", big.NewInt(100), nil)
		assert.ErrorIs(t, err, constant.ErrInvalidAddress)
	})
}

func TestWallet_StableCoin_BurnNoWait_Validation(t *testing.T) {
	contractAddress := "0x1234567890123456789012345678901234567890"

	t.Run("nil amount returns ErrNilAmount", func(t *testing.T) {
		w, _ := New(testPrivHex)
		_, err := w.StableCoin().BurnNoWait(contractAddress, nil, nil)
		assert.ErrorIs(t, err, constant.ErrNilAmount)
	})

	t.Run("negative amount returns ErrNegativeAmount", func(t *testing.T) {
		w, _ := New(testPrivHex)
		_, err := w.StableCoin().BurnNoWait(contractAddress, big.NewInt(-5), nil)
		assert.ErrorIs(t, err, constant.ErrNegativeAmount)
	})
}

func TestWallet_StableCoin_ConfigureMinterNoWait_Validation(t *testing.T) {
	contractAddress := "0x1234567890123456789012345678901234567890"
	validAddr := "0xE25583099BA105D9ec0A67f5Ae86D90e50036425"

	t.Run("nil allowance returns ErrNilAmount", func(t *testing.T) {
		w, _ := New(testPrivHex)
		_, err := w.StableCoin().ConfigureMinterNoWait(contractAddress, validAddr, nil, nil)
		assert.ErrorIs(t, err, constant.ErrNilAmount)
	})

	t.Run("invalid minter address returns ErrInvalidAddress", func(t *testing.T) {
		w, _ := New(testPrivHex)
		_, err := w.StableCoin().ConfigureMinterNoWait(contractAddress, "bad", big.NewInt(1), nil)
		assert.ErrorIs(t, err, constant.ErrInvalidAddress)
	})
}
