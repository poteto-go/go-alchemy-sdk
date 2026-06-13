---
sidebar_position: 2
---

# Simulated Backend

## Overview

The `simulatedBackend` allows you to run high-speed tests without launching a real blockchain network. It is fully integrated with the `go-alchemy-sdk`, enabling you to use the SDK's namespaces (Core, Transact, Nft, etc.) in an in-process, deterministic environment.

## Capabilities

Using `gas.NewSimulatedAlchemy(backend)`, you can interact with a simulated environment as if it were a real network:
- **Contract Deployment:** Deploy and interact with contracts.
- **Transact:** Send transactions and mine them on-demand.
- **ERC20/ERC721/StableCoin:** Full support for standard contract operations.

## Limitations

Because the simulated backend runs in-process and does not have an HTTP JSON-RPC layer:
- **RPC Methods:** Any method that relies on `provider.Send` or direct `rpc.Client` calls (e.g., `Core.GetBalance`, `Debug.Snapshot`, `batch.Batcher`) is **unsupported**.
- **HTTP/Transport:** Custom `http.RoundTripper` configurations are not applicable.

## Tutorial: Integration in Unit Tests

The following example demonstrates how to set up an in-process `simulated.Backend` to test your smart contract interactions.

```go
func TestMyContractInteraction(t *testing.T) {
	// 1. Setup funded simulated backend
	balance := new(big.Int).Mul(big.NewInt(1_000), big.NewInt(1_000_000_000_000_000_000)) // 1000 ETH
	initAddress := common.HexToAddress("0x...") // Your funded test address
	
	backend := simulated.NewBackend(gethTypes.GenesisAlloc{
		initAddress: {Balance: balance},
	})
	defer backend.Close() // Essential cleanup

	// 2. Initialize SimulatedAlchemy
	alchemy, err := gas.NewSimulatedAlchemy(backend)
	assert.NoError(t, err)

	// 3. Now use the 'alchemy' object for operations
	// e.g., deploy, transact, or call contract methods
}
```

For more comprehensive usage, refer to the examples in `e2e/simulated_test.go`.
