# Go-Alchemy-Sdk

golang sdk for alchemy, inspired by https://github.com/alchemyplatform/alchemy-sdk-js

## QuickStart

```go
package main

import (
  "fmt"

  "github.com/poteto-go/go-alchemy-sdk/alchemy"
  "github.com/poteto-go/go-alchemy-sdk/types"
)

func main() {
  setting := alchemy.AlchemySetting{
    ApiKey:  "<api-key>",
    Network: types.EthMainnet,
  }

  alchemy := alchemy.NewAlchemy(setting)
  res, _ := alchemy.Core.GetBlockNumber()
  fmt.Println(res)
}
```
