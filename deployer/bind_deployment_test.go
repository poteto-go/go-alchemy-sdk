package deployer_test

import (
	"math/big"
	"testing"

	"github.com/poteto-go/go-alchemy-sdk/_fixture/artifacts"
	"github.com/poteto-go/go-alchemy-sdk/deployer"
	"github.com/stretchr/testify/assert"
)

func Test_BindDeploymentMetadata(t *testing.T) {
	t.Run("can bind metadata", func(t *testing.T) {
		// Arrange
		metadata := &artifacts.ERC20MetaData
		binData := metadata.Bin

		// Act
		err := deployer.BindDeploymentMetadata(metadata, big.NewInt(10))

		// Assert
		assert.NoError(t, err)
		assert.NotEqual(t, binData, metadata.Bin)
	})

	t.Run("pack error", func(t *testing.T) {
		// Arrange
		metadata := &artifacts.ERC20MetaData
		metadata.Bin = "invalid"

		// Act
		err := deployer.BindDeploymentMetadata(metadata, 1)

		// Assert
		assert.Error(t, err)
	})

	t.Run("json parse error", func(t *testing.T) {
		// Arrange
		metadata := &artifacts.ERC20MetaData
		metadata.ABI = "invalid"

		// Act
		err := deployer.BindDeploymentMetadata(metadata, 1)

		// Assert
		assert.Error(t, err)
	})

	t.Run("do nothing if args are not provided", func(t *testing.T) {
		// Arrange
		metadata := &artifacts.ERC20MetaData
		binData := metadata.Bin

		// Act
		err := deployer.BindDeploymentMetadata(metadata)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, metadata.Bin, binData)
	})

	t.Run("w/o metadata is error", func(t *testing.T) {
		// Act & Assert
		assert.Error(t, deployer.BindDeploymentMetadata(nil))
	})
}
