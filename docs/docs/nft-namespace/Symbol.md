![](https://img.shields.io/badge/go-geth-lightblue)

Get the symbol of the NFT collection (ERC-721 `symbol()`).

```go
func Symbol(
    contractAddress string,
) (symbol string, err error)
```

```go
func main() {
    ...
    alchemy := gas.NewAlchemy(setting)
    ...
    symbol, err := alchemy.Nft.Symbol(contractAddress)
}
```
