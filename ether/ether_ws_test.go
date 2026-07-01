package ether_test

import (
	"context"
	"encoding/json"
	"math/big"
	"net/http"
	"testing"
	"time"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/poteto-go/go-alchemy-sdk/alchemymock"
	"github.com/poteto-go/go-alchemy-sdk/constant"
	"github.com/poteto-go/go-alchemy-sdk/ether"
	eth "github.com/poteto-go/go-alchemy-sdk/ether"
	"github.com/poteto-go/go-alchemy-sdk/gas"
	"github.com/poteto-go/go-alchemy-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var utWsAlchemySetting = gas.AlchemySetting{
	ApiKey:  "hoge",
	Network: "fuga",
	BackoffConfig: &types.BackoffConfig{
		MaxRetries: 0,
	},
	UseWebsocket: true,
}

// newEtherWsApiForTest builds a ws Ether pointed at the derived Alchemy endpoint.
// That endpoint is unreachable in tests, so use it to exercise dial-failure paths.
func newEtherWsApiForTest() *eth.Ether {
	provider := newProviderForTest()
	config, err := gas.NewAlchemyConfig(utWsAlchemySetting)
	if err != nil {
		panic(err)
	}

	return ether.NewEtherApi(
		provider,
		eth.NewEtherApiConfig(
			config.GetUrl(),
			0,
			time.Duration(1*time.Second),
			nil,
			[]http.Header{},
			[]byte(""),
			5<<20,
			nil,
		),
	).(*eth.Ether)
}

// newEtherWsApiForTestWithUrl builds an Ether dialed at an explicit ws url
// (e.g. the one served by alchemymock.AlchemyWsMock) instead of the derived
// Alchemy endpoint.
func newEtherWsApiForTestWithUrl(wsUrl string) *eth.Ether {
	provider := newProviderForTest()
	return ether.NewEtherApi(
		provider,
		eth.NewEtherApiConfig(
			wsUrl,
			0,
			2*time.Second,
			&types.DefaultBackoffConfig,
			[]http.Header{},
			nil,
			5<<20,
			nil,
		),
	).(*eth.Ether)
}

// newWsMinimalHeader returns a gethTypes.Header with all JSON-required fields
// populated so it survives a Header.UnmarshalJSON round-trip through the mock.
func newWsMinimalHeader(number int64) *gethTypes.Header {
	return &gethTypes.Header{
		Number:      big.NewInt(number),
		Difficulty:  big.NewInt(0),
		Extra:       []byte{},
		UncleHash:   gethTypes.EmptyUncleHash,
		TxHash:      gethTypes.EmptyTxsHash,
		ReceiptHash: gethTypes.EmptyReceiptsHash,
	}
}

// newWsMinimalLog returns a gethTypes.Log with the JSON-required fields
// (address, topics, data) populated so it survives a Log.UnmarshalJSON
// round-trip through the mock.
func newWsMinimalLog(address common.Address) gethTypes.Log {
	return gethTypes.Log{
		Address: address,
		Topics:  []common.Hash{},
		Data:    []byte{},
	}
}

// newWsMinimalReceipt returns a gethTypes.Receipt with the JSON-required fields
// (cumulativeGasUsed, logsBloom, logs, transactionHash, gasUsed) populated so it
// survives a Receipt.UnmarshalJSON round-trip through the mock.
func newWsMinimalReceipt(txHash common.Hash) *gethTypes.Receipt {
	return &gethTypes.Receipt{
		TxHash: txHash,
		Logs:   []*gethTypes.Log{},
	}
}

func Test_EtherWsClientLifeCycle(t *testing.T) {
	// Arrange
	e := newEtherWsApiForTest()

	t.Run("failed create client if not ws server connection", func(t *testing.T) {
		// Act & Assert
		assert.Error(t, e.SetEthClient())

		// do nothing on nil client
		assert.NotPanics(t, func() {
			e.Shutdown()
		})
	})

	t.Run("ws server up", func(t *testing.T) {
		// Arrange: stand up an in-process WebSocket JSON-RPC mock.
		mock := alchemymock.NewAlchemyWsMock(utWsAlchemySetting, t)
		wsEther := newEtherWsApiForTestWithUrl(mock.URL())

		// Act & Assert: dial happens eagerly inside DialOptions, so a successful
		// SetEthClient proves the ws socket is established.
		require.NoError(t, wsEther.SetEthClient())
		require.NotNil(t, wsEther.Client())

		// do nothing on double set
		require.NoError(t, wsEther.SetEthClient())

		// Close is a no-op for ws: the persistent socket survives `defer Close`.
		wsEther.Close()
		require.NotNil(t, wsEther.Client())

		// round-trip a real call over the ws socket.
		mock.RegisterResponderOnce("eth_blockNumber", `{"jsonrpc":"2.0","id":1,"result":"0x42"}`)
		bn, err := wsEther.Client().BlockNumber(context.Background())
		require.NoError(t, err)
		assert.Equal(t, uint64(0x42), bn)

		// Shutdown tears the persistent ws client down explicitly.
		wsEther.Shutdown()
		assert.Nil(t, wsEther.Client())
	})
}

func Test_EtherWsSubscribe(t *testing.T) {
	t.Run("returns error if provider is not websocket", func(t *testing.T) {
		// Arrange: a plain (http) Ether cannot open eth_subscribe streams.
		e := newEtherApiForTest()

		// Act
		_, err := e.Subscribe(context.Background(), make(chan *gethTypes.Header), "newHeads")

		// Assert
		assert.ErrorIs(t, err, constant.ErrUnsupportedNotWebsocketProvider)
	})

	t.Run("returns error over a simulated backend", func(t *testing.T) {
		// Arrange: a simulated backend is dialed with an empty (non-ws) url, so it
		// cannot open eth_subscribe streams.
		e, cleanup := newSimulatedEtherForTest(t)
		defer cleanup()

		// Act
		_, err := e.Subscribe(context.Background(), make(chan *gethTypes.Header), "newHeads")

		// Assert
		assert.ErrorIs(t, err, constant.ErrUnsupportedNotWebsocketProvider)
	})

	t.Run("streams subscription notifications over the ws socket", func(t *testing.T) {
		// Arrange
		mock := alchemymock.NewAlchemyWsMock(utWsAlchemySetting, t)
		wsEther := newEtherWsApiForTestWithUrl(mock.URL())

		ch := make(chan *gethTypes.Header, 1)
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		// Act: open the stream, then push a header from the mock.
		sub, err := wsEther.Subscribe(ctx, ch, "newHeads")
		require.NoError(t, err)
		defer sub.Unsubscribe()

		mock.EmitNewHeads(newWsMinimalHeader(0x42))

		// Assert
		select {
		case got := <-ch:
			assert.Equal(t, big.NewInt(0x42), got.Number)
		case err := <-sub.Err():
			t.Fatalf("subscription errored: %v", err)
		case <-time.After(2 * time.Second):
			t.Fatal("timed out waiting for subscription notification")
		}
	})

	t.Run("propagates a dial error when the ws server is unreachable", func(t *testing.T) {
		// Arrange: nothing listens on port 1, so the ws dial inside SetEthClient fails.
		wsEther := newEtherWsApiForTestWithUrl("ws://127.0.0.1:1")

		// Act
		_, err := wsEther.Subscribe(context.Background(), make(chan *gethTypes.Header), "newHeads")

		// Assert
		assert.Error(t, err)
	})

	t.Run("propagates an EthSubscribe error when the context is cancelled", func(t *testing.T) {
		// Arrange: the socket is live, but the subscribe context is already dead.
		mock := alchemymock.NewAlchemyWsMock(utWsAlchemySetting, t)
		wsEther := newEtherWsApiForTestWithUrl(mock.URL())

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// Act: EthSubscribe observes the cancelled context and returns its error.
		_, err := wsEther.Subscribe(ctx, make(chan *gethTypes.Header), "newHeads")

		// Assert
		assert.ErrorIs(t, err, context.Canceled)
	})
}

func Test_EtherWsSubscribeNewHead(t *testing.T) {
	t.Run("returns error if provider is not websocket", func(t *testing.T) {
		// Arrange: a plain (http) Ether cannot open eth_subscribe streams.
		e := newEtherApiForTest()

		// Act
		_, err := e.SubscribeNewHead(context.Background(), make(chan *gethTypes.Header))

		// Assert
		assert.ErrorIs(t, err, constant.ErrUnsupportedNotWebsocketProvider)
	})

	t.Run("returns error over a simulated backend", func(t *testing.T) {
		// Arrange: a simulated backend is dialed with an empty (non-ws) url, so it
		// cannot open eth_subscribe streams.
		e, cleanup := newSimulatedEtherForTest(t)
		defer cleanup()

		// Act
		_, err := e.SubscribeNewHead(context.Background(), make(chan *gethTypes.Header))

		// Assert
		assert.ErrorIs(t, err, constant.ErrUnsupportedNotWebsocketProvider)
	})

	t.Run("streams new head notifications over the ws socket", func(t *testing.T) {
		// Arrange
		mock := alchemymock.NewAlchemyWsMock(utWsAlchemySetting, t)
		wsEther := newEtherWsApiForTestWithUrl(mock.URL())

		ch := make(chan *gethTypes.Header, 1)
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		// Act: open the stream, then push a header from the mock.
		sub, err := wsEther.SubscribeNewHead(ctx, ch)
		require.NoError(t, err)
		defer sub.Unsubscribe()

		mock.EmitNewHeads(newWsMinimalHeader(0x42))

		// Assert
		select {
		case got := <-ch:
			assert.Equal(t, big.NewInt(0x42), got.Number)
		case err := <-sub.Err():
			t.Fatalf("subscription errored: %v", err)
		case <-time.After(2 * time.Second):
			t.Fatal("timed out waiting for subscription notification")
		}
	})

	t.Run("propagates a dial error when the ws server is unreachable", func(t *testing.T) {
		// Arrange: nothing listens on port 1, so the ws dial inside SetEthClient fails.
		wsEther := newEtherWsApiForTestWithUrl("ws://127.0.0.1:1")

		// Act
		_, err := wsEther.SubscribeNewHead(context.Background(), make(chan *gethTypes.Header))

		// Assert
		assert.Error(t, err)
	})

	t.Run("propagates an EthSubscribe error when the context is cancelled", func(t *testing.T) {
		// Arrange: the socket is live, but the subscribe context is already dead.
		mock := alchemymock.NewAlchemyWsMock(utWsAlchemySetting, t)
		wsEther := newEtherWsApiForTestWithUrl(mock.URL())

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// Act: SubscribeNewHead observes the cancelled context and returns its error.
		_, err := wsEther.SubscribeNewHead(ctx, make(chan *gethTypes.Header))

		// Assert
		assert.ErrorIs(t, err, context.Canceled)
	})
}

func Test_EtherWsSubscribeFilterLogs(t *testing.T) {
	query := ethereum.FilterQuery{}

	t.Run("returns error if provider is not websocket", func(t *testing.T) {
		// Arrange: a plain (http) Ether cannot open eth_subscribe streams.
		e := newEtherApiForTest()

		// Act
		_, err := e.SubscribeFilterLogs(context.Background(), query, make(chan gethTypes.Log))

		// Assert
		assert.ErrorIs(t, err, constant.ErrUnsupportedNotWebsocketProvider)
	})

	t.Run("returns error over a simulated backend", func(t *testing.T) {
		// Arrange: a simulated backend is dialed with an empty (non-ws) url, so it
		// cannot open eth_subscribe streams.
		e, cleanup := newSimulatedEtherForTest(t)
		defer cleanup()

		// Act
		_, err := e.SubscribeFilterLogs(context.Background(), query, make(chan gethTypes.Log))

		// Assert
		assert.ErrorIs(t, err, constant.ErrUnsupportedNotWebsocketProvider)
	})

	t.Run("streams log notifications over the ws socket", func(t *testing.T) {
		// Arrange
		mock := alchemymock.NewAlchemyWsMock(utWsAlchemySetting, t)
		wsEther := newEtherWsApiForTestWithUrl(mock.URL())

		ch := make(chan gethTypes.Log, 1)
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		// Act: open the stream, then push a log from the mock.
		sub, err := wsEther.SubscribeFilterLogs(ctx, query, ch)
		require.NoError(t, err)
		defer sub.Unsubscribe()

		addr := common.HexToAddress("0xdeadbeef00000000000000000000000000000042")
		logData, err := json.Marshal(newWsMinimalLog(addr))
		require.NoError(t, err)
		mock.Emit("logs", logData)

		// Assert
		select {
		case got := <-ch:
			assert.Equal(t, addr, got.Address)
		case err := <-sub.Err():
			t.Fatalf("subscription errored: %v", err)
		case <-time.After(2 * time.Second):
			t.Fatal("timed out waiting for subscription notification")
		}
	})

	t.Run("propagates a dial error when the ws server is unreachable", func(t *testing.T) {
		// Arrange: nothing listens on port 1, so the ws dial inside SetEthClient fails.
		wsEther := newEtherWsApiForTestWithUrl("ws://127.0.0.1:1")

		// Act
		_, err := wsEther.SubscribeFilterLogs(context.Background(), query, make(chan gethTypes.Log))

		// Assert
		assert.Error(t, err)
	})

	t.Run("propagates an EthSubscribe error when the context is cancelled", func(t *testing.T) {
		// Arrange: the socket is live, but the subscribe context is already dead.
		mock := alchemymock.NewAlchemyWsMock(utWsAlchemySetting, t)
		wsEther := newEtherWsApiForTestWithUrl(mock.URL())

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// Act: SubscribeFilterLogs observes the cancelled context and returns its error.
		_, err := wsEther.SubscribeFilterLogs(ctx, query, make(chan gethTypes.Log))

		// Assert
		assert.ErrorIs(t, err, context.Canceled)
	})
}

func Test_EtherWsSubscribeTxReceipts(t *testing.T) {
	q := &ethereum.TransactionReceiptsQuery{}

	t.Run("returns error if provider is not websocket", func(t *testing.T) {
		// Arrange: a plain (http) Ether cannot open eth_subscribe streams.
		e := newEtherApiForTest()

		// Act
		_, err := e.SubscribeTxReceipts(context.Background(), q, make(chan []*gethTypes.Receipt))

		// Assert
		assert.ErrorIs(t, err, constant.ErrUnsupportedNotWebsocketProvider)
	})

	t.Run("returns error over a simulated backend", func(t *testing.T) {
		// Arrange: a simulated backend is dialed with an empty (non-ws) url, so it
		// cannot open eth_subscribe streams.
		e, cleanup := newSimulatedEtherForTest(t)
		defer cleanup()

		// Act
		_, err := e.SubscribeTxReceipts(context.Background(), q, make(chan []*gethTypes.Receipt))

		// Assert
		assert.ErrorIs(t, err, constant.ErrUnsupportedNotWebsocketProvider)
	})

	t.Run("streams tx receipt notifications over the ws socket", func(t *testing.T) {
		// Arrange
		mock := alchemymock.NewAlchemyWsMock(utWsAlchemySetting, t)
		wsEther := newEtherWsApiForTestWithUrl(mock.URL())

		ch := make(chan []*gethTypes.Receipt, 1)
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		// Act: open the stream, then push a batch of receipts from the mock.
		sub, err := wsEther.SubscribeTxReceipts(ctx, q, ch)
		require.NoError(t, err)
		defer sub.Unsubscribe()

		txHash := common.HexToHash("0xdeadbeef00000000000000000000000000000000000000000000000000000042")
		data, err := json.Marshal([]*gethTypes.Receipt{newWsMinimalReceipt(txHash)})
		require.NoError(t, err)
		mock.Emit("transactionReceipts", data)

		// Assert
		select {
		case got := <-ch:
			require.Len(t, got, 1)
			assert.Equal(t, txHash, got[0].TxHash)
		case err := <-sub.Err():
			t.Fatalf("subscription errored: %v", err)
		case <-time.After(2 * time.Second):
			t.Fatal("timed out waiting for subscription notification")
		}
	})

	t.Run("propagates a dial error when the ws server is unreachable", func(t *testing.T) {
		// Arrange: nothing listens on port 1, so the ws dial inside SetEthClient fails.
		wsEther := newEtherWsApiForTestWithUrl("ws://127.0.0.1:1")

		// Act
		_, err := wsEther.SubscribeTxReceipts(context.Background(), q, make(chan []*gethTypes.Receipt))

		// Assert
		assert.Error(t, err)
	})

	t.Run("propagates an EthSubscribe error when the context is cancelled", func(t *testing.T) {
		// Arrange: the socket is live, but the subscribe context is already dead.
		mock := alchemymock.NewAlchemyWsMock(utWsAlchemySetting, t)
		wsEther := newEtherWsApiForTestWithUrl(mock.URL())

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// Act: SubscribeTxReceipts observes the cancelled context and returns its error.
		_, err := wsEther.SubscribeTxReceipts(ctx, q, make(chan []*gethTypes.Receipt))

		// Assert
		assert.ErrorIs(t, err, context.Canceled)
	})
}
