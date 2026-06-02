ref: [Wallet-StableCoin-Currency](../wallet/StableCoin.md#currency)

![](https://img.shields.io/badge/go-geth-lightblue)

Get the currency identifier of the token (e.g. `"USD"`). FiatToken/USDC compatibility.

```go
func Currency(
    contractAddress string,
) (string, error)
```

```go
func main() {
    ...
    alchemy := gas.NewAlchemy(setting)
    ...
    currency, err := w.StableCoin().Currency(contractAddress)
}
```
