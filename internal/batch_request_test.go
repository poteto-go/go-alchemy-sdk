package internal

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/poteto-go/go-alchemy-sdk/types"
	"github.com/poteto-go/go-alchemy-sdk/utils"
	"github.com/stretchr/testify/assert"
)

func TestRequestBatcher_QueueRequest(t *testing.T) {
	var (
		mu           sync.Mutex
		requestCount int
	)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		requestCount++
		mu.Unlock()
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"jsonrpc":"2.0","id":1,"result":"0x1"},{"jsonrpc":"2.0","id":2,"result":"0x2"}]`))
	}))
	defer server.Close()

	config := BatcherConfig{
		MaxBatchSize: 2,
		MaxBatchTime: time.Millisecond * 100,
		Fetch:        utils.AlchemyBatchFetch,
	}

	requestConfig := types.RequestConfig{
		Timeout: time.Second,
	}

	batcher := NewRequestBatcher(context.Background(), config, requestConfig)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		req, _ := http.NewRequest("POST", server.URL, nil)
		request := types.AlchemyRequest{
			Request: req,
		}
		body, _ := utils.CreateRequestBodyToBytes(1, "method", []string{})
		res, err := batcher.QueueRequest(context.Background(), request, body)
		assert.NoError(t, err)
		assert.Equal(t, "0x1", res.Result)
	}()

	go func() {
		defer wg.Done()
		req, _ := http.NewRequest("POST", server.URL, nil)
		request := types.AlchemyRequest{
			Request: req,
		}
		body, _ := utils.CreateRequestBodyToBytes(2, "method", []string{})
		res, err := batcher.QueueRequest(context.Background(), request, body)
		assert.NoError(t, err)
		assert.Equal(t, "0x2", res.Result)
	}()

	wg.Wait()

	time.Sleep(time.Millisecond * 200)

	assert.Equal(t, 1, requestCount)
}

func TestRequestBatcher_QueueRequest_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	config := BatcherConfig{
		MaxBatchSize: 1,
		MaxBatchTime: time.Millisecond * 100,
		Fetch:        utils.AlchemyBatchFetch,
	}

	requestConfig := types.RequestConfig{
		Timeout: time.Second,
	}

	batcher := NewRequestBatcher(context.Background(), config, requestConfig)

	req, _ := http.NewRequest("POST", server.URL, nil)
	request := types.AlchemyRequest{
		Request: req,
	}
	body, _ := utils.CreateRequestBodyToBytes(1, "method", []string{})

	_, err := batcher.QueueRequest(context.Background(), request, body)
	assert.Error(t, err)
}

func TestRequestBatcher_Context_Cancel(t *testing.T) {
	config := BatcherConfig{
		MaxBatchSize: 1,
		MaxBatchTime: time.Millisecond * 100,
		Fetch:        utils.AlchemyBatchFetch,
	}

	requestConfig := types.RequestConfig{
		Timeout: time.Second,
	}

	ctx, cancel := context.WithCancel(context.Background())
	batcher := NewRequestBatcher(ctx, config, requestConfig)

	cancel()

	req, _ := http.NewRequest("POST", "", nil)
	request := types.AlchemyRequest{
		Request: req,
	}
	body, _ := utils.CreateRequestBodyToBytes(1, "method", []string{})

	_, err := batcher.QueueRequest(ctx, request, body)
	assert.Error(t, err)
}

func TestRequestBatcher_Flush_Fetch_Error(t *testing.T) {
	config := BatcherConfig{
		MaxBatchSize: 1,
		MaxBatchTime: time.Millisecond * 100,
		Fetch: func(reqs []types.AlchemyRequest, config types.RequestConfig, bodies [][]byte) ([]types.AlchemyResponse, error) {
			return nil, errors.New("fetch error")
		},
	}

	requestConfig := types.RequestConfig{
		Timeout: time.Second,
	}

	batcher := NewRequestBatcher(context.Background(), config, requestConfig)

	req, _ := http.NewRequest("POST", "", nil)
	request := types.AlchemyRequest{
		Request: req,
	}
	body, _ := utils.CreateRequestBodyToBytes(1, "method", []string{})

	_, err := batcher.QueueRequest(context.Background(), request, body)
	assert.Error(t, err)
}

func TestRequestBatcher_Flush_No_Response(t *testing.T) {
	config := BatcherConfig{
		MaxBatchSize: 1,
		MaxBatchTime: time.Millisecond * 100,
		Fetch: func(reqs []types.AlchemyRequest, config types.RequestConfig, bodies [][]byte) ([]types.AlchemyResponse, error) {
			return []types.AlchemyResponse{}, nil
		},
	}

	requestConfig := types.RequestConfig{
		Timeout: time.Second,
	}

	batcher := NewRequestBatcher(context.Background(), config, requestConfig)

	req, _ := http.NewRequest("POST", "", nil)
	request := types.AlchemyRequest{
		Request: req,
	}
	body, _ := utils.CreateRequestBodyToBytes(1, "method", []string{})

	_, err := batcher.QueueRequest(context.Background(), request, body)
	assert.Error(t, err)
}

func TestRequestBatcher_ProcessQueue_Timeout(t *testing.T) {
	config := BatcherConfig{
		MaxBatchSize: 1,
		MaxBatchTime: time.Millisecond * 10,
		Fetch: func(reqs []types.AlchemyRequest, config types.RequestConfig, bodies [][]byte) ([]types.AlchemyResponse, error) {
			return []types.AlchemyResponse{}, nil
		},
	}

	requestConfig := types.RequestConfig{
		Timeout: time.Second,
	}

	NewRequestBatcher(context.Background(), config, requestConfig)

	time.Sleep(time.Millisecond * 20)
}
