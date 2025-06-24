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
				return val >= 0 && val <= 1
			},
		)
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
