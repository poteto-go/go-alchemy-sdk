package alchemy

import (
	"context"
	"errors"
	"io"
	"net/http"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey"
	"github.com/goccy/go-json"
	"github.com/jarcoal/httpmock"
	"github.com/poteto-go/go-alchemy-sdk/core"
	"github.com/poteto-go/go-alchemy-sdk/utils"
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

func TestAlchemyProvider_GetBlockNumber(t *testing.T) {
	// Arrange
	provider := newProviderForTest()

	t.Run("normal case", func(t *testing.T) {
		t.Run("success request", func(t *testing.T) {
			httpmock.Activate(t)
			patches := gomonkey.NewPatches()
			defer func() {
				httpmock.DeactivateAndReset()
				patches.Reset()
			}()

			// Mock
			httpmock.RegisterResponder(
				"POST",
				provider.config.GetUrl(),
				httpmock.NewStringResponder(200, `{"jsonrpc":"2.0","id":1,"result":"0x1234"}`),
			)

			patches.ApplyFunc(
				utils.FromHex,
				func(s string) (int, error) {
					return 1234, nil
				},
			)
			// Act
			result, err := provider.GetBlockNumber()

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, 1234, result)
		})
	})

	t.Run("error case", func(t *testing.T) {
		t.Run("if failed to send request -> core.ErrFailedToConnect", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Act
			_, err := provider.GetBlockNumber()

			// Assert
			assert.ErrorIs(t, core.ErrFailedToConnect, err)
		})

		t.Run("if failed from hex -> error", func(t *testing.T) {
			httpmock.Activate(t)
			patches := gomonkey.NewPatches()
			defer func() {
				httpmock.DeactivateAndReset()
				patches.Reset()
			}()

			// Mock
			httpmock.RegisterResponder(
				"POST",
				provider.config.GetUrl(),
				httpmock.NewStringResponder(200, `{"jsonrpc":"2.0","id":1,"result":"0x1234"}`),
			)

			patches.ApplyFunc(
				utils.FromHex,
				func(s string) (int, error) {
					return 0, core.ErrInvalidHexString
				},
			)
			// Act
			_, err := provider.GetBlockNumber()

			// Assert
			assert.ErrorIs(t, core.ErrInvalidHexString, err)
		})
	})
}

func TestAlchemyProvider_Send(t *testing.T) {
	// Arrange
	provider := newProviderForTest()

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
		t.Run("if failed to marshal parameter -> core.ErrFailedToMarshalParameter", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Mock
			patches.ApplyFunc(
				json.Marshal,
				func(v interface{}) ([]byte, error) {
					return nil, errors.New("error")
				},
			)

			// Act
			_, err := provider.Send("hoge")

			// Assert
			assert.ErrorIs(t, core.ErrFailedToMarshalParameter, err)
		})

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

		t.Run("if failed to request -> core.ErrFailedToConnect", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Mock
			patches.ApplyMethod(
				reflect.TypeOf(http.DefaultClient),
				"Do",
				func(c *http.Client, req *http.Request) (*http.Response, error) {
					return nil, errors.New("error")
				},
			)

			// Act
			_, err := provider.Send("hoge")

			// Assert
			assert.ErrorIs(t, core.ErrFailedToConnect, err)
		})
	})
}
