package utils

import "math/big"

func ToBlockNumber(blockTag string) (*big.Int, error) {
	if blockTag == "latest" {
		return nil, nil
	}

	blockNumber, err := FromBigHex(blockTag)
	if err != nil {
		return nil, err
	}

	return blockNumber, nil
}
