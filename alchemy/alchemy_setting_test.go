package alchemy_test

import (
	"testing"

	"github.com/poteto-go/go-alchemy-sdk/alchemy"
	"github.com/stretchr/testify/assert"
)

func TestDefineAlchemySetting(t *testing.T) {
	setting := alchemy.AlchemySetting{
		ApiKey:  "api-key",
		Network: "network",
	}

	assert.Equal(t, setting.ApiKey, "api-key")
	assert.Equal(t, setting.Network, "network")
}
