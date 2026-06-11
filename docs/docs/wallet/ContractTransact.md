---
sidebar_position: 10
---

![](https://img.shields.io/badge/go-geth-lightblue)

ContractTransact executes a transaction on a deployed contract.
It waits for the transaction to be mined and returns the transaction receipt.

cf.) [`wallet.ContractTransactNoWait`](./ContractTransactNoWait.md)

:::warning

- It requires connected wallet.
- It does not work on non-Ethereum compatible networks.

:::

```go
func ContractTransact(
	ctx context.Context,
	contractAddress string,
	data []byte,
) (txReceipt *gethTypes.Receipt, err error)
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

	contractAddress := "0x1234567890123456789012345678901234567890"

	// Prepare transaction data (e.g., encoded function call)
	data := abi.PackXXX(<data>)

	// Execute transaction
	receipt, err := w.ContractTransact(context.Background(), contractAddress, data)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Transaction hash: %s\n", receipt.TxHash)
	fmt.Printf("Status: %d\n", receipt.Status)
}
```
