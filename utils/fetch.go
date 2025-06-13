package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/poteto-go/go-alchemy-sdk/core"
	"github.com/poteto-go/go-alchemy-sdk/types"
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
