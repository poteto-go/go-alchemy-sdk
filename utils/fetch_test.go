package utils_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey"
	"github.com/jarcoal/httpmock"
	"github.com/poteto-go/go-alchemy-sdk/constant"
	"github.com/poteto-go/go-alchemy-sdk/types"
	"github.com/poteto-go/go-alchemy-sdk/utils"
	"github.com/stretchr/testify/assert"
)

func TestAlchemyFetch(t *testing.T) {
	// Arrange
	targetUrl := "example.com"

	body, _ := utils.CreateRequestBodyToBytes(
		1,
		"method",
		func() types.RequestArgs {
			var params []any
			for _, p := range []string{"param1", "param2"} {
				params = append(params, p)
			}
			return params
		}(),
	)

	t.Run("normal case:", func(t *testing.T) {
		httpmock.Activate(t)
		defer httpmock.DeactivateAndReset()

		// Arrange
		req, _ := http.NewRequest("POST", targetUrl, nil)
		request := types.AlchemyRequest{
			Request: req,
		}
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

		// Act
		result, err := utils.AlchemyFetch(&http.Client{}, request, body)

		// Assert
		assert.Nil(t, err)
		assert.Equal(t, mockResult, result)
	})

	t.Run("error case:", func(t *testing.T) {

		t.Run("if failed to request -> constant.ErrFailedToConnect", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Arrange
			req, _ := http.NewRequest("POST", targetUrl, nil)
			request := types.AlchemyRequest{
				Request: req,
			}

			// Mock
			patches.ApplyMethod(
				reflect.TypeOf(http.DefaultClient),
				"Do",
				func(c *http.Client, req *http.Request) (*http.Response, error) {
					return nil, errors.New("error")
				},
			)

			// Act
			_, err := utils.AlchemyFetch(&http.Client{}, request, body)

			// Assert
			assert.ErrorIs(t, constant.ErrFailedToConnect, err)
		})

		t.Run("if failed to unmarshal response -> constant.ErrFailedToUnmarshalResponse", func(t *testing.T) {
			httpmock.Activate(t)
			patches := gomonkey.NewPatches()
			defer func() {
				patches.Reset()
				httpmock.DeactivateAndReset()
			}()

			// Arrange
			req, _ := http.NewRequest("POST", targetUrl, nil)
			request := types.AlchemyRequest{
				Request: req,
			}
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
			patches.ApplyFunc(
				json.Unmarshal,
				func(data []byte, v interface{}) error {
					return errors.New("error")
				},
			)

			// Act
			_, err := utils.AlchemyFetch(&http.Client{}, request, body)

			// Assert
			assert.ErrorIs(t, constant.ErrFailedToUnmarshalResponse, err)
		})
	})
}

func TestAlchemyFetch_ResponseSizeLimit(t *testing.T) {
	targetUrl := "example.com"

	body, _ := utils.CreateRequestBodyToBytes(
		1,
		"method",
		types.RequestArgs{},
	)

	t.Run("response within limit succeeds", func(t *testing.T) {
		httpmock.Activate(t)
		defer httpmock.DeactivateAndReset()

		req, _ := http.NewRequest("POST", targetUrl, nil)
		request := types.AlchemyRequest{Request: req}
		mockResult := types.AlchemyResponse{Jsonrpc: "2.0", Id: 1, Result: "0x1"}
		resultJson, _ := json.Marshal(mockResult)

		httpmock.RegisterResponder("POST", targetUrl, httpmock.NewStringResponder(200, string(resultJson)))

		// Client with limit larger than the response: succeeds.
		client := utils.NewSharedHTTPClient(int64(len(resultJson)+1), 0, nil)
		_, err := utils.AlchemyFetch(client, request, body)

		assert.Nil(t, err)
	})

	t.Run("response exceeds limit -> ErrFailedToUnmarshalResponse (truncated body fails JSON decode)", func(t *testing.T) {
		httpmock.Activate(t)
		defer httpmock.DeactivateAndReset()

		req, _ := http.NewRequest("POST", targetUrl, nil)
		request := types.AlchemyRequest{Request: req}
		largeBody := string(make([]byte, 100))

		httpmock.RegisterResponder("POST", targetUrl, httpmock.NewStringResponder(200, largeBody))

		// Client with a 10-byte limit: transport truncates body, JSON decode fails.
		client := utils.NewSharedHTTPClient(10, 0, nil)
		_, err := utils.AlchemyFetch(client, request, body)

		assert.ErrorIs(t, err, constant.ErrFailedToUnmarshalResponse)
	})
}

func TestAlchemyBatchFetch_ResponseSizeLimit(t *testing.T) {
	targetUrl := "example.com"

	body1, _ := utils.CreateRequestBodyToBytes(1, "method", types.RequestArgs{})
	body2, _ := utils.CreateRequestBodyToBytes(2, "method", types.RequestArgs{})
	req, _ := http.NewRequest("POST", targetUrl, nil)

	t.Run("single body: response exceeds limit -> ErrFailedToUnmarshalResponse", func(t *testing.T) {
		httpmock.Activate(t)
		defer httpmock.DeactivateAndReset()

		largeBody := string(make([]byte, 100))
		httpmock.RegisterResponder("POST", targetUrl, httpmock.NewStringResponder(200, largeBody))

		client := utils.NewSharedHTTPClient(10, 0, nil)
		_, err := utils.AlchemyBatchFetch(
			client,
			[]types.AlchemyRequest{{Request: req}},
			[][]byte{body1},
		)

		assert.ErrorIs(t, err, constant.ErrFailedToUnmarshalResponse)
	})

	t.Run("batch body: response exceeds limit -> ErrFailedToUnmarshalResponse", func(t *testing.T) {
		httpmock.Activate(t)
		defer httpmock.DeactivateAndReset()

		largeBody := string(make([]byte, 100))
		httpmock.RegisterResponder("POST", targetUrl, httpmock.NewStringResponder(200, largeBody))

		client := utils.NewSharedHTTPClient(10, 0, nil)
		_, err := utils.AlchemyBatchFetch(
			client,
			[]types.AlchemyRequest{{Request: req}, {Request: req}},
			[][]byte{body1, body2},
		)

		assert.ErrorIs(t, err, constant.ErrFailedToUnmarshalResponse)
	})
}

func TestAlchemyBatchFetch(t *testing.T) {
	// Arrange
	targetUrl := "example.com"

	t.Run("normal case:", func(t *testing.T) {
		t.Run("batched", func(t *testing.T) {
			httpmock.Activate(t)
			defer httpmock.DeactivateAndReset()

			// Arrange
			mockResult := []types.AlchemyResponse{
				{
					Jsonrpc: "2.0",
					Id:      1,
					Result:  "0x1234",
				},
			}
			resultJson, _ := json.Marshal(mockResult)

			req, _ := http.NewRequest("POST", targetUrl, nil)
			request1 := types.AlchemyRequest{
				Request: req,
			}
			body1, _ := utils.CreateRequestBodyToBytes(
				1,
				"method",
				types.RequestArgs{
					[]string{"param1", "param2"},
				},
			)
			request2 := types.AlchemyRequest{
				Request: req,
			}
			body2, _ := utils.CreateRequestBodyToBytes(
				2,
				"method",
				types.RequestArgs{
					[]string{"param3", "param4"},
				},
			)
			requests := []types.AlchemyRequest{request1, request2}
			bodies := [][]byte{body1, body2}

			// Mock
			httpmock.RegisterResponder(
				"POST",
				targetUrl,
				httpmock.NewStringResponder(200, string(resultJson)),
			)

			// Act
			result, err := utils.AlchemyBatchFetch(&http.Client{}, requests, bodies)

			// Assert
			assert.Nil(t, err)
			assert.Equal(t, mockResult, result)
		})

		t.Run("not batched", func(t *testing.T) {
			httpmock.Activate(t)
			defer httpmock.DeactivateAndReset()

			// Arrange
			mockResult := types.AlchemyResponse{
				Jsonrpc: "2.0",
				Id:      1,
				Result:  "0x1234",
			}
			resultJson, _ := json.Marshal(mockResult)
			req, _ := http.NewRequest("POST", targetUrl, nil)
			request1 := types.AlchemyRequest{
				Request: req,
			}
			body1, _ := utils.CreateRequestBodyToBytes(
				1,
				"method",
				types.RequestArgs{
					[]string{"param1", "param2"},
				},
			)
			requests := []types.AlchemyRequest{request1}

			// Mock
			httpmock.RegisterResponder(
				"POST",
				targetUrl,
				httpmock.NewStringResponder(200, string(resultJson)),
			)

			// Act
			result, err := utils.AlchemyBatchFetch(&http.Client{}, requests, [][]byte{body1})

			// Assert
			assert.Nil(t, err)
			assert.Equal(t, []types.AlchemyResponse{mockResult}, result)
		})
	})

	t.Run("error case:", func(t *testing.T) {
		t.Run("failed batched request -> constant.ErrFailedToConnect", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Arrange
			req, _ := http.NewRequest("POST", targetUrl, nil)
			request1 := types.AlchemyRequest{
				Request: req,
			}
			body1, _ := utils.CreateRequestBodyToBytes(
				1,
				"method",
				types.RequestArgs{
					[]string{"param1", "param2"},
				},
			)
			request2 := types.AlchemyRequest{
				Request: req,
			}
			body2, _ := utils.CreateRequestBodyToBytes(
				2,
				"method",
				types.RequestArgs{
					[]string{"param3", "param4"},
				},
			)
			requests := []types.AlchemyRequest{request1, request2}
			bodies := [][]byte{body1, body2}

			// Mock
			patches.ApplyMethod(
				reflect.TypeOf(http.DefaultClient),
				"Do",
				func(c *http.Client, req *http.Request) (*http.Response, error) {
					return nil, errors.New("error")
				},
			)

			// Act
			_, err := utils.AlchemyBatchFetch(&http.Client{}, requests, bodies)

			// Assert
			assert.ErrorIs(t, constant.ErrFailedToConnect, err)
		})

		t.Run("failed not batched request -> constant.ErrFailedToConnect", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Arrange
			req, _ := http.NewRequest("POST", targetUrl, nil)
			request1 := types.AlchemyRequest{
				Request: req,
			}
			body1, _ := utils.CreateRequestBodyToBytes(
				1,
				"method",
				types.RequestArgs{
					[]string{"param1", "param2"},
				},
			)
			requests := []types.AlchemyRequest{request1}

			// Mock
			patches.ApplyMethod(
				reflect.TypeOf(http.DefaultClient),
				"Do",
				func(c *http.Client, req *http.Request) (*http.Response, error) {
					return nil, errors.New("error")
				},
			)

			// Act
			_, err := utils.AlchemyBatchFetch(&http.Client{}, requests, [][]byte{body1})

			// Assert
			assert.ErrorIs(t, constant.ErrFailedToConnect, err)
		})

		t.Run("failed batched unmarshal -> constant.ErrFailedToUnmarshalResponse", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			httpmock.Activate(t)
			defer func() {
				httpmock.DeactivateAndReset()
				patches.Reset()
			}()

			// Arrange
			mockResult := types.AlchemyResponse{
				Jsonrpc: "2.0",
				Id:      1,
				Result:  "0x1234",
			}
			resultJson, _ := json.Marshal(mockResult)
			req, _ := http.NewRequest("POST", targetUrl, nil)
			request1 := types.AlchemyRequest{
				Request: req,
			}
			body1, _ := utils.CreateRequestBodyToBytes(
				1,
				"method",
				types.RequestArgs{
					[]string{"param1", "param2"},
				},
			)
			request2 := types.AlchemyRequest{
				Request: req,
			}
			body2, _ := utils.CreateRequestBodyToBytes(
				2,
				"method",
				types.RequestArgs{
					[]string{"param3", "param4"},
				},
			)
			requests := []types.AlchemyRequest{request1, request2}
			bodies := [][]byte{body1, body2}

			// Mock
			patches.ApplyFunc(
				json.Unmarshal,
				func(_ []byte, _ any) error {
					return errors.New("error")
				},
			)
			httpmock.RegisterResponder(
				"POST",
				targetUrl,
				httpmock.NewStringResponder(200, string(resultJson)),
			)

			// Act
			_, err := utils.AlchemyBatchFetch(&http.Client{}, requests, bodies)

			// Assert
			assert.ErrorIs(t, constant.ErrFailedToUnmarshalResponse, err)
		})

		t.Run("failed unmarshal -> constant.ErrFailedToUnmarshalResponse", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			httpmock.Activate(t)
			defer func() {
				httpmock.DeactivateAndReset()
				patches.Reset()
			}()

			// Arrange
			mockResult := types.AlchemyResponse{
				Jsonrpc: "2.0",
				Id:      1,
				Result:  "0x1234",
			}
			resultJson, _ := json.Marshal(mockResult)
			req, _ := http.NewRequest("POST", targetUrl, nil)
			request1 := types.AlchemyRequest{
				Request: req,
			}
			body1, _ := utils.CreateRequestBodyToBytes(
				1,
				"method",
				types.RequestArgs{
					[]string{"param1", "param2"},
				},
			)
			requests := []types.AlchemyRequest{request1}

			// Mock
			patches.ApplyFunc(
				json.Unmarshal,
				func(_ []byte, _ any) error {
					return errors.New("error")
				},
			)
			httpmock.RegisterResponder(
				"POST",
				targetUrl,
				httpmock.NewStringResponder(200, string(resultJson)),
			)

			// Act
			_, err := utils.AlchemyBatchFetch(&http.Client{}, requests, [][]byte{body1})

			// Assert
			assert.ErrorIs(t, constant.ErrFailedToUnmarshalResponse, err)
		})
	})
}

// Ensure timeout set on the client is respected.
func TestAlchemyFetch_ClientTimeout(t *testing.T) {
	httpmock.Activate(t)
	defer httpmock.DeactivateAndReset()

	targetUrl := "example.com"
	req, _ := http.NewRequest("POST", targetUrl, nil)
	body, _ := utils.CreateRequestBodyToBytes(1, "method", types.RequestArgs{})

	mockResult := types.AlchemyResponse{Jsonrpc: "2.0", Id: 1, Result: "ok"}
	resultJson, _ := json.Marshal(mockResult)
	httpmock.RegisterResponder("POST", targetUrl, httpmock.NewStringResponder(200, string(resultJson)))

	client := utils.NewSharedHTTPClient(0, 10*time.Second, nil)
	result, err := utils.AlchemyFetch(client, types.AlchemyRequest{Request: req}, body)

	assert.Nil(t, err)
	assert.Equal(t, mockResult, result)
}
