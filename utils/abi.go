package utils

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/poteto-go/go-alchemy-sdk/constant"
	"github.com/poteto-go/go-alchemy-sdk/validate"
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
//
// ABI-encoded string layout:
//
//	[0x00..20  (32 bytes)] offset  - points to the start of the length field (always 0x20)
//	[0x00..09  (32 bytes)] length  - byte length of the string (e.g. 9 for "TestToken")
//	[data..00  (N bytes) ] data    - UTF-8 string bytes, zero-padded to a 32-byte boundary
//
// output[64 : 64+length] is the actual string data.
func DecodeABIString(output []byte) (string, error) {
	if err := validate.ABIString(output); err != nil {
		return "", err
	}

	// Safe: validate.ABIString guarantees the declared length fits within output.
	length := new(big.Int).SetBytes(output[constant.ABIWordSize : constant.ABIWordSize*2]).Int64()
	return string(output[constant.ABIStringHeaderSize : constant.ABIStringHeaderSize+length]), nil
}
