package internal

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
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
		requestConfig := types.RequestConfig{
			Timeout: 10 * time.Second,
		}
		mockHandler := func(request types.AlchemyRequest, _ types.RequestConfig, _ []byte) (types.AlchemyResponse, error) {
			return types.AlchemyResponse{}, nil
		}
		request := types.AlchemyRequest{}
		body := []byte{}

		// Act
		response, err := RequestHttpWithBackoff(backoffConfig, requestConfig, mockHandler, request, body)

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
		requestConfig := types.RequestConfig{
			Timeout: 10 * time.Second,
		}
		callCount := 0
		mockHandler := func(request types.AlchemyRequest, _ types.RequestConfig, _ []byte) (types.AlchemyResponse, error) {
			callCount++
			if callCount < 3 {
				return types.AlchemyResponse{}, errors.New("test error")
			}
			return types.AlchemyResponse{}, nil
		}
		request := types.AlchemyRequest{}
		body := []byte{}

		// Act
		response, err := RequestHttpWithBackoff(backoffConfig, requestConfig, mockHandler, request, body)

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
		requestConfig := types.RequestConfig{
			Timeout: 10 * time.Second,
		}
		mockHandler := func(request types.AlchemyRequest, _ types.RequestConfig, _ []byte) (types.AlchemyResponse, error) {
			return types.AlchemyResponse{}, errors.New("test error")
		}
		request := types.AlchemyRequest{}
		body := []byte{}

		// Act
		_, err := RequestHttpWithBackoff(backoffConfig, requestConfig, mockHandler, request, body)

		// Assert
		assert.Error(t, err)
	})
}

func TestGethRequestMsgWithBackOff(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Arrange
		backoffConfig := &BackoffConfig{
			Mode:           "exponential",
			MaxRetries:     3,
			InitialDelayMs: 10,
			MaxDelayMs:     30,
		}
		mockHandler := func(
			context.Context, ethereum.CallMsg,
		) (int, error) {
			return 1, nil
		}
		msg := ethereum.CallMsg{}

		// Act
		result, err := GethRequestMsgWithBackOff(backoffConfig, 10*time.Second, mockHandler, msg)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, result, 1)
	})

	t.Run("if backoffConfig is nil, use DefaultBackoffConfig", func(t *testing.T) {
		// Arrange
		mockHandler := func(
			context.Context, ethereum.CallMsg,
		) (int, error) {
			return 1, nil
		}
		msg := ethereum.CallMsg{}

		// Act
		result, err := GethRequestMsgWithBackOff(nil, 10*time.Second, mockHandler, msg)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, result, 1)
	})

	t.Run("backoff first 1 is error, & success", func(t *testing.T) {
		// Arrange
		backoffConfig := &BackoffConfig{
			Mode:           "exponential",
			MaxRetries:     3,
			InitialDelayMs: 10,
			MaxDelayMs:     30,
		}
		callCount := 0
		mockHandler := func(
			context.Context, ethereum.CallMsg,
		) (int, error) {
			callCount++
			if callCount < 2 {
				return 0, errors.New("test error")
			}
			return 1, nil
		}
		msg := ethereum.CallMsg{}

		// Act
		result, err := GethRequestMsgWithBackOff(backoffConfig, 10*time.Second, mockHandler, msg)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, result, 1)
	})

	t.Run("max retries exceeded", func(t *testing.T) {
		// Arrange
		backoffConfig := &BackoffConfig{
			Mode:           "exponential",
			MaxRetries:     3,
			InitialDelayMs: 10,
			MaxDelayMs:     30,
		}
		mockHandler := func(
			context.Context, ethereum.CallMsg,
		) (int, error) {
			return 0, errors.New("test error")
		}
		msg := ethereum.CallMsg{}

		// Act
		_, err := GethRequestMsgWithBackOff(backoffConfig, 10*time.Second, mockHandler, msg)

		// Assert
		assert.Error(t, err)
	})
}

func TestGethRequestWithBackOff(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Arrange
		backoffConfig := &BackoffConfig{
			Mode:           "exponential",
			MaxRetries:     3,
			InitialDelayMs: 10,
			MaxDelayMs:     30,
		}
		mockHandler := func(
			context.Context,
		) (int, error) {
			return 1, nil
		}

		// Act
		result, err := GethRequestWithBackOff(backoffConfig, 10*time.Second, mockHandler)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, result, 1)
	})

	t.Run("if backoffConfig is nil, use DefaultBackoffConfig", func(t *testing.T) {
		// Arrange
		mockHandler := func(
			context.Context,
		) (int, error) {
			return 1, nil
		}

		// Act
		result, err := GethRequestWithBackOff(nil, 10*time.Second, mockHandler)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, result, 1)
	})

	t.Run("backoff first 1 is error, & success", func(t *testing.T) {
		// Arrange
		backoffConfig := &BackoffConfig{
			Mode:           "exponential",
			MaxRetries:     3,
			InitialDelayMs: 10,
			MaxDelayMs:     30,
		}
		callCount := 0
		mockHandler := func(
			context.Context,
		) (int, error) {
			callCount++
			if callCount < 2 {
				return 0, errors.New("test error")
			}
			return 1, nil
		}
		// Act
		result, err := GethRequestWithBackOff(backoffConfig, 10*time.Second, mockHandler)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, result, 1)
	})

	t.Run("max retries exceeded", func(t *testing.T) {
		// Arrange
		backoffConfig := &BackoffConfig{
			Mode:           "exponential",
			MaxRetries:     3,
			InitialDelayMs: 10,
			MaxDelayMs:     30,
		}
		mockHandler := func(
			context.Context,
		) (int, error) {
			return 0, errors.New("test error")
		}

		// Act
		_, err := GethRequestWithBackOff(backoffConfig, 10*time.Second, mockHandler)

		// Assert
		assert.Error(t, err)
	})
}
