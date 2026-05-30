![](https://img.shields.io/badge/go-geth-lightblue)

WaitMined waits for a transaction with the provided hash and
returns the transaction receipt when it is mined.
It stops waiting when ctx is canceled.

```go
func WaitMined(ctx context.Context, txHash string) (txReceipt *gethTypes.Receipt, err error)
```

```go
func main() {
	...
	alchemy := gas.NewAlchemy(setting)
	txReceipt, err := alchemy.Transact.WaitMined(context.Background(), "<txHash>")
}
```

Pass a cancelable context to stop waiting early:

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
txReceipt, err := alchemy.Transact.WaitMined(ctx, "<txHash>")
```
