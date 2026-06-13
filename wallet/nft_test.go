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
	"github.com/poteto-go/go-alchemy-sdk/encode"
	"github.com/poteto-go/go-alchemy-sdk/ether"
	"github.com/poteto-go/go-alchemy-sdk/namespace"
	"github.com/poteto-go/go-alchemy-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestWallet_NftReadMethods(t *testing.T) {
	contractAddress := "0x1234567890123456789012345678901234567890"
	ownerAddress := "0xE25583099BA105D9ec0A67f5Ae86D90e50036425"
	operatorAddress := "0xAbcdef1234567890abcdef1234567890AbCdEf12"
	tokenId := big.NewInt(1)

	t.Run("can get owner of token", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		w := createConnectedWallet()
		expected := "0xe25583099ba105d9ec0a67f5ae86d90e50036425"

		// Mock
		patches.ApplyMethod(reflect.TypeOf(w.nft), "OwnerOf", func(_ *namespace.Nft, _ string, _ *big.Int) (string, error) {
			return expected, nil
		})

		// Act
		res, err := w.Nft().OwnerOf(contractAddress, tokenId)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("can get token URI", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		w := createConnectedWallet()
		expected := "https://example.com/nft/1"

		// Mock
		patches.ApplyMethod(reflect.TypeOf(w.nft), "TokenURI", func(_ *namespace.Nft, _ string, _ *big.Int) (string, error) {
			return expected, nil
		})

		// Act
		res, err := w.Nft().TokenURI(contractAddress, tokenId)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("can get name", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		w := createConnectedWallet()
		expected := "TestNFT"

		// Mock
		patches.ApplyMethod(reflect.TypeOf(w.nft), "Name", func(_ *namespace.Nft, _ string) (string, error) {
			return expected, nil
		})

		// Act
		res, err := w.Nft().Name(contractAddress)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("can get symbol", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		w := createConnectedWallet()
		expected := "TNFT"

		// Mock
		patches.ApplyMethod(reflect.TypeOf(w.nft), "Symbol", func(_ *namespace.Nft, _ string) (string, error) {
			return expected, nil
		})

		// Act
		res, err := w.Nft().Symbol(contractAddress)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("can get approved address", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		w := createConnectedWallet()
		expected := "0xabcdef1234567890abcdef1234567890abcdef12"

		// Mock
		patches.ApplyMethod(reflect.TypeOf(w.nft), "GetApproved", func(_ *namespace.Nft, _ string, _ *big.Int) (string, error) {
			return expected, nil
		})

		// Act
		res, err := w.Nft().GetApproved(contractAddress, tokenId)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("can get isApprovedForAll", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		w := createConnectedWallet()

		// Mock
		patches.ApplyMethod(reflect.TypeOf(w.nft), "IsApprovedForAll", func(_ *namespace.Nft, _, _, _ string) (bool, error) {
			return true, nil
		})

		// Act
		res, err := w.Nft().IsApprovedForAll(contractAddress, ownerAddress, operatorAddress)

		// Assert
		assert.NoError(t, err)
		assert.True(t, res)
	})

	t.Run("error w/o connect wallet on OwnerOf", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.Nft().OwnerOf(contractAddress, tokenId)

		assert.ErrorIs(t, err, constant.ErrWalletIsNotConnected)
	})

	t.Run("error w/o connect wallet on TokenURI", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.Nft().TokenURI(contractAddress, tokenId)

		assert.ErrorIs(t, err, constant.ErrWalletIsNotConnected)
	})

	t.Run("error w/o connect wallet on Name", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.Nft().Name(contractAddress)

		assert.ErrorIs(t, err, constant.ErrWalletIsNotConnected)
	})

	t.Run("error w/o connect wallet on Symbol", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.Nft().Symbol(contractAddress)

		assert.ErrorIs(t, err, constant.ErrWalletIsNotConnected)
	})

	t.Run("error w/o connect wallet on GetApproved", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.Nft().GetApproved(contractAddress, tokenId)

		assert.ErrorIs(t, err, constant.ErrWalletIsNotConnected)
	})

	t.Run("error w/o connect wallet on IsApprovedForAll", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.Nft().IsApprovedForAll(contractAddress, ownerAddress, operatorAddress)

		assert.ErrorIs(t, err, constant.ErrWalletIsNotConnected)
	})
}

func TestWallet_NftTransferFrom(t *testing.T) {
	contractAddress := "0x1234567890123456789012345678901234567890"
	fromAddress := "0xE25583099BA105D9ec0A67f5Ae86D90e50036425"
	toAddress := "0xAbcdef1234567890abcdef1234567890AbCdEf12"
	tokenId := big.NewInt(1)
	expectedHash := common.HexToHash("0x123")

	t.Run("can transfer from", func(t *testing.T) {
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

		_, err := w.Nft().TransferFrom(context.Background(), contractAddress, fromAddress, toAddress, tokenId, nil)

		assert.Nil(t, err)
	})

	t.Run("handle error on transfer from", func(t *testing.T) {
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

		_, err := w.Nft().TransferFrom(context.Background(), contractAddress, fromAddress, toAddress, tokenId, nil)

		assert.Error(t, err)
	})

	t.Run("error w/o connect wallet", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.Nft().TransferFrom(context.Background(), contractAddress, fromAddress, toAddress, tokenId, nil)

		assert.ErrorIs(t, err, constant.ErrWalletIsNotConnected)
	})
}

func TestWallet_NftTransferFromNoWait(t *testing.T) {
	contractAddress := "0x1234567890123456789012345678901234567890"
	fromAddress := "0xE25583099BA105D9ec0A67f5Ae86D90e50036425"
	toAddress := "0xAbcdef1234567890abcdef1234567890AbCdEf12"
	tokenId := big.NewInt(1)
	expectedHash := common.HexToHash("0x123")

	t.Run("can transfer from no wait & encodes transferFrom calldata", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		w := createConnectedWallet()

		var captured types.TransactionRequest
		patches.ApplyMethod(
			reflect.TypeOf(w),
			"SendTransaction",
			func(_ *wallet, req types.TransactionRequest) (common.Hash, error) {
				captured = req
				return expectedHash, nil
			},
		)

		txHash, err := w.Nft().TransferFromNoWait(contractAddress, fromAddress, toAddress, tokenId, nil)

		assert.Nil(t, err)
		assert.Equal(t, expectedHash, txHash)

		expectedData := encode.ReadCalldata(
			constant.TransferFromFnSignature,
			encode.ABIAddress(fromAddress),
			encode.ABIAddress(toAddress),
			encode.ABIUint256(tokenId),
		)
		assert.Equal(t, expectedData, captured.Data)
		assert.Equal(t, contractAddress, captured.To)
		assert.Equal(t, uint64(300000), captured.GasLimit)
	})

	t.Run("can transfer from no wait w/ custom gasLimit", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		w := createConnectedWallet()

		var captured types.TransactionRequest
		patches.ApplyMethod(
			reflect.TypeOf(w),
			"SendTransaction",
			func(_ *wallet, req types.TransactionRequest) (common.Hash, error) {
				captured = req
				return expectedHash, nil
			},
		)

		gasLimit := uint64(500000)
		txHash, err := w.Nft().TransferFromNoWait(contractAddress, fromAddress, toAddress, tokenId, &gasLimit)

		assert.Nil(t, err)
		assert.Equal(t, expectedHash, txHash)
		assert.Equal(t, uint64(500000), captured.GasLimit)
	})

	t.Run("error w/o connect wallet", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.Nft().TransferFromNoWait(contractAddress, fromAddress, toAddress, tokenId, nil)

		assert.ErrorIs(t, err, constant.ErrWalletIsNotConnected)
	})

	t.Run("invalid from-address returns ErrInvalidAddress", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.Nft().TransferFromNoWait(contractAddress, "invalid", toAddress, tokenId, nil)

		assert.ErrorIs(t, err, constant.ErrInvalidAddress)
	})

	t.Run("invalid to-address returns ErrInvalidAddress", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.Nft().TransferFromNoWait(contractAddress, fromAddress, "invalid", tokenId, nil)

		assert.ErrorIs(t, err, constant.ErrInvalidAddress)
	})

	t.Run("nil tokenId returns ErrNilAmount", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.Nft().TransferFromNoWait(contractAddress, fromAddress, toAddress, nil, nil)

		assert.ErrorIs(t, err, constant.ErrNilAmount)
	})

	t.Run("invalid contractAddress returns ErrInvalidAddress", func(t *testing.T) {
		w := createConnectedWallet()

		_, err := w.Nft().TransferFromNoWait("invalid", fromAddress, toAddress, tokenId, nil)

		assert.ErrorIs(t, err, constant.ErrInvalidAddress)
	})
}

func TestWallet_NftSafeTransferFrom(t *testing.T) {
	contractAddress := "0x1234567890123456789012345678901234567890"
	fromAddress := "0xE25583099BA105D9ec0A67f5Ae86D90e50036425"
	toAddress := "0xAbcdef1234567890abcdef1234567890AbCdEf12"
	tokenId := big.NewInt(1)
	expectedHash := common.HexToHash("0x123")

	t.Run("can safe transfer from", func(t *testing.T) {
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

		_, err := w.Nft().SafeTransferFrom(context.Background(), contractAddress, fromAddress, toAddress, tokenId, nil)

		assert.Nil(t, err)
	})

	t.Run("handle error on safe transfer from", func(t *testing.T) {
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

		_, err := w.Nft().SafeTransferFrom(context.Background(), contractAddress, fromAddress, toAddress, tokenId, nil)

		assert.Error(t, err)
	})

	t.Run("error w/o connect wallet", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.Nft().SafeTransferFrom(context.Background(), contractAddress, fromAddress, toAddress, tokenId, nil)

		assert.ErrorIs(t, err, constant.ErrWalletIsNotConnected)
	})
}

func TestWallet_NftSafeTransferFromNoWait(t *testing.T) {
	contractAddress := "0x1234567890123456789012345678901234567890"
	fromAddress := "0xE25583099BA105D9ec0A67f5Ae86D90e50036425"
	toAddress := "0xAbcdef1234567890abcdef1234567890AbCdEf12"
	tokenId := big.NewInt(1)
	expectedHash := common.HexToHash("0x123")

	t.Run("can safe transfer from no wait & encodes safeTransferFrom calldata", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		w := createConnectedWallet()

		var captured types.TransactionRequest
		patches.ApplyMethod(
			reflect.TypeOf(w),
			"SendTransaction",
			func(_ *wallet, req types.TransactionRequest) (common.Hash, error) {
				captured = req
				return expectedHash, nil
			},
		)

		txHash, err := w.Nft().SafeTransferFromNoWait(contractAddress, fromAddress, toAddress, tokenId, nil)

		assert.Nil(t, err)
		assert.Equal(t, expectedHash, txHash)

		expectedData := encode.ReadCalldata(
			constant.SafeTransferFromFnSignature,
			encode.ABIAddress(fromAddress),
			encode.ABIAddress(toAddress),
			encode.ABIUint256(tokenId),
		)
		assert.Equal(t, expectedData, captured.Data)
	})

	t.Run("error w/o connect wallet", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.Nft().SafeTransferFromNoWait(contractAddress, fromAddress, toAddress, tokenId, nil)

		assert.ErrorIs(t, err, constant.ErrWalletIsNotConnected)
	})

	t.Run("invalid from-address returns ErrInvalidAddress", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.Nft().SafeTransferFromNoWait(contractAddress, "invalid", toAddress, tokenId, nil)

		assert.ErrorIs(t, err, constant.ErrInvalidAddress)
	})

	t.Run("nil tokenId returns ErrNilAmount", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.Nft().SafeTransferFromNoWait(contractAddress, fromAddress, toAddress, nil, nil)

		assert.ErrorIs(t, err, constant.ErrNilAmount)
	})

	t.Run("invalid contractAddress returns ErrInvalidAddress", func(t *testing.T) {
		w := createConnectedWallet()

		_, err := w.Nft().SafeTransferFromNoWait("invalid", fromAddress, toAddress, tokenId, nil)

		assert.ErrorIs(t, err, constant.ErrInvalidAddress)
	})
}

func TestWallet_NftSafeTransferFromWithData(t *testing.T) {
	contractAddress := "0x1234567890123456789012345678901234567890"
	fromAddress := "0xE25583099BA105D9ec0A67f5Ae86D90e50036425"
	toAddress := "0xAbcdef1234567890abcdef1234567890AbCdEf12"
	tokenId := big.NewInt(1)
	payload := []byte{0xde, 0xad, 0xbe, 0xef}
	expectedHash := common.HexToHash("0x123")

	t.Run("can safe transfer from with data", func(t *testing.T) {
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

		_, err := w.Nft().SafeTransferFromWithData(context.Background(), contractAddress, fromAddress, toAddress, tokenId, payload, nil)

		assert.Nil(t, err)
	})

	t.Run("error w/o connect wallet", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.Nft().SafeTransferFromWithData(context.Background(), contractAddress, fromAddress, toAddress, tokenId, payload, nil)

		assert.ErrorIs(t, err, constant.ErrWalletIsNotConnected)
	})
}

func TestWallet_NftSafeTransferFromWithDataNoWait(t *testing.T) {
	contractAddress := "0x1234567890123456789012345678901234567890"
	fromAddress := "0xE25583099BA105D9ec0A67f5Ae86D90e50036425"
	toAddress := "0xAbcdef1234567890abcdef1234567890AbCdEf12"
	tokenId := big.NewInt(1)
	payload := []byte{0xde, 0xad, 0xbe, 0xef}
	expectedHash := common.HexToHash("0x123")

	t.Run("encodes safeTransferFrom(...,bytes) calldata with offset + data", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		w := createConnectedWallet()

		var captured types.TransactionRequest
		patches.ApplyMethod(
			reflect.TypeOf(w),
			"SendTransaction",
			func(_ *wallet, req types.TransactionRequest) (common.Hash, error) {
				captured = req
				return expectedHash, nil
			},
		)

		txHash, err := w.Nft().SafeTransferFromWithDataNoWait(contractAddress, fromAddress, toAddress, tokenId, payload, nil)

		assert.Nil(t, err)
		assert.Equal(t, expectedHash, txHash)

		// head: from, to, tokenId, offset(0x80) -> 4 words; tail: ABIBytes(payload).
		expectedData := encode.ReadCalldata(
			constant.SafeTransferFromWithDataFnSignature,
			encode.ABIAddress(fromAddress),
			encode.ABIAddress(toAddress),
			encode.ABIUint256(tokenId),
			encode.ABIUint256(big.NewInt(4*constant.ABIWordSize)),
			encode.ABIBytes(payload),
		)
		assert.Equal(t, expectedData, captured.Data)
	})

	t.Run("error w/o connect wallet", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.Nft().SafeTransferFromWithDataNoWait(contractAddress, fromAddress, toAddress, tokenId, payload, nil)

		assert.ErrorIs(t, err, constant.ErrWalletIsNotConnected)
	})

	t.Run("invalid to-address returns ErrInvalidAddress", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.Nft().SafeTransferFromWithDataNoWait(contractAddress, fromAddress, "invalid", tokenId, payload, nil)

		assert.ErrorIs(t, err, constant.ErrInvalidAddress)
	})

	t.Run("nil tokenId returns ErrNilAmount", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.Nft().SafeTransferFromWithDataNoWait(contractAddress, fromAddress, toAddress, nil, payload, nil)

		assert.ErrorIs(t, err, constant.ErrNilAmount)
	})

	t.Run("invalid contractAddress returns ErrInvalidAddress", func(t *testing.T) {
		w := createConnectedWallet()

		_, err := w.Nft().SafeTransferFromWithDataNoWait("invalid", fromAddress, toAddress, tokenId, payload, nil)

		assert.ErrorIs(t, err, constant.ErrInvalidAddress)
	})
}
