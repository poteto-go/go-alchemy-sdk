package internal

import (
	"context"
	"errors"
	"time"

	"github.com/poteto-go/go-alchemy-sdk/types"
)

type BatchSubscriber struct {
	Responses   []types.AlchemyResponse
	JoinedError error
}

type IRequestBatcher interface {
	// Subscribe wait for responses
	//
	// reset BatchSubscriber
	Subscribe(<-chan struct{}) BatchSubscriber

	// Put BatchSubscriber
	//
	// - reset BatchSubscriber
	put() BatchSubscriber

	// Record request
	//
	// func main() {
	//   batcher := NewRequestBatcher(
	//     context.Background(),
	//     BatcherConfig{
	//       MaxBatchSize: 100,
	//       MaxBatchTime: time.Millisecond * 10,
	//       Fetch:        utils.AlchemyBatchFetch,
	//     },
	//   )
	//
	//   done := make (chan struct{}, 1)
	//   requestChan := make(chan types.AlchemyRequest, 1)
	//   go func() {
	//     batcher.RecordRequest(requestChan)
	//     done <- struct{}{}
	//   } ()
	//
	//   for i := 0; i<100; i++ {
	//     body := types.AlchemyRequestBody{
	//       Jsonrpc: "2.0",
	//       Method:  "method",
	//       Params:  []string{"param1", "param2"},
	//       Id:      i,
	//     }
	//     req, _ := http.NewRequest("POST", targetUrl, nil)
	//
	//     request := types.AlchemyRequest{
	//       Request: req,
	//       Body:    body,
	//     }
	//     requestChan <- request
	//   }
	//
	//   subscribers := batcher.Subscribe(done)
	// }
	RecordRequest(<-chan types.AlchemyRequest)

	send()
	release()
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
	Subscribers BatchSubscriber
}

func NewRequestBatcher(
	ctx context.Context,
	config BatcherConfig,
) IRequestBatcher {
	return &RequestBatcher{
		config:   config,
		IsRun:    true,
		requests: []types.AlchemyRequest{},
		Subscribers: BatchSubscriber{
			Responses:   []types.AlchemyResponse{},
			JoinedError: nil,
		},
	}
}

func (b *RequestBatcher) Subscribe(done <-chan struct{}) BatchSubscriber {
	for {
		select {
		case <-done:
			return b.put()
		default:
			continue
		}
	}
}

func (b *RequestBatcher) put() BatchSubscriber {
	defer func() {
		b.Subscribers = BatchSubscriber{
			Responses:   []types.AlchemyResponse{},
			JoinedError: nil,
		}
	}()

	return b.Subscribers
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
				b.release()
				return
			}
		case <-time.After(b.config.MaxBatchTime):
			if len(b.requests) == 0 {
				return
			}

			b.send()
			b.release()
			return
		}
	}
}

func (b *RequestBatcher) send() {
	batchedRes, err := b.config.Fetch(b.requests)
	for i := len(b.requests) - 1; i >= 0; i-- {
		b.Subscribers.Responses = append(b.Subscribers.Responses, batchedRes...)
		if err != nil {
			b.Subscribers.JoinedError = errors.Join(b.Subscribers.JoinedError, err)
		}
	}
}

func (b *RequestBatcher) release() {
	b.requests = []types.AlchemyRequest{}
}
