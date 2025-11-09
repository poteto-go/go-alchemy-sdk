package gas

import (
	"time"

	"github.com/poteto-go/go-alchemy-sdk/types"
)

type AlchemySetting struct {
	ApiKey         string                `yaml:"api_key"`
	Network        types.Network         `yaml:"network"`
	MaxRetries     int                   `yaml:"max_retries"`
	IsRequestBatch bool                  `yaml:"is_request_batch"`
	BackoffConfig  *types.BackoffConfig  `yaml:"backoff_config"`
	RequestTimeout time.Duration         `yaml:"request_timeout"`
}
