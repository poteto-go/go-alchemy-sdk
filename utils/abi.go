package utils

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/poteto-go/go-alchemy-sdk/constant"
)

// EncodeABIString encodes a string into ABI format (offset + length + data).
func EncodeABIString(s string) []byte {
	dataLen := len(s)
	paddedDataLen := ((dataLen + constant.ABIWordSize - 1) / constant.ABIWordSize) * constant.ABIWordSize
	b := make([]byte, constant.ABIStringHeaderSize+paddedDataLen)
	b[constant.ABIWordSize-1] = byte(constant.ABIWordSize) // offset pointing to length field
	b[constant.ABIStringHeaderSize-1] = byte(dataLen)      // string length
	copy(b[constant.ABIStringHeaderSize:], s)
	return b
}

// DecodeABIAddress decodes an ABI-encoded address (left-padded 32-byte word).
func DecodeABIAddress(output []byte) common.Address {
	if len(output) < constant.ABIWordSize {
		return common.Address{}
	}
	return common.BytesToAddress(output[constant.ABIAddressOffset:constant.ABIWordSize])
}

// DecodeABIString decodes an ABI-encoded string (offset, length, data).
func DecodeABIString(output []byte) (string, error) {
	if len(output) < constant.ABIStringHeaderSize {
		return "", fmt.Errorf("invalid ABI string: output too short, got %d bytes", len(output))
	}

	length := new(big.Int).SetBytes(output[constant.ABIWordSize : constant.ABIWordSize*2]).Int64()

	if int64(len(output)) < constant.ABIStringHeaderSize+length {
		return "", fmt.Errorf("invalid ABI string: declared length %d exceeds output size %d", length, len(output)-constant.ABIStringHeaderSize)
	}

	return string(output[constant.ABIStringHeaderSize : constant.ABIStringHeaderSize+length]), nil
}
