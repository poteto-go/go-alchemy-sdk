package gas

import (
	"testing"

	"github.com/poteto-go/go-alchemy-sdk/ether"
	"github.com/poteto-go/go-alchemy-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestNewAlchemyConfig(t *testing.T) {
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
