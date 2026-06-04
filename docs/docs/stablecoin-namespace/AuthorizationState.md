![](https://img.shields.io/badge/go-geth-lightblue)

Return whether the EIP-3009 authorization identified by `(authorizer, nonce)` has been used or cancelled.

```go
func AuthorizationState(
    contractAddress string,
    authorizer string,
    nonce [32]byte, // the bytes32 nonce that was chosen when the authorization was created
) (bool, error)
```

`nonce` is a required parameter because the EIP-3009 contract stores a `mapping(address => mapping(bytes32 => bool))`. Both the authorizer address and the nonce are needed to identify a specific authorization.

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
