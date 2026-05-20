package utils

import (
	"fmt"
	"math/big"
)

// DecodeABIString decodes an ABI-encoded string (offset, length, data).
func DecodeABIString(output []byte) (string, error) {
	if len(output) < 64 {
		return "", fmt.Errorf("invalid ABI string: output too short, got %d bytes", len(output))
	}
	
	// The offset is at output[0:32], but we already know we're reading a string.
	// We read length from output[32:64]
	length := new(big.Int).SetBytes(output[32:64]).Int64()
	
	if int64(len(output)) < 64+length {
		return "", fmt.Errorf("invalid ABI string: declared length %d exceeds output size %d", length, len(output)-64)
	}
	
	return string(output[64 : 64+length]), nil
}
