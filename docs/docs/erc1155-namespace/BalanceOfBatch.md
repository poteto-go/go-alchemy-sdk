![](https://img.shields.io/badge/go-geth-lightblue)

Get the balances of multiple `(account, tokenId)` pairs in a single call (ERC-1155 `balanceOfBatch(address[],uint256[])`).

`accounts` and `tokenIds` must have the same length, otherwise `ErrMismatchedArrayLength` is returned.

```go
func BalanceOfBatch(
    contractAddress string,
    accounts []string,
    tokenIds []*big.Int,
) (balances []*big.Int, err error)
```

```go
func main() {
    ...
    alchemy := gas.NewAlchemy(setting)
    ...
    balances, err := alchemy.Erc1155.BalanceOfBatch(
        contractAddress,
        []string{account1, account2},
        []*big.Int{big.NewInt(1), big.NewInt(2)},
    )
}
```
