![](https://img.shields.io/badge/go-geth-lightblue)

WaitDeployed waits for a contract deployment transaction with the provided hash and returns the contract address.
It stops waiting when ctx is canceled.

```go
func WaitDeployed(ctx context.Context, txHash string) (contractAddress common.Address, err error)
```

```go
func main() {
	...
	alchemy := gas.NewAlchemy(setting)
	...
	contractAddress, err := alchemy.Transact.WaitDeployed(context.Background(), "<deployedHash>")
}
```

Pass a cancelable context to stop waiting early:

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
contractAddress, err := alchemy.Transact.WaitDeployed(ctx, "<deployedHash>")
```
