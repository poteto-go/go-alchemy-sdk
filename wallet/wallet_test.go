package wallet

import (
	"crypto/ecdsa"
	"errors"
	"math/big"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/poteto-go/go-alchemy-sdk/_fixture/artifacts"
	"github.com/poteto-go/go-alchemy-sdk/constant"
	"github.com/poteto-go/go-alchemy-sdk/ether"
	"github.com/poteto-go/go-alchemy-sdk/gas"
	"github.com/poteto-go/go-alchemy-sdk/internal"
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

func TestWallet_GetBalance(t *testing.T) {
	t.Run("can get balance", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		w := createConnectedWallet()

		// Mock & Assert
		patches.ApplyMethod(
			reflect.TypeOf(w.provider.Eth()),
			"GetBalance",
			func(_ *ether.Ether, address string, blockTag string) (*big.Int, error) {
				assert.Equal(t, address, w.GetAddress())
				assert.Equal(t, blockTag, "latest")
				return big.NewInt(100), nil
			},
		)

		// Act
		balance, err := w.GetBalance()

		// Assert
		assert.Nil(t, err)
		assert.Equal(t, balance.Cmp(big.NewInt(100)), 0)
	})

	t.Run("if wallet is not connected, return error", func(t *testing.T) {
		// Arrange
		w, _ := New(testPrivHex)

		// Act
		balance, err := w.GetBalance()

		// Assert
		assert.ErrorIs(t, err, constant.ErrWalletIsNotConnected)
		assert.Nil(t, balance)
	})

	t.Run("if failed get balance, return error", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		w := createConnectedWallet()

		// Mock & Assert
		patches.ApplyMethod(
			reflect.TypeOf(w.provider.Eth()),
			"GetBalance",
			func(_ *ether.Ether, address string, blockTag string) (*big.Int, error) {
				assert.Equal(t, address, w.GetAddress())
				assert.Equal(t, blockTag, "latest")
				return nil, errors.New("error")
			},
		)

		// Act
		balance, err := w.GetBalance()

		// Assert
		assert.Error(t, err)
		assert.Nil(t, balance)
	})
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

func TestWallet_Provider(t *testing.T) {
	t.Run("can get provider", func(t *testing.T) {
		// Arrange
		w := createConnectedWallet()

		// Act
		provider := w.Provider()

		// Assert
		assert.NotNil(t, provider)
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

		patches.ApplyMethod(
			reflect.TypeOf(w.provider.Eth()),
			"ChainID",
			func(_ *ether.Ether) (*big.Int, error) {
				return big.NewInt(1), nil
			},
		)

		// Act
		signedTx, err := w.SignTx(txRequest)

		// Assert
		assert.Nil(t, err)
		assert.Equal(t, txRequest.ChainID, signedTx.ChainId())
		assert.Equal(t, txRequest.GasLimit, signedTx.Gas())
		assert.Equal(t, estimatedGasPrice, signedTx.GasPrice())
		assert.Equal(t, reservedNonce, signedTx.Nonce())
		assert.Equal(t, "0x0000000000000000000000000000000000000123", signedTx.To().Hex())
		v, r, s := signedTx.RawSignatureValues()
		assert.Equal(t, v.Cmp(big.NewInt(3)), 1)
		assert.Equal(t, common.BigToHash(r).Hex(), "0x75ea94b2f16c4cd67f1a2cc8f35659a434ea7b6f7704dd062f6686d865547fab")
		assert.Equal(t, common.BigToHash(s).Hex(), "0x2af447088ccf07c8742b3cd428d77580ab21e8e88340fbfb20ae7a747d86ca21")
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

	t.Run("if error occur on Eth.ChainId", func(t *testing.T) {
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

		patches.ApplyMethod(
			reflect.TypeOf(w.provider.Eth()),
			"ChainID",
			func(_ *ether.Ether) (*big.Int, error) {
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

		patches.ApplyMethod(
			reflect.TypeOf(w.provider.Eth()),
			"ChainID",
			func(_ *ether.Ether) (*big.Int, error) {
				return big.NewInt(1), nil
			},
		)

		patches.ApplyFunc(
			utils.TransformTxRequestToGethTxData,
			func(_ types.TransactionRequest) (gethTypes.TxData, error) {
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
		patches.ApplyMethod(
			reflect.TypeOf(w.provider.Eth()),
			"ChainID",
			func(_ *ether.Ether) (*big.Int, error) {
				return big.NewInt(1), nil
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

	t.Run("if failed to sign, return error", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		w := createConnectedWallet()
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
				return 0, nil
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
		txHash, err := w.SendTransaction(txRequest)

		// Assert
		assert.Nil(t, err)
		assert.Equal(t, signedTx.Hash(), txHash)
	})

	t.Run("if wallet is not connected, return error", func(t *testing.T) {
		// Arrange
		w, _ := New(testPrivHex)

		// Act
		_, err := w.SendTransaction(txRequest)

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
		_, err := w.SendTransaction(txRequest)

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
		_, err := w.SendTransaction(types.TransactionRequest{})

		// Assert
		assert.Error(t, err)
	})
}

func TestWallet_DeployContract(t *testing.T) {
	metaData := artifacts.PotetoStorageMetaData

	t.Run("can deploy contract", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		w := createConnectedWallet()
		expectedAddr := common.HexToAddress("0x123")
		expectedTx := gethTypes.NewTx(&gethTypes.LegacyTx{})
		expectedDeployRes := &bind.DeploymentResult{
			Txs: map[string]*gethTypes.Transaction{
				metaData.ID: expectedTx,
			},
		}

		// Mock
		patches.ApplyMethod(
			reflect.TypeOf(w.provider.Eth()),
			"ChainID",
			func(_ *ether.Ether) (*big.Int, error) {
				return big.NewInt(1), nil
			},
		)
		patches.ApplyMethod(
			reflect.TypeOf(w.provider.Eth()),
			"DeployContract",
			func(
				_ *ether.Ether,
				auth *bind.TransactOpts,
				metaData *bind.MetaData,
			) (*bind.DeploymentResult, error) {
				return expectedDeployRes, nil
			},
		)
		patches.ApplyMethod(
			reflect.TypeOf(w.provider.Eth()),
			"WaitDeployed",
			func(
				_ *ether.Ether,
				_ common.Hash,
			) (common.Address, error) {
				return expectedAddr, nil
			},
		)

		// Act
		addr, err := w.DeployContract(&metaData)

		// Assert
		assert.Nil(t, err)
		assert.Equal(t, expectedAddr, addr)
	})

	t.Run("if failed to wait deployed, return error", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		w := createConnectedWallet()
		expectedTx := gethTypes.NewTx(&gethTypes.LegacyTx{})
		expectedDeployRes := &bind.DeploymentResult{
			Txs: map[string]*gethTypes.Transaction{
				metaData.ID: expectedTx,
			},
		}

		// Mock
		patches.ApplyMethod(
			reflect.TypeOf(w),
			"DeployContractNoWait",
			func(_ *wallet, metaData *bind.MetaData) (*bind.DeploymentResult, error) {
				return expectedDeployRes, nil
			},
		)
		patches.ApplyMethod(
			reflect.TypeOf(w.provider.Eth()),
			"WaitDeployed",
			func(
				_ *ether.Ether,
				_ common.Hash,
			) (common.Address, error) {
				return common.Address{}, errors.New("error")
			},
		)

		// Act
		_, err := w.DeployContract(&metaData)

		// Assert
		assert.Error(t, err)
	})

	t.Run("if wallet is not connected, return error", func(t *testing.T) {
		w, _ := New(testPrivHex)

		// Act
		_, err := w.DeployContract(&metaData)

		// Assert
		assert.ErrorIs(t, err, constant.ErrWalletIsNotConnected)
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
			"DeployContractNoWait",
			func(_ *wallet, metaData *bind.MetaData) (*bind.DeploymentResult, error) {
				return nil, errors.New("error")
			},
		)

		// Act
		_, err := w.DeployContract(&metaData)

		// Assert
		assert.Error(t, err)
	})
}

func TestWallet_DeployContractNoWait(t *testing.T) {
	metaData := artifacts.PotetoStorageMetaData

	t.Run("can transact of deployment contract", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		w := createConnectedWallet()
		expectedTx := gethTypes.NewTx(&gethTypes.LegacyTx{})
		expectedDeployRes := &bind.DeploymentResult{
			Txs: map[string]*gethTypes.Transaction{
				metaData.ID: expectedTx,
			},
		}

		// Mock
		patches.ApplyMethod(
			reflect.TypeOf(w.provider.Eth()),
			"ChainID",
			func(_ *ether.Ether) (*big.Int, error) {
				return big.NewInt(1), nil
			},
		)
		patches.ApplyMethod(
			reflect.TypeOf(w.provider.Eth()),
			"DeployContract",
			func(
				_ *ether.Ether,
				auth *bind.TransactOpts,
				metaData *bind.MetaData,
			) (*bind.DeploymentResult, error) {
				return expectedDeployRes, nil
			},
		)

		// Act
		res, err := w.DeployContractNoWait(&metaData)

		// Assert
		assert.Nil(t, err)
		assert.Equal(t, res, expectedDeployRes)
	})

	t.Run("if wallet is not connected, return error", func(t *testing.T) {
		w, _ := New(testPrivHex)

		// Act
		_, err := w.DeployContractNoWait(&metaData)

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
		_, err := w.DeployContractNoWait(&metaData)

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
			reflect.TypeOf(w.provider.Eth()),
			"DeployContract",
			func(
				_ *ether.Ether,
				auth *bind.TransactOpts,
				metaData *bind.MetaData,
			) (*bind.DeploymentResult, error) {
				return nil, errors.New("error")
			},
		)

		// Act
		_, err := w.DeployContractNoWait(&metaData)

		// Assert
		assert.Error(t, err)
	})

	t.Run("if chain doesn't support EIP-1559, use call suggestGasPrice", func(t *testing.T) {
		t.Run("can deploy contract", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Arrange
			w := createConnectedWallet()
			isCalledSuggestedGasPrice := false
			expectedTx := gethTypes.NewTx(&gethTypes.LegacyTx{})
			expectedDeployRes := &bind.DeploymentResult{
				Txs: map[string]*gethTypes.Transaction{
					metaData.ID: expectedTx,
				},
			}

			// Mock
			patches.ApplyMethod(
				reflect.TypeOf(w.provider.Eth()),
				"ChainID",
				func(_ *ether.Ether) (*big.Int, error) {
					return big.NewInt(internal.ChainListNotSupportEIP1559[0]), nil
				},
			)
			patches.ApplyMethod(
				reflect.TypeOf(w.provider.Eth()),
				"SuggestGasPrice",
				func(_ *ether.Ether) (*big.Int, error) {
					isCalledSuggestedGasPrice = true
					return big.NewInt(1), nil
				},
			)
			patches.ApplyMethod(
				reflect.TypeOf(w.provider.Eth()),
				"DeployContract",
				func(
					_ *ether.Ether,
					auth *bind.TransactOpts,
					metaData *bind.MetaData,
				) (*bind.DeploymentResult, error) {
					return expectedDeployRes, nil
				},
			)

			// Act
			res, err := w.DeployContractNoWait(&metaData)

			// Assert
			assert.Nil(t, err)
			assert.Equal(t, res, expectedDeployRes)
			assert.True(t, isCalledSuggestedGasPrice)
		})

		t.Run("if failed on suggestedGasPrice", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Arrange
			w := createConnectedWallet()

			// Mock
			patches.ApplyMethod(
				reflect.TypeOf(w.provider.Eth()),
				"ChainID",
				func(_ *ether.Ether) (*big.Int, error) {
					return big.NewInt(internal.ChainListNotSupportEIP1559[0]), nil
				},
			)
			patches.ApplyMethod(
				reflect.TypeOf(w.provider.Eth()),
				"SuggestGasPrice",
				func(_ *ether.Ether) (*big.Int, error) {
					return nil, errors.New("error")
				},
			)

			// Act
			_, err := w.DeployContractNoWait(&metaData)

			// Assert
			assert.Error(t, err)
		})
	})
}

func TestWallet_ContractTransact(t *testing.T) {
	contract := artifacts.NewPotetoStorage()
	contractAddress := "0x1234567890123456789012345678901234567890"
	data := []byte("test data")

	t.Run("can transact with contract & wait to be mined", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		w := createConnectedWallet()
		expectedReceipt := &gethTypes.Receipt{
			Status: 1,
		}
		txData := &gethTypes.AccessListTx{
			To:       &common.Address{},
			ChainID:  big.NewInt(1),
			Nonce:    0,
			GasPrice: big.NewInt(1),
			Gas:      0,
			Data:     []byte("data"),
		}
		expectedTx := gethTypes.NewTx(txData)

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
			"ContractTransactNoWait",
			func(
				_ *wallet,
				_ types.ContractInstance,
				_ string,
				_ []byte,
			) (*gethTypes.Transaction, error) {
				return expectedTx, nil
			},
		)
		patches.ApplyMethod(
			reflect.TypeOf(w.provider.Eth()),
			"WaitMined",
			func(_ *ether.Ether, _ common.Hash) (*gethTypes.Receipt, error) {
				return expectedReceipt, nil
			},
		)

		// Act
		txReceipt, err := w.ContractTransact(contract, contractAddress, data)

		//Assert
		assert.Nil(t, err)
		assert.Equal(t, txReceipt, expectedReceipt)
	})

	t.Run("if the wallet isn't connected, return err", func(t *testing.T) {
		w, _ := New(testPrivHex)

		// Act
		_, err := w.ContractTransact(contract, contractAddress, data)

		// Assert
		assert.ErrorIs(t, err, constant.ErrWalletIsNotConnected)
	})

	t.Run("if failed to transact, return error", func(t *testing.T) {
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
			"ContractTransactNoWait",
			func(
				_ *wallet,
				_ types.ContractInstance,
				_ string,
				_ []byte,
			) (*gethTypes.Transaction, error) {
				return nil, errors.New("error")
			},
		)

		// Act
		_, err := w.ContractTransact(contract, contractAddress, data)

		//Assert
		assert.Error(t, err)
	})

	t.Run("if it failed to wait to be mined, return err", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		w := createConnectedWallet()
		txData := &gethTypes.AccessListTx{
			To:       &common.Address{},
			ChainID:  big.NewInt(1),
			Nonce:    0,
			GasPrice: big.NewInt(1),
			Gas:      0,
			Data:     []byte("data"),
		}
		expectedTx := gethTypes.NewTx(txData)

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
			"ContractTransactNoWait",
			func(
				_ *wallet,
				_ types.ContractInstance,
				_ string,
				_ []byte,
			) (*gethTypes.Transaction, error) {
				return expectedTx, nil
			},
		)
		patches.ApplyMethod(
			reflect.TypeOf(w.provider.Eth()),
			"WaitMined",
			func(_ *ether.Ether, _ common.Hash) (*gethTypes.Receipt, error) {
				return nil, errors.New("error")
			},
		)

		// Act
		_, err := w.ContractTransact(contract, contractAddress, data)

		// Assert
		assert.Error(t, err)
	})
}

func TestWallet_ContractTransactNoWait(t *testing.T) {
	contract := artifacts.NewPotetoStorage()
	contractAddress := "0x1234567890123456789012345678901234567890"
	data := []byte("test data")

	t.Run("can transact with contract", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		w := createConnectedWallet()
		txData := &gethTypes.AccessListTx{
			To:       &common.Address{},
			ChainID:  big.NewInt(1),
			Nonce:    0,
			GasPrice: big.NewInt(1),
			Gas:      0,
			Data:     []byte("data"),
		}
		expectedTx := gethTypes.NewTx(txData)

		// Mock
		patches.ApplyMethod(
			reflect.TypeOf(w.provider.Eth()),
			"ChainID",
			func(_ *ether.Ether) (*big.Int, error) {
				return big.NewInt(1), nil
			},
		)
		patches.ApplyMethod(
			reflect.TypeOf(w.provider.Eth()),
			"ContractTransact",
			func(
				_ *ether.Ether,
				auth *bind.TransactOpts,
				contract types.ContractInstance,
				contractAddress string,
				data []byte,
			) (*gethTypes.Transaction, error) {
				return expectedTx, nil
			},
		)

		// Act
		tx, err := w.ContractTransactNoWait(contract, contractAddress, data)

		// Assert
		assert.Nil(t, err)
		assert.Equal(t, tx, expectedTx)
	})

	t.Run("if wallet is not connected, return error", func(t *testing.T) {
		w, _ := New(testPrivHex)

		// Act
		_, err := w.ContractTransactNoWait(contract, contractAddress, data)

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
		_, err := w.ContractTransactNoWait(contract, contractAddress, data)

		// Assert
		assert.Error(t, err)
	})

	t.Run("if failed to transact with contract, return error", func(t *testing.T) {
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
			reflect.TypeOf(w.provider.Eth()),
			"ContractTransact",
			func(
				_ *ether.Ether,
				auth *bind.TransactOpts,
				contract types.ContractInstance,
				contractAddress string,
				data []byte,
			) (*gethTypes.Transaction, error) {
				return nil, errors.New("error")
			},
		)

		// Act
		_, err := w.ContractTransactNoWait(contract, contractAddress, data)

		// Assert
		assert.Error(t, err)
	})
}

func TestWallet_ContractCall(t *testing.T) {
	contract := artifacts.NewPotetoStorage()
	contractAddress := "0x1234567890123456789012345678901234567890"
	data := []byte("test data")
	unpack := func(b []byte) (any, error) { return nil, nil }
	callOpts := &bind.CallOpts{}

	t.Run("can call contract", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		w := createConnectedWallet()
		expectedRes := "result"

		// Mock
		patches.ApplyMethod(
			reflect.TypeOf(w.provider.Eth()),
			"ContractCall",
			func(
				_ *ether.Ether,
				_ types.ContractInstance,
				_ common.Address,
				_ *bind.CallOpts,
				_ []byte,
				_ func([]byte) (any, error),
			) (any, error) {
				return expectedRes, nil
			},
		)

		// Act
		res, err := w.ContractCall(contract, contractAddress, callOpts, data, unpack)

		// Assert
		assert.Nil(t, err)
		assert.Equal(t, expectedRes, res)
	})

	t.Run("if wallet is not connected, return error", func(t *testing.T) {
		w, _ := New(testPrivHex)

		// Act
		_, err := w.ContractCall(contract, contractAddress, callOpts, data, unpack)

		// Assert
		assert.ErrorIs(t, err, constant.ErrWalletIsNotConnected)
	})

	t.Run("if failed to call contract, return error", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		w := createConnectedWallet()

		// Mock
		patches.ApplyMethod(
			reflect.TypeOf(w.provider.Eth()),
			"ContractCall",
			func(
				_ *ether.Ether,
				_ types.ContractInstance,
				_ common.Address,
				_ *bind.CallOpts,
				_ []byte,
				_ func([]byte) (any, error),
			) (any, error) {
				return nil, errors.New("error")
			},
		)

		// Act
		_, err := w.ContractCall(contract, contractAddress, callOpts, data, unpack)

		// Assert
		assert.Error(t, err)
	})
}

func TestWallet_ResetPool(t *testing.T) {
	txData := &gethTypes.AccessListTx{
		To:       &common.Address{},
		ChainID:  big.NewInt(1),
		Nonce:    0,
		GasPrice: big.NewInt(1),
		Gas:      0,
		Data:     []byte("data"),
	}
	expectedTx := gethTypes.NewTx(txData)

	t.Run("can reset cached values", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		w := createConnectedWallet()
		contract := artifacts.NewPotetoStorage()
		contractAddress := "0x1234567890123456789012345678901234567890"
		data := []byte("test data")

		callCount := 0

		// Mock
		patches.ApplyMethod(
			reflect.TypeOf(w.provider.Eth()),
			"ChainID",
			func(_ *ether.Ether) (*big.Int, error) {
				callCount++
				return big.NewInt(1), nil
			},
		)
		patches.ApplyMethod(
			reflect.TypeOf(w.provider.Eth()),
			"ContractTransact",
			func(
				_ *ether.Ether,
				auth *bind.TransactOpts,
				contract types.ContractInstance,
				contractAddress string,
				data []byte,
			) (*gethTypes.Transaction, error) {
				return expectedTx, nil
			},
		)

		// Act - First call should cache
		_, err := w.ContractTransactNoWait(contract, contractAddress, data)
		assert.NoError(t, err)
		assert.Equal(t, 1, callCount, "ChainID should be called once")

		// Act - Second call should use cache
		_, err = w.ContractTransactNoWait(contract, contractAddress, data)
		assert.NoError(t, err)
		assert.Equal(t, 1, callCount, "ChainID should still be called only once (cached)")

		// Act - Reset pool
		w.ResetPool()

		// Act - Third call should fetch ChainID again
		_, err = w.ContractTransactNoWait(contract, contractAddress, data)
		assert.NoError(t, err)
		assert.Equal(t, 2, callCount, "ChainID should be called again after reset")
	})

	t.Run("clears both chainID and auth", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		w := createConnectedWallet()
		contract := artifacts.NewPotetoStorage()
		contractAddress := "0x1234567890123456789012345678901234567890"
		data := []byte("test data")

		// Mock
		patches.ApplyMethod(
			reflect.TypeOf(w.provider.Eth()),
			"ChainID",
			func(_ *ether.Ether) (*big.Int, error) {
				return big.NewInt(1), nil
			},
		)
		patches.ApplyMethod(
			reflect.TypeOf(w.provider.Eth()),
			"ContractTransact",
			func(
				_ *ether.Ether,
				auth *bind.TransactOpts,
				contract types.ContractInstance,
				contractAddress string,
				data []byte,
			) (*gethTypes.Transaction, error) {
				return expectedTx, nil
			},
		)

		// Act - Create cache
		_, _ = w.ContractTransactNoWait(contract, contractAddress, data)
		assert.NotNil(t, w.cachedChainID, "ChainID should be cached")
		assert.NotNil(t, w.cachedAuth, "Auth should be cached")

		// Act - Reset
		w.ResetPool()

		// Assert
		assert.Nil(t, w.cachedChainID, "ChainID should be cleared")
		assert.Nil(t, w.cachedAuth, "Auth should be cleared")
	})
}
