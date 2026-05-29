package utils

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/assert"
)

func TestToBlockNumber(t *testing.T) {
	t.Run("if latest return nil, nil", func(t *testing.T) {
		// Act
		res, err := ToBlockNumber("latest")

		// Assert
		assert.Nil(t, res)
		assert.Nil(t, err)
	})

	t.Run("if from big hex, return bigNumber", func(t *testing.T) {
		// Act
		res, err := ToBlockNumber("0x123")

		// Assert
		assert.NotNil(t, res)
		assert.Nil(t, err)
	})

	t.Run("if failed to from bigHex, return err", func(t *testing.T) {
		// Act
		_, err := ToBlockNumber("unxpected")

		// Assert
		assert.Error(t, err)
	})

	t.Run("named tags map to rpc.BlockNumber constants", func(t *testing.T) {
		cases := []struct {
			tag      string
			expected *big.Int
		}{
			{"safe", big.NewInt(int64(rpc.SafeBlockNumber))},
			{"finalized", big.NewInt(int64(rpc.FinalizedBlockNumber))},
			{"pending", big.NewInt(int64(rpc.PendingBlockNumber))},
			{"earliest", big.NewInt(int64(rpc.EarliestBlockNumber))},
		}
		for _, c := range cases {
			t.Run(c.tag, func(t *testing.T) {
				// Act
				res, err := ToBlockNumber(c.tag)

				// Assert
				assert.Nil(t, err)
				assert.Equal(t, c.expected, res)
			})
		}
	})
}
