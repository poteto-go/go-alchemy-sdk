![](https://img.shields.io/badge/go-geth-lightblue)

Returns the block from the network based on the provided block number or hash.
Transactions on the block are represented as an array of transaction hashes.

- types: `*types.Block`

  - refs: https://github.com/ethereum/go-ethereum/blob/master/core/types/block.go#L206

```go
func GetBlock(blockHashOrBlockTag types.BlockTagOrHash) (block *types.Block, err error)
```

```go
func main() {
	...
	alchemy := gas.NewAlchemy(setting)
	res, _ := alchemy.Core.GetBlock(types.BlockTagOrHash{
		BlockHash: "0xf7756d836b6716aaeffc2139c032752ba5acf02fe94acb65743f0d177554b2e2",
	})

	res, _ = alchemy.Core.GetBlock(types.BlockTagOrHash{
		BlockTag: "0x123",
	})
}
```
