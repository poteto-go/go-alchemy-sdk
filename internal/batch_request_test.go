package internal

import (
	"context"
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
			Body:    types.AlchemyRequestBody{Id: 1},
		}
		res, err := batcher.QueueRequest(request)
		assert.NoError(t, err)
		assert.Equal(t, "0x1", res.Result)
	}()

	go func() {
		defer wg.Done()
		req, _ := http.NewRequest("POST", server.URL, nil)
		request := types.AlchemyRequest{
			Request: req,
			Body:    types.AlchemyRequestBody{Id: 2},
		}
		res, err := batcher.QueueRequest(request)
		assert.NoError(t, err)
		assert.Equal(t, "0x2", res.Result)
	}()

	wg.Wait()

	time.Sleep(time.Millisecond * 200)

	assert.Equal(t, 1, requestCount)
}