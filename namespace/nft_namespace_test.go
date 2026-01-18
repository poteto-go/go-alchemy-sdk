package namespace_test

import (
	"testing"

	"github.com/poteto-go/go-alchemy-sdk/namespace"
	"github.com/stretchr/testify/assert"
)

func TestNewNftNamespace(t *testing.T) {
	// Arrange
	ether := newEtherApi()

	// Act
	nft := namespace.NewNftNamespace(ether)

	// Assert
	assert.NotNil(t, nft)
}
