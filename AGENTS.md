# Go-Alchemy-Sdk

Golang sdk for alchemy, inspired by https://github.com/alchemyplatform/alchemy-sdk-js.

This project aims to be a **bridge between alchemy api and geth** objects.

It supports not only Alchemy, but also other EVM chains.

## Rule

- Use TDD for development.
- This project internal use go-ethereum client.
  - If you can use go-ethereum, you should use go-ethereum.
- If you create public API, you should update docs(`/docs`) & e2e testing (`/e2e`)

## Testing

### Mocking with alchemymock

This project uses `alchemymock` to mock Alchemy RPC responses (which `ethclient` uses internally) for unit tests.

**Usage:**

1.  **Initialize Mock:**
    Use `newAlchemyMockOnEtherTest(t)` (helper in `ether` package tests) or `alchemymock.NewAlchemyHttpMock`.

    ```go
    alchemyMock := newAlchemyMockOnEtherTest(t)
    defer alchemyMock.DeactivateAndReset()
    ```

2.  **Register Responder:**
    Use `RegisterResponder` to define the JSON-RPC response for a specific method.

    ```go
    // Example: Mocking eth_call
    expectedRes := "0x0000...01" // Hex result
    alchemyMock.RegisterResponder("eth_call", `{"jsonrpc":"2.0","id":1,"result":"`+expectedRes+`"}`)
    ```

    The response string must be a valid JSON-RPC response format.

**Note:** This allows testing functions that interact with `ethclient` without making actual network requests.

### e2e testing

```bash
$ just k-port
$ just k-up # if not detected port
$ just ci-e2e <PORT>
```
