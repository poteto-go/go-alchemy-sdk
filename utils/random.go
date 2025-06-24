package utils

import (
	"crypto/rand"
	"math/big"
)

func RandomF64(max float64) float64 {
	nBig, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		panic(err)
	}
	return float64(nBig.Int64() / int64(max))
}
