package utils

import (
	"crypto/rand"
	"math/big"
)

func RandomF64(max float64) float64 {
	const denom = 1 << 53
	nBig, err := rand.Int(rand.Reader, big.NewInt(denom))
	if err != nil {
		panic(err)
	}
	return float64(nBig.Int64()) / float64(denom) * max
}
