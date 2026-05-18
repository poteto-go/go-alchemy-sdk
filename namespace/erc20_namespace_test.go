package namespace_test

import (
	"fmt"
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
	// Arrange
	contractAddress := "0x1234567890abcdef1234567890abcdef12345678"
	walletAddress := "0xabcdef1234567890abcdef1234567890abcdef12"

	t.Run("can get balance of erc20", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		eth := newEtherApi()
		erc20 := namespace.NewERC20Namespace(eth)
		tmp := new(big.Int)
		expected := tmp.SetBytes([]byte("0x1"))

		// Mock
		patches.ApplyMethod(
			reflect.TypeOf(eth),
			"CallContract",
			func(
				_ *ether.Ether,
				_ ethereum.CallMsg,
				_ string,
			) ([]byte, error) {
				return []byte("0x1"), nil
			},
		)

		// Act
		balance, err := erc20.BalanceOf(contractAddress, walletAddress)
		fmt.Println(balance)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, balance.Cmp(expected), 0)
	})

	t.Run("returns error if contract call fails", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		eth := newEtherApi()
		erc20 := namespace.NewERC20Namespace(eth)

		// Mock
		patches.ApplyMethod(
			reflect.TypeOf(eth),
			"CallContract",
			func(
				_ *ether.Ether,
				_ ethereum.CallMsg,
				_ string,
			) ([]byte, error) {
				return []byte(""), assert.AnError
			},
		)

		// Act
		balance, err := erc20.BalanceOf(contractAddress, walletAddress)

		// Assert
		assert.ErrorIs(t, err, assert.AnError)
		assert.Nil(t, balance)
	})
}
