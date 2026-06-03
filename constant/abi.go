package constant

// ABI encoding parameters.
const (
	// ABIWordSize is the byte length of a single ABI-encoded word.
	ABIWordSize = 32

	// ABIAddressOffset is the byte offset within an ABI word where a 20-byte
	// Ethereum address begins (left-padded with 12 zero bytes).
	ABIAddressOffset = ABIWordSize - 20

	// ABIStringHeaderSize is the byte length of the offset + length header
	// for an ABI-encoded dynamic string.
	ABIStringHeaderSize = ABIWordSize * 2

	// ECDSALegacyVOffset is added to the ECDSA recovery id (0 or 1) returned
	// by crypto.Sign to produce the Ethereum v value (27 or 28).
	ECDSALegacyVOffset = 27
)
