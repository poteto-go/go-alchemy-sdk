package wallet

import (
	"crypto/ecdsa"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

var testAddrHex = "970e8128ab834e8eac17ab8e3812f010678cf791"
var testPrivHex = "289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232032"

func TestNewWallet(t *testing.T) {
	t.Run("if can hex to ECDSA, return wallet", func(t *testing.T) {
		// Arrange
		expectedP8Key, _ := crypto.HexToECDSA(testPrivHex)
		expectedPublicKey := expectedP8Key.Public().(*ecdsa.PublicKey)

		// Act
		w, err := New(testPrivHex)

		// Assert
		assert.Nil(t, err)
		assert.Equal(t, expectedP8Key, w.(*wallet).privateKey)
		assert.Equal(t, expectedPublicKey, w.(*wallet).publicKey)
	})

	t.Run("if failed hexToECDSA, return err", func(t *testing.T) {
		// Act
		_, err := New("key")

		// Assert
		assert.Error(t, err)
	})
}

func TestWallet_GetAddress(t *testing.T) {
	// Arrange
	expectedAddr := common.HexToAddress(testAddrHex)

	// Act
	w, _ := New(testPrivHex)
	addr := w.GetAddress()

	// Assert
	assert.Equal(t, expectedAddr, addr)
}
