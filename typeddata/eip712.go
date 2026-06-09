package typeddata

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/poteto-go/go-alchemy-sdk/constant"
	"github.com/poteto-go/go-alchemy-sdk/encode"
	"github.com/poteto-go/go-alchemy-sdk/types"
)

// EncodeWords ABI-encodes each argument as a 32-byte word for use with SignEIP712.
//
// Supported types:
//   - []byte    — written as-is (e.g. keccak256 type-hash, already 32 bytes)
//   - string    — interpreted as a hex address, left-padded to 32 bytes
//   - *big.Int  — left-padded to 32 bytes
//   - [32]byte  — written as-is
//   - uint8     — left-padded to 32 bytes
func EncodeWords(args ...any) []byte {
	result := make([]byte, 0, len(args)*constant.ABIWordSize)
	for _, arg := range args {
		switch v := arg.(type) {
		case []byte:
			result = append(result, v...)
		case string:
			result = append(result, common.LeftPadBytes(common.HexToAddress(v).Bytes(), constant.ABIWordSize)...)
		case *big.Int:
			result = append(result, common.LeftPadBytes(v.Bytes(), constant.ABIWordSize)...)
		case [32]byte:
			result = append(result, v[:]...)
		case uint8:
			result = append(result, common.LeftPadBytes([]byte{v}, constant.ABIWordSize)...)
		}
	}
	return result
}

/*
EIP-712 signing by str key.

EIP-712 is a standard for hashing and signing of typed structured data, as opposed to just arbitrary bytes.
It is used to prevent signing of unintended data and to make the signed data more human-readable.

refs: https://eips.ethereum.org/EIPS/eip-712
*/
func SignEIP712Str(
	privateKeyStr string, domainSeparator [32]byte, encoded []byte,
) (types.Signature, error) {
	privateKey, err := encode.PrivateKey(privateKeyStr)
	if err != nil {
		return types.Signature{}, err
	}

	return SignEIP712(privateKey, domainSeparator, encoded)
}

/*
EIP-712 signing.

EIP-712 is a standard for hashing and signing of typed structured data, as opposed to just arbitrary bytes.
It is used to prevent signing of unintended data and to make the signed data more human-readable.

refs: https://eips.ethereum.org/EIPS/eip-712
*/
func SignEIP712(
	privateKey *ecdsa.PrivateKey, domainSeparator [32]byte, encoded []byte,
) (types.Signature, error) {
	structHash := crypto.Keccak256(encoded)

	msg := make([]byte, 0, 2+constant.ABIWordSize*2)
	msg = append(msg, constant.EIP191DataPrefix, constant.EIP712StructuredDataVersion)
	msg = append(msg, domainSeparator[:]...)
	msg = append(msg, structHash...)
	hash := crypto.Keccak256(msg)

	sig, err := crypto.Sign(hash, privateKey)
	if err != nil {
		return types.Signature{}, err
	}

	var r, s [32]byte
	copy(r[:], sig[:constant.ABIWordSize])
	copy(s[:], sig[constant.ABIWordSize:constant.ABIWordSize*2])
	v := sig[constant.ABIWordSize*2] + constant.ECDSALegacyVOffset

	return types.Signature{
		V: v,
		R: r,
		S: s,
	}, nil
}
