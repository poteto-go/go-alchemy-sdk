package utils_test

import (
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/poteto-go/go-alchemy-sdk/utils"
	"github.com/stretchr/testify/assert"
)

func TestNewSharedHTTPClient(t *testing.T) {
	t.Run("returns non-nil client", func(t *testing.T) {
		client := utils.NewSharedHTTPClient(0, 0)
		assert.NotNil(t, client)
	})

	t.Run("uses limitedTransport, not nil", func(t *testing.T) {
		client := utils.NewSharedHTTPClient(100, 0)
		assert.NotNil(t, client.Transport)
	})

	t.Run("two calls return different client instances", func(t *testing.T) {
		c1 := utils.NewSharedHTTPClient(0, 0)
		c2 := utils.NewSharedHTTPClient(0, 0)
		assert.NotSame(t, c1, c2)
	})

	t.Run("timeout is set on the client", func(t *testing.T) {
		client := utils.NewSharedHTTPClient(0, 5*time.Second)
		assert.Equal(t, 5*time.Second, client.Timeout)
	})

	t.Run("limits response body to maxBytes", func(t *testing.T) {
		httpmock.Activate(t)
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder(
			"GET", "http://example.com/",
			httpmock.NewStringResponder(200, string(make([]byte, 100))),
		)

		client := utils.NewSharedHTTPClient(10, 0)
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

		client := utils.NewSharedHTTPClient(0, 0)
		req, _ := http.NewRequest("GET", "http://example.com/", nil)
		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		assert.Equal(t, len(payload), len(body), "small response should pass through default limit")
	})
}
