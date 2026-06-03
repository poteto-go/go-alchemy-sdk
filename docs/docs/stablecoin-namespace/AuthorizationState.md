![](https://img.shields.io/badge/go-geth-lightblue)

Return whether the EIP-3009 authorization identified by `(authorizer, nonce)` has been used or cancelled.

```go
func AuthorizationState(
    contractAddress string,
    authorizer string,
    nonce [32]byte,
) (bool, error)
```

```go
func main() {
    ...
    alchemy := gas.NewAlchemy(setting)
    ...
    var nonce [32]byte
    // populate nonce with the bytes32 value used when submitting the authorization
    used, err := alchemy.StableCoin.AuthorizationState(contractAddress, authorizerAddress, nonce)
}
```
