package namespace_test

import (
	"errors"
	"math/big"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey"
	"github.com/poteto-go/go-alchemy-sdk/ether"
	"github.com/poteto-go/go-alchemy-sdk/namespace"
	"github.com/stretchr/testify/assert"
)

func TestNewDebugNamespace(t *testing.T) {
	// Arrange
	ether := newEtherApi()

	// Act
	debug := namespace.NewDebugNamespace(ether)

	// Assert
	assert.NotNil(t, debug)
}

func TestDebug_Snapshot(t *testing.T) {
	// Arrange
	api := newEtherApi()
	debug := namespace.NewDebugNamespace(api).(*namespace.Debug)

	t.Run("call ether.Snapshot & return snapshot id", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		expectedId := big.NewInt(1)

		// Mock
		patches.ApplyMethod(
			reflect.TypeOf(api),
			"Snapshot",
			func(_ *ether.Ether) (*big.Int, error) {
				return expectedId, nil
			},
		)

		// Act
		snapshotId, err := debug.Snapshot()

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedId, snapshotId)
	})

	t.Run("if error occur, return error", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		expectedErr := errors.New("error")

		// Mock
		patches.ApplyMethod(
			reflect.TypeOf(api),
			"Snapshot",
			func(_ *ether.Ether) (*big.Int, error) {
				return nil, expectedErr
			},
		)

		// Act
		_, err := debug.Snapshot()

		// Assert
		assert.ErrorIs(t, err, expectedErr)
	})
}

func TestDebug_RevertTo(t *testing.T) {
	// Arrange
	api := newEtherApi()
	debug := namespace.NewDebugNamespace(api).(*namespace.Debug)

	t.Run("call ether.RevertTo & return result", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		expectedId := big.NewInt(1)

		// Mock & Assert
		patches.ApplyMethod(
			reflect.TypeOf(api),
			"RevertTo",
			func(_ *ether.Ether, snapshotId *big.Int) (bool, error) {
				assert.Equal(t, expectedId, snapshotId)
				return true, nil
			},
		)

		// Act
		reverted, err := debug.RevertTo(expectedId)

		// Assert
		assert.NoError(t, err)
		assert.True(t, reverted)
	})

	t.Run("if error occur, return error", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		expectedErr := errors.New("error")

		// Mock
		patches.ApplyMethod(
			reflect.TypeOf(api),
			"RevertTo",
			func(_ *ether.Ether, snapshotId *big.Int) (bool, error) {
				return false, expectedErr
			},
		)

		// Act
		_, err := debug.RevertTo(big.NewInt(1))

		// Assert
		assert.ErrorIs(t, err, expectedErr)
	})
}
