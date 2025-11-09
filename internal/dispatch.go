package internal

import (
	"context"
	"time"

	"github.com/poteto-go/go-alchemy-sdk/types"
)

func requestWithBackoffError(
	backoffConfig types.BackoffConfig,
	operation func() error,
) error {
	var lastHttpError error
	backoffManager := NewBackoffManager(backoffConfig)
	for {
		err := operation()
		if err == nil {
			return nil
		}

		lastHttpError = err
		if err := backoffManager.Backoff(); err != nil {
			return lastHttpError
		}
	}
}

func requestWithBackoff[T any](
	backoffConfig types.BackoffConfig,
	operation func() (T, error),
) (T, error) {
	var lastHttpError error
	backoffManager := NewBackoffManager(backoffConfig)
	for {
		result, err := operation()
		if err == nil {
			return result, nil
		}

		lastHttpError = err
		if err := backoffManager.Backoff(); err != nil {
			var zero T
			return zero, lastHttpError
		}
	}
}

func requestWithBackoffTuple[T any, O any](
	backoffConfig types.BackoffConfig,
	operation func() (T, O, error),
) (T, O, error) {
	var lastHttpError error
	backoffManager := NewBackoffManager(backoffConfig)
	for {
		result, other, err := operation()
		if err == nil {
			return result, other, nil
		}

		lastHttpError = err
		if err := backoffManager.Backoff(); err != nil {
			var zeroT T
			var zeroO O
			return zeroT, zeroO, lastHttpError
		}
	}
}

func RequestHttpWithBackoff(
	backoffConfig types.BackoffConfig,
	requestConfig types.RequestConfig,
	handler types.AlchemyFetchHandler,
	request types.AlchemyRequest,
	body []byte,
) (types.AlchemyResponse, error) {
	operation := func() (types.AlchemyResponse, error) {
		return handler(request, requestConfig, body)
	}
	return requestWithBackoff(backoffConfig, operation)
}

func GethRequestArgWithBackOff[T any, A any](
	backoffConfig *types.BackoffConfig,
	timeout time.Duration,
	handler func(
		context.Context, A,
	) (T, error),
	arg A,
) (T, error) {
	if backoffConfig == nil {
		backoffConfig = &types.DefaultBackoffConfig
	}
	operation := func() (T, error) {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		return handler(ctx, arg)
	}
	return requestWithBackoff(*backoffConfig, operation)
}

func GethRequestTwoArgWithBackOff[T any, A any, B any](
	backoffConfig *types.BackoffConfig,
	timeout time.Duration,
	handler func(
		context.Context, A, B,
	) (T, error),
	arg1 A,
	arg2 B,
) (T, error) {
	if backoffConfig == nil {
		backoffConfig = &types.DefaultBackoffConfig
	}
	operation := func() (T, error) {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		return handler(ctx, arg1, arg2)
	}
	return requestWithBackoff(*backoffConfig, operation)
}

func GethRequestThreeArgWithBackOff[T any, A any, B any, C any](
	backoffConfig *types.BackoffConfig,
	timeout time.Duration,
	handler func(
		context.Context, A, B, C,
	) (T, error),
	arg1 A,
	arg2 B,
	arg3 C,
) (T, error) {
	if backoffConfig == nil {
		backoffConfig = &types.DefaultBackoffConfig
	}
	operation := func() (T, error) {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		return handler(ctx, arg1, arg2, arg3)
	}
	return requestWithBackoff(*backoffConfig, operation)
}

func GethRequestArgWithBackOffTuple[T any, A any, O any](
	backoffConfig *types.BackoffConfig,
	timeout time.Duration,
	handler func(
		context.Context, A,
	) (T, O, error),
	arg A,
) (T, O, error) {
	if backoffConfig == nil {
		backoffConfig = &types.DefaultBackoffConfig
	}
	operation := func() (T, O, error) {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		return handler(ctx, arg)
	}
	return requestWithBackoffTuple(*backoffConfig, operation)
}

func GethRequestWithBackOff[T any](
	backoffConfig *types.BackoffConfig,
	timeout time.Duration,
	handler func(
		context.Context,
	) (T, error),
) (T, error) {
	if backoffConfig == nil {
		backoffConfig = &types.DefaultBackoffConfig
	}
	operation := func() (T, error) {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		return handler(ctx)
	}
	return requestWithBackoff(*backoffConfig, operation)
}

func GethRequestSingleErrorWithBackOff[A any](
	backoffConfig *types.BackoffConfig,
	timeout time.Duration,
	handler func(
		context.Context, A,
	) error,
	arg A,
) error {
	if backoffConfig == nil {
		backoffConfig = &types.DefaultBackoffConfig
	}
	operation := func() error {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		return handler(ctx, arg)
	}
	return requestWithBackoffError(*backoffConfig, operation)
}
