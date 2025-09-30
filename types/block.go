package types

import "math/big"

type Withdrawal struct {
	Index     string `json:"index"`
	Validator string `json:"validator"`
	Address   string `json:"address"`
	Amount    string `json:"amount"`
}

type BlockResponse struct {
	Hash                  string       `json:"hash"`
	ParentHash            string       `json:"parentHash"`
	Sha3Uncles            string       `json:"sha3Uncles"`
	Miner                 string       `json:"miner"`
	StateRoot             string       `json:"stateRoot"`
	TransactionsRoot      string       `json:"transactionsRoot"`
	ReceiptsRoot          string       `json:"receiptsRoot"`
	LogsBloom             string       `json:"logsBloom"`
	Number                string       `json:"number"`
	GasLimit              string       `json:"gasLimit"`
	GasUsed               string       `json:"gasUsed"`
	Timestamp             string       `json:"timestamp"`
	ExtraData             string       `json:"extraData"`
	Nonce                 string       `json:"nonce"`
	Size                  string       `json:"size"`
	MixHash               string       `json:"mixHash"`
	Transactions          []string     `json:"transactions"`
	Uncles                []string     `json:"uncles"`
	Difficulty            string       `json:"difficulty,omitempty"`
	BaseFeePerGas         string       `json:"baseFeePerGas,omitempty"`
	WithdrawalsRoot       string       `json:"withdrawalsRoot,omitempty"`
	BlobGasUsed           string       `json:"blobGasUsed,omitempty"`
	ExcessBlobGas         string       `json:"excessBlobGas,omitempty"`
	ParentBeaconBlockRoot string       `json:"parentBeaconBlockRoot,omitempty"`
	Withdrawals           []Withdrawal `json:"withdrawals,omitempty"`
}

type BlockHead struct {
	TotalDifficulty string `json:"totalDifficulty"`
	BlockResponse
}

// refs: https://www.alchemy.com/docs/reference/sdk-getblock
type Block struct {
	Hash       string   `json:"hash"`
	ParentHash string   `json:"parentHash"`
	Number     *big.Int `json:"number"`
	Timestamp  uint64   `json:"timestamp"`
	Nonce      uint64   `json:"nonce"`
	Difficulty *big.Int `json:"difficulty"`
	GasLimit   uint64   `json:"gasLimit"`
	GasUsed    uint64   `json:"gasUsed"`
	Miner      string   `json:"miner"`

	// hash of transactions
	Transactions []string `json:"transactions"`
}
