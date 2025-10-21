package wallet

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type Wallet interface {
	GetAddress() common.Address
}

type wallet struct {
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
	// mu         sync.RWMutex
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
