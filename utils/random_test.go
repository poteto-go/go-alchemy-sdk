package utils_test

import (
	"crypto/rand"
	"errors"
	"io"
	"math/big"
	"testing"

	"github.com/agiledragon/gomonkey"
	"github.com/poteto-go/go-alchemy-sdk/utils"
	"github.com/stretchr/testify/assert"
)

func TestRandomF64(t *testing.T) {
	t.Run("normal case:", func(t *testing.T) {
		// Act
		val := utils.RandomF64(1)

		// Assert
		assert.Condition(
			t,
			func() bool {
				return val >= 0 && val < 1
			},
		)
	})

	t.Run("returns distributed values whose mean approaches max/2", func(t *testing.T) {
		// Arrange
		const iterations = 1000
		const max = 1.0
		sum := 0.0

		// Act
		for i := 0; i < iterations; i++ {
			sum += utils.RandomF64(max)
		}
		mean := sum / float64(iterations)

		// Assert: mean should be close to max/2 (0.5) with reasonable tolerance
		assert.InDelta(t, max/2, mean, 0.05)
	})

	t.Run("scales with max", func(t *testing.T) {
		// Arrange
		const iterations = 1000
		const max = 10.0
		hasNonZero := false

		// Act & Assert: values should be in [0, max) and not all zero
		for i := 0; i < iterations; i++ {
			val := utils.RandomF64(max)
			assert.GreaterOrEqual(t, val, 0.0)
			assert.Less(t, val, max)
			if val > 0 {
				hasNonZero = true
			}
		}
		assert.True(t, hasNonZero, "expected at least one non-zero value")
	})

	t.Run("panic case:", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Mock
		patches.ApplyFunc(
			rand.Int,
			func(rand io.Reader, max *big.Int) (n *big.Int, err error) {
				return nil, errors.New("error")
			},
		)

		// Act
		assert.Panics(t, func() {
			utils.RandomF64(1)
		})
	})
}
