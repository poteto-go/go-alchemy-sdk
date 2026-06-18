![](https://img.shields.io/badge/alchemy-only-orange)

Fetches asset transfer history for an address using Alchemy's `alchemy_getAssetTransfers` endpoint.
Supports pagination via `PageKey`.

> **Note:** This is an Alchemy-specific API. It is **not** available on non-Alchemy endpoints such as simulated backends.

- types: `types.AssetTransfersResponse`
- method: `alchemy_getAssetTransfers`
  - refs: https://docs.alchemy.com/reference/alchemy-getassettransfers

```go
func GetAssetTransfers(params types.AssetTransfersParams) (types.AssetTransfersResponse, error)
```

### AssetTransfersParams

| Field               | Type       | Description                                                                 |
| ------------------- | ---------- | --------------------------------------------------------------------------- |
| `FromBlock`         | `string`   | Start block (hex or `"latest"`). Optional.                                  |
| `ToBlock`           | `string`   | End block (hex or `"latest"`). Optional.                                    |
| `FromAddress`       | `string`   | Filter by sender address. Optional.                                         |
| `ToAddress`         | `string`   | Filter by recipient address. Optional.                                      |
| `ContractAddresses` | `[]string` | Filter by contract addresses (ERC-20/721/1155 only). Optional.             |
| `Category`          | `[]string` | **Required.** Transfer categories: `"external"`, `"internal"`, `"erc20"`, `"erc721"`, `"erc1155"`, `"specialnft"`. |
| `WithMetadata`      | `bool`     | Include block timestamp in `Metadata`. Defaults to `false`.                |
| `ExcludeZeroValue`  | `bool`     | Exclude transfers with zero value. Defaults to `false`.                    |
| `MaxCount`          | `int`      | Max results per page (`0` uses the API default of 1000).                   |
| `PageKey`           | `string`   | Pagination cursor from a previous response. Optional.                       |

```go
func main() {
	...
	alchemy := gas.NewAlchemy(setting)

	// Fetch all external and ERC-20 transfers from an address
	res, err := alchemy.Core.GetAssetTransfers(
		types.AssetTransfersParams{
			FromAddress:      "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045",
			Category:         []string{"external", "erc20"},
			ExcludeZeroValue: true,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	for _, transfer := range res.Transfers {
		fmt.Printf("%s: %s -> %s (%s)\n", transfer.BlockNum, transfer.From, transfer.To, transfer.Asset)
	}

	// Paginate using PageKey
	if res.PageKey != "" {
		next, _ := alchemy.Core.GetAssetTransfers(
			types.AssetTransfersParams{
				FromAddress: "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045",
				Category:    []string{"external", "erc20"},
				PageKey:     res.PageKey,
			},
		)
		_ = next
	}
}
```
