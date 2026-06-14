![](https://img.shields.io/badge/go-geth-lightblue)

Get the metadata URI for the given `tokenId` (ERC-1155 `uri(uint256)`).

ERC-1155 uses a single URI template that may contain an `{id}` placeholder shared by all tokens, rather than a per-token URI.

```go
func Uri(
    contractAddress string,
    tokenId *big.Int,
) (uri string, err error)
```

```go
func main() {
    ...
    alchemy := gas.NewAlchemy(setting)
    ...
    uri, err := alchemy.Erc1155.Uri(contractAddress, tokenId)
}
```
