package wallet

import (
	"context"
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

func TestWallet_ERC1155ReadMethods(t *testing.T) {
	contractAddress := "0x1234567890123456789012345678901234567890"
	account := "0xE25583099BA105D9ec0A67f5Ae86D90e50036425"
	tokenId := big.NewInt(1)

	t.Run("can get balance of token", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		w := createConnectedWallet()
		expected := big.NewInt(42)

		// Mock
		patches.ApplyMethod(reflect.TypeOf(w.erc1155), "BalanceOfToken", func(_ *namespace.Erc1155, _, _ string, _ *big.Int) (*big.Int, error) {
			return expected, nil
		})

		// Act
		res, err := w.ERC1155().BalanceOfToken(contractAddress, account, tokenId)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("can get batch balances", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		w := createConnectedWallet()
		expected := []*big.Int{big.NewInt(1), big.NewInt(2)}

		// Mock
		patches.ApplyMethod(reflect.TypeOf(w.erc1155), "BalanceOfBatch", func(_ *namespace.Erc1155, _ string, _ []string, _ []*big.Int) ([]*big.Int, error) {
			return expected, nil
		})

		// Act
		res, err := w.ERC1155().BalanceOfBatch(contractAddress, []string{account}, []*big.Int{tokenId})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("can get uri", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		w := createConnectedWallet()
		expected := "https://example.com/erc1155/{id}.json"

		// Mock
		patches.ApplyMethod(reflect.TypeOf(w.erc1155), "Uri", func(_ *namespace.Erc1155, _ string, _ *big.Int) (string, error) {
			return expected, nil
		})

		// Act
		res, err := w.ERC1155().Uri(contractAddress, tokenId)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("error w/o connect wallet on BalanceOf", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.ERC1155().BalanceOfToken(contractAddress, account, tokenId)

		assert.ErrorIs(t, err, constant.ErrWalletIsNotConnected)
	})

	t.Run("error w/o connect wallet on BalanceOfBatch", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.ERC1155().BalanceOfBatch(contractAddress, []string{account}, []*big.Int{tokenId})

		assert.ErrorIs(t, err, constant.ErrWalletIsNotConnected)
	})

	t.Run("error w/o connect wallet on Uri", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.ERC1155().Uri(contractAddress, tokenId)

		assert.ErrorIs(t, err, constant.ErrWalletIsNotConnected)
	})
}

func TestWallet_ERC1155SafeTransferFrom(t *testing.T) {
	contractAddress := "0x1234567890123456789012345678901234567890"
	fromAddress := "0xE25583099BA105D9ec0A67f5Ae86D90e50036425"
	toAddress := "0xAbcdef1234567890abcdef1234567890AbCdEf12"
	id := big.NewInt(1)
	amount := big.NewInt(5)
	data := []byte{0xde, 0xad}
	expectedHash := common.HexToHash("0x123")

	t.Run("can safe transfer from (wait)", func(t *testing.T) {
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

		_, err := w.ERC1155().SafeTransferFrom(context.Background(), contractAddress, fromAddress, toAddress, id, amount, data, nil)

		assert.Nil(t, err)
	})

	t.Run("handle error on safe transfer from (wait)", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		w := createConnectedWallet()

		patches.ApplyMethod(
			reflect.TypeOf(w),
			"SendTransaction",
			func(_ *wallet, _ types.TransactionRequest) (common.Hash, error) {
				return common.Hash{}, assert.AnError
			},
		)

		_, err := w.ERC1155().SafeTransferFrom(context.Background(), contractAddress, fromAddress, toAddress, id, amount, data, nil)

		assert.Error(t, err)
	})

	t.Run("error w/o connect wallet", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.ERC1155().SafeTransferFrom(context.Background(), contractAddress, fromAddress, toAddress, id, amount, data, nil)

		assert.ErrorIs(t, err, constant.ErrWalletIsNotConnected)
	})
}

func TestWallet_ERC1155SafeTransferFromNoWait(t *testing.T) {
	contractAddress := "0x1234567890123456789012345678901234567890"
	fromAddress := "0xE25583099BA105D9ec0A67f5Ae86D90e50036425"
	toAddress := "0xAbcdef1234567890abcdef1234567890AbCdEf12"
	id := big.NewInt(1)
	amount := big.NewInt(5)
	data := []byte{0xde, 0xad}
	expectedHash := common.HexToHash("0x123")

	t.Run("encodes ERC-1155 safeTransferFrom calldata", func(t *testing.T) {
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

		txHash, err := w.ERC1155().SafeTransferFromNoWait(contractAddress, fromAddress, toAddress, id, amount, data, nil)

		assert.Nil(t, err)
		assert.Equal(t, expectedHash, txHash)

		offsetWord := encode.ABIUint256(big.NewInt(constant.Erc1155SafeTransferFromHeadSize))
		expectedData := encode.ReadCalldata(
			constant.Erc1155SafeTransferFromFnSignature,
			encode.ABIAddress(fromAddress),
			encode.ABIAddress(toAddress),
			encode.ABIUint256(id),
			encode.ABIUint256(amount),
			offsetWord,
			encode.ABIBytes(data),
		)
		assert.Equal(t, expectedData, captured.Data)
		assert.Equal(t, contractAddress, captured.To)
		assert.Equal(t, uint64(300000), captured.GasLimit)
	})

	t.Run("custom gasLimit is forwarded", func(t *testing.T) {
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
		_, err := w.ERC1155().SafeTransferFromNoWait(contractAddress, fromAddress, toAddress, id, amount, data, &gasLimit)

		assert.Nil(t, err)
		assert.Equal(t, uint64(500000), captured.GasLimit)
	})

	t.Run("invalid from-address returns ErrInvalidAddress", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.ERC1155().SafeTransferFromNoWait(contractAddress, "invalid", toAddress, id, amount, data, nil)

		assert.ErrorIs(t, err, constant.ErrInvalidAddress)
	})

	t.Run("invalid to-address returns ErrInvalidAddress", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.ERC1155().SafeTransferFromNoWait(contractAddress, fromAddress, "invalid", id, amount, data, nil)

		assert.ErrorIs(t, err, constant.ErrInvalidAddress)
	})

	t.Run("nil id returns ErrNilAmount", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.ERC1155().SafeTransferFromNoWait(contractAddress, fromAddress, toAddress, nil, amount, data, nil)

		assert.ErrorIs(t, err, constant.ErrNilAmount)
	})

	t.Run("nil amount returns ErrNilAmount", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.ERC1155().SafeTransferFromNoWait(contractAddress, fromAddress, toAddress, id, nil, data, nil)

		assert.ErrorIs(t, err, constant.ErrNilAmount)
	})

	t.Run("invalid contractAddress returns ErrInvalidAddress", func(t *testing.T) {
		w := createConnectedWallet()

		_, err := w.ERC1155().SafeTransferFromNoWait("invalid", fromAddress, toAddress, id, amount, data, nil)

		assert.ErrorIs(t, err, constant.ErrInvalidAddress)
	})

	t.Run("error w/o connect wallet", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.ERC1155().SafeTransferFromNoWait(contractAddress, fromAddress, toAddress, id, amount, data, nil)

		assert.ErrorIs(t, err, constant.ErrWalletIsNotConnected)
	})
}

func TestWallet_ERC1155SafeBatchTransferFrom(t *testing.T) {
	contractAddress := "0x1234567890123456789012345678901234567890"
	fromAddress := "0xE25583099BA105D9ec0A67f5Ae86D90e50036425"
	toAddress := "0xAbcdef1234567890abcdef1234567890AbCdEf12"
	ids := []*big.Int{big.NewInt(1), big.NewInt(2)}
	amounts := []*big.Int{big.NewInt(5), big.NewInt(10)}
	data := []byte{0xca, 0xfe}
	expectedHash := common.HexToHash("0x456")

	t.Run("can safe batch transfer from (wait)", func(t *testing.T) {
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

		_, err := w.ERC1155().SafeBatchTransferFrom(context.Background(), contractAddress, fromAddress, toAddress, ids, amounts, data, nil)

		assert.Nil(t, err)
	})

	t.Run("handle error on safe batch transfer from (wait)", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		w := createConnectedWallet()

		patches.ApplyMethod(
			reflect.TypeOf(w),
			"SendTransaction",
			func(_ *wallet, _ types.TransactionRequest) (common.Hash, error) {
				return common.Hash{}, assert.AnError
			},
		)

		_, err := w.ERC1155().SafeBatchTransferFrom(context.Background(), contractAddress, fromAddress, toAddress, ids, amounts, data, nil)

		assert.Error(t, err)
	})

	t.Run("error w/o connect wallet", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.ERC1155().SafeBatchTransferFrom(context.Background(), contractAddress, fromAddress, toAddress, ids, amounts, data, nil)

		assert.ErrorIs(t, err, constant.ErrWalletIsNotConnected)
	})
}

func TestWallet_ERC1155SafeBatchTransferFromNoWait(t *testing.T) {
	contractAddress := "0x1234567890123456789012345678901234567890"
	fromAddress := "0xE25583099BA105D9ec0A67f5Ae86D90e50036425"
	toAddress := "0xAbcdef1234567890abcdef1234567890AbCdEf12"
	ids := []*big.Int{big.NewInt(1), big.NewInt(2)}
	amounts := []*big.Int{big.NewInt(5), big.NewInt(10)}
	data := []byte{0xca, 0xfe}
	expectedHash := common.HexToHash("0x456")

	t.Run("encodes ERC-1155 safeBatchTransferFrom calldata", func(t *testing.T) {
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

		txHash, err := w.ERC1155().SafeBatchTransferFromNoWait(contractAddress, fromAddress, toAddress, ids, amounts, data, nil)

		assert.Nil(t, err)
		assert.Equal(t, expectedHash, txHash)

		// Manually build expected calldata:
		// head: from(static), to(static), offsetIds, offsetAmounts, offsetData
		// tails: idsTail, amountsTail, dataTail
		idsTail := encode.ABIUint256Array(ids)
		amountsTail := encode.ABIUint256Array(amounts)
		headSize := constant.Erc1155SafeTransferFromHeadSize
		offsetIds := encode.ABIUint256(big.NewInt(int64(headSize)))
		offsetAmounts := encode.ABIUint256(big.NewInt(int64(headSize + len(idsTail))))
		offsetData := encode.ABIUint256(big.NewInt(int64(headSize + len(idsTail) + len(amountsTail))))
		var args []byte
		args = append(args, encode.ABIAddress(fromAddress)...)
		args = append(args, encode.ABIAddress(toAddress)...)
		args = append(args, offsetIds...)
		args = append(args, offsetAmounts...)
		args = append(args, offsetData...)
		args = append(args, idsTail...)
		args = append(args, amountsTail...)
		args = append(args, encode.ABIBytes(data)...)
		expectedData := encode.ReadCalldata(constant.SafeBatchTransferFromFnSignature, args)
		assert.Equal(t, expectedData, captured.Data)
		assert.Equal(t, contractAddress, captured.To)
		assert.Equal(t, uint64(300000), captured.GasLimit)
	})

	t.Run("custom gasLimit is forwarded", func(t *testing.T) {
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

		gasLimit := uint64(600000)
		_, err := w.ERC1155().SafeBatchTransferFromNoWait(contractAddress, fromAddress, toAddress, ids, amounts, data, &gasLimit)

		assert.Nil(t, err)
		assert.Equal(t, uint64(600000), captured.GasLimit)
	})

	t.Run("mismatched ids/amounts length returns ErrMismatchedArrayLength", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.ERC1155().SafeBatchTransferFromNoWait(contractAddress, fromAddress, toAddress, ids, []*big.Int{big.NewInt(1)}, data, nil)

		assert.ErrorIs(t, err, constant.ErrMismatchedArrayLength)
	})

	t.Run("invalid from-address returns ErrInvalidAddress", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.ERC1155().SafeBatchTransferFromNoWait(contractAddress, "invalid", toAddress, ids, amounts, data, nil)

		assert.ErrorIs(t, err, constant.ErrInvalidAddress)
	})

	t.Run("invalid to-address returns ErrInvalidAddress", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.ERC1155().SafeBatchTransferFromNoWait(contractAddress, fromAddress, "invalid", ids, amounts, data, nil)

		assert.ErrorIs(t, err, constant.ErrInvalidAddress)
	})

	t.Run("nil id in slice returns ErrNilAmount", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.ERC1155().SafeBatchTransferFromNoWait(contractAddress, fromAddress, toAddress, []*big.Int{big.NewInt(1), nil}, amounts, data, nil)

		assert.ErrorIs(t, err, constant.ErrNilAmount)
	})

	t.Run("nil amount in slice returns ErrNilAmount", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.ERC1155().SafeBatchTransferFromNoWait(contractAddress, fromAddress, toAddress, ids, []*big.Int{big.NewInt(5), nil}, data, nil)

		assert.ErrorIs(t, err, constant.ErrNilAmount)
	})

	t.Run("invalid contractAddress returns ErrInvalidAddress", func(t *testing.T) {
		w := createConnectedWallet()

		_, err := w.ERC1155().SafeBatchTransferFromNoWait("invalid", fromAddress, toAddress, ids, amounts, data, nil)

		assert.ErrorIs(t, err, constant.ErrInvalidAddress)
	})

	t.Run("error w/o connect wallet", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.ERC1155().SafeBatchTransferFromNoWait(contractAddress, fromAddress, toAddress, ids, amounts, data, nil)

		assert.ErrorIs(t, err, constant.ErrWalletIsNotConnected)
	})
}
