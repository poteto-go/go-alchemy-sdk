package namespace_test

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey"
	"github.com/ethereum/go-ethereum"
	"github.com/poteto-go/go-alchemy-sdk/ether"
	"github.com/poteto-go/go-alchemy-sdk/namespace"
	"github.com/stretchr/testify/assert"
)

func TestNewERC20Namespace(t *testing.T) {
	// Arrange
	ether := newEtherApi()

	// Act
	erc20 := namespace.NewERC20Namespace(ether)

	// Assert
	assert.NotNil(t, erc20)
}

func TestERC20_ReadMethods(t *testing.T) {
	contractAddress := "0x1234567890abcdef1234567890abcdef12345678"
	
	t.Run("TotalSupply", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()
		eth := newEtherApi()
		erc20 := namespace.NewERC20Namespace(eth)
		expected := big.NewInt(1000)

		patches.ApplyMethod(reflect.TypeOf(eth), "CallContract", func(_ *ether.Ether, _ ethereum.CallMsg, _ string) ([]byte, error) {
			return expected.Bytes(), nil
		})

		res, err := erc20.TotalSupply(contractAddress)
		assert.NoError(t, err)
		assert.Equal(t, expected.Cmp(res), 0)
	})

	t.Run("Allowance", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()
		eth := newEtherApi()
		erc20 := namespace.NewERC20Namespace(eth)
		expected := big.NewInt(500)

		patches.ApplyMethod(reflect.TypeOf(eth), "CallContract", func(_ *ether.Ether, _ ethereum.CallMsg, _ string) ([]byte, error) {
			return expected.Bytes(), nil
		})

		res, err := erc20.Allowance(contractAddress, "0xowner", "0xspender")
		assert.NoError(t, err)
		assert.Equal(t, expected.Cmp(res), 0)
	})

	t.Run("Name", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()
		eth := newEtherApi()
		erc20 := namespace.NewERC20Namespace(eth)
		expected := "TestToken"

		patches.ApplyMethod(reflect.TypeOf(eth), "CallContract", func(_ *ether.Ether, _ ethereum.CallMsg, _ string) ([]byte, error) {
			return []byte(expected), nil
		})

		res, err := erc20.Name(contractAddress)
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("Symbol", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()
		eth := newEtherApi()
		erc20 := namespace.NewERC20Namespace(eth)
		expected := "TEST"

		patches.ApplyMethod(reflect.TypeOf(eth), "CallContract", func(_ *ether.Ether, _ ethereum.CallMsg, _ string) ([]byte, error) {
			return []byte(expected), nil
		})

		res, err := erc20.Symbol(contractAddress)
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("Decimals", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()
		eth := newEtherApi()
		erc20 := namespace.NewERC20Namespace(eth)
		expected := uint8(18)

		patches.ApplyMethod(reflect.TypeOf(eth), "CallContract", func(_ *ether.Ether, _ ethereum.CallMsg, _ string) ([]byte, error) {
			return big.NewInt(int64(expected)).Bytes(), nil
		})

		res, err := erc20.Decimals(contractAddress)
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})
}
