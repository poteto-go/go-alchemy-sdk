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
		transfer erc20 token from another address (requires prior approval)
			- wait for mined
			- gas limit is 300000 for default
	*/
	TransferFrom(contractAddress, fromAddress, toAddress string, amount *big.Int, gasLimit *uint64) (*gethTypes.Receipt, error)

	/*
		transfer erc20 token from another address (requires prior approval)
			- gas limit is 300000 for default
	*/
	TransferFromNoWait(contractAddress, fromAddress, toAddress string, amount *big.Int, gasLimit *uint64) (common.Hash, error)

	/*
		approve spender to spend erc20 token
			- wait for mined
			- gas limit is 300000 for default
	*/
	Approve(contractAddress, spenderAddress string, amount *big.Int, gasLimit *uint64) (*gethTypes.Receipt, error)

	/*
		approve spender to spend erc20 token
			- gas limit is 300000 for default
	*/
	ApproveNoWait(contractAddress, spenderAddress string, amount *big.Int, gasLimit *uint64) (common.Hash, error)

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
	provider, _ := api.w.snapshot()
	if provider == nil {
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

	return provider.Eth().WaitMined(txHash)
}

func buildERC20TxData(signature []byte, params ...[]byte) ([]byte, error) {
	hash := sha3.NewLegacyKeccak256()
	if _, err := hash.Write(signature); err != nil {
		return nil, err
	}
	methodID := hash.Sum(nil)[:4]

	size := 4
	for _, p := range params {
		size += len(p)
	}
	data := make([]byte, 0, size)
	data = append(data, methodID...)
	for _, p := range params {
		data = append(data, p...)
	}
	return data, nil
}

func resolveGasLimit(gasLimit *uint64) uint64 {
	if gasLimit == nil {
		return 300000
	}
	return *gasLimit
}

func (api *walletERC20) TransferNoWait(contractAddress, toAddress string, amount *big.Int, gasLimit *uint64) (common.Hash, error) {
	provider, _ := api.w.snapshot()
	if provider == nil {
		return common.Hash{}, constant.ErrWalletIsNotConnected
	}

	data, err := buildERC20TxData(
		constant.TransferFnSignature,
		common.LeftPadBytes(common.HexToAddress(toAddress).Bytes(), 32),
		common.LeftPadBytes(amount.Bytes(), 32),
	)
	if err != nil {
		return common.Hash{}, err
	}

	txRequest := types.TransactionRequest{
		From:     api.w.GetAddress(),
		To:       contractAddress,
		Value:    "0x0",
		GasLimit: resolveGasLimit(gasLimit),
		Data:     data,
	}

	return api.w.SendTransaction(txRequest)
}

func (api *walletERC20) ApproveNoWait(contractAddress, spenderAddress string, amount *big.Int, gasLimit *uint64) (common.Hash, error) {
	provider, _ := api.w.snapshot()
	if provider == nil {
		return common.Hash{}, constant.ErrWalletIsNotConnected
	}

	data, err := buildERC20TxData(
		constant.ApproveFnSignature,
		common.LeftPadBytes(common.HexToAddress(spenderAddress).Bytes(), 32),
		common.LeftPadBytes(amount.Bytes(), 32),
	)
	if err != nil {
		return common.Hash{}, err
	}

	txRequest := types.TransactionRequest{
		From:     api.w.GetAddress(),
		To:       contractAddress,
		Value:    "0x0",
		GasLimit: resolveGasLimit(gasLimit),
		Data:     data,
	}

	return api.w.SendTransaction(txRequest)
}

func (api *walletERC20) Approve(contractAddress, spenderAddress string, amount *big.Int, gasLimit *uint64) (*gethTypes.Receipt, error) {
	provider, _ := api.w.snapshot()
	if provider == nil {
		return nil, constant.ErrWalletIsNotConnected
	}

	txHash, err := api.ApproveNoWait(contractAddress, spenderAddress, amount, gasLimit)
	if err != nil {
		return nil, err
	}

	return provider.Eth().WaitMined(txHash)
}

func (api *walletERC20) TransferFromNoWait(contractAddress, fromAddress, toAddress string, amount *big.Int, gasLimit *uint64) (common.Hash, error) {
	provider, _ := api.w.snapshot()
	if provider == nil {
		return common.Hash{}, constant.ErrWalletIsNotConnected
	}

	data, err := buildERC20TxData(
		constant.TransferFromFnSignature,
		common.LeftPadBytes(common.HexToAddress(fromAddress).Bytes(), 32),
		common.LeftPadBytes(common.HexToAddress(toAddress).Bytes(), 32),
		common.LeftPadBytes(amount.Bytes(), 32),
	)
	if err != nil {
		return common.Hash{}, err
	}

	txRequest := types.TransactionRequest{
		From:     api.w.GetAddress(),
		To:       contractAddress,
		Value:    "0x0",
		GasLimit: resolveGasLimit(gasLimit),
		Data:     data,
	}

	return api.w.SendTransaction(txRequest)
}

func (api *walletERC20) TransferFrom(contractAddress, fromAddress, toAddress string, amount *big.Int, gasLimit *uint64) (*gethTypes.Receipt, error) {
	provider, _ := api.w.snapshot()
	if provider == nil {
		return nil, constant.ErrWalletIsNotConnected
	}

	txHash, err := api.TransferFromNoWait(contractAddress, fromAddress, toAddress, amount, gasLimit)
	if err != nil {
		return nil, err
	}

	return provider.Eth().WaitMined(txHash)
}

func (api *walletERC20) BalanceOf(contractAddress string) (*big.Int, error) {
	provider, erc20 := api.w.snapshot()
	if provider == nil {
		return nil, constant.ErrWalletIsNotConnected
	}

	return erc20.BalanceOf(
		contractAddress,
		api.w.GetAddress(),
	)
}

func (api *walletERC20) TotalSupply(contractAddress string) (*big.Int, error) {
	provider, erc20 := api.w.snapshot()
	if provider == nil {
		return nil, constant.ErrWalletIsNotConnected
	}

	return erc20.TotalSupply(contractAddress)
}

func (api *walletERC20) Allowance(contractAddress, owner, spender string) (*big.Int, error) {
	provider, erc20 := api.w.snapshot()
	if provider == nil {
		return nil, constant.ErrWalletIsNotConnected
	}

	return erc20.Allowance(contractAddress, owner, spender)
}

func (api *walletERC20) Name(contractAddress string) (string, error) {
	provider, erc20 := api.w.snapshot()
	if provider == nil {
		return "", constant.ErrWalletIsNotConnected
	}

	return erc20.Name(contractAddress)
}

func (api *walletERC20) Symbol(contractAddress string) (string, error) {
	provider, erc20 := api.w.snapshot()
	if provider == nil {
		return "", constant.ErrWalletIsNotConnected
	}

	return erc20.Symbol(contractAddress)
}

func (api *walletERC20) Decimals(contractAddress string) (uint8, error) {
	provider, erc20 := api.w.snapshot()
	if provider == nil {
		return 0, constant.ErrWalletIsNotConnected
	}

	return erc20.Decimals(contractAddress)
}
