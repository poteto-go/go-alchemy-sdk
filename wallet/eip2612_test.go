package wallet

import (
	"context"
	"errors"
	"math/big"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/poteto-go/go-alchemy-sdk/constant"
	"github.com/poteto-go/go-alchemy-sdk/ether"
	"github.com/poteto-go/go-alchemy-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestWallet_EIP2612(t *testing.T) {
	t.Run("returns WalletEIP2612", func(t *testing.T) {
		w := createConnectedWallet()

		eip2612 := w.EIP2612()

		assert.NotNil(t, eip2612)
	})
}

func TestWallet_EIP2612_PermitNoWait(t *testing.T) {
	contractAddress := "0x1234567890123456789012345678901234567890"
	ownerAddress := "0xE25583099BA105D9ec0A67f5Ae86D90e50036425"
	spenderAddress := "0xabcdef1234567890abcdef1234567890abcdef12"
	expectedHash := common.HexToHash("0x123")
	var r, s [32]byte

	t.Run("can submit permit", func(t *testing.T) {
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

		hash, err := w.EIP2612().PermitNoWait(contractAddress, ownerAddress, spenderAddress, big.NewInt(100), big.NewInt(9999999), 27, r, s, nil)

		assert.Nil(t, err)
		assert.Equal(t, expectedHash, hash)
	})

	t.Run("handle error on permit", func(t *testing.T) {
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

		_, err := w.EIP2612().PermitNoWait(contractAddress, ownerAddress, spenderAddress, big.NewInt(100), big.NewInt(9999999), 27, r, s, nil)

		assert.Error(t, err)
	})

	t.Run("error w/o connect wallet", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.EIP2612().PermitNoWait(contractAddress, ownerAddress, spenderAddress, big.NewInt(100), big.NewInt(9999999), 27, r, s, nil)

		assert.ErrorIs(t, err, constant.ErrWalletIsNotConnected)
	})
}

func TestWallet_EIP2612_Permit(t *testing.T) {
	contractAddress := "0x1234567890123456789012345678901234567890"
	ownerAddress := "0xE25583099BA105D9ec0A67f5Ae86D90e50036425"
	spenderAddress := "0xabcdef1234567890abcdef1234567890abcdef12"
	var r, s [32]byte

	t.Run("can permit and wait", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		w := createConnectedWallet()
		expected := &gethTypes.Receipt{TxHash: common.HexToHash("0x123")}

		patches.ApplyMethod(
			reflect.TypeOf(w),
			"SendTransaction",
			func(_ *wallet, _ types.TransactionRequest) (common.Hash, error) {
				return expected.TxHash, nil
			},
		)
		patches.ApplyMethod(
			reflect.TypeOf(w.snapshot().Eth()),
			"WaitMined",
			func(_ *ether.Ether, _ context.Context, _ common.Hash) (*gethTypes.Receipt, error) {
				return expected, nil
			},
		)

		receipt, err := w.EIP2612().Permit(context.Background(), contractAddress, ownerAddress, spenderAddress, big.NewInt(100), big.NewInt(9999999), 27, r, s, nil)

		assert.Nil(t, err)
		assert.Equal(t, expected, receipt)
	})

	t.Run("error w/o connect wallet", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.EIP2612().Permit(context.Background(), contractAddress, ownerAddress, spenderAddress, big.NewInt(100), big.NewInt(9999999), 27, r, s, nil)

		assert.ErrorIs(t, err, constant.ErrWalletIsNotConnected)
	})
}

func TestWallet_EIP2612_Nonces(t *testing.T) {
	contractAddress := "0x1234567890123456789012345678901234567890"
	ownerAddress := "0xE25583099BA105D9ec0A67f5Ae86D90e50036425"

	t.Run("returns nonce from namespace", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		w := createConnectedWallet()
		expected := make([]byte, 32)
		expected[31] = 3

		patches.ApplyMethod(reflect.TypeOf(w.snapshot().Eth()), "CallContract", func(_ *ether.Ether, _ ethereum.CallMsg, _ string) ([]byte, error) {
			return expected, nil
		})

		result, err := w.EIP2612().Nonces(contractAddress, ownerAddress)

		assert.NoError(t, err)
		assert.Equal(t, int64(3), result.Int64())
	})

	t.Run("error w/o connect wallet", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.EIP2612().Nonces(contractAddress, ownerAddress)

		assert.ErrorIs(t, err, constant.ErrWalletIsNotConnected)
	})
}

func TestWallet_EIP2612_DomainSeparator(t *testing.T) {
	contractAddress := "0x1234567890123456789012345678901234567890"

	t.Run("returns domain separator from namespace", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		w := createConnectedWallet()
		raw := make([]byte, 32)
		raw[0] = 0xde
		raw[31] = 0xad

		patches.ApplyMethod(reflect.TypeOf(w.snapshot().Eth()), "CallContract", func(_ *ether.Ether, _ ethereum.CallMsg, _ string) ([]byte, error) {
			return raw, nil
		})

		result, err := w.EIP2612().DomainSeparator(contractAddress)

		assert.NoError(t, err)
		assert.Equal(t, raw, result[:])
	})

	t.Run("error w/o connect wallet", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.EIP2612().DomainSeparator(contractAddress)

		assert.ErrorIs(t, err, constant.ErrWalletIsNotConnected)
	})
}
