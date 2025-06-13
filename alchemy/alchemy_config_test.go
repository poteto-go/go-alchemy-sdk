package alchemy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAlchemyConfig(t *testing.T) {
	// Act
	config := NewAlchemyConfig(
		AlchemySetting{
			ApiKey:  "api-key",
			Network: "network",
		},
	)

	// Assert
	assert.Equal(t, config.apiKey, "api-key")
	assert.Equal(t, config.network, "network")
	assert.Equal(t, config.url, "https://network.g.alchemy.com/v2/api-key")
}

func TestAlchemyConfig_GetUrl(t *testing.T) {
	// Arrange
	config := NewAlchemyConfig(
		AlchemySetting{
			ApiKey:  "api-key",
			Network: "network",
		},
	)

	// Act
	url := config.GetUrl()

	// Assert
	assert.Equal(t, url, "https://network.g.alchemy.com/v2/api-key")
}
