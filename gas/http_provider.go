package gas

import (
	"context"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/poteto-go/go-alchemy-sdk/constant"
	"github.com/poteto-go/go-alchemy-sdk/internal"
	"github.com/poteto-go/go-alchemy-sdk/types"
	"github.com/poteto-go/go-alchemy-sdk/utils"
)

type AlchemyProvider struct {
	config  AlchemyConfig
	id      atomic.Int64
	batcher *internal.RequestBatcher
	eth     types.EtherApi
	client  *http.Client // shared across all Send calls
}

func NewAlchemyProvider(config AlchemyConfig) types.IAlchemyProvider {
	provider := &AlchemyProvider{
		config: config,
		client: utils.NewSharedHTTPClient(config.maxResponseBytes),
	}
	provider.id.Store(1)

	if config.maxRetries > 0 {
		sharedClient := provider.client
		rc := types.RequestConfig{
			Timeout:          config.requestTimeout,
			MaxResponseBytes: config.maxResponseBytes,
		}
		provider.batcher = internal.NewRequestBatcher(
			context.Background(),
			internal.BatcherConfig{
				MaxBatchSize: 100,
				MaxBatchTime: time.Millisecond * 10,
				Fetch: func(reqs []types.AlchemyRequest, cfg types.RequestConfig, bodies [][]byte) ([]types.AlchemyResponse, error) {
					return utils.AlchemyBatchFetch(sharedClient, reqs, cfg, bodies)
				},
			},
			rc,
		)
	}

	return provider
}

func (provider *AlchemyProvider) SetEth(eth types.EtherApi) {
	provider.eth = eth
}

func (provider *AlchemyProvider) Eth() types.EtherApi {
	return provider.eth
}

func (provider *AlchemyProvider) CustomHeaders() []http.Header {
	return provider.config.customHeaders
}

/* Send raw transaction */
func (provider *AlchemyProvider) Send(method string, params types.RequestArgs) (any, error) {
	// fetch-and-add: take the current id for this request, then advance the counter atomically.
	id := provider.id.Add(1) - 1
	body, err := utils.CreateRequestBodyToBytes(int(id), method, params)
	if err != nil {
		return nil, err
	}
	return send(provider, body)
}

func send(provider *AlchemyProvider, body []byte) (any, error) {
	req, err := generateAlchemyRequest(provider.config)
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
			return fmt.Sprintf("%v", response.Result), nil
		}
	*/

	sharedClient := provider.client
	rc := types.RequestConfig{
		Timeout:          provider.config.requestTimeout,
		MaxResponseBytes: provider.config.maxResponseBytes,
	}
	response, err := internal.RequestHttpWithBackoff(
		*provider.config.backoffConfig,
		rc,
		func(req types.AlchemyRequest, cfg types.RequestConfig, b []byte) (types.AlchemyResponse, error) {
			return utils.AlchemyFetch(sharedClient, req, cfg, b)
		},
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

	return result, nil
}

func generateAlchemyRequest(config AlchemyConfig) (*http.Request, error) {
	req, err := http.NewRequest("POST", config.GetUrl(), nil)
	if err != nil {
		return &http.Request{}, constant.ErrFailedToCreateRequest
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Alchemy-Ethers-Sdk-Method", "send")

	for _, header := range config.customHeaders {
		for key, values := range header {
			for _, value := range values {
				req.Header.Set(key, value)
			}
		}
	}

	return req, nil
}
