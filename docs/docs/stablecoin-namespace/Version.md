ref: [Wallet-StableCoin-Version](../wallet/StableCoin.md#version)

![](https://img.shields.io/badge/go-geth-lightblue)

Get the contract version string. FiatToken/USDC compatibility.

```go
func Version(
    contractAddress string,
) (string, error)
```

```go
func main() {
    ...
    alchemy := gas.NewAlchemy(setting)
    ...
    version, err := w.StableCoin().Version(contractAddress)
}
```
