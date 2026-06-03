package constant

// Ethereum ECDSA signature parameters.
const (
	// ECDSALegacyVOffset is added to the ECDSA recovery id (0 or 1) returned
	// by crypto.Sign to produce the Ethereum v value (27 or 28).
	ECDSALegacyVOffset = 27
)
