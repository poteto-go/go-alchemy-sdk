---
sidebar_position: 9
---

![](https://img.shields.io/badge/go-geth-lightblue)

ContractTransact executes a transaction on a deployed contract.
It waits for the transaction to be mined and returns the transaction receipt.

:::warning

- It requires connected wallet.
- It does not work on non-Ethereum compatible networks.

:::

```go
func ContractTransact(
	contract types.ContractInstance,
	contractAddress string,
	data []byte,
) (*gethTypes.Receipt, error)
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

	// Create contract instance
	contract := abi.NewYourContract()
	contractAddress := "0x1234567890123456789012345678901234567890"
	
	// Prepare transaction data (e.g., encoded function call)
	data := abi,PackXXX(<data>)

	// Execute transaction
	receipt, err := w.ContractTransact(contract, contractAddress, data)
	if err != nil {
		panic(err)
	}
	
	fmt.Printf("Transaction hash: %s\n", receipt.TxHash)
	fmt.Printf("Status: %d\n", receipt.Status)
}
```
