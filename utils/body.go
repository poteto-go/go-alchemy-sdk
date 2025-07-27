package utils

import (
	"github.com/goccy/go-json"

	"github.com/poteto-go/go-alchemy-sdk/core"
	"github.com/poteto-go/go-alchemy-sdk/types"
)

func CreateRequestBodyToBytes(id int, method string, params types.RequestArgs) ([]byte, error) {
	body := types.AlchemyRequestBody{
		Id:      id,
		Jsonrpc: "2.0",
		Method:  method,
		Params:  params,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, core.ErrFailedToMarshalParameter
	}
	return jsonBody, nil
}
