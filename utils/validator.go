package utils

import "github.com/poteto-go/go-alchemy-sdk/constant"

// https://docs.ethers.org/v5/api/providers/types/
func ValidateBlockTag(blockTag string) error {
	if _, err := ToBlockNumber(blockTag); err != nil {
		return constant.ErrInvalidBlockTag
	}
	return nil
}
