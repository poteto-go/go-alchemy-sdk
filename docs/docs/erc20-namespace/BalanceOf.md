ref: [Wallet-ERC20-BalanceOf](../wallet/ERC20.md#balanceof)

![](https://img.shields.io/badge/go-geth-lightblue)

Get ERC20 token Balance of provided walletAddress.

```go
func BalanceOf(
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
