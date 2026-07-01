---
sidebar_position: 1
---

The websocket namespace handles characteristic methods such as Subscribe using Geth's internal implementation.

To use this, you need to set UseWebsocket is `true`.

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

or your dev (private) chain.

```go
func main() {
	setting := gas.AlchemySetting{
		PrivateNetworkConfig: gas.PrivateNetworkConfig{
			Url: "ws://127.0.0.1:8545",
		},
	}

	alchemy := gas.NewAlchemy(setting)
}
```
