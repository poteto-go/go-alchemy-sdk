package namespace_test

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
	"github.com/poteto-go/go-alchemy-sdk/namespace"
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

func newEtherWsApiForTest() *eth.Ether {
	config, err := gas.NewAlchemyConfig(utWsAlchemySetting)
	if err != nil {
		panic(err)
	}

	return ether.NewEtherApi(
		nil,
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

// newWsNamespaceForTestWithUrl builds an IWS backed by an Ether dialed at an
// explicit ws url (e.g. the one served by alchemymock.AlchemyWsMock).
func newWsNamespaceForTestWithUrl(wsUrl string) namespace.IWS {
	config, err := gas.NewAlchemyConfig(utWsAlchemySetting)
	if err != nil {
		panic(err)
	}
	provider := gas.NewWsAlchemyProvider(config)
	e := ether.NewEtherApi(
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
	)
	return namespace.NewWSNamespace(e)
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

func TestNewWSNamespace(t *testing.T) {
	// Arrange
	ether := newEtherApi()

	// Act
	ws := namespace.NewWSNamespace(ether)

	// Assert
	assert.NotNil(t, ws)
}

func TestWS_Subscribe(t *testing.T) {
	t.Run("returns error if provider is not websocket", func(t *testing.T) {
		// Arrange: a plain (http) Ether cannot open eth_subscribe streams.
		ws := namespace.NewWSNamespace(newEtherApi())

		// Act
		_, err := ws.Subscribe(context.Background(), make(chan *gethTypes.Header), "newHeads")

		// Assert
		assert.ErrorIs(t, err, constant.ErrUnsupportedNotWebsocketProvider)
	})

	t.Run("streams subscription notifications over the ws socket", func(t *testing.T) {
		// Arrange
		mock := alchemymock.NewAlchemyWsMock(utWsAlchemySetting, t)
		ws := newWsNamespaceForTestWithUrl(mock.URL())

		ch := make(chan *gethTypes.Header, 1)
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		// Act
		sub, err := ws.Subscribe(ctx, ch, "newHeads")
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
}

func TestWS_SubscribeNewHead(t *testing.T) {
	t.Run("returns error if provider is not websocket", func(t *testing.T) {
		// Arrange
		ws := namespace.NewWSNamespace(newEtherApi())

		// Act
		_, err := ws.SubscribeNewHead(context.Background(), make(chan *gethTypes.Header))

		// Assert
		assert.ErrorIs(t, err, constant.ErrUnsupportedNotWebsocketProvider)
	})

	t.Run("streams new head notifications over the ws socket", func(t *testing.T) {
		// Arrange
		mock := alchemymock.NewAlchemyWsMock(utWsAlchemySetting, t)
		ws := newWsNamespaceForTestWithUrl(mock.URL())

		ch := make(chan *gethTypes.Header, 1)
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		// Act
		sub, err := ws.SubscribeNewHead(ctx, ch)
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
}

func TestWS_SubscribeLogs(t *testing.T) {
	query := ethereum.FilterQuery{}

	t.Run("returns error if provider is not websocket", func(t *testing.T) {
		// Arrange
		ws := namespace.NewWSNamespace(newEtherApi())

		// Act
		_, err := ws.SubscribeLogs(context.Background(), query, make(chan gethTypes.Log))

		// Assert
		assert.ErrorIs(t, err, constant.ErrUnsupportedNotWebsocketProvider)
	})

	t.Run("streams log notifications over the ws socket", func(t *testing.T) {
		// Arrange
		mock := alchemymock.NewAlchemyWsMock(utWsAlchemySetting, t)
		ws := newWsNamespaceForTestWithUrl(mock.URL())

		ch := make(chan gethTypes.Log, 1)
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		// Act
		sub, err := ws.SubscribeLogs(ctx, query, ch)
		require.NoError(t, err)
		defer sub.Unsubscribe()

		addr := common.HexToAddress("0xdeadbeef00000000000000000000000000000042")
		data, err := json.Marshal(newWsMinimalLog(addr))
		require.NoError(t, err)
		mock.Emit("logs", data)

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
}

func TestWS_SubscribeContractLogs(t *testing.T) {
	contractAddress := common.HexToAddress("0xdeadbeef00000000000000000000000000000042")

	t.Run("returns error if provider is not websocket", func(t *testing.T) {
		// Arrange
		ws := namespace.NewWSNamespace(newEtherApi())

		// Act
		_, err := ws.SubscribeContractLogs(context.Background(), contractAddress, make(chan gethTypes.Log))

		// Assert
		assert.ErrorIs(t, err, constant.ErrUnsupportedNotWebsocketProvider)
	})

	t.Run("streams contract log notifications over the ws socket", func(t *testing.T) {
		// Arrange
		mock := alchemymock.NewAlchemyWsMock(utWsAlchemySetting, t)
		ws := newWsNamespaceForTestWithUrl(mock.URL())

		ch := make(chan gethTypes.Log, 1)
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		// Act
		sub, err := ws.SubscribeContractLogs(ctx, contractAddress, ch)
		require.NoError(t, err)
		defer sub.Unsubscribe()

		data, err := json.Marshal(newWsMinimalLog(contractAddress))
		require.NoError(t, err)
		mock.Emit("logs", data)

		// Assert
		select {
		case got := <-ch:
			assert.Equal(t, contractAddress, got.Address)
		case err := <-sub.Err():
			t.Fatalf("subscription errored: %v", err)
		case <-time.After(2 * time.Second):
			t.Fatal("timed out waiting for subscription notification")
		}
	})
}

func TestWS_SubscribeTxReceipts(t *testing.T) {
	query := &ethereum.TransactionReceiptsQuery{}

	t.Run("returns error if provider is not websocket", func(t *testing.T) {
		// Arrange
		ws := namespace.NewWSNamespace(newEtherApi())

		// Act
		_, err := ws.SubscribeTxReceipts(context.Background(), query, make(chan []*gethTypes.Receipt))

		// Assert
		assert.ErrorIs(t, err, constant.ErrUnsupportedNotWebsocketProvider)
	})

	t.Run("streams tx receipt notifications over the ws socket", func(t *testing.T) {
		// Arrange
		mock := alchemymock.NewAlchemyWsMock(utWsAlchemySetting, t)
		ws := newWsNamespaceForTestWithUrl(mock.URL())

		ch := make(chan []*gethTypes.Receipt, 1)
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		// Act
		sub, err := ws.SubscribeTxReceipts(ctx, query, ch)
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
}
