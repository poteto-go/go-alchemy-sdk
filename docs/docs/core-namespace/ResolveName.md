![](https://img.shields.io/badge/go-geth-lightblue)

Resolves an ENS name to a lowercase hex address.

## ResolveName

Uses the default ENS registry address for the current network (see `constant.ENSRegistryByNetwork`).
If the input is already a valid hex address it is returned as-is (lowercased).
Returns `ErrENSNotSupportedOnNetwork` on chains without a known ENS deployment.

```go
func ResolveName(name string) (string, error)
```

```go
func main() {
	alchemy := gas.NewAlchemy(setting) // network: "eth-mainnet"

	// Resolve ENS name → address
	addr, _ := alchemy.Core.ResolveName("vitalik.eth")
	// addr == "0xd8da6bf26964af9d7eed9e03e53415d37aa96045"

	// Hex address is returned as-is (lowercased)
	addr, _ = alchemy.Core.ResolveName("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045")
	// addr == "0xd8da6bf26964af9d7eed9e03e53415d37aa96045"
}
```

## ResolveNameBy

Same as `ResolveName` but uses an explicitly provided ENS registry contract address.
Useful when targeting a non-standard ENS deployment or a testnet registry not in the default map.

```go
func ResolveNameBy(registryAddress string, name string) (string, error)
```

```go
func main() {
	alchemy := gas.NewAlchemy(setting)

	const ensRegistry = "0x00000000000C2E074eC69A0dFb2997BA6C7d2e1e"
	addr, _ := alchemy.Core.ResolveNameBy(ensRegistry, "vitalik.eth")
	// addr == "0xd8da6bf26964af9d7eed9e03e53415d37aa96045"
}
```

## Name resolution flow

```
ResolveName("vitalik.eth")
  │
  ├─ 1. Compute namehash("vitalik.eth")   [EIP-137 keccak256 recursion]
  │
  ├─ 2. Call ENS registry: resolver(namehash)
  │       returns the resolver contract address for that name
  │
  └─ 3. Call resolver contract: addr(namehash)
          returns the Ethereum address registered to the name
```
