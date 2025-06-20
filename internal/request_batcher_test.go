package internal

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/poteto-go/go-alchemy-sdk/types"
	"github.com/poteto-go/go-alchemy-sdk/utils"
	"github.com/stretchr/testify/assert"
)

func timeoutFetch(reqs []types.AlchemyRequest) ([]types.AlchemyResponse, error) {
	results := make([]types.AlchemyResponse, len(reqs))
	time.Sleep(time.Millisecond * 30)
	return results, nil
}

func errorFetch(reqs []types.AlchemyRequest) ([]types.AlchemyResponse, error) {
	return nil, errors.New("error")
}

func newBatcher() *RequestBatcher {
	ctx := context.Background()
	return NewRequestBatcher(
		ctx,
		BatcherConfig{
			MaxBatchSize: 2,
			MaxBatchTime: time.Millisecond * 10,
			Fetch:        utils.AlchemyBatchFetch,
		},
	).(*RequestBatcher)
}

func Test_RecordRequest(t *testing.T) {
	// Arrange
	targetUrl := "example.com"
	body := types.AlchemyRequestBody{
		Jsonrpc: "2.0",
		Method:  "method",
		Params:  []string{"param1", "param2"},
		Id:      1,
	}
	req, _ := http.NewRequest("POST", targetUrl, nil)
	httpmock.Activate(t)
	defer httpmock.DeactivateAndReset()
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

	t.Run("normal case:", func(t *testing.T) {
		t.Run("wait for max batch size", func(t *testing.T) {
			batcher := newBatcher()

			// Act
			requestChan := make(chan types.AlchemyRequest, 1)
			done := make(chan struct{}, 1)
			go func() {
				batcher.RecordRequest(requestChan)
				done <- struct{}{}
			}()
			for range 2 {
				request := types.AlchemyRequest{
					Request: req,
					Body:    body,
				}
				requestChan <- request
			}
			subscribers := batcher.Subscribe(done)

			assert.Equal(t, 2, len(subscribers.Responses))
			assert.Nil(t, subscribers.JoinedError)
			// test reset
			assert.Equal(t, 0, len(batcher.Subscribers.Responses))
		})

		t.Run("wait for max batch time", func(t *testing.T) {
			batcher := newBatcher()
			batcher.config.Fetch = timeoutFetch
			defer func() {
				batcher.config.Fetch = utils.AlchemyBatchFetch
			}()

			// Act
			requestChan := make(chan types.AlchemyRequest, 1)
			done := make(chan struct{}, 1)
			go func() {
				batcher.RecordRequest(requestChan)
				done <- struct{}{}
			}()

			request := types.AlchemyRequest{
				Request: req,
				Body:    body,
			}
			requestChan <- request

			time.Sleep(time.Millisecond * 20)

			subscribers := batcher.Subscribe(done)

			assert.Equal(t, 1, len(subscribers.Responses))
			assert.Nil(t, batcher.Subscribers.JoinedError)
		})

		t.Run("batcher error", func(t *testing.T) {
			batcher := newBatcher()
			batcher.config.Fetch = errorFetch
			defer func() {
				batcher.config.Fetch = utils.AlchemyBatchFetch
			}()

			requestChan := make(chan types.AlchemyRequest, 1)
			done := make(chan struct{}, 1)
			go func() {
				batcher.RecordRequest(requestChan)
				done <- struct{}{}
			}()

			for range 2 {
				request := types.AlchemyRequest{
					Request: req,
					Body:    body,
				}
				requestChan <- request
			}

			subscribers := batcher.Subscribe(done)

			assert.NotNil(t, subscribers.JoinedError)
			assert.Nil(t, batcher.Subscribers.JoinedError)
		})

		t.Run("no request for batch timeout", func(t *testing.T) {
			batcher := newBatcher()

			requestChan := make(chan types.AlchemyRequest, 1)
			done := make(chan struct{}, 1)
			go func() {
				batcher.RecordRequest(requestChan)
				done <- struct{}{}
			}()

			subscribers := batcher.Subscribe(done)

			assert.Equal(t, 0, len(subscribers.Responses))
		})
	})

	t.Run("abnormal case:", func(t *testing.T) {
		t.Run("not running", func(t *testing.T) {
			batcher := newBatcher()
			batcher.IsRun = false

			requestEvent := make(chan types.AlchemyRequest, 1)
			go batcher.RecordRequest(requestEvent)

			assert.Equal(t, 0, len(batcher.requests))
		})
	})
}
