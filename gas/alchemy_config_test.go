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
		config := NewAlchemyConfig(
			AlchemySetting{
				ApiKey:  "api-key",
				Network: types.MaticMainnet,
			},
		)

		// Assert
		assert.Equal(t, config.apiKey, "api-key")
		assert.Equal(t, string(config.network), "matic-mainnet")
		assert.Equal(t, config.url, "https://matic-mainnet.g.alchemy.com/v2/api-key")
	})

	t.Run("p8 mode", func(t *testing.T) {
		// Act
		config := NewAlchemyConfig(
			AlchemySetting{
				PrivateNetworkConfig: PrivateNetworkConfig{
					Host: "127.0.0.1",
					Port: 8080,
				},
			},
		)

		// Assert
		assert.Equal(t, config.url, "http://127.0.0.1:8080")
	})

	t.Run("use private network func", func(t *testing.T) {
		// Arrange
		env := "prod"
		networkSelector := func(setting AlchemySetting) bool {
			return env == "develop"
		}

		// Act
		config := NewAlchemyConfig(
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
		assert.Equal(t, config.apiKey, "api-key")
		assert.Equal(t, string(config.network), "matic-mainnet")
		assert.Equal(t, config.url, "https://matic-mainnet.g.alchemy.com/v2/api-key")
	})
}

func TestAlchemyConfig_GetUrl(t *testing.T) {
	// Arrange
	config := NewAlchemyConfig(
		AlchemySetting{
			ApiKey:  "api-key",
			Network: types.MaticMainnet,
		},
	)

	// Act
	url := config.GetUrl()

	// Assert
	assert.Equal(t, url, "https://matic-mainnet.g.alchemy.com/v2/api-key")
}

func TestAlchemyConfig_toEtherApiConfig(t *testing.T) {
	// Arrange
	config := NewAlchemyConfig(
		AlchemySetting{
			ApiKey:  "api-key",
			Network: types.MaticMainnet,
		},
	)

	// Act
	etherConfig := config.toEtherApiConfig()

	// Assert
	assert.IsType(t, etherConfig, ether.EtherApiConfig{})
}
