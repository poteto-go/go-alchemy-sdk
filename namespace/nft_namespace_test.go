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

func TestNewNftNamespace(t *testing.T) {
	// Arrange
	ether := newEtherApi()

	// Act
	nft := namespace.NewNftNamespace(ether)

	// Assert
	assert.NotNil(t, nft)
}

func TestNft_BalanceOf(t *testing.T) {
	contractAddress := "0x1234567890abcdef1234567890abcdef12345678"
	owner := "0xabcdef1234567890abcdef1234567890abcdef12"

	t.Run("can get balance", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		eth := newEtherApi()
		nft := namespace.NewNftNamespace(eth)

		patches.ApplyMethod(reflect.TypeOf(eth), "CallContract", func(_ *ether.Ether, _ ethereum.CallMsg, _ string) ([]byte, error) {
			return encode.ABIUint256(big.NewInt(5)), nil
		})

		balance, err := nft.BalanceOf(contractAddress, owner)

		assert.NoError(t, err)
		assert.Equal(t, "5", balance.String())
	})

	t.Run("returns error if contract call fails", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		eth := newEtherApi()
		nft := namespace.NewNftNamespace(eth)

		patches.ApplyMethod(reflect.TypeOf(eth), "CallContract", func(_ *ether.Ether, _ ethereum.CallMsg, _ string) ([]byte, error) {
			return nil, assert.AnError
		})

		_, err := nft.BalanceOf(contractAddress, owner)

		assert.Error(t, err)
	})

	t.Run("returns error for invalid contractAddress", func(t *testing.T) {
		eth := newEtherApi()
		nft := namespace.NewNftNamespace(eth)

		_, err := nft.BalanceOf("invalid", owner)

		assert.ErrorIs(t, err, constant.ErrInvalidAddress)
	})

	t.Run("returns error for invalid owner", func(t *testing.T) {
		eth := newEtherApi()
		nft := namespace.NewNftNamespace(eth)

		_, err := nft.BalanceOf(contractAddress, "invalid")

		assert.ErrorIs(t, err, constant.ErrInvalidAddress)
	})
}

func TestNft_OwnerOf(t *testing.T) {
	contractAddress := "0x1234567890abcdef1234567890abcdef12345678"
	tokenId := big.NewInt(1)
	expectedOwner := "0xabcdef1234567890abcdef1234567890abcdef12"

	t.Run("can get owner of token", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		eth := newEtherApi()
		nft := namespace.NewNftNamespace(eth)

		patches.ApplyMethod(reflect.TypeOf(eth), "CallContract", func(_ *ether.Ether, _ ethereum.CallMsg, _ string) ([]byte, error) {
			return encode.ABIAddress(expectedOwner), nil
		})

		owner, err := nft.OwnerOf(contractAddress, tokenId)

		assert.NoError(t, err)
		assert.Equal(t, expectedOwner, owner)
	})

	t.Run("returns error if contract call fails", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		eth := newEtherApi()
		nft := namespace.NewNftNamespace(eth)

		patches.ApplyMethod(reflect.TypeOf(eth), "CallContract", func(_ *ether.Ether, _ ethereum.CallMsg, _ string) ([]byte, error) {
			return nil, assert.AnError
		})

		_, err := nft.OwnerOf(contractAddress, tokenId)

		assert.Error(t, err)
	})

	t.Run("returns error for invalid contractAddress", func(t *testing.T) {
		eth := newEtherApi()
		nft := namespace.NewNftNamespace(eth)

		_, err := nft.OwnerOf("invalid", tokenId)

		assert.ErrorIs(t, err, constant.ErrInvalidAddress)
	})

	t.Run("returns error for nil tokenId", func(t *testing.T) {
		eth := newEtherApi()
		nft := namespace.NewNftNamespace(eth)

		_, err := nft.OwnerOf(contractAddress, nil)

		assert.ErrorIs(t, err, constant.ErrNilAmount)
	})

	t.Run("returns error for negative tokenId", func(t *testing.T) {
		eth := newEtherApi()
		nft := namespace.NewNftNamespace(eth)

		_, err := nft.OwnerOf(contractAddress, big.NewInt(-1))

		assert.ErrorIs(t, err, constant.ErrNegativeAmount)
	})

	t.Run("returns error for tokenId exceeding uint256", func(t *testing.T) {
		eth := newEtherApi()
		nft := namespace.NewNftNamespace(eth)

		oversized := new(big.Int).Lsh(big.NewInt(1), 256)
		_, err := nft.OwnerOf(contractAddress, oversized)

		assert.ErrorIs(t, err, constant.ErrAmountExceedsUint256)
	})
}

func TestNft_TokenURI(t *testing.T) {
	contractAddress := "0x1234567890abcdef1234567890abcdef12345678"
	tokenId := big.NewInt(1)

	t.Run("can get token URI", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		eth := newEtherApi()
		nft := namespace.NewNftNamespace(eth)
		expected := "https://example.com/token/1"

		patches.ApplyMethod(reflect.TypeOf(eth), "CallContract", func(_ *ether.Ether, _ ethereum.CallMsg, _ string) ([]byte, error) {
			return encode.ABIString(expected), nil
		})

		uri, err := nft.TokenURI(contractAddress, tokenId)

		assert.NoError(t, err)
		assert.Equal(t, expected, uri)
	})

	t.Run("returns error if contract call fails", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		eth := newEtherApi()
		nft := namespace.NewNftNamespace(eth)

		patches.ApplyMethod(reflect.TypeOf(eth), "CallContract", func(_ *ether.Ether, _ ethereum.CallMsg, _ string) ([]byte, error) {
			return nil, assert.AnError
		})

		_, err := nft.TokenURI(contractAddress, tokenId)

		assert.Error(t, err)
	})

	t.Run("returns error for invalid contractAddress", func(t *testing.T) {
		eth := newEtherApi()
		nft := namespace.NewNftNamespace(eth)

		_, err := nft.TokenURI("invalid", tokenId)

		assert.ErrorIs(t, err, constant.ErrInvalidAddress)
	})

	t.Run("returns error for nil tokenId", func(t *testing.T) {
		eth := newEtherApi()
		nft := namespace.NewNftNamespace(eth)

		_, err := nft.TokenURI(contractAddress, nil)

		assert.ErrorIs(t, err, constant.ErrNilAmount)
	})

	t.Run("returns error for negative tokenId", func(t *testing.T) {
		eth := newEtherApi()
		nft := namespace.NewNftNamespace(eth)

		_, err := nft.TokenURI(contractAddress, big.NewInt(-1))

		assert.ErrorIs(t, err, constant.ErrNegativeAmount)
	})
}

func TestNft_Name(t *testing.T) {
	contractAddress := "0x1234567890abcdef1234567890abcdef12345678"

	t.Run("can get name", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		eth := newEtherApi()
		nft := namespace.NewNftNamespace(eth)
		expected := "MyNFT"

		patches.ApplyMethod(reflect.TypeOf(eth), "CallContract", func(_ *ether.Ether, _ ethereum.CallMsg, _ string) ([]byte, error) {
			return encode.ABIString(expected), nil
		})

		name, err := nft.Name(contractAddress)

		assert.NoError(t, err)
		assert.Equal(t, expected, name)
	})

	t.Run("returns error if contract call fails", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		eth := newEtherApi()
		nft := namespace.NewNftNamespace(eth)

		patches.ApplyMethod(reflect.TypeOf(eth), "CallContract", func(_ *ether.Ether, _ ethereum.CallMsg, _ string) ([]byte, error) {
			return nil, assert.AnError
		})

		_, err := nft.Name(contractAddress)

		assert.Error(t, err)
	})

	t.Run("returns error for invalid contractAddress", func(t *testing.T) {
		eth := newEtherApi()
		nft := namespace.NewNftNamespace(eth)

		_, err := nft.Name("invalid")

		assert.ErrorIs(t, err, constant.ErrInvalidAddress)
	})
}

func TestNft_Symbol(t *testing.T) {
	contractAddress := "0x1234567890abcdef1234567890abcdef12345678"

	t.Run("can get symbol", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		eth := newEtherApi()
		nft := namespace.NewNftNamespace(eth)
		expected := "MNFT"

		patches.ApplyMethod(reflect.TypeOf(eth), "CallContract", func(_ *ether.Ether, _ ethereum.CallMsg, _ string) ([]byte, error) {
			return encode.ABIString(expected), nil
		})

		symbol, err := nft.Symbol(contractAddress)

		assert.NoError(t, err)
		assert.Equal(t, expected, symbol)
	})

	t.Run("returns error if contract call fails", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		eth := newEtherApi()
		nft := namespace.NewNftNamespace(eth)

		patches.ApplyMethod(reflect.TypeOf(eth), "CallContract", func(_ *ether.Ether, _ ethereum.CallMsg, _ string) ([]byte, error) {
			return nil, assert.AnError
		})

		_, err := nft.Symbol(contractAddress)

		assert.Error(t, err)
	})

	t.Run("returns error for invalid contractAddress", func(t *testing.T) {
		eth := newEtherApi()
		nft := namespace.NewNftNamespace(eth)

		_, err := nft.Symbol("invalid")

		assert.ErrorIs(t, err, constant.ErrInvalidAddress)
	})
}

func TestNft_GetApproved(t *testing.T) {
	contractAddress := "0x1234567890abcdef1234567890abcdef12345678"
	tokenId := big.NewInt(1)
	expectedApproved := "0xabcdef1234567890abcdef1234567890abcdef12"

	t.Run("can get approved address", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		eth := newEtherApi()
		nft := namespace.NewNftNamespace(eth)

		patches.ApplyMethod(reflect.TypeOf(eth), "CallContract", func(_ *ether.Ether, _ ethereum.CallMsg, _ string) ([]byte, error) {
			return encode.ABIAddress(expectedApproved), nil
		})

		approved, err := nft.GetApproved(contractAddress, tokenId)

		assert.NoError(t, err)
		assert.Equal(t, expectedApproved, approved)
	})

	t.Run("returns error if contract call fails", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		eth := newEtherApi()
		nft := namespace.NewNftNamespace(eth)

		patches.ApplyMethod(reflect.TypeOf(eth), "CallContract", func(_ *ether.Ether, _ ethereum.CallMsg, _ string) ([]byte, error) {
			return nil, assert.AnError
		})

		_, err := nft.GetApproved(contractAddress, tokenId)

		assert.Error(t, err)
	})

	t.Run("returns error for invalid contractAddress", func(t *testing.T) {
		eth := newEtherApi()
		nft := namespace.NewNftNamespace(eth)

		_, err := nft.GetApproved("invalid", tokenId)

		assert.ErrorIs(t, err, constant.ErrInvalidAddress)
	})

	t.Run("returns error for nil tokenId", func(t *testing.T) {
		eth := newEtherApi()
		nft := namespace.NewNftNamespace(eth)

		_, err := nft.GetApproved(contractAddress, nil)

		assert.ErrorIs(t, err, constant.ErrNilAmount)
	})

	t.Run("returns error for negative tokenId", func(t *testing.T) {
		eth := newEtherApi()
		nft := namespace.NewNftNamespace(eth)

		_, err := nft.GetApproved(contractAddress, big.NewInt(-1))

		assert.ErrorIs(t, err, constant.ErrNegativeAmount)
	})
}

func TestNft_IsApprovedForAll(t *testing.T) {
	contractAddress := "0x1234567890abcdef1234567890abcdef12345678"
	owner := "0xabcdef1234567890abcdef1234567890abcdef12"
	operator := "0x1234567890abcdef1234567890abcdef12345678"

	t.Run("can get isApprovedForAll (true)", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		eth := newEtherApi()
		nft := namespace.NewNftNamespace(eth)

		patches.ApplyMethod(reflect.TypeOf(eth), "CallContract", func(_ *ether.Ether, _ ethereum.CallMsg, _ string) ([]byte, error) {
			result := make([]byte, 32)
			result[31] = 0x01
			return result, nil
		})

		approved, err := nft.IsApprovedForAll(contractAddress, owner, operator)

		assert.NoError(t, err)
		assert.True(t, approved)
	})

	t.Run("can get isApprovedForAll (false)", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		eth := newEtherApi()
		nft := namespace.NewNftNamespace(eth)

		patches.ApplyMethod(reflect.TypeOf(eth), "CallContract", func(_ *ether.Ether, _ ethereum.CallMsg, _ string) ([]byte, error) {
			return make([]byte, 32), nil
		})

		approved, err := nft.IsApprovedForAll(contractAddress, owner, operator)

		assert.NoError(t, err)
		assert.False(t, approved)
	})

	t.Run("returns error if contract call fails", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		eth := newEtherApi()
		nft := namespace.NewNftNamespace(eth)

		patches.ApplyMethod(reflect.TypeOf(eth), "CallContract", func(_ *ether.Ether, _ ethereum.CallMsg, _ string) ([]byte, error) {
			return nil, assert.AnError
		})

		_, err := nft.IsApprovedForAll(contractAddress, owner, operator)

		assert.Error(t, err)
	})

	t.Run("returns error for invalid contractAddress", func(t *testing.T) {
		eth := newEtherApi()
		nft := namespace.NewNftNamespace(eth)

		_, err := nft.IsApprovedForAll("invalid", owner, operator)

		assert.ErrorIs(t, err, constant.ErrInvalidAddress)
	})

	t.Run("returns error for invalid owner", func(t *testing.T) {
		eth := newEtherApi()
		nft := namespace.NewNftNamespace(eth)

		_, err := nft.IsApprovedForAll(contractAddress, "invalid", operator)

		assert.ErrorIs(t, err, constant.ErrInvalidAddress)
	})

	t.Run("returns error for invalid operator", func(t *testing.T) {
		eth := newEtherApi()
		nft := namespace.NewNftNamespace(eth)

		_, err := nft.IsApprovedForAll(contractAddress, owner, "invalid")

		assert.ErrorIs(t, err, constant.ErrInvalidAddress)
	})
}
