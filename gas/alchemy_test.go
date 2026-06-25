package gas

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAlchemy(t *testing.T) {
	// Arrange
	setting := AlchemySetting{
		ApiKey:  "hoge",
		Network: "fuga",
	}
	config, err := NewAlchemyConfig(setting)
	assert.NoError(t, err)

	// Act
	alchemy, err := NewAlchemy(setting)
	assert.NoError(t, err)

	// Assert
	assert.NotNil(t, alchemy)
	assert.Equal(t, alchemy.config, config)
	assert.NotNil(t, alchemy.Core)
	assert.NotNil(t, alchemy.Transact)
	assert.NotNil(t, alchemy.Nft)
	assert.NotNil(t, alchemy.Debug)
}

func TestNewAlchemy_SelectsProviderByScheme(t *testing.T) {
	t.Run("http setting -> AlchemyProvider", func(t *testing.T) {
		alchemy, err := NewAlchemy(AlchemySetting{ApiKey: "hoge", Network: "fuga"})
		assert.NoError(t, err)
		assert.IsType(t, &AlchemyProvider{}, alchemy.GetProvider())
	})

	t.Run("ws setting -> WsAlchemyProvider", func(t *testing.T) {
		alchemy, err := NewAlchemy(AlchemySetting{ApiKey: "hoge", Network: "fuga", UseWebsocket: true})
		assert.NoError(t, err)
		assert.IsType(t, &WsAlchemyProvider{}, alchemy.GetProvider())
	})
}

func TestAlchemy_GetProvider(t *testing.T) {
	// Arrange
	setting := AlchemySetting{
		ApiKey:  "hoge",
		Network: "fuga",
	}
	alchemy, err := NewAlchemy(setting)
	assert.NoError(t, err)

	// Act
	provider := alchemy.GetProvider()

	// Assert
	assert.NotNil(t, provider)
}
