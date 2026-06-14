package namespace_test

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey"
	"github.com/ethereum/go-ethereum"
	"github.com/poteto-go/go-alchemy-sdk/constant"
	"github.com/poteto-go/go-alchemy-sdk/encode"
	"github.com/poteto-go/go-alchemy-sdk/ether"
	"github.com/poteto-go/go-alchemy-sdk/namespace"
	"github.com/stretchr/testify/assert"
)

func TestNewErc1155Namespace(t *testing.T) {
	// Arrange
	ether := newEtherApi()

	// Act
	erc1155 := namespace.NewErc1155Namespace(ether)

	// Assert
	assert.NotNil(t, erc1155)
}

func TestErc1155_BalanceOfToken(t *testing.T) {
	contractAddress := "0x1234567890abcdef1234567890abcdef12345678"
	account := "0xabcdef1234567890abcdef1234567890abcdef12"
	tokenId := big.NewInt(1)

	t.Run("can get balance of token", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		eth := newEtherApi()
		erc1155 := namespace.NewErc1155Namespace(eth)

		patches.ApplyMethod(reflect.TypeOf(eth), "CallContract", func(_ *ether.Ether, _ ethereum.CallMsg, _ string) ([]byte, error) {
			return encode.ABIUint256(big.NewInt(42)), nil
		})

		balance, err := erc1155.BalanceOfToken(contractAddress, account, tokenId)

		assert.NoError(t, err)
		assert.Equal(t, "42", balance.String())
	})

	t.Run("returns error if contract call fails", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		eth := newEtherApi()
		erc1155 := namespace.NewErc1155Namespace(eth)

		patches.ApplyMethod(reflect.TypeOf(eth), "CallContract", func(_ *ether.Ether, _ ethereum.CallMsg, _ string) ([]byte, error) {
			return nil, assert.AnError
		})

		_, err := erc1155.BalanceOfToken(contractAddress, account, tokenId)

		assert.Error(t, err)
	})

	t.Run("returns error for invalid contractAddress", func(t *testing.T) {
		eth := newEtherApi()
		erc1155 := namespace.NewErc1155Namespace(eth)

		_, err := erc1155.BalanceOfToken("invalid", account, tokenId)

		assert.ErrorIs(t, err, constant.ErrInvalidAddress)
	})

	t.Run("returns error for invalid account", func(t *testing.T) {
		eth := newEtherApi()
		erc1155 := namespace.NewErc1155Namespace(eth)

		_, err := erc1155.BalanceOfToken(contractAddress, "invalid", tokenId)

		assert.ErrorIs(t, err, constant.ErrInvalidAddress)
	})

	t.Run("returns error for nil tokenId", func(t *testing.T) {
		eth := newEtherApi()
		erc1155 := namespace.NewErc1155Namespace(eth)

		_, err := erc1155.BalanceOfToken(contractAddress, account, nil)

		assert.ErrorIs(t, err, constant.ErrNilAmount)
	})

	t.Run("returns error for negative tokenId", func(t *testing.T) {
		eth := newEtherApi()
		erc1155 := namespace.NewErc1155Namespace(eth)

		_, err := erc1155.BalanceOfToken(contractAddress, account, big.NewInt(-1))

		assert.ErrorIs(t, err, constant.ErrNegativeAmount)
	})
}

func TestErc1155_BalanceOfBatch(t *testing.T) {
	contractAddress := "0x1234567890abcdef1234567890abcdef12345678"
	accounts := []string{
		"0xabcdef1234567890abcdef1234567890abcdef12",
		"0x1234567890abcdef1234567890abcdef12345678",
	}
	tokenIds := []*big.Int{big.NewInt(1), big.NewInt(2)}

	t.Run("can get batch balances", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		eth := newEtherApi()
		erc1155 := namespace.NewErc1155Namespace(eth)

		patches.ApplyMethod(reflect.TypeOf(eth), "CallContract", func(_ *ether.Ether, _ ethereum.CallMsg, _ string) ([]byte, error) {
			// offset(0x20) + length(2) + items.
			out := encode.ABIUint256(big.NewInt(constant.ABIWordSize))
			out = append(out, encode.ABIUint256Array([]*big.Int{big.NewInt(10), big.NewInt(20)})...)
			return out, nil
		})

		balances, err := erc1155.BalanceOfBatch(contractAddress, accounts, tokenIds)

		assert.NoError(t, err)
		assert.Equal(t, 2, len(balances))
		assert.Equal(t, "10", balances[0].String())
		assert.Equal(t, "20", balances[1].String())
	})

	t.Run("returns error if contract call fails", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		eth := newEtherApi()
		erc1155 := namespace.NewErc1155Namespace(eth)

		patches.ApplyMethod(reflect.TypeOf(eth), "CallContract", func(_ *ether.Ether, _ ethereum.CallMsg, _ string) ([]byte, error) {
			return nil, assert.AnError
		})

		_, err := erc1155.BalanceOfBatch(contractAddress, accounts, tokenIds)

		assert.Error(t, err)
	})

	t.Run("returns error for invalid contractAddress", func(t *testing.T) {
		eth := newEtherApi()
		erc1155 := namespace.NewErc1155Namespace(eth)

		_, err := erc1155.BalanceOfBatch("invalid", accounts, tokenIds)

		assert.ErrorIs(t, err, constant.ErrInvalidAddress)
	})

	t.Run("returns error for mismatched lengths", func(t *testing.T) {
		eth := newEtherApi()
		erc1155 := namespace.NewErc1155Namespace(eth)

		_, err := erc1155.BalanceOfBatch(contractAddress, accounts, []*big.Int{big.NewInt(1)})

		assert.ErrorIs(t, err, constant.ErrMismatchedArrayLength)
	})

	t.Run("returns error for invalid account in slice", func(t *testing.T) {
		eth := newEtherApi()
		erc1155 := namespace.NewErc1155Namespace(eth)

		_, err := erc1155.BalanceOfBatch(contractAddress, []string{"invalid", accounts[1]}, tokenIds)

		assert.ErrorIs(t, err, constant.ErrInvalidAddress)
	})

	t.Run("returns error for nil tokenId in slice", func(t *testing.T) {
		eth := newEtherApi()
		erc1155 := namespace.NewErc1155Namespace(eth)

		_, err := erc1155.BalanceOfBatch(contractAddress, accounts, []*big.Int{big.NewInt(1), nil})

		assert.ErrorIs(t, err, constant.ErrNilAmount)
	})
}

func TestErc1155_Uri(t *testing.T) {
	contractAddress := "0x1234567890abcdef1234567890abcdef12345678"
	tokenId := big.NewInt(1)

	t.Run("can get uri", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		eth := newEtherApi()
		erc1155 := namespace.NewErc1155Namespace(eth)
		expected := "https://example.com/erc1155/{id}.json"

		patches.ApplyMethod(reflect.TypeOf(eth), "CallContract", func(_ *ether.Ether, _ ethereum.CallMsg, _ string) ([]byte, error) {
			return encode.ABIString(expected), nil
		})

		uri, err := erc1155.Uri(contractAddress, tokenId)

		assert.NoError(t, err)
		assert.Equal(t, expected, uri)
	})

	t.Run("returns error if contract call fails", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		eth := newEtherApi()
		erc1155 := namespace.NewErc1155Namespace(eth)

		patches.ApplyMethod(reflect.TypeOf(eth), "CallContract", func(_ *ether.Ether, _ ethereum.CallMsg, _ string) ([]byte, error) {
			return nil, assert.AnError
		})

		_, err := erc1155.Uri(contractAddress, tokenId)

		assert.Error(t, err)
	})

	t.Run("returns error for invalid contractAddress", func(t *testing.T) {
		eth := newEtherApi()
		erc1155 := namespace.NewErc1155Namespace(eth)

		_, err := erc1155.Uri("invalid", tokenId)

		assert.ErrorIs(t, err, constant.ErrInvalidAddress)
	})

	t.Run("returns error for nil tokenId", func(t *testing.T) {
		eth := newEtherApi()
		erc1155 := namespace.NewErc1155Namespace(eth)

		_, err := erc1155.Uri(contractAddress, nil)

		assert.ErrorIs(t, err, constant.ErrNilAmount)
	})
}
