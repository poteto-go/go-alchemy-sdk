![](https://img.shields.io/badge/go-geth-lightblue)

Returns `(maxPriorityFeePerGas, maxFeePerGas)` ready to use in a `TransactionRequest`.

`maxFeePerGas` is derived as `baseFeePerGas * 2 + maxPriorityFeePerGas` (standard EIP-1559 formula).
Returns an error on chains that do not support EIP-1559.

```go
func SuggestEIP1559Fees() (maxPriorityFeePerGas *big.Int, maxFeePerGas *big.Int, err error)
```

```go
func main() {
	...
	alchemy := gas.NewAlchemy(setting)
	tip, maxFee, _ := alchemy.Core.SuggestEIP1559Fees()

	tx := types.TransactionRequest{
		MaxPriorityFeePerGas: tip.String(),
		MaxFeePerGas:         maxFee.String(),
		// ...
	}
}
```
