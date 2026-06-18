package constant

// ENSRegistryByNetwork maps Alchemy network names to their ENS registry contract address.
// Only networks with a known ENS deployment are listed.
// Use the address directly for any network not present in this map.
var ENSRegistryByNetwork = map[string]string{
	"eth-mainnet": "0x00000000000C2E074eC69A0dFb2997BA6C7d2e1e",
	"eth-goerli":  "0x00000000000C2E074eC69A0dFb2997BA6C7d2e1e",
	"eth-sepolia": "0x00000000000C2E074eC69A0dFb2997BA6C7d2e1e",
}
