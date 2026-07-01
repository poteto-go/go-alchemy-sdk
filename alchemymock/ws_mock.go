package alchemymock

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"

	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/gorilla/websocket"
	"github.com/poteto-go/go-alchemy-sdk/gas"
)

// wsRPCRequest is a single JSON-RPC 2.0 request frame.
type wsRPCRequest struct {
	JSONRPC string            `json:"jsonrpc"`
	ID      json.RawMessage   `json:"id"`
	Method  string            `json:"method"`
	Params  []json.RawMessage `json:"params"`
}

type wsRPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *wsRPCError     `json:"error,omitempty"`
}

type wsRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type wsSubNotification struct {
	JSONRPC string           `json:"jsonrpc"`
	Method  string           `json:"method"`
	Params  wsSubNotifParams `json:"params"`
}

type wsSubNotifParams struct {
	Subscription string          `json:"subscription"`
	Result       json.RawMessage `json:"result"`
}

// wsMockConn wraps a single WebSocket connection and its active subscriptions.
type wsMockConn struct {
	mu      sync.Mutex
	conn    *websocket.Conn
	subKind map[string]string // subscription ID -> kind (e.g. "newHeads")
}

// writeJSON serialises v and sends it as a single text frame.
// It is safe to call from multiple goroutines.
func (c *wsMockConn) writeJSON(v any) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.conn.WriteJSON(v)
}

var wsUpgrader = websocket.Upgrader{
	CheckOrigin: func(*http.Request) bool { return true },
}

// AlchemyWsMock is an in-process WebSocket JSON-RPC test double that supports
// both regular calls (via RegisterResponderOnce) and push subscriptions
// (via EmitNewHeads and Emit).
//
// Unlike AlchemyHttpMock — which hooks into the http.RoundTripper layer — this
// mock runs a real httptest WebSocket server because geth's rpc.Client dials a
// real socket; there is no WS equivalent of httpmock's transport intercept.
//
// Usage:
//
//	setting := gas.AlchemySetting{BackoffConfig: &types.BackoffConfig{MaxRetries: 0}}
//	mock := alchemymock.NewAlchemyWsMock(setting, t)
//	a, _ := mock.NewAlchemy()
//	mock.RegisterResponderOnce("eth_blockNumber", `{"jsonrpc":"2.0","id":1,"result":"0x42"}`)
//	bn, _ := a.Core.GetBlockNumber() // -> 66
type AlchemyWsMock struct {
	t          testing.TB
	setting    gas.AlchemySetting
	ts         *httptest.Server
	mu         sync.Mutex
	responders map[string][]json.RawMessage // method -> queued "result" values
	conns      []*wsMockConn
	nextSubID  atomic.Uint64
	handlers   sync.WaitGroup // tracks active serveWS goroutines
}

// NewAlchemyWsMock creates a new in-process WebSocket JSON-RPC mock server and
// registers a cleanup that shuts it down when the test ends.
func NewAlchemyWsMock(setting gas.AlchemySetting, t testing.TB) *AlchemyWsMock {
	m := &AlchemyWsMock{
		t:          t,
		setting:    setting,
		responders: make(map[string][]json.RawMessage),
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", m.serveWS)
	m.ts = httptest.NewServer(mux)
	t.Cleanup(m.Close)
	return m
}

// NewAlchemy returns a gas.Alchemy wired to the mock's WebSocket endpoint,
// using the setting passed to NewAlchemyWsMock with the URL overridden.
// This mirrors how AlchemyHttpMock callers do gas.NewAlchemy(setting).
func (m *AlchemyWsMock) NewAlchemy() (gas.Alchemy, error) {
	s := m.setting
	s.PrivateNetworkConfig = gas.PrivateNetworkConfig{Url: m.URL()}
	return gas.NewAlchemy(s)
}

// URL returns the ws:// endpoint of the mock server.
// Pass it as gas.PrivateNetworkConfig.Url to route an Alchemy instance through
// this mock.
func (m *AlchemyWsMock) URL() string {
	return "ws" + m.ts.URL[len("http"):]
}

// Close shuts down the mock server and waits for active handler goroutines to
// finish. It is idempotent and called automatically by the test cleanup
// registered in NewAlchemyWsMock.
func (m *AlchemyWsMock) Close() {
	m.ts.Close()
	// http.Server.Close does not close hijacked WebSocket connections.
	// Explicitly close them to unblock the ReadMessage loops in serveWS goroutines.
	m.mu.Lock()
	conns := make([]*wsMockConn, len(m.conns))
	copy(conns, m.conns)
	m.mu.Unlock()
	for _, mc := range conns {
		mc.conn.Close()
	}
	m.handlers.Wait()
}

// RegisterResponderOnce queues response as the next JSON-RPC reply for method.
// The format mirrors AlchemyHttpMock.RegisterResponderOnce: pass the full
// JSON-RPC envelope, e.g.:
//
//	mock.RegisterResponderOnce("eth_blockNumber", `{"jsonrpc":"2.0","id":1,"result":"0x42"}`)
//
// The mock extracts the "result" field and re-uses the actual request's id when
// sending the response, so the id in response is ignored.
// It calls t.Fatalf if response is not a valid JSON-RPC envelope with a "result" field.
func (m *AlchemyWsMock) RegisterResponderOnce(method, response string) {
	var envelope struct {
		Result json.RawMessage `json:"result"`
	}
	if err := json.Unmarshal([]byte(response), &envelope); err != nil || len(envelope.Result) == 0 {
		m.t.Fatalf("AlchemyWsMock.RegisterResponderOnce: %q is not a valid JSON-RPC envelope with a result field", response)
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.responders[method] = append(m.responders[method], envelope.Result)
}

// Emit pushes a raw subscription notification of the given kind to all
// connections that hold an active subscription of that kind.
// Use this to push any subscription type not covered by a typed helper:
//
//	mock.Emit("logs", rawLogsJSON)
func (m *AlchemyWsMock) Emit(kind string, result json.RawMessage) {
	m.mu.Lock()
	conns := make([]*wsMockConn, len(m.conns))
	copy(conns, m.conns)
	m.mu.Unlock()

	for _, mc := range conns {
		mc.mu.Lock()
		var matchIDs []string
		for subID, k := range mc.subKind {
			if k == kind {
				matchIDs = append(matchIDs, subID)
			}
		}
		mc.mu.Unlock()

		for _, subID := range matchIDs {
			_ = mc.writeJSON(wsSubNotification{
				JSONRPC: "2.0",
				Method:  "eth_subscription",
				Params:  wsSubNotifParams{Subscription: subID, Result: result},
			})
		}
	}
}

// EmitNewHeads pushes a newHeads subscription notification for each header to
// all connections that have an active newHeads subscription.
func (m *AlchemyWsMock) EmitNewHeads(headers ...*gethTypes.Header) {
	for _, h := range headers {
		data, err := json.Marshal(h)
		if err != nil {
			m.t.Logf("AlchemyWsMock: marshal header: %v", err)
			continue
		}
		m.Emit("newHeads", data)
	}
}

func (m *AlchemyWsMock) serveWS(w http.ResponseWriter, r *http.Request) {
	m.handlers.Add(1)
	defer m.handlers.Done()

	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		m.t.Logf("AlchemyWsMock: upgrade: %v", err)
		return
	}

	mc := &wsMockConn{conn: conn, subKind: make(map[string]string)}

	m.mu.Lock()
	m.conns = append(m.conns, mc)
	m.mu.Unlock()

	defer func() {
		conn.Close()
		m.mu.Lock()
		for i, c := range m.conns {
			if c == mc {
				m.conns = append(m.conns[:i], m.conns[i+1:]...)
				break
			}
		}
		m.mu.Unlock()
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}

		var req wsRPCRequest
		if err := json.Unmarshal(msg, &req); err != nil {
			m.t.Logf("AlchemyWsMock: bad frame: %v", err)
			continue
		}

		switch req.Method {
		case "eth_subscribe":
			m.handleSubscribe(mc, &req)
		case "eth_unsubscribe":
			m.handleUnsubscribe(mc, &req)
		default:
			m.handleCall(mc, &req)
		}
	}
}

func (m *AlchemyWsMock) newSubID() string {
	n := m.nextSubID.Add(1)
	return fmt.Sprintf("0x%x", n)
}

func (m *AlchemyWsMock) handleSubscribe(mc *wsMockConn, req *wsRPCRequest) {
	kind := "unknown"
	if len(req.Params) > 0 {
		_ = json.Unmarshal(req.Params[0], &kind)
	}

	subID := m.newSubID()
	subIDJSON, _ := json.Marshal(subID)

	mc.mu.Lock()
	mc.subKind[subID] = kind
	mc.mu.Unlock()

	_ = mc.writeJSON(wsRPCResponse{JSONRPC: "2.0", ID: req.ID, Result: subIDJSON})
}

func (m *AlchemyWsMock) handleUnsubscribe(mc *wsMockConn, req *wsRPCRequest) {
	var subID string
	if len(req.Params) > 0 {
		_ = json.Unmarshal(req.Params[0], &subID)
	}

	mc.mu.Lock()
	delete(mc.subKind, subID)
	mc.mu.Unlock()

	_ = mc.writeJSON(wsRPCResponse{JSONRPC: "2.0", ID: req.ID, Result: json.RawMessage("true")})
}

func (m *AlchemyWsMock) handleCall(mc *wsMockConn, req *wsRPCRequest) {
	m.mu.Lock()
	var result json.RawMessage
	if results := m.responders[req.Method]; len(results) > 0 {
		result = results[0]
		m.responders[req.Method] = results[1:]
	}
	m.mu.Unlock()

	if result != nil {
		_ = mc.writeJSON(wsRPCResponse{JSONRPC: "2.0", ID: req.ID, Result: result})
	} else {
		_ = mc.writeJSON(wsRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error:   &wsRPCError{Code: -32601, Message: fmt.Sprintf("method not mocked or no more mocks: %s", req.Method)},
		})
	}
}
