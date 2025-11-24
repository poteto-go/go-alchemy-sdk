---
sidebar_position: 100
---

ResetPool clears the cached ChainID and TransactOpts.

The wallet caches ChainID and TransactOpts for performance optimization. 
Call this method when you need to refresh these cached values, such as when switching networks or after a long period of inactivity.

:::info

This method is useful when:
- Switching to a different network
- The ChainID has changed
- You want to ensure fresh authentication credentials

:::

```go
func ResetPool()
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

	// Use wallet for transactions
	contract := artifacts.NewYourContract()
	contractAddress := "0x1234567890123456789012345678901234567890"
	data := []byte("encoded transaction data")
	
	receipt, _ := w.ContractTransact(contract, contractAddress, data)
	fmt.Printf("Transaction 1: %s\n", receipt.TxHash)

	// Reset cache if needed (e.g., after network switch)
	w.ResetPool()

	// Next transaction will fetch fresh ChainID and create new auth
	receipt2, _ := w.ContractTransact(contract, contractAddress, data)
	fmt.Printf("Transaction 2: %s\n", receipt2.TxHash)
}
```
