package lib

import (
	"math/big"

	"github.com/goccy/go-json"
)

type INetwork interface {
	ToJson() ([]byte, error)

	/*
		Returns true if %%other%% matches this network.
		Any chain ID must match, and if no chain ID is present, the name must match.
		This method does not currently check for additional properties,
		such as ENS address or plug-in compatibility.
		TODO: Networkish
		https://github.com/ethers-io/ethers.js/blob/main/src.ts/providers/network.ts#L28
	*/
	Matches(other INetwork) bool

	Name() string
	ChainId() *big.Int

	SetName(name string)
	SetChainId(chainId *big.Int)
}

type Network struct {
	name    string
	chainId *big.Int
}

func NewNetwork(name string, chainId *big.Int) INetwork {
	return &Network{
		name:    name,
		chainId: chainId,
	}
}

func (n *Network) ToJson() ([]byte, error) {
	type NetworkForJson struct {
		Name    string   `json:"name"`
		ChainId *big.Int `json:"chainId"`
	}

	nj := NetworkForJson{
		Name:    n.name,
		ChainId: n.chainId,
	}

	return json.Marshal(nj)
}

// TODO: Networkish
func (n *Network) Matches(other INetwork) bool {
	if other == nil {
		return false
	}

	if other.ChainId() != nil {
		return n.chainId.Cmp(other.ChainId()) == 0
	}

	if other.Name() != "" {
		return n.name == other.Name()
	}

	return false
}

func (n *Network) Name() string {
	return n.name
}

func (n *Network) ChainId() *big.Int {
	return n.chainId
}

func (n *Network) SetName(name string) {
	n.name = name
}

func (n *Network) SetChainId(chainId *big.Int) {
	n.chainId = chainId
}
