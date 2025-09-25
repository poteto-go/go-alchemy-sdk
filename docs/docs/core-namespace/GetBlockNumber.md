![](https://img.shields.io/badge/go-geth-lightblue)

get the number of the most recent block.

```go
func GetBlockNumber() (blockNumber uint64, err error)
```

```go
func main() {
  ...
  alchemy := gas.NewAlchemy(setting)
  res, _ := alchemy.Core.GetBlockNumber()
}
```
