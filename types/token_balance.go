package types

import (
	"encoding/json"
	"errors"
)

type TokenBalance struct {
	ContractAddress string `json:"contractAddress"`
	TokenBalance    string `json:"tokenBalance"`
	Error           error  `json:"-"`
}

func (t *TokenBalance) UnmarshalJSON(data []byte) error {
	type Alias TokenBalance
	alias := &struct {
		Error string `json:"error,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(t),
	}
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}
	if alias.Error != "" {
		t.Error = errors.New(alias.Error)
	}
	return nil
}

type TokenBalanceResponse struct {
	Address       string         `json:"address"`
	TokenBalances []TokenBalance `json:"tokenBalances"`
}
