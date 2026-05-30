package famous

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/poteto-go/go-alchemy-sdk/types"
)

var (
	ErrNotSupportedNetwork = errors.New("not supported network for this stablecoin")
	ErrNotSupportedSymbol  = errors.New("not supported stablecoin symbol")
)

// stablecoinAddresses maps network → token symbol → contract address.
// Addresses are pre-parsed at init time to avoid repeated hex conversion on each lookup.
var stablecoinAddresses = map[types.Network]map[string]common.Address{
	types.EthMainnet: {
		"USDC": common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"),
		"USDT": common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7"),
		"JPYC": common.HexToAddress("0x431D5dfF03120AFA4bDf332c61A6e1766eF37BF9"),
	},
	types.PolygonMainnet: {
		"USDC": common.HexToAddress("0x3c499c542cEF5E3811e1192ce70d8cC03d5c3359"),
		"USDT": common.HexToAddress("0xc2132D05D31c914a87C6611C10748AEb04B58e8F"),
		"JPYC": common.HexToAddress("0x6AE7Dfc73E0dDE2aa99ac063DcF7e8A63265108c"),
	},
	types.PolygonAmoy: {
		"USDC": common.HexToAddress("0x41E94Eb019C0762f9Bfcf9Fb1E58725BfB0e7582"),
	},
}

// ContractAddress returns the contract address of a well-known stablecoin on the given network.
func ContractAddress(network types.Network, symbol string) (common.Address, error) {
	networkMap, ok := stablecoinAddresses[network]
	if !ok {
		return common.Address{}, ErrNotSupportedNetwork
	}

	addr, ok := networkMap[symbol]
	if !ok {
		return common.Address{}, ErrNotSupportedSymbol
	}

	return addr, nil
}
