ref: [Wallet-StableCoin-Permit](../wallet/StableCoin.md#permit--permitnowait)

![](https://img.shields.io/badge/go-geth-lightblue)

Return the current EIP-2612 permit nonce for the given owner address on a StableCoin contract.

```go
func Nonces(
    contractAddress,
    ownerAddress string,
) (*big.Int, error)
```

```go
func main() {
    ...
    alchemy := gas.NewAlchemy(setting)
    ...
    nonce, err := alchemy.StableCoin.Nonces(contractAddress, ownerAddress)
}
```
