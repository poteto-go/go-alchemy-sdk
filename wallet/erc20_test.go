package wallet

import (
	"errors"
	"math/big"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey"
	"github.com/ethereum/go-ethereum/common"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/poteto-go/go-alchemy-sdk/constant"
	"github.com/poteto-go/go-alchemy-sdk/ether"
	"github.com/poteto-go/go-alchemy-sdk/namespace"
	"github.com/poteto-go/go-alchemy-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestWallet_GetERC20Balance(t *testing.T) {
	contractAddress := "0x1234567890123456789012345678901234567890"

	t.Run("can get ERC20 balance", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		w := createConnectedWallet()
		expectedBalance := big.NewInt(1000)

		// Mock
		patches.ApplyMethod(
			reflect.TypeOf(w.erc20),
			"BalanceOf",
			func(
				_ *namespace.ERC20,
				_ string,
				_ string,
			) (*big.Int, error) {
				return expectedBalance, nil
			},
		)

		// Act
		balance, err := w.ERC20().BalanceOf(contractAddress)

		// Assert
		assert.Nil(t, err)
		assert.Equal(t, expectedBalance, balance)
	})

	t.Run("error w/o connect wallet", func(t *testing.T) {
		// Arrange
		w, _ := New(testPrivHex)

		// Act
		_, err := w.ERC20().BalanceOf(contractAddress)

		// Assert
		assert.ErrorIs(t, err, constant.ErrWalletIsNotConnected)
	})
}

func TestWallet_ERC20Transfer(t *testing.T) {
	contractAddress := "0x1234567890123456789012345678901234567890"
	otherAddress := "0xE25583099BA105D9ec0A67f5Ae86D90e50036425"
	expectedHash := common.HexToHash("0x123")

	t.Run("can transfer ERC20", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		w := createConnectedWallet()

		// Mock
		patches.ApplyMethod(
			reflect.TypeOf(w),
			"SendTransaction",
			func(_ *wallet, _ types.TransactionRequest) (common.Hash, error) {
				return expectedHash, nil
			},
		)
		patches.ApplyMethod(
			reflect.TypeOf(w.provider.Eth()),
			"WaitMined",
			func(_ *ether.Ether, _ common.Hash) (*gethTypes.Receipt, error) {
				return &gethTypes.Receipt{}, nil
			},
		)

		// Act
		_, err := w.ERC20().Transfer(
			contractAddress,
			otherAddress,
			big.NewInt(1),
			nil,
		)

		// Assert
		assert.Nil(t, err)
	})

	t.Run("handle error on transfer", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		w := createConnectedWallet()

		// Mock
		patches.ApplyMethod(
			reflect.TypeOf(w),
			"SendTransaction",
			func(_ *wallet, _ types.TransactionRequest) (common.Hash, error) {
				return common.Hash{}, errors.New("error")
			},
		)

		// Act
		_, err := w.ERC20().Transfer(
			contractAddress,
			otherAddress,
			big.NewInt(1),
			nil,
		)

		// Assert
		assert.Error(t, err)
	})

	t.Run("error w/o connect wallet", func(t *testing.T) {
		// Arrange
		w, _ := New(testPrivHex)

		// Act
		_, err := w.ERC20().Transfer(
			contractAddress,
			otherAddress,
			big.NewInt(1),
			nil,
		)

		// Assert
		assert.ErrorIs(t, err, constant.ErrWalletIsNotConnected)
	})
}

func TestWallet_ERC20TransferNoWait(t *testing.T) {
	contractAddress := "0x1234567890123456789012345678901234567890"
	otherAddress := "0xE25583099BA105D9ec0A67f5Ae86D90e50036425"
	expectedHash := common.HexToHash("0x123")

	t.Run("can transfer ERC20", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		w := createConnectedWallet()

		// Mock
		patches.ApplyMethod(
			reflect.TypeOf(w),
			"SendTransaction",
			func(_ *wallet, _ types.TransactionRequest) (common.Hash, error) {
				return expectedHash, nil
			},
		)

		// Act
		txHash, err := w.ERC20().TransferNoWait(
			contractAddress,
			otherAddress,
			big.NewInt(1),
			nil,
		)

		// Assert
		assert.Nil(t, err)
		assert.Equal(t, txHash, expectedHash)
	})

	t.Run("can transfer ERC20 w/ custom gasLimit", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		w := createConnectedWallet()

		// Mock
		patches.ApplyMethod(
			reflect.TypeOf(w),
			"SendTransaction",
			func(_ *wallet, _ types.TransactionRequest) (common.Hash, error) {
				return expectedHash, nil
			},
		)

		// Act
		txHash, err := w.ERC20().TransferNoWait(
			contractAddress,
			otherAddress,
			big.NewInt(1),
			new(uint64(1)),
		)

		// Assert
		assert.Nil(t, err)
		assert.Equal(t, txHash, expectedHash)
	})

	t.Run("error w/o connect wallet", func(t *testing.T) {
		// Arrange
		w, _ := New(testPrivHex)

		// Act
		_, err := w.ERC20().TransferNoWait(
			contractAddress,
			otherAddress,
			big.NewInt(1),
			nil,
		)

		// Assert
		assert.ErrorIs(t, err, constant.ErrWalletIsNotConnected)
	})

	t.Run("handle send tx error", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		w := createConnectedWallet()

		// Mock
		patches.ApplyMethod(
			reflect.TypeOf(w),
			"SendTransaction",
			func(_ *wallet, _ types.TransactionRequest) (common.Hash, error) {
				return common.Hash{}, errors.New("error")
			},
		)

		// Act
		_, err := w.ERC20().Transfer(
			contractAddress,
			otherAddress,
			big.NewInt(1),
			nil,
		)

		// Assert
		assert.Error(t, err)
	})
}

func TestWallet_ERC20ReadMethods(t *testing.T) {
	contractAddress := "0x1234567890123456789012345678901234567890"

	t.Run("can get total supply", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()
		
		// Arrange
		w := createConnectedWallet()
		expected := big.NewInt(1000)

		// Mock
		patches.ApplyMethod(reflect.TypeOf(w.erc20), "TotalSupply", func(_ *namespace.ERC20, _ string) (*big.Int, error) {
			return expected, nil
		})

		// Act
		res, err := w.ERC20().TotalSupply(contractAddress)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("can get allowance", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()
		
		// Arrange
		w := createConnectedWallet()
		expected := big.NewInt(500)

		// Mock
		patches.ApplyMethod(reflect.TypeOf(w.erc20), "Allowance", func(_ *namespace.ERC20, _, _, _ string) (*big.Int, error) {
			return expected, nil
		})

		// Act
		res, err := w.ERC20().Allowance(contractAddress, "owner", "spender")

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("can get name", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()
		
		// Arrange
		w := createConnectedWallet()
		expected := "TestToken"

		// Mock
		patches.ApplyMethod(reflect.TypeOf(w.erc20), "Name", func(_ *namespace.ERC20, _ string) (string, error) {
			return expected, nil
		})

		// Act
		res, err := w.ERC20().Name(contractAddress)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("can get symbol", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()
		
		// Arrange
		w := createConnectedWallet()
		expected := "TEST"

		// Mock
		patches.ApplyMethod(reflect.TypeOf(w.erc20), "Symbol", func(_ *namespace.ERC20, _ string) (string, error) {
			return expected, nil
		})

		// Act
		res, err := w.ERC20().Symbol(contractAddress)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("can get decimals", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()
		
		// Arrange
		w := createConnectedWallet()
		expected := uint8(18)

		// Mock
		patches.ApplyMethod(reflect.TypeOf(w.erc20), "Decimals", func(_ *namespace.ERC20, _ string) (uint8, error) {
			return expected, nil
		})

		// Act
		res, err := w.ERC20().Decimals(contractAddress)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})
}
