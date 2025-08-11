package alchemy

import (
	"context"
	"net/http"
	"time"

	"github.com/poteto-go/go-alchemy-sdk/constant"
	"github.com/poteto-go/go-alchemy-sdk/internal"
	"github.com/poteto-go/go-alchemy-sdk/types"
	"github.com/poteto-go/go-alchemy-sdk/utils"
)

type AlchemyProvider struct {
	config  AlchemyConfig
	id      int
	batcher *internal.RequestBatcher
}

func NewAlchemyProvider(config AlchemyConfig) types.IAlchemyProvider {
	provider := &AlchemyProvider{
		config: config,
		id:     1,
	}

	if config.maxRetries > 0 {
		provider.batcher = internal.NewRequestBatcher(
			context.Background(),
			internal.BatcherConfig{
				MaxBatchSize: 100,
				MaxBatchTime: time.Millisecond * 10,
				Fetch:        utils.AlchemyBatchFetch,
			},
			types.RequestConfig{
				Timeout: config.requestTimeout,
			},
		)
	}

	return provider
}

func (provider *AlchemyProvider) Send(method string, params types.RequestArgs) (any, error) {
	body, err := utils.CreateRequestBodyToBytes(provider.id, method, params)
	if err != nil {
		return nil, err
	}
	return send(provider, body)
}

func send(provider *AlchemyProvider, body []byte) (any, error) {
	req, err := generateAlchemyRequest(provider.config.GetUrl())
	if err != nil {
		return nil, err
	}

	request := types.AlchemyRequest{
		Request: req,
	}

	// TODO: not support generics in batch request for now.
	/*
		if provider.batcher != nil {
			response, err := provider.batcher.QueueRequest(context.Background(), request)
			if err != nil {
				return "", err
			}
			provider.id++
			return fmt.Sprintf("%v", response.Result), nil
		}
	*/

	response, err := internal.RequestHttpWithBackoff(
		*provider.config.backoffConfig,
		types.RequestConfig{
			Timeout: provider.config.requestTimeout,
		},
		utils.AlchemyFetch,
		request,
		body,
	)
	if err != nil {
		return nil, err
	}

	result := response.Result
	if result == nil {
		return nil, constant.ErrResultIsNil
	}

	provider.id++

	return result, nil
}

func generateAlchemyRequest(url string) (*http.Request, error) {
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return &http.Request{}, constant.ErrFailedToCreateRequest
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Alchemy-Ethers-Sdk-Method", "send")

	return req, nil
}
