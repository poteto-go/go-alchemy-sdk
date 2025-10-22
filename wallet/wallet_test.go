package wallet

import (
	"crypto/ecdsa"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/poteto-go/go-alchemy-sdk/ether"
	"github.com/poteto-go/go-alchemy-sdk/gas"
	"github.com/poteto-go/go-alchemy-sdk/types"
	"github.com/stretchr/testify/assert"
)

var testAddrHex = "970e8128ab834e8eac17ab8e3812f010678cf791"
var testAddrHexTo = "970e8128ab834e8eac17ab8e3812f010678cf792"
var testPrivHex = "289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232032"

func createConnectedWallet() *wallet {
	w, _ := New(testPrivHex)

	setting := gas.AlchemySetting{
		ApiKey:  "api-key",
		Network: types.EthMainnet,
	}
	alchemy := gas.NewAlchemy(setting)

	w.Connect(alchemy.GetProvider())

	return w.(*wallet)
}

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
	expected := common.HexToAddress(testAddrHex).String()

	// Act
	w, _ := New(testPrivHex)
	addr := w.GetAddress()

	// Assert
	assert.Equal(t, expected, addr)
}

func TestWallet_Connect(t *testing.T) {
	t.Run("can set the provider to wallet", func(t *testing.T) {
		// Arrange
		setting := gas.AlchemySetting{
			ApiKey:  "api-key",
			Network: types.EthMainnet,
		}
		alchemy := gas.NewAlchemy(setting)

		w, _ := New(testPrivHex)

		// Act
		w.Connect(alchemy.GetProvider())

		// Assert
		assert.Equal(t, alchemy.GetProvider(), w.(*wallet).provider)
	})
}

func TestWallet_PendingNonceAt(t *testing.T) {
	t.Run("can get nonce from address", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		w := createConnectedWallet()

		// Mock & Assert
		patches.ApplyMethod(
			reflect.TypeOf(w.provider.Eth()),
			"PendingNonceAt",
			func(_ *ether.Ether, address string) (uint64, error) {
				assert.Equal(t, address, w.GetAddress())
				return uint64(100), nil
			},
		)

		// Act
		nonce, err := w.PendingNonceAt()

		// Assert
		assert.Nil(t, err)
		assert.Equal(t, uint64(100), nonce)
	})

	t.Run("cannot get nonce from address, return error", func(t *testing.T) {
		// Arrange
		w := createConnectedWallet()

		// Act
		nonce, err := w.PendingNonceAt()

		// Assert
		assert.Error(t, err)
		assert.Equal(t, uint64(0), nonce)
	})
}
