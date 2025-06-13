package core

import "errors"

var (
	ErrInvalidHexString          = errors.New("invalid hex string")
	ErrFailedToMarshalParameter  = errors.New("failed to marshal parameter")
	ErrFailedToCreateRequest     = errors.New("failed to create request")
	ErrFailedToConnect           = errors.New("failed to connect")
	ErrFailedToUnmarshalResponse = errors.New("failed to unmarshal response")
)
