package alchemymock

import (
	"context"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
)

// AlchemyWsMock is an in-process WebSocket test double for subscription-based
// Alchemy code (eth_subscribe newHeads / logs / pending / alchemy_minedTransactions).
//
// Unlike AlchemyHttpMock — which is built on jarcoal/httpmock and intercepts at
// the http.RoundTripper layer (one request -> one response) — geth dials ws via
// gorilla/websocket's raw TCP + HTTP Upgrade and a subscription is one request
// followed by an unbounded server-pushed stream. Neither fits the responder
// model, so this is a separate implementation backed by a real geth rpc.Server
// over httptest. Push canned notifications with the Emit* helpers.
//
// Example:
//
//	mock := alchemymock.NewAlchemyWsMock(t)
//	defer mock.Close()
//
//	a, _ := gas.NewAlchemy(gas.AlchemySetting{
//		PrivateNetworkConfig: gas.PrivateNetworkConfig{Url: mock.URL()},
//	})
//	sub := a.GetProvider().(types.ISubscribeProvider)
//	ch := make(chan *gethTypes.Header)
//	sub.Subscribe(ctx, ch, "newHeads")
//	mock.EmitNewHeads(headers...)
type AlchemyWsMock struct {
	server *httptest.Server
	rpcSrv *rpc.Server
	api    *wsEthAPI
}

// NewAlchemyWsMock stands up a geth rpc.Server exposing a fake "eth" namespace
// over a websocket httptest server. Pair it with a ws-scheme Alchemy setting
// pointed at URL(). Call Close (defer recommended) to tear it down.
func NewAlchemyWsMock(t testing.TB) *AlchemyWsMock {
	api := &wsEthAPI{subs: make(map[rpc.ID]*rpc.Notifier)}

	srv := rpc.NewServer()
	if err := srv.RegisterName("eth", api); err != nil {
		t.Fatal(err)
	}

	// httptest serves over http://; geth's rpc client dials ws:// on the same host.
	server := httptest.NewServer(srv.WebsocketHandler([]string{"*"}))

	return &AlchemyWsMock{server: server, rpcSrv: srv, api: api}
}

// URL returns the ws:// endpoint to feed an Alchemy ws setting.
func (m *AlchemyWsMock) URL() string {
	return "ws" + strings.TrimPrefix(m.server.URL, "http")
}

// Close stops the websocket server and any active subscriptions. Safe to call
// more than once: both httptest.Server.Close and rpc.Server.Stop are idempotent.
func (m *AlchemyWsMock) Close() {
	m.server.Close()
	m.rpcSrv.Stop()
}

// EmitNewHeads pushes the given headers to every active newHeads subscriber, in
// order. Call it after the subscription is established (i.e. after Subscribe has
// returned); with no subscribers it is a no-op.
func (m *AlchemyWsMock) EmitNewHeads(headers ...*gethTypes.Header) {
	for _, h := range headers {
		m.api.emit(h)
	}
}

// wsEthAPI is the fake "eth" namespace served over the in-process ws rpc server.
// Its subscription methods register active subscribers (keyed by subscription id)
// that the Emit* helpers fan notifications out to.
type wsEthAPI struct {
	mu   sync.Mutex
	subs map[rpc.ID]*rpc.Notifier
}

// NewHeads implements eth_subscribe("newHeads"). geth derives the subscription
// name from the method name, so this is reachable via EthSubscribe(ctx, ch,
// "newHeads"). It registers the subscriber and unregisters it when the client
// unsubscribes or the connection drops.
func (api *wsEthAPI) NewHeads(ctx context.Context) (*rpc.Subscription, error) {
	notifier, supported := rpc.NotifierFromContext(ctx)
	if !supported {
		return nil, rpc.ErrNotificationsUnsupported
	}

	rpcSub := notifier.CreateSubscription()
	api.add(rpcSub.ID, notifier)

	go func() {
		<-rpcSub.Err() // fires on Unsubscribe or connection drop
		api.remove(rpcSub.ID)
	}()

	return rpcSub, nil
}

func (api *wsEthAPI) add(id rpc.ID, notifier *rpc.Notifier) {
	api.mu.Lock()
	defer api.mu.Unlock()
	api.subs[id] = notifier
}

func (api *wsEthAPI) remove(id rpc.ID) {
	api.mu.Lock()
	defer api.mu.Unlock()
	delete(api.subs, id)
}

// emit fans a single notification out to every active subscriber.
func (api *wsEthAPI) emit(data any) {
	api.mu.Lock()
	defer api.mu.Unlock()
	for id, notifier := range api.subs {
		notifier.Notify(id, data)
	}
}
