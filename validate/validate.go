package validate

import (
	"math/big"

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
