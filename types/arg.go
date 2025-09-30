package types

type BlockNumberOrHash struct {
	BlockNumber string
	BlockHash   string
}

type BlockTagOrHash struct {
	BlockTag  string
	BlockHash string
}

func (btoh *BlockTagOrHash) IsEmpty() bool {
	return btoh.BlockTag == "" && btoh.BlockHash == ""
}
