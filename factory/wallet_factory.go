package factory

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/poteto-go/go-alchemy-sdk/types"
	"github.com/poteto-go/go-alchemy-sdk/wallet"
)

func ContractCall[T any](
	w wallet.Wallet,
	contract types.ContractInstance,
	contractAddress string,
	opts *bind.CallOpts,
	callData []byte,
	unpack func([]byte) (T, error),
) (T, error) {
	res, err := w.ContractCall(contract, contractAddress, opts, callData, func(b []byte) (any, error) {
		return unpack(b)
	})
	if err != nil {
		return *new(T), err
	}

	val, ok := res.(T)
	if !ok {
		return *new(T), fmt.Errorf("failed to cast result to expected type")
	}

	return val, nil
}
