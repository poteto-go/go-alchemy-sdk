package gas

import (
	"time"

	"github.com/poteto-go/go-alchemy-sdk/types"
)

// If you want to run p8 mode,
// you should define host & port
type PrivateNetworkConfig struct {
	Port int    `yaml:"port"`
	Host string `yaml:"host"`
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
}
