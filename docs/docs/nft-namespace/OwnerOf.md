![](https://img.shields.io/badge/go-geth-lightblue)

Get the owner address of the NFT with the given tokenId (ERC-721 `ownerOf(uint256)`).

```go
func OwnerOf(
    contractAddress string,
    tokenId *big.Int,
) (owner string, err error)
```

```go
func main() {
    ...
    alchemy := gas.NewAlchemy(setting)
    ...
    owner, err := alchemy.Nft.OwnerOf(contractAddress, tokenId)
}
```
