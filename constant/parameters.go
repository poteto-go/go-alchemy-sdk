package constant

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
