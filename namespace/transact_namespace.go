package namespace

import (
	"context"

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
	WaitMined(ctx context.Context, txHash string) (*gethTypes.Receipt, error)

	/*
		WaitDeployed waits for a contract deployment transaction with the provided hash and
		returns the contract address.
		It stops waiting when ctx is canceled.
	*/
	WaitDeployed(ctx context.Context, txHash string) (common.Address, error)
}

type Transact struct {
	ether types.EtherApi
}

func NewTransactNamespace(ether types.EtherApi) ITransact {
	return &Transact{
		ether: ether,
	}
}

func (t *Transact) WaitMined(ctx context.Context, txHash string) (*gethTypes.Receipt, error) {
	txReceipt, err := t.ether.WaitMined(ctx, common.HexToHash(txHash))
	if err != nil {
		return nil, err
	}

	return txReceipt, nil
}

func (t *Transact) WaitDeployed(ctx context.Context, txHash string) (common.Address, error) {
	address, err := t.ether.WaitDeployed(ctx, common.HexToHash(txHash))
	if err != nil {
		return common.Address{}, err
	}

	return address, nil
}
