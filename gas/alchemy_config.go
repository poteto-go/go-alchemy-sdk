package gas

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/poteto-go/go-alchemy-sdk/ether"
	"github.com/poteto-go/go-alchemy-sdk/internal"
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
	customHeaders        []http.Header
	jwtSecret            []byte
	maxResponseBytes     int64
}

func NewAlchemyConfig(setting AlchemySetting) (AlchemyConfig, error) {
	decodedJwt, err := internal.DecodeHex(setting.PrivateNetworkConfig.JwtSecret)
	if err != nil {
		return AlchemyConfig{}, err
	}

	if len(decodedJwt) != 0 && len(decodedJwt) != 32 {
		return AlchemyConfig{}, errors.New("unexpected jwt size: required empty or 32 byte (raw 64)")
	}

	config := AlchemyConfig{
		apiKey:               setting.ApiKey,
		network:              setting.Network,
		url:                  settingToUrl(setting),
		maxRetries:           setting.MaxRetries,
		requestTimeout:       setting.RequestTimeout,
		isRequestBatch:       setting.IsRequestBatch,
		backoffConfig:        setting.BackoffConfig,
		privateNetworkConfig: setting.PrivateNetworkConfig,
		customHeaders:        setting.CustomHeaders,
		jwtSecret:            decodedJwt,
		maxResponseBytes:     setting.MaxResponseBytes,
	}

	if config.requestTimeout == 0 {
		config.requestTimeout = time.Second * 10
	}

	if setting.BackoffConfig == nil {
		config.backoffConfig = &types.DefaultBackoffConfig
	}

	if config.maxResponseBytes == 0 {
		config.maxResponseBytes = types.DefaultMaxResponseBytes
	}

	return config, nil
}

func settingToUrl(setting AlchemySetting) string {
	if isPrivateNetwork(setting) {
		return resolvePrivateNetUrl(setting)
	}
	return "https://" + string(setting.Network) + ".g.alchemy.com/v2/" + setting.ApiKey
}

func resolvePrivateNetUrl(setting AlchemySetting) string {
	if setting.PrivateNetworkConfig.Url != "" {
		return setting.PrivateNetworkConfig.Url
	}
	return "http://" + setting.PrivateNetworkConfig.Host + ":" + strconv.Itoa(setting.PrivateNetworkConfig.Port)
}

func isPrivateNetwork(setting AlchemySetting) bool {
	if setting.IsPrivateNetwork == nil {
		return defaultIsPrivateNetwork(setting)
	}
	return setting.IsPrivateNetwork(setting)
}

func defaultIsPrivateNetwork(setting AlchemySetting) bool {
	if setting.PrivateNetworkConfig.Url != "" {
		return true
	}
	return setting.PrivateNetworkConfig.Host != "" && setting.PrivateNetworkConfig.Port != 0
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
		config.customHeaders,
		config.jwtSecret,
	)
}
