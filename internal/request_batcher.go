package internal

import (
	"context"
	"time"

	"github.com/poteto-go/go-alchemy-sdk/types"
)

type batchSubscriber struct {
	responses []types.AlchemyResponse
	errors    []error
}

type IRequestBatcher interface {
	RecordRequest(newRequest <-chan types.AlchemyRequest)
	send()
	ReleaseBatch()
}

type BatcherConfig struct {
	MaxBatchSize int
	MaxBatchTime time.Duration
	Fetch        func([]types.AlchemyRequest) ([]types.AlchemyResponse, error)
}

type RequestBatcher struct {
	config      BatcherConfig
	IsRun       bool
	requests    []types.AlchemyRequest
	subscribers batchSubscriber
}

func NewRequestBatcher(
	ctx context.Context,
	config BatcherConfig,
) IRequestBatcher {
	batcher := &RequestBatcher{
		config: config,
		IsRun:  true,
	}

	return batcher
}

func (b *RequestBatcher) RecordRequest(newRequest <-chan types.AlchemyRequest) {
	if !b.IsRun {
		return
	}

	for {
		select {
		case request := <-newRequest:
			b.requests = append(b.requests, request)
			if len(b.requests) >= b.config.MaxBatchSize {
				b.send()
				b.ReleaseBatch()
				return
			}
		case <-time.After(b.config.MaxBatchTime):
			if len(b.requests) == 0 {
				continue
			}
			b.send()
			b.ReleaseBatch()
			return
		}
	}
}

func (b *RequestBatcher) send() {
	batchedRes, err := b.config.Fetch(b.requests)
	for i := len(b.requests) - 1; i >= 0; i-- {
		b.subscribers.responses = append(b.subscribers.responses, batchedRes...)
		b.subscribers.errors = append(b.subscribers.errors, err)
	}
}

func (b *RequestBatcher) ReleaseBatch() {
	b.requests = []types.AlchemyRequest{}
	b.subscribers = batchSubscriber{}
}
