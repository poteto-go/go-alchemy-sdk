![](https://img.shields.io/badge/dev-chain_only-orange)

RevertTo reverts the blockchain state to the provided snapshot id with `evm_revert`.
It returns true if the state was reverted.

A snapshot can only be reverted once. After a successful revert, take a new
snapshot if you want to revert to the same state again.

Only supported on development chains (hardhat, anvil, ganache, ...).

```go
func RevertTo(snapshotId *big.Int) (reverted bool, err error)
```

```go
func main() {
	...
	alchemy := gas.NewAlchemy(setting)

	snapshotId, err := alchemy.Debug.Snapshot()

	// mutate blockchain state (send transactions, ...)

	reverted, err := alchemy.Debug.RevertTo(snapshotId)
}
```
