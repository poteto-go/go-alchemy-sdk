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
	GetAddress() common.Address

	// connect provider to wallet
	Connect(provider types.IAlchemyProvider)
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

func (w *wallet) GetAddress() common.Address {
	return crypto.PubkeyToAddress(*w.publicKey)
}

func (w *wallet) Connect(provider types.IAlchemyProvider) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.provider = provider
}
