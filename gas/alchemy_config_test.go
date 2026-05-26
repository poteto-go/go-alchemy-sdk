package gas

import (
	"testing"

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
