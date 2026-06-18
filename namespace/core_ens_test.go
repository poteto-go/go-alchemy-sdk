package namespace_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey"
	"github.com/poteto-go/go-alchemy-sdk/ether"
	"github.com/poteto-go/go-alchemy-sdk/namespace"
	"github.com/stretchr/testify/assert"
)

func TestCore_ResolveName(t *testing.T) {
	t.Run("delegates to EtherApi and returns result", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		api := newEtherApi()
		patches.ApplyMethod(
			reflect.TypeOf(api),
			"ResolveName",
			func(_ *ether.Ether, _ string) (string, error) {
				return "0xd8da6bf26964af9d7eed9e03e53415d37aa96045", nil
			},
		)

		core := namespace.NewCore(api)
		result, err := core.ResolveName("vitalik.eth")
		assert.NoError(t, err)
		assert.Equal(t, "0xd8da6bf26964af9d7eed9e03e53415d37aa96045", result)
	})

	t.Run("propagates error from EtherApi", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		api := newEtherApi()
		patches.ApplyMethod(
			reflect.TypeOf(api),
			"ResolveName",
			func(_ *ether.Ether, _ string) (string, error) {
				return "", errors.New("resolver not found")
			},
		)

		core := namespace.NewCore(api)
		_, err := core.ResolveName("unknown.eth")
		assert.Error(t, err)
	})
}

func TestCore_LookupAddress(t *testing.T) {
	t.Run("delegates to EtherApi and returns result", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		api := newEtherApi()
		patches.ApplyMethod(
			reflect.TypeOf(api),
			"LookupAddress",
			func(_ *ether.Ether, _ string) (string, error) {
				return "vitalik.eth", nil
			},
		)

		core := namespace.NewCore(api)
		result, err := core.LookupAddress("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045")
		assert.NoError(t, err)
		assert.Equal(t, "vitalik.eth", result)
	})

	t.Run("propagates error from EtherApi", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		api := newEtherApi()
		patches.ApplyMethod(
			reflect.TypeOf(api),
			"LookupAddress",
			func(_ *ether.Ether, _ string) (string, error) {
				return "", errors.New("name not found")
			},
		)

		core := namespace.NewCore(api)
		_, err := core.LookupAddress("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045")
		assert.Error(t, err)
	})
}
