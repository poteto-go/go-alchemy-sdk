ref: [Wallet-StableCoin-MinterAllowance](../wallet/StableCoin.md#minterallowance)

![](https://img.shields.io/badge/go-geth-lightblue)

Return the remaining mint allowance for a minter on a StableCoin contract (FiatToken/USDC compatibility).

```go
func MinterAllowance(
    contractAddress,
    address string,
) (*big.Int, error)
```

```go
func main() {
    ...
    alchemy := gas.NewAlchemy(setting)
    ...
    allowance, err := alchemy.StableCoin.MinterAllowance(contractAddress, address)
}
```
