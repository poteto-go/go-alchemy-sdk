package types

type BlockNumberOrHash struct {
	BlockNumber string `json:"blockNumber,omitempty"` // use w/ alchemy method
	BlockHash   string `json:"blockHash,omitempty"`
}

type BlockTagOrHash struct {
	BlockTag  string
	BlockHash string
}

func (btoh *BlockTagOrHash) IsEmpty() bool {
	return btoh.BlockTag == "" && btoh.BlockHash == ""
}
