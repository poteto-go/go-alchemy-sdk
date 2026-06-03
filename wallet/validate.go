package wallet

import (
	"math/big"

	"github.com/poteto-go/go-alchemy-sdk/validate"
)

func validateUint256(v *big.Int) error  { return validate.Uint256(v) }
func validateAddress(addr string) error { return validate.Address(addr) }
