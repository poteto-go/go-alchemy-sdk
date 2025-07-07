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
	"github.com/poteto-go/go-alchemy-sdk/types"
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

func TestCore_IsContractAddress(t *testing.T) {
	// Arrange
	api := newEtherApi()
	core := namespace.NewCore(api).(*namespace.Core)

	t.Run("call with latest & if valid code hexString, return true", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Mock & Assert
		patches.ApplyMethod(
			reflect.TypeOf(api),
			"GetCode",
			func(_ *ether.Ether, _ string, blockTag string) (string, error) {
				assert.Equal(t, "latest", blockTag)
				return "0x123", nil
			},
		)

		// Act
		actual := core.IsContractAddress("address")

		// Assert
		assert.True(t, actual)
	})

	t.Run("call with latest & if 0x, return false", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Mock & Assert
		patches.ApplyMethod(
			reflect.TypeOf(api),
			"GetCode",
			func(_ *ether.Ether, _ string, blockTag string) (string, error) {
				assert.Equal(t, "latest", blockTag)
				return "0x", nil
			},
		)

		// Act
		actual := core.IsContractAddress("address")

		// Assert
		assert.False(t, actual)
	})

	t.Run("call with latest & if error, return false", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Mock & Assert
		patches.ApplyMethod(
			reflect.TypeOf(api),
			"GetCode",
			func(_ *ether.Ether, _ string, blockTag string) (string, error) {
				assert.Equal(t, "latest", blockTag)
				return "", errors.New("error")
			},
		)

		// Act
		actual := core.IsContractAddress("address")

		// Assert
		assert.False(t, actual)
	})
}

func TestCore_GetStorageAt(t *testing.T) {
	// Arrange
	api := newEtherApi()
	core := namespace.NewCore(api).(*namespace.Core)

	t.Run("call Ether.GetStorageAt & return result", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		expected := "0x123"

		// Mock
		patches.ApplyMethod(
			reflect.TypeOf(api),
			"GetStorageAt",
			func(_ *ether.Ether, _ string, _ string, _ string) (string, error) {
				return expected, nil
			},
		)

		// Act
		actual, _ := core.GetStorageAt("address", "0", "latest")

		// Assert
		assert.Equal(t, expected, actual)
	})

	t.Run("call Ether.GetStorageAt & return internal error", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		expectedErr := errors.New("error")

		// Mock
		patches.ApplyMethod(
			reflect.TypeOf(api),
			"GetStorageAt",
			func(_ *ether.Ether, _ string, _ string, _ string) (string, error) {
				return "", expectedErr
			},
		)

		// Act
		_, err := core.GetStorageAt("address", "0", "latest")

		// Assert
		assert.ErrorIs(t, expectedErr, err)
	})
}

func TestCore_GetTokenBalances_WithOnlyAddress(t *testing.T) {
	// Arrange
	api := newEtherApi()
	coreNamespace := namespace.NewCore(api).(*namespace.Core)

	expected := types.TokenBalanceResponse{
		Address: "0x123",
		TokenBalances: []types.TokenBalance{
			{
				ContractAddress: "0x456",
				TokenBalance:    "0x0000000000000000000000000000000000000000000000000000000000000000",
				Error:           nil,
			},
			{
				ContractAddress: "0x789",
				TokenBalance:    "0x0000000000000000000000000000000000000000000000000000000000000000",
				Error:           nil,
			},
		},
	}

	t.Run("normal case:", func(t *testing.T) {
		t.Run("call Ether.GetTokenBalances with just address & return result", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Arrange
			callAddress := "0x123"

			// Mock & Assert
			patches.ApplyMethod(
				reflect.TypeOf(api),
				"GetTokenBalances",
				func(_ *ether.Ether, address string, params ...string) (types.TokenBalanceResponse, error) {
					assert.Equal(t, address, callAddress)
					assert.Equal(t, 0, len(params))
					return expected, nil
				},
			)

			// Act
			actual, _ := coreNamespace.GetTokenBalances(callAddress, nil)

			// Assert
			assert.Equal(t, expected, actual)
		})

		t.Run("if response includes error, can mapping", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Arrange
			callAddress := "0x123"
			expected := types.TokenBalanceResponse{
				Address: "0x123",
				TokenBalances: []types.TokenBalance{
					{
						ContractAddress: "0x456",
						TokenBalance:    "0x0000000000000000000000000000000000000000000000000000000000000000",
						Error:           errors.New("error"),
					},
				},
			}

			// Mock & Assert
			patches.ApplyMethod(
				reflect.TypeOf(api),
				"GetTokenBalances",
				func(_ *ether.Ether, address string, params ...string) (types.TokenBalanceResponse, error) {
					assert.Equal(t, address, callAddress)
					assert.Equal(t, 0, len(params))
					return expected, nil
				},
			)

			// Act
			actual, _ := coreNamespace.GetTokenBalances(callAddress, nil)

			// Assert
			assert.Equal(t, expected, actual)
		})
	})

	t.Run("error case:", func(t *testing.T) {
		t.Run("call `ether.GetTokenBalances` with just address & return internal error", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Arrange
			expectedErr := errors.New("error")

			// Mock
			patches.ApplyMethod(
				reflect.TypeOf(api),
				"GetTokenBalances",
				func(_ *ether.Ether, address string, params ...string) (types.TokenBalanceResponse, error) {
					return types.TokenBalanceResponse{}, expectedErr
				},
			)

			// Act
			_, err := coreNamespace.GetTokenBalances("0x123", nil)

			// Assert
			assert.Equal(t, expectedErr, err)
		})
	})
}

func TestCore_GetTokenBalances_WithAddressNContracts(t *testing.T) {
	// Arrange
	api := newEtherApi()
	core := namespace.NewCore(api).(*namespace.Core)

	t.Run("call `ether.GetTokenBalances` with address & contracts & return filtered result", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		callAddress := "0x123"
		contracts := []string{"0x456"}
		option := &types.TokenBalanceOption{
			ContractAddresses: contracts,
		}
		expected := types.TokenBalanceResponse{
			Address: "0x123",
			TokenBalances: []types.TokenBalance{
				{
					ContractAddress: "0x456",
					TokenBalance:    "0x0000000000000000000000000000000000000000000000000000000000000000",
					Error:           nil,
				},
			},
		}

		// Mock & Assert
		patches.ApplyMethod(
			reflect.TypeOf(api),
			"GetTokenBalances",
			func(_ *ether.Ether, address string, params ...string) (types.TokenBalanceResponse, error) {
				assert.Equal(t, address, callAddress)
				assert.Equal(t, contracts, params)
				return expected, nil
			},
		)

		// Act
		actual, _ := core.GetTokenBalances(callAddress, option)

		// Assert
		assert.Equal(t, expected, actual)
	})
}

func TestCore_EstimateGas(t *testing.T) {
	// Arrange
	api := newEtherApi()
	core := namespace.NewCore(api).(*namespace.Core)

	transaction := types.TransactionRequest{
		From:  "0x1234",
		To:    "0x2345",
		Value: big.NewInt(1000),
	}

	t.Run("call ether.EstimateGas & return result", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		expected := big.NewInt(1234)

		// Mock & Assert
		patches.ApplyMethod(
			reflect.TypeOf(api),
			"EstimateGas",
			func(_ *ether.Ether, tx types.TransactionRequest) (*big.Int, error) {
				assert.Equal(t, transaction, tx)
				return expected, nil
			},
		)

		// Act
		actual, _ := core.EstimateGas(transaction)

		// Assert
		assert.Equal(t, expected, actual)
	})

	t.Run("if error occur in ether.EstimateGas & return internal error", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		expectedErr := errors.New("error")

		// Mock & Assert
		patches.ApplyMethod(
			reflect.TypeOf(api),
			"EstimateGas",
			func(_ *ether.Ether, tx types.TransactionRequest) (*big.Int, error) {
				assert.Equal(t, transaction, tx)
				return big.NewInt(0), expectedErr
			},
		)

		// Act
		_, err := core.EstimateGas(transaction)

		// Assert
		assert.Equal(t, expectedErr, err)
	})
}
