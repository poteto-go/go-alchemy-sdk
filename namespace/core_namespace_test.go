package namespace_test

import (
	"errors"
	"math/big"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey"
	"github.com/poteto-go/go-alchemy-sdk/alchemy"
	"github.com/poteto-go/go-alchemy-sdk/ether"
	"github.com/poteto-go/go-alchemy-sdk/internal"
	"github.com/poteto-go/go-alchemy-sdk/namespace"
	"github.com/stretchr/testify/assert"
)

func newEtherApi() *ether.Ether {
	setting := alchemy.AlchemySetting{
		ApiKey:  "hoge",
		Network: "fuga",
		BackoffConfig: &internal.BackoffConfig{
			MaxRetries: 0,
		},
	}
	config := alchemy.NewAlchemyConfig(setting)
	provider := alchemy.NewAlchemyProvider(config)
	return ether.NewEtherApi(provider).(*ether.Ether)
}

func TestNewCore(t *testing.T) {
	// Arrange
	ether := newEtherApi()

	// Act
	core := namespace.NewCore(ether)

	// Assert
	assert.NotNil(t, core)
}

func TestCore_GetBlockNumber(t *testing.T) {
	// Arrange
	api := newEtherApi()
	core := namespace.NewCore(api).(*namespace.Core)

	t.Run("normal case:", func(t *testing.T) {
		t.Run("return block number", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Arrange
			expectedNumber := 100

			// Mock
			patches.ApplyMethod(
				reflect.TypeOf(api),
				"GetBlockNumber",
				func(_ *ether.Ether) (int, error) {
					return expectedNumber, nil
				},
			)

			// Act
			blockNumber, err := core.GetBlockNumber()

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, expectedNumber, blockNumber)
		})
	})

	t.Run("error case:", func(t *testing.T) {
		t.Run("return 0 & provider error", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Arrange
			errExpected := errors.New("error")

			// Mock
			patches.ApplyMethod(
				reflect.TypeOf(api),
				"GetBlockNumber",
				func(_ *ether.Ether) (int, error) {
					return 0, errExpected
				},
			)

			// Act
			blockNumber, err := core.GetBlockNumber()

			// Assert
			assert.ErrorIs(t, errExpected, err)
			assert.Equal(t, 0, blockNumber)
		})
	})
}

func TestCore_GetGasPrice(t *testing.T) {
	// Arrange
	api := newEtherApi()
	core := namespace.NewCore(api).(*namespace.Core)

	t.Run("normal case:", func(t *testing.T) {
		t.Run("return block number", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Arrange
			expectedNumber := 100

			// Mock
			patches.ApplyMethod(
				reflect.TypeOf(api),
				"GetGasPrice",
				func(_ *ether.Ether) (int, error) {
					return expectedNumber, nil
				},
			)

			// Act
			price, err := core.GetGasPrice()

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, expectedNumber, price)
		})
	})

	t.Run("error case:", func(t *testing.T) {
		t.Run("return 0 & provider error", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Arrange
			errExpected := errors.New("error")

			// Mock
			patches.ApplyMethod(
				reflect.TypeOf(api),
				"GetGasPrice",
				func(_ *ether.Ether) (int, error) {
					return 0, errExpected
				},
			)

			// Act
			price, err := core.GetGasPrice()

			// Assert
			assert.ErrorIs(t, errExpected, err)
			assert.Equal(t, 0, price)
		})
	})
}

func TestCore_GetBalance(t *testing.T) {
	// Arrange
	api := newEtherApi()
	core := namespace.NewCore(api).(*namespace.Core)

	t.Run("normal case:", func(t *testing.T) {
		t.Run("return balance", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Arrange
			expectedBalance := big.NewInt(1000000000000000000) // 1 ETH
			address := "0x1234567890abcdef1234567890abcdef12345678"
			blockTag := "latest"

			// Mock
			patches.ApplyMethod(
				reflect.TypeOf(api),
				"GetBalance",
				func(_ *ether.Ether, _address string, _blockTag string) (*big.Int, error) {
					assert.Equal(t, address, _address)
					assert.Equal(t, blockTag, _blockTag)
					return expectedBalance, nil
				},
			)

			// Act
			balance, err := core.GetBalance(address, blockTag)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, expectedBalance, balance)
		})
	})

	t.Run("error case:", func(t *testing.T) {
		t.Run("return 0 & provider error", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Arrange
			errExpected := errors.New("error")
			address := "0x1234567890abcdef1234567890abcdef12345678"
			blockTag := "latest"

			// Mock
			patches.ApplyMethod(
				reflect.TypeOf(api),
				"GetBalance",
				func(_ *ether.Ether, _address string, _blockTag string) (*big.Int, error) {
					assert.Equal(t, address, _address)
					assert.Equal(t, blockTag, _blockTag)
					return big.NewInt(0), errExpected
				},
			)

			// Act
			balance, err := core.GetBalance(address, blockTag)

			// Assert
			assert.ErrorIs(t, errExpected, err)
			assert.Equal(t, big.NewInt(0), balance)
		})
	})
}

func TestCore_GetCode(t *testing.T) {
	// Arrange
	api := newEtherApi()
	core := namespace.NewCore(api).(*namespace.Core)

	t.Run("normal case:", func(t *testing.T) {
		t.Run("call ether.GetCode", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Arrange
			expected := "0x123"

			// Mock
			patches.ApplyMethod(
				reflect.TypeOf(api),
				"GetCode",
				func(_ *ether.Ether, _address string, _blockTag string) (string, error) {
					return expected, nil
				},
			)

			// Act
			actual, err := core.GetCode("0x123", "latest")

			// Assert
			assert.Nil(t, err)
			assert.Equal(t, expected, actual)
		})
	})

	t.Run("error case:", func(t *testing.T) {
		t.Run("if ether occur error, return empty string & error", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Arrange
			expected := errors.New("error")

			// Mock
			patches.ApplyMethod(
				reflect.TypeOf(api),
				"GetCode",
				func(_ *ether.Ether, _address string, _blockTag string) (string, error) {
					return "", expected
				},
			)

			// Act
			actual, err := core.GetCode("0x123", "latest")

			// Assert
			assert.ErrorIs(t, expected, err)
			assert.Equal(t, "", actual)
		})
	})
}
