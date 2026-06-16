---
sidebar_position: 2
---

# Famous NFT

`famous` also provides contract addresses for well-known NFT collections across supported networks.

## NftSymbol

Typed constant for well-known NFT collection symbols. Using these constants instead of raw strings prevents typos at compile time.

```go
type NftSymbol string

const (
    BAYC          NftSymbol = "BAYC"          // Bored Ape Yacht Club
    MAYC          NftSymbol = "MAYC"          // Mutant Ape Yacht Club
    CryptoPunks   NftSymbol = "CryptoPunks"   // CryptoPunks
    Azuki         NftSymbol = "Azuki"         // Azuki
    Doodles       NftSymbol = "Doodles"       // Doodles
    PudgyPenguins NftSymbol = "PudgyPenguins" // Pudgy Penguins
)
```

## NftContractAddress

Returns the contract address of a well-known NFT collection on the given network.

```go
func NftContractAddress(network types.Network, symbol NftSymbol) (common.Address, error)
```

### Supported Collections

| Symbol        | Networks   |
|---------------|------------|
| BAYC          | EthMainnet |
| MAYC          | EthMainnet |
| CryptoPunks   | EthMainnet |
| Azuki         | EthMainnet |
| Doodles       | EthMainnet |
| PudgyPenguins | EthMainnet |

### Errors

- `famous.ErrNotSupportedNftNetwork` — the network has no registered NFT collection addresses
- `famous.ErrNotSupportedNftSymbol` — the symbol is not registered for the given network

### Example

```go
import (
    "github.com/poteto-go/go-alchemy-sdk/famous"
    "github.com/poteto-go/go-alchemy-sdk/types"
)

func main() {
    addr, err := famous.NftContractAddress(types.EthMainnet, famous.BAYC)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(addr.Hex())
}
```

## NftSupportedNetworks

Returns all networks that have at least one registered NFT collection address.

```go
func NftSupportedNetworks() []types.Network
```

### Example

```go
networks := famous.NftSupportedNetworks()
for _, n := range networks {
    fmt.Println(n)
}
```

## NftSupportedSymbols

Returns all NFT collection symbols registered for the given network. Returns an empty slice if the network is not supported.

```go
func NftSupportedSymbols(network types.Network) []NftSymbol
```

### Example

```go
symbols := famous.NftSupportedSymbols(types.EthMainnet)
for _, s := range symbols {
    fmt.Println(s)
}
```
