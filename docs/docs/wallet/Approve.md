---
sidebar_position: 13
---

![](https://img.shields.io/badge/go-geth-lightblue)

Approve a spender to spend ERC20 tokens on behalf of the connected wallet.

:::warning
It requires connected wallet.
:::

## Approve

Wait for the transaction to be mined and return the receipt.

```go
func Approve(contractAddress, spenderAddress string, amount *big.Int, gasLimit *uint64) (*types.Receipt, error)
```

```go
func main() {
	// ... setup ...

	receipt, err := w.ERC20().Approve(
		"<contractAddress>",
		"<spenderAddress>",
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

## ApproveNoWait

Send the transaction without waiting for it to be mined. Returns the transaction hash immediately.

```go
func ApproveNoWait(contractAddress, spenderAddress string, amount *big.Int, gasLimit *uint64) (common.Hash, error)
```

```go
func main() {
	// ... setup ...

	txHash, err := w.ERC20().ApproveNoWait(
		"<contractAddress>",
		"<spenderAddress>",
		big.NewInt(100),
		nil,
	)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Transaction hash: %s\n", txHash)
}
```
