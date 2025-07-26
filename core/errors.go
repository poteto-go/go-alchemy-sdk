package core

import "errors"

var (
	ErrInvalidHexString             = errors.New("invalid hex string")
	ErrFailedToMarshalParameter     = errors.New("failed to marshal parameter")
	ErrFailedToCreateRequest        = errors.New("failed to create request")
	ErrFailedToConnect              = errors.New("failed to connect")
	ErrFailedToUnmarshalResponse    = errors.New("failed to unmarshal response")
	ErrFailedToUnmarshalTransaction = errors.New("failed to unmarshal transaction")
	ErrFailedToMapTransaction       = errors.New("failed to map transaction")
	ErrFailedToMapTokenResponse     = errors.New("failed to map token response")
	ErrBatcherNotRunning            = errors.New("batcher not running")
	ErrRequestTimeoutTooShort       = errors.New("request timeout should be longer than batch timeout")
	ErrRequestTimeout               = errors.New("request timeout")
	ErrOverMaxRetries               = errors.New("over max retries")
	ErrInvalidBlockTag              = errors.New("invalid block tag")
	ErrFailedToTransformBlockNumber = errors.New("failed to transform block number")
	ErrFailedToTransformType        = errors.New("failed to transform type")
	ErrFailedToTransformNonce       = errors.New("failed to transform nonce")
	ErrFailedToTransformGasPrice    = errors.New("failed to transform gas price")
	ErrFailedToTransformGasLimit    = errors.New("failed to transform gas limit")
	ErrFailedToTransformValue       = errors.New("failed to transform value")
	ErrFailedToTransformChainId     = errors.New("failed to transform chain id")
	ErrFailedToTransformV           = errors.New("failed to transform v")
	ErrResultIsNil                  = errors.New("result is nil")
)
