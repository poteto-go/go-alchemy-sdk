ref: [Wallet-StableCoin-Paused](../wallet/StableCoin.md#paused)

![](https://img.shields.io/badge/go-geth-lightblue)

Check whether the StableCoin contract is currently paused (FiatToken/USDC compatibility).

```go
func Paused(
    contractAddress string,
) (bool, error)
```

```go
func main() {
    ...
    alchemy := gas.NewAlchemy(setting)
    ...
    paused, err := alchemy.StableCoin.Paused(contractAddress)
}
```
