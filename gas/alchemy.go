package gas

import (
	"github.com/poteto-go/go-alchemy-sdk/ether"
	"github.com/poteto-go/go-alchemy-sdk/namespace"
)

type Alchemy struct {
	config AlchemyConfig
	Core   namespace.ICore
}

func NewAlchemy(setting AlchemySetting) Alchemy {
	alchemyConfig := NewAlchemyConfig(setting)
	alchemyProvider := NewAlchemyProvider(alchemyConfig)
	ether := ether.NewEtherApi(alchemyProvider, alchemyConfig.GetUrl())
	coreNamespace := namespace.NewCore(ether)

	return Alchemy{
		config: alchemyConfig,
		Core:   coreNamespace,
	}
}
