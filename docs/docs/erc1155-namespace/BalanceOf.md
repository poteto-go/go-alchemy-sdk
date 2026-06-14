![](https://img.shields.io/badge/go-geth-lightblue)

Get the amount of tokens of the given `tokenId` owned by `account` (ERC-1155 `balanceOf(address,uint256)`).

```go
func BalanceOf(
    contractAddress,
    account string,
    tokenId *big.Int,
) (balance *big.Int, err error)
```

```go
func main() {
    ...
    alchemy := gas.NewAlchemy(setting)
    ...
    balance, err := alchemy.Erc1155.BalanceOf(contractAddress, account, tokenId)
}
```
