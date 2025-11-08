---
sidebar_position: 2
---

Connect provider to wallet

```go
func Connect(provider types.IAlchemyProvider)
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
}
```
