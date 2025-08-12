package lib_test

import (
	"math/big"
	"testing"

	"github.com/poteto-go/go-alchemy-sdk/ether/lib"
	"github.com/stretchr/testify/assert"
)

func newNetwork() *lib.Network {
	return lib.NewNetwork(
		"test-network", big.NewInt(1),
	).(*lib.Network)
}

func TestNewNetwork(t *testing.T) {
	network := newNetwork()
	assert.NotNil(t, network)
}

func TestNetwork_Matches(t *testing.T) {
	// Arrange
	network := newNetwork()

	t.Run("if other is nil, return false", func(t *testing.T) {
		// Act & Assert
		assert.False(t, network.Matches(nil))
	})

	t.Run("if other.ChainId is same, return true", func(t *testing.T) {
		// Arrange
		other := lib.NewNetwork(
			"", network.ChainId(),
		).(*lib.Network)

		// Act & Assert
		assert.True(t, network.Matches(other))
	})

	t.Run("if other.ChainId is different, return false", func(t *testing.T) {
		// Arrange
		other := lib.NewNetwork(
			"", big.NewInt(2),
		).(*lib.Network)

		// Act & Assert
		assert.False(t, network.Matches(other))
	})

	t.Run("if other.Name is same, return true", func(t *testing.T) {
		// Arrange
		other := lib.NewNetwork(
			network.Name(), nil,
		)

		// Act & Assert
		assert.True(t, network.Matches(other))
	})

	t.Run("if other.Name is different, return false", func(t *testing.T) {
		// Arrange
		other := lib.NewNetwork(
			"different-name", nil,
		)

		// Act & Assert
		assert.False(t, network.Matches(other))
	})

	t.Run("if other.ChainId is nil & Name is empty, return false", func(t *testing.T) {
		// Arrange
		other := lib.NewNetwork(
			"", nil,
		)

		// Act & Assert
		assert.False(t, network.Matches(other))
	})
}

func TestNetwork_ToJson(t *testing.T) {
	// Arrange
	network := newNetwork()

	// Act
	json, err := network.ToJson()

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, `{"name":"test-network","chainId":1}`, string(json))
}

func TestNetwork_SetterAndGetter(t *testing.T) {
	// Act
	network := newNetwork()

	t.Run("Name", func(t *testing.T) {
		// Arrange
		newName := "new-name"

		// Act
		network.SetName(newName)

		// Act & Assert
		assert.Equal(t, newName, network.Name())
	})

	t.Run("ChainId", func(t *testing.T) {
		// Arrange
		newChainId := big.NewInt(2)

		// Act
		network.SetChainId(newChainId)

		// Act & Assert
		assert.Equal(t, newChainId, network.ChainId())
	})
}
