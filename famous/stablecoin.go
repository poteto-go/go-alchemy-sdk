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

// StableCoinSymbol is a typed constant for well-known stablecoin token symbols.
type StableCoinSymbol string

const (
	USDC StableCoinSymbol = "USDC"
	USDT StableCoinSymbol = "USDT"
	JPYC StableCoinSymbol = "JPYC"
)

// stablecoinAddresses maps network → token symbol → contract address.
// Addresses are pre-parsed at init time to avoid repeated hex conversion on each lookup.
var stablecoinAddresses = map[types.Network]map[StableCoinSymbol]common.Address{
	types.EthMainnet: {
		// https://etherscan.io/token/0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48
		USDC: common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"),
		// https://etherscan.io/token/0xdac17f958d2ee523a2206206994597c13d831ec7
		USDT: common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7"),
		// https://etherscan.io/token/0xe7c3d8c9a439fede00d2600032d5db0be71c3c29
		JPYC: common.HexToAddress("0xE7C3D8C9a439feDe00D2600032D5dB0Be71C3c29"),
	},
	types.PolygonMainnet: {
		// https://polygonscan.com/token/0x3c499c542cef5e3811e1192ce70d8cc03d5c3359 (native USDC)
		USDC: common.HexToAddress("0x3c499c542cEF5E3811e1192ce70d8cC03d5c3359"),
		// https://polygonscan.com/token/0xc2132d05d31c914a87c6611c10748aeb04b58e8f
		USDT: common.HexToAddress("0xc2132D05D31c914a87C6611C10748AEb04B58e8F"),
		// https://polygonscan.com/token/0xe7c3d8c9a439fede00d2600032d5db0be71c3c29
		JPYC: common.HexToAddress("0xE7C3D8C9a439feDe00D2600032D5dB0Be71C3c29"),
	},
	types.PolygonAmoy: {
		// https://amoy.polygonscan.com/token/0x41e94eb019c0762f9bfcf9fb1e58725bfb0e7582 (testnet USDC)
		USDC: common.HexToAddress("0x41E94Eb019C0762f9Bfcf9Fb1E58725BfB0e7582"),
	},
}

// ContractAddress returns the contract address of a well-known stablecoin on the given network.
func ContractAddress(network types.Network, symbol StableCoinSymbol) (common.Address, error) {
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

// SupportedNetworks returns all networks that have at least one registered stablecoin address.
func SupportedNetworks() []types.Network {
	networks := make([]types.Network, 0, len(stablecoinAddresses))
	for n := range stablecoinAddresses {
		networks = append(networks, n)
	}
	return networks
}

// SupportedSymbols returns all stablecoin symbols registered for the given network.
// Returns an empty slice if the network is not supported.
func SupportedSymbols(network types.Network) []StableCoinSymbol {
	networkMap, ok := stablecoinAddresses[network]
	if !ok {
		return nil
	}

	symbols := make([]StableCoinSymbol, 0, len(networkMap))
	for s := range networkMap {
		symbols = append(symbols, s)
	}
	return symbols
}
