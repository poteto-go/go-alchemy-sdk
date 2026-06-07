package ether

import (
	"net/http"
	"time"

	"github.com/poteto-go/go-alchemy-sdk/types"
)

// cf. gsk.AlchemyConfig
// It is for avoiding circle import & keeping api gsk.AlchemyConfig
type EtherApiConfig struct {
	url              string
	maxRetries       int
	requestTimeout   time.Duration
	backoffConfig    *types.BackoffConfig
	customHeaders    []http.Header
	jwtSecret        []byte
	maxResponseBytes int64
	transport        http.RoundTripper
}

func NewEtherApiConfig(
	url string,
	maxRetries int,
	requestTimeout time.Duration,
	backoffConfig *types.BackoffConfig,
	customHeaders []http.Header,
	jwtSecret []byte,
	maxResponseBytes int64,
	transport http.RoundTripper,
) EtherApiConfig {
	return EtherApiConfig{
		url:              url,
		maxRetries:       maxRetries,
		requestTimeout:   requestTimeout,
		backoffConfig:    backoffConfig,
		customHeaders:    customHeaders,
		jwtSecret:        jwtSecret,
		maxResponseBytes: maxResponseBytes,
		transport:        transport,
	}
}

func (config *EtherApiConfig) JwtSecret() []byte {
	return config.jwtSecret
}
