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
	t.Run("normal case", func(t *testing.T) {
		t.Run("success request", func(t *testing.T) {
			// Arrange
			ether := newEtherApiForTest()
			alchemyMock := newAlchemyMockOnEtherTest(t)
			defer alchemyMock.DeactivateAndReset()

			// Mock
			alchemyMock.RegisterResponder("eth_getTransactionCount", `{"jsonrpc":"2.0","id":1,"result":"0x10"}`)

			// Act
			result, err := ether.PendingNonceAt("0xa7d9ddbe1f17865597fbd27ec712455208b6b76d")

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, uint64(0x10), result)
		})
	})

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
