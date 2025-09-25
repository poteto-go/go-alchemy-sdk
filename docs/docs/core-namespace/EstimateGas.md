![](https://img.shields.io/badge/go-geth-lightblue)

Returns an estimate of the amount of gas that would be required to submit transaction to the network.
An estimate may not be accurate since there could be another transaction on the network that was not accounted for,
but after being mined affects the relevant state.

```go
func EstimateGas(tx types.TransactionRequest) (price *big.Int, err error)
```

```go
func main() {
  ...
  alchemy := gas.NewAlchemy(setting)
  res, _ := alchemy.Core.EstimateGas(
    types.TransactionRequest{
      From:  "0x44aa93095d6749a706051658b970b941c72c1d53",
      To:    "0xfe3b557e8fb62b89f4916b721be55ceb828dbd73",
      Value: "0x1",
    },
  )
}
```
