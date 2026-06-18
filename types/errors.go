package types

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// RpcError wraps a JSON-RPC error with method context so callers can inspect
// the method name and error code programmatically via errors.As.
type RpcError struct {
	Method  string
	Code    int
	Message string
	Err     error
}

func (e *RpcError) Error() string {
	s := fmt.Sprintf("rpc error on %s (code=%d): %s", e.Method, e.Code, e.Message)
	if e.Err != nil {
		s += ": " + e.Err.Error()
	}
	return s
}

func (e *RpcError) Unwrap() error { return e.Err }

// TxError is returned when a transaction-related operation fails.
// Callers can use errors.As to extract the TxHash and ChainID.
type TxError struct {
	TxHash  common.Hash
	ChainID *big.Int
	Err     error
}

func (e *TxError) Error() string {
	chain := "unknown"
	if e.ChainID != nil {
		chain = e.ChainID.String()
	}
	return fmt.Sprintf("tx error (chain=%s, tx=%s): %v", chain, e.TxHash.Hex(), e.Err)
}

func (e *TxError) Unwrap() error { return e.Err }
