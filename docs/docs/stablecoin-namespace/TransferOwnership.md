ref: [Wallet-StableCoin-TransferOwnership](../wallet/StableCoin.md#transferownership--transferownershipnowait)

![](https://img.shields.io/badge/go-geth-lightblue)

Transfer contract ownership to a new address (FiatToken/USDC compatibility). Requires the caller to be the current owner.

```go
func TransferOwnership(
    ctx context.Context,
    contractAddress,
    newOwner string,
    gasLimit *uint64,
) (*types.Receipt, error)

func TransferOwnershipNoWait(
    contractAddress,
    newOwner string,
    gasLimit *uint64,
) (common.Hash, error)
```

```go
func main() {
    ...
    w, _ := wallet.New("<privateKey>")
    w.Connect(alchemy.GetProvider())
    ...
    receipt, err := w.StableCoin().TransferOwnership(ctx, contractAddress, newOwnerAddress, nil)
}
```
