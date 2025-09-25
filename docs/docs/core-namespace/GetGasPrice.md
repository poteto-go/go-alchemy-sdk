![](https://img.shields.io/badge/go-geth-lightblue)

Returns the best guess of the current gas price to use in a transaction.

```go
func GetGasPrice() (price *big.Int, err error)
```

```go
func main() {
  ...
  alchemy := gas.NewAlchemy(setting)
  res, _ := alchemy.Core.GetGasPrice()
}
```
