package batch_test

import (
	"math/big"
	"testing"

	"github.com/poteto-go/go-alchemy-sdk/alchemymock"
	"github.com/poteto-go/go-alchemy-sdk/batch"
	"github.com/stretchr/testify/assert"
)

func TestERC1155Batch_AllMethods(t *testing.T) {
	mock := alchemymock.NewAlchemyHttpMock(batchSetting, t)
	defer mock.DeactivateAndReset()

	tokenId := big.NewInt(42)
	b := batch.NewBatcher(newBatchEther())
	balance := b.ERC1155.BalanceOfToken(contractAddr, walletAddr, tokenId)
	uri := b.ERC1155.Uri(contractAddr, tokenId)

	mock.RegisterBatchResponderOnce(resp(
		uintWord(7),
		abiStrWord("https://example.com/token/{id}"),
	))

	assert.NoError(t, b.Send())

	assertUnwrapStr(t, balance, "7")
	assertUnwrap(t, uri, "https://example.com/token/{id}")
}

func TestERC1155Batch_InvalidInputs(t *testing.T) {
	b := batch.NewBatcher(newBatchEther())
	bad := "not-an-address"
	tokenId := big.NewInt(1)

	cases := []struct {
		name string
		run  func() error
	}{
		{"BalanceOfToken bad contract", func() error { _, e := b.ERC1155.BalanceOfToken(bad, walletAddr, tokenId).Unwrap(); return e }},
		{"BalanceOfToken bad account", func() error { _, e := b.ERC1155.BalanceOfToken(contractAddr, bad, tokenId).Unwrap(); return e }},
		{"BalanceOfToken nil tokenId", func() error { _, e := b.ERC1155.BalanceOfToken(contractAddr, walletAddr, nil).Unwrap(); return e }},
		{"Uri bad contract", func() error { _, e := b.ERC1155.Uri(bad, tokenId).Unwrap(); return e }},
		{"Uri nil tokenId", func() error { _, e := b.ERC1155.Uri(contractAddr, nil).Unwrap(); return e }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Error(t, tc.run())
		})
	}
}
