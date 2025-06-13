package utils

import (
	"strconv"

	"github.com/poteto-go/go-alchemy-sdk/core"
)

func FromHex(hexString string) (int, error) {
	if len(hexString) <= 1 {
		return 0, core.ErrInvalidHexString
	}

	if hexString[0:2] != "0x" {
		return 0, core.ErrInvalidHexString
	}

	// remove 0x
	numString := hexString[2:]
	num, err := strconv.ParseInt(numString, 16, 64)
	if err != nil {
		return 0, core.ErrInvalidHexString
	}

	return int(num), nil
}
