package wallet

import (
	"context"
	"errors"
	"math/big"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey"
	"github.com/ethereum/go-ethereum/common"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/poteto-go/go-alchemy-sdk/constant"
	"github.com/poteto-go/go-alchemy-sdk/ether"
	"github.com/poteto-go/go-alchemy-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestWallet_StableCoin(t *testing.T) {
	t.Run("returns WalletStableCoin", func(t *testing.T) {
		w := createConnectedWallet()

		sc := w.StableCoin()

		assert.NotNil(t, sc)
	})
}

func TestWallet_StableCoin_MintNoWait(t *testing.T) {
	contractAddress := "0x1234567890123456789012345678901234567890"
	toAddress := "0xE25583099BA105D9ec0A67f5Ae86D90e50036425"
	expectedHash := common.HexToHash("0x123")

	t.Run("can mint stablecoin", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		w := createConnectedWallet()

		patches.ApplyMethod(
			reflect.TypeOf(w),
			"SendTransaction",
			func(_ *wallet, _ types.TransactionRequest) (common.Hash, error) {
				return expectedHash, nil
			},
		)

		hash, err := w.StableCoin().MintNoWait(contractAddress, toAddress, big.NewInt(100), nil)

		assert.Nil(t, err)
		assert.Equal(t, expectedHash, hash)
	})

	t.Run("handle error on mint", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		w := createConnectedWallet()

		patches.ApplyMethod(
			reflect.TypeOf(w),
			"SendTransaction",
			func(_ *wallet, _ types.TransactionRequest) (common.Hash, error) {
				return common.Hash{}, errors.New("error")
			},
		)

		_, err := w.StableCoin().MintNoWait(contractAddress, toAddress, big.NewInt(100), nil)

		assert.Error(t, err)
	})

	t.Run("error w/o connect wallet", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.StableCoin().MintNoWait(contractAddress, toAddress, big.NewInt(100), nil)

		assert.ErrorIs(t, err, constant.ErrWalletIsNotConnected)
	})
}

func TestWallet_StableCoin_Mint(t *testing.T) {
	contractAddress := "0x1234567890123456789012345678901234567890"
	toAddress := "0xE25583099BA105D9ec0A67f5Ae86D90e50036425"
	expectedHash := common.HexToHash("0x123")

	t.Run("can mint stablecoin and wait", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		w := createConnectedWallet()

		patches.ApplyMethod(
			reflect.TypeOf(w),
			"SendTransaction",
			func(_ *wallet, _ types.TransactionRequest) (common.Hash, error) {
				return expectedHash, nil
			},
		)
		patches.ApplyMethod(
			reflect.TypeOf(w.provider.Eth()),
			"WaitMined",
			func(_ *ether.Ether, _ context.Context, _ common.Hash) (*gethTypes.Receipt, error) {
				return &gethTypes.Receipt{}, nil
			},
		)

		_, err := w.StableCoin().Mint(context.Background(), contractAddress, toAddress, big.NewInt(100), nil)

		assert.Nil(t, err)
	})

	t.Run("error w/o connect wallet", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.StableCoin().Mint(context.Background(), contractAddress, toAddress, big.NewInt(100), nil)

		assert.ErrorIs(t, err, constant.ErrWalletIsNotConnected)
	})
}

func TestWallet_StableCoin_BurnNoWait(t *testing.T) {
	contractAddress := "0x1234567890123456789012345678901234567890"
	expectedHash := common.HexToHash("0x123")

	t.Run("can burn stablecoin", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		w := createConnectedWallet()

		patches.ApplyMethod(
			reflect.TypeOf(w),
			"SendTransaction",
			func(_ *wallet, _ types.TransactionRequest) (common.Hash, error) {
				return expectedHash, nil
			},
		)

		hash, err := w.StableCoin().BurnNoWait(contractAddress, big.NewInt(50), nil)

		assert.Nil(t, err)
		assert.Equal(t, expectedHash, hash)
	})

	t.Run("handle error on burn", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		w := createConnectedWallet()

		patches.ApplyMethod(
			reflect.TypeOf(w),
			"SendTransaction",
			func(_ *wallet, _ types.TransactionRequest) (common.Hash, error) {
				return common.Hash{}, errors.New("error")
			},
		)

		_, err := w.StableCoin().BurnNoWait(contractAddress, big.NewInt(50), nil)

		assert.Error(t, err)
	})

	t.Run("error w/o connect wallet", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.StableCoin().BurnNoWait(contractAddress, big.NewInt(50), nil)

		assert.ErrorIs(t, err, constant.ErrWalletIsNotConnected)
	})
}

func TestWallet_StableCoin_Burn(t *testing.T) {
	contractAddress := "0x1234567890123456789012345678901234567890"
	expectedHash := common.HexToHash("0x123")

	t.Run("can burn stablecoin and wait", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		w := createConnectedWallet()

		patches.ApplyMethod(
			reflect.TypeOf(w),
			"SendTransaction",
			func(_ *wallet, _ types.TransactionRequest) (common.Hash, error) {
				return expectedHash, nil
			},
		)
		patches.ApplyMethod(
			reflect.TypeOf(w.provider.Eth()),
			"WaitMined",
			func(_ *ether.Ether, _ context.Context, _ common.Hash) (*gethTypes.Receipt, error) {
				return &gethTypes.Receipt{}, nil
			},
		)

		_, err := w.StableCoin().Burn(context.Background(), contractAddress, big.NewInt(50), nil)

		assert.Nil(t, err)
	})

	t.Run("error w/o connect wallet", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.StableCoin().Burn(context.Background(), contractAddress, big.NewInt(50), nil)

		assert.ErrorIs(t, err, constant.ErrWalletIsNotConnected)
	})
}
