ref: [Wallet-StableCoin-BalanceOf](../wallet/StableCoin.md#balanceof)

![](https://img.shields.io/badge/go-geth-lightblue)

Get StableCoin token balance of the provided walletAddress.

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
    balance, err := alchemy.StableCoin.BalanceOf(contractAddress, walletAddress)
}
```
