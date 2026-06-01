package e2e

import (
	"context"
	"math/big"
	"os"
	"strconv"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/poteto-go/go-alchemy-sdk/_fixture/artifacts"
	"github.com/poteto-go/go-alchemy-sdk/deployer"
	"github.com/poteto-go/go-alchemy-sdk/gas"
	"github.com/poteto-go/go-alchemy-sdk/types"
	"github.com/poteto-go/go-alchemy-sdk/wallet"
	"github.com/stretchr/testify/assert"
)

// These are Kurtosis ethereum-package genesis accounts (public test keys, not credentials)
var initPrivateKey = "bcdf20249abf0ed6d944c0288fad489e33f66b3960d9e6229c1cd214ed3bbe31"
var initAddress = "0x8943545177806ED17B9F23F0a21ee5948eCaa776"
var otherPrivateKey = "39725efee3fb28614de3bacaffe4cc4bd8c436257e2c8bb887c4b5c4be45e76d"
var otherAddress = "0xE25583099BA105D9ec0A67f5Ae86D90e50036425"
var deployedContractAddress common.Address
var alchemy gas.Alchemy

func TestMain(m *testing.M) {
	setup()

	m.Run()

	teardown()
}

func setup() {
	port, err := strconv.Atoi(os.Getenv("RPC_PORT"))
	if err != nil {
		panic(err)
	}

	setting := gas.AlchemySetting{
		PrivateNetworkConfig: gas.PrivateNetworkConfig{
			Host: "127.0.0.1",
			Port: port,
		},
	}
	var errA error
	alchemy, errA = gas.NewAlchemy(setting)
	if errA != nil {
		panic(errA)
	}
}

func teardown() {
}

func TestSenario_BaseMethod(t *testing.T) {
	t.Run("GetGasPrice", func(t *testing.T) {
		gasPrice, err := alchemy.Core.GetGasPrice()

		assert.Nil(t, err)
		assert.Equal(t, gasPrice.Cmp(big.NewInt(0)), 1)
	})

	t.Run("estimate gas", func(t *testing.T) {
		gas, err := alchemy.Core.EstimateGas(types.TransactionRequest{
			From:  initAddress,
			To:    "0x0",
			Value: "0x0",
		})

		assert.Nil(t, err)
		assert.Equal(t, gas.Cmp(big.NewInt(0)), 1)
	})

	t.Run("PeerCount", func(t *testing.T) {
		peerCount, err := alchemy.Core.PeerCount()

		assert.Nil(t, err)
		assert.GreaterOrEqual(t, peerCount, uint64(0))
	})
}

func TestScenario_GetBalance(t *testing.T) {
	t.Run("core namespace case", func(t *testing.T) {
		balance, err := alchemy.Core.GetBalance(initAddress, "latest")

		assert.Nil(t, err)
		assert.Equal(t, balance.Cmp(big.NewInt(0)), 1)
	})

	t.Run("1. can create wallet 2. connect wallet 3. can get balance", func(t *testing.T) {
		w, err := wallet.New(initPrivateKey)

		assert.Nil(t, err)

		w.Connect(alchemy.GetProvider())

		balance, err := w.GetBalance()

		assert.Nil(t, err)
		assert.Equal(t, balance.Cmp(big.NewInt(0)), 1)
	})
}

func TestSenario_DeployContract(t *testing.T) {
	t.Run("1. can create wallet 2. connect wallet 3. can deploy contract", func(t *testing.T) {
		w, err := wallet.New(initPrivateKey)

		assert.Nil(t, err)

		w.Connect(alchemy.GetProvider())

		contractAddress, err := w.DeployContract(context.Background(), &artifacts.PotetoStorageMetaData)

		assert.Nil(t, err)
		assert.NotEqual(t, contractAddress, common.HexToAddress("0x0"))
		deployedContractAddress = contractAddress

		t.Run("IsContractAddress is true", func(t *testing.T) {
			isContractAddress := alchemy.Core.IsContractAddress(deployedContractAddress.Hex())

			assert.True(t, isContractAddress)
		})

		t.Run("can get Code", func(t *testing.T) {
			code, err := alchemy.Core.GetCode(deployedContractAddress.Hex(), types.BlockTagOrHash{
				BlockTag: "latest",
			})

			assert.Nil(t, err)
			assert.Equal(t, code[:10], artifacts.PotetoStorageMetaData.Bin[:10])

			block, err := alchemy.Core.GetBlock(types.BlockTagOrHash{
				BlockTag: "latest",
			})

			assert.Nil(t, err)

			blockHash := block.Hash()
			code, err = alchemy.Core.GetCode(deployedContractAddress.Hex(), types.BlockTagOrHash{
				BlockHash: blockHash.Hex(),
			})

			assert.Nil(t, err)
			assert.Equal(t, code[:10], artifacts.PotetoStorageMetaData.Bin[:10])

			t.Run("can get transaction & its receipt form deployed block", func(t *testing.T) {
				txHash := block.Body().Transactions[0].Hash()
				tx, isPending, err := alchemy.Core.GetTransaction(txHash.Hex())

				assert.Nil(t, err)
				assert.False(t, isPending)
				assert.Equal(t, tx.Hash().Hex(), txHash.Hex())

				txReceipt, err := alchemy.Core.GetTransactionReceipt(txHash.Hex())

				assert.Nil(t, err)
				assert.Equal(t, txReceipt.TxHash.Hex(), txHash.Hex())
			})
		})

		t.Run("blockNumber > 0", func(t *testing.T) {
			blockNumber, err := alchemy.Core.GetBlockNumber()

			assert.Nil(t, err)
			assert.Greater(t, blockNumber, uint64(0))
		})
	})
}

func TestScenario_StableCoin(t *testing.T) {
	t.Run("1. can create wallet 2. connect wallet 3. can deploy erc20 contract as stablecoin", func(t *testing.T) {
		w, err := wallet.New(initPrivateKey)

		assert.Nil(t, err)

		w.Connect(alchemy.GetProvider())

		erc20Metadata := &artifacts.ERC20MetaData
		deployer.BindDeploymentMetadata(erc20Metadata, big.NewInt(1000))
		contractAddress, err := w.DeployContract(context.Background(), erc20Metadata)

		assert.Nil(t, err)
		assert.NotEqual(t, contractAddress, common.HexToAddress("0x0"))
		contractHex := contractAddress.Hex()

		t.Run("can get balance via StableCoin", func(t *testing.T) {
			balance, err := w.StableCoin().BalanceOf(contractHex)

			assert.Nil(t, err)
			assert.Equal(t, balance.Cmp(big.NewInt(10)), 1)
		})

		t.Run("can get other info via StableCoin", func(t *testing.T) {
			totalSupply, err := w.StableCoin().TotalSupply(contractHex)
			assert.Nil(t, err)
			assert.Equal(t, totalSupply.Cmp(big.NewInt(1000)), 0)

			name, err := w.StableCoin().Name(contractHex)
			assert.Nil(t, err)
			assert.Equal(t, "Minimal Token", name)

			symbol, err := w.StableCoin().Symbol(contractHex)
			assert.Nil(t, err)
			assert.Equal(t, "MTK", symbol)

			decimals, err := w.StableCoin().Decimals(contractHex)
			assert.Nil(t, err)
			assert.Equal(t, uint8(18), decimals)

			allowance, err := w.StableCoin().Allowance(contractHex, initAddress, otherAddress)
			assert.Nil(t, err)
			assert.Equal(t, allowance.Cmp(big.NewInt(0)), 0)
		})
	})
}

func TestScenario_ERC20(t *testing.T) {
	t.Run("1. can create wallet 2. connect wallet 3. can deploy erc20 contract", func(t *testing.T) {
		w, err := wallet.New(initPrivateKey)

		assert.Nil(t, err)

		w.Connect(alchemy.GetProvider())

		erc20Metadata := &artifacts.ERC20MetaData
		deployer.BindDeploymentMetadata(erc20Metadata, big.NewInt(1000))
		contractAddress, err := w.DeployContract(context.Background(), erc20Metadata)

		assert.Nil(t, err)
		assert.NotEqual(t, contractAddress, common.HexToAddress("0x0"))
		deployedContractAddress = contractAddress
		contractHex := contractAddress.Hex()

		t.Run("can get balance", func(t *testing.T) {
			balance, err := w.ERC20().BalanceOf(contractHex)

			assert.Nil(t, err)
			assert.Equal(t, balance.Cmp(big.NewInt(10)), 1)
		})

		t.Run("can get other info", func(t *testing.T) {
			totalSupply, err := w.ERC20().TotalSupply(contractHex)
			assert.Nil(t, err)
			assert.Equal(t, totalSupply.Cmp(big.NewInt(1000)), 0)

			name, err := w.ERC20().Name(contractHex)
			assert.Nil(t, err)
			assert.Equal(t, "Minimal Token", name)

			symbol, err := w.ERC20().Symbol(contractHex)
			assert.Nil(t, err)
			assert.Equal(t, "MTK", symbol)

			decimals, err := w.ERC20().Decimals(contractHex)
			assert.Nil(t, err)
			assert.Equal(t, uint8(18), decimals)

			allowance, err := w.ERC20().Allowance(contractHex, initAddress, otherAddress)
			assert.Nil(t, err)
			assert.Equal(t, allowance.Cmp(big.NewInt(0)), 0)
		})

		t.Run("Approve: allowance is set to approved amount", func(t *testing.T) {
			approveAmount := big.NewInt(100)

			_, err := w.ERC20().Approve(context.Background(), contractHex, otherAddress, approveAmount, nil)
			assert.Nil(t, err)

			allowance, err := w.ERC20().Allowance(contractHex, initAddress, otherAddress)
			assert.Nil(t, err)
			assert.Equal(t, 0, allowance.Cmp(approveAmount))
		})

		t.Run("ApproveNoWait: returns txHash, allowance is set after mined", func(t *testing.T) {
			approveAmount := big.NewInt(200)

			txHash, err := w.ERC20().ApproveNoWait(contractHex, otherAddress, approveAmount, nil)
			assert.Nil(t, err)
			assert.NotEqual(t, txHash.Hex(), "0x0000000000000000000000000000000000000000000000000000000000000000")

			_, err = alchemy.Transact.WaitMined(context.Background(), txHash.Hex())
			assert.Nil(t, err)

			allowance, err := w.ERC20().Allowance(contractHex, initAddress, otherAddress)
			assert.Nil(t, err)
			assert.Equal(t, 0, allowance.Cmp(approveAmount))
		})

		t.Run("TransferFrom: otherWallet transfers from initAddress after approval", func(t *testing.T) {
			transferAmount := big.NewInt(50)

			_, err := w.ERC20().Approve(context.Background(), contractHex, otherAddress, transferAmount, nil)
			assert.Nil(t, err)

			balanceBefore, err := w.ERC20().BalanceOf(contractHex)
			assert.Nil(t, err)

			otherWallet, err := wallet.New(otherPrivateKey)
			assert.Nil(t, err)
			otherWallet.Connect(alchemy.GetProvider())

			_, err = otherWallet.ERC20().TransferFrom(context.Background(), contractHex, initAddress, otherAddress, transferAmount, nil)
			assert.Nil(t, err)

			balanceAfter, err := w.ERC20().BalanceOf(contractHex)
			assert.Nil(t, err)
			assert.Equal(t, 0, new(big.Int).Sub(balanceBefore, transferAmount).Cmp(balanceAfter))
		})

		t.Run("TransferFromNoWait: otherWallet transfers from initAddress, balance decreases after mined", func(t *testing.T) {
			transferAmount := big.NewInt(30)

			_, err := w.ERC20().Approve(context.Background(), contractHex, otherAddress, transferAmount, nil)
			assert.Nil(t, err)

			balanceBefore, err := w.ERC20().BalanceOf(contractHex)
			assert.Nil(t, err)

			otherWallet, err := wallet.New(otherPrivateKey)
			assert.Nil(t, err)
			otherWallet.Connect(alchemy.GetProvider())

			txHash, err := otherWallet.ERC20().TransferFromNoWait(contractHex, initAddress, otherAddress, transferAmount, nil)
			assert.Nil(t, err)
			assert.NotEqual(t, txHash.Hex(), "0x0000000000000000000000000000000000000000000000000000000000000000")

			_, err = alchemy.Transact.WaitMined(context.Background(), txHash.Hex())
			assert.Nil(t, err)

			balanceAfter, err := w.ERC20().BalanceOf(contractHex)
			assert.Nil(t, err)
			assert.Equal(t, 0, new(big.Int).Sub(balanceBefore, transferAmount).Cmp(balanceAfter))
		})
	})
}

func TestScenario_FiatToken(t *testing.T) {
	t.Run("1. can create wallet 2. connect wallet 3. can deploy fiat token contract", func(t *testing.T) {
		w, err := wallet.New(initPrivateKey)

		assert.Nil(t, err)

		w.Connect(alchemy.GetProvider())

		ownerAddr := common.HexToAddress(initAddress)
		fiatTokenMetadata := &artifacts.FiatTokenMetaData
		err = deployer.BindDeploymentMetadata(fiatTokenMetadata,
			"USD Coin",
			"USDC",
			"USD",
			uint8(6),
			ownerAddr,
			ownerAddr,
			ownerAddr,
			ownerAddr,
		)
		assert.Nil(t, err)

		contractAddress, err := w.DeployContract(context.Background(), fiatTokenMetadata)

		assert.Nil(t, err)
		assert.NotEqual(t, contractAddress, common.HexToAddress("0x0"))

		t.Run("IsContractAddress is true", func(t *testing.T) {
			isContractAddress := alchemy.Core.IsContractAddress(contractAddress.Hex())

			assert.True(t, isContractAddress)
		})
	})
}

func TestScenario_SendTransaction(t *testing.T) {
	w, err := wallet.New(initPrivateKey)

	assert.Nil(t, err)
	w.Connect(alchemy.GetProvider())

	t.Run("can get pending nonce", func(t *testing.T) {
		pendingNonce, err := w.PendingNonceAt()

		assert.Nil(t, err)
		assert.NotEqual(t, pendingNonce, uint64(0)) // first transaction
	})

	t.Run("can send transaciton", func(t *testing.T) {
		txRequest := types.TransactionRequest{
			From:     initAddress,
			To:       otherAddress,
			Value:    "0x123",
			GasLimit: 300000,
		}

		txHash, err := w.SendTransaction(txRequest)

		assert.Nil(t, err)
		assert.NotEqual(t, txHash.Hex(), "0x0000000000000000000000000000000000000000000000000000000000000000")

		// Verify transaction is pending
		tx, isPending, err := alchemy.Core.GetTransaction(txHash.Hex())
		assert.Nil(t, err)
		assert.True(t, isPending, "transaction should be pending immediately after sending")
		assert.Equal(t, txHash, tx.Hash())

		// wait for transact finish
		txReceipt, err := alchemy.Transact.WaitMined(context.Background(), txHash.Hex())
		assert.Nil(t, err)
		assert.Equal(t, txReceipt.TxHash.Hex(), txHash.Hex())

		// Verify transaction is not pending
		tx, isPending, err = alchemy.Core.GetTransaction(txHash.Hex())
		assert.Nil(t, err)
		assert.False(t, isPending, "transaction should be finished after waitMined")
		assert.Equal(t, txHash, tx.Hash())
	})
}
