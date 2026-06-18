package batch_test

import (
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/poteto-go/go-alchemy-sdk/alchemymock"
	"github.com/poteto-go/go-alchemy-sdk/batch"
	"github.com/stretchr/testify/assert"
)

func TestNftBatch_AllMethods(t *testing.T) {
	mock := alchemymock.NewAlchemyHttpMock(batchSetting, t)
	defer mock.DeactivateAndReset()

	tokenId := big.NewInt(1)
	b := batch.NewBatcher(newBatchEther())
	balance := b.Nft.BalanceOf(contractAddr, walletAddr)
	ownerOf := b.Nft.OwnerOf(contractAddr, tokenId)
	tokenURI := b.Nft.TokenURI(contractAddr, tokenId)
	name := b.Nft.Name(contractAddr)
	symbol := b.Nft.Symbol(contractAddr)
	getApproved := b.Nft.GetApproved(contractAddr, tokenId)

	mock.RegisterBatchResponderOnce(resp(
		uintWord(5),
		addrWord(ownerAddr),
		abiStrWord("ipfs://token1"),
		abiStrWord("CryptoKitties"),
		abiStrWord("CK"),
		addrWord(spenderAddr),
	))

	assert.NoError(t, b.Send())

	assertUnwrapStr(t, balance, "5")
	assertUnwrap(t, ownerOf, strings.ToLower(common.HexToAddress(ownerAddr).Hex()))
	assertUnwrap(t, tokenURI, "ipfs://token1")
	assertUnwrap(t, name, "CryptoKitties")
	assertUnwrap(t, symbol, "CK")
	assertUnwrap(t, getApproved, strings.ToLower(common.HexToAddress(spenderAddr).Hex()))
}

func TestNftBatch_DecodeError(t *testing.T) {
	mock := alchemymock.NewAlchemyHttpMock(batchSetting, t)
	defer mock.DeactivateAndReset()

	b := batch.NewBatcher(newBatchEther())
	ownerOf := b.Nft.OwnerOf(contractAddr, big.NewInt(1))

	// Too-short response (2 bytes): decode.ABIAddress errors on < 32 bytes.
	mock.RegisterBatchResponderOnce(resp("1234"))

	assert.NoError(t, b.Send())

	_, err := ownerOf.Unwrap()
	assert.Error(t, err)
}

func TestNftBatch_InvalidInputs(t *testing.T) {
	b := batch.NewBatcher(newBatchEther())
	bad := "not-an-address"
	tokenId := big.NewInt(1)

	cases := []struct {
		name string
		run  func() error
	}{
		{"BalanceOf bad contract", func() error { _, e := b.Nft.BalanceOf(bad, walletAddr).Unwrap(); return e }},
		{"BalanceOf bad wallet", func() error { _, e := b.Nft.BalanceOf(contractAddr, bad).Unwrap(); return e }},
		{"OwnerOf bad contract", func() error { _, e := b.Nft.OwnerOf(bad, tokenId).Unwrap(); return e }},
		{"OwnerOf nil tokenId", func() error { _, e := b.Nft.OwnerOf(contractAddr, nil).Unwrap(); return e }},
		{"TokenURI bad contract", func() error { _, e := b.Nft.TokenURI(bad, tokenId).Unwrap(); return e }},
		{"TokenURI nil tokenId", func() error { _, e := b.Nft.TokenURI(contractAddr, nil).Unwrap(); return e }},
		{"Name bad contract", func() error { _, e := b.Nft.Name(bad).Unwrap(); return e }},
		{"Symbol bad contract", func() error { _, e := b.Nft.Symbol(bad).Unwrap(); return e }},
		{"GetApproved bad contract", func() error { _, e := b.Nft.GetApproved(bad, tokenId).Unwrap(); return e }},
		{"GetApproved nil tokenId", func() error { _, e := b.Nft.GetApproved(contractAddr, nil).Unwrap(); return e }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Error(t, tc.run())
		})
	}
}
