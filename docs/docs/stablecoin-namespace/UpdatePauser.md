ref: [Wallet-StableCoin-UpdatePauser](../wallet/StableCoin.md#updatepauser--updatepausernowait)

![](https://img.shields.io/badge/go-geth-lightblue)

Update the pauser address (FiatToken/USDC compatibility). Requires the caller to be the current owner.

```go
func UpdatePauser(
    ctx context.Context,
    contractAddress,
    newPauser string,
    gasLimit *uint64,
) (*types.Receipt, error)

func UpdatePauserNoWait(
    contractAddress,
    newPauser string,
    gasLimit *uint64,
) (common.Hash, error)
```

```go
func main() {
    ...
    w, _ := wallet.New("<privateKey>")
    w.Connect(alchemy.GetProvider())
    ...
    receipt, err := w.StableCoin().UpdatePauser(ctx, contractAddress, newPauserAddress, nil)
}
```
