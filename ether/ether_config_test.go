package ether

import (
	"testing"
	"time"

	"github.com/poteto-go/go-alchemy-sdk/internal"
	"github.com/stretchr/testify/assert"
)

func TestNewEtherApiConfig(t *testing.T) {
	// Arrange
	url := "url"
	maxRetries := 1
	requestTimeout := time.Duration(1)
	backoffConfig := internal.BackoffConfig{
		MaxRetries: 1,
	}

	// Act
	config := NewEtherApiConfig(
		url,
		maxRetries,
		requestTimeout,
		&backoffConfig,
	)

	// Assert
	assert.Equal(t, config.url, url)
	assert.Equal(t, config.maxRetries, maxRetries)
	assert.Equal(t, config.requestTimeout, requestTimeout)
	assert.Equal(t, config.backoffConfig, &backoffConfig)
}
