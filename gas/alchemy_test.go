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
	config := NewAlchemyConfig(setting)

	// Act
	alchemy := NewAlchemy(setting)

	// Assert
	assert.NotNil(t, alchemy)
	assert.Equal(t, alchemy.config, config)
	assert.NotNil(t, alchemy.Core)
	assert.NotNil(t, alchemy.Transact)
	assert.NotNil(t, alchemy.Nft)
}

func TestAlchemy_GetProvider(t *testing.T) {
	// Arrange
	setting := AlchemySetting{
		ApiKey:  "hoge",
		Network: "fuga",
	}
	alchemy := NewAlchemy(setting)

	// Act
	provider := alchemy.GetProvider()

	// Assert
	assert.NotNil(t, provider)
}
