---
sidebar_position: 17
---

![](https://img.shields.io/badge/go-geth-lightblue)

get name of erc20 token

:::warning
It requires connected wallet.
:::

```go
func Name(contractAddress string) (string, error)
```

```go
func main() {
	// ... setup ...
	
	name, err := w.ERC20().Name("<contractAddress>")
	if err != nil {
		panic(err)
	}
	fmt.Println(name)
}
```
