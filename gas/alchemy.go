package gas

import (
	"github.com/poteto-go/go-alchemy-sdk/ether"
	"github.com/poteto-go/go-alchemy-sdk/namespace"
	"github.com/poteto-go/go-alchemy-sdk/types"
)

type Alchemy struct {
	config   AlchemyConfig
	Core     namespace.ICore
	Transact namespace.ITransact
	Nft      namespace.INft
	ERC20    namespace.IERC20
	provider types.IAlchemyProvider
}

func NewAlchemy(setting AlchemySetting) (Alchemy, error) {
	alchemyConfig, err := NewAlchemyConfig(setting)
	if err != nil {
		return Alchemy{}, err
	}

	alchemyProvider := NewAlchemyProvider(alchemyConfig)
	ether := ether.NewEtherApi(
		alchemyProvider,
		alchemyConfig.toEtherApiConfig(),
	)
	alchemyProvider.SetEth(ether)
	coreNamespace := namespace.NewCore(ether)
	transactNamespace := namespace.NewTransactNamespace(ether)
	nftNamespace := namespace.NewNftNamespace(ether)
	erc20Namespace := namespace.NewERC20Namespace(ether)

	return Alchemy{
		config:   alchemyConfig,
		Core:     coreNamespace,
		Transact: transactNamespace,
		Nft:      nftNamespace,
		ERC20:    erc20Namespace,
		provider: alchemyProvider,
	}, nil
}

func (gas *Alchemy) GetProvider() types.IAlchemyProvider {
	return gas.provider
}
