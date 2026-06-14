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
	Erc1155    namespace.IErc1155
	ERC20      namespace.IERC20
	StableCoin namespace.IStableCoin
	Debug      namespace.IDebug
	provider   types.IAlchemyProvider
}

func NewAlchemy(setting AlchemySetting) (Alchemy, error) {
	alchemyConfig, err := NewAlchemyConfig(setting)
	if err != nil {
		return Alchemy{}, err
	}

	alchemyProvider := NewAlchemyProvider(alchemyConfig)
	eth := ether.NewEtherApi(
		alchemyProvider,
		alchemyConfig.toEtherApiConfig(),
	)
	alchemyProvider.SetEth(eth)
	coreNamespace := namespace.NewCore(eth)
	transactNamespace := namespace.NewTransactNamespace(eth)
	nftNamespace := namespace.NewNftNamespace(eth)
	erc1155Namespace := namespace.NewErc1155Namespace(eth)
	erc20Namespace := namespace.NewERC20Namespace(eth)
	stableCoinNamespace := namespace.NewStableCoinNamespace(eth)
	debugNamespace := namespace.NewDebugNamespace(eth)

	return Alchemy{
		config:     alchemyConfig,
		Core:       coreNamespace,
		Transact:   transactNamespace,
		Nft:        nftNamespace,
		Erc1155:    erc1155Namespace,
		ERC20:      erc20Namespace,
		StableCoin: stableCoinNamespace,
		Debug:      debugNamespace,
		provider:   alchemyProvider,
	}, nil
}

func (gas *Alchemy) GetProvider() types.IAlchemyProvider {
	return gas.provider
}
