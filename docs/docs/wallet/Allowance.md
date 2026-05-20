---
sidebar_position: 16
---

![](https://img.shields.io/badge/go-geth-lightblue)

get allowance of erc20 token

:::warning
It requires connected wallet.
:::

```go
func Allowance(contractAddress, owner, spender string) (*big.Int, error)
```

```go
func main() {
	// ... setup ...
	
	allowance, err := w.ERC20().Allowance("<contractAddress>", "<ownerAddress>", "<spenderAddress>")
	if err != nil {
		panic(err)
	}
	fmt.Println(allowance)
}
```
