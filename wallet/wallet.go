package wallet

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"slices"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	gethTypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/poteto-go/go-alchemy-sdk/constant"
	"github.com/poteto-go/go-alchemy-sdk/internal"
	"github.com/poteto-go/go-alchemy-sdk/namespace"
	"github.com/poteto-go/go-alchemy-sdk/types"
	"github.com/poteto-go/go-alchemy-sdk/utils"
)

type wallet struct {
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
	provider   types.IAlchemyProvider
	mu         sync.RWMutex

	// Cache for performance (chainID and legacy-chain flag are immutable per network)
	cachedChainID *big.Int
	legacyChain   bool

	// for ERC20 / StableCoin
	erc20      namespace.IERC20
	stablecoin namespace.IStableCoin
}

func New(privateKeyStr string) (types.Wallet, error) {
	privateKeyStr = strings.TrimPrefix(privateKeyStr, "0x")

	privateKey, err := crypto.HexToECDSA(privateKeyStr)
	if err != nil {
		return &wallet{}, err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA := publicKey.(*ecdsa.PublicKey)
	return &wallet{
		privateKey: privateKey,
		publicKey:  publicKeyECDSA,
	}, nil
}

func (w *wallet) GetAddress() string {
	address := crypto.PubkeyToAddress(*w.publicKey)
	return common.HexToAddress(address.Hex()).String()
}

// snapshot returns the current provider under read lock.
// Callers should use the returned value for the remainder of the call so that
// a concurrent Connect cannot tear the interface header mid-read.
func (w *wallet) snapshot() types.IAlchemyProvider {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.provider
}

// snapshotERC20 returns the current erc20 namespace under read lock.
// Same rationale as snapshot — keeping it as its own getter so future
// namespaces (erc721, etc.) can be added without piling tuple returns
// onto a single helper.
func (w *wallet) snapshotERC20() namespace.IERC20 {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.erc20
}

func (w *wallet) snapshotStableCoin() namespace.IStableCoin {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.stablecoin
}

func (w *wallet) GetBalance() (*big.Int, error) {
	provider := w.snapshot()
	if provider == nil {
		return nil, constant.ErrWalletIsNotConnected
	}

	balance, err := provider.Eth().GetBalance(w.GetAddress(), "latest")
	if err != nil {
		return nil, err
	}
	return balance, nil
}

func (w *wallet) Connect(provider types.IAlchemyProvider) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.provider = provider
	w.erc20 = namespace.NewERC20Namespace(provider.Eth())
	w.stablecoin = namespace.NewStableCoinNamespace(provider.Eth())
}

func (w *wallet) PendingNonceAt() (uint64, error) {
	provider := w.snapshot()
	if provider == nil {
		return 0, constant.ErrWalletIsNotConnected
	}

	nonce, err := provider.Eth().PendingNonceAt(w.GetAddress())
	if err != nil {
		return nonce, err
	}

	return nonce, nil
}

// sign Transaction by wallet's p8 key
func (w *wallet) SignTx(txRequest types.TransactionRequest) (*gethTypes.Transaction, error) {
	provider := w.snapshot()
	if provider == nil {
		return nil, constant.ErrWalletIsNotConnected
	}

	nonce, err := w.PendingNonceAt()
	if err != nil {
		return nil, err
	}
	txRequest.Nonce = nonce

	estimatedGas, err := provider.Eth().EstimateGas(txRequest)
	if err != nil {
		return nil, err
	}
	if txRequest.GasLimit < estimatedGas.Uint64() {
		return nil, fmt.Errorf(
			"gasLimit(%d) is less than estimated gas %d",
			txRequest.GasLimit,
			estimatedGas.Uint64(),
		)
	}

	// MaxFeePerGas or MaxPriorityFeePerGas being set means the caller intends an
	// EIP-1559 (DynamicFeeTx) transaction, which uses those fields instead of
	// GasPrice. Only inject a legacy GasPrice when neither EIP-1559 field is set
	// and the caller has not supplied their own GasPrice.
	if txRequest.MaxFeePerGas == nil && txRequest.MaxPriorityFeePerGas == nil && txRequest.GasPrice == nil {
		gasPrice, err := provider.Eth().SuggestGasPrice()
		if err != nil {
			return nil, err
		}
		txRequest.GasPrice = gasPrice
	}

	chainID, _, err := w.chainID()
	if err != nil {
		return nil, err
	}
	txRequest.ChainID = chainID

	txData, err := utils.TransformTxRequestToGethTxData(txRequest)
	if err != nil {
		return nil, err
	}
	tx := gethTypes.NewTx(txData)

	latestEIP155Signer := gethTypes.LatestSignerForChainID(txRequest.ChainID)
	signedTx, err := gethTypes.SignTx(tx, latestEIP155Signer, w.privateKey)
	if err != nil {
		return nil, err
	}

	return signedTx, nil
}

func (w *wallet) SendTransaction(txRequest types.TransactionRequest) (common.Hash, error) {
	provider := w.snapshot()
	if provider == nil {
		return common.Hash{}, constant.ErrWalletIsNotConnected
	}

	signedTx, err := w.SignTx(txRequest)
	if err != nil {
		return common.Hash{}, err
	}

	if err := provider.Eth().SendRawTransaction(signedTx); err != nil {
		return common.Hash{}, err
	}

	return signedTx.Hash(), nil
}

func (w *wallet) DeployContract(ctx context.Context, metaData *bind.MetaData) (common.Address, error) {
	provider := w.snapshot()
	if provider == nil {
		return common.Address{}, constant.ErrWalletIsNotConnected
	}

	deployRes, err := w.DeployContractNoWait(metaData)
	if err != nil {
		return common.Address{}, err
	}

	tx := deployRes.Txs[metaData.ID]
	address, err := provider.Eth().WaitDeployed(ctx, tx.Hash())
	if err != nil {
		return common.Address{}, err
	}
	return address, nil
}

func (w *wallet) DeployContractNoWait(metaData *bind.MetaData) (*bind.DeploymentResult, error) {
	provider := w.snapshot()
	if provider == nil {
		return nil, constant.ErrWalletIsNotConnected
	}

	auth, err := w.buildAuth()
	if err != nil {
		return nil, err
	}

	deployRes, err := provider.Eth().DeployContract(auth, metaData)
	if err != nil {
		return nil, err
	}

	return deployRes, nil
}

func (w *wallet) ContractTransact(
	ctx context.Context,
	contract types.ContractInstance,
	contractAddress string,
	data []byte,
) (*gethTypes.Receipt, error) {
	provider := w.snapshot()
	if provider == nil {
		return nil, constant.ErrWalletIsNotConnected
	}

	tx, err := w.ContractTransactNoWait(
		contract,
		contractAddress,
		data,
	)
	if err != nil {
		return nil, err
	}

	return provider.Eth().WaitMined(ctx, tx.Hash())
}

func (w *wallet) ContractTransactNoWait(
	contract types.ContractInstance,
	contractAddress string,
	data []byte,
) (*gethTypes.Transaction, error) {
	provider := w.snapshot()
	if provider == nil {
		return nil, constant.ErrWalletIsNotConnected
	}

	auth, err := w.buildAuth()
	if err != nil {
		return nil, err
	}

	tx, err := provider.Eth().ContractTransact(auth, contract, contractAddress, data)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (w *wallet) ContractCall(
	contract types.ContractInstance,
	contractAddress string,
	opts *bind.CallOpts,
	callData []byte,
	unpack func([]byte) (any, error),
) (any, error) {
	provider := w.snapshot()
	if provider == nil {
		return nil, constant.ErrWalletIsNotConnected
	}

	addr := common.HexToAddress(contractAddress)
	return provider.Eth().ContractCall(contract, addr, opts, callData, unpack)
}

func (w *wallet) ERC20() types.WalletERC20 {
	return &walletERC20{w: w}
}

func (w *wallet) StableCoin() types.WalletStableCoin {
	return &walletStableCoin{walletERC20{w: w}}
}

func (w *wallet) buildAuth() (*bind.TransactOpts, error) {
	chainID, legacyChain, err := w.chainID()
	if err != nil {
		return nil, err
	}

	auth := bind.NewKeyedTransactor(w.privateKey, chainID)

	if legacyChain {
		gasPrice, err := w.snapshot().Eth().SuggestGasPrice()
		if err != nil {
			return nil, err
		}
		auth.GasPrice = gasPrice
	}

	return auth, nil
}

// chainID returns the cached chainID and legacy-chain flag, fetching them from
// the network on first use. It uses a double-checked pattern so the network
// round-trip happens outside the lock: a concurrent buildAuth / SignTx no longer
// blocks on the write lock for the duration of the RPC.
func (w *wallet) chainID() (*big.Int, bool, error) {
	w.mu.RLock()
	if w.cachedChainID != nil {
		chainID, legacyChain := w.cachedChainID, w.legacyChain
		w.mu.RUnlock()
		return chainID, legacyChain, nil
	}
	provider := w.provider
	w.mu.RUnlock()

	if provider == nil {
		return nil, false, constant.ErrWalletIsNotConnected
	}

	chainID, err := provider.Eth().ChainID()
	if err != nil {
		return nil, false, err
	}
	legacyChain := slices.Contains(internal.ChainListNotSupportEIP1559, chainID.Int64())

	w.mu.Lock()
	defer w.mu.Unlock()
	// Another goroutine may have populated the cache while we were fetching;
	// prefer the already-cached value so all callers observe the same chainID.
	if w.cachedChainID == nil {
		w.cachedChainID = chainID
		w.legacyChain = legacyChain
	}
	return w.cachedChainID, w.legacyChain, nil
}

func (w *wallet) ResetPool() {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.cachedChainID = nil
	w.legacyChain = false
}
