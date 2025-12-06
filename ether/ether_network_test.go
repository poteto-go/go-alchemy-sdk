package ether_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey"
	eth "github.com/poteto-go/go-alchemy-sdk/ether"
	"github.com/stretchr/testify/assert"
)

func TestEther_PeerCount(t *testing.T) {
	t.Run("normal case", func(t *testing.T) {
		t.Run("success request", func(t *testing.T) {
			// Arrange
			ether := newEtherApiForTest()
			alchemyMock := newAlchemyMockOnEtherTest(t)
			defer alchemyMock.DeactivateAndReset()

			// Mock
			alchemyMock.RegisterResponder("net_peerCount", `{"jsonrpc":"2.0","id":1,"result":"0x19"}`)

			// Act
			result, err := ether.PeerCount()

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, uint64(0x19), result)
		})
	})

	t.Run("error case", func(t *testing.T) {
		t.Run("if cannot create ethClient, return err", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Arrange
			ether := newEtherApiForTest()

			// Mock
			patches.ApplyMethod(
				reflect.TypeOf(ether),
				"SetEthClient",
				func(_ *eth.Ether) error {
					return errors.New("error")
				},
			)

			// Act
			_, err := ether.PeerCount()

			// Assert
			assert.Error(t, err)
		})

		t.Run("if failed to get peer count, return error", func(t *testing.T) {
			// Arrange
			ether := newEtherApiForTest()

			// Act
			_, err := ether.PeerCount()

			// Assert
			assert.Error(t, err)
		})
	})
}
