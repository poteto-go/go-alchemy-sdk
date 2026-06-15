package e2e

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/poteto-go/go-alchemy-sdk/_fixture/artifacts"
	"github.com/poteto-go/go-alchemy-sdk/namespace"
	"github.com/poteto-go/go-alchemy-sdk/types"
	"github.com/stretchr/testify/assert"
)

func mintERC1155(t *testing.T, contract *artifacts.ERC1155, w types.Wallet, contractHex string, transact namespace.ITransact, id, amount *big.Int) {
	t.Helper()
	data := contract.PackMint(common.HexToAddress(initAddress), id, amount)
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
