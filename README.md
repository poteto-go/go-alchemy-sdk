# Go-Alchemy-Sdk

golang sdk for alchemy, inspired by https://github.com/alchemyplatform/alchemy-sdk-js.

This project is aim to be **bridge between alchemy api and geth** objects.

## Documentation

- https://go-alchemy-sdk.poteto-mahiro.com/

## Features

🚀 easily switch private network(`geth`) to Eth-Public Network

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
> This project will breaking-change till `geth` migration is done.

## Major Milestone

- [x] Deploy Contract to Eth: done for `v0.1.0`
- [ ] Private Geth Support `in-progress`
- [ ] Non-Ether Chain Support

## Alchemy Sdk Namespace Support

- [ ] `Core`: in-progress
- [ ] `Nft`
- [ ] `Debug`
- [ ] `Notify`
- [ ] `Portfolio`
- [ ] `Prices`
- [ ] `Transact`
