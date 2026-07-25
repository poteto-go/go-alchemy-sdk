package gas

import (
	"github.com/poteto-go/go-alchemy-sdk/ether"
	"github.com/poteto-go/go-alchemy-sdk/namespace"
	"github.com/poteto-go/go-alchemy-sdk/types"
)

type Alchemy struct {
	config     AlchemyConfig
	Core       namespace.ICore
	Transact   namespace.ITransact
	Nft        namespace.INft
	ERC1155    namespace.IErc1155
	ERC20      namespace.IERC20
	StableCoin namespace.IStableCoin
	Debug      namespace.IDebug
	WS         namespace.IWS
	provider   types.IAlchemyProvider
}

func NewAlchemy(setting AlchemySetting) (Alchemy, error) {
	alchemyConfig, err := NewAlchemyConfig(setting)
	if err != nil {
		return Alchemy{}, err
	}

	alchemyProvider := newProvider(alchemyConfig)
	eth := ether.NewEtherApi(
		alchemyProvider,
		alchemyConfig.toEtherApiConfig(),
	)

	// WS is left nil: subscriptions need a websocket endpoint, use NewWsAlchemy.
	return newAlchemy(alchemyConfig, alchemyProvider, eth), nil
}

func NewWsAlchemy(setting AlchemySetting) (Alchemy, error) {
	alchemyConfig, err := NewAlchemyConfig(setting)
	if err != nil {
		return Alchemy{}, err
	}

	alchemyProvider := newProvider(alchemyConfig)
	eth := ether.NewWsEtherApi(
		alchemyProvider,
		alchemyConfig.toEtherApiConfig(),
	)

	alchemy := newAlchemy(alchemyConfig, alchemyProvider, eth)
	alchemy.WS = namespace.NewWSNamespace(eth)
	return alchemy, nil
}

// newAlchemy wires every transport-agnostic namespace onto eth. Transport
// specific extras (WS) are added by the caller.
func newAlchemy(
	config AlchemyConfig,
	provider types.IAlchemyProvider,
	eth types.EtherApi,
) Alchemy {
	provider.SetEth(eth)

	return Alchemy{
		config:     config,
		Core:       namespace.NewCore(eth),
		Transact:   namespace.NewTransactNamespace(eth),
		Nft:        namespace.NewNftNamespace(eth),
		ERC1155:    namespace.NewErc1155Namespace(eth),
		ERC20:      namespace.NewERC20Namespace(eth),
		StableCoin: namespace.NewStableCoinNamespace(eth),
		Debug:      namespace.NewDebugNamespace(eth),
		provider:   provider,
	}
}

func (gas *Alchemy) GetProvider() types.IAlchemyProvider {
	return gas.provider
}

// newProvider picks the transport-appropriate provider: ws/wss endpoints route
// over the persistent websocket socket, everything else over HTTP.
func newProvider(config AlchemyConfig) types.IAlchemyProvider {
	if config.isWebSocket() {
		return NewWsAlchemyProvider(config)
	}
	return NewAlchemyProvider(config)
}
