package utils

import (
	"fmt"
	"math/big"

	"github.com/poteto-go/go-alchemy-sdk/constant"
)

// https://docs.ethers.org/v5/api/providers/types/
func ValidateBlockTag(blockTag string) error {
	if _, err := ToBlockNumber(blockTag); err != nil {
		return constant.ErrInvalidBlockTag
	}
	return nil
}

// ValidateABIString validates an ABI-encoded string's header and declared
// length against the available output, returning the validated byte length.
func ValidateABIString(output []byte) (int64, error) {
	if len(output) < constant.ABIStringHeaderSize {
		return 0, fmt.Errorf("%w: output too short, got %d bytes", constant.ErrInvalidABIString, len(output))
	}

	// Compare as *big.Int to avoid Int64() overflow: a length whose lower 64
	// bits are >= 2^63 would become negative and slip past the bounds check,
	// causing a slice out-of-range panic (DoS on untrusted contract/RPC data).
	length := new(big.Int).SetBytes(output[constant.ABIWordSize : constant.ABIWordSize*2])
	maxLength := big.NewInt(int64(len(output) - constant.ABIStringHeaderSize))
	if length.Cmp(maxLength) > 0 {
		return 0, fmt.Errorf("%w: declared length %s exceeds output size %d", constant.ErrInvalidABIString, length.String(), len(output)-constant.ABIStringHeaderSize)
	}

	return length.Int64(), nil
}
