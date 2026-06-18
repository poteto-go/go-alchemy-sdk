![](https://img.shields.io/badge/go-geth-lightblue)

Performs a reverse ENS lookup: resolves an Ethereum address to its registered ENS name.

## LookupAddress

Uses the default ENS registry address for the current network (see `constant.ENSRegistryByNetwork`).
Returns `ErrENSNotSupportedOnNetwork` on chains without a known ENS deployment.
Returns `ErrENSNameNotFound` when no reverse record is registered for the address.

```go
func LookupAddress(address string) (string, error)
```

```go
func main() {
	alchemy := gas.NewAlchemy(setting) // network: "eth-mainnet"

	name, _ := alchemy.Core.LookupAddress("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045")
	// name == "vitalik.eth"
}
```

## LookupAddressBy

Same as `LookupAddress` but uses an explicitly provided ENS registry contract address.
Useful when targeting a non-standard ENS deployment or a testnet registry not in the default map.

```go
func LookupAddressBy(registryAddress string, address string) (string, error)
```

```go
func main() {
	alchemy := gas.NewAlchemy(setting)

	const ensRegistry = "0x00000000000C2E074eC69A0dFb2997BA6C7d2e1e"
	name, _ := alchemy.Core.LookupAddressBy(ensRegistry, "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045")
	// name == "vitalik.eth"
}
```

## Reverse resolution flow

ENS reverse lookup uses the `.addr.reverse` subdomain convention (EIP-181).

```
LookupAddress("0xd8da...045")
  │
  ├─ 1. Build reverse name: "d8da...045.addr.reverse"
  │
  ├─ 2. Compute namehash of the reverse name  [EIP-137]
  │
  ├─ 3. Call ENS registry: resolver(namehash)
  │       returns the reverse resolver contract address
  │
  └─ 4. Call reverse resolver: name(namehash)
          returns the ENS name registered to that address
```
