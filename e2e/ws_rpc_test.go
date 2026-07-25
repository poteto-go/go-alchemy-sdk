package e2e

import (
	"context"
	"math/big"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/poteto-go/go-alchemy-sdk/_fixture/artifacts"
	"github.com/poteto-go/go-alchemy-sdk/gas"
	"github.com/poteto-go/go-alchemy-sdk/types"
	"github.com/poteto-go/go-alchemy-sdk/wallet"
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

	a, err := gas.NewWsAlchemy(gas.AlchemySetting{
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

// subscribeTimeout bounds how long a subscription assertion waits for its
// event, so a subscription that never fires fails the test instead of hanging
// until the go test panic timeout.
const subscribeTimeout = 30 * time.Second

func TestSenario_Ws_Subscribe(t *testing.T) {
	wsAlchemy := newWsAlchemy(t)
	defer wsAlchemy.GetProvider().Eth().Shutdown()

	w, _ := wallet.New(initPrivateKey)
	w.Connect(wsAlchemy.GetProvider())

	t.Run("can subscribe new head", func(t *testing.T) {
		headers := make(chan *gethTypes.Header)

		sub, err := wsAlchemy.WS.SubscribeNewHead(
			context.Background(),
			headers,
		)
		require.NoError(t, err)
		defer sub.Unsubscribe()

		txRequest := types.TransactionRequest{
			From:     initAddress,
			To:       otherAddress,
			Value:    "0x123",
			GasLimit: 300000,
		}

		txHash, err := w.SendTransaction(txRequest)
		require.NoError(t, err)

		txReceipt, err := wsAlchemy.Core.GetTransactionReceipt(txHash.Hex())
		require.NoError(t, err)

		for {
			select {
			case err := <-sub.Err():
				require.FailNow(t, "unexpected err: "+err.Error())
			case <-time.After(subscribeTimeout):
				require.FailNow(t, "timed out waiting for a new head")
			case header := <-headers:
				assert.Equal(t, txReceipt.BlockNumber.Uint64(), header.Number.Uint64())
				return
			}
		}
	})

	t.Run("can subscribe contract event by subscribe log", func(t *testing.T) {
		contract := artifacts.NewPotetoStorage()
		contractAddress, err := w.DeployContract(context.Background(), &artifacts.PotetoStorageMetaData)
		require.NoError(t, err)

		logs := make(chan gethTypes.Log)
		sub, err := wsAlchemy.WS.SubscribeLogs(
			context.Background(),
			ethereum.FilterQuery{
				Addresses: []common.Address{
					contractAddress,
				},
			},
			logs,
		)
		require.NoError(t, err)
		defer sub.Unsubscribe()

		// store() emits Stored(address indexed sender, uint256 value),
		// so this tx is what pushes a log down the subscription.
		data := contract.PackStore(big.NewInt(42))
		receipt, err := w.ContractTransact(context.Background(), contractAddress.Hex(), data)
		require.NoError(t, err)

		for {
			select {
			case err := <-sub.Err():
				require.FailNow(t, "unexpected err: "+err.Error())
			case <-time.After(subscribeTimeout):
				require.FailNow(t, "timed out waiting for the contract event")
			case log := <-logs:
				assert.Equal(t, receipt.BlockNumber.Uint64(), log.BlockNumber)

				event, err := contract.UnpackStoredEvent(&log)
				require.NoError(t, err)
				assert.Equal(t, common.HexToAddress(initAddress), event.Sender)
				assert.Equal(t, uint64(42), event.Value.Uint64())

				return
			}
		}
	})

	t.Run("can subscribe contract event by subscribe contract log", func(t *testing.T) {
		contract := artifacts.NewPotetoStorage()
		contractAddress, err := w.DeployContract(context.Background(), &artifacts.PotetoStorageMetaData)
		require.NoError(t, err)

		logs := make(chan gethTypes.Log)
		sub, err := wsAlchemy.WS.SubscribeContractLogs(
			context.Background(),
			contractAddress,
			logs,
		)
		require.NoError(t, err)
		defer sub.Unsubscribe()

		// store() emits Stored(address indexed sender, uint256 value),
		// so this tx is what pushes a log down the subscription.
		data := contract.PackStore(big.NewInt(42))
		receipt, err := w.ContractTransact(context.Background(), contractAddress.Hex(), data)
		require.NoError(t, err)

		for {
			select {
			case err := <-sub.Err():
				require.FailNow(t, "unexpected err: "+err.Error())
			case <-time.After(subscribeTimeout):
				require.FailNow(t, "timed out waiting for the contract event")
			case log := <-logs:
				assert.Equal(t, receipt.BlockNumber.Uint64(), log.BlockNumber)

				event, err := contract.UnpackStoredEvent(&log)
				require.NoError(t, err)
				assert.Equal(t, common.HexToAddress(initAddress), event.Sender)
				assert.Equal(t, uint64(42), event.Value.Uint64())

				return
			}
		}
	})
}
