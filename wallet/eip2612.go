package wallet

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/poteto-go/go-alchemy-sdk/constant"
)

type walletEIP2612 struct {
	walletERC20
}

func (api *walletEIP2612) PermitNoWait(contractAddress, ownerAddress, spenderAddress string, value, deadline *big.Int, v uint8, r, s [32]byte, gasLimit *uint64) (common.Hash, error) {
	return api.sendERC20Tx(contractAddress, gasLimit, constant.PermitFnSignature,
		common.LeftPadBytes(common.HexToAddress(ownerAddress).Bytes(), constant.ABIWordSize),
		common.LeftPadBytes(common.HexToAddress(spenderAddress).Bytes(), constant.ABIWordSize),
		common.LeftPadBytes(value.Bytes(), constant.ABIWordSize),
		common.LeftPadBytes(deadline.Bytes(), constant.ABIWordSize),
		common.LeftPadBytes([]byte{v}, constant.ABIWordSize),
		r[:],
		s[:],
	)
}

func (api *walletEIP2612) Permit(ctx context.Context, contractAddress, ownerAddress, spenderAddress string, value, deadline *big.Int, v uint8, r, s [32]byte, gasLimit *uint64) (*gethTypes.Receipt, error) {
	return api.waitMined(ctx, func() (common.Hash, error) {
		return api.PermitNoWait(contractAddress, ownerAddress, spenderAddress, value, deadline, v, r, s, gasLimit)
	})
}

func (api *walletEIP2612) Nonces(contractAddress, ownerAddress string) (*big.Int, error) {
	eip2612 := api.w.snapshotEIP2612()
	if eip2612 == nil {
		return nil, constant.ErrWalletIsNotConnected
	}
	return eip2612.Nonces(contractAddress, ownerAddress)
}

func (api *walletEIP2612) DomainSeparator(contractAddress string) ([32]byte, error) {
	eip2612 := api.w.snapshotEIP2612()
	if eip2612 == nil {
		return [32]byte{}, constant.ErrWalletIsNotConnected
	}
	return eip2612.DomainSeparator(contractAddress)
}
