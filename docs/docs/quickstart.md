---
sidebar_position: 1
slug: /
---

# Quick Start

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
