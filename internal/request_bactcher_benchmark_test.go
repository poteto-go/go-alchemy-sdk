package internal

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/poteto-go/go-alchemy-sdk/types"
	"github.com/poteto-go/go-alchemy-sdk/utils"
	"github.com/stretchr/testify/assert"
)

func TestNewRequestBatcher(t *testing.T) {
	// Arrange
	ctx := context.Background()

	// Act
	batcher := NewRequestBatcher(
		ctx,
		BatcherConfig{
			MaxBatchSize: 10,
			MaxBatchTime: time.Millisecond * 10,
			Fetch:        utils.AlchemyBatchFetch,
		},
		types.RequestConfig{
			Timeout: time.Second * 10,
		},
	).(*RequestBatcher)

	// Assert
	assert.True(t, batcher.IsRun)
}

func newBenchmarkBatcher() *RequestBatcher {
	ctx := context.Background()
	return NewRequestBatcher(
		ctx,
		BatcherConfig{
			MaxBatchSize: 100,
			MaxBatchTime: time.Millisecond * 10,
			Fetch:        utils.AlchemyBatchFetch,
		},
		types.RequestConfig{
			Timeout: time.Second * 10,
		},
	).(*RequestBatcher)
}

func BenchmarkRequest_Batch(b *testing.B) {
	httpmock.Activate(b)
	defer httpmock.DeactivateAndReset()

	// Arrange
	targetUrl := "example.com"
	body := types.AlchemyRequestBody{
		Jsonrpc: "2.0",
		Method:  "method",
		Params:  []string{"param1", "param2"},
		Id:      1,
	}
	req, _ := http.NewRequest("POST", targetUrl, nil)
	mockResult := []types.AlchemyResponse{
		{
			Jsonrpc: "2.0",
			Id:      1,
			Result:  "0x1234",
		},
	}
	resultJson, _ := json.Marshal(mockResult)

	// Mock
	httpmock.RegisterResponder(
		"POST",
		targetUrl,
		httpmock.NewStringResponder(200, string(resultJson)),
	)

	b.ResetTimer()
	for b.Loop() {
		batcher := newBenchmarkBatcher()

		requestEvent := make(chan types.AlchemyRequest, 1)
		done := make(chan struct{}, 1)
		go func() {
			batcher.RecordRequest(requestEvent)
			done <- struct{}{}
		}()

		for range 100 {
			request := types.AlchemyRequest{
				Request: req,
				Body:    body,
			}
			requestEvent <- request
		}

		subscribers := batcher.Subscribe(done)
		assert.Equal(b, 100, len(subscribers.Responses))
	}
}

func BenchmarkRequest_Parallel(b *testing.B) {
	httpmock.Activate(b)
	defer httpmock.DeactivateAndReset()

	// Arrange
	targetUrl := "example.com"
	body := types.AlchemyRequestBody{
		Jsonrpc: "2.0",
		Method:  "method",
		Params:  []string{"param1", "param2"},
		Id:      1,
	}
	req, _ := http.NewRequest("POST", targetUrl, nil)
	mockResult := types.AlchemyResponse{
		Jsonrpc: "2.0",
		Id:      1,
		Result:  "0x1234",
	}
	resultJson, _ := json.Marshal(mockResult)

	// Mock
	httpmock.RegisterResponder(
		"POST",
		targetUrl,
		httpmock.NewStringResponder(200, string(resultJson)),
	)

	b.ResetTimer()
	for b.Loop() {
		for range 100 {
			request := types.AlchemyRequest{
				Request: req,
				Body:    body,
			}
			utils.AlchemyFetch(request, types.RequestConfig{
				Timeout: time.Second * 10,
			})
		}
	}
}
