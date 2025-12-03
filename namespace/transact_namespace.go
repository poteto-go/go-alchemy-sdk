package namespace

import (
	"github.com/ethereum/go-ethereum/common"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/poteto-go/go-alchemy-sdk/types"
)

type ITransact interface {
	/*
		WaitMined waits for a transaction with the provided hash and
		returns the transaction receipt when it is mined.
		It stops waiting when ctx is canceled.
	*/
	WaitMined(txHash string) (*gethTypes.Receipt, error)

	/*
		WaitDeployed waits for a contract deployment transaction with the provided hash and
		returns the contract address
		It stops waiting when ctx is canceled.
	*/
	WaitDeployed(txHash string) (common.Address, error)
}

type Transact struct {
	ether types.EtherApi
}

func NewTransactNamespace(ether types.EtherApi) ITransact {
	return &Transact{
		ether: ether,
	}
}

func (t *Transact) WaitMined(txHash string) (*gethTypes.Receipt, error) {
	txReceipt, err := t.ether.WaitMined(common.HexToHash(txHash))
	if err != nil {
		return nil, err
	}

	return txReceipt, nil
}

func (t *Transact) WaitDeployed(txHash string) (common.Address, error) {
	address, err := t.ether.WaitDeployed(common.HexToHash(txHash))
	if err != nil {
		return common.Address{}, err
	}

	return address, nil
}
