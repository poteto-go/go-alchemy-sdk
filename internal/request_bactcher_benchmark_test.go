package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/poteto-go/go-alchemy-sdk/core"
	"github.com/poteto-go/go-alchemy-sdk/types"
	"github.com/poteto-go/go-alchemy-sdk/utils"
	"github.com/poteto-go/tslice"
	"github.com/stretchr/testify/assert"
)

func alchemyBatchFetch(reqs []types.AlchemyRequest) ([]types.AlchemyResponse, error) {
	request := reqs[0].Request
	bodies := tslice.Map(reqs, func(req types.AlchemyRequest) types.AlchemyRequestBody {
		return req.Body
	})

	if len(bodies) == 1 {
		paramJson, _ := json.Marshal(bodies[0])

		request.Body = io.NopCloser(bytes.NewBuffer(paramJson))
		res, err := http.DefaultClient.Do(request)
		if err != nil {
			return []types.AlchemyResponse{}, core.ErrFailedToConnect
		}
		defer res.Body.Close()

		body, _ := io.ReadAll(res.Body)
		result := types.AlchemyResponse{}
		if err := json.Unmarshal(body, &result); err != nil {
			return []types.AlchemyResponse{}, core.ErrFailedToUnmarshalResponse
		}

		return []types.AlchemyResponse{result}, nil
	}

	paramJson, _ := json.Marshal(bodies)

	request.Body = io.NopCloser(bytes.NewBuffer(paramJson))
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		return []types.AlchemyResponse{}, core.ErrFailedToConnect
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	results := []types.AlchemyResponse{}
	if err := json.Unmarshal(body, &results); err != nil {
		return []types.AlchemyResponse{}, core.ErrFailedToUnmarshalResponse
	}

	return results, nil
}

func timeoutFetch(reqs []types.AlchemyRequest) ([]types.AlchemyResponse, error) {
	results := make([]types.AlchemyResponse, len(reqs))
	time.Sleep(time.Millisecond * 30)
	return results, nil
}

func errorFetch(reqs []types.AlchemyRequest) ([]types.AlchemyResponse, error) {
	return nil, errors.New("error")
}

func TestNewRequestBatcher(t *testing.T) {
	// Arrange
	ctx := context.Background()

	// Act
	batcher := NewRequestBatcher(
		ctx,
		BatcherConfig{
			MaxBatchSize: 10,
			MaxBatchTime: time.Millisecond * 10,
			Fetch:        alchemyBatchFetch,
		},
	).(*RequestBatcher)

	// Assert
	assert.True(t, batcher.IsRun)
}

func newBatcher() *RequestBatcher {
	ctx := context.Background()
	return NewRequestBatcher(
		ctx,
		BatcherConfig{
			MaxBatchSize: 100,
			MaxBatchTime: time.Millisecond * 10,
			Fetch:        alchemyBatchFetch,
		},
	).(*RequestBatcher)
}

func Benchmark_BatchRequest(b *testing.B) {
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

	b.ResetTimer()
	for range b.N {
		batcher := newBatcher()

		requestEvent := make(chan types.AlchemyRequest, 1)
		go batcher.RecordRequest(requestEvent)

		for j := 0; j < 100; j++ {
			request := types.AlchemyRequest{
				Request: req,
				Body:    body,
			}
			requestEvent <- request
		}
	}
}

func Benchmark_Parallel(b *testing.B) {
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

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 100; j++ {
			request := types.AlchemyRequest{
				Request: req,
				Body:    body,
			}
			utils.AlchemyFetch(request)
		}
	}
}
