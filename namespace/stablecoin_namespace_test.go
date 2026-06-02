package namespace_test

import (
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/poteto-go/go-alchemy-sdk/ether"
	"github.com/poteto-go/go-alchemy-sdk/namespace"
	"github.com/poteto-go/go-alchemy-sdk/utils"
	"github.com/stretchr/testify/assert"
)

func TestNewStableCoinNamespace(t *testing.T) {
	eth := newEtherApi()

	sc := namespace.NewStableCoinNamespace(eth)

	assert.NotNil(t, sc)
}

func TestStableCoin_IsBlacklisted(t *testing.T) {
	contractAddress := "0x1234567890abcdef1234567890abcdef12345678"
	walletAddress := "0xabcdef1234567890abcdef1234567890abcdef12"

	t.Run("returns true when address is blacklisted", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		eth := newEtherApi()
		sc := namespace.NewStableCoinNamespace(eth)
		expected := make([]byte, 32)
		expected[31] = 1

		patches.ApplyMethod(reflect.TypeOf(eth), "CallContract", func(_ *ether.Ether, _ ethereum.CallMsg, _ string) ([]byte, error) {
			return expected, nil
		})

		result, err := sc.IsBlacklisted(contractAddress, walletAddress)

		assert.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("returns false when address is not blacklisted", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		eth := newEtherApi()
		sc := namespace.NewStableCoinNamespace(eth)
		expected := make([]byte, 32)

		patches.ApplyMethod(reflect.TypeOf(eth), "CallContract", func(_ *ether.Ether, _ ethereum.CallMsg, _ string) ([]byte, error) {
			return expected, nil
		})

		result, err := sc.IsBlacklisted(contractAddress, walletAddress)

		assert.NoError(t, err)
		assert.False(t, result)
	})

	t.Run("returns error if contract call fails", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		eth := newEtherApi()
		sc := namespace.NewStableCoinNamespace(eth)

		patches.ApplyMethod(reflect.TypeOf(eth), "CallContract", func(_ *ether.Ether, _ ethereum.CallMsg, _ string) ([]byte, error) {
			return nil, assert.AnError
		})

		result, err := sc.IsBlacklisted(contractAddress, walletAddress)

		assert.Error(t, err)
		assert.False(t, result)
	})
}

func TestStableCoin_Currency(t *testing.T) {
	contractAddress := "0x1234567890abcdef1234567890abcdef12345678"

	t.Run("returns currency string", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		eth := newEtherApi()
		sc := namespace.NewStableCoinNamespace(eth)
		expected := "USD"

		patches.ApplyMethod(reflect.TypeOf(eth), "CallContract", func(_ *ether.Ether, _ ethereum.CallMsg, _ string) ([]byte, error) {
			return utils.EncodeABIString(expected), nil
		})

		result, err := sc.Currency(contractAddress)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("returns error if contract call fails", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		eth := newEtherApi()
		sc := namespace.NewStableCoinNamespace(eth)

		patches.ApplyMethod(reflect.TypeOf(eth), "CallContract", func(_ *ether.Ether, _ ethereum.CallMsg, _ string) ([]byte, error) {
			return nil, assert.AnError
		})

		result, err := sc.Currency(contractAddress)

		assert.Error(t, err)
		assert.Empty(t, result)
	})
}

func TestStableCoin_Version(t *testing.T) {
	contractAddress := "0x1234567890abcdef1234567890abcdef12345678"

	t.Run("returns version string", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		eth := newEtherApi()
		sc := namespace.NewStableCoinNamespace(eth)
		expected := "1"

		patches.ApplyMethod(reflect.TypeOf(eth), "CallContract", func(_ *ether.Ether, _ ethereum.CallMsg, _ string) ([]byte, error) {
			return utils.EncodeABIString(expected), nil
		})

		result, err := sc.Version(contractAddress)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("returns error if contract call fails", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		eth := newEtherApi()
		sc := namespace.NewStableCoinNamespace(eth)

		patches.ApplyMethod(reflect.TypeOf(eth), "CallContract", func(_ *ether.Ether, _ ethereum.CallMsg, _ string) ([]byte, error) {
			return nil, assert.AnError
		})

		result, err := sc.Version(contractAddress)

		assert.Error(t, err)
		assert.Empty(t, result)
	})
}

func TestStableCoin_Paused(t *testing.T) {
	contractAddress := "0x1234567890abcdef1234567890abcdef12345678"

	t.Run("returns true when contract is paused", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		eth := newEtherApi()
		sc := namespace.NewStableCoinNamespace(eth)
		expected := make([]byte, 32)
		expected[31] = 1

		patches.ApplyMethod(reflect.TypeOf(eth), "CallContract", func(_ *ether.Ether, _ ethereum.CallMsg, _ string) ([]byte, error) {
			return expected, nil
		})

		result, err := sc.Paused(contractAddress)

		assert.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("returns false when contract is not paused", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		eth := newEtherApi()
		sc := namespace.NewStableCoinNamespace(eth)
		expected := make([]byte, 32)

		patches.ApplyMethod(reflect.TypeOf(eth), "CallContract", func(_ *ether.Ether, _ ethereum.CallMsg, _ string) ([]byte, error) {
			return expected, nil
		})

		result, err := sc.Paused(contractAddress)

		assert.NoError(t, err)
		assert.False(t, result)
	})

	t.Run("returns error if contract call fails", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		eth := newEtherApi()
		sc := namespace.NewStableCoinNamespace(eth)

		patches.ApplyMethod(reflect.TypeOf(eth), "CallContract", func(_ *ether.Ether, _ ethereum.CallMsg, _ string) ([]byte, error) {
			return nil, assert.AnError
		})

		result, err := sc.Paused(contractAddress)

		assert.Error(t, err)
		assert.False(t, result)
	})
}

func TestStableCoin_Owner(t *testing.T) {
	contractAddress := "0x1234567890abcdef1234567890abcdef12345678"
	ownerAddress := common.HexToAddress("0xabcdef1234567890abcdef1234567890abcdef12")

	t.Run("returns owner address", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		eth := newEtherApi()
		sc := namespace.NewStableCoinNamespace(eth)
		expected := make([]byte, 32)
		copy(expected[12:], ownerAddress.Bytes())

		patches.ApplyMethod(reflect.TypeOf(eth), "CallContract", func(_ *ether.Ether, _ ethereum.CallMsg, _ string) ([]byte, error) {
			return expected, nil
		})

		result, err := sc.Owner(contractAddress)

		assert.NoError(t, err)
		assert.Equal(t, ownerAddress, result)
	})

	t.Run("returns error if contract call fails", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		eth := newEtherApi()
		sc := namespace.NewStableCoinNamespace(eth)

		patches.ApplyMethod(reflect.TypeOf(eth), "CallContract", func(_ *ether.Ether, _ ethereum.CallMsg, _ string) ([]byte, error) {
			return nil, assert.AnError
		})

		_, err := sc.Owner(contractAddress)

		assert.Error(t, err)
	})
}
