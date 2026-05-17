---
sidebar_position: 12
---

![](https://img.shields.io/badge/go-geth-lightblue)

Transfer ERC20 token from wallet.
Wait for mined.

cf.) [`wallet.ERC20TransferNoWait`](./ERC20TransferNoWait.md)

:::warning

- It requires connected wallet.
- It does not work on non-Ethereum compatible networks.

:::

```go
func ERC20Transfer(
	contractAddress,
	toAddress string,
	amount *big.Int,
	gasLimit *uint64,
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

	// Create contract instance
	contractAddress := "0x1234567890123456789012345678901234567890"

	// Execute transaction
	receipt, err := w.ERC20Transfer(
		contractAddress,
		<toAddress>,
		big.NewInt(100),
		nil,
	)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Transaction hash: %s\n", receipt.TxHash)
	fmt.Printf("Status: %d\n", receipt.Status)
}
```
