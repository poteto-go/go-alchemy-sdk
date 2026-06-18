![](https://img.shields.io/badge/go-geth-lightblue)

Performs a reverse ENS lookup (address → name) using the default ENS registry for the current network.
Returns `ErrENSNotSupportedOnNetwork` on chains without a known ENS deployment.
Returns `ErrENSNameNotFound` when no reverse record is registered for the address.

```go
func LookupAddress(address string) (string, error)
```

```go
func main() {
	...
	alchemy := gas.NewAlchemy(setting) // network: "eth-mainnet"

	name, _ := alchemy.Core.LookupAddress("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045")
	// name == "vitalik.eth"
}
```
