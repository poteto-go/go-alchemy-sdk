package namespace_test

import (
	"errors"
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey"
	"github.com/ethereum/go-ethereum/common"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/poteto-go/go-alchemy-sdk/constant"
	"github.com/poteto-go/go-alchemy-sdk/ether"
	"github.com/poteto-go/go-alchemy-sdk/gas"
	"github.com/poteto-go/go-alchemy-sdk/internal"
	"github.com/poteto-go/go-alchemy-sdk/namespace"
	"github.com/poteto-go/go-alchemy-sdk/types"
	"github.com/poteto-go/go-alchemy-sdk/utils"
	"github.com/stretchr/testify/assert"
)

func newEtherApi() *ether.Ether {
	setting := gas.AlchemySetting{
		ApiKey:  "hoge",
		Network: "fuga",
		BackoffConfig: &internal.BackoffConfig{
			MaxRetries: 0,
		},
	}
	config := gas.NewAlchemyConfig(setting)
	provider := gas.NewAlchemyProvider(config)
	return ether.NewEtherApi(provider, ether.NewEtherApiConfig(
		config.GetUrl(),
		0,
		time.Duration(0),
		nil,
	)).(*ether.Ether)
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
			expectedNumber := uint64(100)

			// Mock
			patches.ApplyMethod(
				reflect.TypeOf(api),
				"BlockNumber",
				func(_ *ether.Ether) (uint64, error) {
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
				"BlockNumber",
				func(_ *ether.Ether) (uint64, error) {
					return uint64(0), errExpected
				},
			)

			// Act
			blockNumber, err := core.GetBlockNumber()

			// Assert
			assert.ErrorIs(t, errExpected, err)
			assert.Equal(t, uint64(0), blockNumber)
		})
	})
}

func TestCore_GetGasPrice(t *testing.T) {
	// Arrange
	api := newEtherApi()
	core := namespace.NewCore(api).(*namespace.Core)

	t.Run("normal case:", func(t *testing.T) {
		t.Run("return gas price", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Arrange
			expectedNumber := big.NewInt(100)

			// Mock
			patches.ApplyMethod(
				reflect.TypeOf(api),
				"GasPrice",
				func(_ *ether.Ether) (*big.Int, error) {
					return expectedNumber, nil
				},
			)

			// Act
			price, err := core.GetGasPrice()

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, price.Cmp(expectedNumber), 0)
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
				"GasPrice",
				func(_ *ether.Ether) (*big.Int, error) {
					return big.NewInt(0), errExpected
				},
			)

			// Act
			_, err := core.GetGasPrice()

			// Assert
			assert.ErrorIs(t, errExpected, err)
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

	t.Run("if arg is empty return error", func(t *testing.T) {
		// Act
		_, err := core.GetCode("", types.BlockTagOrHash{})

		// Assert
		assert.ErrorIs(t, constant.ErrInvalidArgs, err)
	})

	t.Run("blockHash case:", func(t *testing.T) {
		t.Run("call ether.CodeAtHash", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Arrange
			expected := "0x123"

			// Mock
			patches.ApplyMethod(
				reflect.TypeOf(api),
				"CodeAtHash",
				func(_ *ether.Ether, _ string, _ string) (string, error) {
					return expected, nil
				},
			)

			// Act
			actual, err := core.GetCode("0x123", types.BlockTagOrHash{
				BlockHash: "hash",
			})

			// Assert
			assert.Nil(t, err)
			assert.Equal(t, expected, actual)
		})

		t.Run("if internal error, return error", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Arrange
			expectedErr := errors.New("error")

			// Mock
			patches.ApplyMethod(
				reflect.TypeOf(api),
				"CodeAtHash",
				func(_ *ether.Ether, _ string, _ string) (string, error) {
					return "", expectedErr
				},
			)

			// Act
			actual, err := core.GetCode("0x123", types.BlockTagOrHash{
				BlockHash: "hash",
			})

			// Assert
			assert.ErrorIs(t, expectedErr, err)
			assert.Equal(t, "", actual)
		})
	})

	t.Run("blockTag case:", func(t *testing.T) {
		t.Run("call ether.GetCode", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Arrange
			expected := "0x123"

			// Mock
			patches.ApplyMethod(
				reflect.TypeOf(api),
				"CodeAt",
				func(_ *ether.Ether, _ string, _ string) (string, error) {
					return expected, nil
				},
			)

			// Act
			actual, err := core.GetCode("0x123", types.BlockTagOrHash{
				BlockTag: "latest",
			})

			// Assert
			assert.Nil(t, err)
			assert.Equal(t, expected, actual)
		})

		t.Run("if ether occur error, return empty string & error", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Arrange
			expected := errors.New("error")

			// Mock
			patches.ApplyMethod(
				reflect.TypeOf(api),
				"CodeAt",
				func(_ *ether.Ether, _ string, _ string) (string, error) {
					return "", expected
				},
			)

			// Act
			actual, err := core.GetCode("0x123", types.BlockTagOrHash{
				BlockTag: "latest",
			})

			// Assert
			assert.ErrorIs(t, expected, err)
			assert.Equal(t, "", actual)
		})
	})
}

func TestCore_IsContractAddress(t *testing.T) {
	t.Run("call core.GetCode & if error occur return false", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		api := newEtherApi()
		core := namespace.NewCore(api).(*namespace.Core)

		// Mock
		patches.ApplyMethod(
			reflect.TypeOf(api),
			"CodeAt",
			func(_ *ether.Ether, _ string, _ string) (string, error) {
				return "", errors.New("error")
			},
		)

		// Act
		isContractAddress := core.IsContractAddress("0x")

		// Assert
		assert.False(t, isContractAddress)
	})

	t.Run("call core.GetCode & if starts w/ 0x is return false", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		api := newEtherApi()
		core := namespace.NewCore(api).(*namespace.Core)

		// Mock
		patches.ApplyMethod(
			reflect.TypeOf(api),
			"CodeAt",
			func(_ *ether.Ether, _ string, _ string) (string, error) {
				return "0x", nil
			},
		)

		// Act
		isContractAddress := core.IsContractAddress("0x")

		// Assert
		assert.False(t, isContractAddress)
	})

	t.Run("call core.GetCode & if not starts w/ 0x is return true", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		api := newEtherApi()
		core := namespace.NewCore(api).(*namespace.Core)

		// Mock
		patches.ApplyMethod(
			reflect.TypeOf(api),
			"CodeAt",
			func(_ *ether.Ether, _ string, _ string) (string, error) {
				return "contract", nil
			},
		)

		// Act
		isContractAddress := core.IsContractAddress("0x")

		// Assert
		assert.True(t, isContractAddress)
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

func TestCore_GetTokenMetadata(t *testing.T) {
	// Arrange
	api := newEtherApi()
	coreNamespace := namespace.NewCore(api).(*namespace.Core)

	expected := types.TokenMetadataResponse{
		Name:     "name",
		Symbol:   "symbol",
		Decimals: 18,
		Logo:     "https://test.example.com",
	}
	t.Run("call ether.GetTokenMetadata & return result", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		callAddress := "0x123"

		// Mock & Assert
		patches.ApplyMethod(
			reflect.TypeOf(api),
			"GetTokenMetadata",
			func(_ *ether.Ether, address string) (types.TokenMetadataResponse, error) {
				assert.Equal(t, address, callAddress)
				return expected, nil
			},
		)

		// Act
		actual, _ := coreNamespace.GetTokenMetadata(callAddress)

		// Assert
		assert.Equal(t, expected, actual)
	})

	t.Run("call ether.GetTokenMetadata & return internal error", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		callAddress := "0x123"
		expectedErr := errors.New("error")

		// Mock
		patches.ApplyMethod(
			reflect.TypeOf(api),
			"GetTokenMetadata",
			func(_ *ether.Ether, _ string) (types.TokenMetadataResponse, error) {
				return types.TokenMetadataResponse{}, expectedErr
			},
		)

		// Act
		_, err := coreNamespace.GetTokenMetadata(callAddress)

		// Assert
		assert.ErrorIs(t, expectedErr, err)
	})
}

func TestCore_GetLogs(t *testing.T) {
	// Arrange
	api := newEtherApi()
	coreNamespace := namespace.NewCore(api).(*namespace.Core)
	filter := types.Filter{
		FromBlock: "0x1",
		ToBlock:   "0x2",
		Address:   "0x3",
		Topics:    []string{"0x4"},
	}

	t.Run("call ether.GetLogs & return result", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		expected := []types.LogResponse{
			{
				LogIndex:         "0x0",
				Removed:          false,
				BlockNumber:      "0x233",
				BlockHash:        "0xfc139f5e2edee9e9c888d8df9a2d2226133a9bd87c88ccbd9c930d3d4c9f9ef5",
				TransactionHash:  "0x66e7a140c8fa27fe98fde923defea7562c3ca2d6bb89798aabec65782c08f63d",
				TransactionIndex: "0x0",
				Address:          "0x42699a7612a82f1d9c36148af9c77354759b210b",
				Data:             "0x0000000000000000000000000000000000000000000000000000000000000004",
				Topics: []string{
					"0x04474795f5b996ff80cb47c148d4c5ccdbe09ef27551820caa9c2f8ed149cce3",
				},
			},
			{
				LogIndex:         "0x0",
				Removed:          false,
				BlockNumber:      "0x233",
				BlockHash:        "0xfc139f5e2edee9e9c888d8df9a2d2226133a9bd87c88ccbd9c930d3d4c9f9ef5",
				TransactionHash:  "0x66e7a140c8fa27fe98fde923defea7562c3ca2d6bb89798aabec65782c08f63d",
				TransactionIndex: "0x0",
				Address:          "0x42699a7612a82f1d9c36148af9c77354759b210b",
				Data:             "0x0000000000000000000000000000000000000000000000000000000000000004",
				Topics: []string{
					"0x04474795f5b996ff80cb47c148d4c5ccdbe09ef27551820caa9c2f8ed149cce3",
				},
			},
		}

		// Mock & Assert
		patches.ApplyMethod(
			reflect.TypeOf(api),
			"GetLogs",
			func(_ *ether.Ether, fil types.Filter) ([]types.LogResponse, error) {
				assert.Equal(t, filter, fil)
				return expected, nil
			},
		)

		// Act
		actual, _ := coreNamespace.GetLogs(filter)

		// Assert
		assert.Equal(t, expected, actual)
	})

	t.Run("call ether.GetLogs & return internal error", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		expectedErr := errors.New("error")

		// Mock
		patches.ApplyMethod(
			reflect.TypeOf(api),
			"GetLogs",
			func(_ *ether.Ether, _ types.Filter) ([]types.LogResponse, error) {
				return []types.LogResponse{}, expectedErr
			},
		)

		// Act
		_, err := coreNamespace.GetLogs(filter)

		// Assert
		assert.ErrorIs(t, expectedErr, err)
	})
}

func TestCore_EstimateGas(t *testing.T) {
	// Arrange
	api := newEtherApi()
	core := namespace.NewCore(api).(*namespace.Core)

	transaction := types.TransactionRequest{
		From:  "0x1234",
		To:    "0x2345",
		Value: "0x1",
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

func TestCore_GetTransaction(t *testing.T) {
	// Arrange
	api := newEtherApi()
	core := namespace.NewCore(api).(*namespace.Core)

	t.Run("call ether.GetTransaction & return result", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		txHash := "hash"
		expectedIsPending := false
		expected := gethTypes.NewTransaction(
			uint64(0),
			common.HexToAddress("0x1"),
			big.NewInt(0),
			uint64(0),
			big.NewInt(0),
			nil,
		)

		// Mock & Assert
		patches.ApplyMethod(
			reflect.TypeOf(api),
			"GetTransaction",
			func(_ *ether.Ether, hash string) (*gethTypes.Transaction, bool, error) {
				assert.Equal(t, hash, txHash)
				return expected, expectedIsPending, nil
			},
		)

		// Act
		actual, isPending, _ := core.GetTransaction(txHash)

		// Assert
		assert.Equal(t, expected, actual)
		assert.Equal(t, expectedIsPending, isPending)
	})

	t.Run("if error occur in ether.GetTransaction & return internal error", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		txHash := "hash"
		expectedErr := errors.New("error")

		// Mock & Assert
		patches.ApplyMethod(
			reflect.TypeOf(api),
			"GetTransaction",
			func(_ *ether.Ether, hash string) (*gethTypes.Transaction, bool, error) {
				assert.Equal(t, hash, txHash)
				return nil, false, expectedErr
			},
		)

		// Act
		_, _, err := core.GetTransaction(txHash)

		// Assert
		assert.Equal(t, expectedErr, err)
	})
}

func TestCore_Call(t *testing.T) {
	// Arrange
	api := newEtherApi()
	core := namespace.NewCore(api).(*namespace.Core)

	transaction := types.TransactionRequest{
		To:    "0x2345",
		Value: "0x1",
	}

	t.Run("call ether.Call & return result", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		expected := "0x123"

		// Mock & Assert
		patches.ApplyMethod(
			reflect.TypeOf(api),
			"Call",
			func(_ *ether.Ether, tx types.TransactionRequest, tag string) (string, error) {
				assert.Equal(t, transaction, tx)
				assert.Equal(t, "latest", tag)
				return expected, nil
			},
		)

		// Act
		actual, _ := core.Call(transaction, "latest")

		// Assert
		assert.Equal(t, expected, actual)
	})

	t.Run("if error occur in ether.Call & return internal error", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		expectedErr := errors.New("error")

		// Mock & Assert
		patches.ApplyMethod(
			reflect.TypeOf(api),
			"Call",
			func(_ *ether.Ether, tx types.TransactionRequest, tag string) (string, error) {
				assert.Equal(t, transaction, tx)
				assert.Equal(t, "latest", tag)
				return "", expectedErr
			},
		)

		// Act
		_, err := core.Call(transaction, "latest")

		// Assert
		assert.Equal(t, expectedErr, err)
	})
}

func TestCore_GetTransactionReceipt(t *testing.T) {
	// Arrange
	api := newEtherApi()
	core := namespace.NewCore(api).(*namespace.Core)

	t.Run("call ether.GetTransactionReceipt & return result", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		txHash := "0x0000000000000000000000000000000000000000000000000000000000000123"
		expected := &gethTypes.Receipt{
			TxHash: common.HexToHash(txHash),
		}

		// Mock & Assert
		patches.ApplyMethod(
			reflect.TypeOf(api),
			"GetTransactionReceipt",
			func(_ *ether.Ether, hash string) (*gethTypes.Receipt, error) {
				assert.Equal(t, hash, txHash)
				return expected, nil
			},
		)

		// Act
		actual, _ := core.GetTransactionReceipt(txHash)

		// Assert
		assert.Equal(t, txHash, actual.TxHash.Hex())
	})

	t.Run("if error occur in ether.TransactionReceipts & return internal error", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		expectedErr := errors.New("error")
		txHash := "hash"

		// Mock & Assert
		patches.ApplyMethod(
			reflect.TypeOf(api),
			"GetTransactionReceipt",
			func(_ *ether.Ether, hash string) (*gethTypes.Receipt, error) {
				assert.Equal(t, hash, txHash)
				return nil, expectedErr
			},
		)

		// Act
		_, err := core.GetTransactionReceipt(txHash)

		// Assert
		assert.Equal(t, expectedErr, err)
	})
}

func TestCore_GetTransactionReceipts(t *testing.T) {
	// Arrange
	api := newEtherApi()
	core := namespace.NewCore(api).(*namespace.Core)

	t.Run("call ether.GetTransactionReceipts & return result", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		txHash := "hash"
		txArg := types.BlockNumberOrHash{
			BlockHash: txHash,
		}
		expectedAlchemyReceipt := []types.TransactionReceipt{
			{
				TransactionHash:   "0x504ce587a65bdbdb6414a0c6c16d86a04dd79bfcc4f2950eec9634b30ce5370f",
				TransactionIndex:  "0x0",
				BlockHash:         "0xe7212a92cfb9b06addc80dec2a0dfae9ea94fd344efeb157c41e12994fcad60a",
				BlockNumber:       "0x50",
				From:              "0x627306090abab3a6e1400e9345bc60c78a8bef57",
				To:                "0xf17f52151ebef6c7334fad080c5704d77216b732",
				CumulativeGasUsed: "0x5208",
				GasUsed:           "0x5208",
				Logs:              []types.LogResponse{},
				LogsBloom:         "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
				EffectiveGasPrice: "0x1",
				Type:              "0x0",
				Status:            "0x1",
			},
		}

		expected := make([]*gethTypes.Receipt, 1)
		expected[0], _ = utils.TransformAlchemyReceiptToGeth(expectedAlchemyReceipt[0])

		// Mock & Assert
		patches.ApplyMethod(
			reflect.TypeOf(api),
			"GetTransactionReceipts",
			func(_ *ether.Ether, arg types.BlockNumberOrHash) ([]*gethTypes.Receipt, error) {
				assert.Equal(t, arg, txArg)
				return expected, nil
			},
		)

		// Act
		actual, _ := core.GetTransactionReceipts(txArg)

		// Assert
		assert.Equal(t, expected, actual)
	})

	t.Run("if error occur in ether.GetTransactionReceipts & return internal error", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		expectedErr := errors.New("error")
		txHash := "hash"
		txArg := types.BlockNumberOrHash{
			BlockHash: txHash,
		}

		// Mock & Assert
		patches.ApplyMethod(
			reflect.TypeOf(api),
			"GetTransactionReceipts",
			func(_ *ether.Ether, arg types.BlockNumberOrHash) ([]*gethTypes.Receipt, error) {
				assert.Equal(t, arg, txArg)
				return []*gethTypes.Receipt{}, expectedErr
			},
		)

		// Act
		_, err := core.GetTransactionReceipts(txArg)

		// Assert
		assert.Equal(t, expectedErr, err)
	})
}

func TestCore_GetBlock(t *testing.T) {
	// Arrange
	api := newEtherApi()
	core := namespace.NewCore(api).(*namespace.Core)

	t.Run("blockHash case:", func(t *testing.T) {
		t.Run("normal case:", func(t *testing.T) {
			t.Run("return block", func(t *testing.T) {
				patches := gomonkey.NewPatches()
				defer patches.Reset()

				// Arrange
				header := gethTypes.Header{
					Number: big.NewInt(123),
				}
				expectedBlock := gethTypes.NewBlockWithHeader(&header)

				// Mock
				patches.ApplyMethod(
					reflect.TypeOf(api),
					"GetBlockByHash",
					func(_ *ether.Ether, _ string) (*gethTypes.Block, error) {
						return expectedBlock, nil
					},
				)

				// Act
				block, err := core.GetBlock(types.BlockTagOrHash{
					BlockHash: "hash",
				})

				// Assert
				assert.NoError(t, err)
				assert.Equal(t, expectedBlock, block)
			})
		})

		t.Run("error case:", func(t *testing.T) {
			t.Run("return empty block & provider error", func(t *testing.T) {
				patches := gomonkey.NewPatches()
				defer patches.Reset()

				// Arrange
				errExpected := errors.New("error")

				// Mock
				patches.ApplyMethod(
					reflect.TypeOf(api),
					"GetBlockByHash",
					func(_ *ether.Ether, _ string) (*gethTypes.Block, error) {
						return nil, errExpected
					},
				)

				// Act
				block, err := core.GetBlock(types.BlockTagOrHash{
					BlockHash: "hash",
				})

				// Assert
				assert.ErrorIs(t, errExpected, err)
				assert.Nil(t, block)
			})
		})
	})

	t.Run("blockTag case:", func(t *testing.T) {
		t.Run("normal case:", func(t *testing.T) {
			t.Run("return block", func(t *testing.T) {
				patches := gomonkey.NewPatches()
				defer patches.Reset()

				// Arrange
				header := gethTypes.Header{
					Number: big.NewInt(123),
				}
				expectedBlock := gethTypes.NewBlockWithHeader(&header)

				// Mock
				patches.ApplyMethod(
					reflect.TypeOf(api),
					"GetBlockByNumber",
					func(_ *ether.Ether, _ string) (*gethTypes.Block, error) {
						return expectedBlock, nil
					},
				)

				// Act
				block, err := core.GetBlock(types.BlockTagOrHash{
					BlockTag: "0x123",
				})

				// Assert
				assert.NoError(t, err)
				assert.Equal(t, expectedBlock, block)
			})
		})

		t.Run("error case:", func(t *testing.T) {
			t.Run("return empty block & provider error", func(t *testing.T) {
				patches := gomonkey.NewPatches()
				defer patches.Reset()

				// Arrange
				errExpected := errors.New("error")

				// Mock
				patches.ApplyMethod(
					reflect.TypeOf(api),
					"GetBlockByNumber",
					func(_ *ether.Ether, _ string) (*gethTypes.Block, error) {
						return nil, errExpected
					},
				)

				// Act
				block, err := core.GetBlock(types.BlockTagOrHash{
					BlockTag: "0x123",
				})

				// Assert
				assert.ErrorIs(t, errExpected, err)
				assert.Nil(t, block)
			})
		})
	})

	t.Run("invalid arg case", func(t *testing.T) {
		t.Run("should return core.ErrInvalidArgs on empty args", func(t *testing.T) {
			// Act
			_, err := core.GetBlock(types.BlockTagOrHash{})

			// Assert
			assert.ErrorIs(t, err, constant.ErrInvalidArgs)
		})
	})
}
