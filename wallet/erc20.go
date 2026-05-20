package wallet

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/poteto-go/go-alchemy-sdk/constant"
	"github.com/poteto-go/go-alchemy-sdk/types"
	"golang.org/x/crypto/sha3"
)

// ERC20 interface for wallet.
// This is only defined for UX.
type WalletERC20 interface {
	/*
		transfer erc20 token by provided wallet
			- wait for mined
			- gas limit is 300000 for default
	*/
	Transfer(contractAddress, toAddress string, amount *big.Int, gasLimit *uint64) (*gethTypes.Receipt, error)

	/*
		transfer erc20 token by provided wallet
			- gas limit is 300000 for default
	*/
	TransferNoWait(contractAddress, toAddress string, amount *big.Int, gasLimit *uint64) (common.Hash, error)

	/*
		get balance of provided wallet & erc20 token
	*/
	BalanceOf(contractAddress string) (*big.Int, error)

	/*
		get total supply of erc20 token
	*/
	TotalSupply(contractAddress string) (*big.Int, error)

	/*
		get allowance of erc20 token
	*/
	Allowance(contractAddress, owner, spender string) (*big.Int, error)

	/*
		get name of erc20 token
	*/
	Name(contractAddress string) (string, error)

	/*
		get symbol of erc20 token
	*/
	Symbol(contractAddress string) (string, error)

	/*
		get decimals of erc20 token
	*/
	Decimals(contractAddress string) (uint8, error)
}

type walletERC20 struct {
	w *wallet
}

func (api *walletERC20) Transfer(contractAddress, toAddress string, amount *big.Int, gasLimit *uint64) (*gethTypes.Receipt, error) {
	if api.w.provider == nil {
		return nil, constant.ErrWalletIsNotConnected
	}

	txHash, err := api.TransferNoWait(
		contractAddress,
		toAddress,
		amount,
		gasLimit,
	)
	if err != nil {
		return nil, err
	}

	return api.w.provider.Eth().WaitMined(txHash)
}

func (api *walletERC20) TransferNoWait(contractAddress, toAddress string, amount *big.Int, gasLimit *uint64) (common.Hash, error) {
	if api.w.provider == nil {
		return common.Hash{}, constant.ErrWalletIsNotConnected
	}

	hash := sha3.NewLegacyKeccak256()
	if _, err := hash.Write(constant.TransferFnSignature); err != nil {
		return common.Hash{}, err
	}
	methodID := hash.Sum(nil)[:4]

	paddedAddress := common.LeftPadBytes(common.HexToAddress(toAddress).Bytes(), 32)
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)

	data := make([]byte, 0, 68)
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)

	var txGasLimit uint64
	if gasLimit == nil {
		txGasLimit = 300000
	} else {
		txGasLimit = *gasLimit
	}

	txRequest := types.TransactionRequest{
		From:     api.w.GetAddress(),
		To:       contractAddress,
		Value:    "0x0",
		GasLimit: txGasLimit,
		Data:     data,
	}

	txHash, err := api.w.SendTransaction(
		txRequest,
	)
	if err != nil {
		return common.Hash{}, err
	}

	return txHash, nil
}

func (api *walletERC20) BalanceOf(contractAddress string) (*big.Int, error) {
	if api.w.provider == nil {
		return nil, constant.ErrWalletIsNotConnected
	}

	return api.w.erc20.BalanceOf(
		contractAddress,
		api.w.GetAddress(),
	)
}

func (api *walletERC20) TotalSupply(contractAddress string) (*big.Int, error) {
	if api.w.provider == nil {
		return nil, constant.ErrWalletIsNotConnected
	}

	return api.w.erc20.TotalSupply(contractAddress)
}

func (api *walletERC20) Allowance(contractAddress, owner, spender string) (*big.Int, error) {
	if api.w.provider == nil {
		return nil, constant.ErrWalletIsNotConnected
	}

	return api.w.erc20.Allowance(contractAddress, owner, spender)
}

func (api *walletERC20) Name(contractAddress string) (string, error) {
	if api.w.provider == nil {
		return "", constant.ErrWalletIsNotConnected
	}

	return api.w.erc20.Name(contractAddress)
}

func (api *walletERC20) Symbol(contractAddress string) (string, error) {
	if api.w.provider == nil {
		return "", constant.ErrWalletIsNotConnected
	}

	return api.w.erc20.Symbol(contractAddress)
}

func (api *walletERC20) Decimals(contractAddress string) (uint8, error) {
	if api.w.provider == nil {
		return 0, constant.ErrWalletIsNotConnected
	}

	return api.w.erc20.Decimals(contractAddress)
}
