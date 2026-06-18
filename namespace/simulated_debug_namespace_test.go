package namespace_test

import (
	"errors"
	"math/big"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey"
	"github.com/ethereum/go-ethereum/common"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/poteto-go/go-alchemy-sdk/constant"
	"github.com/poteto-go/go-alchemy-sdk/ether"
	"github.com/poteto-go/go-alchemy-sdk/namespace"
	"github.com/stretchr/testify/assert"
)

func newTestBlock() *gethTypes.Block {
	return gethTypes.NewBlockWithHeader(&gethTypes.Header{Number: big.NewInt(1)})
}

func TestNewSimulatedDebugNamespace(t *testing.T) {
	api := newEtherApi()
	debug := namespace.NewSimulatedDebugNamespace(api)
	assert.NotNil(t, debug)
}

func TestSimulatedDebug_Snapshot(t *testing.T) {
	t.Run("first snapshot returns id 0", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		api := newEtherApi()
		simDebug := namespace.NewSimulatedDebugNamespace(api).(*namespace.SimulatedDebug)

		block := newTestBlock()
		patches.ApplyMethod(
			reflect.TypeOf(api),
			"GetBlockByNumber",
			func(_ *ether.Ether, _ string) (*gethTypes.Block, error) {
				return block, nil
			},
		)

		id, err := simDebug.Snapshot()

		assert.NoError(t, err)
		assert.Equal(t, int64(0), id.Int64())
	})

	t.Run("successive snapshots return incrementing ids", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		api := newEtherApi()
		simDebug := namespace.NewSimulatedDebugNamespace(api).(*namespace.SimulatedDebug)

		block := newTestBlock()
		patches.ApplyMethod(
			reflect.TypeOf(api),
			"GetBlockByNumber",
			func(_ *ether.Ether, _ string) (*gethTypes.Block, error) {
				return block, nil
			},
		)

		id0, _ := simDebug.Snapshot()
		id1, _ := simDebug.Snapshot()

		assert.Equal(t, int64(0), id0.Int64())
		assert.Equal(t, int64(1), id1.Int64())
	})

	t.Run("does not call Commit (block number stays the same)", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		api := newEtherApi()
		simDebug := namespace.NewSimulatedDebugNamespace(api).(*namespace.SimulatedDebug)

		block := newTestBlock()
		patches.ApplyMethod(
			reflect.TypeOf(api),
			"GetBlockByNumber",
			func(_ *ether.Ether, blockNumber string) (*gethTypes.Block, error) {
				assert.Equal(t, "latest", blockNumber)
				return block, nil
			},
		)

		_, err := simDebug.Snapshot()

		assert.NoError(t, err)
	})

	t.Run("if GetBlockByNumber fails, return error", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		api := newEtherApi()
		simDebug := namespace.NewSimulatedDebugNamespace(api).(*namespace.SimulatedDebug)

		expectedErr := errors.New("get block error")
		patches.ApplyMethod(
			reflect.TypeOf(api),
			"GetBlockByNumber",
			func(_ *ether.Ether, _ string) (*gethTypes.Block, error) {
				return nil, expectedErr
			},
		)

		_, err := simDebug.Snapshot()

		assert.ErrorIs(t, err, expectedErr)
	})
}

func TestSimulatedDebug_RevertTo(t *testing.T) {
	t.Run("calls Fork with the correct hash and returns true", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		api := newEtherApi()
		simDebug := namespace.NewSimulatedDebugNamespace(api).(*namespace.SimulatedDebug)

		expectedBlock := newTestBlock()
		expectedHash := expectedBlock.Hash()

		patches.ApplyMethod(
			reflect.TypeOf(api),
			"GetBlockByNumber",
			func(_ *ether.Ether, _ string) (*gethTypes.Block, error) {
				return expectedBlock, nil
			},
		)

		id, _ := simDebug.Snapshot()
		patches.Reset()

		patches = gomonkey.NewPatches()
		patches.ApplyMethod(
			reflect.TypeOf(api),
			"Fork",
			func(_ *ether.Ether, h common.Hash) error {
				assert.Equal(t, expectedHash, h)
				return nil
			},
		)

		reverted, err := simDebug.RevertTo(id)

		assert.NoError(t, err)
		assert.True(t, reverted)
	})

	t.Run("returns ErrUnexpectedSnapshotId for unknown id", func(t *testing.T) {
		api := newEtherApi()
		simDebug := namespace.NewSimulatedDebugNamespace(api).(*namespace.SimulatedDebug)

		_, err := simDebug.RevertTo(big.NewInt(999))

		assert.ErrorIs(t, err, constant.ErrUnexpectedSnapshotId)
	})

	t.Run("if Fork fails, return error", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		api := newEtherApi()
		simDebug := namespace.NewSimulatedDebugNamespace(api).(*namespace.SimulatedDebug)

		block := newTestBlock()
		patches.ApplyMethod(
			reflect.TypeOf(api),
			"GetBlockByNumber",
			func(_ *ether.Ether, _ string) (*gethTypes.Block, error) {
				return block, nil
			},
		)

		id, _ := simDebug.Snapshot()
		patches.Reset()

		patches = gomonkey.NewPatches()
		expectedErr := errors.New("fork error")
		patches.ApplyMethod(
			reflect.TypeOf(api),
			"Fork",
			func(_ *ether.Ether, _ common.Hash) error {
				return expectedErr
			},
		)

		_, err := simDebug.RevertTo(id)

		assert.ErrorIs(t, err, expectedErr)
	})
}
