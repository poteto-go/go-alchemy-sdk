---
sidebar_position: 8
---

![](https://img.shields.io/badge/go-geth-lightblue)

DeployContract creates and submits a deployment transaction based on the
deployer bytecode.
It returns the address and creation transaction of the pending contract,
or an error if the creation failed.

cf.) [`wallet.DeployContractNoWait`](./DeployContractNoWait.md)

:::warning

- It requires connected wallet.
- It does not work on non-Ethernet compatible networks.

:::

```go
func DeployContract(metaData *bind.MetaData) (deployedAddr common.Address, err error)
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

	address, err := w.DeployContract(&<your-metaData>)
	if err != nil {
		panic(err)
	}
	fmt.Println(address)
}
```
