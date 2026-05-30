package famous_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/poteto-go/go-alchemy-sdk/famous"
	"github.com/poteto-go/go-alchemy-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestContractAddress(t *testing.T) {
	t.Run("returns known addresses", func(t *testing.T) {
		tests := []struct {
			name    string
			network types.Network
			symbol  string
			want    common.Address
		}{
			{
				"USDC on Ethereum mainnet",
				types.EthMainnet, "USDC",
				common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"),
			},
			{
				"USDT on Ethereum mainnet",
				types.EthMainnet, "USDT",
				common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7"),
			},
			{
				"JPYC on Ethereum mainnet",
				types.EthMainnet, "JPYC",
				common.HexToAddress("0x431D5dfF03120AFA4bDf332c61A6e1766eF37BF9"),
			},
			{
				"USDC on Polygon mainnet",
				types.PolygonMainnet, "USDC",
				common.HexToAddress("0x3c499c542cEF5E3811e1192ce70d8cC03d5c3359"),
			},
			{
				"USDT on Polygon mainnet",
				types.PolygonMainnet, "USDT",
				common.HexToAddress("0xc2132D05D31c914a87C6611C10748AEb04B58e8F"),
			},
			{
				"JPYC on Polygon mainnet",
				types.PolygonMainnet, "JPYC",
				common.HexToAddress("0x6AE7Dfc73E0dDE2aa99ac063DcF7e8A63265108c"),
			},
			{
				"USDC on Polygon Amoy",
				types.PolygonAmoy, "USDC",
				common.HexToAddress("0x41E94Eb019C0762f9Bfcf9Fb1E58725BfB0e7582"),
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				addr, err := famous.ContractAddress(tt.network, tt.symbol)

				assert.NoError(t, err)
				assert.Equal(t, tt.want, addr)
			})
		}
	})

	t.Run("returns error for unsupported network", func(t *testing.T) {
		_, err := famous.ContractAddress(types.SolanaMainnet, "USDC")

		assert.ErrorIs(t, err, famous.ErrNotSupportedNetwork)
	})

	t.Run("returns error for unsupported symbol", func(t *testing.T) {
		_, err := famous.ContractAddress(types.EthMainnet, "UNKNOWN")

		assert.ErrorIs(t, err, famous.ErrNotSupportedSymbol)
	})
}
