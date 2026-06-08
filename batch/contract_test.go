package batch_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/poteto-go/go-alchemy-sdk/alchemymock"
	"github.com/poteto-go/go-alchemy-sdk/batch"
	"github.com/stretchr/testify/assert"
)

const (
	contractAddr = "0x1111111111111111111111111111111111111111"
	walletAddr   = "0x2222222222222222222222222222222222222222"
)

// ABI-encoded "MTK" string (offset, length, data).
const abiStringMTK = "0000000000000000000000000000000000000000000000000000000000000020" +
	"0000000000000000000000000000000000000000000000000000000000000003" +
	"4d544b0000000000000000000000000000000000000000000000000000000000"

func TestBatcher_ERC20(t *testing.T) {
	mock := alchemymock.NewAlchemyHttpMock(batchSetting, t)
	defer mock.DeactivateAndReset()

	b := batch.NewBatcher(newBatchEther())
	balance := b.ERC20.BalanceOf(contractAddr, walletAddr)
	name := b.ERC20.Name(contractAddr)

	mock.RegisterBatchResponderOnce(
		`[{"jsonrpc":"2.0","id":1,"result":"0x000000000000000000000000000000000000000000000000000000000000000a"},` +
			`{"jsonrpc":"2.0","id":2,"result":"0x` + abiStringMTK + `"}]`,
	)

	assert.NoError(t, b.Send())

	bal, err := balance.Unwrap()
	assert.NoError(t, err)
	assert.Equal(t, "10", bal.String())

	nm, err := name.Unwrap()
	assert.NoError(t, err)
	assert.Equal(t, "MTK", nm)
}

func TestBatcher_StableCoin(t *testing.T) {
	mock := alchemymock.NewAlchemyHttpMock(batchSetting, t)
	defer mock.DeactivateAndReset()

	b := batch.NewBatcher(newBatchEther())
	owner := b.StableCoin.Owner(contractAddr)
	paused := b.StableCoin.Paused(contractAddr)

	mock.RegisterBatchResponderOnce(
		`[{"jsonrpc":"2.0","id":1,"result":"0x000000000000000000000000aabbccddeeff00112233445566778899aabbccdd"},` +
			`{"jsonrpc":"2.0","id":2,"result":"0x0000000000000000000000000000000000000000000000000000000000000001"}]`,
	)

	assert.NoError(t, b.Send())

	o, err := owner.Unwrap()
	assert.NoError(t, err)
	assert.Equal(t, common.HexToAddress("0xaabbccddeeff00112233445566778899aabbccdd"), o)

	p, err := paused.Unwrap()
	assert.NoError(t, err)
	assert.True(t, p)
}

func TestBatcher_StableCoinInheritsERC20(t *testing.T) {
	mock := alchemymock.NewAlchemyHttpMock(batchSetting, t)
	defer mock.DeactivateAndReset()

	b := batch.NewBatcher(newBatchEther())
	// Name is an ERC-20 read promoted onto StableCoin via embedding.
	name := b.StableCoin.Name(contractAddr)

	mock.RegisterBatchResponderOnce(
		`[{"jsonrpc":"2.0","id":1,"result":"0x` + abiStringMTK + `"}]`,
	)

	assert.NoError(t, b.Send())

	nm, err := name.Unwrap()
	assert.NoError(t, err)
	assert.Equal(t, "MTK", nm)
}

func TestBatcher_InvalidAddressFailsResultOnly(t *testing.T) {
	b := batch.NewBatcher(newBatchEther())

	// invalid address: result is settled with the validation error, not queued.
	bal := b.ERC20.BalanceOf("not-an-address", "0xwallet")

	_, err := bal.Unwrap()
	assert.Error(t, err)
}
