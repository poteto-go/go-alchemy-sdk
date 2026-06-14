package wallet

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey"
	"github.com/poteto-go/go-alchemy-sdk/constant"
	"github.com/poteto-go/go-alchemy-sdk/namespace"
	"github.com/stretchr/testify/assert"
)

func TestWallet_ERC1155ReadMethods(t *testing.T) {
	contractAddress := "0x1234567890123456789012345678901234567890"
	account := "0xE25583099BA105D9ec0A67f5Ae86D90e50036425"
	tokenId := big.NewInt(1)

	t.Run("can get balance of token", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		w := createConnectedWallet()
		expected := big.NewInt(42)

		// Mock
		patches.ApplyMethod(reflect.TypeOf(w.erc1155), "BalanceOfToken", func(_ *namespace.Erc1155, _, _ string, _ *big.Int) (*big.Int, error) {
			return expected, nil
		})

		// Act
		res, err := w.ERC1155().BalanceOfToken(contractAddress, account, tokenId)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("can get batch balances", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		w := createConnectedWallet()
		expected := []*big.Int{big.NewInt(1), big.NewInt(2)}

		// Mock
		patches.ApplyMethod(reflect.TypeOf(w.erc1155), "BalanceOfBatch", func(_ *namespace.Erc1155, _ string, _ []string, _ []*big.Int) ([]*big.Int, error) {
			return expected, nil
		})

		// Act
		res, err := w.ERC1155().BalanceOfBatch(contractAddress, []string{account}, []*big.Int{tokenId})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("can get uri", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		w := createConnectedWallet()
		expected := "https://example.com/erc1155/{id}.json"

		// Mock
		patches.ApplyMethod(reflect.TypeOf(w.erc1155), "Uri", func(_ *namespace.Erc1155, _ string, _ *big.Int) (string, error) {
			return expected, nil
		})

		// Act
		res, err := w.ERC1155().Uri(contractAddress, tokenId)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("error w/o connect wallet on BalanceOf", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.ERC1155().BalanceOfToken(contractAddress, account, tokenId)

		assert.ErrorIs(t, err, constant.ErrWalletIsNotConnected)
	})

	t.Run("error w/o connect wallet on BalanceOfBatch", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.ERC1155().BalanceOfBatch(contractAddress, []string{account}, []*big.Int{tokenId})

		assert.ErrorIs(t, err, constant.ErrWalletIsNotConnected)
	})

	t.Run("error w/o connect wallet on Uri", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.ERC1155().Uri(contractAddress, tokenId)

		assert.ErrorIs(t, err, constant.ErrWalletIsNotConnected)
	})
}
