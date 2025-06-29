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

	/*
		Returns the contract code of the provided address at the block.
		If there is no contract deployed, the result is 0x.
	*/
	GetCode(address, blockTag string) (string, error)

	/* Checks if the provided address is a smart contract. */
	IsContractAddress(address string) bool
}

type Core struct {
	ether ether.EtherApi
}

func NewCore(ether ether.EtherApi) ICore {
	return &Core{
		ether: ether,
	}
}

/* get  the number of the most recent block. */
func (c *Core) GetBlockNumber() (int, error) {
	blockNumber, err := c.ether.GetBlockNumber()
	if err != nil {
		return 0, err
	}
	return blockNumber, nil
}

/* Returns the best guess of the current gas price to use in a transaction. */
func (c *Core) GetGasPrice() (int, error) {
	price, err := c.ether.GetGasPrice()
	if err != nil {
		return 0, err
	}
	return price, nil
}

/* Returns the balance of a given address as of the provided block. */
func (c *Core) GetBalance(address string, blockTag string) (*big.Int, error) {
	balance, err := c.ether.GetBalance(address, blockTag)
	if err != nil {
		return big.NewInt(0), err
	}
	return balance, nil
}

/*
Returns the contract code of the provided address at the block.
If there is no contract deployed, the result is 0x.
*/
func (c *Core) GetCode(address, blockTag string) (string, error) {
	hexCode, err := c.ether.GetCode(address, blockTag)
	if err != nil {
		return "", err
	}
	return hexCode, nil
}

/* Checks if the provided address is a smart contract. */
func (c *Core) IsContractAddress(address string) bool {
	hexCode, err := c.GetCode(address, "latest")
	if err != nil {
		return false
	}

	return hexCode != "0x"
}
