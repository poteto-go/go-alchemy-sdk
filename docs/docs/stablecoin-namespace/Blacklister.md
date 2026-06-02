![](https://img.shields.io/badge/go-geth-lightblue)

Return the blacklister address of a StableCoin contract (FiatToken/USDC compatibility).

```go
func Blacklister(
    contractAddress string,
) (common.Address, error)
```

```go
func main() {
    ...
    alchemy := gas.NewAlchemy(setting)
    ...
    blacklister, err := alchemy.StableCoin.Blacklister(contractAddress)
}
```
