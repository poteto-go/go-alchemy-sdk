package encode_test

import (
	"testing"

	"github.com/poteto-go/go-alchemy-sdk/encode"
	"github.com/stretchr/testify/assert"
)

const testPrivateKey = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"

func TestPrivateKey(t *testing.T) {
	t.Run("valid key without 0x prefix", func(t *testing.T) {
		pk, err := encode.PrivateKey(testPrivateKey)
		assert.NoError(t, err)
		assert.NotNil(t, pk)
	})

	t.Run("valid key with 0x prefix", func(t *testing.T) {
		pk, err := encode.PrivateKey("0x" + testPrivateKey)
		assert.NoError(t, err)
		assert.NotNil(t, pk)
	})

	t.Run("invalid key returns error", func(t *testing.T) {
		_, err := encode.PrivateKey("invalid")
		assert.Error(t, err)
	})
}
