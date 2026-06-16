package famous

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/poteto-go/go-alchemy-sdk/types"
)

var (
	ErrNotSupportedNftNetwork = errors.New("not supported network for this nft collection")
	ErrNotSupportedNftSymbol  = errors.New("not supported nft collection symbol")
)

// NftSymbol is a typed constant for well-known NFT collection symbols.
type NftSymbol string

const (
	BAYC          NftSymbol = "BAYC"          // Bored Ape Yacht Club
	MAYC          NftSymbol = "MAYC"          // Mutant Ape Yacht Club
	CryptoPunks   NftSymbol = "CryptoPunks"   // CryptoPunks
	Azuki         NftSymbol = "Azuki"         // Azuki
	Doodles       NftSymbol = "Doodles"       // Doodles
	PudgyPenguins NftSymbol = "PudgyPenguins" // Pudgy Penguins
)

// nftAddresses maps network → collection symbol → contract address.
// Addresses are pre-parsed at init time to avoid repeated hex conversion on each lookup.
var nftAddresses = map[types.Network]map[NftSymbol]common.Address{
	types.EthMainnet: {
		// https://etherscan.io/address/0xbc4ca0eda7647a8ab7c2061c2e118a18a936f13d
		BAYC: common.HexToAddress("0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D"),
		// https://etherscan.io/address/0x60e4d786628fea6478f785a6d7e704777c86a7c6
		MAYC: common.HexToAddress("0x60E4d786628Fea6478F785A6d7e704777c86a7c6"),
		// https://etherscan.io/address/0xb47e3cd837ddf8e4c57f05d70ab865de6e193bbb
		CryptoPunks: common.HexToAddress("0xb47e3cd837dDF8e4c57F05d70Ab865de6e193BBB"),
		// https://etherscan.io/address/0xed5af388653567af2f388e6224dc7c4b3241c544
		Azuki: common.HexToAddress("0xED5AF388653567Af2F388E6224dC7C4b3241C544"),
		// https://etherscan.io/address/0x8a90cab2b38dba80c64b7734e58ee1db38b8992e
		Doodles: common.HexToAddress("0x8a90CAb2b38dba80c64b7734e58Ee1dB38B8992e"),
		// https://etherscan.io/address/0xbd3531da5cf5857e7cfaa92426877b022e612cf8
		PudgyPenguins: common.HexToAddress("0xBd3531dA5CF5857e7CfAA92426877b022e612cf8"),
	},
}

// NftContractAddress returns the contract address of a well-known NFT collection on the given network.
func NftContractAddress(network types.Network, symbol NftSymbol) (common.Address, error) {
	networkMap, ok := nftAddresses[network]
	if !ok {
		return common.Address{}, ErrNotSupportedNftNetwork
	}

	addr, ok := networkMap[symbol]
	if !ok {
		return common.Address{}, ErrNotSupportedNftSymbol
	}

	return addr, nil
}

// NftSupportedNetworks returns all networks that have at least one registered NFT collection address.
func NftSupportedNetworks() []types.Network {
	networks := make([]types.Network, 0, len(nftAddresses))
	for n := range nftAddresses {
		networks = append(networks, n)
	}
	return networks
}

// NftSupportedSymbols returns all NFT collection symbols registered for the given network.
// Returns an empty slice if the network is not supported.
func NftSupportedSymbols(network types.Network) []NftSymbol {
	networkMap, ok := nftAddresses[network]
	if !ok {
		return nil
	}

	symbols := make([]NftSymbol, 0, len(networkMap))
	for s := range networkMap {
		symbols = append(symbols, s)
	}
	return symbols
}
