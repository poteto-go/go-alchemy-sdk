package batch_test

import (
	"testing"

	"github.com/poteto-go/go-alchemy-sdk/alchemymock"
	"github.com/poteto-go/go-alchemy-sdk/batch"
	"github.com/stretchr/testify/assert"
)

func TestCoreBatch_AllMethods(t *testing.T) {
	mock := alchemymock.NewAlchemyHttpMock(batchSetting, t)
	defer mock.DeactivateAndReset()

	b := batch.NewBatcher(newBatchEther())
	blockNumber := b.Core.BlockNumber()
	gasPrice := b.Core.GasPrice()
	chainID := b.Core.ChainID()
	peerCount := b.Core.PeerCount()
	balance := b.Core.Balance(walletAddr, "latest")
	code := b.Core.Code(contractAddr, "latest")
	storage := b.Core.StorageAt(contractAddr, "0x0", "latest")

	mock.RegisterBatchResponderOnce(resp(
		"10",   // blockNumber -> 16
		"100",  // gasPrice -> 256
		"1",    // chainId -> 1
		"5",    // peerCount -> 5
		"1234", // balance -> 4660
		"abcd", // code -> 0xabcd
		"00ff", // storageAt -> 00ff
	))

	assert.NoError(t, b.Send())

	assertUnwrap(t, blockNumber, uint64(16))
	assertUnwrapStr(t, gasPrice, "256")
	assertUnwrapStr(t, chainID, "1")
	assertUnwrap(t, peerCount, uint64(5))
	assertUnwrapStr(t, balance, "4660")
	assertUnwrap(t, code, "0xabcd")
	assertUnwrap(t, storage, "00ff")
}
