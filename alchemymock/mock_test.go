package alchemymock_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/poteto-go/go-alchemy-sdk/alchemymock"
	"github.com/poteto-go/go-alchemy-sdk/batch"
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
		alchemy, err := gas.NewAlchemy(setting)
		assert.NoError(t, err)

		// Act
		alchemyMock.RegisterResponderOnce("eth_getBalance", `{"jsonrpc":"2.0","id":1,"result":"0x1234"}`)
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
		alchemy, err := gas.NewAlchemy(setting)
		assert.NoError(t, err)

		// Act
		alchemyMock.RegisterResponderOnce("eth_unexpected", `{"jsonrpc":"2.0","id":1,"result":"0x1234"}`)
		_, err = alchemy.Core.GetBalance("0x", "latest")

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
		alchemy, err := gas.NewAlchemy(setting)
		assert.NoError(t, err)

		// Act
		alchemyMock.RegisterResponderOnce("eth_getBalance", `{"jsonrpc":"2.0","id":1,"result":"0x1234"}`)
		_, err = alchemy.Core.GetBalance("0x", "latest")

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
		alchemyMock.RegisterResponderOnce("eth_getBalance", `{"jsonrpc":"2.0","id":1,"result":"0x1234"}`)
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

func TestAlchemyMock_MultipleResponders(t *testing.T) {
	t.Run("can register multiple responders", func(t *testing.T) {
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
		alchemy, err := gas.NewAlchemy(setting)
		assert.NoError(t, err)

		// Act
		// Register first responder
		alchemyMock.RegisterResponderOnce("eth_getBalance", `{"jsonrpc":"2.0","id":1,"result":"0x1234"}`)
		// Register second responder
		alchemyMock.RegisterResponderOnce("eth_blockNumber", `{"jsonrpc":"2.0","id":1,"result":"0x10"}`)

		// Assert
		// Call first method
		balance, err1 := alchemy.Core.GetBalance("0x", "latest")
		assert.NoError(t, err1)
		assert.Equal(t, "4660", balance.String())

		// Call second method
		blockNumber, err2 := alchemy.Core.GetBlockNumber()
		assert.NoError(t, err2)
		assert.Equal(t, uint64(16), blockNumber)
	})
}

func TestAlchemyMock_SequenceResponders(t *testing.T) {
	t.Run("responders return values in sequence", func(t *testing.T) {
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
		alchemy, err := gas.NewAlchemy(setting)
		assert.NoError(t, err)

		// Act
		// Register sequence of responders for the same method
		alchemyMock.RegisterResponderOnce("eth_getBalance", `{"jsonrpc":"2.0","id":1,"result":"0x1"}`) // 1
		alchemyMock.RegisterResponderOnce("eth_getBalance", `{"jsonrpc":"2.0","id":1,"result":"0x2"}`) // 2

		// Assert
		// First call should get first result
		balance1, err3 := alchemy.Core.GetBalance("0x", "latest")
		assert.NoError(t, err3)
		assert.Equal(t, "1", balance1.String())

		// Second call should get second result
		balance2, err4 := alchemy.Core.GetBalance("0x", "latest")
		assert.NoError(t, err4)
		assert.Equal(t, "2", balance2.String())
	})
}

var batchMockSetting = gas.AlchemySetting{
	ApiKey:  "hoge",
	Network: "fuga",
	BackoffConfig: &types.BackoffConfig{
		MaxRetries: 0,
	},
}

// newBatcherForMock builds a real Batcher (via ether) so the batch responder is
// exercised through actual batch traffic rather than a hand-crafted request.
func newBatcherForMock() *batch.Batcher {
	config, err := gas.NewAlchemyConfig(batchMockSetting)
	if err != nil {
		panic(err)
	}
	provider := gas.NewAlchemyProvider(config)
	e := ether.NewEtherApi(provider, ether.NewEtherApiConfig(
		config.GetUrl(),
		0,
		2*time.Second,
		&types.BackoffConfig{MaxRetries: 0},
		[]http.Header{},
		nil,
		0,
		nil,
	))
	return batch.NewBatcher(e)
}

func TestAlchemyMock_BatchResponder(t *testing.T) {
	t.Run("a batch sent via the Batcher is served by the batch responder", func(t *testing.T) {
		// Arrange
		alchemyMock := alchemymock.NewAlchemyHttpMock(batchMockSetting, t)
		defer alchemyMock.DeactivateAndReset()

		b := newBatcherForMock()
		balance := b.Core.Balance("0xabc", "latest")
		blockNumber := b.Core.BlockNumber()

		// ids are assigned sequentially (1, 2, ...) by a fresh geth rpc.Client.
		alchemyMock.RegisterBatchResponderOnce(
			`[{"jsonrpc":"2.0","id":1,"result":"0x1234"},{"jsonrpc":"2.0","id":2,"result":"0x10"}]`,
		)

		// Act
		err := b.Send()

		// Assert
		assert.NoError(t, err)

		bal, err := balance.Unwrap()
		assert.NoError(t, err)
		assert.Equal(t, "4660", bal.String())

		bn, err := blockNumber.Unwrap()
		assert.NoError(t, err)
		assert.Equal(t, uint64(16), bn)
	})

	t.Run("if no batch responder registered, Send returns error", func(t *testing.T) {
		// Arrange
		alchemyMock := alchemymock.NewAlchemyHttpMock(batchMockSetting, t)
		defer alchemyMock.DeactivateAndReset()

		b := newBatcherForMock()
		b.Core.BlockNumber()

		// Act: no batch responder registered.
		err := b.Send()

		// Assert
		assert.Error(t, err)
	})
}
