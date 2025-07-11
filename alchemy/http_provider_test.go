package alchemy

import (
	"context"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/agiledragon/gomonkey"
	"github.com/goccy/go-json"
	"github.com/jarcoal/httpmock"
	"github.com/poteto-go/go-alchemy-sdk/core"
	"github.com/stretchr/testify/assert"
)

func TestNewAlchemyProvider(t *testing.T) {
	// Arrange
	config := NewAlchemyConfig(
		AlchemySetting{
			ApiKey:  "hoge",
			Network: "fuga",
		},
	)

	// Act
	provider := NewAlchemyProvider(config).(*AlchemyProvider)

	// Assert
	assert.Equal(t, config, provider.config)
	assert.Equal(t, 1, provider.id)
}

func newProviderForTest() *AlchemyProvider {
	config := NewAlchemyConfig(
		AlchemySetting{
			ApiKey:  "hoge",
			Network: "fuga",
		},
	)
	return NewAlchemyProvider(config).(*AlchemyProvider)
}

func TestAlchemyProvider_Send(t *testing.T) {
	// Arrange
	provider := newProviderForTest()
	provider.config.backoffConfig.MaxRetries = 0

	t.Run("normal case", func(t *testing.T) {
		t.Run("success request & increment id", func(t *testing.T) {
			httpmock.Activate(t)
			defer httpmock.DeactivateAndReset()

			// Mock
			httpmock.RegisterResponder(
				"POST",
				provider.config.GetUrl(),
				httpmock.NewStringResponder(200, `{"jsonrpc":"2.0","id":1,"result":"0x1234"}`),
			)

			// Act
			result, err := provider.Send("hoge")

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, "0x1234", result)
			assert.Equal(t, 2, provider.id)
		})
	})

	t.Run("error case", func(t *testing.T) {
		t.Run("if failed to create request -> core.ErrFailedToCreateRequest", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Mock
			patches.ApplyFunc(
				http.NewRequestWithContext,
				func(ctx context.Context, method, url string, body io.Reader) (*http.Request, error) {
					return nil, errors.New("error")
				},
			)

			// Act
			_, err := provider.Send("hoge")

			// Assert
			assert.ErrorIs(t, core.ErrFailedToCreateRequest, err)
		})

		t.Run("error on AlchemyFetch", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Mock
			// NOTE: cannot mock generic func
			patches.ApplyFunc(
				json.Marshal,
				func(v any) ([]byte, error) {
					return nil, errors.New("error")
				},
			)

			// Act
			_, err := provider.Send("hoge")

			// Assert
			assert.ErrorIs(t, core.ErrFailedToMarshalParameter, err)
		})
	})
}

func TestAlchemyProvider_SendTransaction(t *testing.T) {
	// Arrange
	provider := newProviderForTest()
	provider.config.backoffConfig.MaxRetries = 0

	t.Run("normal case", func(t *testing.T) {
		t.Run("success request & increment id", func(t *testing.T) {
			httpmock.Activate(t)
			defer httpmock.DeactivateAndReset()

			// Mock
			httpmock.RegisterResponder(
				"POST",
				provider.config.GetUrl(),
				httpmock.NewStringResponder(200, `{"jsonrpc":"2.0","id":1,"result":"0x1234"}`),
			)

			// Act
			result, err := provider.SendTransaction("hoge")

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, "0x1234", result)
			assert.Equal(t, 2, provider.id)
		})
	})
}

func TestAlchemyProvider_SendFilter(t *testing.T) {
	// Arrange
	provider := newProviderForTest()
	provider.config.backoffConfig.MaxRetries = 0

	t.Run("normal case", func(t *testing.T) {
		t.Run("success request & increment id", func(t *testing.T) {
			httpmock.Activate(t)
			defer httpmock.DeactivateAndReset()

			// Mock
			httpmock.RegisterResponder(
				"POST",
				provider.config.GetUrl(),
				httpmock.NewStringResponder(200, `{"jsonrpc":"2.0","id":1,"result":"0x1234"}`),
			)

			// Act
			result, err := provider.SendFilter("hoge")

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, "0x1234", result)
			assert.Equal(t, 2, provider.id)
		})
	})
}
