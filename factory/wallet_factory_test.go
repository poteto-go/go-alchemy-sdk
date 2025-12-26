package factory

import (
	"errors"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/poteto-go/go-alchemy-sdk/types"
	"github.com/poteto-go/go-alchemy-sdk/wallet"
	"github.com/stretchr/testify/assert"
)

// Mock wallet implementation for test
type mockWallet struct {
	wallet.Wallet
}

func (m *mockWallet) ContractCall(
	contract types.ContractInstance,
	contractAddress string,
	opts *bind.CallOpts,
	callData []byte,
	unpack func([]byte) (any, error),
) (any, error) {
	return unpack([]byte("result"))
}

func TestContractCall(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Arrange
		w := &mockWallet{}
		unpack := func(b []byte) (string, error) {
			return string(b), nil
		}

		// Act
		res, err := ContractCall(w, nil, "0x123", nil, nil, unpack)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "result", res)
	})

	t.Run("error from wallet", func(t *testing.T) {
		// Arrange
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		w := &mockWallet{}
		patches.ApplyMethod(
			reflect.TypeOf(w),
			"ContractCall",
			func(
				_ *mockWallet,
				_ types.ContractInstance,
				_ string,
				_ *bind.CallOpts,
				_ []byte,
				_ func([]byte) (any, error),
			) (any, error) {
				return nil, errors.New("error")
			},
		)

		unpack := func(b []byte) (string, error) {
			return string(b), nil
		}

		// Act
		_, err := ContractCall(w, nil, "0x123", nil, nil, unpack)

		// Assert
		assert.Error(t, err)
	})

	t.Run("type assertion error", func(t *testing.T) {
		// Arrange
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		w := &mockWallet{}
		patches.ApplyMethod(
			reflect.TypeOf(w),
			"ContractCall",
			func(
				_ *mockWallet,
				_ types.ContractInstance,
				_ string,
				_ *bind.CallOpts,
				_ []byte,
				_ func([]byte) (any, error),
			) (any, error) {
				return 123, nil // Return int, expected string
			},
		)

		unpack := func(b []byte) (string, error) {
			return string(b), nil
		}

		// Act
		_, err := ContractCall(w, nil, "0x123", nil, nil, unpack)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to cast result")
	})
}
