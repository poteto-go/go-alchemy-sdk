package wallet

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	gethTypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/poteto-go/go-alchemy-sdk/constant"
	"github.com/poteto-go/go-alchemy-sdk/types"
	"github.com/poteto-go/go-alchemy-sdk/utils"
)

// Wallet class inherits Signer and can sign transactions and messages using
type Wallet interface {
	// get address of wallet
	GetAddress() string

	// get balance of native token
	GetBalance() (balance *big.Int, err error)

	// connect provider to wallet
	Connect(provider types.IAlchemyProvider)

	/*
		PendingNonceAt returns the account nonce of the given account in the pending state.
		This is the nonce that should be used for the next transaction.

		internal call geth
	*/
	PendingNonceAt() (nonce uint64, err error)

	/*
		sign Transaction by wallet's p8 key
		using latest EIP155Signer

		EIP155Signer sign w/ ChainID to protect replay-attack
	*/
	SignTx(txRequest types.TransactionRequest) (signedTx *gethTypes.Transaction, err error)

	// Signs tx and sends it to the pending pool for execution
	SendTransaction(txRequest types.TransactionRequest) (err error)

	/*
		DeployContract creates and submits a deployment transaction based on the
		deployer bytecode.
		It returns the address and creation transaction of the pending contract,
		or an error if the creation failed.

		It does not work on non-Ethernet compatible networks.
	*/
	DeployContract(metaData *bind.MetaData) (common.Address, error)
}

type wallet struct {
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
	provider   types.IAlchemyProvider
	mu         sync.RWMutex
}

func New(privateKeyStr string) (Wallet, error) {
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

func (w *wallet) GetBalance() (*big.Int, error) {
	if w.provider == nil {
		return big.NewInt(0), constant.ErrWalletIsNotConnected
	}

	balance, err := w.provider.Eth().GetBalance(w.GetAddress(), "latest")
	if err != nil {
		return nil, err
	}
	return balance, nil
}

func (w *wallet) Connect(provider types.IAlchemyProvider) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.provider = provider
}

func (w *wallet) PendingNonceAt() (uint64, error) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	if w.provider == nil {
		return 0, constant.ErrWalletIsNotConnected
	}

	nonce, err := w.provider.Eth().PendingNonceAt(w.GetAddress())
	if err != nil {
		return nonce, err
	}

	return nonce, nil
}

// sign Transaction by wallet's p8 key
func (w *wallet) SignTx(txRequest types.TransactionRequest) (*gethTypes.Transaction, error) {
	if w.provider == nil {
		return nil, constant.ErrWalletIsNotConnected
	}

	nonce, err := w.PendingNonceAt()
	if err != nil {
		return nil, err
	}
	txRequest.Nonce = nonce

	gasPrice, err := w.provider.Eth().EstimateGas(txRequest)
	if err != nil {
		return nil, err
	}
	txRequest.GasPrice = gasPrice

	if txRequest.GasLimit <= txRequest.GasPrice.Uint64() {
		return nil, fmt.Errorf(
			"gasLimit(%d) is expected over estimated gasPrice %d",
			txRequest.GasLimit,
			txRequest.GasPrice.Uint64(),
		)
	}

	txData, err := utils.TransformTxRequestToGethTxData(txRequest)
	if err != nil {
		return nil, err
	}
	tx := gethTypes.NewTx(txData)

	latestEIP155Signer := gethTypes.LatestSignerForChainID(txData.ChainID)
	signedTx, err := gethTypes.SignTx(tx, latestEIP155Signer, w.privateKey)
	if err != nil {
		return nil, err
	}

	return signedTx, nil
}

func (w *wallet) SendTransaction(txRequest types.TransactionRequest) error {
	if w.provider == nil {
		return constant.ErrWalletIsNotConnected
	}

	signedTx, err := w.SignTx(txRequest)
	if err != nil {
		return err
	}

	if err := w.provider.Eth().SendRawTransaction(signedTx); err != nil {
		return err
	}

	return nil
}

func (w *wallet) DeployContract(metaData *bind.MetaData) (common.Address, error) {
	if w.provider == nil {
		return common.Address{}, constant.ErrWalletIsNotConnected
	}

	chainID, err := w.provider.Eth().ChainID()
	if err != nil {
		return common.Address{}, err
	}
	auth := bind.NewKeyedTransactor(w.privateKey, chainID)
	address, err := w.provider.Eth().DeployContract(auth, metaData)
	if err != nil {
		return common.Address{}, err
	}
	return address, nil
}
