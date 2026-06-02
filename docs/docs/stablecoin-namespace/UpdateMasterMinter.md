ref: [Wallet-StableCoin-UpdateMasterMinter](../wallet/StableCoin.md#updatemasterminter--updatemasterminternowait)

![](https://img.shields.io/badge/go-geth-lightblue)

Update the master minter address (FiatToken/USDC compatibility). Requires the caller to be the current owner.

```go
func UpdateMasterMinter(
    ctx context.Context,
    contractAddress,
    newMasterMinter string,
    gasLimit *uint64,
) (*types.Receipt, error)

func UpdateMasterMinterNoWait(
    contractAddress,
    newMasterMinter string,
    gasLimit *uint64,
) (common.Hash, error)
```

```go
func main() {
    ...
    w, _ := wallet.New("<privateKey>")
    w.Connect(alchemy.GetProvider())
    ...
    receipt, err := w.StableCoin().UpdateMasterMinter(ctx, contractAddress, newMasterMinterAddress, nil)
}
```
