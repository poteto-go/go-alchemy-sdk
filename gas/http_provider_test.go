package gas

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sync"
	"testing"

	"github.com/agiledragon/gomonkey"
	"github.com/jarcoal/httpmock"
	"github.com/poteto-go/go-alchemy-sdk/constant"
	"github.com/poteto-go/go-alchemy-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestNewAlchemyProvider(t *testing.T) {
	// Arrange
	customHeaders := []http.Header{
		{
			"hello": []string{"world"},
		},
	}
	config, err := NewAlchemyConfig(
		AlchemySetting{
			ApiKey:        "hoge",
			Network:       "fuga",
			CustomHeaders: customHeaders,
		},
	)
	assert.NoError(t, err)

	// Act
	provider := NewAlchemyProvider(config).(*AlchemyProvider)

	// Assert
	assert.Equal(t, config, provider.config)
	assert.Equal(t, int64(1), provider.id.Load())
	assert.Equal(t, config.customHeaders, customHeaders)
	assert.NotNil(t, provider.client, "shared http.Client must be created at construction")
}

func TestNewAlchemyProvider_ClientHasLimitedTransport(t *testing.T) {
	config, _ := NewAlchemyConfig(AlchemySetting{ApiKey: "k", Network: "n"})
	provider := NewAlchemyProvider(config).(*AlchemyProvider)

	// Transport must be non-nil (set by NewSharedHTTPClient), not a bare &http.Client{}.
	// A nil Transport would mean the size-cap limitedTransport was not installed.
	assert.NotNil(t, provider.client.Transport, "client must use limitedTransport from NewSharedHTTPClient")
}

func newProviderForTest() *AlchemyProvider {
	config, _ := NewAlchemyConfig(
		AlchemySetting{
			ApiKey:  "hoge",
			Network: "fuga",
			CustomHeaders: []http.Header{
				{
					"hello": []string{"world"},
				},
			},
		},
	)
	return NewAlchemyProvider(config).(*AlchemyProvider)
}

func TestAlchemyProvider_Network(t *testing.T) {
	t.Run("returns the network from config", func(t *testing.T) {
		provider := newProviderForTest()
		assert.Equal(t, types.Network("fuga"), provider.Network())
	})
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
			result, err := provider.Send("hoge", types.RequestArgs{})

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, "0x1234", result)
			assert.Equal(t, int64(2), provider.id.Load())
		})

		t.Run("concurrent Send produces unique JSON-RPC ids", func(t *testing.T) {
			provider := newProviderForTest()
			provider.config.backoffConfig.MaxRetries = 0

			httpmock.Activate(t)
			defer httpmock.DeactivateAndReset()

			var (
				mu  sync.Mutex
				ids []int
			)
			httpmock.RegisterResponder(
				"POST",
				provider.config.GetUrl(),
				func(req *http.Request) (*http.Response, error) {
					var body types.AlchemyRequestBody
					if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
						return nil, err
					}
					mu.Lock()
					ids = append(ids, body.Id)
					mu.Unlock()
					return httpmock.NewStringResponse(
						200,
						`{"jsonrpc":"2.0","id":1,"result":"0x1"}`,
					), nil
				},
			)

			const goroutines = 50
			var wg sync.WaitGroup
			for i := 0; i < goroutines; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					_, _ = provider.Send("hoge", types.RequestArgs{})
				}()
			}
			wg.Wait()

			assert.Len(t, ids, goroutines)
			seen := map[int]bool{}
			for _, id := range ids {
				assert.Falsef(t, seen[id], "duplicate JSON-RPC id %d generated under concurrency", id)
				seen[id] = true
			}
		})
	})

	t.Run("error case", func(t *testing.T) {
		t.Run("if failed to create request -> constant.ErrFailedToCreateRequest", func(t *testing.T) {
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
			_, err := provider.Send("hoge", types.RequestArgs{})

			// Assert
			assert.ErrorIs(t, constant.ErrFailedToCreateRequest, err)
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
			_, err := provider.Send("hoge", types.RequestArgs{})

			// Assert
			assert.ErrorIs(t, constant.ErrFailedToMarshalParameter, err)
		})

		t.Run("if not error, but result is nil, return constant.ErrResultIsNil", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Mock
			httpmock.Activate(t)
			defer httpmock.DeactivateAndReset()

			// Mock
			httpmock.RegisterResponder(
				"POST",
				provider.config.GetUrl(),
				httpmock.NewStringResponder(200, `{"jsonrpc":"2.0","id":1,"result":null}`),
			)

			// Act
			_, err := provider.Send("hoge", types.RequestArgs{})

			// Assert
			assert.ErrorIs(t, constant.ErrResultIsNil, err)
		})
	})
}
