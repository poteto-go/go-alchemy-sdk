package types

type BlockNumberOrHash struct {
	BlockNumber string `json:"blockNumber,omitempty"`
	BlockHash   string `json:"blockHash,omitempty"`
}

type BlockTagOrHash struct {
	BlockTag  string
	BlockHash string
}
