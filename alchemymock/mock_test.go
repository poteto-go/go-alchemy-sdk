package alchemymock_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/poteto-go/go-alchemy-sdk/alchemymock"
	"github.com/poteto-go/go-alchemy-sdk/ether"
	"github.com/poteto-go/go-alchemy-sdk/gas"
	"github.com/poteto-go/go-alchemy-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestNewAlchemyHttpMockAndDefer(t *testing.T) {
	t.Run("can crate instance", func(t *testing.T) {
		// Act
		alchemyMock := alchemymock.NewAlchemyHttpMock(gas.AlchemySetting{
			ApiKey:  "hoge",
			Network: "fuga",
			BackoffConfig: &types.BackoffConfig{
				MaxRetries: 0,
			},
		}, t)
		defer alchemyMock.DeactivateAndReset()

		// Assert
		assert.NotNil(t, alchemyMock)
	})
}

func TestAlchemyMock_RegisterResponder(t *testing.T) {
	t.Run("can register responder w/ expected ethMethod", func(t *testing.T) {
		// Arrange
		setting := gas.AlchemySetting{
			ApiKey:  "hoge",
			Network: "fuga",
			BackoffConfig: &types.BackoffConfig{
				MaxRetries: 0,
			},
		}
		alchemyMock := alchemymock.NewAlchemyHttpMock(setting, t)
		defer alchemyMock.DeactivateAndReset()

		config := gas.NewAlchemyConfig(setting)
		provider := gas.NewAlchemyProvider(config)
		eth := ether.NewEtherApi(
			provider,
			ether.NewEtherApiConfig(
				"https://fuga.g.alchemy.com/v2/hoge",
				0,
				time.Duration(0),
				nil,
			),
		).(*ether.Ether)

		// Act
		alchemyMock.RegisterResponder("eth_getBalance", `{"jsonrpc":"2.0","id":1,"result":"0x1234"}`)
		balance, err := eth.GetBalance("hoge", "latest")

		// Assert
		assert.NotNil(t, alchemyMock)
		assert.NoError(t, err)
		assert.Equal(t, "4660", balance.String())
	})

	t.Run("if ethMethod is not match, return error", func(t *testing.T) {
		// Arrange
		setting := gas.AlchemySetting{
			ApiKey:  "hoge",
			Network: "fuga",
			BackoffConfig: &types.BackoffConfig{
				MaxRetries: 0,
			},
		}
		alchemyMock := alchemymock.NewAlchemyHttpMock(setting, t)
		defer alchemyMock.DeactivateAndReset()

		config := gas.NewAlchemyConfig(setting)
		provider := gas.NewAlchemyProvider(config)
		eth := ether.NewEtherApi(
			provider,
			ether.NewEtherApiConfig(
				"https://fuga.g.alchemy.com/v2/hoge",
				0,
				time.Duration(0),
				nil,
			),
		).(*ether.Ether)

		// Act
		alchemyMock.RegisterResponder("eth_unexpected", `{"jsonrpc":"2.0","id":1,"result":"0x1234"}`)
		_, err := eth.GetBalance("hoge", "latest")

		// Assert
		assert.Error(t, err)
	})

	t.Run("if deactivated, not mock work", func(t *testing.T) {
		// Arrange
		setting := gas.AlchemySetting{
			ApiKey:  "hoge",
			Network: "fuga",
			BackoffConfig: &types.BackoffConfig{
				MaxRetries: 0,
			},
		}
		alchemyMock := alchemymock.NewAlchemyHttpMock(setting, t)
		alchemyMock.DeactivateAndReset()

		config := gas.NewAlchemyConfig(setting)
		provider := gas.NewAlchemyProvider(config)
		eth := ether.NewEtherApi(
			provider,
			ether.NewEtherApiConfig(
				"https://fuga.g.alchemy.com/v2/hoge",
				0,
				time.Duration(0),
				nil,
			),
		).(*ether.Ether)

		// Act
		alchemyMock.RegisterResponder("eth_getBalance", `{"jsonrpc":"2.0","id":1,"result":"0x1234"}`)
		_, err := eth.GetBalance("hoge", "latest")

		// Assert
		assert.Error(t, err)
	})

	t.Run("if not json rpc, return error", func(t *testing.T) {
		// Arrange
		setting := gas.AlchemySetting{
			ApiKey:  "hoge",
			Network: "fuga",
			BackoffConfig: &types.BackoffConfig{
				MaxRetries: 0,
			},
		}
		alchemyMock := alchemymock.NewAlchemyHttpMock(setting, t)
		defer alchemyMock.DeactivateAndReset()

		// Act
		alchemyMock.RegisterResponder("eth_getBalance", `{"jsonrpc":"2.0","id":1,"result":"0x1234"}`)
		req, _ := http.NewRequest(
			http.MethodPost,
			"https://fuga.g.alchemy.com/v2/hoge",
			nil,
		)
		_, err := http.DefaultClient.Do(req)

		// Assert
		assert.Error(t, err)
	})
}
