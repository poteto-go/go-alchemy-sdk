package namespace

import "github.com/poteto-go/go-alchemy-sdk/types"

type IStableCoin interface {
	IERC20
}

type stableCoin struct {
	*ERC20
}

func NewStableCoinNamespace(ether types.EtherApi) IStableCoin {
	return &stableCoin{ERC20: &ERC20{ether: ether}}
}
