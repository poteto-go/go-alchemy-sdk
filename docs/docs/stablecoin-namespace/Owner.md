ref: [Wallet-StableCoin-TransferOwnership](../wallet/StableCoin.md#transferownership--transferownershipnowait)

![](https://img.shields.io/badge/go-geth-lightblue)

Return the current owner address of a StableCoin contract (FiatToken/USDC compatibility).

```go
func Owner(
    contractAddress string,
) (common.Address, error)
```

```go
func main() {
    ...
    alchemy := gas.NewAlchemy(setting)
    ...
    owner, err := alchemy.StableCoin.Owner(contractAddress)
}
```
