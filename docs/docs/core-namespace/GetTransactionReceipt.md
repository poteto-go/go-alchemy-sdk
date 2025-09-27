![](https://img.shields.io/badge/go-geth-lightblue)

Null if the tx has not been mined.
Returns the transaction receipt for hash.
To stall until the transaction has been mined, consider the waitForTransaction method below.

refs: https://github.com/ethereum/go-ethereum/blob/master/core/types/receipt.go#L53

```go
func GetTransactionReceipt(hash string) (txReceipt *gethTypes.Receipt, err error)
```

```go
func main() {
  ...
  alchemy := gas.NewAlchemy(setting)
  res, _ := alchemy.Core.GetTransactionReceipt(
    "0xc11dacdf03d9fd9297e3a005560e8855608dde8534d9b1053f6608b8541623b8",
  )
}
```
