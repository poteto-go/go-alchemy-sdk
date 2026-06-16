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

func TestNftContractAddress(t *testing.T) {
	t.Run("returns known addresses using typed symbols", func(t *testing.T) {
		tests := []struct {
			name    string
			network types.Network
			symbol  famous.NftSymbol
			want    common.Address
		}{
			{
				"BAYC on Ethereum mainnet",
				types.EthMainnet, famous.BAYC,
				common.HexToAddress("0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D"),
			},
			{
				"MAYC on Ethereum mainnet",
				types.EthMainnet, famous.MAYC,
				common.HexToAddress("0x60E4d786628Fea6478F785A6d7e704777c86a7c6"),
			},
			{
				"CryptoPunks on Ethereum mainnet",
				types.EthMainnet, famous.CryptoPunks,
				common.HexToAddress("0xb47e3cd837dDF8e4c57F05d70Ab865de6e193BBB"),
			},
			{
				"Azuki on Ethereum mainnet",
				types.EthMainnet, famous.Azuki,
				common.HexToAddress("0xED5AF388653567Af2F388E6224dC7C4b3241C544"),
			},
			{
				"Doodles on Ethereum mainnet",
				types.EthMainnet, famous.Doodles,
				common.HexToAddress("0x8a90CAb2b38dba80c64b7734e58Ee1dB38B8992e"),
			},
			{
				"PudgyPenguins on Ethereum mainnet",
				types.EthMainnet, famous.PudgyPenguins,
				common.HexToAddress("0xBd3531dA5CF5857e7CfAA92426877b022e612cf8"),
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				addr, err := famous.NftContractAddress(tt.network, tt.symbol)

				assert.NoError(t, err)
				assert.Equal(t, tt.want, addr)
			})
		}
	})

	t.Run("returns error for unsupported network", func(t *testing.T) {
		_, err := famous.NftContractAddress(types.SolanaMainnet, famous.BAYC)

		assert.ErrorIs(t, err, famous.ErrNotSupportedNftNetwork)
	})

	t.Run("returns error for unsupported symbol", func(t *testing.T) {
		_, err := famous.NftContractAddress(types.EthMainnet, "UNKNOWN")

		assert.ErrorIs(t, err, famous.ErrNotSupportedNftSymbol)
	})
}

func TestNftSupportedNetworks(t *testing.T) {
	networks := famous.NftSupportedNetworks()

	assert.NotEmpty(t, networks)
	assert.True(t, slices.Contains(networks, types.EthMainnet))
}

func TestNftSupportedSymbols(t *testing.T) {
	t.Run("returns symbols for EthMainnet", func(t *testing.T) {
		symbols := famous.NftSupportedSymbols(types.EthMainnet)

		got := make([]string, len(symbols))
		for i, s := range symbols {
			got[i] = string(s)
		}
		sort.Strings(got)

		assert.Equal(t, []string{"Azuki", "BAYC", "CryptoPunks", "Doodles", "MAYC", "PudgyPenguins"}, got)
	})

	t.Run("returns empty slice for unsupported network", func(t *testing.T) {
		symbols := famous.NftSupportedSymbols(types.SolanaMainnet)

		assert.Empty(t, symbols)
	})
}
