package utils_test

import (
	"errors"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/goccy/go-json"

	"github.com/agiledragon/gomonkey"
	"github.com/jarcoal/httpmock"
	"github.com/poteto-go/go-alchemy-sdk/core"
	"github.com/poteto-go/go-alchemy-sdk/types"
	"github.com/poteto-go/go-alchemy-sdk/utils"
	"github.com/stretchr/testify/assert"
)

func TestAlchemyFetch(t *testing.T) {
	// Arrange
	targetUrl := "example.com"

	body, _ := utils.CreateRequestBodyToBytes(1, "method", []string{"param1", "param2"})

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
		result, err := utils.AlchemyFetch(request, types.RequestConfig{
			Timeout: 10 * time.Second,
		}, body)

		// Assert
		assert.Nil(t, err)
		assert.Equal(t, mockResult, result)
	})

	t.Run("error case:", func(t *testing.T) {

		t.Run("if failed to request -> core.ErrFailedToConnect", func(t *testing.T) {
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
			_, err := utils.AlchemyFetch(request, types.RequestConfig{
				Timeout: 10 * time.Second,
			}, body)

			// Assert
			assert.ErrorIs(t, core.ErrFailedToConnect, err)
		})

		t.Run("if failed to unmarshal response -> core.ErrFailedToUnmarshalResponse", func(t *testing.T) {
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
			_, err := utils.AlchemyFetch(request, types.RequestConfig{
				Timeout: 10 * time.Second,
			}, body)

			// Assert
			assert.ErrorIs(t, core.ErrFailedToUnmarshalResponse, err)
		})
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
			body1, _ := utils.CreateRequestBodyToBytes(1, "method", []string{"param1", "param2"})
			request2 := types.AlchemyRequest{
				Request: req,
			}
			body2, _ := utils.CreateRequestBodyToBytes(2, "method", []string{"param3", "param4"})
			requests := []types.AlchemyRequest{request1, request2}
			bodies := [][]byte{body1, body2}

			// Mock
			httpmock.RegisterResponder(
				"POST",
				targetUrl,
				httpmock.NewStringResponder(200, string(resultJson)),
			)

			// Act
			result, err := utils.AlchemyBatchFetch(requests, types.RequestConfig{
				Timeout: 10 * time.Second,
			}, bodies)

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
			body1, _ := utils.CreateRequestBodyToBytes(1, "method", []string{"param1", "param2"})
			requests := []types.AlchemyRequest{request1}

			// Mock
			httpmock.RegisterResponder(
				"POST",
				targetUrl,
				httpmock.NewStringResponder(200, string(resultJson)),
			)

			// Act
			result, err := utils.AlchemyBatchFetch(requests, types.RequestConfig{
				Timeout: 10 * time.Second,
			}, [][]byte{body1})

			// Assert
			assert.Nil(t, err)
			assert.Equal(t, []types.AlchemyResponse{mockResult}, result)
		})
	})

	t.Run("error case:", func(t *testing.T) {
		t.Run("failed batched request -> core.ErrFailedToConnect", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Arrange
			req, _ := http.NewRequest("POST", targetUrl, nil)
			request1 := types.AlchemyRequest{
				Request: req,
			}
			body1, _ := utils.CreateRequestBodyToBytes(1, "method", []string{"param1", "param2"})
			request2 := types.AlchemyRequest{
				Request: req,
			}
			body2, _ := utils.CreateRequestBodyToBytes(2, "method", []string{"param3", "param4"})
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
			_, err := utils.AlchemyBatchFetch(requests, types.RequestConfig{
				Timeout: 10 * time.Second,
			}, bodies)

			// Assert
			assert.ErrorIs(t, core.ErrFailedToConnect, err)
		})

		t.Run("failed not batched request -> core.ErrFailedToConnect", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Arrange
			req, _ := http.NewRequest("POST", targetUrl, nil)
			request1 := types.AlchemyRequest{
				Request: req,
			}
			body1, _ := utils.CreateRequestBodyToBytes(1, "method", []string{"param1", "param2"})
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
			_, err := utils.AlchemyBatchFetch(requests, types.RequestConfig{
				Timeout: 10 * time.Second,
			}, [][]byte{body1})

			// Assert
			assert.ErrorIs(t, core.ErrFailedToConnect, err)
		})

		t.Run("failed batched unmarshal -> core.ErrFailedToUnmarshalResponse", func(t *testing.T) {
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
			body1, _ := utils.CreateRequestBodyToBytes(1, "method", []string{"param1", "param2"})
			request2 := types.AlchemyRequest{
				Request: req,
			}
			body2, _ := utils.CreateRequestBodyToBytes(2, "method", []string{"param3", "param4"})
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
			_, err := utils.AlchemyBatchFetch(requests, types.RequestConfig{
				Timeout: 10 * time.Second,
			}, bodies)

			// Assert
			assert.ErrorIs(t, core.ErrFailedToUnmarshalResponse, err)
		})

		t.Run("failed unmarshal -> core.ErrFailedToUnmarshalResponse", func(t *testing.T) {
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
			body1, _ := utils.CreateRequestBodyToBytes(1, "method", []string{"param1", "param2"})
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
			_, err := utils.AlchemyBatchFetch(requests, types.RequestConfig{
				Timeout: 10 * time.Second,
			}, [][]byte{body1})

			// Assert
			assert.ErrorIs(t, core.ErrFailedToUnmarshalResponse, err)
		})
	})
}
