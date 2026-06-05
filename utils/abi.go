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
	// Write the length as a big-endian 32-byte word so strings of 256+ bytes
	// encode correctly (a single byte() would truncate the length).
	big.NewInt(int64(dataLen)).FillBytes(b[constant.ABIWordSize:constant.ABIStringHeaderSize])
	copy(b[constant.ABIStringHeaderSize:], s)
	return b
}

// DecodeABIAddress decodes an ABI-encoded address (left-padded 32-byte word).
func DecodeABIAddress(output []byte) (common.Address, error) {
	if len(output) < constant.ABIWordSize {
		return common.Address{}, fmt.Errorf("invalid ABI address: output too short, got %d bytes", len(output))
	}
	return common.BytesToAddress(output[constant.ABIAddressOffset:constant.ABIWordSize]), nil
}

// DecodeABIString decodes an ABI-encoded string (offset, length, data).
func DecodeABIString(output []byte) (string, error) {
	if len(output) < constant.ABIStringHeaderSize {
		return "", fmt.Errorf("invalid ABI string: output too short, got %d bytes", len(output))
	}

	// Compare as *big.Int to avoid Int64() overflow: a length whose lower 64
	// bits are >= 2^63 would become negative and slip past the bounds check,
	// causing a slice out-of-range panic (DoS on untrusted contract/RPC data).
	length := new(big.Int).SetBytes(output[constant.ABIWordSize : constant.ABIWordSize*2])
	maxLength := big.NewInt(int64(len(output) - constant.ABIStringHeaderSize))
	if length.Cmp(maxLength) > 0 {
		return "", fmt.Errorf("invalid ABI string: declared length %s exceeds output size %d", length.String(), len(output)-constant.ABIStringHeaderSize)
	}

	l := length.Int64()
	return string(output[constant.ABIStringHeaderSize : constant.ABIStringHeaderSize+l]), nil
}
