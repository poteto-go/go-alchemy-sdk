package alchemymock

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/poteto-go/go-alchemy-sdk/gas"
)

// AlchemyHttpMock is a helper for mocking Alchemy API responses in tests.
type AlchemyHttpMock struct {
	baseUrl string
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
	return &AlchemyHttpMock{
		baseUrl: config.GetUrl(),
	}
}

// DeactivateAndReset deactivates httpmock and resets its state.
func (m *AlchemyHttpMock) DeactivateAndReset() {
	httpmock.DeactivateAndReset()
}

type jsonRpcRequest struct {
	Method string `json:"method"`
}

// assert you can call jsonrpc w/ your expected eth method
func (am *AlchemyHttpMock) RegisterResponder(ethMethod, response string) {
	am.registerResponderWithCode(
		http.StatusOK,
		ethMethod,
		response,
	)
}

func (m *AlchemyHttpMock) registerResponderWithCode(statusCode int, ethMethod, response string) {
	responder := httpmock.NewStringResponder(statusCode, response)
	matcherFunc := func(req *http.Request) bool {
		var request jsonRpcRequest
		if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
			return false
		}

		return request.Method == ethMethod
	}

	httpmock.RegisterMatcherResponder(
		"POST",
		m.baseUrl,
		httpmock.NewMatcher("match eth method", matcherFunc),
		responder,
	)
}
