package utils

import (
	"strings"

	"github.com/poteto-go/go-alchemy-sdk/constant"
)

// https://docs.ethers.org/v5/api/providers/types/
func ValidateBlockTag(blockTag string) error {
	if _, ok := constant.BlockTagNumbers[blockTag]; ok {
		return nil
	}

	if strings.HasPrefix(blockTag, "0x") {
		_, err := FromBigHex(blockTag)
		if err != nil {
			return constant.ErrInvalidBlockTag
		}
		return nil
	}

	return constant.ErrInvalidBlockTag
}
