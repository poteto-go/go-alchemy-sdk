package types_test

import (
	"errors"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/poteto-go/go-alchemy-sdk/types"
	"github.com/stretchr/testify/assert"
)

var sentinel = errors.New("underlying error")

func TestRpcError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  *types.RpcError
		want string
	}{
		{
			name: "formats with method and code",
			err:  &types.RpcError{Method: "eth_call", Code: -32000, Message: "execution reverted", Err: sentinel},
			want: "rpc error on eth_call (code=-32000): execution reverted: underlying error",
		},
		{
			name: "zero code",
			err:  &types.RpcError{Method: "eth_blockNumber", Code: 0, Message: "ok", Err: nil},
			want: "rpc error on eth_blockNumber (code=0): ok",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.err.Error())
		})
	}
}

func TestRpcError_Unwrap(t *testing.T) {
	e := &types.RpcError{Err: sentinel}
	assert.True(t, errors.Is(e, sentinel))
}

func TestRpcError_ErrorsAs(t *testing.T) {
	wrapped := fmt.Errorf("outer: %w", &types.RpcError{Method: "eth_call", Code: -32000, Message: "reverted", Err: sentinel})

	var rpcErr *types.RpcError
	assert.True(t, errors.As(wrapped, &rpcErr))
	assert.Equal(t, "eth_call", rpcErr.Method)
	assert.Equal(t, -32000, rpcErr.Code)
}

func TestTxError_Error(t *testing.T) {
	hash := common.HexToHash("0xdeadbeef")

	tests := []struct {
		name string
		err  *types.TxError
		want string
	}{
		{
			name: "formats with hash and chainID",
			err:  &types.TxError{TxHash: hash, ChainID: big.NewInt(1), Err: sentinel},
			want: "tx error (chain=1, tx=0x00000000000000000000000000000000000000000000000000000000deadbeef): underlying error",
		},
		{
			name: "nil chainID",
			err:  &types.TxError{TxHash: hash, ChainID: nil, Err: sentinel},
			want: "tx error (chain=unknown, tx=0x00000000000000000000000000000000000000000000000000000000deadbeef): underlying error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.err.Error())
		})
	}
}

func TestTxError_Unwrap(t *testing.T) {
	e := &types.TxError{Err: sentinel}
	assert.True(t, errors.Is(e, sentinel))
}

func TestTxError_ErrorsAs(t *testing.T) {
	hash := common.HexToHash("0xdeadbeef")
	wrapped := fmt.Errorf("outer: %w", &types.TxError{TxHash: hash, ChainID: big.NewInt(137), Err: sentinel})

	var txErr *types.TxError
	assert.True(t, errors.As(wrapped, &txErr))
	assert.Equal(t, hash, txErr.TxHash)
	assert.Equal(t, big.NewInt(137), txErr.ChainID)
}
