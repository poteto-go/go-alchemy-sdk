package alchemy

import (
	"time"

	"github.com/poteto-go/go-alchemy-sdk/internal"
	"github.com/poteto-go/go-alchemy-sdk/types"
)

type AlchemyConfig struct {
	apiKey         string
	network        types.Network
	url            string
	maxRetries     int
	requestTimeout time.Duration
	isRequestBatch bool
	backoffConfig  *internal.BackoffConfig
}

func NewAlchemyConfig(setting AlchemySetting) AlchemyConfig {
	config := AlchemyConfig{
		apiKey:         setting.ApiKey,
		network:        setting.Network,
		url:            settingToUrl(setting),
		maxRetries:     setting.MaxRetries,
		requestTimeout: setting.RequestTimeout,
		isRequestBatch: setting.IsRequestBatch,
		backoffConfig:  setting.BackoffConfig,
	}

	if config.requestTimeout == 0 {
		config.requestTimeout = time.Second * 10
	}

	if setting.BackoffConfig == nil {
		config.backoffConfig = &internal.DefaultBackoffConfig
	}

	return config
}

func settingToUrl(setting AlchemySetting) string {
	return "https://" + string(setting.Network) + ".g.alchemy.com/v2/" + setting.ApiKey
}

func (config *AlchemyConfig) GetUrl() string {
	return config.url
}