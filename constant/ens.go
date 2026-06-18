package constant

import "github.com/poteto-go/go-alchemy-sdk/types"

// ENSRegistryByNetwork maps network to the ENS registry contract address.
// Only networks with a known ENS deployment are listed.
var ENSRegistryByNetwork = map[types.Network]string{
	types.EthMainnet: "0x00000000000C2E074eC69A0dFb2997BA6C7d2e1e",
	types.EthGoerli:  "0x00000000000C2E074eC69A0dFb2997BA6C7d2e1e",
	types.EthSepolia: "0x00000000000C2E074eC69A0dFb2997BA6C7d2e1e",
}
