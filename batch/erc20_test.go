package batch_test

import (
	"testing"

	"github.com/poteto-go/go-alchemy-sdk/alchemymock"
	"github.com/poteto-go/go-alchemy-sdk/batch"
	"github.com/stretchr/testify/assert"
)

func TestERC20Batch_BalanceOfAndName(t *testing.T) {
	mock := alchemymock.NewAlchemyHttpMock(batchSetting, t)
	defer mock.DeactivateAndReset()

	b := batch.NewBatcher(newBatchEther())
	balance := b.ERC20.BalanceOf(contractAddr, walletAddr)
	name := b.ERC20.Name(contractAddr)

	mock.RegisterBatchResponderOnce(
		`[{"jsonrpc":"2.0","id":1,"result":"0x000000000000000000000000000000000000000000000000000000000000000a"},` +
			`{"jsonrpc":"2.0","id":2,"result":"0x` + abiStringMTK + `"}]`,
	)

	assert.NoError(t, b.Send())

	bal, err := balance.Unwrap()
	assert.NoError(t, err)
	assert.Equal(t, "10", bal.String())

	nm, err := name.Unwrap()
	assert.NoError(t, err)
	assert.Equal(t, "MTK", nm)
}

func TestERC20Batch_AllMethods(t *testing.T) {
	mock := alchemymock.NewAlchemyHttpMock(batchSetting, t)
	defer mock.DeactivateAndReset()

	b := batch.NewBatcher(newBatchEther())
	balance := b.ERC20.BalanceOf(contractAddr, walletAddr)
	supply := b.ERC20.TotalSupply(contractAddr)
	allowance := b.ERC20.Allowance(contractAddr, ownerAddr, spenderAddr)
	name := b.ERC20.Name(contractAddr)
	symbol := b.ERC20.Symbol(contractAddr)
	decimals := b.ERC20.Decimals(contractAddr)

	mock.RegisterBatchResponderOnce(resp(
		uintWord(10),
		uintWord(1000),
		uintWord(5),
		abiStrWord("MTK"),
		abiStrWord("SYM"),
		uintWord(18),
	))

	assert.NoError(t, b.Send())

	assertUnwrapStr(t, balance, "10")
	assertUnwrapStr(t, supply, "1000")
	assertUnwrapStr(t, allowance, "5")
	assertUnwrap(t, name, "MTK")
	assertUnwrap(t, symbol, "SYM")
	assertUnwrap(t, decimals, uint8(18))
}

func TestERC20Batch_DecodeError(t *testing.T) {
	mock := alchemymock.NewAlchemyHttpMock(batchSetting, t)
	defer mock.DeactivateAndReset()

	b := batch.NewBatcher(newBatchEther())
	decimals := b.ERC20.Decimals(contractAddr)

	// 256 overflows uint8 -> DecodeUint8 errors during finalize.
	mock.RegisterBatchResponderOnce(resp(uintWord(256)))

	assert.NoError(t, b.Send())

	_, err := decimals.Unwrap()
	assert.Error(t, err)
}

func TestERC20Batch_InvalidAddresses(t *testing.T) {
	b := batch.NewBatcher(newBatchEther())
	bad := "not-an-address"

	cases := []struct {
		name string
		run  func() error
	}{
		{"BalanceOf", func() error { _, e := b.ERC20.BalanceOf(bad, bad).Unwrap(); return e }},
		{"TotalSupply", func() error { _, e := b.ERC20.TotalSupply(bad).Unwrap(); return e }},
		{"Allowance", func() error { _, e := b.ERC20.Allowance(bad, bad, bad).Unwrap(); return e }},
		{"Name", func() error { _, e := b.ERC20.Name(bad).Unwrap(); return e }},
		{"Symbol", func() error { _, e := b.ERC20.Symbol(bad).Unwrap(); return e }},
		{"Decimals", func() error { _, e := b.ERC20.Decimals(bad).Unwrap(); return e }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Error(t, tc.run())
		})
	}
}
