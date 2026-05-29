package utils

import (
	"math/big"

	"github.com/poteto-go/go-alchemy-sdk/constant"
)

func ToBlockNumber(blockTag string) (*big.Int, error) {
	if blockTag == "" {
		return nil, constant.ErrInvalidBlockTag
	}

	if num, ok := constant.BlockTagNumbers[blockTag]; ok {
		if num == nil {
			return nil, nil
		}
		return new(big.Int).Set(num), nil
	}

	blockNumber, err := FromBigHex(blockTag)
	if err != nil {
		return nil, err
	}

	return blockNumber, nil
}
