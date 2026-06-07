package utils_test

import (
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/poteto-go/go-alchemy-sdk/utils"
	"github.com/stretchr/testify/assert"
)

func TestNewSharedHTTPClient(t *testing.T) {
	t.Run("returns non-nil client", func(t *testing.T) {
		client := utils.NewSharedHTTPClient(0, 0, nil)
		assert.NotNil(t, client)
	})

	t.Run("uses limitedTransport, not nil", func(t *testing.T) {
		client := utils.NewSharedHTTPClient(100, 0, nil)
		assert.NotNil(t, client.Transport)
	})

	t.Run("two calls return different client instances", func(t *testing.T) {
		c1 := utils.NewSharedHTTPClient(0, 0, nil)
		c2 := utils.NewSharedHTTPClient(0, 0, nil)
		assert.NotSame(t, c1, c2)
	})

	t.Run("timeout is set on the client", func(t *testing.T) {
		client := utils.NewSharedHTTPClient(0, 5*time.Second, nil)
		assert.Equal(t, 5*time.Second, client.Timeout)
	})

	t.Run("uses custom transport when provided", func(t *testing.T) {
		custom := &recordingTransport{}
		client := utils.NewSharedHTTPClient(0, 0, custom)

		req, _ := http.NewRequest("GET", "http://example.com/", nil)
		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.True(t, custom.called, "custom transport should be used as the underlying RoundTripper")
	})

	t.Run("custom transport response is still capped at maxBytes", func(t *testing.T) {
		custom := &recordingTransport{body: string(make([]byte, 100))}
		client := utils.NewSharedHTTPClient(10, 0, custom)

		req, _ := http.NewRequest("GET", "http://example.com/", nil)
		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		assert.Equal(t, 10, len(body), "custom transport body should still be capped")
	})

	t.Run("limits response body to maxBytes", func(t *testing.T) {
		httpmock.Activate(t)
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder(
			"GET", "http://example.com/",
			httpmock.NewStringResponder(200, string(make([]byte, 100))),
		)

		client := utils.NewSharedHTTPClient(10, 0, nil)
		req, _ := http.NewRequest("GET", "http://example.com/", nil)
		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		assert.Equal(t, 10, len(body), "body should be capped at maxBytes")
	})

	t.Run("uses default maxResponseBytes when maxBytes is 0", func(t *testing.T) {
		httpmock.Activate(t)
		defer httpmock.DeactivateAndReset()

		payload := string(make([]byte, 100))
		httpmock.RegisterResponder(
			"GET", "http://example.com/",
			httpmock.NewStringResponder(200, payload),
		)

		client := utils.NewSharedHTTPClient(0, 0, nil)
		req, _ := http.NewRequest("GET", "http://example.com/", nil)
		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		assert.Equal(t, len(payload), len(body), "small response should pass through default limit")
	})
}

// recordingTransport is a test RoundTripper that records invocation and returns
// a fixed body, used to verify that a caller-supplied transport is actually used.
type recordingTransport struct {
	called bool
	body   string
}

func (rt *recordingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	rt.called = true
	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(rt.body)),
		Header:     make(http.Header),
		Request:    req,
	}, nil
}
