package namespace_test

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/poteto-go/go-alchemy-sdk/_fixture/artifacts"
	"github.com/poteto-go/go-alchemy-sdk/ether"
	"github.com/poteto-go/go-alchemy-sdk/namespace"
	"github.com/poteto-go/go-alchemy-sdk/types"
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
	contract := artifacts.NewERC20()
	contractAddress := "0x1234567890abcdef1234567890abcdef12345678"
	walletAddress := "0xabcdef1234567890abcdef1234567890abcdef12"

	t.Run("can get balance of erc20", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		eth := newEtherApi()
		erc20 := namespace.NewERC20Namespace(eth)

		// Mock
		patches.ApplyMethod(
			reflect.TypeOf(eth),
			"ContractCall",
			func(
				_ *ether.Ether,
				_ types.ContractInstance,
				_ common.Address,
				_ *bind.CallOpts,
				_ []byte,
				_ func([]byte) (any, error),
			) (any, error) {
				return big.NewInt(1), nil
			},
		)

		// Act
		balance, err := erc20.BalanceOf(contract, contractAddress, walletAddress)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, balance.Cmp(big.NewInt(1)), 0)
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
			"ContractCall",
			func(
				_ *ether.Ether,
				_ types.ContractInstance,
				_ common.Address,
				_ *bind.CallOpts,
				_ []byte,
				_ func([]byte) (any, error),
			) (any, error) {
				return nil, assert.AnError
			},
		)

		// Act
		balance, err := erc20.BalanceOf(contract, contractAddress, walletAddress)

		// Assert
		assert.ErrorIs(t, err, assert.AnError)
		assert.Nil(t, balance)
	})
}
