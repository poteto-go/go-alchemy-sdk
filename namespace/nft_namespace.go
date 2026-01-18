package namespace

import "github.com/poteto-go/go-alchemy-sdk/types"

type INft interface{}

type Nft struct {
	ether types.EtherApi
}

func NewNftNamespace(ether types.EtherApi) INft {
	return &Nft{
		ether: ether,
	}
}
