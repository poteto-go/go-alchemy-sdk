package alchemy

import (
	"github.com/poteto-go/go-alchemy-sdk/internal"
	"github.com/poteto-go/go-alchemy-sdk/types"
)

type AlchemyConfig struct {
	apiKey         string
	network        types.Network
	url            string
	isRequestBatch bool
	backoffConfig  internal.BackoffConfig
}

func NewAlchemyConfig(setting AlchemySetting) AlchemyConfig {
	config := AlchemyConfig{
		apiKey:         setting.ApiKey,
		network:        setting.Network,
		url:            settingToUrl(setting),
		isRequestBatch: setting.IsRequestBatch,
	}

	if setting.BackoffConfig != nil {
		config.backoffConfig = *setting.BackoffConfig
	} else {
		config.backoffConfig = internal.DefaultBackoffConfig
	}

	return config
}

func settingToUrl(setting AlchemySetting) string {
	return "https://" + string(setting.Network) + ".g.alchemy.com/v2/" + setting.ApiKey
}

func (config *AlchemyConfig) GetUrl() string {
	return config.url
}
