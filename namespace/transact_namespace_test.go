package namespace_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey"
	"github.com/ethereum/go-ethereum/common"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/poteto-go/go-alchemy-sdk/ether"
	"github.com/poteto-go/go-alchemy-sdk/namespace"
	"github.com/stretchr/testify/assert"
)

func TestNewTransactNamespace(t *testing.T) {
	// Arrange
	ether := newEtherApi()

	// Act
	core := namespace.NewTransactNamespace(ether)

	// Assert
	assert.NotNil(t, core)
}

func Test_WaitMined(t *testing.T) {
	// Arrange
	api := newEtherApi()
	transact := namespace.NewTransactNamespace(api).(*namespace.Transact)

	t.Run("call ether.WaitMined & return result", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		hexHash := "0x123"
		hash := common.HexToHash(hexHash)

		// Mock & Assert
		patches.ApplyMethod(
			reflect.TypeOf(api),
			"WaitMined",
			func(_ *ether.Ether, txHash common.Hash) (*gethTypes.Receipt, error) {
				assert.Equal(t, txHash, hash)
				return &gethTypes.Receipt{}, nil
			},
		)

		// Act
		_, err := transact.WaitMined(hexHash)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("if error occur, return error", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		hexHash := "0x123"
		hash := common.HexToHash(hexHash)
		expectedErr := errors.New("error")

		// Mock & Assert
		patches.ApplyMethod(
			reflect.TypeOf(api),
			"WaitMined",
			func(_ *ether.Ether, txHash common.Hash) (*gethTypes.Receipt, error) {
				assert.Equal(t, txHash, hash)
				return &gethTypes.Receipt{}, expectedErr
			},
		)

		// Act
		_, err := transact.WaitMined(hexHash)

		// Assert
		assert.ErrorIs(t, err, expectedErr)
	})
}

func Test_WaitDeployed(t *testing.T) {
	// Arrange
	api := newEtherApi()
	transact := namespace.NewTransactNamespace(api).(*namespace.Transact)

	t.Run("call ether.WaitDeployed & return result", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		hexHash := "0x123"
		hash := common.HexToHash(hexHash)

		// Mock
		patches.ApplyMethod(
			reflect.TypeOf(api),
			"WaitDeployed",
			func(_ *ether.Ether, txHash common.Hash) (common.Address, error) {
				assert.Equal(t, txHash, hash)
				return common.Address{}, nil
			},
		)

		// Act
		_, err := transact.WaitDeployed(hexHash)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("if error occur, return error", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		hexHash := "0x123"
		hash := common.HexToHash(hexHash)
		expectedErr := errors.New("error")

		// Mock & Assert
		patches.ApplyMethod(
			reflect.TypeOf(api),
			"WaitDeployed",
			func(_ *ether.Ether, txHash common.Hash) (common.Address, error) {
				assert.Equal(t, txHash, hash)
				return common.Address{}, expectedErr
			},
		)

		// Act
		_, err := transact.WaitDeployed(hexHash)

		// Assert
		assert.ErrorIs(t, err, expectedErr)
	})
}
