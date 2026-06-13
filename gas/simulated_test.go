package gas

import (
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient/simulated"
	"github.com/stretchr/testify/assert"
)

func TestNewSimulated(t *testing.T) {
	t.Run("can create from simulated backend", func(t *testing.T) {
		// Arrange
		sim := simulated.NewBackend(types.GenesisAlloc{})
		defer sim.Close()

		// Act
		alchemy, err := NewSimulatedAlchemy(sim)

		// Assert
		assert.Nil(t, err)
		assert.NotNil(t, alchemy)
		assert.NotNil(t, alchemy.GetProvider())
	})

	t.Run("cannot create from nil", func(t *testing.T) {
		// Act
		_, err := NewSimulatedAlchemy(nil)

		// Assert
		assert.Error(t, err)
	})
}
