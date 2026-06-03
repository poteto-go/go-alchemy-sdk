package namespace_test

import (
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey"
	"github.com/ethereum/go-ethereum"
	"github.com/poteto-go/go-alchemy-sdk/ether"
	"github.com/poteto-go/go-alchemy-sdk/namespace"
	"github.com/stretchr/testify/assert"
)

func TestNewEIP2612Namespace(t *testing.T) {
	eth := newEtherApi()

	eip2612 := namespace.NewEIP2612Namespace(eth)

	assert.NotNil(t, eip2612)
}

func TestEIP2612_Nonces(t *testing.T) {
	contractAddress := "0x1234567890abcdef1234567890abcdef12345678"
	ownerAddress := "0xabcdef1234567890abcdef1234567890abcdef12"

	t.Run("returns nonce for owner", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		eth := newEtherApi()
		eip2612 := namespace.NewEIP2612Namespace(eth)
		expected := make([]byte, 32)
		expected[31] = 5

		patches.ApplyMethod(reflect.TypeOf(eth), "CallContract", func(_ *ether.Ether, _ ethereum.CallMsg, _ string) ([]byte, error) {
			return expected, nil
		})

		result, err := eip2612.Nonces(contractAddress, ownerAddress)

		assert.NoError(t, err)
		assert.Equal(t, int64(5), result.Int64())
	})

	t.Run("returns zero nonce for new owner", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		eth := newEtherApi()
		eip2612 := namespace.NewEIP2612Namespace(eth)
		expected := make([]byte, 32)

		patches.ApplyMethod(reflect.TypeOf(eth), "CallContract", func(_ *ether.Ether, _ ethereum.CallMsg, _ string) ([]byte, error) {
			return expected, nil
		})

		result, err := eip2612.Nonces(contractAddress, ownerAddress)

		assert.NoError(t, err)
		assert.Equal(t, int64(0), result.Int64())
	})

	t.Run("returns error if contract call fails", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		eth := newEtherApi()
		eip2612 := namespace.NewEIP2612Namespace(eth)

		patches.ApplyMethod(reflect.TypeOf(eth), "CallContract", func(_ *ether.Ether, _ ethereum.CallMsg, _ string) ([]byte, error) {
			return nil, assert.AnError
		})

		result, err := eip2612.Nonces(contractAddress, ownerAddress)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestEIP2612_DomainSeparator(t *testing.T) {
	contractAddress := "0x1234567890abcdef1234567890abcdef12345678"

	t.Run("returns domain separator", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		eth := newEtherApi()
		eip2612 := namespace.NewEIP2612Namespace(eth)
		expected := make([]byte, 32)
		expected[0] = 0xab
		expected[31] = 0xcd

		patches.ApplyMethod(reflect.TypeOf(eth), "CallContract", func(_ *ether.Ether, _ ethereum.CallMsg, _ string) ([]byte, error) {
			return expected, nil
		})

		result, err := eip2612.DomainSeparator(contractAddress)

		assert.NoError(t, err)
		assert.Equal(t, expected, result[:])
	})

	t.Run("returns error when output is shorter than 32 bytes", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		eth := newEtherApi()
		eip2612 := namespace.NewEIP2612Namespace(eth)

		patches.ApplyMethod(reflect.TypeOf(eth), "CallContract", func(_ *ether.Ether, _ ethereum.CallMsg, _ string) ([]byte, error) {
			return []byte{0x01, 0x02}, nil
		})

		result, err := eip2612.DomainSeparator(contractAddress)

		assert.Error(t, err)
		assert.Equal(t, [32]byte{}, result)
	})

	t.Run("returns error if contract call fails", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		eth := newEtherApi()
		eip2612 := namespace.NewEIP2612Namespace(eth)

		patches.ApplyMethod(reflect.TypeOf(eth), "CallContract", func(_ *ether.Ether, _ ethereum.CallMsg, _ string) ([]byte, error) {
			return nil, assert.AnError
		})

		result, err := eip2612.DomainSeparator(contractAddress)

		assert.Error(t, err)
		assert.Equal(t, [32]byte{}, result)
	})
}
