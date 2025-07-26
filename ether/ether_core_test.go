package ether_test

import (
	"errors"
	"math/big"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey"
	"github.com/go-viper/mapstructure/v2"
	"github.com/goccy/go-json"
	"github.com/jarcoal/httpmock"
	"github.com/poteto-go/go-alchemy-sdk/alchemy"
	"github.com/poteto-go/go-alchemy-sdk/core"
	"github.com/poteto-go/go-alchemy-sdk/ether"
	"github.com/poteto-go/go-alchemy-sdk/internal"
	"github.com/poteto-go/go-alchemy-sdk/types"
	"github.com/poteto-go/go-alchemy-sdk/utils"
	"github.com/stretchr/testify/assert"
)

func newEtherApiForTest() *ether.Ether {
	provider := newProviderForTest()
	return ether.NewEtherApi(provider).(*ether.Ether)
}

func newProviderForTest() *alchemy.AlchemyProvider {
	config := alchemy.NewAlchemyConfig(
		alchemy.AlchemySetting{
			ApiKey:  "hoge",
			Network: "fuga",
			BackoffConfig: &internal.BackoffConfig{
				MaxRetries: 0,
			},
		},
	)
	return alchemy.NewAlchemyProvider(config).(*alchemy.AlchemyProvider)
}

func TestEther_GetBlockNumber(t *testing.T) {
	// Arrange
	ether := newEtherApiForTest()

	t.Run("normal case", func(t *testing.T) {
		t.Run("success request", func(t *testing.T) {
			httpmock.Activate(t)
			patches := gomonkey.NewPatches()
			defer func() {
				httpmock.DeactivateAndReset()
				patches.Reset()
			}()

			// Mock
			httpmock.RegisterResponder(
				"POST",
				"https://fuga.g.alchemy.com/v2/hoge",
				httpmock.NewStringResponder(200, `{"jsonrpc":"2.0","id":1,"result":"0x1234"}`),
			)

			patches.ApplyFunc(
				utils.FromHex,
				func(s string) (int, error) {
					return 1234, nil
				},
			)
			// Act
			result, err := ether.GetBlockNumber()

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, 1234, result)
		})
	})

	t.Run("error case", func(t *testing.T) {
		t.Run("if failed to send request -> core.ErrFailedToConnect", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Act
			_, err := ether.GetBlockNumber()

			// Assert
			assert.ErrorIs(t, core.ErrFailedToConnect, err)
		})

		t.Run("if failed from hex -> error", func(t *testing.T) {
			httpmock.Activate(t)
			patches := gomonkey.NewPatches()
			defer func() {
				httpmock.DeactivateAndReset()
				patches.Reset()
			}()

			// Mock
			httpmock.RegisterResponder(
				"POST",
				"https://fuga.g.alchemy.com/v2/hoge",
				httpmock.NewStringResponder(200, `{"jsonrpc":"2.0","id":1,"result":"0x1234"}`),
			)

			patches.ApplyFunc(
				utils.FromHex,
				func(s string) (int, error) {
					return 0, core.ErrInvalidHexString
				},
			)
			// Act
			_, err := ether.GetBlockNumber()

			// Assert
			assert.ErrorIs(t, core.ErrInvalidHexString, err)
		})
	})
}

func TestEther_GetGasPrice(t *testing.T) {
	// Arrange
	ether := newEtherApiForTest()

	t.Run("normal case", func(t *testing.T) {
		t.Run("success request", func(t *testing.T) {
			httpmock.Activate(t)
			patches := gomonkey.NewPatches()
			defer func() {
				httpmock.DeactivateAndReset()
				patches.Reset()
			}()

			// Mock
			httpmock.RegisterResponder(
				"POST",
				"https://fuga.g.alchemy.com/v2/hoge",
				httpmock.NewStringResponder(200, `{"jsonrpc":"2.0","id":1,"result":"0x1234"}`),
			)

			patches.ApplyFunc(
				utils.FromHex,
				func(s string) (int, error) {
					return 1234, nil
				},
			)
			// Act
			result, err := ether.GetGasPrice()

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, 1234, result)
		})
	})

	t.Run("error case:", func(t *testing.T) {
		t.Run("if failed to send request -> core.ErrFailedToConnect", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Act
			_, err := ether.GetGasPrice()

			// Assert
			assert.ErrorIs(t, core.ErrFailedToConnect, err)
		})

		t.Run("if failed from hex -> error", func(t *testing.T) {
			httpmock.Activate(t)
			patches := gomonkey.NewPatches()
			defer func() {
				httpmock.DeactivateAndReset()
				patches.Reset()
			}()

			// Mock
			httpmock.RegisterResponder(
				"POST",
				"https://fuga.g.alchemy.com/v2/hoge",
				httpmock.NewStringResponder(200, `{"jsonrpc":"2.0","id":1,"result":"0x1234"}`),
			)

			patches.ApplyFunc(
				utils.FromHex,
				func(s string) (int, error) {
					return 0, core.ErrInvalidHexString
				},
			)
			// Act
			_, err := ether.GetGasPrice()

			// Assert
			assert.ErrorIs(t, core.ErrInvalidHexString, err)
		})
	})
}

func TestEther_GetBalance(t *testing.T) {
	// Arrange
	ether := newEtherApiForTest()

	t.Run("normal case", func(t *testing.T) {
		t.Run("success request", func(t *testing.T) {
			httpmock.Activate(t)
			patches := gomonkey.NewPatches()
			defer func() {
				httpmock.DeactivateAndReset()
				patches.Reset()
			}()

			// Mock
			httpmock.RegisterResponder(
				"POST",
				"https://fuga.g.alchemy.com/v2/hoge",
				httpmock.NewStringResponder(200, `{"jsonrpc":"2.0","id":1,"result":"0x1234"}`),
			)

			patches.ApplyFunc(
				utils.FromBigHex,
				func(s string) (*big.Int, error) {
					return big.NewInt(1234), nil
				},
			)
			// Act
			result, err := ether.GetBalance("hoge", "latest")

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, big.NewInt(1234), result)
		})
	})

	t.Run("error case:", func(t *testing.T) {
		t.Run("if failed to validate block tag -> core.ErrInvalidBlockTag", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Act
			_, err := ether.GetBalance("hoge", "unxpected")

			// Assert
			assert.ErrorIs(t, core.ErrInvalidBlockTag, err)
		})

		t.Run("if failed to send request -> core.ErrFailedToConnect", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Act
			_, err := ether.GetBalance("hoge", "latest")

			// Assert
			assert.ErrorIs(t, core.ErrFailedToConnect, err)
		})

		t.Run("if failed from hex -> error", func(t *testing.T) {
			httpmock.Activate(t)
			patches := gomonkey.NewPatches()
			defer func() {
				httpmock.DeactivateAndReset()
				patches.Reset()
			}()

			// Mock
			httpmock.RegisterResponder(
				"POST",
				"https://fuga.g.alchemy.com/v2/hoge",
				httpmock.NewStringResponder(200, `{"jsonrpc":"2.0","id":1,"result":"0x1234"}`),
			)

			patches.ApplyFunc(
				utils.FromBigHex,
				func(s string) (*big.Int, error) {
					return big.NewInt(0), core.ErrInvalidHexString
				},
			)
			// Act
			_, err := ether.GetBalance("hoge", "latest")

			// Assert
			assert.ErrorIs(t, core.ErrInvalidHexString, err)
		})
	})
}

func TestEther_GetCode(t *testing.T) {
	// Arrange
	provider := newProviderForTest()
	ether := newEtherApiForTest()

	t.Run("normal case:", func(t *testing.T) {
		t.Run("call eth_getCode & if contract exist, return hex string of code", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Arrange
			expected := "0x60806040526004361061020f57600035"

			// Mock & Assert
			patches.ApplyMethod(
				reflect.TypeOf(provider),
				"Send",
				func(_ *alchemy.AlchemyProvider, method string, _ ...string) (any, error) {
					assert.Equal(t, core.Eth_GetCode, method)
					return expected, nil
				},
			)

			// Act
			actual, err := ether.GetCode("hoge", "latest")

			// Assert
			assert.Nil(t, err)
			assert.Equal(t, expected, actual)
		})

		t.Run("call eth_getCode & if not contract exists, return 0x", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Arrange
			expected := "0x"

			// Mock & Assert
			patches.ApplyMethod(
				reflect.TypeOf(provider),
				"Send",
				func(_ *alchemy.AlchemyProvider, method string, _ ...string) (any, error) {
					assert.Equal(t, core.Eth_GetCode, method)
					return "0x", nil
				},
			)

			// Act
			actual, err := ether.GetCode("hoge", "latest")

			// Assert
			assert.Nil(t, err)
			assert.Equal(t, expected, actual)
		})
	})

	t.Run("errir case:", func(t *testing.T) {
		t.Run("if invalid blockTag provided, throw core.ErrInvalidBlockTag", func(t *testing.T) {
			// Act
			_, err := ether.GetCode("hoge", "unxpected")

			// Assert
			assert.ErrorIs(t, err, core.ErrInvalidBlockTag)
		})

		t.Run("if invalid send, throw error", func(t *testing.T) {
			// Act
			_, err := ether.GetCode("hoge", "latest")

			// Assert
			assert.Error(t, err)
		})
	})
}

func TestEther_GetTransaction(t *testing.T) {
	// Arrange
	provider := newProviderForTest()
	ether := ether.NewEtherApi(provider).(*ether.Ether)

	t.Run("normal case:", func(t *testing.T) {
		t.Run("call eth_getTransactionByHash & return transaction", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Arrange
			expectedTransaction := types.TransactionResponse{
				BlockHash:   "0x1",
				Index:       1,
				BlockNumber: 2,
				From:        "0x3",
				To:          "0x4",
				GasLimit:    big.NewInt(5),
				GasPrice:    big.NewInt(6),
				Hash:        "0x7",
				Data:        "0x8",
				Nonce:       9,
				Type:        17,
				Value:       big.NewInt(10),
				ChainID:     11,
				Signature: types.Signature{
					R: "0xd",
					S: "0xe",
					V: big.NewInt(12),
				},
				MaxPriorityFeePerGas: big.NewInt(0),
				MaxFeePerGas:         big.NewInt(0),
				AccessList:           []string{},
				BlobVersionedHashes:  []string{},
				AuthorizationList:    []string{},
			}

			// Mock & Assert
			patches.ApplyMethod(
				reflect.TypeOf(provider),
				"Send",
				func(_ *alchemy.AlchemyProvider, method string, _ ...string) (any, error) {
					assert.Equal(t, core.Eth_GetTransactionByHash, method)
					return `{"hello": "world"}`, nil
				},
			)
			patches.ApplyFunc(
				mapstructure.Decode,
				func(_ any, _ any) error {
					return nil
				},
			)
			patches.ApplyFunc(
				utils.TransformTransaction,
				func(rawTx types.TransactionRawResponse) (types.TransactionResponse, error) {
					return expectedTransaction, nil
				},
			)

			// Act
			actual, err := ether.GetTransaction("hoge")

			// Assert
			assert.Nil(t, err)
			assert.Equal(t, expectedTransaction, actual)
		})
	})

	t.Run("error case:", func(t *testing.T) {
		t.Run("if invalid send, throw error", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Arrange
			expectedErr := errors.New("error")

			// Mock & Assert
			patches.ApplyMethod(
				reflect.TypeOf(provider),
				"Send",
				func(_ *alchemy.AlchemyProvider, method string, _ ...string) (any, error) {
					assert.Equal(t, core.Eth_GetTransactionByHash, method)
					return "", expectedErr
				},
			)

			// Act
			_, err := ether.GetTransaction("hoge")

			// Assert
			assert.ErrorIs(t, err, expectedErr)
		})

		t.Run("if error on mapstructure, throw core.ErrFailedToMapTransaction", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Mock
			patches.ApplyMethod(
				reflect.TypeOf(provider),
				"Send",
				func(_ *alchemy.AlchemyProvider, method string, _ ...string) (any, error) {
					return `invalid json`, nil
				},
			)
			patches.ApplyFunc(
				mapstructure.Decode,
				func(_ any, _ any) error {
					return errors.New("error")
				},
			)

			// Act
			_, err := ether.GetTransaction("hoge")

			// Assert
			assert.ErrorIs(t, err, core.ErrFailedToMapTransaction)
		})

		t.Run("if error on transform, throw error", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Arrange
			expectedErr := errors.New("error")

			// Mock
			patches.ApplyMethod(
				reflect.TypeOf(provider),
				"Send",
				func(_ *alchemy.AlchemyProvider, method string, _ ...string) (any, error) {
					return `{"hello": "world"}`, nil
				},
			)
			patches.ApplyFunc(
				mapstructure.Decode,
				func(_ any, _ any) error {
					return nil
				},
			)
			patches.ApplyFunc(
				utils.TransformTransaction,
				func(rawTx types.TransactionRawResponse) (types.TransactionResponse, error) {
					return types.TransactionResponse{}, expectedErr
				},
			)

			// Act
			_, err := ether.GetTransaction("hoge")

			// Assert
			assert.ErrorIs(t, err, expectedErr)
		})
	})
}

func TestEther_GetStorageAt(t *testing.T) {
	// Arrange
	provider := newProviderForTest()
	ether := ether.NewEtherApi(provider).(*ether.Ether)

	t.Run("normal case:", func(t *testing.T) {
		t.Run("call eth_getStorageAt & return provided block", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Arrange
			expected := "0xffff"

			// Mock & Assert
			patches.ApplyMethod(
				reflect.TypeOf(provider),
				"Send",
				func(_ *alchemy.AlchemyProvider, method string, _ ...string) (any, error) {
					assert.Equal(t, core.Eth_GetStorageAt, method)
					return expected, nil
				},
			)

			// Act
			actual, _ := ether.GetStorageAt("0x", "0x", "latest")

			// Assert
			assert.Equal(t, expected, actual)
		})
	})

	t.Run("error case:", func(t *testing.T) {
		t.Run("if error occur, return error internal error", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Arrange
			expectedErr := errors.New("error")

			// Mock
			patches.ApplyMethod(
				reflect.TypeOf(provider),
				"Send",
				func(_ *alchemy.AlchemyProvider, method string, _ ...string) (any, error) {
					return "", expectedErr
				},
			)

			// Act
			_, err := ether.GetStorageAt("0x", "0x", "latest")

			// Assert
			assert.ErrorIs(t, expectedErr, err)
		})

		t.Run("if invalid blockTag, return error", func(t *testing.T) {
			// Act
			_, err := ether.GetStorageAt("0x", "0x", "unxpected")

			// Assert
			assert.ErrorIs(t, core.ErrInvalidBlockTag, err)
		})
	})
}

func TestEther_GetTokenBalances(t *testing.T) {
	// Arrange
	provider := newProviderForTest()
	ether := ether.NewEtherApi(provider).(*ether.Ether)
	expectedResponse := map[string]any{
		"address": "0x123",
		"tokenBalances": []map[string]any{
			{
				"contractAddress": "0x456",
				"tokenBalance":    "0x0000000000000000000000000000000000000000000000000000000000000000",
				"error":           nil,
			},
		},
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

	t.Run("normal case:", func(t *testing.T) {
		t.Run("call with alchemy_getTokenBalances and params & return provided block", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Arrange
			params := []string{"0x111"}

			// Mock & Assert
			patches.ApplyMethod(
				reflect.TypeOf(provider),
				"Send",
				func(_ *alchemy.AlchemyProvider, method string, _ ...string) (any, error) {
					assert.Equal(t, core.Alchemy_GetTokenBalances, method)
					return expectedResponse, nil
				},
			)

			// Act
			actual, _ := ether.GetTokenBalances("0x123", params...)

			// Assert
			assert.Equal(t, expected, actual)
		})

		t.Run("call with alchemy_getTokenBalances and no params & return provided block", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Mock & Assert
			patches.ApplyMethod(
				reflect.TypeOf(provider),
				"Send",
				func(_ *alchemy.AlchemyProvider, method string, _ ...string) (any, error) {
					assert.Equal(t, core.Alchemy_GetTokenBalances, method)
					return expectedResponse, nil
				},
			)

			// Act
			actual, _ := ether.GetTokenBalances("0x123")

			// Assert
			assert.Equal(t, expected, actual)
		})

		t.Run("mapping included error", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Arrange
			expectedErrResponse := map[string]any{
				"address": "0x123",
				"tokenBalances": []map[string]any{
					{
						"contractAddress": "0x456",
						"tokenBalance":    "0x0000000000000000000000000000000000000000000000000000000000000000",
						"error":           "error",
					},
				},
			}
			expectedWithErr := types.TokenBalanceResponse{
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
				reflect.TypeOf(provider),
				"Send",
				func(_ *alchemy.AlchemyProvider, method string, _ ...string) (any, error) {
					assert.Equal(t, core.Alchemy_GetTokenBalances, method)
					return expectedErrResponse, nil
				},
			)

			// Act
			actual, _ := ether.GetTokenBalances("0x123")

			// Assert
			assert.Equal(t, expectedWithErr, actual)
		})
	})

	t.Run("error case:", func(t *testing.T) {
		t.Run("if error occur in send, return internal error", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Arrange
			expectedErr := errors.New("error")

			// Mock
			patches.ApplyMethod(
				reflect.TypeOf(provider),
				"Send",
				func(_ *alchemy.AlchemyProvider, method string, _ ...string) (any, error) {
					return "", expectedErr
				},
			)

			// Act
			_, err := ether.GetTokenBalances("0x123")

			// Assert
			assert.ErrorIs(t, expectedErr, err)
		})

		t.Run("mapstructure error, return core.ErrFailedToMapTokenResponse", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Mock
			patches.ApplyMethod(
				reflect.TypeOf(provider),
				"Send",
				func(_ *alchemy.AlchemyProvider, method string, _ ...string) (any, error) {
					assert.Equal(t, core.Alchemy_GetTokenBalances, method)
					return expectedResponse, nil
				},
			)
			patches.ApplyFunc(
				mapstructure.Decode,
				func(_ any, _ any) error {
					return errors.New("error")
				},
			)

			// Act
			_, err := ether.GetTokenBalances("0x123")

			// Assert
			assert.ErrorIs(t, core.ErrFailedToMapTokenResponse, err)
		})
	})
}

func TestEther_GetTokenMetadata(t *testing.T) {
	// Arrange
	provider := newProviderForTest()
	ether := ether.NewEtherApi(provider).(*ether.Ether)
	expectedResponse := map[string]any{
		"name":     "USD Coin",
		"symbol":   "USDC",
		"decimals": 6,
		"logo":     "https://static.alchemyapi.io/images/assets/3408.png",
	}

	t.Run("normal case:", func(t *testing.T) {
		t.Run("call with alchemy_getTokenMetadata & return token balance", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Mock & Assert
			patches.ApplyMethod(
				reflect.TypeOf(provider),
				"Send",
				func(_ *alchemy.AlchemyProvider, method string, _ ...string) (any, error) {
					assert.Equal(t, core.Alchemy_GetTokenMetadata, method)
					return expectedResponse, nil
				},
			)

			// Act
			actual, _ := ether.GetTokenMetadata("0x123")

			// Assert
			assert.Equal(t, types.TokenMetadataResponse{
				Name:     "USD Coin",
				Symbol:   "USDC",
				Decimals: 6,
				Logo:     "https://static.alchemyapi.io/images/assets/3408.png",
			}, actual)
		})
	})

	t.Run("error case:", func(t *testing.T) {
		t.Run("if error occur on send, return internal error", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Arrange
			expectedErr := errors.New("error")

			// Mock
			patches.ApplyMethod(
				reflect.TypeOf(provider),
				"Send",
				func(_ *alchemy.AlchemyProvider, method string, _ ...string) (any, error) {
					return expectedResponse, expectedErr
				},
			)

			// Act
			_, err := ether.GetTokenMetadata("0x123")

			// Assert
			assert.ErrorIs(t, expectedErr, err)
		})

		t.Run("if failed mapstructure, return core.ErrFailedToMapTokenResponse", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Mock & Assert
			patches.ApplyMethod(
				reflect.TypeOf(provider),
				"Send",
				func(_ *alchemy.AlchemyProvider, method string, _ ...string) (any, error) {
					assert.Equal(t, core.Alchemy_GetTokenMetadata, method)
					return expectedResponse, nil
				},
			)
			patches.ApplyFunc(
				mapstructure.Decode,
				func(_ any, _ any) error {
					return errors.New("error")
				},
			)

			// Act
			_, err := ether.GetTokenMetadata("0x123")

			// Assert
			assert.ErrorIs(t, core.ErrFailedToMapTokenResponse, err)
		})
	})
}

func TestEther_GetLogs(t *testing.T) {
	// Arrange
	provider := newProviderForTest()
	ether := ether.NewEtherApi(provider).(*ether.Ether)
	expectedRes := []any{
		map[string]any{
			"logIndex":         "0x0",
			"removed":          false,
			"blockNumber":      "0x233",
			"blockHash":        "0xfc139f5e2edee9e9c888d8df9a2d2226133a9bd87c88ccbd9c930d3d4c9f9ef5",
			"transactionHash":  "0x66e7a140c8fa27fe98fde923defea7562c3ca2d6bb89798aabec65782c08f63d",
			"transactionIndex": "0x0",
			"address":          "0x42699a7612a82f1d9c36148af9c77354759b210b",
			"data":             "0x0000000000000000000000000000000000000000000000000000000000000004",
			"topics": []string{
				"0x04474795f5b996ff80cb47c148d4c5ccdbe09ef27551820caa9c2f8ed149cce3",
			},
		},
		map[string]any{
			"logIndex":         "0x0",
			"removed":          false,
			"blockNumber":      "0x233",
			"blockHash":        "0xfc139f5e2edee9e9c888d8df9a2d2226133a9bd87c88ccbd9c930d3d4c9f9ef5",
			"transactionHash":  "0x66e7a140c8fa27fe98fde923defea7562c3ca2d6bb89798aabec65782c08f63d",
			"transactionIndex": "0x0",
			"address":          "0x42699a7612a82f1d9c36148af9c77354759b210b",
			"data":             "0x0000000000000000000000000000000000000000000000000000000000000004",
			"topics": []string{
				"0x04474795f5b996ff80cb47c148d4c5ccdbe09ef27551820caa9c2f8ed149cce3",
			},
		},
	}

	t.Run("normal case:", func(t *testing.T) {
		t.Run("call eth_getLogs & return logs", func(t *testing.T) {
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

			// Mock
			patches.ApplyMethod(
				reflect.TypeOf(provider),
				"SendFilter",
				func(_ *alchemy.AlchemyProvider, method string, params ...types.Filter) (any, error) {
					assert.Equal(t, core.Eth_GetLogs, method)
					return expectedRes, nil
				},
			)

			// Act
			actual, _ := ether.GetLogs(types.Filter{})

			// Assert
			assert.Equal(t, expected, actual)
		})
	})

	t.Run("error case", func(t *testing.T) {
		t.Run("if error occur in SendFilter, return internal server error", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Arrange
			expectedErr := errors.New("error")

			// Mock
			patches.ApplyMethod(
				reflect.TypeOf(provider),
				"SendFilter",
				func(_ *alchemy.AlchemyProvider, method string, params ...types.Filter) (any, error) {
					assert.Equal(t, core.Eth_GetLogs, method)
					return expectedRes, expectedErr
				},
			)

			// Act
			_, err := ether.GetLogs(types.Filter{})

			// Assert
			assert.ErrorIs(t, expectedErr, err)
		})

		// TODO: why this is not mocked
		/*
			t.Run("if failed mapstructure, return core.ErrFailedToMapTokenResponse", func(t *testing.T) {
				patches := gomonkey.NewPatches()
				defer patches.Reset()

				// Mock
				patches.ApplyMethod(
					reflect.TypeOf(provider),
					"SendFilter",
					func(_ *alchemy.AlchemyProvider, method string, params ...types.Filter) (any, error) {
						assert.Equal(t, core.Eth_GetLogs, method)
						return expectedRes, nil
					},
				)

				patches.ApplyFunc(
					mapstructure.Decode,
					func(_ any, _ any) error {
						return errors.New("error")
					},
				)

				// Act
				_, err := ether.GetTokenMetadata("0x123")

				// Assert
				assert.ErrorIs(t, core.ErrFailedToMapTokenResponse, err)
			})
		*/
	})
}

func TestEther_EstimateGas(t *testing.T) {
	// Arrange
	provider := newProviderForTest()
	ether := ether.NewEtherApi(provider).(*ether.Ether)

	transaction := types.TransactionRequest{
		From:  "0x1234",
		To:    "0x2345",
		Value: "0x1",
	}

	t.Run("normal case", func(t *testing.T) {
		t.Run("call eth_estimateGas & estimate gas", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Arrange
			expectedRes := "0x1234"
			expected, _ := utils.FromBigHex(expectedRes)

			// Mock
			patches.ApplyMethod(
				reflect.TypeOf(provider),
				"SendTransaction",
				func(_ *alchemy.AlchemyProvider, method string, _ ...types.TransactionRequest) (any, error) {
					assert.Equal(t, core.Eth_EstimateGas, method)
					return expectedRes, nil
				},
			)

			// Act
			actual, _ := ether.EstimateGas(transaction)

			// Assert
			assert.Equal(t, expected, actual)
		})
	})

	t.Run("error case", func(t *testing.T) {
		t.Run("if error occur in marshaling parameter, return core.ErrFailedToMarshalParameter", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Mock
			patches.ApplyFunc(
				json.Marshal,
				func(v any) ([]byte, error) {
					return nil, errors.New("error")
				},
			)

			// Act
			_, err := ether.EstimateGas(transaction)

			// Assert
			assert.ErrorIs(t, err, core.ErrFailedToMarshalParameter)
		})

		t.Run("if error occur in Send, return internal error", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Arrange
			expectedErr := errors.New("error")

			// Mock
			patches.ApplyMethod(
				reflect.TypeOf(provider),
				"SendTransaction",
				func(_ *alchemy.AlchemyProvider, method string, _ ...types.TransactionRequest) (any, error) {
					return "", expectedErr
				},
			)

			// Act
			_, err := ether.EstimateGas(transaction)

			// Assert
			assert.ErrorIs(t, err, expectedErr)
		})

		t.Run("if error occur in FromBigHex, return internal error", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Arrange
			expectedErr := errors.New("error")
			expectedRes := "0x1234"

			// Mock
			patches.ApplyMethod(
				reflect.TypeOf(provider),
				"SendTransaction",
				func(_ *alchemy.AlchemyProvider, method string, _ ...types.TransactionRequest) (any, error) {
					assert.Equal(t, core.Eth_EstimateGas, method)
					return expectedRes, nil
				},
			)
			patches.ApplyFunc(
				utils.FromBigHex,
				func(_ string) (*big.Int, error) {
					return big.NewInt(0), expectedErr
				},
			)

			// Act
			_, err := ether.EstimateGas(transaction)

			// Assert
			assert.ErrorIs(t, err, expectedErr)
		})
	})
}

func Test_Call(t *testing.T) {
	provider := newProviderForTest()
	ether := ether.NewEtherApi(provider).(*ether.Ether)

	transaction := types.TransactionRequest{
		To:    "0x2345",
		Value: "0x1",
	}

	t.Run("normal case:", func(t *testing.T) {
		t.Run("call eth_call, & return result", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Arrange
			expectedRes := "0x1234"

			// Mock
			patches.ApplyMethod(
				reflect.TypeOf(provider),
				"SendTransaction",
				func(_ *alchemy.AlchemyProvider, method string, _ ...types.TransactionRequest) (any, error) {
					assert.Equal(t, core.Eth_Call, method)
					return expectedRes, nil
				},
			)

			// Act
			res, err := ether.Call(transaction)

			// Assert
			assert.Nil(t, err)
			assert.Equal(t, expectedRes, res)
		})
	})

	t.Run("error case: ", func(t *testing.T) {
		t.Run("if error occur in Send, return internal error", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Arrange
			expectedErr := errors.New("error")

			// Mock
			patches.ApplyMethod(
				reflect.TypeOf(provider),
				"SendTransaction",
				func(_ *alchemy.AlchemyProvider, method string, _ ...types.TransactionRequest) (any, error) {
					return "", expectedErr
				},
			)

			// Act
			_, err := ether.Call(transaction)

			// Assert
			assert.ErrorIs(t, err, expectedErr)
		})
	})
}
