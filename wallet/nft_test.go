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

func TestWallet_NftReadMethods(t *testing.T) {
	contractAddress := "0x1234567890123456789012345678901234567890"
	ownerAddress := "0xE25583099BA105D9ec0A67f5Ae86D90e50036425"
	operatorAddress := "0xAbcdef1234567890abcdef1234567890AbCdEf12"
	tokenId := big.NewInt(1)

	t.Run("can get owner of token", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		w := createConnectedWallet()
		expected := "0xe25583099ba105d9ec0a67f5ae86d90e50036425"

		// Mock
		patches.ApplyMethod(reflect.TypeOf(w.nft), "OwnerOf", func(_ *namespace.Nft, _ string, _ *big.Int) (string, error) {
			return expected, nil
		})

		// Act
		res, err := w.Nft().OwnerOf(contractAddress, tokenId)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("can get token URI", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		w := createConnectedWallet()
		expected := "https://example.com/nft/1"

		// Mock
		patches.ApplyMethod(reflect.TypeOf(w.nft), "TokenURI", func(_ *namespace.Nft, _ string, _ *big.Int) (string, error) {
			return expected, nil
		})

		// Act
		res, err := w.Nft().TokenURI(contractAddress, tokenId)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("can get name", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		w := createConnectedWallet()
		expected := "TestNFT"

		// Mock
		patches.ApplyMethod(reflect.TypeOf(w.nft), "Name", func(_ *namespace.Nft, _ string) (string, error) {
			return expected, nil
		})

		// Act
		res, err := w.Nft().Name(contractAddress)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("can get symbol", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		w := createConnectedWallet()
		expected := "TNFT"

		// Mock
		patches.ApplyMethod(reflect.TypeOf(w.nft), "Symbol", func(_ *namespace.Nft, _ string) (string, error) {
			return expected, nil
		})

		// Act
		res, err := w.Nft().Symbol(contractAddress)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("can get approved address", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		w := createConnectedWallet()
		expected := "0xabcdef1234567890abcdef1234567890abcdef12"

		// Mock
		patches.ApplyMethod(reflect.TypeOf(w.nft), "GetApproved", func(_ *namespace.Nft, _ string, _ *big.Int) (string, error) {
			return expected, nil
		})

		// Act
		res, err := w.Nft().GetApproved(contractAddress, tokenId)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("can get isApprovedForAll", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		w := createConnectedWallet()

		// Mock
		patches.ApplyMethod(reflect.TypeOf(w.nft), "IsApprovedForAll", func(_ *namespace.Nft, _, _, _ string) (bool, error) {
			return true, nil
		})

		// Act
		res, err := w.Nft().IsApprovedForAll(contractAddress, ownerAddress, operatorAddress)

		// Assert
		assert.NoError(t, err)
		assert.True(t, res)
	})

	t.Run("error w/o connect wallet on OwnerOf", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.Nft().OwnerOf(contractAddress, tokenId)

		assert.ErrorIs(t, err, constant.ErrWalletIsNotConnected)
	})

	t.Run("error w/o connect wallet on TokenURI", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.Nft().TokenURI(contractAddress, tokenId)

		assert.ErrorIs(t, err, constant.ErrWalletIsNotConnected)
	})

	t.Run("error w/o connect wallet on Name", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.Nft().Name(contractAddress)

		assert.ErrorIs(t, err, constant.ErrWalletIsNotConnected)
	})

	t.Run("error w/o connect wallet on Symbol", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.Nft().Symbol(contractAddress)

		assert.ErrorIs(t, err, constant.ErrWalletIsNotConnected)
	})

	t.Run("error w/o connect wallet on GetApproved", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.Nft().GetApproved(contractAddress, tokenId)

		assert.ErrorIs(t, err, constant.ErrWalletIsNotConnected)
	})

	t.Run("error w/o connect wallet on IsApprovedForAll", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.Nft().IsApprovedForAll(contractAddress, ownerAddress, operatorAddress)

		assert.ErrorIs(t, err, constant.ErrWalletIsNotConnected)
	})
}
