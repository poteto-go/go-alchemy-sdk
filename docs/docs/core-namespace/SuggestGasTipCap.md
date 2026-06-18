![](https://img.shields.io/badge/go-geth-lightblue)

Returns the suggested `maxPriorityFeePerGas` (EIP-1559 tip) via `eth_maxPriorityFeePerGas`.

```go
func SuggestGasTipCap() (tip *big.Int, err error)
```

```go
func main() {
	...
	alchemy := gas.NewAlchemy(setting)
	tip, _ := alchemy.Core.SuggestGasTipCap()
}
```
