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
