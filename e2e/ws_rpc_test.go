package e2e

import (
	"math/big"
	"os"
	"strconv"
	"testing"

	"github.com/poteto-go/go-alchemy-sdk/gas"
	"github.com/poteto-go/go-alchemy-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newWsAlchemy builds an Alchemy whose geth client dials anvil over WebSocket.
// anvil serves JSON-RPC over both http and ws on the same port, so
// ws://127.0.0.1:RPC_PORT reaches the same node as the http e2e.
func newWsAlchemy(t *testing.T) gas.Alchemy {
	t.Helper()

	port, err := strconv.Atoi(os.Getenv("RPC_PORT"))
	require.NoError(t, err)

	a, err := gas.NewAlchemy(gas.AlchemySetting{
		// A private-network url with the ws scheme makes Ether.isWebSocket() true,
		// so the geth client dials a persistent WebSocket instead of per-call HTTP.
		PrivateNetworkConfig: gas.PrivateNetworkConfig{
			Url: "ws://127.0.0.1:" + strconv.Itoa(port),
		},
	})
	require.NoError(t, err)
	return a
}

// TestScenario_Ws_BaseMethod verifies that both geth-client-backed RPC calls and
// provider.Send-based methods round-trip over a single persistent WebSocket
// connection to anvil — a ws Alchemy serves the whole surface over one socket.
func TestScenario_Ws_BaseMethod(t *testing.T) {
	wsAlchemy := newWsAlchemy(t)
	// Close() is a no-op on ws, so the socket is persistent; shut it down explicitly.
	defer wsAlchemy.GetProvider().Eth().Shutdown()

	t.Run("GetBlockNumber over ws", func(t *testing.T) {
		bn, err := wsAlchemy.Core.GetBlockNumber()

		assert.Nil(t, err)
		assert.GreaterOrEqual(t, bn, uint64(0))
	})

	t.Run("GetGasPrice over ws", func(t *testing.T) {
		gasPrice, err := wsAlchemy.Core.GetGasPrice()

		assert.Nil(t, err)
		assert.Equal(t, 1, gasPrice.Cmp(big.NewInt(0)))
	})

	t.Run("SuggestGasTipCap over ws", func(t *testing.T) {
		tip, err := wsAlchemy.Core.SuggestGasTipCap()

		assert.Nil(t, err)
		assert.NotNil(t, tip)
	})

	t.Run("SuggestEIP1559Fees over ws", func(t *testing.T) {
		tip, maxFee, err := wsAlchemy.Core.SuggestEIP1559Fees()

		assert.Nil(t, err)
		assert.NotNil(t, tip)
		assert.NotNil(t, maxFee)
		assert.True(t, maxFee.Cmp(tip) >= 0)
	})

	t.Run("EstimateGas over ws", func(t *testing.T) {
		gasLimit, err := wsAlchemy.Core.EstimateGas(types.TransactionRequest{
			From:  initAddress,
			To:    "0x0",
			Value: "0x0",
		})

		assert.Nil(t, err)
		assert.Equal(t, 1, gasLimit.Cmp(big.NewInt(0)))
	})

	t.Run("GetBalance over ws (provider.Send routed over the socket)", func(t *testing.T) {
		balance, err := wsAlchemy.Core.GetBalance(initAddress, "latest")

		assert.Nil(t, err)
		assert.Equal(t, 1, balance.Cmp(big.NewInt(0)))
	})

	t.Run("persistent ws client survives a per-call Close", func(t *testing.T) {
		eth := wsAlchemy.GetProvider().Eth()

		require.NoError(t, eth.SetEthClient())
		eth.Close() // no-op on ws
		require.NotNil(t, eth.Client())
	})
}
