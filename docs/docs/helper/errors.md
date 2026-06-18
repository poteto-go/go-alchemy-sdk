---
sidebar_position: 3
---

# Structured Error Types

The SDK provides two typed error structs in the `types` package for programmatic error inspection via `errors.As`.

## RpcError

Wraps a JSON-RPC error with method context.

```go
type RpcError struct {
    Method  string
    Code    int
    Message string
    Err     error
}
```

### Usage

```go
result, err := ether.SomeRpcCall(ctx)
if err != nil {
    var rpcErr *types.RpcError
    if errors.As(err, &rpcErr) {
        fmt.Printf("RPC method %s failed with code %d: %s\n",
            rpcErr.Method, rpcErr.Code, rpcErr.Message)
    }
}
```

The wrapped error is accessible via `errors.Is` and `errors.Unwrap`.

## TxError

Returned when a transaction-related operation fails.

```go
type TxError struct {
    TxHash  common.Hash
    ChainID *big.Int
    Err     error
}
```

### Usage

```go
hash, err := wallet.SendTransaction(ctx, tx)
if err != nil {
    var txErr *types.TxError
    if errors.As(err, &txErr) {
        fmt.Printf("tx %s on chain %s failed: %v\n",
            txErr.TxHash.Hex(), txErr.ChainID, txErr.Err)
    }
}
```

## Sentinel errors

The existing sentinel errors in `constant/errors.go` are preserved for `errors.Is` compatibility. Structured error types wrap them as the inner `Err` field where appropriate.
