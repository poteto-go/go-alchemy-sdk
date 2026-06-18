package namespace_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey"
	"github.com/poteto-go/go-alchemy-sdk/constant"
	"github.com/poteto-go/go-alchemy-sdk/ether"
	"github.com/poteto-go/go-alchemy-sdk/namespace"
	"github.com/poteto-go/go-alchemy-sdk/types"
	"github.com/stretchr/testify/assert"
)

const testENSRegistry = "0x00000000000C2E074eC69A0dFb2997BA6C7d2e1e"

func TestCore_ResolveName(t *testing.T) {
	t.Run("returns ErrENSNotSupportedOnNetwork when network has no ENS registry", func(t *testing.T) {
		api := newEtherApi()
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		patches.ApplyMethod(
			reflect.TypeOf(api),
			"Network",
			func(_ *ether.Ether) types.Network { return "unknown-network" },
		)

		core := namespace.NewCore(api)
		_, err := core.ResolveName("vitalik.eth")
		assert.ErrorIs(t, err, constant.ErrENSNotSupportedOnNetwork)
	})

	t.Run("delegates to ResolveNameBy using network registry", func(t *testing.T) {
		api := newEtherApi()
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		patches.ApplyMethod(
			reflect.TypeOf(api),
			"Network",
			func(_ *ether.Ether) types.Network { return types.EthMainnet },
		)
		patches.ApplyMethod(
			reflect.TypeOf(api),
			"ResolveNameBy",
			func(_ *ether.Ether, _ string, _ string) (string, error) {
				return "0xd8da6bf26964af9d7eed9e03e53415d37aa96045", nil
			},
		)

		core := namespace.NewCore(api)
		result, err := core.ResolveName("vitalik.eth")
		assert.NoError(t, err)
		assert.Equal(t, "0xd8da6bf26964af9d7eed9e03e53415d37aa96045", result)
	})

	t.Run("propagates error from ResolveNameBy", func(t *testing.T) {
		api := newEtherApi()
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		patches.ApplyMethod(
			reflect.TypeOf(api),
			"Network",
			func(_ *ether.Ether) types.Network { return types.EthMainnet },
		)
		patches.ApplyMethod(
			reflect.TypeOf(api),
			"ResolveNameBy",
			func(_ *ether.Ether, _ string, _ string) (string, error) {
				return "", errors.New("resolver not found")
			},
		)

		core := namespace.NewCore(api)
		_, err := core.ResolveName("unknown.eth")
		assert.Error(t, err)
	})
}

func TestCore_ResolveNameBy(t *testing.T) {
	t.Run("delegates to EtherApi with provided registry", func(t *testing.T) {
		api := newEtherApi()
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		patches.ApplyMethod(
			reflect.TypeOf(api),
			"ResolveNameBy",
			func(_ *ether.Ether, _ string, _ string) (string, error) {
				return "0xd8da6bf26964af9d7eed9e03e53415d37aa96045", nil
			},
		)

		core := namespace.NewCore(api)
		result, err := core.ResolveNameBy(testENSRegistry, "vitalik.eth")
		assert.NoError(t, err)
		assert.Equal(t, "0xd8da6bf26964af9d7eed9e03e53415d37aa96045", result)
	})

	t.Run("propagates error from EtherApi", func(t *testing.T) {
		api := newEtherApi()
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		patches.ApplyMethod(
			reflect.TypeOf(api),
			"ResolveNameBy",
			func(_ *ether.Ether, _ string, _ string) (string, error) {
				return "", errors.New("resolver not found")
			},
		)

		core := namespace.NewCore(api)
		_, err := core.ResolveNameBy(testENSRegistry, "unknown.eth")
		assert.Error(t, err)
	})
}

func TestCore_LookupAddress(t *testing.T) {
	t.Run("returns ErrENSNotSupportedOnNetwork when network has no ENS registry", func(t *testing.T) {
		api := newEtherApi()
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		patches.ApplyMethod(
			reflect.TypeOf(api),
			"Network",
			func(_ *ether.Ether) types.Network { return "unknown-network" },
		)

		core := namespace.NewCore(api)
		_, err := core.LookupAddress("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045")
		assert.ErrorIs(t, err, constant.ErrENSNotSupportedOnNetwork)
	})

	t.Run("delegates to LookupAddressBy using network registry", func(t *testing.T) {
		api := newEtherApi()
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		patches.ApplyMethod(
			reflect.TypeOf(api),
			"Network",
			func(_ *ether.Ether) types.Network { return types.EthMainnet },
		)
		patches.ApplyMethod(
			reflect.TypeOf(api),
			"LookupAddressBy",
			func(_ *ether.Ether, _ string, _ string) (string, error) {
				return "vitalik.eth", nil
			},
		)

		core := namespace.NewCore(api)
		result, err := core.LookupAddress("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045")
		assert.NoError(t, err)
		assert.Equal(t, "vitalik.eth", result)
	})
}

func TestCore_LookupAddressBy(t *testing.T) {
	t.Run("delegates to EtherApi with provided registry", func(t *testing.T) {
		api := newEtherApi()
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		patches.ApplyMethod(
			reflect.TypeOf(api),
			"LookupAddressBy",
			func(_ *ether.Ether, _ string, _ string) (string, error) {
				return "vitalik.eth", nil
			},
		)

		core := namespace.NewCore(api)
		result, err := core.LookupAddressBy(testENSRegistry, "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045")
		assert.NoError(t, err)
		assert.Equal(t, "vitalik.eth", result)
	})

	t.Run("propagates error from EtherApi", func(t *testing.T) {
		api := newEtherApi()
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		patches.ApplyMethod(
			reflect.TypeOf(api),
			"LookupAddressBy",
			func(_ *ether.Ether, _ string, _ string) (string, error) {
				return "", errors.New("name not found")
			},
		)

		core := namespace.NewCore(api)
		_, err := core.LookupAddressBy(testENSRegistry, "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045")
		assert.Error(t, err)
	})
}
