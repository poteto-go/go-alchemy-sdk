package internal

import "github.com/poteto-go/go-alchemy-sdk/types"

func RequestHttpWithBackoff(
	backoffConfig BackoffConfig,
	requestConfig types.RequestConfig,
	handler types.AlchemyFetchHandler,
	request types.AlchemyRequest,
) (types.AlchemyResponse, error) {
	var lastHttpError error

	backoffManager := NewBackoffManager(backoffConfig)
	for {
		response, err := handler(request, requestConfig)
		if err == nil {
			return response, nil
		}

		lastHttpError = err

		if err := backoffManager.Backoff(); err != nil {
			return types.AlchemyResponse{}, lastHttpError
		}
	}
}
