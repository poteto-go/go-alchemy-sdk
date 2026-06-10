package ether_test

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEther_Snapshot(t *testing.T) {
	t.Run("normal case:", func(t *testing.T) {
		t.Run("return snapshot id", func(t *testing.T) {
			// Arrange
			e := newEtherApiForTest()
			alchemyMock := newAlchemyMockOnEtherTest(t)
			defer alchemyMock.DeactivateAndReset()

			// Mock
			alchemyMock.RegisterResponderOnce(
				"evm_snapshot",
				`{"jsonrpc":"2.0","id":1,"result":"0x1"}`,
			)

			// Act
			snapshotId, err := e.Snapshot()

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, big.NewInt(1), snapshotId)
		})
	})

	t.Run("error case:", func(t *testing.T) {
		t.Run("if rpc call fails, return error", func(t *testing.T) {
			// Arrange
			e := newEtherApiForTest()
			alchemyMock := newAlchemyMockOnEtherTest(t)
			defer alchemyMock.DeactivateAndReset()

			// Mock
			alchemyMock.RegisterResponderOnce(
				"evm_snapshot",
				`{"jsonrpc":"2.0","id":1,"error":{"code":-32601,"message":"method not found"}}`,
			)

			// Act
			_, err := e.Snapshot()

			// Assert
			assert.Error(t, err)
		})

		t.Run("if result is not hex string, return error", func(t *testing.T) {
			// Arrange
			e := newEtherApiForTest()
			alchemyMock := newAlchemyMockOnEtherTest(t)
			defer alchemyMock.DeactivateAndReset()

			// Mock
			alchemyMock.RegisterResponderOnce(
				"evm_snapshot",
				`{"jsonrpc":"2.0","id":1,"result":12}`,
			)

			// Act
			_, err := e.Snapshot()

			// Assert
			assert.Error(t, err)
		})
	})
}

func TestEther_RevertTo(t *testing.T) {
	t.Run("normal case:", func(t *testing.T) {
		t.Run("return true if reverted", func(t *testing.T) {
			// Arrange
			e := newEtherApiForTest()
			alchemyMock := newAlchemyMockOnEtherTest(t)
			defer alchemyMock.DeactivateAndReset()

			// Mock
			alchemyMock.RegisterResponderOnce(
				"evm_revert",
				`{"jsonrpc":"2.0","id":1,"result":true}`,
			)

			// Act
			reverted, err := e.RevertTo(big.NewInt(1))

			// Assert
			assert.NoError(t, err)
			assert.True(t, reverted)
		})

		t.Run("return false if snapshot does not exist", func(t *testing.T) {
			// Arrange
			e := newEtherApiForTest()
			alchemyMock := newAlchemyMockOnEtherTest(t)
			defer alchemyMock.DeactivateAndReset()

			// Mock
			alchemyMock.RegisterResponderOnce(
				"evm_revert",
				`{"jsonrpc":"2.0","id":1,"result":false}`,
			)

			// Act
			reverted, err := e.RevertTo(big.NewInt(99))

			// Assert
			assert.NoError(t, err)
			assert.False(t, reverted)
		})
	})

	t.Run("error case:", func(t *testing.T) {
		t.Run("if snapshot id is nil, return error", func(t *testing.T) {
			// Arrange
			e := newEtherApiForTest()

			// Act
			_, err := e.RevertTo(nil)

			// Assert
			assert.Error(t, err)
		})

		t.Run("if rpc call fails, return error", func(t *testing.T) {
			// Arrange
			e := newEtherApiForTest()
			alchemyMock := newAlchemyMockOnEtherTest(t)
			defer alchemyMock.DeactivateAndReset()

			// Mock
			alchemyMock.RegisterResponderOnce(
				"evm_revert",
				`{"jsonrpc":"2.0","id":1,"error":{"code":-32601,"message":"method not found"}}`,
			)

			// Act
			_, err := e.RevertTo(big.NewInt(1))

			// Assert
			assert.Error(t, err)
		})

		t.Run("if result is not bool, return error", func(t *testing.T) {
			// Arrange
			e := newEtherApiForTest()
			alchemyMock := newAlchemyMockOnEtherTest(t)
			defer alchemyMock.DeactivateAndReset()

			// Mock
			alchemyMock.RegisterResponderOnce(
				"evm_revert",
				`{"jsonrpc":"2.0","id":1,"result":"yes"}`,
			)

			// Act
			_, err := e.RevertTo(big.NewInt(1))

			// Assert
			assert.Error(t, err)
		})
	})
}
