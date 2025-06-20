package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/poteto-go/go-alchemy-sdk/core"
	"github.com/poteto-go/go-alchemy-sdk/types"
	"github.com/poteto-go/tslice"
)

func AlchemyFetch(req types.AlchemyRequest) (types.AlchemyResponse, error) {
	paramJson, err := json.Marshal(req.Body)
	if err != nil {
		return types.AlchemyResponse{}, core.ErrFailedToMarshalParameter
	}

	req.Request.Body = io.NopCloser(bytes.NewBuffer(paramJson))
	res, err := http.DefaultClient.Do(req.Request)
	if err != nil {
		return types.AlchemyResponse{}, core.ErrFailedToConnect
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	result := types.AlchemyResponse{}
	if err := json.Unmarshal(body, &result); err != nil {
		return types.AlchemyResponse{}, core.ErrFailedToUnmarshalResponse
	}

	return result, nil
}

func AlchemyBatchFetch(reqs []types.AlchemyRequest) ([]types.AlchemyResponse, error) {
	request := reqs[0].Request
	bodies := tslice.Map(reqs, func(req types.AlchemyRequest) types.AlchemyRequestBody {
		return req.Body
	})

	if len(bodies) == 1 {
		paramJson, _ := json.Marshal(bodies[0])

		request.Body = io.NopCloser(bytes.NewBuffer(paramJson))
		res, err := http.DefaultClient.Do(request)
		if err != nil {
			return []types.AlchemyResponse{}, core.ErrFailedToConnect
		}
		defer res.Body.Close()

		body, _ := io.ReadAll(res.Body)
		result := types.AlchemyResponse{}
		if err := json.Unmarshal(body, &result); err != nil {
			return []types.AlchemyResponse{}, core.ErrFailedToUnmarshalResponse
		}

		return []types.AlchemyResponse{result}, nil
	}

	paramJson, _ := json.Marshal(bodies)

	request.Body = io.NopCloser(bytes.NewBuffer(paramJson))
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		return []types.AlchemyResponse{}, core.ErrFailedToConnect
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	results := []types.AlchemyResponse{}
	if err := json.Unmarshal(body, &results); err != nil {
		return []types.AlchemyResponse{}, core.ErrFailedToUnmarshalResponse
	}

	return results, nil
}
