---
sidebar_position: 12
---

![](https://img.shields.io/badge/go-geth-lightblue)

ContractCall calls a contract method.
It is used for read-only methods.

### Factory Pattern

:::info

you can use factory pattern for more type safety.

:::

```go
func ContractCall[T any](
	w wallet.Wallet,
	contract types.ContractInstance,
	contractAddress string,
	opts *bind.CallOpts,
	callData []byte,
	unpack func([]byte) (T, error),
) (T, error)
```

```go
func main() {
	setting := gas.AlchemySetting{
		ApiKey:  "api-key",
		Network: types.EthMainnet,
	}
	alchemy := gas.NewAlchemy(setting)

	w, _ := wallet.New("<privateKey>")
	w.Connect(alchemy.GetProvider())

	// Create contract instance
	contract := abi.NewYourContract()
	contractAddress := "0x1234567890123456789012345678901234567890"

	// Prepare call data
	data := abi.PackXXX(<data>)

	// Unpack function
	unpack := func(b []byte) (*big.Int, error) {
		return contract.UnpackXXX(b)
	}

	// Execute call
	result, err := factory.ContractCall(w, contract, contractAddress, nil, data, unpack)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Result: %v\n", result) // 100% *big.Int
}
```

### Call from Wallet

:::warning

- It requires connected wallet.
- It does not work on non-Ethereum compatible networks.

:::

```go
func ContractCall(
	contract types.ContractInstance,
	contractAddress string,
	opts *bind.CallOpts,
	callData []byte,
	unpack func([]byte) (any, error),
) (any, error)
```

```go
func main() {
	setting := gas.AlchemySetting{
		ApiKey:  "api-key",
		Network: types.EthMainnet,
	}
	alchemy := gas.NewAlchemy(setting)

	w, _ := wallet.New("<privateKey>")
	w.Connect(alchemy.GetProvider())

	// Create contract instance
	contract := abi.NewYourContract()
	contractAddress := "0x1234567890123456789012345678901234567890"

	// Prepare call data
	data := abi.PackXXX(<data>)

	// Unpack function
	unpack := func(b []byte) (any, error) {
		return contract.UnpackXXX(b)
	}

	// Execute call
	result, err := w.ContractCall(contract, contractAddress, nil, data, unpack)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Result: %v\n", result)
}
```
