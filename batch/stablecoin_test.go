package batch_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/poteto-go/go-alchemy-sdk/alchemymock"
	"github.com/poteto-go/go-alchemy-sdk/batch"
	"github.com/stretchr/testify/assert"
)

func TestStableCoinBatch_OwnerAndPaused(t *testing.T) {
	mock := alchemymock.NewAlchemyHttpMock(batchSetting, t)
	defer mock.DeactivateAndReset()

	b := batch.NewBatcher(newBatchEther())
	owner := b.StableCoin.Owner(contractAddr)
	paused := b.StableCoin.Paused(contractAddr)

	mock.RegisterBatchResponderOnce(
		`[{"jsonrpc":"2.0","id":1,"result":"0x000000000000000000000000aabbccddeeff00112233445566778899aabbccdd"},` +
			`{"jsonrpc":"2.0","id":2,"result":"0x0000000000000000000000000000000000000000000000000000000000000001"}]`,
	)

	assert.NoError(t, b.Send())

	o, err := owner.Unwrap()
	assert.NoError(t, err)
	assert.Equal(t, common.HexToAddress("0xaabbccddeeff00112233445566778899aabbccdd"), o)

	p, err := paused.Unwrap()
	assert.NoError(t, err)
	assert.True(t, p)
}

func TestStableCoinBatch_InheritsERC20(t *testing.T) {
	mock := alchemymock.NewAlchemyHttpMock(batchSetting, t)
	defer mock.DeactivateAndReset()

	b := batch.NewBatcher(newBatchEther())
	// Name is an ERC-20 read promoted onto StableCoin via embedding.
	name := b.StableCoin.Name(contractAddr)

	mock.RegisterBatchResponderOnce(
		`[{"jsonrpc":"2.0","id":1,"result":"0x` + abiStringMTK + `"}]`,
	)

	assert.NoError(t, b.Send())

	nm, err := name.Unwrap()
	assert.NoError(t, err)
	assert.Equal(t, "MTK", nm)
}

func TestStableCoinBatch_AllMethods(t *testing.T) {
	mock := alchemymock.NewAlchemyHttpMock(batchSetting, t)
	defer mock.DeactivateAndReset()

	var nonce [32]byte
	nonce[0] = 0xab

	b := batch.NewBatcher(newBatchEther())
	isBlacklisted := b.StableCoin.IsBlacklisted(contractAddr, walletAddr)
	paused := b.StableCoin.Paused(contractAddr)
	owner := b.StableCoin.Owner(contractAddr)
	masterMinter := b.StableCoin.MasterMinter(contractAddr)
	pauser := b.StableCoin.Pauser(contractAddr)
	blacklister := b.StableCoin.Blacklister(contractAddr)
	currency := b.StableCoin.Currency(contractAddr)
	version := b.StableCoin.Version(contractAddr)
	isMinter := b.StableCoin.IsMinter(contractAddr, walletAddr)
	minterAllowance := b.StableCoin.MinterAllowance(contractAddr, walletAddr)
	nonces := b.StableCoin.Nonces(contractAddr, ownerAddr)
	domainSeparator := b.StableCoin.DomainSeparator(contractAddr)
	authState := b.StableCoin.AuthorizationState(contractAddr, ownerAddr, nonce)

	domainWord := uintWord(0)[:62] + "ff" // a non-zero 32-byte word

	mock.RegisterBatchResponderOnce(resp(
		boolWord(true),
		boolWord(false),
		addrWord(ownerAddr),
		addrWord(spenderAddr),
		addrWord(walletAddr),
		addrWord(contractAddr),
		abiStrWord("USD"),
		abiStrWord("1"),
		boolWord(true),
		uintWord(100),
		uintWord(7),
		domainWord,
		boolWord(true),
	))

	assert.NoError(t, b.Send())

	assertUnwrap(t, isBlacklisted, true)
	assertUnwrap(t, paused, false)
	assertUnwrap(t, owner, common.HexToAddress(ownerAddr))
	assertUnwrap(t, masterMinter, common.HexToAddress(spenderAddr))
	assertUnwrap(t, pauser, common.HexToAddress(walletAddr))
	assertUnwrap(t, blacklister, common.HexToAddress(contractAddr))
	assertUnwrap(t, currency, "USD")
	assertUnwrap(t, version, "1")
	assertUnwrap(t, isMinter, true)
	assertUnwrapStr(t, minterAllowance, "100")
	assertUnwrapStr(t, nonces, "7")

	ds, err := domainSeparator.Unwrap()
	assert.NoError(t, err)
	assert.Equal(t, byte(0xff), ds[31])

	assertUnwrap(t, authState, true)
}

func TestStableCoinBatch_InvalidAddresses(t *testing.T) {
	b := batch.NewBatcher(newBatchEther())
	bad := "not-an-address"
	var nonce [32]byte

	cases := []struct {
		name string
		run  func() error
	}{
		{"IsBlacklisted", func() error { _, e := b.StableCoin.IsBlacklisted(bad, bad).Unwrap(); return e }},
		{"Paused", func() error { _, e := b.StableCoin.Paused(bad).Unwrap(); return e }},
		{"Owner", func() error { _, e := b.StableCoin.Owner(bad).Unwrap(); return e }},
		{"MasterMinter", func() error { _, e := b.StableCoin.MasterMinter(bad).Unwrap(); return e }},
		{"Pauser", func() error { _, e := b.StableCoin.Pauser(bad).Unwrap(); return e }},
		{"Blacklister", func() error { _, e := b.StableCoin.Blacklister(bad).Unwrap(); return e }},
		{"Currency", func() error { _, e := b.StableCoin.Currency(bad).Unwrap(); return e }},
		{"Version", func() error { _, e := b.StableCoin.Version(bad).Unwrap(); return e }},
		{"IsMinter", func() error { _, e := b.StableCoin.IsMinter(bad, bad).Unwrap(); return e }},
		{"MinterAllowance", func() error { _, e := b.StableCoin.MinterAllowance(bad, bad).Unwrap(); return e }},
		{"Nonces", func() error { _, e := b.StableCoin.Nonces(bad, bad).Unwrap(); return e }},
		{"DomainSeparator", func() error { _, e := b.StableCoin.DomainSeparator(bad).Unwrap(); return e }},
		{"AuthorizationState", func() error { _, e := b.StableCoin.AuthorizationState(bad, bad, nonce).Unwrap(); return e }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Error(t, tc.run())
		})
	}
}
