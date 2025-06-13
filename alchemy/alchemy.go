package alchemy

import (
	"github.com/poteto-go/go-alchemy-sdk/namespace"
)

type Alchemy struct {
	config AlchemyConfig
	Core   namespace.ICore
}

func NewAlchemy(setting AlchemySetting) Alchemy {
	alchemyConfig := NewAlchemyConfig(setting)
	alchemyProvider := NewAlchemyProvider(alchemyConfig)
	coreNamespace := namespace.NewCore(alchemyProvider)

	return Alchemy{
		config: alchemyConfig,
		Core:   coreNamespace,
	}
}
