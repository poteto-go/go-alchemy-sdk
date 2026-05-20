---
sidebar_position: 18
---

![](https://img.shields.io/badge/go-geth-lightblue)

get symbol of erc20 token

:::warning
It requires connected wallet.
:::

```go
func Symbol(contractAddress string) (string, error)
```

```go
func main() {
	// ... setup ...
	
	symbol, err := w.ERC20().Symbol("<contractAddress>")
	if err != nil {
		panic(err)
	}
	fmt.Println(symbol)
}
```
