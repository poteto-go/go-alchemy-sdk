![](https://img.shields.io/badge/go-geth-lightblue)

Get the name of the NFT collection (ERC-721 `name()`).

```go
func Name(
    contractAddress string,
) (name string, err error)
```

```go
func main() {
    ...
    alchemy := gas.NewAlchemy(setting)
    ...
    name, err := alchemy.Nft.Name(contractAddress)
}
```
