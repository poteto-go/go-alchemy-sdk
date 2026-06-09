package gas

import (
	"net/http"
	"testing"

	"github.com/poteto-go/go-alchemy-sdk/constant"
	"github.com/poteto-go/go-alchemy-sdk/ether"
	"github.com/poteto-go/go-alchemy-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestNewAlchemyConfig(t *testing.T) {
	t.Run("public mode", func(t *testing.T) {
		// Act
		config, err := NewAlchemyConfig(
			AlchemySetting{
				ApiKey:  "api-key",
				Network: types.MaticMainnet,
			},
		)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, config.apiKey, "api-key")
		assert.Equal(t, string(config.network), "matic-mainnet")
		assert.Equal(t, config.url, "https://matic-mainnet.g.alchemy.com/v2/api-key")
	})

	t.Run("p8 mode", func(t *testing.T) {
		// Act
		config, err := NewAlchemyConfig(
			AlchemySetting{
				PrivateNetworkConfig: PrivateNetworkConfig{
					Host: "127.0.0.1",
					Port: 8080,
				},
			},
		)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, config.url, "http://127.0.0.1:8080")
	})

	t.Run("use private network func", func(t *testing.T) {
		// Arrange
		env := "prod"
		networkSelector := func(setting AlchemySetting) bool {
			return env == "develop"
		}

		// Act
		config, err := NewAlchemyConfig(
			AlchemySetting{
				ApiKey:  "api-key",
				Network: types.MaticMainnet,
				PrivateNetworkConfig: PrivateNetworkConfig{
					Host: "127.0.0.1",
					Port: 8080,
				},
				IsPrivateNetwork: networkSelector,
			},
		)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, config.apiKey, "api-key")
		assert.Equal(t, string(config.network), "matic-mainnet")
		assert.Equal(t, config.url, "https://matic-mainnet.g.alchemy.com/v2/api-key")
	})
}

func TestNewAlchemyConfig_PrivateNetworkUrlValidation(t *testing.T) {
	t.Run("valid http url is accepted", func(t *testing.T) {
		_, err := NewAlchemyConfig(AlchemySetting{
			PrivateNetworkConfig: PrivateNetworkConfig{
				Url: "http://localhost:8545",
			},
		})
		assert.NoError(t, err)
	})

	t.Run("valid https url is accepted", func(t *testing.T) {
		_, err := NewAlchemyConfig(AlchemySetting{
			PrivateNetworkConfig: PrivateNetworkConfig{
				Url: "https://my-rpc.example.com",
			},
		})
		assert.NoError(t, err)
	})

	t.Run("invalid scheme returns error", func(t *testing.T) {
		_, err := NewAlchemyConfig(AlchemySetting{
			PrivateNetworkConfig: PrivateNetworkConfig{
				Url: "ftp://bad-scheme.com",
			},
		})
		assert.ErrorIs(t, err, constant.ErrInvalidPrivateNetworkUrl)
	})

	t.Run("missing scheme returns error", func(t *testing.T) {
		_, err := NewAlchemyConfig(AlchemySetting{
			PrivateNetworkConfig: PrivateNetworkConfig{
				Url: "localhost:8545",
			},
		})
		assert.ErrorIs(t, err, constant.ErrInvalidPrivateNetworkUrl)
	})

	t.Run("empty host returns error", func(t *testing.T) {
		_, err := NewAlchemyConfig(AlchemySetting{
			PrivateNetworkConfig: PrivateNetworkConfig{
				Url: "http://",
			},
		})
		assert.ErrorIs(t, err, constant.ErrInvalidPrivateNetworkUrl)
	})

	t.Run("empty url is not validated (no private network)", func(t *testing.T) {
		_, err := NewAlchemyConfig(AlchemySetting{
			ApiKey:  "api-key",
			Network: types.MaticMainnet,
		})
		assert.NoError(t, err)
	})

	t.Run("host+port path with empty host returns error", func(t *testing.T) {
		_, err := NewAlchemyConfig(AlchemySetting{
			PrivateNetworkConfig: PrivateNetworkConfig{
				Host: "",
				Port: 8545,
			},
			IsPrivateNetwork: func(AlchemySetting) bool { return true },
		})
		assert.ErrorIs(t, err, constant.ErrInvalidPrivateNetworkUrl)
	})
}

func TestNewAlchemyConfig_MaxResponseBytes(t *testing.T) {
	t.Run("defaults to DefaultMaxResponseBytes when not set", func(t *testing.T) {
		config, err := NewAlchemyConfig(AlchemySetting{
			ApiKey:  "api-key",
			Network: types.MaticMainnet,
		})

		assert.NoError(t, err)
		assert.Equal(t, types.DefaultMaxResponseBytes, config.maxResponseBytes)
	})

	t.Run("uses configured value when set", func(t *testing.T) {
		const custom int64 = 1024 * 1024 // 1 MiB
		config, err := NewAlchemyConfig(AlchemySetting{
			ApiKey:           "api-key",
			Network:          types.MaticMainnet,
			MaxResponseBytes: custom,
		})

		assert.NoError(t, err)
		assert.Equal(t, custom, config.maxResponseBytes)
	})
}

func TestNewAlchemyConfig_Transport(t *testing.T) {
	t.Run("nil when not set", func(t *testing.T) {
		config, err := NewAlchemyConfig(AlchemySetting{
			ApiKey:  "api-key",
			Network: types.MaticMainnet,
		})

		assert.NoError(t, err)
		assert.Nil(t, config.transport)
	})

	t.Run("propagates custom transport from setting", func(t *testing.T) {
		custom := &http.Transport{}
		config, err := NewAlchemyConfig(AlchemySetting{
			ApiKey:    "api-key",
			Network:   types.MaticMainnet,
			Transport: custom,
		})

		assert.NoError(t, err)
		assert.Same(t, custom, config.transport)
	})
}

func TestAlchemyConfig_GetUrl(t *testing.T) {
	t.Run("can resolve alchemy rpc url", func(t *testing.T) {
		// Arrange
		config, err := NewAlchemyConfig(
			AlchemySetting{
				ApiKey:  "api-key",
				Network: types.MaticMainnet,
			},
		)
		assert.NoError(t, err)

		// Act
		url := config.GetUrl()

		// Assert
		assert.Equal(t, url, "https://matic-mainnet.g.alchemy.com/v2/api-key")
	})

	t.Run("can resolve private (means not alchemy) rpc by url", func(t *testing.T) {
		// Arrange
		config, err := NewAlchemyConfig(
			AlchemySetting{
				PrivateNetworkConfig: PrivateNetworkConfig{
					Url: "http://custom-rpc.com",
				},
			},
		)
		assert.NoError(t, err)

		// Act
		url := config.GetUrl()

		// Assert
		assert.Equal(t, url, "http://custom-rpc.com")
	})
}

func TestAlchemyConfig_toEtherApiConfig(t *testing.T) {
	// Arrange
	config, err := NewAlchemyConfig(
		AlchemySetting{
			ApiKey:  "api-key",
			Network: types.MaticMainnet,
		},
	)
	assert.NoError(t, err)

	// Act
	etherConfig := config.toEtherApiConfig()

	// Assert
	assert.IsType(t, etherConfig, ether.EtherApiConfig{})
}
