package utils

import (
	"fmt"
	"math/big"
)

const (
	// ABIWordSize is the byte length of one ABI-encoded word.
	ABIWordSize = 32
	// ABIStringHeaderSize is the byte length of the offset + length header for an ABI string.
	ABIStringHeaderSize = ABIWordSize * 2
)

// EncodeABIString encodes a string into ABI format (offset + length + data).
func EncodeABIString(s string) []byte {
	dataLen := len(s)
	paddedDataLen := ((dataLen + ABIWordSize - 1) / ABIWordSize) * ABIWordSize
	b := make([]byte, ABIStringHeaderSize+paddedDataLen)
	b[ABIWordSize-1] = byte(ABIWordSize)     // offset pointing to length field
	b[ABIStringHeaderSize-1] = byte(dataLen) // string length
	copy(b[ABIStringHeaderSize:], s)
	return b
}

// DecodeABIString decodes an ABI-encoded string (offset, length, data).
func DecodeABIString(output []byte) (string, error) {
	if len(output) < ABIStringHeaderSize {
		return "", fmt.Errorf("invalid ABI string: output too short, got %d bytes", len(output))
	}

	length := new(big.Int).SetBytes(output[ABIWordSize : ABIWordSize*2]).Int64()

	if int64(len(output)) < ABIStringHeaderSize+length {
		return "", fmt.Errorf("invalid ABI string: declared length %d exceeds output size %d", length, len(output)-ABIStringHeaderSize)
	}

	return string(output[ABIStringHeaderSize : ABIStringHeaderSize+length]), nil
}
