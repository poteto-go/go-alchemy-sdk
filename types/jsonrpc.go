package types

import (
	"net/http"
	"time"
)

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

type AlchemyFetchHandler func(AlchemyRequest, RequestConfig) (AlchemyResponse, error)

type BatchAlchemyFetchHandler func([]AlchemyRequest, RequestConfig) ([]AlchemyResponse, error)

type RequestConfig struct {
	Timeout time.Duration
}
