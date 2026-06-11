![](https://img.shields.io/badge/go-geth-lightblue)

Get the address approved to transfer the NFT with the given tokenId (ERC-721 `getApproved(uint256)`).

```go
func GetApproved(
    contractAddress string,
    tokenId *big.Int,
) (approved string, err error)
```

```go
func main() {
    ...
    alchemy := gas.NewAlchemy(setting)
    ...
    approved, err := alchemy.Nft.GetApproved(contractAddress, tokenId)
}
```
