package gas

import (
	"net/http"
	"time"

	"github.com/poteto-go/go-alchemy-sdk/types"
)

// If you want to run p8 mode,
// you should define host & port
type PrivateNetworkConfig struct {
	Url       string `yaml:"url"`
	Port      int    `yaml:"port"`
	Host      string `yaml:"host"`
	JwtSecret string `yaml:"jwt_secret"`
}

type AlchemySetting struct {
	ApiKey     string        `yaml:"api_key"`
	Network    types.Network `yaml:"network"`
	MaxRetries int           `yaml:"max_retries"`

	// currently not working
	IsRequestBatch bool `yaml:"is_request_batch"`

	// config for backoff retry
	BackoffConfig  *types.BackoffConfig `yaml:"backoff_config"`
	RequestTimeout time.Duration        `yaml:"request_timeout"`

	// You should set if you want to use p8 network
	PrivateNetworkConfig PrivateNetworkConfig `yaml:"private_network_config"`

	CustomHeaders []http.Header `yaml:"custom_headers"`

	// Maximum bytes to read from an RPC response body (default: 32 MiB).
	// Set to 0 to use the default.
	MaxResponseBytes int64 `yaml:"max_response_bytes"`

	// Transport is a caller-supplied http.RoundTripper used for the actual HTTP
	// communication of every RPC call. Use it to plug in connection pooling,
	// retry/backoff, request tracing, latency/error-rate metrics, or to
	// benchmark different (private) RPC providers. If nil, requests delegate to
	// http.DefaultTransport. The SDK always applies its response-size cap on top of it.
	Transport http.RoundTripper `yaml:"-"`

	/*
		return true => p8net is selected

		return false => public is selected
	*/
	IsPrivateNetwork func(setting AlchemySetting) bool
}
