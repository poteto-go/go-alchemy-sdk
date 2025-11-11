package gas

import (
	"strconv"
	"time"

	"github.com/poteto-go/go-alchemy-sdk/ether"
	"github.com/poteto-go/go-alchemy-sdk/types"
)

type AlchemyConfig struct {
	apiKey               string
	network              types.Network
	url                  string
	maxRetries           int
	requestTimeout       time.Duration
	isRequestBatch       bool
	backoffConfig        *types.BackoffConfig
	privateNetworkConfig PrivateNetworkConfig
}

func NewAlchemyConfig(setting AlchemySetting) AlchemyConfig {
	config := AlchemyConfig{
		apiKey:               setting.ApiKey,
		network:              setting.Network,
		url:                  settingToUrl(setting),
		maxRetries:           setting.MaxRetries,
		requestTimeout:       setting.RequestTimeout,
		isRequestBatch:       setting.IsRequestBatch,
		backoffConfig:        setting.BackoffConfig,
		privateNetworkConfig: setting.PrivateNetworkConfig,
	}

	if config.requestTimeout == 0 {
		config.requestTimeout = time.Second * 10
	}

	if setting.BackoffConfig == nil {
		config.backoffConfig = &types.DefaultBackoffConfig
	}

	return config
}

func settingToUrl(setting AlchemySetting) string {
	if isPrivateNetwork(setting) {
		return "http://" + setting.PrivateNetworkConfig.Host + ":" + strconv.Itoa(setting.PrivateNetworkConfig.Port)
	}
	return "https://" + string(setting.Network) + ".g.alchemy.com/v2/" + setting.ApiKey
}

func isPrivateNetwork(setting AlchemySetting) bool {
	if setting.PrivateNetworkConfig.Port == 0 || setting.PrivateNetworkConfig.Host == "" {
		return false
	}
	return true
}

func (config *AlchemyConfig) GetUrl() string {
	return config.url
}

// To avoid circle import
func (config *AlchemyConfig) toEtherApiConfig() ether.EtherApiConfig {
	return ether.NewEtherApiConfig(
		config.url,
		config.maxRetries,
		config.requestTimeout,
		config.backoffConfig,
	)
}
