package batch_test

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/poteto-go/go-alchemy-sdk/alchemymock"
	"github.com/poteto-go/go-alchemy-sdk/batch"
	"github.com/poteto-go/go-alchemy-sdk/constant"
	"github.com/poteto-go/go-alchemy-sdk/encode"
	"github.com/poteto-go/go-alchemy-sdk/ether"
	"github.com/poteto-go/go-alchemy-sdk/gas"
	"github.com/poteto-go/go-alchemy-sdk/types"
	"github.com/stretchr/testify/assert"
)

// --- shared test settings and addresses --------------------------------------

var batchSetting = gas.AlchemySetting{
	ApiKey:  "hoge",
	Network: "fuga",
	BackoffConfig: &types.BackoffConfig{
		MaxRetries: 0,
	},
}

const (
	contractAddr = "0x1111111111111111111111111111111111111111"
	walletAddr   = "0x2222222222222222222222222222222222222222"
	ownerAddr    = "0x3333333333333333333333333333333333333333"
	spenderAddr  = "0x4444444444444444444444444444444444444444"
)

// ABI-encoded "MTK" string (offset, length, data).
const abiStringMTK = "0000000000000000000000000000000000000000000000000000000000000020" +
	"0000000000000000000000000000000000000000000000000000000000000003" +
	"4d544b0000000000000000000000000000000000000000000000000000000000"

// --- shared constructor -------------------------------------------------------

// newBatchEther builds an EtherApi that exercises the real HTTP path with a
// non-zero timeout and no retries.
func newBatchEther() types.EtherApi {
	config, err := gas.NewAlchemyConfig(batchSetting)
	if err != nil {
		panic(err)
	}
	provider := gas.NewAlchemyProvider(config)
	return ether.NewEtherApi(provider, ether.NewEtherApiConfig(
		config.GetUrl(),
		0,
		2*time.Second,
		&types.BackoffConfig{MaxRetries: 0},
		[]http.Header{},
		nil,
		0,
		nil,
	))
}

// --- ABI word helpers --------------------------------------------------------

func uintWord(n int64) string {
	b := make([]byte, 32)
	big.NewInt(n).FillBytes(b)
	return hex.EncodeToString(b)
}

func boolWord(v bool) string {
	if v {
		return uintWord(1)
	}
	return uintWord(0)
}

func addrWord(a string) string {
	b := make([]byte, 32)
	copy(b[12:], common.HexToAddress(a).Bytes())
	return hex.EncodeToString(b)
}

func abiStrWord(s string) string {
	return hex.EncodeToString(encode.ABIString(s))
}

// resp builds a JSON-RPC batch response, assigning sequential ids 1..N.
func resp(results ...string) string {
	parts := make([]string, len(results))
	for i, r := range results {
		parts[i] = fmt.Sprintf(`{"jsonrpc":"2.0","id":%d,"result":"0x%s"}`, i+1, r)
	}
	return "[" + strings.Join(parts, ",") + "]"
}

// --- assert helpers ----------------------------------------------------------

func assertUnwrap[T comparable](t *testing.T, r *batch.Result[T], want T) {
	t.Helper()
	got, err := r.Unwrap()
	assert.NoError(t, err)
	assert.Equal(t, want, got)
}

func assertUnwrapStr(t *testing.T, r *batch.Result[*big.Int], want string) {
	t.Helper()
	got, err := r.Unwrap()
	assert.NoError(t, err)
	assert.Equal(t, want, got.String())
}

// --- Batcher infrastructure tests --------------------------------------------

func TestBatcher_Send(t *testing.T) {
	t.Run("typed scalar results are decoded after Send in a single round-trip", func(t *testing.T) {
		mock := alchemymock.NewAlchemyHttpMock(batchSetting, t)
		defer mock.DeactivateAndReset()

		b := batch.NewBatcher(newBatchEther())
		blockNumber := b.Core.BlockNumber()
		gasPrice := b.Core.GasPrice()
		balance := b.Core.Balance("0xABC", "latest")

		mock.RegisterBatchResponderOnce(
			`[{"jsonrpc":"2.0","id":1,"result":"0x10"},` +
				`{"jsonrpc":"2.0","id":2,"result":"0x100"},` +
				`{"jsonrpc":"2.0","id":3,"result":"0x1234"}]`,
		)

		err := b.Send()
		assert.NoError(t, err)

		n, err := blockNumber.Unwrap()
		assert.NoError(t, err)
		assert.Equal(t, uint64(16), n)

		gp, err := gasPrice.Unwrap()
		assert.NoError(t, err)
		assert.Equal(t, "256", gp.String())

		bal, err := balance.Unwrap()
		assert.NoError(t, err)
		assert.Equal(t, "4660", bal.String())
	})

	t.Run("AddCall decodes a contract eth_call result", func(t *testing.T) {
		mock := alchemymock.NewAlchemyHttpMock(batchSetting, t)
		defer mock.DeactivateAndReset()

		b := batch.NewBatcher(newBatchEther())
		balance := batch.AddCall(
			b,
			"0xcontract",
			[]byte("balanceOf(address)"),
			func(out []byte) (*big.Int, error) { return new(big.Int).SetBytes(out), nil },
			common.LeftPadBytes(common.HexToAddress("0xabc").Bytes(), constant.ABIWordSize),
		)

		// eth_call returns a 32-byte word -> 0x0a == 10.
		mock.RegisterBatchResponderOnce(
			`[{"jsonrpc":"2.0","id":1,"result":"0x000000000000000000000000000000000000000000000000000000000000000a"}]`,
		)

		err := b.Send()
		assert.NoError(t, err)

		bal, err := balance.Unwrap()
		assert.NoError(t, err)
		assert.Equal(t, "10", bal.String())
	})

	t.Run("per-request RPC error surfaces only on that result", func(t *testing.T) {
		mock := alchemymock.NewAlchemyHttpMock(batchSetting, t)
		defer mock.DeactivateAndReset()

		b := batch.NewBatcher(newBatchEther())
		blockNumber := b.Core.BlockNumber()
		gasPrice := b.Core.GasPrice()

		mock.RegisterBatchResponderOnce(
			`[{"jsonrpc":"2.0","id":1,"result":"0x10"},` +
				`{"jsonrpc":"2.0","id":2,"error":{"code":-32000,"message":"boom"}}]`,
		)

		err := b.Send()
		assert.NoError(t, err)

		n, err := blockNumber.Unwrap()
		assert.NoError(t, err)
		assert.Equal(t, uint64(16), n)

		_, err = gasPrice.Unwrap()
		assert.Error(t, err)
	})

	t.Run("Unwrap before Send returns ErrBatchNotSent", func(t *testing.T) {
		b := batch.NewBatcher(newBatchEther())
		blockNumber := b.Core.BlockNumber()

		_, err := blockNumber.Unwrap()

		assert.ErrorIs(t, err, constant.ErrBatchNotSent)
	})

	t.Run("Send twice returns ErrBatchAlreadySent", func(t *testing.T) {
		mock := alchemymock.NewAlchemyHttpMock(batchSetting, t)
		defer mock.DeactivateAndReset()

		b := batch.NewBatcher(newBatchEther())
		b.Core.BlockNumber()
		mock.RegisterBatchResponderOnce(`[{"jsonrpc":"2.0","id":1,"result":"0x10"}]`)

		assert.NoError(t, b.Send())
		assert.ErrorIs(t, b.Send(), constant.ErrBatchAlreadySent)
	})

	t.Run("empty batch Send is a no-op", func(t *testing.T) {
		b := batch.NewBatcher(newBatchEther())

		assert.NoError(t, b.Send())
	})

	t.Run("I/O failure is returned by Send and results stay unsettled", func(t *testing.T) {
		mock := alchemymock.NewAlchemyHttpMock(batchSetting, t)
		defer mock.DeactivateAndReset()

		b := batch.NewBatcher(newBatchEther())
		blockNumber := b.Core.BlockNumber()

		err := b.Send()

		assert.Error(t, err)
		_, unwrapErr := blockNumber.Unwrap()
		assert.ErrorIs(t, unwrapErr, constant.ErrBatchNotSent)
	})
}
