---
sidebar_position: 11
---

![](https://img.shields.io/badge/go-geth-lightblue)

ContractTransact executes a transaction on a deployed contract.

You can wait deployment using deployRes.

cf.) [`wallet.ContractTransact`](./ContractTransact.md)

:::warning

- It requires connected wallet.
- It does not work on non-Ethereum compatible networks.

:::

```go
func ContractTransactNoWait(
	contract types.ContractInstance,
	contractAddress string,
	data []byte,
) (tx *gethTypes.Transaction, err error)
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
	data := abi.PackXXX(<data>)

	// Execute transaction
	tx, err := w.ContractTransactNoWait(contract, contractAddress, data)
	if err != nil {
		panic(err)
	}

	// wait to be mined
	receipt, err := alchemy.Transact.WaitMined(tx.Hash().Hex())

	fmt.Printf("Transaction hash: %s\n", receipt.TxHash)
	fmt.Printf("Status: %d\n", receipt.Status)
}
```
