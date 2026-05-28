package utils

import (
	"math/big"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/poteto-go/go-alchemy-sdk/constant"
)

// namedTagNumbers maps EIP-1898 named block tags to their geth rpc.BlockNumber sentinel
// values. Both ToBlockNumber and ValidateBlockTag reference this map so the tag list
// stays in one place.
var namedTagNumbers = map[string]*big.Int{
	"safe":      big.NewInt(int64(rpc.SafeBlockNumber)),
	"finalized": big.NewInt(int64(rpc.FinalizedBlockNumber)),
	"pending":   big.NewInt(int64(rpc.PendingBlockNumber)),
	"earliest":  big.NewInt(int64(rpc.EarliestBlockNumber)),
}

func ToBlockNumber(blockTag string) (*big.Int, error) {
	if blockTag == "latest" {
		return nil, nil
	}
	if num, ok := namedTagNumbers[blockTag]; ok {
		return new(big.Int).Set(num), nil
	}
	blockNumber, err := FromBigHex(blockTag)
	if err != nil {
		return nil, constant.ErrInvalidBlockTag
	}
	return blockNumber, nil
}
