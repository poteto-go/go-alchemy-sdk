package constant

import (
	"math/big"

	"github.com/ethereum/go-ethereum/rpc"
)

// BlockTagNumbers maps EIP-1898 named block tags to their geth rpc.BlockNumber sentinel
// values. "latest" maps to nil (geth toBlockNumArg treats nil as "latest").
// Non-nil values are negative sentinels that toBlockNumArg serialises via rpc.BlockNumber.String().
var BlockTagNumbers = map[string]*big.Int{
	"latest":    nil,
	"safe":      big.NewInt(int64(rpc.SafeBlockNumber)),
	"finalized": big.NewInt(int64(rpc.FinalizedBlockNumber)),
	"pending":   big.NewInt(int64(rpc.PendingBlockNumber)),
	"earliest":  big.NewInt(int64(rpc.EarliestBlockNumber)),
}
