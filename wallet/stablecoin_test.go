package wallet

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWallet_StableCoin(t *testing.T) {
	t.Run("returns WalletStableCoin", func(t *testing.T) {
		w := createConnectedWallet()

		sc := w.StableCoin()

		assert.NotNil(t, sc)
	})
}
