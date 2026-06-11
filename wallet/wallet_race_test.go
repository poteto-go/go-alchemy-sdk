package wallet

import (
	"math/big"
	"reflect"
	"sync"
	"testing"

	"github.com/agiledragon/gomonkey"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/poteto-go/go-alchemy-sdk/ether"
	"github.com/poteto-go/go-alchemy-sdk/gas"
	"github.com/poteto-go/go-alchemy-sdk/namespace"
	"github.com/poteto-go/go-alchemy-sdk/types"
)

// TestWallet_NoRaceConnectVsReaders verifies that concurrent calls to
// Connect and reader methods do not race on w.provider / w.erc20.
// Must be run with -race to catch unsynchronized field access.
// Refs: https://github.com/poteto-go/go-alchemy-sdk/issues/332
func TestWallet_NoRaceConnectVsReaders(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	setting := gas.AlchemySetting{
		ApiKey:  "api-key",
		Network: types.EthMainnet,
	}
	alchemy, err := gas.NewAlchemy(setting)
	if err != nil {
		t.Fatal(err)
	}
	provider := alchemy.GetProvider()

	// Mock everything readers touch so we don't make network calls.
	patches.ApplyMethod(
		reflect.TypeOf(provider.Eth()),
		"GetBalance",
		func(_ *ether.Ether, _ string, _ string) (*big.Int, error) {
			return big.NewInt(0), nil
		},
	)
	patches.ApplyMethod(
		reflect.TypeOf(provider.Eth()),
		"PendingNonceAt",
		func(_ *ether.Ether, _ string) (uint64, error) {
			return 0, nil
		},
	)
	patches.ApplyMethod(
		reflect.TypeOf(provider.Eth()),
		"EstimateGas",
		func(_ *ether.Ether, _ types.TransactionRequest) (*big.Int, error) {
			return big.NewInt(1), nil
		},
	)
	patches.ApplyMethod(
		reflect.TypeOf(provider.Eth()),
		"ChainID",
		func(_ *ether.Ether) (*big.Int, error) {
			return big.NewInt(1), nil
		},
	)
	patches.ApplyMethod(
		reflect.TypeOf(provider.Eth()),
		"SendRawTransaction",
		func(_ *ether.Ether, _ *gethTypes.Transaction) error {
			return nil
		},
	)
	patches.ApplyMethod(
		reflect.TypeOf(provider.Eth()),
		"ContractCall",
		func(
			_ *ether.Ether,
			_ common.Address,
			_ *bind.CallOpts,
			_ []byte,
			_ func([]byte) (any, error),
		) (any, error) {
			return nil, nil
		},
	)
	patches.ApplyMethod(
		reflect.TypeOf(provider.Eth()),
		"DeployContract",
		func(
			_ *ether.Ether,
			_ *bind.TransactOpts,
			_ *bind.MetaData,
		) (*bind.DeploymentResult, error) {
			return &bind.DeploymentResult{
				Txs: map[string]*gethTypes.Transaction{},
			}, nil
		},
	)
	patches.ApplyMethod(
		reflect.TypeOf(provider.Eth()),
		"ContractTransact",
		func(
			_ *ether.Ether,
			_ *bind.TransactOpts,
			_ string,
			_ []byte,
		) (*gethTypes.Transaction, error) {
			return gethTypes.NewTx(&gethTypes.LegacyTx{}), nil
		},
	)

	w, _ := New(testPrivHex)
	w.Connect(provider) // ensure erc20 is set before readers start

	contractAddress := "0x1234567890123456789012345678901234567890"
	callData := []byte("call data")
	callOpts := &bind.CallOpts{}
	unpack := func(b []byte) (any, error) { return nil, nil }

	const iterations = 200
	var wg sync.WaitGroup

	// Writer: re-Connect repeatedly.
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			w.Connect(provider)
		}
	}()

	// Reader: GetBalance.
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			_, _ = w.GetBalance()
		}
	}()

	// Reader: PendingNonceAt.
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			_, _ = w.PendingNonceAt()
		}
	}()

	// Reader: ContractCall.
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			_, _ = w.ContractCall(contractAddress, callOpts, callData, unpack)
		}
	}()

	// Reader: ContractTransactNoWait (also touches getOrCreateAuth).
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			_, _ = w.ContractTransactNoWait(contractAddress, callData)
		}
	}()

	wg.Wait()
}

// TestWallet_NoRaceConnectVsERC20Readers verifies that concurrent Connect
// calls do not race with walletERC20 readers that touch w.provider / w.erc20.
// Refs: https://github.com/poteto-go/go-alchemy-sdk/issues/332
func TestWallet_NoRaceConnectVsERC20Readers(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	setting := gas.AlchemySetting{
		ApiKey:  "api-key",
		Network: types.EthMainnet,
	}
	alchemy, err := gas.NewAlchemy(setting)
	if err != nil {
		t.Fatal(err)
	}
	provider := alchemy.GetProvider()

	patches.ApplyMethod(
		reflect.TypeOf(&namespace.ERC20{}),
		"BalanceOf",
		func(_ *namespace.ERC20, _ string, _ string) (*big.Int, error) {
			return big.NewInt(0), nil
		},
	)
	patches.ApplyMethod(
		reflect.TypeOf(&namespace.ERC20{}),
		"TotalSupply",
		func(_ *namespace.ERC20, _ string) (*big.Int, error) {
			return big.NewInt(0), nil
		},
	)
	patches.ApplyMethod(
		reflect.TypeOf(&namespace.ERC20{}),
		"Allowance",
		func(_ *namespace.ERC20, _, _, _ string) (*big.Int, error) {
			return big.NewInt(0), nil
		},
	)
	patches.ApplyMethod(
		reflect.TypeOf(&namespace.ERC20{}),
		"Name",
		func(_ *namespace.ERC20, _ string) (string, error) {
			return "", nil
		},
	)
	patches.ApplyMethod(
		reflect.TypeOf(&namespace.ERC20{}),
		"Symbol",
		func(_ *namespace.ERC20, _ string) (string, error) {
			return "", nil
		},
	)
	patches.ApplyMethod(
		reflect.TypeOf(&namespace.ERC20{}),
		"Decimals",
		func(_ *namespace.ERC20, _ string) (uint8, error) {
			return 0, nil
		},
	)

	w, _ := New(testPrivHex)
	w.Connect(provider)

	contractAddress := "0x1234567890123456789012345678901234567890"

	const iterations = 200
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			w.Connect(provider)
		}
	}()

	readers := []func(){
		func() { _, _ = w.ERC20().BalanceOf(contractAddress) },
		func() { _, _ = w.ERC20().TotalSupply(contractAddress) },
		func() { _, _ = w.ERC20().Allowance(contractAddress, contractAddress, contractAddress) },
		func() { _, _ = w.ERC20().Name(contractAddress) },
		func() { _, _ = w.ERC20().Symbol(contractAddress) },
		func() { _, _ = w.ERC20().Decimals(contractAddress) },
	}

	for _, r := range readers {
		wg.Add(1)
		go func(read func()) {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				read()
			}
		}(r)
	}

	wg.Wait()
}
