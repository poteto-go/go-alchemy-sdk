---
sidebar_position: 19
---

![](https://img.shields.io/badge/go-geth-lightblue)

get decimals of erc20 token

:::warning
It requires connected wallet.
:::

```go
func Decimals(contractAddress string) (uint8, error)
```

```go
func main() {
	// ... setup ...
	
	decimals, err := w.ERC20().Decimals("<contractAddress>")
	if err != nil {
		panic(err)
	}
	fmt.Println(decimals)
}
```
