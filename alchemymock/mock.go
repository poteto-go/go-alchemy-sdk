package alchemymock

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sync"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/poteto-go/go-alchemy-sdk/gas"
)

// AlchemyHttpMock is a helper for mocking Alchemy API responses in tests.
type AlchemyHttpMock struct {
	baseUrl    string
	responders map[string][]httpmock.Responder
	mu         sync.Mutex
}

// NewAlchemyHttpMock creates a new AlchemyHttpMock and activates httpmock.
// Create alchemymock's instance & alchemy provider w/ same setting.
// It's recommended to call the returned object's DeactivateAndReset method using defer.
//
// Example:
//
//	mock := alchemymock.NewAlchemyHttpMock(setting)
//	defer mock.DeactivateAndReset()
//	alchemy := gas.NewAlchemy(setting)
func NewAlchemyHttpMock(setting gas.AlchemySetting, t testing.TB) *AlchemyHttpMock {
	config := gas.NewAlchemyConfig(setting)
	httpmock.Activate(t)
	mock := &AlchemyHttpMock{
		baseUrl:    config.GetUrl(),
		responders: make(map[string][]httpmock.Responder),
	}
	mock.registerMasterResponder()
	return mock
}

// DeactivateAndReset deactivates httpmock and resets its state.
func (m *AlchemyHttpMock) DeactivateAndReset() {
	httpmock.DeactivateAndReset()
}

type jsonRpcRequest struct {
	Method string `json:"method"`
}

// assert you can call jsonrpc w/ your expected eth method
func (am *AlchemyHttpMock) RegisterResponderOnce(ethMethod, response string) {
	am.registerResponderWithCode(
		http.StatusOK,
		ethMethod,
		response,
	)
}

func (m *AlchemyHttpMock) registerResponderWithCode(statusCode int, ethMethod, response string) {
	responder := httpmock.NewStringResponder(statusCode, response)
	m.mu.Lock()
	defer m.mu.Unlock()
	m.responders[ethMethod] = append(m.responders[ethMethod], responder)
}

func (m *AlchemyHttpMock) registerMasterResponder() {
	httpmock.RegisterResponder("POST", m.baseUrl, func(req *http.Request) (*http.Response, error) {
		var request jsonRpcRequest

		// Read the body
		if req.Body == nil {
			return nil, errors.New("body is nil")
		}
		bodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, errors.New("cannot read body")
		}

		// Restore the body so other matchers can read it
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		if err := json.Unmarshal(bodyBytes, &request); err != nil {
			return nil, errors.New("invalid json")
		}

		m.mu.Lock()
		defer m.mu.Unlock()

		responders, ok := m.responders[request.Method]
		if !ok || len(responders) == 0 {
			return nil, errors.New("method not mocked or no more mocks available")
		}

		// Always pop the first responder (FIFO)
		responder := responders[0]
		m.responders[request.Method] = responders[1:]

		return responder(req)
	})
}
