package wallet

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/poteto-go/go-alchemy-sdk/constant"
	"github.com/poteto-go/go-alchemy-sdk/encode"
	"github.com/poteto-go/go-alchemy-sdk/types"
)

type walletERC20 struct {
	w *wallet
}

func (api *walletERC20) waitMined(ctx context.Context, send func() (common.Hash, error)) (*gethTypes.Receipt, error) {
	provider := api.w.snapshot()
	if provider == nil {
		return nil, constant.ErrWalletIsNotConnected
	}
	txHash, err := send()
	if err != nil {
		return nil, err
	}
	return provider.Eth().WaitMined(ctx, txHash)
}

func (api *walletERC20) Transfer(ctx context.Context, contractAddress, toAddress string, amount *big.Int, gasLimit *uint64) (*gethTypes.Receipt, error) {
	return api.waitMined(ctx, func() (common.Hash, error) {
		return api.TransferNoWait(contractAddress, toAddress, amount, gasLimit)
	})
}

func resolveGasLimit(gasLimit *uint64) uint64 {
	if gasLimit == nil {
		return 0 // 0 is the sentinel for "auto": SignTx fills it from eth_estimateGas
	}
	return *gasLimit
}

func (api *walletERC20) sendERC20Tx(contractAddress string, gasLimit *uint64, sig []byte, params ...[]byte) (common.Hash, error) {
	if err := validateAddress(contractAddress); err != nil {
		return common.Hash{}, err
	}
	if api.w.snapshot() == nil {
		return common.Hash{}, constant.ErrWalletIsNotConnected
	}
	return api.w.SendTransaction(types.TransactionRequest{
		From:     api.w.GetAddress(),
		To:       contractAddress,
		Value:    "0x0",
		GasLimit: resolveGasLimit(gasLimit),
		Data:     encode.ReadCalldata(sig, params...),
	})
}

func (api *walletERC20) TransferNoWait(contractAddress, toAddress string, amount *big.Int, gasLimit *uint64) (common.Hash, error) {
	if err := validateAddress(toAddress); err != nil {
		return common.Hash{}, err
	}
	if err := validateUint256(amount); err != nil {
		return common.Hash{}, err
	}
	return api.sendERC20Tx(contractAddress, gasLimit, constant.TransferFnSignature,
		common.LeftPadBytes(common.HexToAddress(toAddress).Bytes(), constant.ABIWordSize),
		common.LeftPadBytes(amount.Bytes(), constant.ABIWordSize),
	)
}

func (api *walletERC20) ApproveNoWait(contractAddress, spenderAddress string, amount *big.Int, gasLimit *uint64) (common.Hash, error) {
	if err := validateAddress(spenderAddress); err != nil {
		return common.Hash{}, err
	}
	if err := validateUint256(amount); err != nil {
		return common.Hash{}, err
	}
	return api.sendERC20Tx(contractAddress, gasLimit, constant.ApproveFnSignature,
		common.LeftPadBytes(common.HexToAddress(spenderAddress).Bytes(), constant.ABIWordSize),
		common.LeftPadBytes(amount.Bytes(), constant.ABIWordSize),
	)
}

func (api *walletERC20) Approve(ctx context.Context, contractAddress, spenderAddress string, amount *big.Int, gasLimit *uint64) (*gethTypes.Receipt, error) {
	return api.waitMined(ctx, func() (common.Hash, error) {
		return api.ApproveNoWait(contractAddress, spenderAddress, amount, gasLimit)
	})
}

func (api *walletERC20) TransferFromNoWait(contractAddress, fromAddress, toAddress string, amount *big.Int, gasLimit *uint64) (common.Hash, error) {
	if err := validateAddress(fromAddress); err != nil {
		return common.Hash{}, err
	}
	if err := validateAddress(toAddress); err != nil {
		return common.Hash{}, err
	}
	if err := validateUint256(amount); err != nil {
		return common.Hash{}, err
	}
	return api.sendERC20Tx(contractAddress, gasLimit, constant.TransferFromFnSignature,
		common.LeftPadBytes(common.HexToAddress(fromAddress).Bytes(), constant.ABIWordSize),
		common.LeftPadBytes(common.HexToAddress(toAddress).Bytes(), constant.ABIWordSize),
		common.LeftPadBytes(amount.Bytes(), constant.ABIWordSize),
	)
}

func (api *walletERC20) TransferFrom(ctx context.Context, contractAddress, fromAddress, toAddress string, amount *big.Int, gasLimit *uint64) (*gethTypes.Receipt, error) {
	return api.waitMined(ctx, func() (common.Hash, error) {
		return api.TransferFromNoWait(contractAddress, fromAddress, toAddress, amount, gasLimit)
	})
}

func (api *walletERC20) BalanceOf(contractAddress string) (*big.Int, error) {
	erc20 := api.w.snapshotERC20()
	if erc20 == nil {
		return nil, constant.ErrWalletIsNotConnected
	}

	return erc20.BalanceOf(
		contractAddress,
		api.w.GetAddress(),
	)
}

func (api *walletERC20) TotalSupply(contractAddress string) (*big.Int, error) {
	erc20 := api.w.snapshotERC20()
	if erc20 == nil {
		return nil, constant.ErrWalletIsNotConnected
	}

	return erc20.TotalSupply(contractAddress)
}

func (api *walletERC20) Allowance(contractAddress, owner, spender string) (*big.Int, error) {
	erc20 := api.w.snapshotERC20()
	if erc20 == nil {
		return nil, constant.ErrWalletIsNotConnected
	}

	return erc20.Allowance(contractAddress, owner, spender)
}

func (api *walletERC20) Name(contractAddress string) (string, error) {
	erc20 := api.w.snapshotERC20()
	if erc20 == nil {
		return "", constant.ErrWalletIsNotConnected
	}

	return erc20.Name(contractAddress)
}

func (api *walletERC20) Symbol(contractAddress string) (string, error) {
	erc20 := api.w.snapshotERC20()
	if erc20 == nil {
		return "", constant.ErrWalletIsNotConnected
	}

	return erc20.Symbol(contractAddress)
}

func (api *walletERC20) Decimals(contractAddress string) (uint8, error) {
	erc20 := api.w.snapshotERC20()
	if erc20 == nil {
		return 0, constant.ErrWalletIsNotConnected
	}

	return erc20.Decimals(contractAddress)
}
