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

const DefaultMaxResponseBytes int64 = 32 * 1024 * 1024 // 32 MiB

type RequestConfig struct {
	Timeout          time.Duration
	MaxResponseBytes int64
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
	CustomHeaders() []http.Header

	/* Send raw transaction */
	Send(method string, params RequestArgs) (any, error)
}

type AlchemyFetchHandler func(AlchemyRequest, RequestConfig, []byte) (AlchemyResponse, error)

type BatchAlchemyFetchHandler func([]AlchemyRequest, RequestConfig, [][]byte) ([]AlchemyResponse, error)
