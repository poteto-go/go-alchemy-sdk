package alchemy_test

import (
	"testing"

	"github.com/poteto-go/go-alchemy-sdk/alchemy"
	"github.com/poteto-go/go-alchemy-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestDefineAlchemySetting(t *testing.T) {
	setting := alchemy.AlchemySetting{
		ApiKey:  "api-key",
		Network: types.MaticMainnet,
	}

	assert.Equal(t, setting.ApiKey, "api-key")
	assert.Equal(t, string(setting.Network), "matic-mainnet")
}
