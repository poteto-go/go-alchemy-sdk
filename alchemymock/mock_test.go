package alchemymock_test

import (
	"net/http"
	"testing"

	"github.com/poteto-go/go-alchemy-sdk/alchemymock"
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
		alchemy := gas.NewAlchemy(setting)

		// Act
		alchemyMock.RegisterResponder("eth_getBalance", `{"jsonrpc":"2.0","id":1,"result":"0x1234"}`)
		balance, err := alchemy.Core.GetBalance("0x", "latest")

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
		alchemy := gas.NewAlchemy(setting)

		// Act
		alchemyMock.RegisterResponder("eth_unexpected", `{"jsonrpc":"2.0","id":1,"result":"0x1234"}`)
		_, err := alchemy.Core.GetBalance("0x", "latest")

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
		alchemy := gas.NewAlchemy(setting)

		// Act
		alchemyMock.RegisterResponder("eth_getBalance", `{"jsonrpc":"2.0","id":1,"result":"0x1234"}`)
		_, err := alchemy.Core.GetBalance("0x", "latest")

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
