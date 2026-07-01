package alchemymock_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/poteto-go/go-alchemy-sdk/alchemymock"
	"github.com/poteto-go/go-alchemy-sdk/gas"
	"github.com/poteto-go/go-alchemy-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newMinimalHeader returns a gethTypes.Header with all JSON-required fields
// populated so it survives a Header.UnmarshalJSON round-trip.
func newMinimalHeader(number int64) *gethTypes.Header {
	return &gethTypes.Header{
		Number:      big.NewInt(number),
		Difficulty:  big.NewInt(0),
		Extra:       []byte{},
		UncleHash:   gethTypes.EmptyUncleHash,
		TxHash:      gethTypes.EmptyTxsHash,
		ReceiptHash: gethTypes.EmptyReceiptsHash,
	}
}

var wsMockSetting = gas.AlchemySetting{
	BackoffConfig: &types.BackoffConfig{MaxRetries: 0},
}

func TestNewAlchemyWsMock(t *testing.T) {
	mock := alchemymock.NewAlchemyWsMock(wsMockSetting, t)
	assert.NotEmpty(t, mock.URL())
}

func TestAlchemyWsMock_RegisterResponderOnce(t *testing.T) {
	t.Run("serves a regular call over the ws socket", func(t *testing.T) {
		mock := alchemymock.NewAlchemyWsMock(wsMockSetting, t)
		a, err := mock.NewAlchemy()
		require.NoError(t, err)

		mock.RegisterResponderOnce("eth_blockNumber", `{"jsonrpc":"2.0","id":1,"result":"0x10"}`)

		bn, err := a.Core.GetBlockNumber()

		require.NoError(t, err)
		assert.Equal(t, uint64(16), bn)
	})

	t.Run("returns error when method is not mocked", func(t *testing.T) {
		mock := alchemymock.NewAlchemyWsMock(wsMockSetting, t)
		a, err := mock.NewAlchemy()
		require.NoError(t, err)

		// No responder registered -> mock returns -32601.
		_, err = a.Core.GetBlockNumber()
		assert.Error(t, err)
	})

	t.Run("serves responses in FIFO order for the same method", func(t *testing.T) {
		mock := alchemymock.NewAlchemyWsMock(wsMockSetting, t)
		a, err := mock.NewAlchemy()
		require.NoError(t, err)

		mock.RegisterResponderOnce("eth_blockNumber", `{"jsonrpc":"2.0","id":1,"result":"0x1"}`)
		mock.RegisterResponderOnce("eth_blockNumber", `{"jsonrpc":"2.0","id":1,"result":"0x2"}`)

		bn1, err1 := a.Core.GetBlockNumber()
		require.NoError(t, err1)

		bn2, err2 := a.Core.GetBlockNumber()
		require.NoError(t, err2)

		assert.Equal(t, uint64(1), bn1)
		assert.Equal(t, uint64(2), bn2)
	})

	t.Run("multiple different methods can be mocked independently", func(t *testing.T) {
		mock := alchemymock.NewAlchemyWsMock(wsMockSetting, t)
		a, err := mock.NewAlchemy()
		require.NoError(t, err)

		mock.RegisterResponderOnce("eth_blockNumber", `{"jsonrpc":"2.0","id":1,"result":"0x42"}`)
		mock.RegisterResponderOnce("eth_gasPrice", `{"jsonrpc":"2.0","id":1,"result":"0x3b9aca00"}`)

		bn, err := a.Core.GetBlockNumber()
		require.NoError(t, err)
		assert.Equal(t, uint64(0x42), bn)

		gp, err := a.Core.GetGasPrice()
		require.NoError(t, err)
		assert.Equal(t, big.NewInt(1_000_000_000), gp)
	})
}

func TestAlchemyWsMock_EmitNewHeads(t *testing.T) {
	mock := alchemymock.NewAlchemyWsMock(wsMockSetting, t)
	a, err := mock.NewAlchemy()
	require.NoError(t, err)

	ch := make(chan *gethTypes.Header, 4)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	subscription, err := a.WS.Subscribe(ctx, ch, "newHeads")
	require.NoError(t, err)
	defer subscription.Unsubscribe()

	header := newMinimalHeader(0x42)
	mock.EmitNewHeads(header)

	select {
	case got := <-ch:
		assert.Equal(t, big.NewInt(0x42), got.Number)
	case err := <-subscription.Err():
		t.Fatalf("subscription error: %v", err)
	case <-time.After(3 * time.Second):
		t.Fatal("timed out waiting for newHeads notification")
	}
}

func TestAlchemyWsMock_EmitNewHeads_MultipleHeaders(t *testing.T) {
	mock := alchemymock.NewAlchemyWsMock(wsMockSetting, t)
	a, err := mock.NewAlchemy()
	require.NoError(t, err)

	ch := make(chan *gethTypes.Header, 4)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	subscription, err := a.WS.Subscribe(ctx, ch, "newHeads")
	require.NoError(t, err)
	defer subscription.Unsubscribe()

	h1 := newMinimalHeader(1)
	h2 := newMinimalHeader(2)
	mock.EmitNewHeads(h1, h2)

	for _, wantN := range []*big.Int{big.NewInt(1), big.NewInt(2)} {
		select {
		case got := <-ch:
			assert.Equal(t, wantN, got.Number)
		case err := <-subscription.Err():
			t.Fatalf("subscription error: %v", err)
		case <-time.After(3 * time.Second):
			t.Fatal("timed out waiting for newHeads notification")
		}
	}
}
