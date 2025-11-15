---
sidebar_position: 6
---

![](https://img.shields.io/badge/go-geth-lightblue)

sign Transaction by wallet's p8 key
using latest EIP155Signer.

using pending nonce & estimate gas.

:::info
EIP155Signer: sign w/ ChainID to protect replay-attack
:::

:::warning
It requires connected wallet.
:::

```go
func SignTx(txRequest types.TransactionRequest) (signedTx *gethTypes.Transaction, err error)
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

	txRequest := types.TransactionRequest{
		To:       "0x123",
		ChainID:  big.NewInt(1),
		GasLimit: 1000,
		Data:     "0x123",
	}

	signedTx, _ := w.SignTx(txRequest)
}
```
