package utils

import (
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/poteto-go/go-alchemy-sdk/types"
)

func TransformAlchemyBlock(gethBlock *gethTypes.Block) types.Block {
	return types.Block{
		Hash:         gethBlock.Hash().Hex(),
		ParentHash:   gethBlock.ParentHash().Hex(),
		Number:       gethBlock.Number(),
		Timestamp:    gethBlock.Time(),
		Nonce:        gethBlock.Nonce(),
		Difficulty:   gethBlock.Difficulty(),
		GasLimit:     gethBlock.GasLimit(),
		GasUsed:      gethBlock.GasUsed(),
		Miner:        gethBlock.Coinbase().Hex(),
		Transactions: make([]string, len(gethBlock.Transactions())),
	}
}
