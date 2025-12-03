package ether_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	eth "github.com/poteto-go/go-alchemy-sdk/ether"
	"github.com/stretchr/testify/assert"
)

func TestEther_WaitMined(t *testing.T) {
	txHash := common.HexToHash("0x123")

	t.Run("normal case", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		ether := newEtherApiForTest()
		expectedReceipt := &gethTypes.Receipt{}

		// Mock
		patches.ApplyFunc(
			bind.WaitMined,
			func(ctx context.Context, b bind.DeployBackend, hash common.Hash) (*gethTypes.Receipt, error) {
				return expectedReceipt, nil
			},
		)

		// Act
		receipt, err := ether.WaitMined(txHash)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, receipt, expectedReceipt)
	})

	t.Run("error case", func(t *testing.T) {
		t.Run("if failed to create connection, return error", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Arrange
			ether := newEtherApiForTest()

			// Mock
			patches.ApplyMethod(
				reflect.TypeOf(ether),
				"GetEthClient",
				func(_ *eth.Ether) (*ethclient.Client, error) {
					return nil, errors.New("error")
				},
			)

			// Act
			receipt, err := ether.WaitMined(txHash)

			// Assert
			assert.Error(t, err)
			assert.Nil(t, receipt)
		})

		t.Run("if error occur on wait mined, return error", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Arrange
			ether := newEtherApiForTest()

			// Mock
			patches.ApplyFunc(
				bind.WaitMined,
				func(ctx context.Context, b bind.DeployBackend, hash common.Hash) (*gethTypes.Receipt, error) {
					return nil, errors.New("error")
				},
			)

			// Act
			receipt, err := ether.WaitMined(txHash)

			// Assert
			assert.Error(t, err)
			assert.Nil(t, receipt)
		})
	})
}

func TestEther_WaitDeployed(t *testing.T) {
	txHash := common.HexToHash("0x123")

	t.Run("normal case", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		ether := newEtherApiForTest()
		expectedAddress := common.HexToAddress("0xabc")

		// Mock
		patches.ApplyFunc(
			bind.WaitDeployed,
			func(ctx context.Context, b bind.DeployBackend, hash common.Hash) (common.Address, error) {
				return expectedAddress, nil
			},
		)

		// Act
		address, err := ether.WaitDeployed(txHash)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, address, expectedAddress)
	})

	t.Run("error case", func(t *testing.T) {
		t.Run("if failed to create connection, return error", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Arrange
			ether := newEtherApiForTest()

			// Mock
			patches.ApplyMethod(
				reflect.TypeOf(ether),
				"GetEthClient",
				func(_ *eth.Ether) (*ethclient.Client, error) {
					return nil, errors.New("error")
				},
			)

			// Act
			address, err := ether.WaitDeployed(txHash)

			// Assert
			assert.Error(t, err)
			assert.Equal(t, address, common.Address{})
		})

		t.Run("if error occur on wait deployed, return error", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Arrange
			ether := newEtherApiForTest()

			// Mock
			patches.ApplyFunc(
				bind.WaitDeployed,
				func(ctx context.Context, b bind.DeployBackend, hash common.Hash) (common.Address, error) {
					return common.Address{}, errors.New("error")
				},
			)

			// Act
			address, err := ether.WaitDeployed(txHash)

			// Assert
			assert.Error(t, err)
			assert.Equal(t, address, common.Address{})
		})
	})
}
