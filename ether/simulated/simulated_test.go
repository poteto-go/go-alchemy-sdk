package simulated_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	gethSimulated "github.com/ethereum/go-ethereum/ethclient/simulated"
	"github.com/stretchr/testify/assert"

	"github.com/poteto-go/go-alchemy-sdk/ether/simulated"
)

func TestNewSimulatedApi(t *testing.T) {
	t.Run("can create ether api from simulated backend", func(t *testing.T) {
		// Arrange
		sim := gethSimulated.NewBackend(types.GenesisAlloc{})
		defer sim.Close()

		// Act
		eth := simulated.NewSimulatedApi(sim)

		// Assert
		assert.NotNil(t, eth)
		assert.NotNil(t, eth.Client())
	})

	t.Run("commit is routed to the backend", func(t *testing.T) {
		// Arrange
		sim := gethSimulated.NewBackend(types.GenesisAlloc{})
		defer sim.Close()
		eth := simulated.NewSimulatedApi(sim)

		// Act
		hash, err := eth.Commit()

		// Assert
		assert.Nil(t, err)
		assert.NotEqual(t, "", hash.Hex())
	})
}

func TestNewSimulated(t *testing.T) {
	t.Run("can create from simulated backend", func(t *testing.T) {
		// Arrange
		sim := gethSimulated.NewBackend(types.GenesisAlloc{})
		defer sim.Close()

		// Act
		alchemy, err := simulated.NewSimulatedAlchemy(sim)

		// Assert
		assert.Nil(t, err)
		assert.NotNil(t, alchemy)
		assert.NotNil(t, alchemy.GetProvider())
	})

	t.Run("cannot create from nil", func(t *testing.T) {
		// Act
		_, err := simulated.NewSimulatedAlchemy(nil)

		// Assert
		assert.Error(t, err)
	})
}
