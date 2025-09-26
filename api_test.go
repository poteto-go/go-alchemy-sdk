package test

import (
	"fmt"
	"math/big"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/poteto-go/go-alchemy-sdk/gas"
	"github.com/poteto-go/go-alchemy-sdk/types"
	"github.com/stretchr/testify/assert"
)

/*
Check API Request
Core Namespace
*/
var setting gas.AlchemySetting
var address string

func TestMain(m *testing.M) {
	setup()

	m.Run()
}

func setup() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("this doesn't run on local")
	}

	setting = gas.AlchemySetting{
		ApiKey:  os.Getenv("API_KEY"),
		Network: types.EthMainnet,
	}

	address = os.Getenv("ADDRESS")
}

func TestAPI_Core_GetTokenMetadata(t *testing.T) {
	alchemy := gas.NewAlchemy(setting)

	t.Run("usdc response", func(t *testing.T) {
		res, err := alchemy.Core.GetTokenMetadata(
			"0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
		)

		assert.Nil(t, err)
		assert.Equal(
			t,
			types.TokenMetadataResponse{
				Name:     "USDC",
				Symbol:   "USDC",
				Decimals: 6,
				Logo:     "https://static.alchemyapi.io/images/assets/3408.png",
			},
			res,
		)
	})
}

func TestAPI_Core_EstimateGas(t *testing.T) {
	alchemy := gas.NewAlchemy(setting)

	t.Run("over 0 response", func(t *testing.T) {
		res, err := alchemy.Core.EstimateGas(
			types.TransactionRequest{
				From:  "0x44aa93095d6749a706051658b970b941c72c1d53",
				To:    "0xfe3b557e8fb62b89f4916b721be55ceb828dbd73",
				Value: "0x1",
			},
		)

		assert.Nil(t, err)
		assert.Equal(
			t,
			res.Cmp(big.NewInt(0)),
			1,
		)
	})
}

func TestAPI_Core_GetLogs(t *testing.T) {
	alchemy := gas.NewAlchemy(setting)

	t.Run("get logs", func(t *testing.T) {
		res, err := alchemy.Core.GetLogs(
			types.Filter{
				FromBlock: "0x137d3c2",
				ToBlock:   "0x137d3c3",
				Address:   "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
				Topics:    []string{},
			},
		)

		assert.Nil(t, err)
		assert.Greater(t, len(res), 0)
	})
}

func TestAPI_Core_GetBlockNumber(t *testing.T) {
	alchemy := gas.NewAlchemy(setting)

	t.Run("get block number", func(t *testing.T) {
		res, err := alchemy.Core.GetBlockNumber()

		assert.Nil(t, err)
		assert.Greater(t, res, uint64(0))
	})
}

func TestAPI_Core_GetGasPrice(t *testing.T) {
	alchemy := gas.NewAlchemy(setting)

	t.Run("get gas price", func(t *testing.T) {
		res, err := alchemy.Core.GetGasPrice()

		assert.Nil(t, err)
		assert.Equal(t, res.Cmp(big.NewInt(0)), 1)
	})
}

func TestAPI_Core_GetBalance(t *testing.T) {
	setting.Network = types.PolygonAmoy // I don't have on eth
	alchemy := gas.NewAlchemy(setting)

	t.Run("get balance", func(t *testing.T) {
		res, err := alchemy.Core.GetBalance(
			address, "latest",
		)

		assert.Nil(t, err)
		assert.Equal(
			t,
			res.Cmp(big.NewInt(0)),
			1,
		)
	})
}

func TestAPI_Core_GetTransaction(t *testing.T) {
	alchemy := gas.NewAlchemy(setting)

	t.Run("get transaction", func(t *testing.T) {
		res, err := alchemy.Core.GetTransaction(
			"0x591b59017dc8b5b154dbca6b27811206e2794f636c7a9cb26a6b26afe0526eb1",
		)

		assert.Nil(t, err)
		assert.Equal(
			t,
			"0x591b59017dc8b5b154dbca6b27811206e2794f636c7a9cb26a6b26afe0526eb1",
			res.Hash,
		)
	})
}

func TestAPI_Core_GetStorageAt(t *testing.T) {
	alchemy := gas.NewAlchemy(setting)

	t.Run("get storage at", func(t *testing.T) {
		res, err := alchemy.Core.GetStorageAt(
			"0xfe3b557e8fb62b89f4916b721be55ceb828dbd73",
			"0x0",
			"latest",
		)

		assert.Nil(t, err)
		assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000000", res)
	})
}

func TestAPI_Core_Call(t *testing.T) {
	alchemy := gas.NewAlchemy(setting)

	t.Run("call", func(t *testing.T) {
		res, err := alchemy.Core.Call(
			types.TransactionRequest{
				To:    address,
				Value: "0x1",
			},
			"latest",
		)

		assert.Nil(t, err)
		assert.Equal(t, "0x", res)
	})
}

func TestAPI_Core_GetTransactionReceipt(t *testing.T) {
	setting.Network = types.EthMainnet
	alchemy := gas.NewAlchemy(setting)

	t.Run("get transaction receipt", func(t *testing.T) {
		txHash := "0xc11dacdf03d9fd9297e3a005560e8855608dde8534d9b1053f6608b8541623b8"
		res, err := alchemy.Core.GetTransactionReceipt(
			txHash,
		)

		assert.Nil(t, err)
		assert.Equal(t, res.TransactionHash, txHash)
	})
}

func TestAPI_Core_GetTransactionReceipts(t *testing.T) {
	setting.Network = types.EthMainnet
	alchemy := gas.NewAlchemy(setting)

	t.Run("get transaction receipts", func(t *testing.T) {
		blockNumber := "0xF1D1C6"
		res, err := alchemy.Core.GetTransactionReceipts(
			types.TransactionReceiptsArg{
				BlockNumber: blockNumber,
			},
		)

		assert.Nil(t, err)
		assert.Greater(t, len(res), 0)
	})
}

func TestAPI_Core_GetBlock(t *testing.T) {
	setting.Network = types.EthMainnet
	alchemy := gas.NewAlchemy(setting)

	t.Run("get block by block hash", func(t *testing.T) {
		blockHash := "0xf7756d836b6716aaeffc2139c032752ba5acf02fe94acb65743f0d177554b2e2"
		res, err := alchemy.Core.GetBlock(
			types.BlockHashOrBlockTag{
				BlockHash: blockHash,
			},
		)

		assert.Nil(t, err)
		assert.Equal(t, res.Hash, blockHash)
	})

	t.Run("get block by block hash, but result is nil error", func(t *testing.T) {
		blockHash := "0x123"
		_, err := alchemy.Core.GetBlock(
			types.BlockHashOrBlockTag{
				BlockHash: blockHash,
			},
		)

		assert.Error(t, err)
	})

	t.Run("get block by block number", func(t *testing.T) {
		blockNumber := "0x68b3"
		res, err := alchemy.Core.GetBlock(
			types.BlockHashOrBlockTag{
				BlockTag: blockNumber,
			},
		)

		assert.Nil(t, err)
		assert.Equal(t, res.Number.Cmp(big.NewInt(26803)), 0)
	})

	t.Run("get block by block number, but result is nil error", func(t *testing.T) {
		blockNumber := "0x9999999999999999999999999999"
		_, err := alchemy.Core.GetBlock(
			types.BlockHashOrBlockTag{
				BlockTag: blockNumber,
			},
		)

		assert.Error(t, err)
	})
}

/*
func TestAPI_Core_GetTokenBalance(t *testing.T) {
	setting.Network = types.PolygonAmoy // I don't have on eth
	alchemy := gas.NewAlchemy(setting)

	t.Run("get token balance", func(t *testing.T) {
		option := &types.TokenBalanceOption{
			ContractAddresses: []string{"0x41e94eb019c0762f9bfcf9fb1e58725bfb0e7582"},
		}

		_, err := alchemy.Core.GetTokenBalances(
			address, option,
		)

		assert.Nil(t, err)
	})
}
*/
