package alchemy

type AlchemyConfig struct {
	apiKey  string
	network string
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
	return "https://" + setting.Network + ".g.alchemy.com/v2/" + setting.ApiKey
}

func (config *AlchemyConfig) GetUrl() string {
	return config.url
}
