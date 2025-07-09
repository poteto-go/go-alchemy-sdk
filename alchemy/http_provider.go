package alchemy

import (
	"context"
	"net/http"
	"time"

	"github.com/poteto-go/go-alchemy-sdk/core"
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
				Fetch:        utils.AlchemyBatchFetch[string],
			},
			types.RequestConfig{
				Timeout: config.requestTimeout,
			},
		)
	}

	return provider
}

func (provider *AlchemyProvider) Send(method string, params ...string) (any, error) {
	return send(provider, method, params...)
}

func send[T string | types.TransactionRequest](provider *AlchemyProvider, method string, params ...T) (any, error) {
	if len(params) == 0 {
		params = []T{}
	}

	body := types.AlchemyRequestBody[T]{
		Jsonrpc: "2.0",
		Method:  method,
		Params:  params,
		Id:      provider.id,
	}

	req, err := generateAlchemyRequest(provider.config.GetUrl())
	if err != nil {
		return "", err
	}

	request := types.AlchemyRequest[T]{
		Body:    body,
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
		utils.AlchemyFetch[T],
		request,
	)
	if err != nil {
		return "", err
	}

	provider.id++

	return response.Result, nil
}

func (provider *AlchemyProvider) SendTransaction(method string, params ...types.TransactionRequest) (any, error) {
	return send(provider, method, params...)
}

func generateAlchemyRequest(url string) (*http.Request, error) {
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return &http.Request{}, core.ErrFailedToCreateRequest
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Alchemy-Ethers-Sdk-Method", "send")

	return req, nil
}
