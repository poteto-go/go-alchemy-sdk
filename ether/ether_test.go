package ether_test

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey"
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

func TestGetBlockNumber(t *testing.T) {
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

func TestGetGasPrice(t *testing.T) {
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

func TestGetBalance(t *testing.T) {
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
	ether := ether.NewEtherApi(provider).(*ether.Ether)

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
				func(_ *alchemy.AlchemyProvider, method string, _ ...string) (string, error) {
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
				func(_ *alchemy.AlchemyProvider, method string, _ ...string) (string, error) {
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
			expectedTransaction := &types.TransactionResponse{
				BlockNumber:          500,
				BlockHash:            "0xb1112ef37861f39ff395a245eb962791e11eae26f94b50bb95e3e31378ef3d25",
				Index:                180,
				Hash:                 "0xb1fac2cb5074a4eda8296faebe3b5a3c10b48947dd9a738b2fdf859be0e1fbaf",
				Type:                 2,
				To:                   "0x111111517e4929d3dcbdfa7cce55d30d4b6bc4d6",
				From:                 "0x2d218ce7d8892fc6b391b614f84278d12decae52",
				Nonce:                1000,
				GasLimit:             big.NewInt(62584),
				GasPrice:             big.NewInt(24999999999),
				MaxPriorityFeePerGas: big.NewInt(380000000),
				MaxFeePerGas:         big.NewInt(26950000000),
				Data:                 "0x",
				Value:                big.NewInt(11001000),
				ChainID:              1,
				Signature: types.Signature{
					R: "0x1b5e176d927f8e9ab405058b2d2457392da3",
					S: "0x4ba69724e8f69de52f0125ad8b3c5c2cef33",
					V: big.NewInt(27),
				},
				AccessList:          []string{},
				BlobVersionedHashes: []string{},
				AuthorizationList:   []string{},
			}
			hash := "0xb1fac2cb5074a4eda8296faebe3b5a3c10b48947dd9a738b2fdf859be0e1fbaf"
			jsonResponse := `{
				"blockHash": "0xb1112ef37861f39ff395a245eb962791e11eae26f94b50bb95e3e31378ef3d25",
				"blockNumber": "0xfd27df",
				"hash": "0xb1fac2cb5074a4eda8296faebe3b5a3c10b48947dd9a738b2fdf859be0e1fbaf",
				"accessList": [],
				"transactionIndex": "0xb4",
				"type": "0x2",
				"nonce": "0x2768",
				"input": "0xa9059cbb000000000000000000000000b6ae07829376a5b704bb46a0869f383555097c29000000000000000000000000000000000000000000000034df6db862352c72d0",
				"r": "0x4dec2c2ab964f28385d31cd203fe5960e001ccd110db816ad462d411cf496548",
				"s": "0x62ffcab5b6ae1cf4a59d32dd39a92f14eadea5fbbb7587c1a845a3d0d8621253",
				"chainId": "0x1",
				"v": "0x1",
				"gas": "0xf478",
				"maxPriorityFeePerGas": "0x173eed80",
				"from": "0x2d218ce7d8892fc6b391b614f84278d12decae52",
				"to": "0x111111517e4929d3dcbdfa7cce55d30d4b6bc4d6",
				"maxFeePerGas": "0x645a4b0a6",
				"value": "0x0",
				"gasPrice": "0x5bcdcacee"
			}`

			// Mock & Assert
			patches.ApplyMethod(
				reflect.TypeOf(provider),
				"Send",
				func(_ *alchemy.AlchemyProvider, method string, _ ...string) (string, error) {
					assert.Equal(t, core.Eth_GetTransactionByHash, method)
					return jsonResponse, nil
				},
			)

			// Act
			actual, err := ether.GetTransaction(hash)

			// Assert
			assert.Nil(t, err)
			assert.Equal(t, expectedTransaction, actual)
		})
	})

	t.Run("error case:", func(t *testing.T) {
		t.Run("if invalid send, throw error", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Mock & Assert
			patches.ApplyMethod(
				reflect.TypeOf(provider),
				"Send",
				func(_ *alchemy.AlchemyProvider, method string, _ ...string) (string, error) {
					assert.Equal(t, core.Eth_GetTransactionByHash, method)
					return "", core.ErrFailedToConnect
				},
			)

			// Act
			actual, err := ether.GetTransaction("hoge")

			// Assert
			assert.ErrorIs(t, err, core.ErrFailedToConnect)
			assert.Nil(t, actual)
		})
	})
}
