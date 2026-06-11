![](https://img.shields.io/badge/go-geth-lightblue)

Check whether an operator is approved to manage all NFTs of the given owner (ERC-721 `isApprovedForAll(address,address)`).

```go
func IsApprovedForAll(
    contractAddress,
    owner,
    operator string,
) (approved bool, err error)
```

```go
func main() {
    ...
    alchemy := gas.NewAlchemy(setting)
    ...
    approved, err := alchemy.Nft.IsApprovedForAll(contractAddress, owner, operator)
}
```
