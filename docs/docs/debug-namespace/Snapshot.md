![](https://img.shields.io/badge/dev-chain_only-orange)

Snapshot takes a snapshot of the current blockchain state with `evm_snapshot`
and returns the snapshot id.

Only supported on development chains (hardhat, anvil, ganache, ...).

```go
func Snapshot() (snapshotId *big.Int, err error)
```

```go
func main() {
	...
	alchemy := gas.NewAlchemy(setting)
	snapshotId, err := alchemy.Debug.Snapshot()
}
```
