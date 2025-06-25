package alchemy

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/poteto-go/go-alchemy-sdk/internal"
	"github.com/poteto-go/go-alchemy-sdk/types"
	"github.com/poteto-go/go-alchemy-sdk/utils"
)

func TestRequestBatcher_QueueRequest(t *testing.T) {
	var (
		mu          sync.Mutex
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

	config := internal.BatcherConfig{
		MaxBatchSize: 2,
		MaxBatchTime: time.Millisecond * 100,
		Fetch:        utils.AlchemyBatchFetch,
	}

	requestConfig := types.RequestConfig{
		Timeout: time.Second,
	}

	batcher := internal.NewRequestBatcher(context.Background(), config, requestConfig)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		req, _ := http.NewRequest("POST", server.URL, nil)
		request := types.AlchemyRequest{
			Request: req,
			Body:    types.AlchemyRequestBody{Id: 1},
		}
		batcher.QueueRequest(request)
	}()

	go func() {
		defer wg.Done()
		req, _ := http.NewRequest("POST", server.URL, nil)
		request := types.AlchemyRequest{
			Request: req,
			Body:    types.AlchemyRequestBody{Id: 2},
		}
		batcher.QueueRequest(request)
	}()

	wg.Wait()

	time.Sleep(time.Millisecond * 200)

	if requestCount != 1 {
		t.Errorf("Expected 1 request, but got %d", requestCount)
	}
}
