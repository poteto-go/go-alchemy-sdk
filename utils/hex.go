package utils

import (
	"math/big"
	"strconv"

	"github.com/poteto-go/go-alchemy-sdk/constant"
)

func FromHex(hexString string) (int, error) {
	if hexString == "" {
		return 0, nil
	}

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

func FromHexU64(hexString string) (uint64, error) {
	if hexString == "" {
		return uint64(0), nil
	}

	if len(hexString) <= 1 {
		return uint64(0), constant.ErrInvalidHexString
	}

	if hexString[0:2] != "0x" {
		return 0, constant.ErrInvalidHexString
	}

	// remove 0x
	numString := hexString[2:]
	num, err := strconv.ParseUint(numString, 16, 64)
	if err != nil {
		return 0, constant.ErrInvalidHexString
	}

	return num, nil
}

func FromBigHex(hexString string) (*big.Int, error) {
	if hexString == "" {
		return big.NewInt(0), nil
	}

	if len(hexString) <= 1 {
		return nil, constant.ErrInvalidHexString
	}

	if hexString[0:2] != "0x" {
		return nil, constant.ErrInvalidHexString
	}

	// remove 0x
	i := new(big.Int)
	num, ok := i.SetString(hexString[2:], 16)
	if !ok {
		return nil, constant.ErrInvalidHexString
	}

	if num.Cmp(big.NewInt(0)) == 0 {
		return big.NewInt(0), nil
	}

	return num, nil
}
