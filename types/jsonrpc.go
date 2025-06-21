package types

import "net/http"

type AlchemyRequestBody struct {
	Method  string   `json:"method"`
	Params  []string `json:"params"`
	Id      int      `json:"id"`
	Jsonrpc string   `json:"jsonrpc"`
}

type AlchemyRequest struct {
	Body    AlchemyRequestBody `json:"-"`
	Request *http.Request      `json:"-"`
}

type AlchemyResponse struct {
	Result  string `json:"result"`
	Id      int    `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
}

type AlchemyFetchHandler func(AlchemyRequest) (AlchemyResponse, error)
