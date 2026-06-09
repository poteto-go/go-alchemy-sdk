package encode

import (
	"encoding/binary"

	"github.com/ethereum/go-ethereum/common"
	"github.com/poteto-go/go-alchemy-sdk/constant"
	"golang.org/x/crypto/sha3"
)

// ReadCalldata builds eth_call calldata for a contract read: the 4-byte
// selector (keccak256(signature)[:4]) followed by the ABI-encoded args.
func ReadCalldata(signature []byte, args ...[]byte) []byte {
	hash := sha3.NewLegacyKeccak256()
	// keccak Write never returns an error.
	hash.Write(signature)
	methodID := hash.Sum(nil)[:4]

	total := 4
	for _, arg := range args {
		total += len(arg)
	}
	data := make([]byte, 0, total)
	data = append(data, methodID...)
	for _, arg := range args {
		data = append(data, arg...)
	}
	return data
}

// ABIString encodes a string into ABI format (offset + length + data).
func ABIString(s string) []byte {
	dataLen := len(s)
	paddedDataLen := ((dataLen + constant.ABIWordSize - 1) / constant.ABIWordSize) * constant.ABIWordSize
	b := make([]byte, constant.ABIStringHeaderSize+paddedDataLen)
	b[constant.ABIWordSize-1] = byte(constant.ABIWordSize) // offset pointing to length field
	binary.BigEndian.PutUint64(b[constant.ABIStringHeaderSize-8:constant.ABIStringHeaderSize], uint64(dataLen))
	copy(b[constant.ABIStringHeaderSize:], s)
	return b
}

// ABIAddress left-pads an address to a 32-byte ABI word.
func ABIAddress(address string) []byte {
	return common.LeftPadBytes(common.HexToAddress(address).Bytes(), constant.ABIWordSize)
}
