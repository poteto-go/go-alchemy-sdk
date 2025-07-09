package internal

import (
	"context"
	"sync"
	"time"

	"github.com/poteto-go/go-alchemy-sdk/types"
)

type QueuedRequest struct {
	Request  types.AlchemyRequest[string]
	Response chan types.AlchemyResponse
}

type RequestBatcher struct {
	config        BatcherConfig
	requestConfig types.RequestConfig
	requestQueue  chan QueuedRequest
	mutex         sync.Mutex
}

type BatcherConfig struct {
	MaxBatchSize int
	MaxBatchTime time.Duration
	Fetch        types.BatchAlchemyFetchHandler[string]
}

func NewRequestBatcher(
	ctx context.Context,
	config BatcherConfig,
	requestConfig types.RequestConfig,
) *RequestBatcher {
	batcher := &RequestBatcher{
		config:        config,
		requestConfig: requestConfig,
		requestQueue:  make(chan QueuedRequest, config.MaxBatchSize),
	}
	go batcher.processQueue(ctx)
	return batcher
}

func (b *RequestBatcher) QueueRequest(ctx context.Context, request types.AlchemyRequest[string]) (types.AlchemyResponse, error) {
	responseChan := make(chan types.AlchemyResponse, 1)
	select {
	case <-ctx.Done():
		return types.AlchemyResponse{}, ctx.Err()
	case b.requestQueue <- QueuedRequest{
		Request:  request,
		Response: responseChan,
	}:
	}

	select {
	case <-ctx.Done():
		return types.AlchemyResponse{}, ctx.Err()
	case response := <-responseChan:
		if response.Error != nil {
			return types.AlchemyResponse{}, response.Error
		}
		return response, nil
	}
}

func (b *RequestBatcher) processQueue(ctx context.Context) {
	batch := make([]QueuedRequest, 0, b.config.MaxBatchSize)
	timer := time.NewTimer(b.config.MaxBatchTime)

	for {
		select {
		case <-ctx.Done():
			b.flushWithError(batch, ctx.Err())
			return
		case req := <-b.requestQueue:
			batch = append(batch, req)
			if len(batch) >= b.config.MaxBatchSize {
				b.flush(batch)
				batch = make([]QueuedRequest, 0, b.config.MaxBatchSize)
				timer.Reset(b.config.MaxBatchTime)
			}
		case <-timer.C:
			if len(batch) > 0 {
				b.flush(batch)
				batch = make([]QueuedRequest, 0, b.config.MaxBatchSize)
			}
			timer.Reset(b.config.MaxBatchTime)
		}
	}
}

func (b *RequestBatcher) flush(batch []QueuedRequest) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	requests := make([]types.AlchemyRequest[string], len(batch))
	for i, req := range batch {
		requests[i] = req.Request
	}

	responses, err := b.config.Fetch(requests, b.requestConfig)
	if err != nil {
		for _, req := range batch {
			req.Response <- types.AlchemyResponse{Error: err}
		}
		return
	}

	responseMap := make(map[int]types.AlchemyResponse)
	for _, res := range responses {
		responseMap[res.Id] = res
	}

	for _, req := range batch {
		if res, ok := responseMap[req.Request.Body.Id]; ok {
			req.Response <- res
		} else {
			req.Response <- types.AlchemyResponse{Error: types.ErrNoResultFound}
		}
	}
}

func (b *RequestBatcher) flushWithError(batch []QueuedRequest, err error) {
	for _, req := range batch {
		req.Response <- types.AlchemyResponse{Error: err}
	}
}
