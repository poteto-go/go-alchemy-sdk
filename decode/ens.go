package decode

import (
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
)

// ENSNamehash implements EIP-137: recursively hashes ENS name labels with keccak256.
// An empty name returns the 32-byte zero node.
func ENSNamehash(name string) [32]byte {
	node := [32]byte{}
	if name == "" {
		return node
	}
	var buf [64]byte
	labels := strings.Split(name, ".")
	for i := len(labels) - 1; i >= 0; i-- {
		copy(buf[:32], node[:])
		copy(buf[32:], crypto.Keccak256([]byte(labels[i])))
		node = crypto.Keccak256Hash(buf[:])
	}
	return node
}
