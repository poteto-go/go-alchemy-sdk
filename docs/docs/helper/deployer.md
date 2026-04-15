---
sidebar_position: 2
---

`deployer` is a helper that simplifies the configuration required for deployment.

## BindDeploymentMetadata

For example, to deploy a contract that has a constructor requiring arguments, you need to store the values in the binary.

It has a helper to do that easily.

```go title="erc20deploy.go"
func main() {
    ...

    erc20Metadata := &artifacts.ERC20MetaData

    // easily bind data
    deployer.BindDeploymentMetadata(erc20Metadata, big.NewInt(1000))
    contractAddress, err := w.DeployContract(erc20Metadata)

    ...
}
```
