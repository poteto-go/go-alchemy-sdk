ref: [Wallet-StableCoin-Permit](../wallet/StableCoin.md#permit--permitnowait)

![](https://img.shields.io/badge/go-geth-lightblue)

Return the EIP-712 domain separator for the StableCoin contract.

```go
func DomainSeparator(
    contractAddress string,
) ([32]byte, error)
```

```go
func main() {
    ...
    alchemy := gas.NewAlchemy(setting)
    ...
    domainSeparator, err := alchemy.StableCoin.DomainSeparator(contractAddress)
}
```
