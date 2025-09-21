package internal

import (
	"context"

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

func GethRequestWithBackOff[T any](
	BackoffConfig BackoffConfig,
	handler func(
		context.Context, ethereum.CallMsg,
	) (T, error),
	ctx context.Context,
	msg ethereum.CallMsg,
) (T, error) {
	var lastHttpError error

	backoffManager := NewBackoffManager(BackoffConfig)
	for {
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
