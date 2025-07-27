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

type AlchemyRequestBody[
	T string | TransactionRequest | Filter | TransactionRequestWithBlockTag,
] struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  []T    `json:"params,omitempty"`
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

	/* Send transaction w blockTag */
	SendTransactionWithBlockTag(method string, params ...TransactionRequestWithBlockTag) (any, error)

	/* Send filter */
	SendFilter(method string, params ...Filter) (any, error)
}

type AlchemyFetchHandler func(AlchemyRequest, RequestConfig, []byte) (AlchemyResponse, error)

type BatchAlchemyFetchHandler func([]AlchemyRequest, RequestConfig, [][]byte) ([]AlchemyResponse, error)

type RequestConfig struct {
	Timeout time.Duration
}
