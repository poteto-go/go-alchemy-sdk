---
sidebar_position: 1
---

# Famous

`famous` is a helper package that provides contract addresses for well-known stablecoins across supported networks.

## StableCoinSymbol

Typed constant for well-known stablecoin token symbols. Using these constants instead of raw strings prevents typos at compile time.

```go
type StableCoinSymbol string

const (
    USDC StableCoinSymbol = "USDC"
    USDT StableCoinSymbol = "USDT"
    JPYC StableCoinSymbol = "JPYC"
)
```

## ContractAddress

Returns the contract address of a well-known stablecoin on the given network.

```go
func ContractAddress(network types.Network, symbol StableCoinSymbol) (common.Address, error)
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
    addr, err := famous.ContractAddress(types.PolygonMainnet, famous.JPYC)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(addr.Hex())
}
```

## SupportedNetworks

Returns all networks that have at least one registered stablecoin address.

```go
func SupportedNetworks() []types.Network
```

### Example

```go
networks := famous.SupportedNetworks()
for _, n := range networks {
    fmt.Println(n)
}
```

## SupportedSymbols

Returns all stablecoin symbols registered for the given network. Returns an empty slice if the network is not supported.

```go
func SupportedSymbols(network types.Network) []StableCoinSymbol
```

### Example

```go
symbols := famous.SupportedSymbols(types.EthMainnet)
for _, s := range symbols {
    fmt.Println(s)
}
```
