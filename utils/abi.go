package utils

import (
	"fmt"
	"math/big"
)

const (
	abiWordSize         = 32
	abiStringHeaderSize = 64 // offset(32) + length(32)
)

// DecodeABIString decodes an ABI-encoded string (offset, length, data).
func DecodeABIString(output []byte) (string, error) {
	if len(output) < abiStringHeaderSize {
		return "", fmt.Errorf("invalid ABI string: output too short, got %d bytes", len(output))
	}

	length := new(big.Int).SetBytes(output[abiWordSize : abiWordSize*2]).Int64()

	if int64(len(output)) < abiStringHeaderSize+length {
		return "", fmt.Errorf("invalid ABI string: declared length %d exceeds output size %d", length, len(output)-abiStringHeaderSize)
	}

	return string(output[abiStringHeaderSize : abiStringHeaderSize+length]), nil
}
