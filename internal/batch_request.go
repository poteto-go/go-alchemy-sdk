package internal

import (
	"context"
	"sync"
	"time"

	"github.com/goccy/go-json"
	"github.com/poteto-go/go-alchemy-sdk/types"
)

type QueuedRequest struct {
	Request  types.AlchemyRequest
	Body     []byte
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
	Fetch        types.BatchAlchemyFetchHandler
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

func (b *RequestBatcher) QueueRequest(ctx context.Context, request types.AlchemyRequest, body []byte) (types.AlchemyResponse, error) {
	responseChan := make(chan types.AlchemyResponse, 1)
	select {
	case <-ctx.Done():
		return types.AlchemyResponse{}, ctx.Err()
	case b.requestQueue <- QueuedRequest{
		Request:  request,
		Body:     body,
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

	requests := make([]types.AlchemyRequest, len(batch))
	bodies := make([][]byte, len(batch))
	for i, req := range batch {
		requests[i] = req.Request
		bodies[i] = req.Body
	}

	responses, err := b.config.Fetch(requests, b.requestConfig, bodies)
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
		var requestBody types.AlchemyRequestBody[string]
		json.Unmarshal(req.Body, &requestBody)
		if res, ok := responseMap[requestBody.Id]; ok {
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
