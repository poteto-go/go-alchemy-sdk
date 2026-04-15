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

- [x] Deploy Contract to Eth: done for `v0.1.0`
- [x] Private Geth Support `v0.1.4`
- [ ] geth migration `in-progress`
- [ ] Non-EVM Chain Support

## Alchemy Sdk Namespace Support

- [ ] `Core`: in-progress
- [ ] `ERC20`(not on alchemy-sdk-js): in-progress
- [ ] `Nft`
- [ ] `Debug`
- [ ] `Notify`
- [ ] `Portfolio`
- [ ] `Prices`
- [ ] `Transact`: in-progress
