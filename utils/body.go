package utils

import (
	"github.com/goccy/go-json"

	"github.com/poteto-go/go-alchemy-sdk/core"
	"github.com/poteto-go/go-alchemy-sdk/types"
)

func CreateRequestBodyToBytes[
	T string | types.TransactionRequest | types.Filter | types.TransactionRequestWithBlockTag,
](id int, method string, params []T) ([]byte, error) {
	body := types.AlchemyRequestBody[T]{
		Id:      id,
		Jsonrpc: "2.0",
		Method:  method,
	}

	if len(params) > 0 {
		body.Params = params
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, core.ErrFailedToMarshalParameter
	}
	return jsonBody, nil
}
