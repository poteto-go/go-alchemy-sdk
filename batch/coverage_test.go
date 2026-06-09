package batch_test

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/poteto-go/go-alchemy-sdk/alchemymock"
	"github.com/poteto-go/go-alchemy-sdk/batch"
	"github.com/poteto-go/go-alchemy-sdk/utils"
	"github.com/stretchr/testify/assert"
)

const (
	ownerAddr   = "0x3333333333333333333333333333333333333333"
	spenderAddr = "0x4444444444444444444444444444444444444444"
)

// --- result word helpers -----------------------------------------------------

func uintWord(n int64) string {
	b := make([]byte, 32)
	big.NewInt(n).FillBytes(b)
	return hex.EncodeToString(b)
}

func boolWord(v bool) string {
	if v {
		return uintWord(1)
	}
	return uintWord(0)
}

func addrWord(a string) string {
	b := make([]byte, 32)
	copy(b[12:], common.HexToAddress(a).Bytes())
	return hex.EncodeToString(b)
}

func abiStrWord(s string) string {
	return hex.EncodeToString(utils.EncodeABIString(s))
}

// resp builds a JSON-RPC batch response, assigning sequential ids 1..N.
func resp(results ...string) string {
	parts := make([]string, len(results))
	for i, r := range results {
		parts[i] = fmt.Sprintf(`{"jsonrpc":"2.0","id":%d,"result":"0x%s"}`, i+1, r)
	}
	return "[" + strings.Join(parts, ",") + "]"
}

func TestCoreBatch_AllMethods(t *testing.T) {
	mock := alchemymock.NewAlchemyHttpMock(batchSetting, t)
	defer mock.DeactivateAndReset()

	b := batch.NewBatcher(newBatchEther())
	blockNumber := b.Core.BlockNumber()
	gasPrice := b.Core.GasPrice()
	chainID := b.Core.ChainID()
	peerCount := b.Core.PeerCount()
	balance := b.Core.Balance(walletAddr, "latest")
	code := b.Core.Code(contractAddr, "latest")
	storage := b.Core.StorageAt(contractAddr, "0x0", "latest")

	mock.RegisterBatchResponderOnce(resp(
		"10",   // blockNumber -> 16
		"100",  // gasPrice -> 256
		"1",    // chainId -> 1
		"5",    // peerCount -> 5
		"1234", // balance -> 4660
		"abcd", // code -> 0xabcd
		"00ff", // storageAt -> 00ff
	))

	assert.NoError(t, b.Send())

	assertUnwrap(t, blockNumber, uint64(16))
	assertUnwrapStr(t, gasPrice, "256")
	assertUnwrapStr(t, chainID, "1")
	assertUnwrap(t, peerCount, uint64(5))
	assertUnwrapStr(t, balance, "4660")
	assertUnwrap(t, code, "0xabcd")
	assertUnwrap(t, storage, "00ff")
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

func TestBatch_DecodeErrorSurfacesOnResult(t *testing.T) {
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

func TestBatch_InvalidAddressIsRejectedPerMethod(t *testing.T) {
	b := batch.NewBatcher(newBatchEther())
	bad := "not-an-address"
	var nonce [32]byte

	cases := []struct {
		name string
		run  func() error
	}{
		{"ERC20.BalanceOf", func() error { _, e := b.ERC20.BalanceOf(bad, bad).Unwrap(); return e }},
		{"ERC20.TotalSupply", func() error { _, e := b.ERC20.TotalSupply(bad).Unwrap(); return e }},
		{"ERC20.Allowance", func() error { _, e := b.ERC20.Allowance(bad, bad, bad).Unwrap(); return e }},
		{"ERC20.Name", func() error { _, e := b.ERC20.Name(bad).Unwrap(); return e }},
		{"ERC20.Symbol", func() error { _, e := b.ERC20.Symbol(bad).Unwrap(); return e }},
		{"ERC20.Decimals", func() error { _, e := b.ERC20.Decimals(bad).Unwrap(); return e }},
		{"StableCoin.IsBlacklisted", func() error { _, e := b.StableCoin.IsBlacklisted(bad, bad).Unwrap(); return e }},
		{"StableCoin.Paused", func() error { _, e := b.StableCoin.Paused(bad).Unwrap(); return e }},
		{"StableCoin.Owner", func() error { _, e := b.StableCoin.Owner(bad).Unwrap(); return e }},
		{"StableCoin.MasterMinter", func() error { _, e := b.StableCoin.MasterMinter(bad).Unwrap(); return e }},
		{"StableCoin.Pauser", func() error { _, e := b.StableCoin.Pauser(bad).Unwrap(); return e }},
		{"StableCoin.Blacklister", func() error { _, e := b.StableCoin.Blacklister(bad).Unwrap(); return e }},
		{"StableCoin.Currency", func() error { _, e := b.StableCoin.Currency(bad).Unwrap(); return e }},
		{"StableCoin.Version", func() error { _, e := b.StableCoin.Version(bad).Unwrap(); return e }},
		{"StableCoin.IsMinter", func() error { _, e := b.StableCoin.IsMinter(bad, bad).Unwrap(); return e }},
		{"StableCoin.MinterAllowance", func() error { _, e := b.StableCoin.MinterAllowance(bad, bad).Unwrap(); return e }},
		{"StableCoin.Nonces", func() error { _, e := b.StableCoin.Nonces(bad, bad).Unwrap(); return e }},
		{"StableCoin.DomainSeparator", func() error { _, e := b.StableCoin.DomainSeparator(bad).Unwrap(); return e }},
		{"StableCoin.AuthorizationState", func() error { _, e := b.StableCoin.AuthorizationState(bad, bad, nonce).Unwrap(); return e }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Error(t, tc.run())
		})
	}
}

// --- assert helpers ----------------------------------------------------------

func assertUnwrap[T comparable](t *testing.T, r *batch.Result[T], want T) {
	t.Helper()
	got, err := r.Unwrap()
	assert.NoError(t, err)
	assert.Equal(t, want, got)
}

func assertUnwrapStr(t *testing.T, r *batch.Result[*big.Int], want string) {
	t.Helper()
	got, err := r.Unwrap()
	assert.NoError(t, err)
	assert.Equal(t, want, got.String())
}
