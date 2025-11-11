---
sidebar_position: 1
---

Setting for `go-alchemy-sdk`.

### Public Network

It connect public network w/ Alchemy.
Supported network list is as below.
[link](https://github.com/poteto-go/go-alchemy-sdk/blob/main/types/network.go)

```go
func main() {
	setting := gas.AlchemySetting{
		ApiKey:  "<alchemy-api-key>",
		Network: types.EthSepolia,
	}

	alchemy := gas.NewAlchemy(setting)
}
```

### Private Network

If you want to run w/ private network,
you should define host & port.

```go
func main() {
	setting := gas.AlchemySetting{
		PrivateNetworkConfig: gas.PrivateNetworkConfig{
			Host: "127.0.0.1",
			Port: 32770,
		},
	}

	alchemy := gas.NewAlchemy(setting)
}
```

### More Configuration

```go
func main() {
	setting := gas.AlchemySetting{
		ApiKey:  "<alchemy-api-key>",
		Network: types.EthSepolia,
		BackoffConfig: *types.BackoffConfig{
			Mode:           "exponential",
			MaxRetries:     1,
			InitialDelayMs: 1000,
			MaxDelayMs:     30000,
		},
		MaxRetries: 3,
		RequestTimeout: time.Second * 10,
	}

	alchemy := gas.NewAlchemy(setting)
}
```
