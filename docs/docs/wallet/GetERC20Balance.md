---
sidebar_position: 13
---

ref: [ERC20-BalanceOf](../erc20-namespace/BalanceOf.md)

![](https://img.shields.io/badge/go-geth-lightblue)

ContractCall wrapper getBalanceOf ERC20 contract.

:::warning

- It requires connected wallet.
- It does not work on non-Ethereum compatible networks.

:::

```go
func main() {
    ...
    alchemy = gas.NewAlchemy(setting)
    w, err := wallet.New(initPrivateKey)
    w.Connect(alchemy.GetProvider())
    contract := artifacts.NewERC20()
    balance, err := w.GetERC20Balance(contract, contractAddress.Hex())
}
```
