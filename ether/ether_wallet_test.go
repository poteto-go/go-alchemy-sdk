package ether_test

import (
	"context"
	"errors"
	"math/big"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/poteto-go/go-alchemy-sdk/_fixture/artifacts"
	eth "github.com/poteto-go/go-alchemy-sdk/ether"
	"github.com/stretchr/testify/assert"
)

func TestEther_PendingNonceAt(t *testing.T) {
	t.Run("normal case", func(t *testing.T) {
		t.Run("success request", func(t *testing.T) {
			// Arrange
			ether := newEtherApiForTest()
			alchemyMock := newAlchemyMockOnEtherTest(t)
			defer alchemyMock.DeactivateAndReset()

			// Mock
			alchemyMock.RegisterResponder("eth_getTransactionCount", `{"jsonrpc":"2.0","id":1,"result":"0x10"}`)

			// Act
			result, err := ether.PendingNonceAt("0xa7d9ddbe1f17865597fbd27ec712455208b6b76d")

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, uint64(0x10), result)
		})
	})

	t.Run("error case", func(t *testing.T) {
		t.Run("if cannot create ethClient, return err", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Arrange
			ether := newEtherApiForTest()
			address := "0x123"

			// Mock
			patches.ApplyMethod(
				reflect.TypeOf(ether),
				"GetEthClient",
				func(_ *eth.Ether) (*ethclient.Client, error) {
					return nil, errors.New("error")
				},
			)

			// Act
			_, err := ether.PendingNonceAt(address)

			// Assert
			assert.Error(t, err)
		})

		t.Run("if failed get pending nonce, return error", func(t *testing.T) {
			// Arrange
			ether := newEtherApiForTest()
			address := "0x123"

			// Act
			_, err := ether.PendingNonceAt(address)

			// Assert
			assert.Error(t, err)
		})
	})
}

func Test_DeployContract(t *testing.T) {
	metaData := &artifacts.PotetoStorageMetaData

	t.Run("can deploy contract", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Arrange
		ether := newEtherApiForTest()
		expectedAddr := common.HexToAddress("0x123")
		expectedTx := gethTypes.NewTx(&gethTypes.LegacyTx{})

		// Mock
		patches.ApplyFunc(
			bind.LinkAndDeploy,
			func(params *bind.DeploymentParams, deploy bind.DeployFn) (*bind.DeploymentResult, error) {
				return &bind.DeploymentResult{
					Txs: map[string]*gethTypes.Transaction{
						metaData.ID: expectedTx,
					},
				}, nil
			},
		)
		patches.ApplyFunc(
			bind.WaitDeployed,
			func(ctx context.Context, b bind.DeployBackend, hash common.Hash) (common.Address, error) {
				return expectedAddr, nil
			},
		)

		// Act
		addr, err := ether.DeployContract(nil, metaData)

		// Assert
		assert.Nil(t, err)
		assert.Equal(t, expectedAddr, addr)
	})

	t.Run("error case", func(t *testing.T) {
		t.Run("if cannot create ethClient, return err", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Arrange
			ether := newEtherApiForTest()

			// Mock
			patches.ApplyMethod(
				reflect.TypeOf(ether),
				"GetEthClient",
				func(_ *eth.Ether) (*ethclient.Client, error) {
					return nil, errors.New("error")
				},
			)

			// Act
			_, err := ether.DeployContract(nil, metaData)

			// Assert
			assert.Error(t, err)
		})

		t.Run("if fail to deploy, return error", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Arrange
			ether := newEtherApiForTest()

			// Mock
			patches.ApplyFunc(
				bind.LinkAndDeploy,
				func(params *bind.DeploymentParams, deploy bind.DeployFn) (*bind.DeploymentResult, error) {
					return nil, errors.New("error")
				},
			)

			// Act
			_, err := ether.DeployContract(nil, metaData)

			// Assert
			assert.Error(t, err)
		})

		t.Run("if failed to wait for deployed, return error", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Arrange
			ether := newEtherApiForTest()
			expectedTx := gethTypes.NewTx(&gethTypes.LegacyTx{})

			// Mock
			patches.ApplyFunc(
				bind.LinkAndDeploy,
				func(params *bind.DeploymentParams, deploy bind.DeployFn) (*bind.DeploymentResult, error) {
					return &bind.DeploymentResult{
						Txs: map[string]*gethTypes.Transaction{
							metaData.ID: expectedTx,
						},
					}, nil
				},
			)
			patches.ApplyFunc(
				bind.WaitDeployed,
				func(ctx context.Context, b bind.DeployBackend, hash common.Hash) (common.Address, error) {
					return common.Address{}, errors.New("error")
				},
			)

			// Act
			_, err := ether.DeployContract(nil, metaData)

			// Assert
			assert.Error(t, err)
		})
	})
}

func TestEther_ContractTransact(t *testing.T) {
	txData := &gethTypes.AccessListTx{
		To:       &common.Address{},
		ChainID:  big.NewInt(1),
		Nonce:    0,
		GasPrice: big.NewInt(1),
		Gas:      0,
		Data:     []byte("data"),
	}

	t.Run("normal case", func(t *testing.T) {
		t.Run("success transact by contract", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Arrange
			ether := newEtherApiForTest()
			contract := artifacts.NewPotetoStorage()

			signedTx := gethTypes.NewTx(txData)

			// Mock
			patches.ApplyFunc(
				bind.Transact,
				func(c *bind.BoundContract, auth *bind.TransactOpts, data []byte) (*gethTypes.Transaction, error) {
					return signedTx, nil
				},
			)
			patches.ApplyFunc(
				bind.WaitMined,
				func(ctx context.Context, b bind.DeployBackend, hash common.Hash) (*gethTypes.Receipt, error) {
					return &gethTypes.Receipt{}, nil
				},
			)

			// Act
			_, err := ether.ContractTransact(nil, contract, "", []byte(""))

			// Assert
			assert.NoError(t, err)
		})
	})

	t.Run("error case", func(t *testing.T) {
		t.Run("if contract is nil, return error", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Arrange
			ether := newEtherApiForTest()

			// Act
			_, err := ether.ContractTransact(nil, nil, "", []byte(""))

			// Assert
			assert.Error(t, err)
		})

		t.Run("if failed to create connection, retunr error", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Arrange
			ether := newEtherApiForTest()
			contract := artifacts.NewPotetoStorage()

			// Mock
			patches.ApplyMethod(
				reflect.TypeOf(ether),
				"GetEthClient",
				func(_ *eth.Ether) (*ethclient.Client, error) {
					return nil, errors.New("error")
				},
			)

			// Act
			_, err := ether.ContractTransact(nil, contract, "", []byte(""))

			// Assert
			assert.Error(t, err)
		})

		t.Run("if failed to transact, return error", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Arrange
			ether := newEtherApiForTest()
			contract := artifacts.NewPotetoStorage()

			// Mock
			patches.ApplyFunc(
				bind.Transact,
				func(c *bind.BoundContract, auth *bind.TransactOpts, data []byte) (*gethTypes.Transaction, error) {
					return nil, errors.New("error")
				},
			)

			// Act
			_, err := ether.ContractTransact(nil, contract, "", []byte(""))

			// Assert
			assert.Error(t, err)
		})

		t.Run("if failed to wait for mained, return error", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Arrange
			ether := newEtherApiForTest()
			contract := artifacts.NewPotetoStorage()
			signedTx := gethTypes.NewTx(txData)

			// Mock
			patches.ApplyFunc(
				bind.Transact,
				func(c *bind.BoundContract, auth *bind.TransactOpts, data []byte) (*gethTypes.Transaction, error) {
					return signedTx, nil
				},
			)
			patches.ApplyFunc(
				bind.WaitMined,
				func(ctx context.Context, b bind.DeployBackend, hash common.Hash) (*gethTypes.Receipt, error) {
					return nil, errors.New("error")
				},
			)

			// Act
			_, err := ether.ContractTransact(nil, contract, "", []byte(""))

			// Assert
			assert.Error(t, err)
		})
	})
}
