---
sidebar_position: 15
---

![](https://img.shields.io/badge/go-geth-lightblue)

get total supply of erc20 token

:::warning
It requires connected wallet.
:::

```go
func TotalSupply(contractAddress string) (*big.Int, error)
```

```go
func main() {
	// ... setup ...
	
	balance, err := w.ERC20().TotalSupply("<contractAddress>")
	if err != nil {
		panic(err)
	}
	fmt.Println(balance)
}
```
