---
sidebar_position: 4
---

![](https://img.shields.io/badge/go-geth-lightblue)

get balance of native token

:::warning
It requires connected wallet.
:::

```go
func GetBalance() (balance *big.Int, err error)
```

```go
func main() {
	setting := gas.AlchemySetting{
		ApiKey:  "api-key",
		Network: types.EthMainnet,
	}
	alchemy := gas.NewAlchemy(setting)

	w, _ := wallet.New("<privateKey>")
	w.Connect(alchemy.GetProvider())

	balance, err := w.GetBalance()
	if err != nil {
		panic(err)
	}
	fmt.Println(balance)
}
```
