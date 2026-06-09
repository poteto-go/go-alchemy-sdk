package utils

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/poteto-go/go-alchemy-sdk/constant"
	"github.com/poteto-go/go-alchemy-sdk/validate"
	"golang.org/x/crypto/sha3"
)

// EncodeReadCalldata builds eth_call calldata for a contract read: the 4-byte
// selector (keccak256(signature)[:4]) followed by the ABI-encoded args.
func EncodeReadCalldata(signature []byte, args ...[]byte) []byte {
	hash := sha3.NewLegacyKeccak256()
	// keccak Write never returns an error.
	hash.Write(signature)
	methodID := hash.Sum(nil)[:4]

	data := make([]byte, 0, len(methodID)+len(args)*constant.ABIWordSize)
	data = append(data, methodID...)
	for _, arg := range args {
		data = append(data, arg...)
	}
	return data
}

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

// EncodeABIAddress left-pads an address to a 32-byte ABI word.
func EncodeABIAddress(address string) []byte {
	return common.LeftPadBytes(common.HexToAddress(address).Bytes(), constant.ABIWordSize)
}

// DecodeUint256 decodes a 32-byte ABI word as an unsigned integer.
func DecodeUint256(output []byte) (*big.Int, error) {
	return new(big.Int).SetBytes(output), nil
}

// DecodeBool decodes an ABI bool (non-zero last byte is true).
func DecodeBool(output []byte) (bool, error) {
	return len(output) > 0 && output[len(output)-1] == 1, nil
}

// DecodeUint8 decodes a 32-byte ABI word as a uint8, erroring on overflow.
func DecodeUint8(output []byte) (uint8, error) {
	out := new(big.Int).SetBytes(output)
	if out.BitLen() > 8 {
		return 0, fmt.Errorf("uint8 overflow: %s", out.String())
	}
	b := out.Bytes()
	if len(b) == 0 {
		return 0, nil
	}
	return b[len(b)-1], nil
}

// DecodeBytes32 decodes the first 32-byte ABI word as a fixed [32]byte.
func DecodeBytes32(output []byte) ([32]byte, error) {
	if len(output) < constant.ABIWordSize {
		return [32]byte{}, fmt.Errorf("unexpected output length: %d", len(output))
	}
	var result [32]byte
	copy(result[:], output[:constant.ABIWordSize])
	return result, nil
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
