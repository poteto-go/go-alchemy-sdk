package e2e

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/poteto-go/go-alchemy-sdk/_fixture/artifacts"
	"github.com/poteto-go/go-alchemy-sdk/namespace"
	"github.com/poteto-go/go-alchemy-sdk/types"
	"github.com/poteto-go/go-alchemy-sdk/wallet"
	"github.com/stretchr/testify/assert"
)

const erc1155Uri = "https://example.com/erc1155/{id}.json"

// runErc1155ReadScenario deploys the ERC1155 fixture, mints two token types to
// initAddress, then exercises the Erc1155 read methods. It is shared by the RPC
// and simulated-backend e2e suites since the read behaviour is identical.
func runErc1155ReadScenario(t *testing.T, erc1155 namespace.IErc1155, transact namespace.ITransact, provider types.IAlchemyProvider) {
	t.Helper()

	w, err := wallet.New(initPrivateKey)
	assert.Nil(t, err)
	w.Connect(provider)

	// ERC1155 has a no-arg constructor; no BindDeploymentMetadata needed.
	contractAddress, err := w.DeployContract(context.Background(), &artifacts.ERC1155MetaData)
	assert.Nil(t, err)
	assert.NotEqual(t, contractAddress, common.HexToAddress("0x0"))
	contractHex := contractAddress.Hex()

	erc1155Contract := artifacts.NewERC1155()
	tokenId1 := big.NewInt(1)
	tokenId2 := big.NewInt(2)

	mint := func(id, amount *big.Int) {
		data := erc1155Contract.PackMint(common.HexToAddress(initAddress), id, amount)
		txHash, err := w.SendTransaction(types.TransactionRequest{
			From:     initAddress,
			To:       contractHex,
			Value:    "0x0",
			GasLimit: 300000,
			Data:     data,
		})
		assert.Nil(t, err)
		_, err = transact.WaitMined(context.Background(), txHash.Hex())
		assert.Nil(t, err)
	}
	mint(tokenId1, big.NewInt(10))
	mint(tokenId2, big.NewInt(20))

	t.Run("can get uri via Erc1155 namespace", func(t *testing.T) {
		uri, err := erc1155.Uri(contractHex, tokenId1)
		assert.Nil(t, err)
		assert.Equal(t, erc1155Uri, uri)
	})

	t.Run("can get balanceOf minted token", func(t *testing.T) {
		balance, err := erc1155.BalanceOf(contractHex, initAddress, tokenId1)
		assert.Nil(t, err)
		assert.Equal(t, "10", balance.String())
	})

	t.Run("can get balanceOfBatch", func(t *testing.T) {
		balances, err := erc1155.BalanceOfBatch(
			contractHex,
			[]string{initAddress, initAddress},
			[]*big.Int{tokenId1, tokenId2},
		)
		assert.Nil(t, err)
		assert.Equal(t, 2, len(balances))
		assert.Equal(t, "10", balances[0].String())
		assert.Equal(t, "20", balances[1].String())
	})

	// The wallet Erc1155 namespace delegates to the Erc1155 namespace (verified
	// in detail above), so a couple of representative reads prove the wiring.
	t.Run("can call read methods via wallet Erc1155 namespace", func(t *testing.T) {
		balance, err := w.Erc1155().BalanceOf(contractHex, initAddress, tokenId1)
		assert.Nil(t, err)
		assert.Equal(t, "10", balance.String())

		uri, err := w.Erc1155().Uri(contractHex, tokenId1)
		assert.Nil(t, err)
		assert.Equal(t, erc1155Uri, uri)
	})
}
