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

func TestERC20_BalanceOf(t *testing.T) {
	contractAddress := "0x1234567890abcdef1234567890abcdef12345678"
	walletAddress := "0xabcdef1234567890abcdef1234567890abcdef12"

	t.Run("can get balance of erc20", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		eth := newEtherApi()
		erc20 := namespace.NewERC20Namespace(eth)
		expected := big.NewInt(1)

		patches.ApplyMethod(reflect.TypeOf(eth), "CallContract", func(_ *ether.Ether, _ ethereum.CallMsg, _ string) ([]byte, error) {
			return expected.Bytes(), nil
		})

		balance, err := erc20.BalanceOf(contractAddress, walletAddress)

		assert.NoError(t, err)
		assert.Equal(t, balance.Cmp(expected), 0)
	})

	t.Run("returns error if contract call fails", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		eth := newEtherApi()
		erc20 := namespace.NewERC20Namespace(eth)

		patches.ApplyMethod(reflect.TypeOf(eth), "CallContract", func(_ *ether.Ether, _ ethereum.CallMsg, _ string) ([]byte, error) {
			return nil, assert.AnError
		})

		balance, err := erc20.BalanceOf(contractAddress, walletAddress)

		assert.Error(t, err)
		assert.Nil(t, balance)
	})
}

func TestERC20_TotalSupply(t *testing.T) {
	contractAddress := "0x1234567890abcdef1234567890abcdef12345678"

	t.Run("can get total supply", func(t *testing.T) {
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

	t.Run("returns error if fails", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()
		eth := newEtherApi()
		erc20 := namespace.NewERC20Namespace(eth)

		patches.ApplyMethod(reflect.TypeOf(eth), "CallContract", func(_ *ether.Ether, _ ethereum.CallMsg, _ string) ([]byte, error) {
			return nil, assert.AnError
		})

		_, err := erc20.TotalSupply(contractAddress)
		assert.Error(t, err)
	})
}

func TestERC20_Allowance(t *testing.T) {
	contractAddress := "0x1234567890abcdef1234567890abcdef12345678"
	owner := "0xowner"
	spender := "0xspender"

	t.Run("can get allowance", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()
		eth := newEtherApi()
		erc20 := namespace.NewERC20Namespace(eth)
		expected := big.NewInt(500)

		patches.ApplyMethod(reflect.TypeOf(eth), "CallContract", func(_ *ether.Ether, _ ethereum.CallMsg, _ string) ([]byte, error) {
			return expected.Bytes(), nil
		})

		res, err := erc20.Allowance(contractAddress, owner, spender)
		assert.NoError(t, err)
		assert.Equal(t, expected.Cmp(res), 0)
	})

	t.Run("returns error if fails", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()
		eth := newEtherApi()
		erc20 := namespace.NewERC20Namespace(eth)

		patches.ApplyMethod(reflect.TypeOf(eth), "CallContract", func(_ *ether.Ether, _ ethereum.CallMsg, _ string) ([]byte, error) {
			return nil, assert.AnError
		})

		_, err := erc20.Allowance(contractAddress, owner, spender)
		assert.Error(t, err)
	})
}

func TestERC20_Name(t *testing.T) {
	contractAddress := "0x1234567890abcdef1234567890abcdef12345678"

	t.Run("can get name", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()
		eth := newEtherApi()
		erc20 := namespace.NewERC20Namespace(eth)
		expected := "TestToken"

		// ABI encoded string: offset (32) + length + data
		encoded := make([]byte, 96)
		encoded[31] = 0x20
		encoded[63] = byte(len(expected))
		copy(encoded[64:], []byte(expected))

		patches.ApplyMethod(reflect.TypeOf(eth), "CallContract", func(_ *ether.Ether, _ ethereum.CallMsg, _ string) ([]byte, error) {
			return encoded, nil
		})

		res, err := erc20.Name(contractAddress)
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("returns error if fails", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()
		eth := newEtherApi()
		erc20 := namespace.NewERC20Namespace(eth)

		patches.ApplyMethod(reflect.TypeOf(eth), "CallContract", func(_ *ether.Ether, _ ethereum.CallMsg, _ string) ([]byte, error) {
			return nil, assert.AnError
		})

		_, err := erc20.Name(contractAddress)
		assert.Error(t, err)
	})
}

func TestERC20_Symbol(t *testing.T) {
	contractAddress := "0x1234567890abcdef1234567890abcdef12345678"

	t.Run("can get symbol", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()
		eth := newEtherApi()
		erc20 := namespace.NewERC20Namespace(eth)
		expected := "TEST"

		// ABI encoded string
		encoded := make([]byte, 96)
		encoded[31] = 0x20
		encoded[63] = byte(len(expected))
		copy(encoded[64:], []byte(expected))

		patches.ApplyMethod(reflect.TypeOf(eth), "CallContract", func(_ *ether.Ether, _ ethereum.CallMsg, _ string) ([]byte, error) {
			return encoded, nil
		})

		res, err := erc20.Symbol(contractAddress)
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("returns error if fails", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()
		eth := newEtherApi()
		erc20 := namespace.NewERC20Namespace(eth)

		patches.ApplyMethod(reflect.TypeOf(eth), "CallContract", func(_ *ether.Ether, _ ethereum.CallMsg, _ string) ([]byte, error) {
			return nil, assert.AnError
		})

		_, err := erc20.Symbol(contractAddress)
		assert.Error(t, err)
	})
}

func TestERC20_Decimals(t *testing.T) {
	contractAddress := "0x1234567890abcdef1234567890abcdef12345678"

	t.Run("can get decimals", func(t *testing.T) {
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

	t.Run("returns error if fails", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()
		eth := newEtherApi()
		erc20 := namespace.NewERC20Namespace(eth)

		patches.ApplyMethod(reflect.TypeOf(eth), "CallContract", func(_ *ether.Ether, _ ethereum.CallMsg, _ string) ([]byte, error) {
			return nil, assert.AnError
		})

		_, err := erc20.Decimals(contractAddress)
		assert.Error(t, err)
	})

	t.Run("returns error when value overflows uint8 via uint64 truncation", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()
		eth := newEtherApi()
		erc20 := namespace.NewERC20Namespace(eth)

		// 2^64+5: Uint64() silently truncates to 5, bypassing the >255 check
		val := new(big.Int).Add(new(big.Int).Lsh(big.NewInt(1), 64), big.NewInt(5))
		patches.ApplyMethod(reflect.TypeOf(eth), "CallContract", func(_ *ether.Ether, _ ethereum.CallMsg, _ string) ([]byte, error) {
			return val.Bytes(), nil
		})

		_, err := erc20.Decimals(contractAddress)
		assert.Error(t, err)
	})
}
