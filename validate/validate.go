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
	if u == nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return constant.ErrInvalidPrivateNetworkUrl
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

	// Compare as *big.Int to avoid Int64() overflow: a length whose lower 64
	// bits are >= 2^63 would become negative and slip past a naive bounds check.
	length := new(big.Int).SetBytes(output[constant.ABIWordSize : constant.ABIWordSize*2])
	maxLength := big.NewInt(int64(len(output) - constant.ABIStringHeaderSize))
	if length.Cmp(maxLength) > 0 {
		return constant.ErrInvalidABIString
	}

	return nil
}
