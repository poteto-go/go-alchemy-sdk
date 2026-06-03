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
)
