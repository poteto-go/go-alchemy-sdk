package internal

import (
	"context"
	"time"

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

func GethRequestTwoArgWithBackOff[T any, A any, B any](
	backoffConfig *BackoffConfig,
	timeout time.Duration,
	handler func(
		context.Context, A, B,
	) (T, error),
	arg1 A,
	arg2 B,
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

		result, err := handler(ctx, arg1, arg2)
		if err == nil {
			return result, nil
		}

		lastHttpError = err
		if err := backoffManager.Backoff(); err != nil {
			return result, lastHttpError
		}
	}
}

func GethRequestArgWithBackOffTuple[T any, A any, O any](
	backoffConfig *BackoffConfig,
	timeout time.Duration,
	handler func(
		context.Context, A,
	) (T, O, error),
	arg A,
) (T, O, error) {
	var lastHttpError error
	if backoffConfig == nil {
		backoffConfig = &DefaultBackoffConfig
	}

	backoffManager := NewBackoffManager(*backoffConfig)
	for {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		result, other, err := handler(ctx, arg)
		if err == nil {
			return result, other, nil
		}

		lastHttpError = err
		if err := backoffManager.Backoff(); err != nil {
			return result, other, lastHttpError
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
