package alchemy

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
}
