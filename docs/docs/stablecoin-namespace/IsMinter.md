ref: [Wallet-StableCoin-IsMinter](../wallet/StableCoin.md#isminter)

![](https://img.shields.io/badge/go-geth-lightblue)

Check whether an address is a configured minter on a StableCoin contract (FiatToken/USDC compatibility).

```go
func IsMinter(
    contractAddress,
    address string,
) (bool, error)
```

```go
func main() {
    ...
    alchemy := gas.NewAlchemy(setting)
    ...
    isMinter, err := alchemy.StableCoin.IsMinter(contractAddress, address)
}
```
