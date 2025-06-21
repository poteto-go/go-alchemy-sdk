package core

import "errors"

var (
	ErrInvalidHexString          = errors.New("invalid hex string")
	ErrFailedToMarshalParameter  = errors.New("failed to marshal parameter")
	ErrFailedToCreateRequest     = errors.New("failed to create request")
	ErrFailedToConnect           = errors.New("failed to connect")
	ErrFailedToUnmarshalResponse = errors.New("failed to unmarshal response")
	ErrBatcherNotRunning         = errors.New("batcher not running")
	ErrRequestTimeoutTooShort    = errors.New("request timeout should be longer than batch timeout")
	ErrRequestTimeout            = errors.New("request timeout")
	ErrOverMaxRetries            = errors.New("over max retries")
)
