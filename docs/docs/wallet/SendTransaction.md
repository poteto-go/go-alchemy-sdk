---
sidebar_position: 5
---

![](https://img.shields.io/badge/go-geth-lightblue)

Signs tx and sends it to the pending pool for execution.

:::warning
It requires connected wallet.
:::

```go
func SendTransaction(txRequest types.TransactionRequest) (err error)
```

```go
func main() {
	setting := gas.AlchemySetting{
		ApiKey:  "api-key",
		Network: types.EthMainnet,
	}
	alchemy := gas.NewAlchemy(setting)

	w, _ := wallet.New("privateKey")
	w.Connect(alchemy.GetProvider())

	txRequest := types.TransactionRequest{
		To:       "0x123",
		ChainID:  big.NewInt(1),
		GasLimit: 1000,
		Data:     "0x123",
	}

	w.SendTransaction(txRequest)
}
```
