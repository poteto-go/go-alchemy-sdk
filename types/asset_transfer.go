package types

// AssetTransfersParams is the input for GetAssetTransfers.
// Category is required; all other fields are optional.
// MaxCount of 0 uses the Alchemy API default (1000).
//
// NOTE: This method is Alchemy-specific and is not available on
// non-Alchemy endpoints such as simulated backends.
type AssetTransfersParams struct {
	FromBlock         string
	ToBlock           string
	FromAddress       string
	ToAddress         string
	ContractAddresses []string
	Category          []string // "external","internal","erc20","erc721","erc1155","specialnft"
	WithMetadata      bool
	ExcludeZeroValue  bool
	MaxCount          int // 0 → API default (1000)
	PageKey           string
}

type Erc1155Metadata struct {
	TokenId string
	Value   string
}

type RawContract struct {
	Value   string
	Address string
	Decimal string
}

type AssetTransferMetadata struct {
	BlockTimestamp string
}

type AssetTransfer struct {
	BlockNum        string
	UniqueId        string
	Hash            string
	From            string
	To              string
	Value           *float64
	Erc721TokenId   string
	Erc1155Metadata []Erc1155Metadata
	TokenId         string
	Asset           string
	Category        string
	RawContract     RawContract
	Metadata        *AssetTransferMetadata
}

type AssetTransfersResponse struct {
	Transfers []AssetTransfer
	PageKey   string
}
