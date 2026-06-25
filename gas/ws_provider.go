package gas

import (
	"context"
	"net/http"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/poteto-go/go-alchemy-sdk/constant"
	"github.com/poteto-go/go-alchemy-sdk/types"
)

// WsAlchemyProvider routes JSON-RPC requests and eth_subscribe streams over the
// persistent websocket rpc.Client owned by Ether (see ether.setWsEthClient).
// Unlike AlchemyProvider it does not dial its own transport: reusing Ether's
// socket keeps a single long-lived ws connection per Alchemy instance.
type WsAlchemyProvider struct {
	config AlchemyConfig
	eth    types.EtherApi
}

// compile-time assertions: WsAlchemyProvider serves both the base provider
// surface and the focused subscribe path.
var (
	_ types.IAlchemyProvider   = (*WsAlchemyProvider)(nil)
	_ types.ISubscribeProvider = (*WsAlchemyProvider)(nil)
)

func NewWsAlchemyProvider(config AlchemyConfig) types.IAlchemyProvider {
	return &WsAlchemyProvider{config: config}
}

func (provider *WsAlchemyProvider) SetEth(eth types.EtherApi) {
	provider.eth = eth
}

func (provider *WsAlchemyProvider) Eth() types.EtherApi {
	return provider.eth
}

func (provider *WsAlchemyProvider) CustomHeaders() []http.Header {
	return provider.config.customHeaders
}

func (provider *WsAlchemyProvider) Network() types.Network {
	return provider.config.network
}

// Send dispatches a single JSON-RPC call over the ws socket, mirroring the
// HTTP provider's contract: it returns the raw result (map/slice/string) or
// ErrResultIsNil when the node answers with a null result.
func (provider *WsAlchemyProvider) Send(method string, params types.RequestArgs) (any, error) {
	client, err := provider.rpcClient()
	if err != nil {
		return nil, err
	}

	ctx, cancel := provider.requestContext()
	defer cancel()

	var result any
	if err := client.CallContext(ctx, &result, method, params...); err != nil {
		return nil, err
	}
	if result == nil {
		return nil, constant.ErrResultIsNil
	}
	return result, nil
}

// Subscribe opens an eth_subscribe stream over the ws socket. The returned
// *rpc.ClientSubscription satisfies ethereum.Subscription.
func (provider *WsAlchemyProvider) Subscribe(ctx context.Context, channel any, params ...any) (ethereum.Subscription, error) {
	client, err := provider.rpcClient()
	if err != nil {
		return nil, err
	}
	return client.EthSubscribe(ctx, channel, params...)
}

// rpcClient establishes (or reuses) Ether's persistent ws client and returns the
// underlying geth rpc.Client used for both calls and subscriptions.
func (provider *WsAlchemyProvider) rpcClient() (*rpc.Client, error) {
	if provider.eth == nil {
		return nil, constant.ErrProviderEthNotSet
	}
	if err := provider.eth.SetEthClient(); err != nil {
		return nil, err
	}
	client, ok := provider.eth.Client().(*ethclient.Client)
	if !ok {
		return nil, constant.ErrUnSupportSimulatedMethod
	}
	return client.Client(), nil
}

func (provider *WsAlchemyProvider) requestContext() (context.Context, context.CancelFunc) {
	if provider.config.requestTimeout > 0 {
		return context.WithTimeout(context.Background(), provider.config.requestTimeout)
	}
	return context.WithCancel(context.Background())
}
