package utils

import (
	"crypto/rand"
	"math/big"
)

func NewAuthorizationNonce() [32]byte {
	var nonce [32]byte
	if _, err := rand.Read(nonce[:]); err != nil {
		panic(err)
	}
	return nonce
}

func RandomF64(max float64) float64 {
	const denom = 1 << 53
	nBig, err := rand.Int(rand.Reader, big.NewInt(denom))
	if err != nil {
		panic(err)
	}
	return float64(nBig.Int64()) / float64(denom) * max
}
