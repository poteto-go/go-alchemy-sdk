![](https://img.shields.io/badge/go-geth-lightblue)

Returns the transaction with hash or null if the transaction is unknown.

If a transaction has not been mined, this method will search the transaction pool.
Various backends may have more restrictive transaction
pool access (e.g. if the gas price is too low or the transaction was only recently sent and not yet indexed) in which case this method may also return null.

- types: `*types.Transaction`

  - refs: https://github.com/ethereum/go-ethereum/blob/master/core/types/transaction.go#L57

```go
func GetTransaction(hash string) (tx *types.Transaction, isPending bool, err error)
```

```go
func main() {
	...
	alchemy := gas.NewAlchemy(setting)
	res, isPending, _ := alchemy.Core.GetTransaction(
		"0x9b300f515857b60d52cd23fb75b56aeae6eb96aa60778ce3758bc8f68db061e3",
	)
}
```
