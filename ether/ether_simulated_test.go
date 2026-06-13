package ether_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient/simulated"
	"github.com/poteto-go/go-alchemy-sdk/_fixture/artifacts"
	eth "github.com/poteto-go/go-alchemy-sdk/ether"
	"github.com/stretchr/testify/assert"
)

// well-known anvil key #0; funds the simulated genesis in these tests.
const simPrivateKeyHex = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"

// simulated.NewBackend always uses chainID 1337.
var simChainID = big.NewInt(1337)

func newSimulatedEtherForTest(t *testing.T) (*eth.Ether, func()) {
	t.Helper()

	key, err := crypto.HexToECDSA(simPrivateKeyHex)
	assert.Nil(t, err)
	from := crypto.PubkeyToAddress(key.PublicKey)

	balance := new(big.Int).Mul(big.NewInt(1_000_000), big.NewInt(1_000_000_000_000_000_000))
	backend := simulated.NewBackend(gethTypes.GenesisAlloc{
		from: {Balance: balance},
	})

	e := eth.NewSimulatedApi(backend).(*eth.Ether)
	return e, func() { _ = backend.Close() }
}

// simSign signs a legacy tx from the funded account and submits it to the
// simulated tx pool, returning its hash. to == nil produces a contract creation.
func simSign(t *testing.T, e *eth.Ether, to *common.Address, data []byte) common.Hash {
	t.Helper()

	key, err := crypto.HexToECDSA(simPrivateKeyHex)
	assert.Nil(t, err)
	from := crypto.PubkeyToAddress(key.PublicKey)

	gasPrice, err := e.SuggestGasPrice()
	assert.Nil(t, err)
	nonce, err := e.PendingNonceAt(from.Hex())
	assert.Nil(t, err)

	gasLimit := uint64(21000)
	if to == nil {
		gasLimit = 1_500_000
	}
	tx := gethTypes.NewTx(&gethTypes.LegacyTx{
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      gasLimit,
		To:       to,
		Value:    big.NewInt(0),
		Data:     data,
	})
	signed, err := gethTypes.SignTx(tx, gethTypes.LatestSignerForChainID(simChainID), key)
	assert.Nil(t, err)
	assert.Nil(t, e.SendRawTransaction(signed))
	return signed.Hash()
}

func simSendValueTx(t *testing.T, e *eth.Ether) common.Hash {
	t.Helper()
	key, _ := crypto.HexToECDSA(simPrivateKeyHex)
	to := crypto.PubkeyToAddress(key.PublicKey)
	return simSign(t, e, &to, nil)
}

func simSendDeployTx(t *testing.T, e *eth.Ether) common.Hash {
	t.Helper()
	return simSign(t, e, nil, common.FromHex(artifacts.PotetoStorageMetaData.Bin))
}

func TestEther_NewSimulatedApi(t *testing.T) {
	t.Run("Client returns a usable simulated client", func(t *testing.T) {
		e, cleanup := newSimulatedEtherForTest(t)
		defer cleanup()

		assert.NotNil(t, e.Client())

		bn, err := e.BlockNumber()
		assert.NoError(t, err)
		assert.Equal(t, uint64(0), bn) // fresh genesis
	})

	t.Run("SetEthClient and Close are no-ops on simulated backend", func(t *testing.T) {
		e, cleanup := newSimulatedEtherForTest(t)
		defer cleanup()

		assert.NoError(t, e.SetEthClient())
		e.Close()

		// the simulated client is not torn down, so reads still work after Close
		_, err := e.BlockNumber()
		assert.NoError(t, err)
	})
}
