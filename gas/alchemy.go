package gas

import (
	"github.com/poteto-go/go-alchemy-sdk/ether"
	"github.com/poteto-go/go-alchemy-sdk/namespace"
	"github.com/poteto-go/go-alchemy-sdk/types"
)

type Alchemy struct {
	config   AlchemyConfig
	Core     namespace.ICore
	provider types.IAlchemyProvider
}

func NewAlchemy(setting AlchemySetting) Alchemy {
	alchemyConfig := NewAlchemyConfig(setting)
	alchemyProvider := NewAlchemyProvider(alchemyConfig)
	ether := ether.NewEtherApi(
		alchemyProvider,
		alchemyConfig.toEtherApiConfig(),
	)
	coreNamespace := namespace.NewCore(ether)

	return Alchemy{
		config:   alchemyConfig,
		Core:     coreNamespace,
		provider: alchemyProvider,
	}
}

func (gas *Alchemy) GetProvider() types.IAlchemyProvider {
	return gas.provider
}
