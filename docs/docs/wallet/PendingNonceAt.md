---
sidebar_position: 3
---

![](https://img.shields.io/badge/go-geth-lightblue)

PendingNonceAt returns the account nonce of the given account in the pending state.
This is the nonce that should be used for the next transaction.

:::warning
It requires connected wallet.
:::

```go
func PendingNonceAt() (nonce uint64, err error)
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
  nonce, _ := w.PendingNonceAt()
}
```
