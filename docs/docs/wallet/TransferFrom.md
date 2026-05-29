---
sidebar_position: 14
---

![](https://img.shields.io/badge/go-geth-lightblue)

Transfer ERC20 tokens from another address to a recipient, using a prior allowance granted via [Approve](./Approve.md).

:::warning
It requires connected wallet.
:::

:::note
The connected wallet must have been approved as a spender by the `fromAddress` before calling this method.
:::

## TransferFrom

Wait for the transaction to be mined and return the receipt.

```go
func TransferFrom(ctx context.Context, contractAddress, fromAddress, toAddress string, amount *big.Int, gasLimit *uint64) (*types.Receipt, error)
```

```go
func main() {
	// ... setup ...

	// Assumes fromAddress has already approved the connected wallet as a spender
	receipt, err := w.ERC20().TransferFrom(
		context.Background(),
		"<contractAddress>",
		"<fromAddress>",
		"<toAddress>",
		big.NewInt(100),
		nil, // use default gas limit (300000)
	)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Transaction hash: %s\n", receipt.TxHash)
	fmt.Printf("Status: %d\n", receipt.Status)
}
```

## TransferFromNoWait

Send the transaction without waiting for it to be mined. Returns the transaction hash immediately.

```go
func TransferFromNoWait(contractAddress, fromAddress, toAddress string, amount *big.Int, gasLimit *uint64) (common.Hash, error)
```

```go
func main() {
	// ... setup ...

	txHash, err := w.ERC20().TransferFromNoWait(
		"<contractAddress>",
		"<fromAddress>",
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
