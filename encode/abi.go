package encode

import (
	"encoding/binary"
	"math/big"

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

// ABIUint256 encodes a *big.Int as a 32-byte ABI word (big-endian, left-padded with zeros).
func ABIUint256(v *big.Int) []byte {
	return common.LeftPadBytes(v.Bytes(), constant.ABIWordSize)
}

// ABIBool encodes a bool as a 32-byte ABI word (0x00..01 for true, all zeros
// for false).
func ABIBool(v bool) []byte {
	b := make([]byte, constant.ABIWordSize)
	if v {
		b[constant.ABIWordSize-1] = 1
	}
	return b
}

// ABIBytes encodes the tail of a dynamic `bytes` argument: a 32-byte length
// word followed by the data right-padded to a word boundary. The caller is
// responsible for emitting the preceding offset word that points at this tail.
func ABIBytes(data []byte) []byte {
	paddedDataLen := ((len(data) + constant.ABIWordSize - 1) / constant.ABIWordSize) * constant.ABIWordSize
	b := make([]byte, constant.ABIWordSize+paddedDataLen)
	// First word is the byte length, encoded as a 32-byte ABI uint256.
	copy(b, ABIUint256(big.NewInt(int64(len(data)))))
	copy(b[constant.ABIWordSize:], data)
	return b
}

// ABIUint256Array encodes the tail of a dynamic `uint256[]` argument: a 32-byte
// length word followed by one ABI word per element. The caller is responsible
// for emitting the preceding offset word that points at this tail.
func ABIUint256Array(vs []*big.Int) []byte {
	b := make([]byte, 0, constant.ABIWordSize*(1+len(vs)))
	b = append(b, ABIUint256(big.NewInt(int64(len(vs))))...)
	for _, v := range vs {
		b = append(b, ABIUint256(v)...)
	}
	return b
}

// ABIAddressArray encodes the tail of a dynamic `address[]` argument: a 32-byte
// length word followed by one left-padded ABI word per address. The caller is
// responsible for emitting the preceding offset word that points at this tail.
func ABIAddressArray(addresses []string) []byte {
	b := make([]byte, 0, constant.ABIWordSize*(1+len(addresses)))
	b = append(b, ABIUint256(big.NewInt(int64(len(addresses))))...)
	for _, addr := range addresses {
		b = append(b, ABIAddress(addr)...)
	}
	return b
}

// ABIDynamicArgs packs a sequence of all-dynamic arguments (e.g. the tails from
// ABIBytes / ABIUint256Array / ABIAddressArray) into a single calldata blob:
// one offset word per argument forming the head, followed by every tail. Each
// offset points at where its tail begins relative to the start of the head.
// Use this for methods whose arguments are all dynamic, such as
// balanceOfBatch(address[],uint256[]).
func ABIDynamicArgs(tails ...[]byte) []byte {
	headSize := len(tails) * constant.ABIWordSize
	total := headSize
	for _, tail := range tails {
		total += len(tail)
	}

	b := make([]byte, 0, total)
	offset := headSize
	for _, tail := range tails {
		b = append(b, ABIUint256(big.NewInt(int64(offset)))...)
		offset += len(tail)
	}
	for _, tail := range tails {
		b = append(b, tail...)
	}
	return b
}
