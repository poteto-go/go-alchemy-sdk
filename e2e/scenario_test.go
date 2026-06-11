package e2e

import (
	"context"
	"math/big"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/poteto-go/go-alchemy-sdk/_fixture/artifacts"
	"github.com/poteto-go/go-alchemy-sdk/batch"
	"github.com/poteto-go/go-alchemy-sdk/constant"
	"github.com/poteto-go/go-alchemy-sdk/deployer"
	"github.com/poteto-go/go-alchemy-sdk/famous"
	"github.com/poteto-go/go-alchemy-sdk/gas"
	"github.com/poteto-go/go-alchemy-sdk/typeddata"
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

	t.Run("Batch fetches blockNumber, gasPrice & chainId in a single round-trip", func(t *testing.T) {
		b := batch.NewBatcher(alchemy.GetProvider().Eth())
		blockNumber := b.Core.BlockNumber()
		gasPrice := b.Core.GasPrice()
		chainId := b.Core.ChainID()

		err := b.Send()
		assert.Nil(t, err)

		bn, err := blockNumber.Unwrap()
		assert.Nil(t, err)
		assert.GreaterOrEqual(t, bn, uint64(0))

		gp, err := gasPrice.Unwrap()
		assert.Nil(t, err)
		assert.Equal(t, 1, gp.Cmp(big.NewInt(0)))

		id, err := chainId.Unwrap()
		assert.Nil(t, err)
		assert.Equal(t, 1, id.Cmp(big.NewInt(0)))
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

// countingTransport wraps an http.RoundTripper and counts how many RPC calls
// pass through it, modeling the metrics/benchmarking use case for a custom
// transport against a blockchain node.
type countingTransport struct {
	base  http.RoundTripper
	count atomic.Int64
}

func (t *countingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.count.Add(1)
	return t.base.RoundTrip(req)
}

func TestScenario_CustomTransport(t *testing.T) {
	t.Run("RPC calls flow through the caller-supplied transport", func(t *testing.T) {
		port, err := strconv.Atoi(os.Getenv("RPC_PORT"))
		assert.Nil(t, err)

		transport := &countingTransport{base: http.DefaultTransport}
		customAlchemy, err := gas.NewAlchemy(gas.AlchemySetting{
			PrivateNetworkConfig: gas.PrivateNetworkConfig{
				Host: "127.0.0.1",
				Port: port,
			},
			Transport: transport,
		})
		assert.Nil(t, err)

		blockNumber, err := customAlchemy.Core.GetBlockNumber()

		assert.Nil(t, err)
		assert.Greater(t, blockNumber, uint64(0))
		assert.Greater(t, transport.count.Load(), int64(0), "the custom transport must handle the RPC round-trip")
	})
}

func TestScenario_Famous_StableCoin(t *testing.T) {
	t.Run("SupportedNetworks returns non-empty list including EthMainnet", func(t *testing.T) {
		networks := famous.SupportedNetworks()

		assert.NotEmpty(t, networks)
		assert.True(t, slices.Contains(networks, types.EthMainnet))
	})

	t.Run("SupportedSymbols returns USDC/USDT/JPYC for EthMainnet", func(t *testing.T) {
		symbols := famous.SupportedSymbols(types.EthMainnet)

		assert.True(t, slices.Contains(symbols, famous.USDC))
		assert.True(t, slices.Contains(symbols, famous.USDT))
		assert.True(t, slices.Contains(symbols, famous.JPYC))
	})

	t.Run("SupportedSymbols returns empty for unsupported network", func(t *testing.T) {
		symbols := famous.SupportedSymbols(types.SolanaMainnet)

		assert.Empty(t, symbols)
	})

	t.Run("ContractAddress resolves typed symbol on EthMainnet", func(t *testing.T) {
		addr, err := famous.ContractAddress(types.EthMainnet, famous.USDC)

		assert.NoError(t, err)
		assert.Equal(t, common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"), addr)
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

func TestScenario_StableCoin_FiatToken(t *testing.T) {
	t.Run("1. can create wallet 2. connect wallet 3. can deploy fiat token contract as stablecoin", func(t *testing.T) {
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
		contractHex := contractAddress.Hex()

		t.Run("IsContractAddress is true", func(t *testing.T) {
			isContractAddress := alchemy.Core.IsContractAddress(contractHex)

			assert.True(t, isContractAddress)
		})

		t.Run("can read stablecoin metadata via StableCoin", func(t *testing.T) {
			name, err := w.StableCoin().Name(contractHex)
			assert.Nil(t, err)
			assert.Equal(t, "USD Coin", name)

			symbol, err := w.StableCoin().Symbol(contractHex)
			assert.Nil(t, err)
			assert.Equal(t, "USDC", symbol)

			decimals, err := w.StableCoin().Decimals(contractHex)
			assert.Nil(t, err)
			assert.Equal(t, uint8(6), decimals)

			totalSupply, err := w.StableCoin().TotalSupply(contractHex)
			assert.Nil(t, err)
			assert.Equal(t, 0, totalSupply.Cmp(big.NewInt(0)))

			currency, err := w.StableCoin().Currency(contractHex)
			assert.Nil(t, err)
			assert.Equal(t, "USD", currency)

			version, err := w.StableCoin().Version(contractHex)
			assert.Nil(t, err)
			assert.Equal(t, "1", version)
		})

		t.Run("can configure minter and mint tokens", func(t *testing.T) {
			fiatToken := artifacts.NewFiatToken()

			// configure minter allowance (owner == initAddress in this test setup)
			configureMinterData := fiatToken.PackConfigureMinter(
				common.HexToAddress(initAddress),
				big.NewInt(1_000_000),
			)
			txHash, err := w.SendTransaction(types.TransactionRequest{
				From:     initAddress,
				To:       contractHex,
				Value:    "0x0",
				GasLimit: 300000,
				Data:     configureMinterData,
			})
			assert.Nil(t, err)
			_, err = alchemy.Transact.WaitMined(context.Background(), txHash.Hex())
			assert.Nil(t, err)

			mintAmount := big.NewInt(500_000)
			receipt, err := w.StableCoin().Mint(
				context.Background(),
				contractHex,
				initAddress,
				mintAmount,
				nil,
			)
			assert.Nil(t, err)
			assert.NotNil(t, receipt)

			totalSupply, err := w.StableCoin().TotalSupply(contractHex)
			assert.Nil(t, err)
			assert.Equal(t, 0, totalSupply.Cmp(mintAmount))
		})

		t.Run("can burn tokens", func(t *testing.T) {
			burnAmount := big.NewInt(100_000)

			supplyBefore, err := w.StableCoin().TotalSupply(contractHex)
			assert.Nil(t, err)

			receipt, err := w.StableCoin().Burn(
				context.Background(),
				contractHex,
				burnAmount,
				nil,
			)
			assert.Nil(t, err)
			assert.NotNil(t, receipt)

			supplyAfter, err := w.StableCoin().TotalSupply(contractHex)
			assert.Nil(t, err)
			assert.Equal(t, 0, new(big.Int).Sub(supplyBefore, burnAmount).Cmp(supplyAfter))
		})

		t.Run("can blacklist and unBlacklist an address", func(t *testing.T) {
			// address should not be blacklisted initially
			isBlacklisted, err := w.StableCoin().IsBlacklisted(contractHex, otherAddress)
			assert.Nil(t, err)
			assert.False(t, isBlacklisted)

			// blacklist the address (initAddress is the blacklister in this test setup)
			receipt, err := w.StableCoin().Blacklist(context.Background(), contractHex, otherAddress, nil)
			assert.Nil(t, err)
			assert.NotNil(t, receipt)

			// address should now be blacklisted
			isBlacklisted, err = w.StableCoin().IsBlacklisted(contractHex, otherAddress)
			assert.Nil(t, err)
			assert.True(t, isBlacklisted)

			// unBlacklist the address
			receipt, err = w.StableCoin().UnBlacklist(context.Background(), contractHex, otherAddress, nil)
			assert.Nil(t, err)
			assert.NotNil(t, receipt)

			// address should no longer be blacklisted
			isBlacklisted, err = w.StableCoin().IsBlacklisted(contractHex, otherAddress)
			assert.Nil(t, err)
			assert.False(t, isBlacklisted)
		})

		t.Run("pause prevents transfer and unpause restores it", func(t *testing.T) {
			// contract should not be paused initially
			paused, err := w.StableCoin().Paused(contractHex)
			assert.Nil(t, err)
			assert.False(t, paused)

			// pause the contract (initAddress is the pauser in this test setup)
			receipt, err := w.StableCoin().Pause(context.Background(), contractHex, nil)
			assert.Nil(t, err)
			assert.NotNil(t, receipt)

			// contract should now be paused
			paused, err = w.StableCoin().Paused(contractHex)
			assert.Nil(t, err)
			assert.True(t, paused)

			// transfer should fail while paused
			_, err = w.StableCoin().Transfer(context.Background(), contractHex, otherAddress, big.NewInt(100), nil)
			assert.Error(t, err, "transfer should fail when contract is paused")

			// unpause the contract
			receipt, err = w.StableCoin().Unpause(context.Background(), contractHex, nil)
			assert.Nil(t, err)
			assert.NotNil(t, receipt)

			// contract should no longer be paused
			paused, err = w.StableCoin().Paused(contractHex)
			assert.Nil(t, err)
			assert.False(t, paused)

			// transfer should succeed after unpause
			balanceBefore, err := w.StableCoin().BalanceOf(contractHex)
			assert.Nil(t, err)

			transferAmount := big.NewInt(100)
			_, err = w.StableCoin().Transfer(context.Background(), contractHex, otherAddress, transferAmount, nil)
			assert.Nil(t, err)

			balanceAfter, err := w.StableCoin().BalanceOf(contractHex)
			assert.Nil(t, err)
			assert.Equal(t, 0, new(big.Int).Sub(balanceBefore, transferAmount).Cmp(balanceAfter))
		})

		t.Run("can configure minter, check allowance, and remove minter", func(t *testing.T) {
			allowance := big.NewInt(5_000_000)

			// address should not be a minter initially
			isMinter, err := alchemy.StableCoin.IsMinter(contractHex, otherAddress)
			assert.Nil(t, err)
			assert.False(t, isMinter)

			// configure minter with allowance (initAddress is the masterMinter in this test setup)
			receipt, err := w.StableCoin().ConfigureMinter(context.Background(), contractHex, otherAddress, allowance, nil)
			assert.Nil(t, err)
			assert.NotNil(t, receipt)

			// address should now be a minter
			isMinter, err = alchemy.StableCoin.IsMinter(contractHex, otherAddress)
			assert.Nil(t, err)
			assert.True(t, isMinter)

			// allowance should match what was configured
			minterAllowance, err := alchemy.StableCoin.MinterAllowance(contractHex, otherAddress)
			assert.Nil(t, err)
			assert.Equal(t, 0, allowance.Cmp(minterAllowance))

			// remove the minter
			receipt, err = w.StableCoin().RemoveMinter(context.Background(), contractHex, otherAddress, nil)
			assert.Nil(t, err)
			assert.NotNil(t, receipt)

			// address should no longer be a minter
			isMinter, err = alchemy.StableCoin.IsMinter(contractHex, otherAddress)
			assert.Nil(t, err)
			assert.False(t, isMinter)
		})

		// UpdateMasterMinter/UpdateBlacklister/UpdatePauser require owner role.
		// Run these before TransferOwnership so initAddress still has the owner role.
		t.Run("can update master minter", func(t *testing.T) {
			masterMinterBefore, err := alchemy.StableCoin.MasterMinter(contractHex)
			assert.Nil(t, err)
			assert.Equal(t, common.HexToAddress(initAddress), masterMinterBefore)

			receipt, err := w.StableCoin().UpdateMasterMinter(context.Background(), contractHex, otherAddress, nil)
			assert.Nil(t, err)
			assert.NotNil(t, receipt)

			masterMinterAfter, err := alchemy.StableCoin.MasterMinter(contractHex)
			assert.Nil(t, err)
			assert.Equal(t, common.HexToAddress(otherAddress), masterMinterAfter)
		})

		t.Run("can update blacklister", func(t *testing.T) {
			blacklisterBefore, err := alchemy.StableCoin.Blacklister(contractHex)
			assert.Nil(t, err)
			assert.Equal(t, common.HexToAddress(initAddress), blacklisterBefore)

			receipt, err := w.StableCoin().UpdateBlacklister(context.Background(), contractHex, otherAddress, nil)
			assert.Nil(t, err)
			assert.NotNil(t, receipt)

			blacklisterAfter, err := alchemy.StableCoin.Blacklister(contractHex)
			assert.Nil(t, err)
			assert.Equal(t, common.HexToAddress(otherAddress), blacklisterAfter)
		})

		t.Run("can update pauser", func(t *testing.T) {
			pauserBefore, err := alchemy.StableCoin.Pauser(contractHex)
			assert.Nil(t, err)
			assert.Equal(t, common.HexToAddress(initAddress), pauserBefore)

			receipt, err := w.StableCoin().UpdatePauser(context.Background(), contractHex, otherAddress, nil)
			assert.Nil(t, err)
			assert.NotNil(t, receipt)

			pauserAfter, err := alchemy.StableCoin.Pauser(contractHex)
			assert.Nil(t, err)
			assert.Equal(t, common.HexToAddress(otherAddress), pauserAfter)
		})

		t.Run("transferOwnership transfers owner to new address", func(t *testing.T) {
			ownerBefore, err := alchemy.StableCoin.Owner(contractHex)
			assert.Nil(t, err)
			assert.Equal(t, common.HexToAddress(initAddress), ownerBefore)

			receipt, err := w.StableCoin().TransferOwnership(context.Background(), contractHex, otherAddress, nil)
			assert.Nil(t, err)
			assert.NotNil(t, receipt)

			ownerAfter, err := alchemy.StableCoin.Owner(contractHex)
			assert.Nil(t, err)
			assert.Equal(t, common.HexToAddress(otherAddress), ownerAfter)
		})

		t.Run("EIP-2612: nonces returns 0 for new owner", func(t *testing.T) {
			nonce, err := alchemy.StableCoin.Nonces(contractHex, initAddress)

			assert.Nil(t, err)
			assert.Equal(t, int64(0), nonce.Int64())
		})

		t.Run("EIP-2612: domain separator returns non-zero value", func(t *testing.T) {
			ds, err := alchemy.StableCoin.DomainSeparator(contractHex)

			assert.Nil(t, err)
			var zero [32]byte
			assert.NotEqual(t, zero, ds)
		})

		t.Run("EIP-2612: permit sets allowance and increments nonce", func(t *testing.T) {
			permitValue := big.NewInt(999)
			deadline := new(big.Int).SetInt64(9999999999)

			allowanceBefore, err := alchemy.StableCoin.Allowance(contractHex, initAddress, otherAddress)
			assert.Nil(t, err)

			receipt, err := w.StableCoin().Permit(
				context.Background(),
				contractHex,
				otherAddress,
				permitValue,
				deadline,
				nil,
			)
			assert.Nil(t, err)
			assert.NotNil(t, receipt)

			allowanceAfter, err := alchemy.StableCoin.Allowance(contractHex, initAddress, otherAddress)
			assert.Nil(t, err)
			assert.NotEqual(t, allowanceBefore.Int64(), allowanceAfter.Int64())
			assert.Equal(t, permitValue.Int64(), allowanceAfter.Int64())

			nonce, err := alchemy.StableCoin.Nonces(contractHex, initAddress)
			assert.Nil(t, err)
			assert.Equal(t, int64(1), nonce.Int64())
		})

		t.Run("EIP-3009: authorizationState returns false for unused nonce", func(t *testing.T) {
			var nonce [32]byte
			nonce[0] = 0xde
			nonce[31] = 0xad

			used, err := alchemy.StableCoin.AuthorizationState(contractHex, initAddress, nonce)

			assert.Nil(t, err)
			assert.False(t, used)
		})

		t.Run("EIP-3009: transferWithAuthorization moves tokens and marks nonce used", func(t *testing.T) {
			transferValue := big.NewInt(1)
			validAfter := big.NewInt(0)
			validBefore := new(big.Int).SetInt64(9999999999)
			var nonce [32]byte
			nonce[0] = 0xEE
			nonce[31] = 0xFF

			domainSeparator, err := alchemy.StableCoin.DomainSeparator(contractHex)
			assert.Nil(t, err)

			sig, err := typeddata.SignEIP712Str(
				initPrivateKey,
				domainSeparator,
				typeddata.EncodeWords(
					constant.TransferWithAuthorizationTypeHash,
					initAddress,
					otherAddress,
					transferValue,
					validAfter,
					validBefore,
					nonce,
				),
			)
			assert.Nil(t, err)

			balanceBefore, err := alchemy.StableCoin.BalanceOf(contractHex, otherAddress)
			assert.Nil(t, err)

			receipt, err := w.StableCoin().TransferWithAuthorization(
				context.Background(),
				contractHex,
				initAddress,
				otherAddress,
				transferValue,
				validAfter,
				validBefore,
				nonce,
				sig,
				nil,
			)
			assert.Nil(t, err)
			assert.NotNil(t, receipt)

			balanceAfter, err := alchemy.StableCoin.BalanceOf(contractHex, otherAddress)
			assert.Nil(t, err)
			assert.Equal(t, new(big.Int).Add(balanceBefore, transferValue).Int64(), balanceAfter.Int64())

			used, err := alchemy.StableCoin.AuthorizationState(contractHex, initAddress, nonce)
			assert.Nil(t, err)
			assert.True(t, used)
		})

		t.Run("EIP-3009: cancelAuthorization marks nonce as used", func(t *testing.T) {
			var nonce [32]byte
			nonce[0] = 0xCA
			nonce[31] = 0xC1

			usedBefore, err := alchemy.StableCoin.AuthorizationState(contractHex, initAddress, nonce)
			assert.Nil(t, err)
			assert.False(t, usedBefore)

			domainSeparator, err := alchemy.StableCoin.DomainSeparator(contractHex)
			assert.Nil(t, err)

			sig, err := typeddata.SignEIP712Str(
				initPrivateKey,
				domainSeparator,
				typeddata.EncodeWords(
					constant.CancelAuthorizationTypeHash,
					initAddress,
					nonce,
				),
			)
			assert.Nil(t, err)

			receipt, err := w.StableCoin().CancelAuthorization(
				context.Background(),
				contractHex,
				initAddress,
				nonce,
				sig,
				nil,
			)
			assert.Nil(t, err)
			assert.NotNil(t, receipt)

			usedAfter, err := alchemy.StableCoin.AuthorizationState(contractHex, initAddress, nonce)
			assert.Nil(t, err)
			assert.True(t, usedAfter)
		})
	})
}

func TestScenario_Nft(t *testing.T) {
	t.Run("1. can create wallet 2. connect wallet 3. deploy erc721 contract", func(t *testing.T) {
		w, err := wallet.New(initPrivateKey)

		assert.Nil(t, err)

		w.Connect(alchemy.GetProvider())

		// ERC721 has a no-arg constructor; no BindDeploymentMetadata needed.
		contractAddress, err := w.DeployContract(context.Background(), &artifacts.ERC721MetaData)

		assert.Nil(t, err)
		assert.NotEqual(t, contractAddress, common.HexToAddress("0x0"))
		contractHex := contractAddress.Hex()

		erc721 := artifacts.NewERC721()
		tokenId := big.NewInt(1)

		// Mint token 1 to initAddress.
		mintData := erc721.PackMint(common.HexToAddress(initAddress), tokenId)
		txHash, err := w.SendTransaction(types.TransactionRequest{
			From:     initAddress,
			To:       contractHex,
			Value:    "0x0",
			GasLimit: 300000,
			Data:     mintData,
		})
		assert.Nil(t, err)
		_, err = alchemy.Transact.WaitMined(context.Background(), txHash.Hex())
		assert.Nil(t, err)

		t.Run("can get name via Nft namespace", func(t *testing.T) {
			name, err := alchemy.Nft.Name(contractHex)
			assert.Nil(t, err)
			assert.Equal(t, "Minimal NFT", name)
		})

		t.Run("can get symbol via Nft namespace", func(t *testing.T) {
			symbol, err := alchemy.Nft.Symbol(contractHex)
			assert.Nil(t, err)
			assert.Equal(t, "MNFT", symbol)
		})

		t.Run("can get ownerOf minted token", func(t *testing.T) {
			owner, err := alchemy.Nft.OwnerOf(contractHex, tokenId)
			assert.Nil(t, err)
			assert.Equal(t, strings.ToLower(initAddress), owner)
		})

		t.Run("can get tokenURI for minted token", func(t *testing.T) {
			uri, err := alchemy.Nft.TokenURI(contractHex, tokenId)
			assert.Nil(t, err)
			assert.Equal(t, "https://example.com/nft/1", uri)
		})

		t.Run("getApproved returns approved address after approve", func(t *testing.T) {
			// No approval yet — should return zero address.
			approvedBefore, err := alchemy.Nft.GetApproved(contractHex, tokenId)
			assert.Nil(t, err)
			assert.Equal(t, strings.ToLower(common.HexToAddress("0x0").Hex()), approvedBefore)

			// Approve otherAddress for token 1.
			approveData := erc721.PackApprove(common.HexToAddress(otherAddress), tokenId)
			txHash, err := w.SendTransaction(types.TransactionRequest{
				From:     initAddress,
				To:       contractHex,
				Value:    "0x0",
				GasLimit: 300000,
				Data:     approveData,
			})
			assert.Nil(t, err)
			_, err = alchemy.Transact.WaitMined(context.Background(), txHash.Hex())
			assert.Nil(t, err)

			approvedAfter, err := alchemy.Nft.GetApproved(contractHex, tokenId)
			assert.Nil(t, err)
			assert.Equal(t, strings.ToLower(otherAddress), approvedAfter)
		})

		t.Run("isApprovedForAll returns false before setApprovalForAll", func(t *testing.T) {
			isApproved, err := alchemy.Nft.IsApprovedForAll(contractHex, initAddress, otherAddress)
			assert.Nil(t, err)
			assert.False(t, isApproved)
		})

		t.Run("isApprovedForAll returns true after setApprovalForAll", func(t *testing.T) {
			setApprovalData := erc721.PackSetApprovalForAll(common.HexToAddress(otherAddress), true)
			txHash, err := w.SendTransaction(types.TransactionRequest{
				From:     initAddress,
				To:       contractHex,
				Value:    "0x0",
				GasLimit: 300000,
				Data:     setApprovalData,
			})
			assert.Nil(t, err)
			_, err = alchemy.Transact.WaitMined(context.Background(), txHash.Hex())
			assert.Nil(t, err)

			isApproved, err := alchemy.Nft.IsApprovedForAll(contractHex, initAddress, otherAddress)
			assert.Nil(t, err)
			assert.True(t, isApproved)
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

func TestScenario_DebugSnapshotRevertTo(t *testing.T) {
	// 1. snapshot
	snapshotId, err := alchemy.Debug.Snapshot()
	if err != nil {
		// evm_snapshot is only supported on development chains
		// (anvil, hardhat, ganache, ...), not on the kurtosis geth network.
		t.Skipf("node does not support evm_snapshot: %v", err)
	}

	balanceBefore, err := alchemy.Core.GetBalance(otherAddress, "latest")
	assert.Nil(t, err)

	// 2. transfer native
	w, err := wallet.New(initPrivateKey)
	assert.Nil(t, err)
	w.Connect(alchemy.GetProvider())

	txHash, err := w.SendTransaction(types.TransactionRequest{
		From:     initAddress,
		To:       otherAddress,
		Value:    "0x123",
		GasLimit: 300000,
	})
	assert.Nil(t, err)

	_, err = alchemy.Transact.WaitMined(context.Background(), txHash.Hex())
	assert.Nil(t, err)

	// 3. balance check
	balanceAfter, err := alchemy.Core.GetBalance(otherAddress, "latest")
	assert.Nil(t, err)
	assert.Equal(t, 1, balanceAfter.Cmp(balanceBefore))

	// 4. revert to
	reverted, err := alchemy.Debug.RevertTo(snapshotId)
	assert.Nil(t, err)
	assert.True(t, reverted)

	// 5. balance check
	balanceReverted, err := alchemy.Core.GetBalance(otherAddress, "latest")
	assert.Nil(t, err)
	assert.Equal(t, 0, balanceReverted.Cmp(balanceBefore))
}
