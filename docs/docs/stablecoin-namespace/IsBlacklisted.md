ref: [Wallet-StableCoin-IsBlacklisted](../wallet/StableCoin.md#isblacklisted)

![](https://img.shields.io/badge/go-geth-lightblue)

Check whether an address is blacklisted on the StableCoin contract (FiatToken/USDC compatibility).

```go
func IsBlacklisted(
    contractAddress,
    address string,
) (bool, error)
```

```go
func main() {
    ...
    alchemy := gas.NewAlchemy(setting)
    ...
    isBlacklisted, err := alchemy.StableCoin.IsBlacklisted(contractAddress, address)
}
```
