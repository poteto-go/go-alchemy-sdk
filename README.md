![coverage](_fixture/coverage.svg)

# Go-Alchemy-Sdk

golang sdk for alchemy, inspired by https://github.com/alchemyplatform/alchemy-sdk-js.

This project aims to be a **bridge between alchemy api and geth** objects.

It supports not only Alchemy, but also other EVM chains.

## Documentation

- https://go-alchemy-sdk.poteto-mahiro.com/

## Features

🖧 easily switch EVM network | private network to Eth-Public Network.

⚙️ Designed to meet the needs of custom private chains, such as those requiring custom headers.

🧪 easily mock rpc responses w/ `alchemymock`.

🪙 StableCoin native support

## QuickStart

```go
package main

import (
  "fmt"

  "github.com/poteto-go/go-alchemy-sdk/gas"
  "github.com/poteto-go/go-alchemy-sdk/types"
)

func main() {
  setting := gas.AlchemySetting{
    ApiKey:  "<api-key>",
    Network: types.EthMainnet,
  }

  alchemy := gas.NewAlchemy(setting)
  res, _ := alchemy.Core.GetBlockNumber()
  fmt.Println(res)
}
```

## Caution

> [!Caution]
> This project will have breaking changes until `geth` migration is done.

## Major Milestone

- [x] Deploy Contract to Eth: `v0.1.0`
- [x] Smart-Contract Tx & call support: `v0.2.0`
- [x] Private Geth Support & `ERC20.Transfer`: `v0.3.0`
- [x] `ERC20` method fully support: `v0.4.0`
- [x] `SC(Stable Coin)` method fully support: `v0.6.0`
- [ ] `NFT(ERC721)` method fully support: `rc0.8.0`
- [ ] stable release: `rc1.0.0`
