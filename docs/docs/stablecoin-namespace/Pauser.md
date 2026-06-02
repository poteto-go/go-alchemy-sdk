![](https://img.shields.io/badge/go-geth-lightblue)

Return the pauser address of a StableCoin contract (FiatToken/USDC compatibility).

```go
func Pauser(
    contractAddress string,
) (common.Address, error)
```

```go
func main() {
    ...
    alchemy := gas.NewAlchemy(setting)
    ...
    pauser, err := alchemy.StableCoin.Pauser(contractAddress)
}
```
