package types

import "math/big"

type Signature struct {
	R string   `json:"r"`
	S string   `json:"s"`
	V *big.Int `json:"v"`
}

type TransactionRawResponse struct {
	BlockHash        string `json:"blockHash"`
	BlockNumber      string `json:"blockNumber"`
	From             string `json:"from"`
	Gas              string `json:"gas"`
	GasPrice         string `json:"gasPrice"`
	Hash             string `json:"hash"`
	Input            string `json:"input"`
	Nonce            string `json:"nonce"`
	To               string `json:"to"`
	TransactionIndex string `json:"transactionIndex"`
	Value            string `json:"value"`
	Type             string `json:"type"`
	ChainId          string `json:"chainId"`
	V                string `json:"v"`
	R                string `json:"r"`
	S                string `json:"s"`
}

type TransactionResponse struct {
	BlockNumber          int       `json:"blockNumber"`
	BlockHash            string    `json:"blockHash"`
	Index                int       `json:"index"`
	Hash                 string    `json:"hash"`
	Type                 int       `json:"type"`
	To                   string    `json:"to"`
	From                 string    `json:"from"`
	Nonce                int       `json:"nonce"`
	GasLimit             *big.Int  `json:"gasLimit"`
	GasPrice             *big.Int  `json:"gasPrice"`
	MaxPriorityFeePerGas *big.Int  `json:"maxPriorityFeePerGas"`
	MaxFeePerGas         *big.Int  `json:"maxFeePerGas"`
	Data                 string    `json:"data"`
	Value                *big.Int  `json:"value"`
	ChainID              int       `json:"chainId"`
	Signature            Signature `json:"signature"`
	AccessList           []string  `json:"accessList"`
	BlobVersionedHashes  []string  `json:"blobVersionedHashes"`
	AuthorizationList    []string  `json:"authorizationList"`
}

type TransactionRequest struct {
	Type                 *int      `json:"type"`
	To                   string    `json:"to"`
	From                 string    `json:"from"`
	Nonce                *string   `json:"nonce"`
	GasLimit             *big.Int  `json:"gasLimit"`
	GasPrice             *big.Int  `json:"gasPrice"`
	MaxPriorityFeePerGas *big.Int  `json:"maxPriorityFeePerGas"`
	MaxFeePerGas         *big.Int  `json:"maxFeePerGas"`
	Data                 *string   `json:"data"`
	Value                *big.Int  `json:"value"`
	ChainID              *int      `json:"chainId"`
	AccessList           *[]string `json:"accessList"`
}
