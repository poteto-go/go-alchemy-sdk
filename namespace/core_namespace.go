package namespace

import (
	"math/big"

	"github.com/poteto-go/go-alchemy-sdk/ether"
)

type ICore interface {
	/* get  the number of the most recent block. */
	GetBlockNumber() (int, error)

	/* Returns the best guess of the current gas price to use in a transaction. */
	GetGasPrice() (int, error)

	/* Returns the balance of a given address as of the provided block. */
	GetBalance(address string, blockTag string) (*big.Int, error)
}

type Core struct {
	ether ether.EtherApi
}

func NewCore(ether ether.EtherApi) ICore {
	return &Core{
		ether: ether,
	}
}

func (c *Core) GetBlockNumber() (int, error) {
	blockNumber, err := c.ether.GetBlockNumber()
	if err != nil {
		return 0, err
	}
	return blockNumber, nil
}

func (c *Core) GetGasPrice() (int, error) {
	price, err := c.ether.GetGasPrice()
	if err != nil {
		return 0, err
	}
	return price, nil
}

func (c *Core) GetBalance(address string, blockTag string) (*big.Int, error) {
	balance, err := c.ether.GetBalance(address, blockTag)
	if err != nil {
		return big.NewInt(0), err
	}
	return balance, nil
}
