![](https://img.shields.io/badge/go-geth-lightblue)

Return the master minter address of a StableCoin contract (FiatToken/USDC compatibility).

```go
func MasterMinter(
    contractAddress string,
) (common.Address, error)
```

```go
func main() {
    ...
    alchemy := gas.NewAlchemy(setting)
    ...
    masterMinter, err := alchemy.StableCoin.MasterMinter(contractAddress)
}
```
