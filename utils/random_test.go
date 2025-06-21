package utils_test

import (
	"testing"

	"github.com/poteto-go/go-alchemy-sdk/utils"
	"github.com/stretchr/testify/assert"
)

func TestRandomF64(t *testing.T) {
	// Act
	val := utils.RandomF64(0, 1)

	// Assert
	assert.Condition(
		t,
		func() bool {
			return val >= 0 && val <= 1
		},
	)
}
