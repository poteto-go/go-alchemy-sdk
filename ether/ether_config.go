package ether

import (
	"time"

	"github.com/poteto-go/go-alchemy-sdk/types"
)

// cf. gsk.AlchemyConfig
// It is for avoiding circle import & keeping api gsk.AlchemyConfig
type EtherApiConfig struct {
	url            string
	maxRetries     int
	requestTimeout time.Duration
	backoffConfig  *types.BackoffConfig
}

func NewEtherApiConfig(url string, maxRetries int, requestTimeout time.Duration, backoffConfig *types.BackoffConfig) EtherApiConfig {
	return EtherApiConfig{
		url:            url,
		maxRetries:     maxRetries,
		requestTimeout: requestTimeout,
		backoffConfig:  backoffConfig,
	}
}
