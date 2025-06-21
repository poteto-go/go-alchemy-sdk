package internal

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/poteto-go/go-alchemy-sdk/types"
)

func TestRequestHttpWithBackoff(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Arrange
		backoffConfig := BackoffConfig{
			Mode:           "exponential",
			MaxRetries:     3,
			InitialDelayMs: 10,
			MaxDelayMs:     30,
		}
		mockHandler := func(request types.AlchemyRequest) (types.AlchemyResponse, error) {
			return types.AlchemyResponse{}, nil
		}
		request := types.AlchemyRequest{}

		// Act
		response, err := RequestHttpWithBackoff(backoffConfig, mockHandler, request)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, types.AlchemyResponse{}, response)
	})

	t.Run("backoff", func(t *testing.T) {
		// Arrange
		backoffConfig := BackoffConfig{
			Mode:           "exponential",
			MaxRetries:     3,
			InitialDelayMs: 10,
			MaxDelayMs:     30,
		}
		callCount := 0
		mockHandler := func(request types.AlchemyRequest) (types.AlchemyResponse, error) {
			callCount++
			if callCount < 3 {
				return types.AlchemyResponse{}, errors.New("test error")
			}
			return types.AlchemyResponse{}, nil
		}
		request := types.AlchemyRequest{}

		// Act
		response, err := RequestHttpWithBackoff(backoffConfig, mockHandler, request)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, types.AlchemyResponse{}, response)
		assert.Equal(t, 3, callCount)
	})

	t.Run("max retries exceeded", func(t *testing.T) {
		// Arrange
		backoffConfig := BackoffConfig{
			Mode:           "exponential",
			MaxRetries:     3,
			InitialDelayMs: 10,
			MaxDelayMs:     30,
		}
		mockHandler := func(request types.AlchemyRequest) (types.AlchemyResponse, error) {
			return types.AlchemyResponse{}, errors.New("test error")
		}
		request := types.AlchemyRequest{}

		// Act
		_, err := RequestHttpWithBackoff(backoffConfig, mockHandler, request)

		// Assert
		assert.Error(t, err)
	})
}
