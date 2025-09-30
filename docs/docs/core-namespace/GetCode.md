![](https://img.shields.io/badge/go-geth-lightblue)

Returns the contract code of the provided address at the block.
If there is no contract deployed, the result is 0x.
BlockTag is latest or BlockNumber.

```go
func GetCode(address string, arg types.BlockTagOrHash) (code string, err error)
```

```go
func main() {
	...
	alchemy := gas.NewAlchemy(setting)
	res, err := alchemy.Core.GetCode(
		"0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
		types.BlockTagOrHash{
			BlockTag: "latest",
		},
	)

	res, err := alchemy.Core.GetCode(
		"0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
		types.BlockTagOrHash{
			BlockHash: "0xbd05b61cc68595a7c30039b2b092ea293c9a2faee20158d578528e399f4d4244",
		},
	)
}
```
