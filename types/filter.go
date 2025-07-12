package types

type Filter struct {
	FromBlock string   `json:"fromBlock,omitempty"`
	ToBlock   string   `json:"toBlock,omitempty"`
	Address   []string `json:"address,omitempty"`
	Topics    []string `json:"topics"`
}

type LogResponse struct {
	LogIndex         string   `json:"logIndex"`
	Topics           []string `json:"topics"`
	Data             string   `json:"data"`
	Address          string   `json:"address"`
	BlockNumber      string   `json:"blockNumber"`
	BlockHash        string   `json:"blockHash"`
	TransactionHash  string   `json:"transactionHash"`
	TransactionIndex string   `json:"transactionIndex"`
	Removed          bool     `json:"removed"`
}
