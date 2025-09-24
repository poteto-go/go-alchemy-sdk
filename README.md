# Go-Alchemy-Sdk

golang sdk for alchemy, inspired by https://github.com/alchemyplatform/alchemy-sdk-js

> [!Important]
> The methods in `geth` are scheduled to undergo major changes with the goal of using `gethclient`.

## Documentation

- https://go-alchemy-sdk.poteto-mahiro.com/

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
