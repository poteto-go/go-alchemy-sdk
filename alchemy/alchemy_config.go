package alchemy

import "github.com/poteto-go/go-alchemy-sdk/types"

type AlchemyConfig struct {
	apiKey  string
	network types.Network
	url     string
}

func NewAlchemyConfig(setting AlchemySetting) AlchemyConfig {
	return AlchemyConfig{
		apiKey:  setting.ApiKey,
		network: setting.Network,
		url:     settingToUrl(setting),
	}
}

func settingToUrl(setting AlchemySetting) string {
	return "https://" + string(setting.Network) + ".g.alchemy.com/v2/" + setting.ApiKey
}

func (config *AlchemyConfig) GetUrl() string {
	return config.url
}
