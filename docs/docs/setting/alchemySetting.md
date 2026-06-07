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

### Custom Transport

`Transport` lets you plug in your own `http.RoundTripper` for the actual HTTP communication of every RPC call. Because a single transport is shared and reused across calls, this is also where you control the connection pool, so you stop paying the TCP/TLS handshake cost on repeated calls to the same host.

Use it for:

- connection pooling settings
- retry / backoff
- request tracing
- latency / error-rate metrics
- provider-level benchmarking

If `Transport` is `nil`, `http.DefaultTransport` is used. The SDK always applies its [response-size cap](#response-size-limit) on top of whatever transport you provide, so you never lose that protection.

```go
func main() {
	setting := gas.AlchemySetting{
		ApiKey:  "<alchemy-api-key>",
		Network: types.EthSepolia,
		// Reuse connections aggressively for back-to-back RPC calls.
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
			IdleConnTimeout:     90 * time.Second,
		},
	}

	alchemy := gas.NewAlchemy(setting)
}
```

#### Use case: benchmark your blockchain node

Wrap a transport to measure per-call latency, then point the SDK at different (private) RPC providers to compare them under identical SDK behavior:

```go
// latencyTransport records how long each RPC round-trip takes.
type latencyTransport struct {
	base http.RoundTripper
}

func (t *latencyTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()
	resp, err := t.base.RoundTrip(req)
	log.Printf("rpc %s took %s", req.URL.Host, time.Since(start))
	return resp, err
}

func main() {
	setting := gas.AlchemySetting{
		PrivateNetworkConfig: gas.PrivateNetworkConfig{
			Url: "<your-node-or-provider-url>",
		},
		Transport: &latencyTransport{base: http.DefaultTransport},
	}

	alchemy := gas.NewAlchemy(setting)

	// Every RPC call now logs its latency, so you can benchmark and
	// compare providers (e.g. Alchemy vs. a private endpoint).
	alchemy.Core.GetBlockNumber()
}
```

### JWT Secret (Engine API Authentication)

geth's [Engine API](https://github.com/ethereum/execution-apis/blob/main/src/engine/authentication.md) requires JWT authentication. Set `JwtSecret` in `PrivateNetworkConfig` to enable it.

`JwtSecret` must be a **hex-encoded 64-character string** (32 raw bytes), matching the `--authrpc.jwtsecret` file used by your geth node.

```go
func main() {
	setting := gas.AlchemySetting{
		PrivateNetworkConfig: gas.PrivateNetworkConfig{
			Host:      "127.0.0.1",
			Port:      8551, // geth default Engine API port
			JwtSecret: "bcdf20249abf0ed6d944c0288fad489e33f66b3960d9e6229c1cd214ed3bbe31",
		},
	}

	alchemy := gas.NewAlchemy(setting)
}
```

When `JwtSecret` is set the SDK:

1. Signs a JWT (`HS256`, `iat` + `exp` claims) and sends it as `Authorization: Bearer <token>` on every request.
2. Automatically recreates the underlying `ethclient` before the token expires — geth accepts tokens within a **60-second** `iat` window, and the SDK refreshes ~5 seconds early to absorb clock skew.

If `JwtSecret` is omitted the `Authorization` header is not sent (suitable for unauthenticated HTTP endpoints).
