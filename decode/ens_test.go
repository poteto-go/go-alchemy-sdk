package decode_test

import (
	"encoding/hex"
	"testing"

	"github.com/poteto-go/go-alchemy-sdk/decode"
	"github.com/stretchr/testify/assert"
)

func TestENSNamehash(t *testing.T) {
	t.Run("empty string returns zero node", func(t *testing.T) {
		got := decode.ENSNamehash("")
		assert.Equal(t, [32]byte{}, got)
	})

	t.Run("eth returns known hash", func(t *testing.T) {
		// namehash("eth") = keccak256(0x00...00 ++ keccak256("eth"))
		// known value: 0x93cdeb708b7545dc668eb9280176169d1c33cfd8ed6f04690a0bcc88a93fc4ae
		want, _ := hex.DecodeString("93cdeb708b7545dc668eb9280176169d1c33cfd8ed6f04690a0bcc88a93fc4ae")
		var wantArr [32]byte
		copy(wantArr[:], want)
		got := decode.ENSNamehash("eth")
		assert.Equal(t, wantArr, got)
	})

	t.Run("vitalik.eth returns known hash", func(t *testing.T) {
		// namehash("vitalik.eth") known value
		want, _ := hex.DecodeString("ee6c4522aab0003e8d14cd40a6af439055fd2577951148c14b6cea9a53475835")
		var wantArr [32]byte
		copy(wantArr[:], want)
		got := decode.ENSNamehash("vitalik.eth")
		assert.Equal(t, wantArr, got)
	})

	t.Run("single label with no dots", func(t *testing.T) {
		// result must not equal zero node
		got := decode.ENSNamehash("foo")
		assert.NotEqual(t, [32]byte{}, got)
	})

	t.Run("deterministic: same input produces same output", func(t *testing.T) {
		a := decode.ENSNamehash("alice.eth")
		b := decode.ENSNamehash("alice.eth")
		assert.Equal(t, a, b)
	})

	t.Run("different names produce different hashes", func(t *testing.T) {
		a := decode.ENSNamehash("alice.eth")
		b := decode.ENSNamehash("bob.eth")
		assert.NotEqual(t, a, b)
	})
}
