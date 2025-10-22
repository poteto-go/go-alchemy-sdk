package ether_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey"
	"github.com/ethereum/go-ethereum/ethclient"
	eth "github.com/poteto-go/go-alchemy-sdk/ether"
	"github.com/stretchr/testify/assert"
)

func TestEther_PendingNonceAt(t *testing.T) {
	t.Run("error case", func(t *testing.T) {
		t.Run("if cannot create ethClient, return err", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Arrange
			ether := newEtherApiForTest()
			address := "0x123"

			// Mock
			patches.ApplyMethod(
				reflect.TypeOf(ether),
				"GetEthClient",
				func(_ *eth.Ether) (*ethclient.Client, error) {
					return nil, errors.New("error")
				},
			)

			// Act
			_, err := ether.PendingNonceAt(address)

			// Assert
			assert.Error(t, err)
		})

		t.Run("if failed get pending nonce, return error", func(t *testing.T) {
			// Arrange
			ether := newEtherApiForTest()
			address := "0x123"

			// Act
			_, err := ether.PendingNonceAt(address)

			// Assert
			assert.Error(t, err)
		})
	})
}
