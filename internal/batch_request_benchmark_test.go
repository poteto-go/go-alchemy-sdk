package internal

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/poteto-go/go-alchemy-sdk/types"
	"github.com/poteto-go/go-alchemy-sdk/utils"
)

func BenchmarkRequestBatcher_QueueRequest(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"jsonrpc":"2.0","id":1,"result":"0x1"}]`))
	}))
	defer server.Close()

	config := BatcherConfig{
		MaxBatchSize: 100,
		MaxBatchTime: time.Millisecond * 10,
		Client:       utils.NewSharedHTTPClient(0, time.Second, nil),
		Fetch:        utils.AlchemyBatchFetch,
	}

	batcher := NewRequestBatcher(context.Background(), config)

	req, _ := http.NewRequest("POST", server.URL, nil)
	request := types.AlchemyRequest{
		Request: req,
	}
	body, _ := utils.CreateRequestBodyToBytes(1, "eth_blockNumber", types.RequestArgs{})

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		go func() {
			batcher.QueueRequest(context.Background(), request, body)
		}()
	}
}
