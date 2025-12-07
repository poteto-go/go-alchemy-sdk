---
sidebar_position: 10
---

![](https://img.shields.io/badge/go-geth-lightblue)

transact of Deployment Contract to tx pool.

You can wait deployment using deployRes.

cf.) [`wallet.DeployContract`](./DeployContract.md)

:::warning

- It requires connected wallet.
- It does not work on non-Ethernet compatible networks.

:::

```go
func DeployContractNoWait(metaData *bind.MetaData) (deployRes *bind.DeploymentResult, err error)
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

	deployRes, err := w.DeployContractNoWait(&<your-metaData>)
	if err != nil {
		panic(err)
	}

	// wait to deployed
	tx := deployRes.Txs[metaData.ID]
	addr, err := alchemy.Transact.WaitDeployed(tx.Hash().Hex())
	if err != nil {
		panic(err)
	}
	fmt.Println(addr)
}
```
