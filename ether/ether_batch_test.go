package ether_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/poteto-go/go-alchemy-sdk/constant"
	"github.com/stretchr/testify/assert"
)

func TestEther_BatchCall(t *testing.T) {
	t.Run("executes multiple geth calls in a single batch request", func(t *testing.T) {
		// Arrange
		alchemyMock := newAlchemyMockOnEtherTest(t)
		defer alchemyMock.DeactivateAndReset()

		// ids are assigned sequentially (1, 2, ...) by a fresh geth rpc.Client.
		alchemyMock.RegisterBatchResponderOnce(
			`[{"jsonrpc":"2.0","id":1,"result":"0x1234"},{"jsonrpc":"2.0","id":2,"result":"0x5678"}]`,
		)

		ether := newEtherApiForTest()

		blockNumber := new(string)
		gasPrice := new(string)
		elems := []rpc.BatchElem{
			{Method: "eth_blockNumber", Args: []any{}, Result: blockNumber},
			{Method: "eth_gasPrice", Args: []any{}, Result: gasPrice},
		}

		// Act
		err := ether.BatchCall(elems)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "0x1234", *blockNumber)
		assert.Equal(t, "0x5678", *gasPrice)
		assert.NoError(t, elems[0].Error)
		assert.NoError(t, elems[1].Error)
	})

	t.Run("per-element error is populated in place without failing the call", func(t *testing.T) {
		// Arrange
		alchemyMock := newAlchemyMockOnEtherTest(t)
		defer alchemyMock.DeactivateAndReset()

		alchemyMock.RegisterBatchResponderOnce(
			`[{"jsonrpc":"2.0","id":1,"result":"0x1234"},{"jsonrpc":"2.0","id":2,"error":{"code":-32000,"message":"boom"}}]`,
		)

		ether := newEtherApiForTest()

		blockNumber := new(string)
		gasPrice := new(string)
		elems := []rpc.BatchElem{
			{Method: "eth_blockNumber", Args: []any{}, Result: blockNumber},
			{Method: "eth_gasPrice", Args: []any{}, Result: gasPrice},
		}

		// Act
		err := ether.BatchCall(elems)

		// Assert: I/O succeeded, the failure is attached to the element only.
		assert.NoError(t, err)
		assert.Equal(t, "0x1234", *blockNumber)
		assert.NoError(t, elems[0].Error)
		assert.Error(t, elems[1].Error)
	})

	t.Run("returns error if client cannot be created", func(t *testing.T) {
		// Arrange: empty url makes SetEthClient fail at createRpcClient.
		ether := newNilEtherApiForTest(newProviderForTest())

		// Act
		err := ether.BatchCall([]rpc.BatchElem{
			{Method: "eth_blockNumber", Args: []any{}, Result: new(string)},
		})

		// Assert
		assert.Error(t, err)
	})

	t.Run("simulated backend does not support BatchCall", func(t *testing.T) {
		// Arrange
		ether, cleanup := newSimulatedEtherForTest(t)
		defer cleanup()

		// Act
		err := ether.BatchCall([]rpc.BatchElem{
			{Method: "eth_blockNumber", Args: []any{}, Result: new(string)},
		})

		// Assert
		assert.ErrorIs(t, constant.ErrUnSupportSimulatedMethod, err)
	})
}
