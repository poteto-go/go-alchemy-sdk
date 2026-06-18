![](https://img.shields.io/badge/go-geth-lightblue)

Performs a reverse ENS lookup (address → name) using an explicitly provided ENS registry contract address.
Returns `ErrENSNameNotFound` when no reverse record is registered for the address.

```go
func LookupAddressBy(registryAddress string, address string) (string, error)
```

```go
func main() {
	...
	alchemy := gas.NewAlchemy(setting)

	const ensRegistry = "0x00000000000C2E074eC69A0dFb2997BA6C7d2e1e"
	name, _ := alchemy.Core.LookupAddressBy(ensRegistry, "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045")
	// name == "vitalik.eth"
}
```
