package ether_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"math/big"

	"github.com/agiledragon/gomonkey"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	eth "github.com/poteto-go/go-alchemy-sdk/ether"
	"github.com/stretchr/testify/assert"
)

func TestEther_ContractCall(t *testing.T) {
	t.Run("normal case", func(t *testing.T) {
		// Arrange
		ether := newEtherApiForTest()
		alchemyMock := newAlchemyMockOnEtherTest(t)
		defer alchemyMock.DeactivateAndReset()

		addr := common.HexToAddress("0x123")

		expectedVal := big.NewInt(1)
		expectedHex := "0x0000000000000000000000000000000000000000000000000000000000000001"
		alchemyMock.RegisterResponderOnce("eth_call", `{"jsonrpc":"2.0","id":1,"result":"`+expectedHex+`"}`)

		unpack := func(b []byte) (any, error) {
			return expectedVal, nil
		}

		callData := []byte{0x12, 0x34}

		// Act
		res, err := ether.ContractCall(addr, nil, callData, unpack)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedVal, res)
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
				"SetEthClient",
				func(_ *eth.Ether) error {
					return errors.New("error")
				},
			)

			// Mock objects
			addr := common.HexToAddress("0x123")
			callData := []byte{}
			unpack := func(b []byte) (any, error) { return nil, nil }

			// Act
			res, err := ether.ContractCall(addr, nil, callData, unpack)

			// Assert
			assert.Error(t, err)
			assert.Nil(t, res)
		})

		t.Run("if error occur on call contract, return error", func(t *testing.T) {
			// Arrange
			ether := newEtherApiForTest()
			alchemyMock := newAlchemyMockOnEtherTest(t)
			defer alchemyMock.DeactivateAndReset()

			// Mock: Respond with an error
			alchemyMock.RegisterResponderOnce("eth_call", `{"jsonrpc":"2.0","id":1,"error":{"code":-32000,"message":"execution reverted"}}`)

			// Mock objects
			addr := common.HexToAddress("0x123")
			callData := []byte{0x12, 0x34}
			unpack := func(b []byte) (any, error) { return nil, nil }

			// Act
			res, err := ether.ContractCall(addr, nil, callData, unpack)

			// Assert
			assert.Error(t, err)
			assert.Nil(t, res)
		})
	})
}

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
		receipt, err := ether.WaitMined(context.Background(), txHash)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, receipt, expectedReceipt)
	})

	t.Run("ctx cancelled, return context.Canceled error", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		ether := newEtherApiForTest()
		ctx, cancel := context.WithCancel(context.Background())

		// Mock: block until ctx is done
		patches.ApplyFunc(
			bind.WaitMined,
			func(ctx context.Context, b bind.DeployBackend, hash common.Hash) (*gethTypes.Receipt, error) {
				<-ctx.Done()
				return nil, ctx.Err()
			},
		)

		cancel()

		// Act
		receipt, err := ether.WaitMined(ctx, txHash)

		// Assert
		assert.ErrorIs(t, err, context.Canceled)
		assert.Nil(t, receipt)
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
				"SetEthClient",
				func(_ *eth.Ether) error {
					return errors.New("error")
				},
			)

			// Act
			receipt, err := ether.WaitMined(context.Background(), txHash)

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
			receipt, err := ether.WaitMined(context.Background(), txHash)

			// Assert
			assert.Error(t, err)
			assert.Nil(t, receipt)
		})
	})

	t.Run("simulated backend mines the pending tx and returns its receipt", func(t *testing.T) {
		// Arrange
		ether, cleanup := newSimulatedEtherForTest(t)
		defer cleanup()
		hash := simSendValueTx(t, ether)

		// Act
		receipt, err := ether.WaitMined(context.Background(), hash)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, uint64(1), receipt.Status)
		assert.Equal(t, hash, receipt.TxHash)
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
		address, err := ether.WaitDeployed(context.Background(), txHash)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, address, expectedAddress)
	})

	t.Run("ctx cancelled, return context.Canceled error", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		ether := newEtherApiForTest()
		ctx, cancel := context.WithCancel(context.Background())

		// Mock: block until ctx is done
		patches.ApplyFunc(
			bind.WaitDeployed,
			func(ctx context.Context, b bind.DeployBackend, hash common.Hash) (common.Address, error) {
				<-ctx.Done()
				return common.Address{}, ctx.Err()
			},
		)

		cancel()

		// Act
		address, err := ether.WaitDeployed(ctx, txHash)

		// Assert
		assert.ErrorIs(t, err, context.Canceled)
		assert.Equal(t, address, common.Address{})
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
				"SetEthClient",
				func(_ *eth.Ether) error {
					return errors.New("error")
				},
			)

			// Act
			address, err := ether.WaitDeployed(context.Background(), txHash)

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
			address, err := ether.WaitDeployed(context.Background(), txHash)

			// Assert
			assert.Error(t, err)
			assert.Equal(t, address, common.Address{})
		})
	})
}
