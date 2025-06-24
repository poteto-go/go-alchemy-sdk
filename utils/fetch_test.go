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
	"github.com/poteto-go/go-alchemy-sdk/core"
	"github.com/poteto-go/go-alchemy-sdk/types"
	"github.com/poteto-go/go-alchemy-sdk/utils"
	"github.com/stretchr/testify/assert"
)

func TestAlchemyFetch(t *testing.T) {
	// Arrange
	targetUrl := "example.com"

	body := types.AlchemyRequestBody{
		Jsonrpc: "2.0",
		Method:  "method",
		Params:  []string{"param1", "param2"},
		Id:      1,
	}

	t.Run("normal case:", func(t *testing.T) {
		httpmock.Activate(t)
		defer httpmock.DeactivateAndReset()

		// Arrange
		req, _ := http.NewRequest("POST", targetUrl, nil)
		request := types.AlchemyRequest{
			Request: req,
			Body:    body,
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
		})

		// Assert
		assert.Nil(t, err)
		assert.Equal(t, mockResult, result)
	})

	t.Run("error case:", func(t *testing.T) {
		t.Run("if failed to marshal parameter -> core.ErrFailedToMarshalParameter", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Arrange
			req, _ := http.NewRequest("POST", targetUrl, nil)
			request := types.AlchemyRequest{
				Request: req,
				Body:    body,
			}

			// Mock
			patches.ApplyFunc(
				json.Marshal,
				func(v any) ([]byte, error) {
					return nil, errors.New("error")
				},
			)

			// Act
			_, err := utils.AlchemyFetch(request, types.RequestConfig{
				Timeout: 10 * time.Second,
			})

			// Assert
			assert.ErrorIs(t, core.ErrFailedToMarshalParameter, err)
		})

		t.Run("if failed to request -> core.ErrFailedToConnect", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Arrange
			req, _ := http.NewRequest("POST", targetUrl, nil)
			request := types.AlchemyRequest{
				Request: req,
				Body:    body,
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
			})

			// Assert
			assert.ErrorIs(t, core.ErrFailedToConnect, err)
		})

		t.Run("if failed to unmarshal response -> core.ErrFailedToUnmarshalResponse", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			httpmock.Activate(t)
			defer func() {
				httpmock.DeactivateAndReset()
				patches.Reset()
			}()

			// Arrange
			req, _ := http.NewRequest("POST", targetUrl, nil)
			request := types.AlchemyRequest{
				Request: req,
				Body:    body,
			}

			// Mock
			httpmock.RegisterResponder(
				"POST",
				targetUrl,
				httpmock.NewStringResponder(200, `ok`),
			)

			// Mock
			patches.ApplyFunc(
				json.Unmarshal,
				func(_ []byte, _ any) error {
					return errors.New("error")
				},
			)

			// Act
			_, err := utils.AlchemyFetch(request, types.RequestConfig{
				Timeout: 10 * time.Second,
			})

			// Assert
			assert.ErrorIs(t, core.ErrFailedToUnmarshalResponse, err)
		})
	})
}

func TestAlchemyBatchFetch(t *testing.T) {
	// Arrange
	targetUrl := "example.com"
	body1 := types.AlchemyRequestBody{
		Jsonrpc: "2.0",
		Method:  "method1",
		Params:  []string{"param1", "param2"},
		Id:      1,
	}
	body2 := types.AlchemyRequestBody{
		Jsonrpc: "2.0",
		Method:  "method2",
		Params:  []string{"param3", "param4"},
		Id:      2,
	}

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
				Body:    body1,
			}
			request2 := types.AlchemyRequest{
				Request: req,
				Body:    body2,
			}
			requests := []types.AlchemyRequest{request1, request2}

			// Mock
			httpmock.RegisterResponder(
				"POST",
				targetUrl,
				httpmock.NewStringResponder(200, string(resultJson)),
			)

			// Act
			result, err := utils.AlchemyBatchFetch(requests, types.RequestConfig{
				Timeout: 10 * time.Second,
			})

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
				Body:    body1,
			}
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
			})

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
				Body:    body1,
			}
			request2 := types.AlchemyRequest{
				Request: req,
				Body:    body2,
			}
			requests := []types.AlchemyRequest{request1, request2}

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
			})

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
				Body:    body1,
			}
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
			})

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
				Body:    body1,
			}
			request2 := types.AlchemyRequest{
				Request: req,
				Body:    body2,
			}
			requests := []types.AlchemyRequest{request1, request2}

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
			})

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
				Body:    body1,
			}
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
			})

			// Assert
			assert.ErrorIs(t, core.ErrFailedToUnmarshalResponse, err)
		})
	})
}
