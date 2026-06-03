package famous_test

import (
	"slices"
	"sort"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/poteto-go/go-alchemy-sdk/famous"
	"github.com/poteto-go/go-alchemy-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestContractAddress(t *testing.T) {
	t.Run("returns known addresses using typed symbols", func(t *testing.T) {
		tests := []struct {
			name    string
			network types.Network
			symbol  famous.StableCoinSymbol
			want    common.Address
		}{
			{
				"USDC on Ethereum mainnet",
				types.EthMainnet, famous.USDC,
				common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"),
			},
			{
				"USDT on Ethereum mainnet",
				types.EthMainnet, famous.USDT,
				common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7"),
			},
			{
				"JPYC on Ethereum mainnet",
				types.EthMainnet, famous.JPYC,
				common.HexToAddress("0x431D5dfF03120AFA4bDf332c61A6e1766eF37BF9"),
			},
			{
				"USDC on Polygon mainnet",
				types.PolygonMainnet, famous.USDC,
				common.HexToAddress("0x3c499c542cEF5E3811e1192ce70d8cC03d5c3359"),
			},
			{
				"USDT on Polygon mainnet",
				types.PolygonMainnet, famous.USDT,
				common.HexToAddress("0xc2132D05D31c914a87C6611C10748AEb04B58e8F"),
			},
			{
				"JPYC on Polygon mainnet",
				types.PolygonMainnet, famous.JPYC,
				common.HexToAddress("0x6AE7Dfc73E0dDE2aa99ac063DcF7e8A63265108c"),
			},
			{
				"USDC on Polygon Amoy",
				types.PolygonAmoy, famous.USDC,
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
		_, err := famous.ContractAddress(types.SolanaMainnet, famous.USDC)

		assert.ErrorIs(t, err, famous.ErrNotSupportedNetwork)
	})

	t.Run("returns error for unsupported symbol", func(t *testing.T) {
		_, err := famous.ContractAddress(types.EthMainnet, "UNKNOWN")

		assert.ErrorIs(t, err, famous.ErrNotSupportedSymbol)
	})
}

func TestSupportedNetworks(t *testing.T) {
	networks := famous.SupportedNetworks()

	assert.NotEmpty(t, networks)
	assert.True(t, slices.Contains(networks, types.EthMainnet))
	assert.True(t, slices.Contains(networks, types.PolygonMainnet))
	assert.True(t, slices.Contains(networks, types.PolygonAmoy))
}

func TestSupportedSymbols(t *testing.T) {
	t.Run("returns symbols for EthMainnet", func(t *testing.T) {
		symbols := famous.SupportedSymbols(types.EthMainnet)

		got := make([]string, len(symbols))
		for i, s := range symbols {
			got[i] = string(s)
		}
		sort.Strings(got)

		assert.Equal(t, []string{"JPYC", "USDC", "USDT"}, got)
	})

	t.Run("returns symbols for PolygonAmoy", func(t *testing.T) {
		symbols := famous.SupportedSymbols(types.PolygonAmoy)

		assert.Len(t, symbols, 1)
		assert.Equal(t, famous.USDC, symbols[0])
	})

	t.Run("returns empty slice for unsupported network", func(t *testing.T) {
		symbols := famous.SupportedSymbols(types.SolanaMainnet)

		assert.Empty(t, symbols)
	})
}
