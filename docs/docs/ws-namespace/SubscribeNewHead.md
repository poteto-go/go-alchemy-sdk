---
sidebar_position: 2
---

**Subscribing to New Blocks:**

set provider

```go
func main() {
	setting := gas.AlchemySetting{
		ApiKey:       "<alchemy-api-key>",
		Network:      types.EthSepolia,
		UseWebsocket: true,
	}

	alchemy := gas.NewAlchemy(setting)
	// The WebSocket socket is persistent; close it when you are done.
	defer alchemy.GetProvider().Eth().Shutdown()
}
```

1. Create a new channel that will be receiving the latest block headers.

```go
headers := make(chan *types.Header)
```

2. Call SubscribeNewHead

```go
sub, err := gas.WS.SubscribeNewHead(
  context.Background(),
  headers,
)
if err != nil {
  log.Fatal(err)
}
```

3. Subscribe

```go
for {
  select {
  case err := <-sub.Err():
    log.Fatal(err)
  case header := <-headers:
    fmt.Println(header.Hash().Hex())
  }
}
```
