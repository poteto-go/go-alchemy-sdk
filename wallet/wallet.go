package wallet

import (
	"crypto/ecdsa"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/poteto-go/go-alchemy-sdk/types"
)

// Wallet class inherits Signer and can sign transactions and messages using
type Wallet interface {
	GetAddress() string

	// connect provider to wallet
	Connect(provider types.IAlchemyProvider)

	/*
		PendingNonceAt returns the account nonce of the given account in the pending state.
		This is the nonce that should be used for the next transaction.

		internal call geth
	*/
	PendingNonceAt() (nonce uint64, err error)
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

func (w *wallet) Connect(provider types.IAlchemyProvider) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.provider = provider
}

func (w *wallet) PendingNonceAt() (uint64, error) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	nonce, err := w.provider.Eth().PendingNonceAt(w.GetAddress())
	if err != nil {
		return nonce, err
	}

	return nonce, nil
}
