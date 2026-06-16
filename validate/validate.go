package validate

import (
	"math/big"
	"net/url"

	"github.com/ethereum/go-ethereum/common"
	"github.com/poteto-go/go-alchemy-sdk/constant"
)

func Uint256(v *big.Int) error {
	if v == nil {
		return constant.ErrNilAmount
	}
	if v.Sign() < 0 {
		return constant.ErrNegativeAmount
	}
	if v.BitLen() > 256 {
		return constant.ErrAmountExceedsUint256
	}
	return nil
}

func Address(addr string) error {
	if !common.IsHexAddress(addr) {
		return constant.ErrInvalidAddress
	}
	return nil
}

func Addresses(addrs ...string) error {
	for _, addr := range addrs {
		if err := Address(addr); err != nil {
			return err
		}
	}
	return nil
}

// BlockTag validates an EIP-1898 block tag: a named tag (latest/safe/...) or a
// hex-encoded block number.
// https://docs.ethers.org/v5/api/providers/types/
func BlockTag(blockTag string) error {
	if blockTag == "" {
		return constant.ErrInvalidBlockTag
	}
	if _, ok := constant.BlockTagNumbers[blockTag]; ok {
		return nil
	}
	if len(blockTag) <= 2 || blockTag[:2] != "0x" {
		return constant.ErrInvalidBlockTag
	}
	if _, ok := new(big.Int).SetString(blockTag[2:], 16); !ok {
		return constant.ErrInvalidBlockTag
	}
	return nil
}

func Url(rawUrl string) error {
	if rawUrl == "" {
		return nil
	}
	u, _ := url.Parse(rawUrl)
	if u == nil || (u.Scheme != "http" && u.Scheme != "https") || u.Hostname() == "" {
		return constant.ErrInvalidPrivateNetworkUrl
	}
	return nil
}

// declaredLengthExceedsAvailable returns true when a declared length (items or
// bytes) would read past the end of the available data region.
// Compare as *big.Int to avoid Int64() overflow: a length whose lower 64
// bits are >= 2^63 would become negative and slip past a naive bounds check.
func declaredLengthExceedsAvailable(declared *big.Int, available int) bool {
	return declared.Cmp(big.NewInt(int64(available))) > 0
}

// abiOffsetIsStandard reports whether the leading 32-byte word equals the
// standard 0x20 offset used for a single dynamic return value. An offset word
// pointing elsewhere signals a malformed/malicious response that would read
// the length from the wrong position. The standard offset is 31 zero bytes
// followed by 0x20 (== constant.ABIWordSize).
func abiOffsetIsStandard(output []byte) bool {
	word := output[:constant.ABIWordSize]
	for _, b := range word[:constant.ABIWordSize-1] {
		if b != 0 {
			return false
		}
	}
	return word[constant.ABIWordSize-1] == constant.ABIWordSize
}

// ABIUint256Array validates an ABI-encoded uint256[] return value
// (offset word + length word + element words) against the available output,
// guarding callers against a slice out-of-range panic.
func ABIUint256Array(output []byte) error {
	if len(output) < constant.ABIWordSize*2 {
		return constant.ErrInvalidABIArray
	}
	if !abiOffsetIsStandard(output) {
		return constant.ErrInvalidABIArray
	}
	dataStart := constant.ABIWordSize * 2
	length := new(big.Int).SetBytes(output[constant.ABIWordSize:dataStart])
	if declaredLengthExceedsAvailable(length, (len(output)-dataStart)/constant.ABIWordSize) {
		return constant.ErrInvalidABIArray
	}
	return nil
}

// ABIString validates an ABI-encoded string's header and declared length
// against the available output, guarding callers against a slice
// out-of-range panic on malicious/corrupt contract or RPC data.
func ABIString(output []byte) error {
	if len(output) < constant.ABIStringHeaderSize {
		return constant.ErrInvalidABIString
	}
	if !abiOffsetIsStandard(output) {
		return constant.ErrInvalidABIString
	}
	length := new(big.Int).SetBytes(output[constant.ABIWordSize : constant.ABIWordSize*2])
	if declaredLengthExceedsAvailable(length, len(output)-constant.ABIStringHeaderSize) {
		return constant.ErrInvalidABIString
	}
	return nil
}
