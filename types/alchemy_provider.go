package types

import (
	"errors"
	"net/http"
	"time"
)

var ErrNoResultFound = errors.New("no result found")

type AlchemyRequest struct {
	Request *http.Request
	Body    AlchemyRequestBody
}

type AlchemyRequestBody struct {
	Jsonrpc string   `json:"jsonrpc"`
	Method  string   `json:"method"`
	Params  []string `json:"params"`
	Id      int      `json:"id"`
}

type AlchemyResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      int    `json:"id"`
	Result  string `json:"result"`
	Error   error  `json:"-"`
}

type IAlchemyProvider interface {
	/* Send raw transaction */
	Send(method string, params ...string) (string, error)
}

type AlchemyFetchHandler func(AlchemyRequest, RequestConfig) (AlchemyResponse, error)

type BatchAlchemyFetchHandler func([]AlchemyRequest, RequestConfig) ([]AlchemyResponse, error)

type RequestConfig struct {
	Timeout time.Duration
}
