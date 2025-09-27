An enhanced API that gets all transaction receipts for a given block by number or block hash.
Returns geth's Receipt.

- types: `*types.Receipt`

  - refs: https://github.com/ethereum/go-ethereum/blob/master/core/types/receipt.go#L53

- method: `alchemy_getTransactionReceipts`

  - refs: https://www.alchemy.com/docs/data/utility-apis/transactions-receipts-endpoints/alchemy-get-transaction-receipts

```go
func GetTransactionReceipts(arg types.TransactionReceiptsArg) (receipts []*types.Receipt, err error)
```

```go
func main() {
	...
	alchemy := gas.NewAlchemy(setting)
	res, err := alchemy.Core.GetTransactionReceipts(
		types.TransactionReceiptsArg{
			BlockNumber: "0xF1D1C6",
		},
	)
}
```
