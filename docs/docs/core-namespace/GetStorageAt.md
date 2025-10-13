![](https://img.shields.io/badge/go-geth-lightblue)

Return the value of the provided position at the provided address, at the provided block in `Bytes32` format.
For inspecting solidity code.

```go
func GetStorageAt(address, position, blockTag string) (value string, err error)
```

```go
func main() {
	...
	alchemy := gas.NewAlchemy(setting)
	res, _ := alchemy.Core.GetStorageAt(
		"0xfe3b557e8fb62b89f4916b721be55ceb828dbd73",
		"0x0",
		"latest",
	)
}
```
