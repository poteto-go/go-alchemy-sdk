package types

type TokenBalance struct {
	ContractAddress string `json:"contractAddress"`
	TokenBalance    string `json:"tokenBalance"`
	Error           error  `json:"-"`
}
type TokenBalanceResponse struct {
	Address       string         `json:"address"`
	PageKey       string         `json:"pageKey"`
	TokenBalances []TokenBalance `json:"tokenBalances"`
}

type TokenBalanceOption struct {
	ContractAddresses []string
}
