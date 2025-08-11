package utils

import (
	"math/big"
	"strconv"

	"github.com/poteto-go/go-alchemy-sdk/constant"
)

func FromHex(hexString string) (int, error) {
	if len(hexString) <= 1 {
		return 0, constant.ErrInvalidHexString
	}

	if hexString[0:2] != "0x" {
		return 0, constant.ErrInvalidHexString
	}

	// remove 0x
	numString := hexString[2:]
	num, err := strconv.ParseInt(numString, 16, 64)
	if err != nil {
		return 0, constant.ErrInvalidHexString
	}

	return int(num), nil
}

func FromBigHex(hexString string) (*big.Int, error) {
	if len(hexString) <= 1 {
		return big.NewInt(0), constant.ErrInvalidHexString
	}

	if hexString[0:2] != "0x" {
		return big.NewInt(0), constant.ErrInvalidHexString
	}

	// remove 0x
	i := new(big.Int)
	num, ok := i.SetString(hexString[2:], 16)
	if !ok {
		return big.NewInt(0), constant.ErrInvalidHexString
	}
	return num, nil
}
