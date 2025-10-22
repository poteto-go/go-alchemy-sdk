package types

import (
	"errors"
	"net/http"
	"time"
)

var ErrNoResultFound = errors.New("no result found")

type AlchemyRequest struct {
	Request *http.Request
}

type RequestConfig struct {
	Timeout time.Duration
}

// for json marshal
type RequestArgs = []any

type AlchemyRequestBody struct {
	Jsonrpc string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  RequestArgs `json:"params,omitzero"`
	Id      int         `json:"id"`
}

type AlchemyResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      int    `json:"id"`
	Result  any    `json:"result"`
	Error   error  `json:"-"`
}

type IAlchemyProvider interface {
	SetEth(eth EtherApi)
	Eth() EtherApi

	/* Send raw transaction */
	Send(method string, params RequestArgs) (any, error)
}

type AlchemyFetchHandler func(AlchemyRequest, RequestConfig, []byte) (AlchemyResponse, error)

type BatchAlchemyFetchHandler func([]AlchemyRequest, RequestConfig, [][]byte) ([]AlchemyResponse, error)
