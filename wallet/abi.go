package wallet

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/poteto-go/go-alchemy-sdk/constant"
)

// abiEncodeWords ABI-encodes each argument as a 32-byte word.
//
// Supported types:
//   - []byte    — written as-is (e.g. keccak256 type-hash, already 32 bytes)
//   - string    — interpreted as a hex address, left-padded to 32 bytes
//   - *big.Int  — left-padded to 32 bytes
//   - [32]byte  — written as-is
//   - uint8     — left-padded to 32 bytes
func abiEncodeWords(args ...any) []byte {
	var result []byte
	for _, arg := range args {
		switch v := arg.(type) {
		case []byte:
			result = append(result, v...)
		case string:
			result = append(result, common.LeftPadBytes(common.HexToAddress(v).Bytes(), constant.ABIWordSize)...)
		case *big.Int:
			result = append(result, common.LeftPadBytes(v.Bytes(), constant.ABIWordSize)...)
		case [32]byte:
			result = append(result, v[:]...)
		case uint8:
			result = append(result, common.LeftPadBytes([]byte{v}, constant.ABIWordSize)...)
		}
	}
	return result
}
