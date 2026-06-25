package types

import (
	"context"
	"errors"
	"net/http"

	"github.com/ethereum/go-ethereum"
)

var ErrNoResultFound = errors.New("no result found")

type AlchemyRequest struct {
	Request *http.Request
}

const DefaultMaxResponseBytes int64 = 32 * 1024 * 1024 // 32 MiB

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
	Network() Network

	/* Send raw transaction */
	Send(method string, params RequestArgs) (any, error)
}

// ISubscribeProvider is the focused subscribe surface implemented only by the
// ws provider. It is kept separate from IAlchemyProvider so HTTP providers (which
// cannot push) are not forced to implement a subscribe path. The subscription
// namespace type-asserts a provider to this interface to open eth_subscribe
// streams over the persistent ws socket.
type ISubscribeProvider interface {
	// Subscribe opens an eth_subscribe stream. params[0] is the subscription
	// name (e.g. "newHeads", "logs", "alchemy_minedTransactions"); notifications
	// are delivered on channel until the returned subscription is unsubscribed.
	Subscribe(ctx context.Context, channel any, params ...any) (ethereum.Subscription, error)
}

type AlchemyFetchHandler func(*http.Client, AlchemyRequest, []byte) (AlchemyResponse, error)

type BatchAlchemyFetchHandler func(*http.Client, []AlchemyRequest, [][]byte) ([]AlchemyResponse, error)
