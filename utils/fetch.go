package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/poteto-go/go-alchemy-sdk/constant"
	"github.com/poteto-go/go-alchemy-sdk/types"
)

func AlchemyFetch(
	client *http.Client,
	req types.AlchemyRequest,
	requestConfig types.RequestConfig,
	body []byte,
) (types.AlchemyResponse, error) {
	httpReq := req.Request
	if requestConfig.Timeout > 0 {
		ctx, cancel := context.WithTimeout(httpReq.Context(), requestConfig.Timeout)
		defer cancel()
		httpReq = httpReq.WithContext(ctx)
	}
	httpReq.Body = io.NopCloser(bytes.NewBuffer(body))

	res, err := client.Do(httpReq)
	if err != nil {
		return types.AlchemyResponse{}, constant.ErrFailedToConnect
	}
	defer res.Body.Close()

	// Response body size is capped by the client's limitedTransport.
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return types.AlchemyResponse{}, constant.ErrFailedToReadResponse
	}

	result := types.AlchemyResponse{}
	if err := json.Unmarshal(resBody, &result); err != nil {
		return types.AlchemyResponse{}, constant.ErrFailedToUnmarshalResponse
	}
	return result, nil
}

func AlchemyBatchFetch(
	client *http.Client,
	reqs []types.AlchemyRequest,
	requestConfig types.RequestConfig,
	bodies [][]byte,
) ([]types.AlchemyResponse, error) {
	request := reqs[0].Request

	if len(bodies) == 1 {
		httpReq := request
		if requestConfig.Timeout > 0 {
			ctx, cancel := context.WithTimeout(httpReq.Context(), requestConfig.Timeout)
			defer cancel()
			httpReq = httpReq.WithContext(ctx)
		}
		httpReq.Body = io.NopCloser(bytes.NewBuffer(bodies[0]))

		res, err := client.Do(httpReq)
		if err != nil {
			return []types.AlchemyResponse{}, constant.ErrFailedToConnect
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return []types.AlchemyResponse{}, constant.ErrFailedToReadResponse
		}

		result := types.AlchemyResponse{}
		if err := json.Unmarshal(body, &result); err != nil {
			return []types.AlchemyResponse{}, constant.ErrFailedToUnmarshalResponse
		}

		return []types.AlchemyResponse{result}, nil
	}

	paramJson, _ := json.Marshal(bodies)

	httpReq := request
	if requestConfig.Timeout > 0 {
		ctx, cancel := context.WithTimeout(httpReq.Context(), requestConfig.Timeout)
		defer cancel()
		httpReq = httpReq.WithContext(ctx)
	}
	httpReq.Body = io.NopCloser(bytes.NewBuffer(paramJson))

	res, err := client.Do(httpReq)
	if err != nil {
		return []types.AlchemyResponse{}, constant.ErrFailedToConnect
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return []types.AlchemyResponse{}, constant.ErrFailedToReadResponse
	}

	results := []types.AlchemyResponse{}
	if err := json.Unmarshal(body, &results); err != nil {
		return []types.AlchemyResponse{}, constant.ErrFailedToUnmarshalResponse
	}

	return results, nil
}
