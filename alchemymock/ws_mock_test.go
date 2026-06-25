package alchemymock_test

import (
	"context"
	"math/big"
	"strings"
	"testing"
	"time"

	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/poteto-go/go-alchemy-sdk/alchemymock"
	"github.com/poteto-go/go-alchemy-sdk/gas"
	"github.com/poteto-go/go-alchemy-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newWsAlchemyForMock builds a ws Alchemy whose geth client dials the in-process
// mock. A ws-scheme private-network url selects the WsAlchemyProvider, so the
// subscribe path is exercised exactly as a real ws Alchemy would.
func newWsAlchemyForMock(t *testing.T, url string) gas.Alchemy {
	t.Helper()
	a, err := gas.NewAlchemy(gas.AlchemySetting{
		PrivateNetworkConfig: gas.PrivateNetworkConfig{Url: url},
	})
	require.NoError(t, err)
	return a
}

func TestNewAlchemyWsMock(t *testing.T) {
	t.Run("URL returns a ws:// endpoint", func(t *testing.T) {
		mock := alchemymock.NewAlchemyWsMock(t)
		defer mock.Close()

		assert.True(t, strings.HasPrefix(mock.URL(), "ws://"), "url should be ws://: %s", mock.URL())
	})

	t.Run("Close is safe to call twice", func(t *testing.T) {
		mock := alchemymock.NewAlchemyWsMock(t)
		mock.Close()
		assert.NotPanics(t, func() { mock.Close() })
	})
}

func TestAlchemyWsMock_EmitNewHeads(t *testing.T) {
	t.Run("emitted headers are pushed to a newHeads subscriber", func(t *testing.T) {
		// Arrange
		mock := alchemymock.NewAlchemyWsMock(t)
		defer mock.Close()

		wsAlchemy := newWsAlchemyForMock(t, mock.URL())
		defer wsAlchemy.GetProvider().Eth().Shutdown()

		sub, ok := wsAlchemy.GetProvider().(types.ISubscribeProvider)
		require.True(t, ok, "ws provider must implement ISubscribeProvider")

		ch := make(chan *gethTypes.Header, 4)
		subscription, err := sub.Subscribe(context.Background(), ch, "newHeads")
		require.NoError(t, err)
		defer subscription.Unsubscribe()

		// Act: push two canned heads after the subscription is established.
		// newHeads notifications carry full headers; geth's Header.UnmarshalJSON
		// requires difficulty (and number), so both are set.
		mock.EmitNewHeads(
			&gethTypes.Header{Number: big.NewInt(0x10), Difficulty: big.NewInt(0)},
			&gethTypes.Header{Number: big.NewInt(0x11), Difficulty: big.NewInt(0)},
		)

		// Assert: they arrive on the subscriber channel in order.
		for _, want := range []int64{0x10, 0x11} {
			select {
			case got := <-ch:
				assert.Equal(t, want, got.Number.Int64())
			case err := <-subscription.Err():
				t.Fatalf("subscription errored: %v", err)
			case <-time.After(2 * time.Second):
				t.Fatalf("timed out waiting for head 0x%x", want)
			}
		}
	})

	t.Run("EmitNewHeads with no subscribers is a no-op", func(t *testing.T) {
		mock := alchemymock.NewAlchemyWsMock(t)
		defer mock.Close()

		assert.NotPanics(t, func() {
			mock.EmitNewHeads(&gethTypes.Header{Number: big.NewInt(1)})
		})
	})
}
