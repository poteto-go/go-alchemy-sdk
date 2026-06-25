package gas

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/poteto-go/go-alchemy-sdk/constant"
	"github.com/poteto-go/go-alchemy-sdk/ether"
	"github.com/poteto-go/go-alchemy-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// wsTestAPI is a fake "eth" namespace served over the in-process websocket rpc
// server. It answers eth_blockNumber (a plain call) and eth_subscribe("ticks")
// (a push subscription) so the ws provider can do real round-trips.
type wsTestAPI struct{}

func (wsTestAPI) BlockNumber() hexutil.Uint64 { return hexutil.Uint64(0x42) }

func (wsTestAPI) Ticks(ctx context.Context) (*rpc.Subscription, error) {
	notifier, supported := rpc.NotifierFromContext(ctx)
	if !supported {
		return nil, rpc.ErrNotificationsUnsupported
	}

	sub := notifier.CreateSubscription()
	go func() {
		notifier.Notify(sub.ID, "0x1")
	}()
	return sub, nil
}

// newWsProviderForTest stands up an in-process JSON-RPC ws server and returns a
// WsAlchemyProvider whose Ether is dialed at that server. The cleanup closes the
// server and the persistent ws socket.
func newWsProviderForTest(t *testing.T) *WsAlchemyProvider {
	t.Helper()

	srv := rpc.NewServer()
	require.NoError(t, srv.RegisterName("eth", wsTestAPI{}))

	ts := httptest.NewServer(srv.WebsocketHandler([]string{"*"}))
	// httptest serves over http://; the rpc client dials ws:// on the same host.
	wsUrl := "ws" + strings.TrimPrefix(ts.URL, "http")

	config, err := NewAlchemyConfig(AlchemySetting{
		ApiKey:       "hoge",
		Network:      "fuga",
		UseWebsocket: true,
	})
	require.NoError(t, err)

	provider := NewWsAlchemyProvider(config).(*WsAlchemyProvider)
	eth := ether.NewEtherApi(
		provider,
		// override the derived alchemy endpoint with the in-process ws url.
		ether.NewEtherApiConfig(wsUrl, 0, 2*time.Second, &types.DefaultBackoffConfig, []http.Header{}, nil, 5<<20, nil),
	)
	provider.SetEth(eth)

	t.Cleanup(func() {
		eth.Shutdown()
		ts.Close()
		srv.Stop()
	})
	return provider
}

func TestNewWsAlchemyProvider(t *testing.T) {
	customHeaders := []http.Header{{"hello": []string{"world"}}}
	config, err := NewAlchemyConfig(AlchemySetting{
		ApiKey:        "hoge",
		Network:       "fuga",
		UseWebsocket:  true,
		CustomHeaders: customHeaders,
	})
	require.NoError(t, err)

	provider := NewWsAlchemyProvider(config).(*WsAlchemyProvider)

	assert.Equal(t, types.Network("fuga"), provider.Network())
	assert.Equal(t, customHeaders, provider.CustomHeaders())
}

func TestWsAlchemyProvider_Send(t *testing.T) {
	t.Run("routes the call over the ws socket", func(t *testing.T) {
		provider := newWsProviderForTest(t)

		result, err := provider.Send("eth_blockNumber", types.RequestArgs{})

		require.NoError(t, err)
		assert.Equal(t, "0x42", result)
	})

	t.Run("returns error if eth client is not set", func(t *testing.T) {
		config, _ := NewAlchemyConfig(AlchemySetting{ApiKey: "k", Network: "n", UseWebsocket: true})
		provider := NewWsAlchemyProvider(config).(*WsAlchemyProvider)

		_, err := provider.Send("eth_blockNumber", types.RequestArgs{})
		assert.ErrorIs(t, err, constant.ErrProviderEthNotSet)
	})
}

func TestWsAlchemyProvider_Subscribe(t *testing.T) {
	provider := newWsProviderForTest(t)

	ch := make(chan string, 1)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	sub, err := provider.Subscribe(ctx, ch, "ticks")
	require.NoError(t, err)
	defer sub.Unsubscribe()

	select {
	case got := <-ch:
		assert.Equal(t, "0x1", got)
	case err := <-sub.Err():
		t.Fatalf("subscription errored: %v", err)
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for subscription notification")
	}
}
