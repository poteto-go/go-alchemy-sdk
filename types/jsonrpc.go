package types

type AlchemyRequest struct {
	Method  string   `json:"method"`
	Params  []string `json:"params"`
	Id      int      `json:"id"`
	Jsonrpc string   `json:"jsonrpc"`
}

type AlchemyResponse struct {
	Result  string `json:"result"`
	Id      int    `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
}
