![](https://img.shields.io/badge/go-geth-lightblue)

Resolves an ENS name to a lowercase hex address using the default ENS registry for the current network.
If the input is already a valid hex address it is returned as-is (lowercased).
Returns `ErrENSNotSupportedOnNetwork` on chains without a known ENS deployment.

```go
func ResolveName(name string) (string, error)
```

```go
func main() {
	...
	alchemy := gas.NewAlchemy(setting) // network: "eth-mainnet"

	// Resolve ENS name → address
	addr, _ := alchemy.Core.ResolveName("vitalik.eth")
	// addr == "0xd8da6bf26964af9d7eed9e03e53415d37aa96045"

	// Hex address is returned as-is (lowercased)
	addr, _ = alchemy.Core.ResolveName("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045")
	// addr == "0xd8da6bf26964af9d7eed9e03e53415d37aa96045"
}
```
