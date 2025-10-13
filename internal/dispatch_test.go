package internal

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/poteto-go/go-alchemy-sdk/types"
)

func TestRequestWithBackoff(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Arrange
		operation := func() (int, error) {
			return 1, nil
		}

		// Act
		result, err := requestWithBackoff(DefaultBackoffConfig, operation)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 1, result)
	})

	t.Run("backoff", func(t *testing.T) {
		// Arrange
		callCount := 0
		operation := func() (int, error) {
			callCount++
			if callCount < 3 {
				return 0, errors.New("test error")
			}
			return 1, nil
		}
		config := BackoffConfig{
			MaxRetries: 3,
		}

		// Act
		result, err := requestWithBackoff(config, operation)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 1, result)
		assert.Equal(t, 3, callCount)
	})

	t.Run("max retries exceeded", func(t *testing.T) {
		// Arrange
		operation := func() (int, error) {
			return 0, errors.New("test error")
		}
		config := BackoffConfig{
			MaxRetries: 3,
		}

		// Act
		_, err := requestWithBackoff(config, operation)

		// Assert
		assert.Error(t, err)
	})
}

func TestRequestWithBackoffTuple(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Arrange
		operation := func() (int, int, error) {
			return 1, 2, nil
		}

		// Act
		result1, result2, err := requestWithBackoffTuple(DefaultBackoffConfig, operation)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 1, result1)
		assert.Equal(t, 2, result2)
	})

	t.Run("backoff", func(t *testing.T) {
		// Arrange
		callCount := 0
		operation := func() (int, int, error) {
			callCount++
			if callCount < 3 {
				return 0, 0, errors.New("test error")
			}
			return 1, 2, nil
		}
		config := BackoffConfig{
			MaxRetries: 3,
		}

		// Act
		result1, result2, err := requestWithBackoffTuple(config, operation)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 1, result1)
		assert.Equal(t, 2, result2)
		assert.Equal(t, 3, callCount)
	})

	t.Run("max retries exceeded", func(t *testing.T) {
		// Arrange
		operation := func() (int, int, error) {
			return 0, 0, errors.New("test error")
		}
		config := BackoffConfig{
			MaxRetries: 3,
		}

		// Act
		_, _, err := requestWithBackoffTuple(config, operation)

		// Assert
		assert.Error(t, err)
	})
}

func TestRequestHttpWithBackoff(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Arrange
		backoffConfig := BackoffConfig{
			MaxRetries: 1,
		}
		requestConfig := types.RequestConfig{}
		mockHandler := func(request types.AlchemyRequest, _ types.RequestConfig, _ []byte) (types.AlchemyResponse, error) {
			return types.AlchemyResponse{Jsonrpc: "2.0"}, nil
		}
		request := types.AlchemyRequest{}
		body := []byte{}

		// Act
		response, err := RequestHttpWithBackoff(backoffConfig, requestConfig, mockHandler, request, body)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "2.0", response.Jsonrpc)
	})
}

func TestGethRequestArgWithBackOff(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Arrange
		backoffConfig := &BackoffConfig{
			MaxRetries: 1,
		}
		mockHandler := func(ctx context.Context, a int) (int, error) {
			return a, nil
		}

		// Act
		result, err := GethRequestArgWithBackOff(backoffConfig, 0, mockHandler, 1)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 1, result)
	})

	t.Run("nil backoffConfig", func(t *testing.T) {
		// Arrange
		mockHandler := func(ctx context.Context, a int) (int, error) {
			return a, nil
		}

		// Act
		result, err := GethRequestArgWithBackOff(nil, 0, mockHandler, 1)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 1, result)
	})
}

func TestGethRequestTwoArgWithBackOff(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Arrange
		backoffConfig := &BackoffConfig{
			MaxRetries: 1,
		}
		mockHandler := func(ctx context.Context, a, b int) (int, error) {
			return a + b, nil
		}

		// Act
		result, err := GethRequestTwoArgWithBackOff(backoffConfig, 0, mockHandler, 1, 2)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 3, result)
	})

	t.Run("nil backoffConfig", func(t *testing.T) {
		// Arrange
		mockHandler := func(ctx context.Context, a, b int) (int, error) {
			return a + b, nil
		}

		// Act
		result, err := GethRequestTwoArgWithBackOff(nil, 0, mockHandler, 1, 2)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 3, result)
	})
}

func TestGethRequestThreeArgWithBackOff(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Arrange
		backoffConfig := &BackoffConfig{
			MaxRetries: 1,
		}
		mockHandler := func(ctx context.Context, a, b, c int) (int, error) {
			return a + b + c, nil
		}

		// Act
		result, err := GethRequestThreeArgWithBackOff(backoffConfig, 0, mockHandler, 1, 2, 3)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 6, result)
	})

	t.Run("nil backoffConfig", func(t *testing.T) {
		// Arrange
		mockHandler := func(ctx context.Context, a, b, c int) (int, error) {
			return a + b + c, nil
		}

		// Act
		result, err := GethRequestThreeArgWithBackOff(nil, 0, mockHandler, 1, 2, 3)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 6, result)
	})
}

func TestGethRequestArgWithBackOffTuple(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Arrange
		backoffConfig := &BackoffConfig{
			MaxRetries: 1,
		}
		mockHandler := func(ctx context.Context, a int) (int, int, error) {
			return a, a, nil
		}

		// Act
		result1, result2, err := GethRequestArgWithBackOffTuple(backoffConfig, 0, mockHandler, 1)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 1, result1)
		assert.Equal(t, 1, result2)
	})

	t.Run("nil backoffConfig", func(t *testing.T) {
		// Arrange
		mockHandler := func(ctx context.Context, a int) (int, int, error) {
			return a, a, nil
		}

		// Act
		result1, result2, err := GethRequestArgWithBackOffTuple(nil, 0, mockHandler, 1)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 1, result1)
		assert.Equal(t, 1, result2)
	})
}

func TestGethRequestWithBackOff(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Arrange
		backoffConfig := &BackoffConfig{
			MaxRetries: 1,
		}
		mockHandler := func(ctx context.Context) (int, error) {
			return 1, nil
		}

		// Act
		result, err := GethRequestWithBackOff(backoffConfig, 0, mockHandler)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 1, result)
	})

	t.Run("nil backoffConfig", func(t *testing.T) {
		// Arrange
		mockHandler := func(ctx context.Context) (int, error) {
			return 1, nil
		}

		// Act
		result, err := GethRequestWithBackOff(nil, 0, mockHandler)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 1, result)
	})
}
