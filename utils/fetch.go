package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/poteto-go/go-alchemy-sdk/constant"
	"github.com/poteto-go/go-alchemy-sdk/types"
)

func AlchemyFetch(
	req types.AlchemyRequest,
	requestConfig types.RequestConfig,
	body []byte,
) (types.AlchemyResponse, error) {
	client := &http.Client{
		Timeout: requestConfig.Timeout,
	}

	req.Request.Body = io.NopCloser(bytes.NewBuffer(body))
	res, err := client.Do(req.Request)
	if err != nil {
		return types.AlchemyResponse{}, constant.ErrFailedToConnect
	}
	defer res.Body.Close()

	maxBytes := requestConfig.MaxResponseBytes
	if maxBytes == 0 {
		maxBytes = types.DefaultMaxResponseBytes
	}
	resBody, err := io.ReadAll(io.LimitReader(res.Body, maxBytes+1))
	if err != nil {
		return types.AlchemyResponse{}, constant.ErrFailedToReadResponse
	}
	if int64(len(resBody)) > maxBytes {
		return types.AlchemyResponse{}, constant.ErrFailedToReadResponse
	}

	result := types.AlchemyResponse{}
	if err := json.Unmarshal(resBody, &result); err != nil {
		return types.AlchemyResponse{}, constant.ErrFailedToUnmarshalResponse
	}
	return result, nil
}

func AlchemyBatchFetch(
	reqs []types.AlchemyRequest,
	requestConfig types.RequestConfig,
	bodies [][]byte,
) ([]types.AlchemyResponse, error) {
	request := reqs[0].Request

	client := &http.Client{
		Timeout: requestConfig.Timeout,
	}

	maxBytes := requestConfig.MaxResponseBytes
	if maxBytes == 0 {
		maxBytes = types.DefaultMaxResponseBytes
	}

	if len(bodies) == 1 {
		request.Body = io.NopCloser(bytes.NewBuffer(bodies[0]))
		res, err := client.Do(request)
		if err != nil {
			return []types.AlchemyResponse{}, constant.ErrFailedToConnect
		}
		defer res.Body.Close()

		body, err := io.ReadAll(io.LimitReader(res.Body, maxBytes+1))
		if err != nil {
			return []types.AlchemyResponse{}, constant.ErrFailedToReadResponse
		}
		if int64(len(body)) > maxBytes {
			return []types.AlchemyResponse{}, constant.ErrFailedToReadResponse
		}

		result := types.AlchemyResponse{}
		if err := json.Unmarshal(body, &result); err != nil {
			return []types.AlchemyResponse{}, constant.ErrFailedToUnmarshalResponse
		}

		return []types.AlchemyResponse{result}, nil
	}

	paramJson, _ := json.Marshal(bodies)

	request.Body = io.NopCloser(bytes.NewBuffer(paramJson))
	res, err := client.Do(request)
	if err != nil {
		return []types.AlchemyResponse{}, constant.ErrFailedToConnect
	}
	defer res.Body.Close()

	body, err := io.ReadAll(io.LimitReader(res.Body, maxBytes+1))
	if err != nil {
		return []types.AlchemyResponse{}, constant.ErrFailedToReadResponse
	}
	if int64(len(body)) > maxBytes {
		return []types.AlchemyResponse{}, constant.ErrFailedToReadResponse
	}

	results := []types.AlchemyResponse{}
	if err := json.Unmarshal(body, &results); err != nil {
		return []types.AlchemyResponse{}, constant.ErrFailedToUnmarshalResponse
	}

	return results, nil
}
