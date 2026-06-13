---
sidebar_position: 20
---

# Nft

You can execute basic NFT (ERC721) read methods from the specified wallet.
Each method delegates to the [Nft namespace](../nft-namespace/OwnerOf.md).

:::note

Since it performs calls via bytecode, it does not require the contract implementation and can be called as long as the address is available.

:::

:::warning

- It requires connected wallet.
- It does not work on non-Ethereum compatible networks.

:::

```go
func main() {
	alchemy = gas.NewAlchemy(setting)
	w, _ := wallet.New("<privateKey>")
	w.Connect(alchemy.GetProvider())

	// call each method
	w.Nft().MethodXXX()
}
```

## Read Methods

You can fetch NFT collection metadata and ownership / approval information:

- [OwnerOf](../nft-namespace/OwnerOf.md)
- [TokenURI](../nft-namespace/TokenURI.md)
- [Name](../nft-namespace/Name.md)
- [Symbol](../nft-namespace/Symbol.md)
- [GetApproved](../nft-namespace/GetApproved.md)
- [IsApprovedForAll](../nft-namespace/IsApprovedForAll.md)

```go
owner, err := w.Nft().OwnerOf("<contractAddress>", big.NewInt(1))
uri, err := w.Nft().TokenURI("<contractAddress>", big.NewInt(1))
name, err := w.Nft().Name("<contractAddress>")
symbol, err := w.Nft().Symbol("<contractAddress>")
approved, err := w.Nft().GetApproved("<contractAddress>", big.NewInt(1))
isApproved, err := w.Nft().IsApprovedForAll("<contractAddress>", "<owner>", "<operator>")
```
