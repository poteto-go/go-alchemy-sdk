package ether_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
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

// wsTestEthAPI is a fake "eth" namespace served over the in-process websocket
// rpc server. It answers eth_blockNumber so the client can do a real round-trip.
type wsTestEthAPI struct{}

func (wsTestEthAPI) BlockNumber() hexutil.Uint64 { return hexutil.Uint64(0x42) }

// newEtherWsApiForTestWithUrl builds an Ether pointed at an explicit ws url
// (e.g. an in-process test server) instead of the derived Alchemy endpoint.
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
			common.FromHex("bcdf20249abf0ed6d944c0288fad489e33f66b3960d9e6229c1cd214ed3bbe31"),
			5<<20,
			nil,
		),
	).(*eth.Ether)
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
		// Arrange: stand up an in-process JSON-RPC server over websocket.
		srv := rpc.NewServer()
		require.NoError(t, srv.RegisterName("eth", wsTestEthAPI{}))

		ts := httptest.NewServer(srv.WebsocketHandler([]string{"*"}))
		// httptest serves over http://; the rpc client dials ws:// on the same host.
		wsUrl := "ws" + strings.TrimPrefix(ts.URL, "http")

		// kill the server (and client) at the end of the lifecycle.
		defer func() {
			ts.Close()
			srv.Stop()
		}()

		wsEther := newEtherWsApiForTestWithUrl(wsUrl)

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
		bn, err := wsEther.Client().BlockNumber(context.Background())
		require.NoError(t, err)
		assert.Equal(t, uint64(0x42), bn)

		// Shutdown tears the persistent ws client down explicitly.
		wsEther.Shutdown()
		assert.Nil(t, wsEther.Client())
	})
}
