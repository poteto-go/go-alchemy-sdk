package wallet

import (
	"crypto/ecdsa"
	"errors"
	"math/big"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/poteto-go/go-alchemy-sdk/_fixture/artifacts"
	"github.com/poteto-go/go-alchemy-sdk/constant"
	"github.com/poteto-go/go-alchemy-sdk/ether"
	"github.com/poteto-go/go-alchemy-sdk/gas"
	"github.com/poteto-go/go-alchemy-sdk/types"
	"github.com/poteto-go/go-alchemy-sdk/utils"
	"github.com/stretchr/testify/assert"
)

var testAddrHex = "970e8128ab834e8eac17ab8e3812f010678cf791"
var testAddrHexTo = "970e8128ab834e8eac17ab8e3812f010678cf792"
var testPrivHex = "289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232032"

func createConnectedWallet() *wallet {
	w, _ := New(testPrivHex)

	setting := gas.AlchemySetting{
		ApiKey:  "api-key",
		Network: types.EthMainnet,
	}
	alchemy := gas.NewAlchemy(setting)

	w.Connect(alchemy.GetProvider())

	return w.(*wallet)
}

func TestNewWallet(t *testing.T) {
	t.Run("if can hex to ECDSA, return wallet", func(t *testing.T) {
		// Arrange
		expectedP8Key, _ := crypto.HexToECDSA(testPrivHex)
		expectedPublicKey := expectedP8Key.Public().(*ecdsa.PublicKey)

		// Act
		w, err := New(testPrivHex)

		// Assert
		assert.Nil(t, err)
		assert.Equal(t, expectedP8Key, w.(*wallet).privateKey)
		assert.Equal(t, expectedPublicKey, w.(*wallet).publicKey)
	})

	t.Run("if failed hexToECDSA, return err", func(t *testing.T) {
		// Act
		_, err := New("key")

		// Assert
		assert.Error(t, err)
	})
}

func TestWallet_GetAddress(t *testing.T) {
	// Arrange
	expected := common.HexToAddress(testAddrHex).String()

	// Act
	w, _ := New(testPrivHex)
	addr := w.GetAddress()

	// Assert
	assert.Equal(t, expected, addr)
}

func TestWallet_Connect(t *testing.T) {
	t.Run("can set the provider to wallet", func(t *testing.T) {
		// Arrange
		setting := gas.AlchemySetting{
			ApiKey:  "api-key",
			Network: types.EthMainnet,
		}
		alchemy := gas.NewAlchemy(setting)

		w, _ := New(testPrivHex)

		// Act
		w.Connect(alchemy.GetProvider())

		// Assert
		assert.Equal(t, alchemy.GetProvider(), w.(*wallet).provider)
	})
}

func TestWallet_PendingNonceAt(t *testing.T) {
	t.Run("can get nonce from address", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		w := createConnectedWallet()

		// Mock & Assert
		patches.ApplyMethod(
			reflect.TypeOf(w.provider.Eth()),
			"PendingNonceAt",
			func(_ *ether.Ether, address string) (uint64, error) {
				assert.Equal(t, address, w.GetAddress())
				return uint64(100), nil
			},
		)

		// Act
		nonce, err := w.PendingNonceAt()

		// Assert
		assert.Nil(t, err)
		assert.Equal(t, uint64(100), nonce)
	})

	t.Run("if wallet is not connected, return error", func(t *testing.T) {
		// Arrange
		w, _ := New(testPrivHex)

		// Act
		_, err := w.PendingNonceAt()

		// Assert
		assert.ErrorIs(t, err, constant.ErrWalletIsNotConnected)
	})

	t.Run("cannot get nonce from address, return error", func(t *testing.T) {
		// Arrange
		w := createConnectedWallet()

		// Act
		nonce, err := w.PendingNonceAt()

		// Assert
		assert.Error(t, err)
		assert.Equal(t, uint64(0), nonce)
	})
}

func TestWallet_SignTx(t *testing.T) {
	t.Run("can sign the transaction", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		w := createConnectedWallet()
		reservedNonce := uint64(100)
		txRequest := types.TransactionRequest{
			To:       "0x123",
			ChainID:  big.NewInt(1),
			Nonce:    0,
			GasPrice: big.NewInt(0),
			GasLimit: 1000,
			Value:    "0x123",
			Data:     "0x123",
		}
		estimatedGasPrice := big.NewInt(100)

		// Mock
		patches.ApplyMethod(
			reflect.TypeOf(w),
			"PendingNonceAt",
			func(_ *wallet) (uint64, error) {
				return reservedNonce, nil
			},
		)

		patches.ApplyMethod(
			reflect.TypeOf(w.provider.Eth()),
			"EstimateGas",
			func(_ *ether.Ether, txRequest types.TransactionRequest) (*big.Int, error) {
				return estimatedGasPrice, nil
			},
		)

		// Act
		signedTx, err := w.SignTx(txRequest)

		// Assert
		assert.Nil(t, err)
		assert.Equal(t, txRequest.GasLimit, signedTx.Gas())
		assert.Equal(t, estimatedGasPrice, signedTx.GasPrice())
		assert.Equal(t, reservedNonce, signedTx.Nonce())
		assert.Equal(t, "0x0000000000000000000000000000000000000123", signedTx.To().Hex())
		v, r, s := signedTx.RawSignatureValues()
		assert.Equal(t, v.Cmp(big.NewInt(1)), 0)
		assert.Equal(t, common.BigToHash(r).Hex(), "0x30282c4886900c1309d12e53d1373fc675d905adc84d54e5a2b4afdda2490c07")
		assert.Equal(t, common.BigToHash(s).Hex(), "0x75129dbf83fd1f473464bb1788ae136059d23c2786220da9d9f65ce8cfabb388")
	})

	t.Run("if wallet is not connected, return error", func(t *testing.T) {
		// Arrange
		w, _ := New(testPrivHex)

		// Act
		_, err := w.SignTx(types.TransactionRequest{})

		// Assert
		assert.ErrorIs(t, err, constant.ErrWalletIsNotConnected)
	})

	t.Run("if error occur on PendingAt, return error", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		w := createConnectedWallet()

		// Mock
		patches.ApplyMethod(
			reflect.TypeOf(w),
			"PendingNonceAt",
			func(_ *wallet) (uint64, error) {
				return uint64(0), errors.New("error")
			},
		)

		// Act
		_, err := w.SignTx(types.TransactionRequest{})

		// Assert
		assert.Error(t, err)
	})

	t.Run("if error occur on Eth.EstimateGas", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		w := createConnectedWallet()

		// Mock
		patches.ApplyMethod(
			reflect.TypeOf(w),
			"PendingNonceAt",
			func(_ *wallet) (uint64, error) {
				return uint64(0), nil
			},
		)

		patches.ApplyMethod(
			reflect.TypeOf(w.provider.Eth()),
			"EstimateGas",
			func(_ *ether.Ether, txRequest types.TransactionRequest) (*big.Int, error) {
				return nil, errors.New("error")
			},
		)

		// Act
		_, err := w.SignTx(types.TransactionRequest{})

		// Assert
		assert.Error(t, err)
	})

	t.Run("if gasLimit < gasPrice, return error", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		w := createConnectedWallet()

		// Mock
		patches.ApplyMethod(
			reflect.TypeOf(w),
			"PendingNonceAt",
			func(_ *wallet) (uint64, error) {
				return uint64(0), nil
			},
		)

		patches.ApplyMethod(
			reflect.TypeOf(w.provider.Eth()),
			"EstimateGas",
			func(_ *ether.Ether, txRequest types.TransactionRequest) (*big.Int, error) {
				return big.NewInt(100), nil
			},
		)

		// Act
		_, err := w.SignTx(types.TransactionRequest{
			GasLimit: 0,
		})

		// Assert
		assert.Error(t, err)
	})

	t.Run("if error occur on transform, return error", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		w := createConnectedWallet()

		// Mock
		patches.ApplyMethod(
			reflect.TypeOf(w),
			"PendingNonceAt",
			func(_ *wallet) (uint64, error) {
				return uint64(0), nil
			},
		)

		patches.ApplyMethod(
			reflect.TypeOf(w.provider.Eth()),
			"EstimateGas",
			func(_ *ether.Ether, txRequest types.TransactionRequest) (*big.Int, error) {
				return big.NewInt(1), nil
			},
		)

		patches.ApplyFunc(
			utils.TransformTxRequestToGethTxData,
			func(_ types.TransactionRequest) (*gethTypes.AccessListTx, error) {
				return nil, errors.New("error")
			},
		)

		// Act
		_, err := w.SignTx(types.TransactionRequest{
			GasLimit: 100,
		})

		// Assert
		assert.Error(t, err)
	})

	t.Run("if error on gethTypes.SignTx, return error", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		w := createConnectedWallet()
		reservedNonce := uint64(100)
		txRequest := types.TransactionRequest{
			To:       "0x123",
			ChainID:  big.NewInt(1),
			Nonce:    0,
			GasPrice: big.NewInt(0),
			GasLimit: 1000,
			Value:    "0x123",
			Data:     "0x123",
		}
		estimatedGasPrice := big.NewInt(100)

		// Mock
		patches.ApplyMethod(
			reflect.TypeOf(w),
			"PendingNonceAt",
			func(_ *wallet) (uint64, error) {
				return reservedNonce, nil
			},
		)

		patches.ApplyMethod(
			reflect.TypeOf(w.provider.Eth()),
			"EstimateGas",
			func(_ *ether.Ether, txRequest types.TransactionRequest) (*big.Int, error) {
				return estimatedGasPrice, nil
			},
		)

		patches.ApplyFunc(
			gethTypes.SignTx,
			func(tx *gethTypes.Transaction, s gethTypes.Signer, prv *ecdsa.PrivateKey) (*gethTypes.Transaction, error) {
				return nil, errors.New("error")
			},
		)

		// Act
		_, err := w.SignTx(txRequest)

		// Assert
		assert.Error(t, err)
	})
}

func TestWallet_SendTransaction(t *testing.T) {
	// Arrange
	txRequest := types.TransactionRequest{
		To:       "0x123",
		ChainID:  big.NewInt(1),
		Nonce:    0,
		GasPrice: big.NewInt(0),
		GasLimit: 1000,
		Value:    "0x123",
		Data:     "0x123",
	}
	address := common.HexToAddress("0x123")
	txData := &gethTypes.AccessListTx{
		To:       &address,
		ChainID:  big.NewInt(1),
		Nonce:    0,
		GasPrice: big.NewInt(1),
		Gas:      0,
		Data:     []byte("data"),
	}
	signedTx := gethTypes.NewTx(txData)

	t.Run("can sign and send transaction", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		w := createConnectedWallet()

		// Mock
		patches.ApplyMethod(
			reflect.TypeOf(w),
			"SignTx",
			func(_ *wallet, txRequest types.TransactionRequest) (*gethTypes.Transaction, error) {
				return signedTx, nil
			},
		)
		patches.ApplyMethod(
			reflect.TypeOf(w.provider.Eth()),
			"SendRawTransaction",
			func(_ *ether.Ether, _ *gethTypes.Transaction) error {
				return nil
			},
		)

		// Act
		err := w.SendTransaction(txRequest)

		// Assert
		assert.Nil(t, err)
	})

	t.Run("if wallet is not connected, return error", func(t *testing.T) {
		// Arrange
		w, _ := New(testPrivHex)

		// Act
		err := w.SendTransaction(txRequest)

		// Assert
		assert.ErrorIs(t, err, constant.ErrWalletIsNotConnected)
	})

	t.Run("if error occur on sign-tx, return error", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		w := createConnectedWallet()

		// Mock
		patches.ApplyMethod(
			reflect.TypeOf(w),
			"SignTx",
			func(_ *wallet, txRequest types.TransactionRequest) (*gethTypes.Transaction, error) {
				return nil, errors.New("error")
			},
		)

		// Act
		err := w.SendTransaction(txRequest)

		// Assert
		assert.Error(t, err)
	})

	t.Run("if error occur on send raw transaction, return error", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		w := createConnectedWallet()

		// Mock
		patches.ApplyMethod(
			reflect.TypeOf(w),
			"SignTx",
			func(_ *wallet, txRequest types.TransactionRequest) (*gethTypes.Transaction, error) {
				return signedTx, nil
			},
		)
		patches.ApplyMethod(
			reflect.TypeOf(w.provider.Eth()),
			"SendRawTransaction",
			func(_ *ether.Ether, _ *gethTypes.Transaction) error {
				return errors.New("error")
			},
		)

		// Act
		err := w.SendTransaction(types.TransactionRequest{})

		// Assert
		assert.Error(t, err)
	})
}

func TestWallet_DeployContract(t *testing.T) {
	bytecode := []byte("binary")
	parsed, _ := artifacts.StorageMetaData.GetAbi()

	t.Run("can deploy contract", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		w := createConnectedWallet()
		expectedAddr := common.HexToAddress("0x123")
		expectedTx := gethTypes.NewTx(&gethTypes.LegacyTx{})

		// Mock
		patches.ApplyMethod(
			reflect.TypeOf(w.provider.Eth()),
			"ChainID",
			func(_ *ether.Ether) (*big.Int, error) {
				return big.NewInt(1), nil
			},
		)
		patches.ApplyMethod(
			reflect.TypeOf(w),
			"PendingNonceAt",
			func(_ *wallet) (uint64, error) {
				return uint64(100), nil
			},
		)
		patches.ApplyMethod(
			reflect.TypeOf(w.provider.Eth()),
			"SuggestGasPrice",
			func(_ *ether.Ether) (*big.Int, error) {
				return big.NewInt(100), nil
			},
		)
		patches.ApplyMethod(
			reflect.TypeOf(w.provider.Eth()),
			"DeployContract",
			func(
				_ *ether.Ether,
				opts *bind.TransactOpts,
				abi abi.ABI,
				bytecode []byte,
				params ...any,
			) (common.Address, *gethTypes.Transaction, *bind.BoundContract, error) {
				return expectedAddr, expectedTx, nil, nil
			},
		)

		// Act
		addr, tx, _, err := w.DeployContract(*parsed, bytecode)

		// Assert
		assert.Nil(t, err)
		assert.Equal(t, expectedAddr, addr)
		assert.Equal(t, expectedTx, tx)
	})

	t.Run("if wallet is not connected, return error", func(t *testing.T) {
		w, _ := New(testPrivHex)

		// Act
		_, _, _, err := w.DeployContract(abi.ABI{}, []byte{})

		// Assert
		assert.ErrorIs(t, err, constant.ErrWalletIsNotConnected)
	})

	t.Run("if failed get chainId, return error", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		w := createConnectedWallet()

		// Mock
		patches.ApplyMethod(
			reflect.TypeOf(w.provider.Eth()),
			"ChainID",
			func(_ *ether.Ether) (*big.Int, error) {
				return big.NewInt(0), errors.New("error")
			},
		)

		// Act
		_, _, _, err := w.DeployContract(abi.ABI{}, []byte{})

		// Assert
		assert.Error(t, err)
	})

	t.Run("if failed get pending nonce, return error", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		w := createConnectedWallet()

		// Mock
		patches.ApplyMethod(
			reflect.TypeOf(w.provider.Eth()),
			"ChainID",
			func(_ *ether.Ether) (*big.Int, error) {
				return big.NewInt(1), nil
			},
		)
		patches.ApplyMethod(
			reflect.TypeOf(w),
			"PendingNonceAt",
			func(_ *wallet) (uint64, error) {
				return 0, errors.New("error")
			},
		)

		// Act
		_, _, _, err := w.DeployContract(abi.ABI{}, []byte{})

		// Assert
		assert.Error(t, err)
	})

	t.Run("if failed get suggest gas price, return error", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		w := createConnectedWallet()

		// Mock
		patches.ApplyMethod(
			reflect.TypeOf(w.provider.Eth()),
			"ChainID",
			func(_ *ether.Ether) (*big.Int, error) {
				return big.NewInt(1), nil
			},
		)
		patches.ApplyMethod(
			reflect.TypeOf(w),
			"PendingNonceAt",
			func(_ *wallet) (uint64, error) {
				return uint64(100), nil
			},
		)
		patches.ApplyMethod(
			reflect.TypeOf(w.provider.Eth()),
			"SuggestGasPrice",
			func(_ *ether.Ether) (*big.Int, error) {
				return big.NewInt(0), errors.New("error")
			},
		)

		// Act
		_, _, _, err := w.DeployContract(abi.ABI{}, []byte{})

		// Assert
		assert.Error(t, err)
	})

	t.Run("if failed to deploy contract, return error", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		w := createConnectedWallet()

		// Mock
		patches.ApplyMethod(
			reflect.TypeOf(w.provider.Eth()),
			"ChainID",
			func(_ *ether.Ether) (*big.Int, error) {
				return big.NewInt(1), nil
			},
		)
		patches.ApplyMethod(
			reflect.TypeOf(w),
			"PendingNonceAt",
			func(_ *wallet) (uint64, error) {
				return uint64(100), nil
			},
		)
		patches.ApplyMethod(
			reflect.TypeOf(w.provider.Eth()),
			"SuggestGasPrice",
			func(_ *ether.Ether) (*big.Int, error) {
				return big.NewInt(100), nil
			},
		)
		patches.ApplyMethod(
			reflect.TypeOf(w.provider.Eth()),
			"DeployContract",
			func(
				_ *ether.Ether,
				opts *bind.TransactOpts,
				abi abi.ABI,
				bytecode []byte,
				params ...any,
			) (common.Address, *gethTypes.Transaction, *bind.BoundContract, error) {
				return common.Address{}, nil, nil, errors.New("error")
			},
		)

		// Act
		_, _, _, err := w.DeployContract(abi.ABI{}, []byte{})

		// Assert
		assert.Error(t, err)
	})
}
