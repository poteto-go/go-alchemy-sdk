package namespace_test

import (
	"testing"

	"github.com/poteto-go/go-alchemy-sdk/namespace"
	"github.com/stretchr/testify/assert"
)

func TestNewStableCoinNamespace(t *testing.T) {
	eth := newEtherApi()

	sc := namespace.NewStableCoinNamespace(eth)

	assert.NotNil(t, sc)
}
