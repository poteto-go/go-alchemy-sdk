package namespace

import (
	"github.com/poteto-go/go-alchemy-sdk/types"
)

type ICore interface {
	GetBlockNumber() (int, error)
}

type Core struct {
	provider types.IAlchemyProvider
}

func NewCore(provider types.IAlchemyProvider) ICore {
	return &Core{
		provider: provider,
	}
}

func (c *Core) GetBlockNumber() (int, error) {
	blockNumber, err := c.provider.GetBlockNumber()
	if err != nil {
		return 0, err
	}
	return blockNumber, nil
}
