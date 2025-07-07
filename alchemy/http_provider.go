package alchemy

import (
	"context"
	"fmt"
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
	return provider.send(method, params...)
}

func (provider *AlchemyProvider) send(method string, params ...string) (any, error) {
	if len(params) == 0 {
		params = []string{}
	}

	body := types.AlchemyRequestBody[string]{
		Jsonrpc: "2.0",
		Method:  method,
		Params:  params,
		Id:      provider.id,
	}

	req, err := http.NewRequest("POST", provider.config.GetUrl(), nil)
	if err != nil {
		return "", core.ErrFailedToCreateRequest
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Alchemy-Ethers-Sdk-Method", "send")

	request := types.AlchemyRequest[string]{
		Body:    body,
		Request: req,
	}

	if provider.batcher != nil {
		response, err := provider.batcher.QueueRequest(context.Background(), request)
		if err != nil {
			return "", err
		}
		provider.id++
		return fmt.Sprintf("%v", response.Result), nil
	}

	response, err := internal.RequestHttpWithBackoff(
		*provider.config.backoffConfig,
		types.RequestConfig{
			Timeout: provider.config.requestTimeout,
		},
		utils.AlchemyFetch[string],
		request,
	)
	if err != nil {
		return "", err
	}

	provider.id++

	return response.Result, nil
}

func (provider *AlchemyProvider) SendTransaction(method string, params ...types.TransactionRequest) (any, error) {
	return provider.sendTransaction(method, params...)
}

func (provider *AlchemyProvider) sendTransaction(method string, params ...types.TransactionRequest) (any, error) {
	if len(params) == 0 {
		params = []types.TransactionRequest{}
	}

	body := types.AlchemyRequestBody[types.TransactionRequest]{
		Jsonrpc: "2.0",
		Method:  method,
		Params:  params,
		Id:      provider.id,
	}

	req, err := http.NewRequest("POST", provider.config.GetUrl(), nil)
	if err != nil {
		return "", core.ErrFailedToCreateRequest
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Alchemy-Ethers-Sdk-Method", "send")

	request := types.AlchemyRequest[types.TransactionRequest]{
		Body:    body,
		Request: req,
	}

	response, err := internal.RequestHttpWithBackoff(
		*provider.config.backoffConfig,
		types.RequestConfig{
			Timeout: provider.config.requestTimeout,
		},
		utils.AlchemyFetch[types.TransactionRequest],
		request,
	)
	if err != nil {
		return "", err
	}

	provider.id++

	return response.Result, nil
}
