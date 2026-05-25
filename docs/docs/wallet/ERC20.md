---
sidebar_position: 12
---

# ERC20

You can execute basic ERC20 methods from the specified wallet.

:::note

Since it performs calls via bytecode, it does not require the contract implementation and can be called as long as the address is available.

:::

:::warning

- It requires connected wallet.
- It does not work on non-Ethereum compatible networks.

:::

```go
func main() {
	alchemy = gas.NewAlchemy(setting)
	w, _ := wallet.New("<privateKey")
	w.Connect(alchemy.GetProvider())

	// call each method
	w.ERC20.MethodXXX()
}
```

## Read Methods

You can fetch token metadata and balance information:

- [TotalSupply](./TotalSupply.md)
- [Allowance](./Allowance.md)
- [Name](./Name.md)
- [Symbol](./Symbol.md)
- [Decimals](./Decimals.md)
- [BalanceOf](./GetBalance.md)

See also: [Namespace Core](../core-namespace/EstimateGas.md) (referencing lower-level implementations).

## Write Methods

### Transfer & TransferNoWait

Transfer ERC20 token from wallet.
Wait for mined or not.

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

	// Execute transaction (wait for mined)
	receipt, err := w.ERC20().Transfer(
		contractAddress,
		"<toAddress>",
		big.NewInt(100),
		nil,
	)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Transaction hash: %s\n", receipt.TxHash)
	fmt.Printf("Status: %d\n", receipt.Status)

	// Execute transaction (no wait)
	txHash, err := w.ERC20().TransferNoWait(
		contractAddress,
		"<toAddress>",
		big.NewInt(100),
		nil,
	)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Transaction hash: %s\n", txHash)
}
```

### Approve & ApproveNoWait

Approve a spender to spend tokens on behalf of the connected wallet.

- [Approve](./Approve.md)

### TransferFrom & TransferFromNoWait

Transfer tokens from another address using a prior allowance.

- [TransferFrom](./TransferFrom.md)
