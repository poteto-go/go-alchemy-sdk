---
sidebar_position: 20
---

# Nft

You can execute basic NFT (ERC721) read and transfer methods from the specified wallet.
Each read method delegates to the [Nft namespace](../nft-namespace/OwnerOf.md).

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

## Write Methods

You can transfer an NFT held by (or approved to) the connected wallet.
Each method has a `NoWait` variant that returns the transaction hash without
waiting for the transaction to be mined. The trailing `gasLimit *uint64` is
optional — pass `nil` to use the default (`300000`).

### TransferFrom & TransferFromNoWait

Transfer the NFT with the given `tokenId` from one address to another.

```go
receipt, err := w.Nft().TransferFrom(
	context.Background(),
	"<contractAddress>",
	"<fromAddress>",
	"<toAddress>",
	big.NewInt(1),
	nil,
)

// or without waiting for the receipt
txHash, err := w.Nft().TransferFromNoWait("<contractAddress>", "<fromAddress>", "<toAddress>", big.NewInt(1), nil)
```

### SafeTransferFrom & SafeTransferFromNoWait

Same as `TransferFrom`, but the contract verifies that the recipient is able to
receive ERC721 tokens (`onERC721Received`).

```go
receipt, err := w.Nft().SafeTransferFrom(
	context.Background(),
	"<contractAddress>",
	"<fromAddress>",
	"<toAddress>",
	big.NewInt(1),
	nil,
)

txHash, err := w.Nft().SafeTransferFromNoWait("<contractAddress>", "<fromAddress>", "<toAddress>", big.NewInt(1), nil)
```

### SafeTransferFromWithData & SafeTransferFromWithDataNoWait

`safeTransferFrom` overload that forwards additional `data []byte` to the
recipient's `onERC721Received` hook.

```go
data := []byte{0xde, 0xad, 0xbe, 0xef}

receipt, err := w.Nft().SafeTransferFromWithData(
	context.Background(),
	"<contractAddress>",
	"<fromAddress>",
	"<toAddress>",
	big.NewInt(1),
	data,
	nil,
)

txHash, err := w.Nft().SafeTransferFromWithDataNoWait("<contractAddress>", "<fromAddress>", "<toAddress>", big.NewInt(1), data, nil)
```
