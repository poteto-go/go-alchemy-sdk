package constant

// ABI encoding parameters.
const (
	// ABIWordSize is the byte length of a single ABI-encoded word.
	ABIWordSize = 32

	// ABIAddressOffset is the byte offset within an ABI word where a 20-byte
	// Ethereum address begins (left-padded with 12 zero bytes).
	ABIAddressOffset = ABIWordSize - 20
)

// JWT/JWS parameters for geth's iat window check.
const (
	// GethJwsIatWindowSec is the tight iat window geth accepts for the
	// JWT-authenticated engine api.
	GethJwsIatWindowSec = 60

	// JwsAliveSafetyRatio leaves a safety margin against clock skew or in-flight
	// latency so the client is always recreated before geth's iat window closes.
	JwsAliveSafetyRatio = 0.95

	// JwsAliveWindowSec is the effective alive window after applying the safety ratio.
	JwsAliveWindowSec = int64(GethJwsIatWindowSec * JwsAliveSafetyRatio)
)
