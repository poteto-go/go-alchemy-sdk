package alchemy

import (
	"time"

	"github.com/poteto-go/go-alchemy-sdk/internal"
	"github.com/poteto-go/go-alchemy-sdk/types"
)

type AlchemySetting struct {
	ApiKey         string                  `yaml:"api_key"`
	Network        types.Network           `yaml:"network"`
	IsRequestBatch bool                    `yaml:"is_request_batch"`
	BackoffConfig  *internal.BackoffConfig `yaml:"backoff_config"`
	RequestTimeout time.Duration           `yaml:"request_timeout"`
}
