package namespace_test

import (
	"testing"

	"github.com/poteto-go/go-alchemy-sdk/namespace"
	"github.com/stretchr/testify/assert"
)

func TestNewTransactNamespace(t *testing.T) {
	// Arrange
	ether := newEtherApi()

	// Act
	core := namespace.NewTransactNamespace(ether)

	// Assert
	assert.NotNil(t, core)
}
