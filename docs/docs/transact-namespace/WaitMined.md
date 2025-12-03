![](https://img.shields.io/badge/go-geth-lightblue)

WaitMined waits for a transaction with the provided hash and
returns the transaction receipt when it is mined.
It stops waiting when ctx is canceled.

```go
func WaitMined(txHash string) (txReceipt *gethTypes.Receipt, err error)
```

```go
func main() {
	...
	alchemy := gas.NewAlchemy(setting)
	txReceipt, err := alchemy.Transact.WaitMined("<txHash>")
}
```
