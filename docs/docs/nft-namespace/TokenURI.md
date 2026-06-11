![](https://img.shields.io/badge/go-geth-lightblue)

Get the metadata URI of the NFT with the given tokenId (ERC-721 `tokenURI(uint256)`).

```go
func TokenURI(
    contractAddress string,
    tokenId *big.Int,
) (uri string, err error)
```

```go
func main() {
    ...
    alchemy := gas.NewAlchemy(setting)
    ...
    uri, err := alchemy.Nft.TokenURI(contractAddress, tokenId)
}
```
