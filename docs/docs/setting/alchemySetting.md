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

### Not Alchemy Provider

```go
func main() {
	setting := gas.AlchemySetting{
		PrivateNetworkConfig: gas.PrivateNetworkConfig{
			Url: <providerUrl>
		},
	}

	alchemy := gas.NewAlchemy(setting)
}
```

### Custom Header

```go
func main() {
	setting := gas.AlchemySetting{
		ApiKey:  "<alchemy-api-key>",
		Network: types.EthSepolia,
		CustomHeaders: []http.Header{
			{
				"X-Custom-Header": []string{"custom value"},
			},
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
		BackoffConfig: &types.BackoffConfig{
			Mode:           "exponential",
			MaxRetries:     1,
			InitialDelayMs: 1000,
			MaxDelayMs:     30000,
		},
		MaxRetries:     3,
		RequestTimeout: time.Second * 10,
	}

	alchemy := gas.NewAlchemy(setting)
}
```

### Response Size Limit

`MaxResponseBytes` caps how many bytes are read from an RPC response body. This prevents a malicious or misbehaving endpoint from exhausting process memory. The default is **32 MiB**. Set to `0` to keep the default.

```go
func main() {
	setting := gas.AlchemySetting{
		ApiKey:  "<alchemy-api-key>",
		Network: types.EthSepolia,
		// Override to 8 MiB
		MaxResponseBytes: 8 * 1024 * 1024,
	}

	alchemy := gas.NewAlchemy(setting)
}
```

The limit is enforced on **all** response paths:

- **Alchemy JSON-RPC** (`AlchemyFetch` / `AlchemyBatchFetch`): returns `constant.ErrFailedToReadResponse` when the body exceeds the limit.
- **geth `ethclient` methods** (e.g. `BlockNumber`, `CallContract`): the underlying `http.Client` uses a `limitedTransport` that wraps every response body with `io.LimitReader`, so oversized responses are also truncated here.
