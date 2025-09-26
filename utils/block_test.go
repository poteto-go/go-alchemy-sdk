package utils_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/poteto-go/go-alchemy-sdk/utils"
	"github.com/stretchr/testify/assert"
)

func TestTransformAlchemyBlock(t *testing.T) {
	// Arrange
	header := gethTypes.Header{
		TxHash: common.HexToHash("0x123"),
	}
	gethBlock := gethTypes.NewBlockWithHeader(&header)

	// Act
	block := utils.TransformAlchemyBlock(gethBlock)

	// Assert
	assert.Equal(t, block.Hash, gethBlock.Hash().Hex())

}
