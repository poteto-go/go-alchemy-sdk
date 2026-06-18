![](https://img.shields.io/badge/go-geth-lightblue)

Resolves an ENS name to a lowercase hex address using an explicitly provided ENS registry contract address.
Useful when targeting a non-standard ENS deployment or a testnet registry.
If the input is already a valid hex address it is returned as-is (lowercased).

```go
func ResolveNameBy(registryAddress string, name string) (string, error)
```

```go
func main() {
	...
	alchemy := gas.NewAlchemy(setting)

	const ensRegistry = "0x00000000000C2E074eC69A0dFb2997BA6C7d2e1e"
	addr, _ := alchemy.Core.ResolveNameBy(ensRegistry, "vitalik.eth")
	// addr == "0xd8da6bf26964af9d7eed9e03e53415d37aa96045"
}
```
