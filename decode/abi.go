package decode

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/poteto-go/go-alchemy-sdk/constant"
	"github.com/poteto-go/go-alchemy-sdk/validate"
)

// Uint256 decodes a 32-byte ABI word as an unsigned integer.
func Uint256(output []byte) (*big.Int, error) {
	return new(big.Int).SetBytes(output), nil
}

// Uint256Array decodes an ABI-encoded dynamic uint256[] return value
// (offset, length, items). The leading offset word points at the length field
// and is assumed to be the standard 0x20 for a single dynamic return.
func Uint256Array(output []byte) ([]*big.Int, error) {
	if err := validate.ABIUint256Array(output); err != nil {
		return nil, err
	}
	dataStart := constant.ABIWordSize * 2
	// Safe: validate.ABIUint256Array guarantees length fits within output.
	n := int(new(big.Int).SetBytes(output[constant.ABIWordSize:dataStart]).Int64())
	result := make([]*big.Int, n)
	for i := 0; i < n; i++ {
		start := dataStart + i*constant.ABIWordSize
		result[i] = new(big.Int).SetBytes(output[start : start+constant.ABIWordSize])
	}
	return result, nil
}

// Bool decodes an ABI bool (non-zero last byte is true).
func Bool(output []byte) (bool, error) {
	return len(output) > 0 && output[len(output)-1] == 1, nil
}

// Uint8 decodes a 32-byte ABI word as a uint8, erroring on overflow.
func Uint8(output []byte) (uint8, error) {
	if len(output) == 0 {
		return 0, nil
	}
	for _, b := range output[:len(output)-1] {
		if b != 0 {
			return 0, fmt.Errorf("uint8 overflow")
		}
	}
	return output[len(output)-1], nil
}

// Bytes32 decodes the first 32-byte ABI word as a fixed [32]byte.
func Bytes32(output []byte) ([32]byte, error) {
	if len(output) < constant.ABIWordSize {
		return [32]byte{}, fmt.Errorf("unexpected output length: %d", len(output))
	}
	var result [32]byte
	copy(result[:], output[:constant.ABIWordSize])
	return result, nil
}

// ABIAddress decodes an ABI-encoded address (left-padded 32-byte word).
func ABIAddress(output []byte) (common.Address, error) {
	if len(output) < constant.ABIWordSize {
		return common.Address{}, fmt.Errorf("invalid ABI address: output too short, got %d bytes", len(output))
	}
	return common.BytesToAddress(output[constant.ABIAddressOffset:constant.ABIWordSize]), nil
}

// ABIString decodes an ABI-encoded string (offset, length, data).
//
// ABI-encoded string layout:
//
//	[0x00..20  (32 bytes)] offset  - points to the start of the length field (always 0x20)
//	[0x00..09  (32 bytes)] length  - byte length of the string (e.g. 9 for "TestToken")
//	[data..00  (N bytes) ] data    - UTF-8 string bytes, zero-padded to a 32-byte boundary
//
// output[64 : 64+length] is the actual string data.
func ABIString(output []byte) (string, error) {
	if err := validate.ABIString(output); err != nil {
		return "", err
	}

	// Safe: validate.ABIString guarantees the declared length fits within output.
	length := new(big.Int).SetBytes(output[constant.ABIWordSize : constant.ABIWordSize*2]).Int64()
	return string(output[constant.ABIStringHeaderSize : constant.ABIStringHeaderSize+length]), nil
}
