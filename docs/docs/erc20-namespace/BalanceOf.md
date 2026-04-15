ref: [Wallet-GetERC20BalanceOf](../wallet/GetERC20Balance.md)

![](https://img.shields.io/badge/go-geth-lightblue)

Get ERC20 token Balance of provided wallet.

```go
func BalanceOf(
    contract types.ERC20ContractInstance,
    contractAddress,
    walletAddress string,
) (balance *big.Int, err error)
```

```go
func main() {
    ...
    alchemy := gas.NewAlchemy(setting)
    ...
    balance, err := erc20.BalanceOf(artifacts.NewERC20(), contractAddress, walletAddress)
}
```
