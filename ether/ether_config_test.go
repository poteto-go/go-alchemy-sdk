package ether

import (
	"net/http"
	"testing"
	"time"

	"github.com/poteto-go/go-alchemy-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestNewEtherApiConfig(t *testing.T) {
	// Arrange
	url := "url"
	maxRetries := 1
	requestTimeout := time.Duration(1)
	backoffConfig := types.BackoffConfig{
		MaxRetries: 1,
	}
	customHeaders := []http.Header{
		{"hello": []string{"world"}},
	}

	// Act
	config := NewEtherApiConfig(
		url,
		maxRetries,
		requestTimeout,
		&backoffConfig,
		customHeaders,
	)

	// Assert
	assert.Equal(t, config.url, url)
	assert.Equal(t, config.maxRetries, maxRetries)
	assert.Equal(t, config.requestTimeout, requestTimeout)
	assert.Equal(t, config.backoffConfig, &backoffConfig)
	assert.Equal(t, config.customHeaders, customHeaders)
}
