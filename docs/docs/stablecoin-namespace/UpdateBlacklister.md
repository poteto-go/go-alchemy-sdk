ref: [Wallet-StableCoin-UpdateBlacklister](../wallet/StableCoin.md#updateblacklister--updateblacklisternowait)

![](https://img.shields.io/badge/go-geth-lightblue)

Update the blacklister address (FiatToken/USDC compatibility). Requires the caller to be the current owner.

```go
func UpdateBlacklister(
    ctx context.Context,
    contractAddress,
    newBlacklister string,
    gasLimit *uint64,
) (*types.Receipt, error)

func UpdateBlacklisterNoWait(
    contractAddress,
    newBlacklister string,
    gasLimit *uint64,
) (common.Hash, error)
```

```go
func main() {
    ...
    w, _ := wallet.New("<privateKey>")
    w.Connect(alchemy.GetProvider())
    ...
    receipt, err := w.StableCoin().UpdateBlacklister(ctx, contractAddress, newBlacklisterAddress, nil)
}
```
