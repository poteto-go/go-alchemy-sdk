package typeddata_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/poteto-go/go-alchemy-sdk/constant"
	"github.com/poteto-go/go-alchemy-sdk/encode"
	"github.com/poteto-go/go-alchemy-sdk/typeddata"
	"github.com/stretchr/testify/assert"
)

const testPrivateKey = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"

func TestEncodeWords(t *testing.T) {
	tests := []struct {
		name string
		args []any
		want int // expected byte length
	}{
		{
			name: "empty args",
			args: []any{},
			want: 0,
		},
		{
			name: "[]byte (32-byte hash)",
			args: []any{make([]byte, constant.ABIWordSize)},
			want: constant.ABIWordSize,
		},
		{
			name: "string address left-padded to 32 bytes",
			args: []any{"0x742d35Cc6634C0532925a3b844Bc454e4438f44e"},
			want: constant.ABIWordSize,
		},
		{
			name: "*big.Int left-padded to 32 bytes",
			args: []any{big.NewInt(1000)},
			want: constant.ABIWordSize,
		},
		{
			name: "[32]byte written as-is",
			args: []any{[32]byte{0x01}},
			want: constant.ABIWordSize,
		},
		{
			name: "uint8 left-padded to 32 bytes",
			args: []any{uint8(27)},
			want: constant.ABIWordSize,
		},
		{
			name: "multiple args concatenated",
			args: []any{
				make([]byte, constant.ABIWordSize),
				big.NewInt(500),
				"0x742d35Cc6634C0532925a3b844Bc454e4438f44e",
			},
			want: constant.ABIWordSize * 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := typeddata.EncodeWords(tt.args...)
			assert.Equal(t, tt.want, len(got))
		})
	}
}

func TestEncodeWords_addressPadding(t *testing.T) {
	addr := "0x742d35Cc6634C0532925a3b844Bc454e4438f44e"
	got := typeddata.EncodeWords(addr)

	// first 12 bytes must be zero (left-padded)
	for i := 0; i < 12; i++ {
		assert.Equal(t, byte(0), got[i], "byte %d should be zero", i)
	}
	assert.Len(t, got, constant.ABIWordSize)
}

func TestEncodeWords_bigIntPadding(t *testing.T) {
	val := big.NewInt(1)
	got := typeddata.EncodeWords(val)

	assert.Len(t, got, constant.ABIWordSize)
	// last byte should be 1
	assert.Equal(t, byte(1), got[constant.ABIWordSize-1])
	// first 31 bytes should be zero
	for i := 0; i < constant.ABIWordSize-1; i++ {
		assert.Equal(t, byte(0), got[i], "byte %d should be zero", i)
	}
}

func TestSignEIP712(t *testing.T) {
	privateKey, err := encode.PrivateKey(testPrivateKey)
	assert.NoError(t, err)

	var domainSeparator [32]byte
	copy(domainSeparator[:], crypto.Keccak256([]byte("TestDomain")))

	encoded := typeddata.EncodeWords(
		crypto.Keccak256([]byte("Transfer(address to,uint256 value)")),
		"0x742d35Cc6634C0532925a3b844Bc454e4438f44e",
		big.NewInt(1000),
	)

	sig, err := typeddata.SignEIP712(privateKey, domainSeparator, encoded)
	assert.NoError(t, err)

	// V must be 27 or 28 (legacy offset)
	assert.True(t, sig.V == 27 || sig.V == 28, "V should be 27 or 28, got %d", sig.V)
	// R and S must be non-zero
	assert.NotEqual(t, [32]byte{}, sig.R)
	assert.NotEqual(t, [32]byte{}, sig.S)
}

func TestSignEIP712_deterministicForSameInput(t *testing.T) {
	privateKey, err := encode.PrivateKey(testPrivateKey)
	assert.NoError(t, err)

	var domainSeparator [32]byte
	copy(domainSeparator[:], crypto.Keccak256([]byte("TestDomain")))
	encoded := typeddata.EncodeWords(big.NewInt(42))

	sig1, err := typeddata.SignEIP712(privateKey, domainSeparator, encoded)
	assert.NoError(t, err)

	sig2, err := typeddata.SignEIP712(privateKey, domainSeparator, encoded)
	assert.NoError(t, err)

	assert.Equal(t, sig1, sig2)
}

func TestSignEIP712_differentDomainProducesDifferentSig(t *testing.T) {
	privateKey, err := encode.PrivateKey(testPrivateKey)
	assert.NoError(t, err)

	encoded := typeddata.EncodeWords(big.NewInt(1))

	var domain1 [32]byte
	copy(domain1[:], crypto.Keccak256([]byte("Domain1")))

	var domain2 [32]byte
	copy(domain2[:], crypto.Keccak256([]byte("Domain2")))

	sig1, err := typeddata.SignEIP712(privateKey, domain1, encoded)
	assert.NoError(t, err)

	sig2, err := typeddata.SignEIP712(privateKey, domain2, encoded)
	assert.NoError(t, err)

	assert.NotEqual(t, sig1, sig2)
}

func TestSignEIP712Str(t *testing.T) {
	var domainSeparator [32]byte
	copy(domainSeparator[:], crypto.Keccak256([]byte("TestDomain")))

	encoded := typeddata.EncodeWords(big.NewInt(1))

	sig, err := typeddata.SignEIP712Str(testPrivateKey, domainSeparator, encoded)
	assert.NoError(t, err)
	assert.True(t, sig.V == 27 || sig.V == 28)
	assert.NotEqual(t, [32]byte{}, sig.R)
	assert.NotEqual(t, [32]byte{}, sig.S)
}

func TestSignEIP712Str_with0xPrefix(t *testing.T) {
	var domainSeparator [32]byte
	copy(domainSeparator[:], crypto.Keccak256([]byte("TestDomain")))
	encoded := typeddata.EncodeWords(big.NewInt(1))

	sigWithout, err := typeddata.SignEIP712Str(testPrivateKey, domainSeparator, encoded)
	assert.NoError(t, err)

	sigWith, err := typeddata.SignEIP712Str("0x"+testPrivateKey, domainSeparator, encoded)
	assert.NoError(t, err)

	assert.Equal(t, sigWithout, sigWith)
}

func TestSignEIP712Str_invalidKey(t *testing.T) {
	var domainSeparator [32]byte
	encoded := typeddata.EncodeWords(big.NewInt(1))

	_, err := typeddata.SignEIP712Str("not-a-valid-hex-key", domainSeparator, encoded)
	assert.Error(t, err)
}

func TestSignEIP712Str_matchesSignEIP712(t *testing.T) {
	privateKey, err := encode.PrivateKey(testPrivateKey)
	assert.NoError(t, err)

	var domainSeparator [32]byte
	copy(domainSeparator[:], crypto.Keccak256([]byte("TestDomain")))
	encoded := typeddata.EncodeWords(big.NewInt(999))

	sigFromKey, err := typeddata.SignEIP712(privateKey, domainSeparator, encoded)
	assert.NoError(t, err)

	sigFromStr, err := typeddata.SignEIP712Str(testPrivateKey, domainSeparator, encoded)
	assert.NoError(t, err)

	assert.Equal(t, sigFromKey, sigFromStr)
}
