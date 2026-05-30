---
sidebar_position: 1
---

# ContractAddress

`famous` is a helper package that provides contract addresses for well-known stablecoins across supported networks.

## ContractAddress

Returns the contract address of a well-known stablecoin on the given network.

```go
func ContractAddress(network types.Network, symbol string) (common.Address, error)
```

### Supported Tokens

| Symbol | Networks |
|--------|----------|
| USDC   | EthMainnet, PolygonMainnet, PolygonAmoy |
| USDT   | EthMainnet, PolygonMainnet |
| JPYC   | EthMainnet, PolygonMainnet |

### Errors

- `famous.ErrNotSupportedNetwork` — the network has no registered stablecoin addresses
- `famous.ErrNotSupportedSymbol` — the symbol is not registered for the given network

### Example

```go
import (
    "github.com/poteto-go/go-alchemy-sdk/famous"
    "github.com/poteto-go/go-alchemy-sdk/types"
)

func main() {
    addr, err := famous.ContractAddress(types.PolygonMainnet, "JPYC")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(addr.Hex())
}
```
