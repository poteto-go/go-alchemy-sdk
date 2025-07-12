package types

import (
	"errors"
	"net/http"
	"time"
)

var ErrNoResultFound = errors.New("no result found")

type AlchemyRequest[T string | TransactionRequest | Filter] struct {
	Request *http.Request
	Body    AlchemyRequestBody[T]
}

type AlchemyRequestBody[T string | TransactionRequest | Filter] struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  []T    `json:"params"`
	Id      int    `json:"id"`
}

type AlchemyResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      int    `json:"id"`
	Result  any    `json:"result"`
	Error   error  `json:"-"`
}

type IAlchemyProvider interface {
	/* Send raw transaction */
	Send(method string, params ...string) (any, error)

	/* Send transaction */
	SendTransaction(method string, params ...TransactionRequest) (any, error)

	/* Send filter */
	SendFilter(method string, params ...Filter) (any, error)
}

type AlchemyFetchHandler[T string | TransactionRequest | Filter] func(AlchemyRequest[T], RequestConfig) (AlchemyResponse, error)

type BatchAlchemyFetchHandler[T string | TransactionRequest | Filter] func([]AlchemyRequest[T], RequestConfig) ([]AlchemyResponse, error)

type RequestConfig struct {
	Timeout time.Duration
}
