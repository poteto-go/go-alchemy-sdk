package internal

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/poteto-go/go-alchemy-sdk/types"
)

func RequestHttpWithBackoff(
	backoffConfig BackoffConfig,
	requestConfig types.RequestConfig,
	handler types.AlchemyFetchHandler,
	request types.AlchemyRequest,
	body []byte,
) (types.AlchemyResponse, error) {
	var lastHttpError error

	backoffManager := NewBackoffManager(backoffConfig)
	for {
		response, err := handler(request, requestConfig, body)
		if err == nil {
			return response, nil
		}

		lastHttpError = err

		if err := backoffManager.Backoff(); err != nil {
			return types.AlchemyResponse{}, lastHttpError
		}
	}
}

func GethRequestMsgWithBackOff[T any](
	backoffConfig *BackoffConfig,
	timeout time.Duration,
	handler func(
		context.Context, ethereum.CallMsg,
	) (T, error),
	msg ethereum.CallMsg,
) (T, error) {
	var lastHttpError error
	if backoffConfig == nil {
		backoffConfig = &DefaultBackoffConfig
	}

	backoffManager := NewBackoffManager(*backoffConfig)
	for {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		result, err := handler(ctx, msg)
		if err == nil {
			return result, nil
		}

		lastHttpError = err
		if err := backoffManager.Backoff(); err != nil {
			return result, lastHttpError
		}
	}
}

func GethRequestArgWithBackOff[T any, A any](
	backoffConfig *BackoffConfig,
	timeout time.Duration,
	handler func(
		context.Context, A,
	) (T, error),
	arg A,
) (T, error) {
	var lastHttpError error
	if backoffConfig == nil {
		backoffConfig = &DefaultBackoffConfig
	}

	backoffManager := NewBackoffManager(*backoffConfig)
	for {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		result, err := handler(ctx, arg)
		if err == nil {
			return result, nil
		}

		lastHttpError = err
		if err := backoffManager.Backoff(); err != nil {
			return result, lastHttpError
		}
	}
}

func GethRequestWithBackOff[T any](
	backoffConfig *BackoffConfig,
	timeout time.Duration,
	handler func(
		context.Context,
	) (T, error),
) (T, error) {
	var lastHttpError error
	if backoffConfig == nil {
		backoffConfig = &DefaultBackoffConfig
	}

	backoffManager := NewBackoffManager(*backoffConfig)
	for {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		result, err := handler(ctx)
		if err == nil {
			return result, nil
		}

		lastHttpError = err
		if err := backoffManager.Backoff(); err != nil {
			return result, lastHttpError
		}
	}
}
