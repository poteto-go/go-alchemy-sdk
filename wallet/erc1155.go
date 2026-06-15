package wallet

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/poteto-go/go-alchemy-sdk/constant"
	"github.com/poteto-go/go-alchemy-sdk/encode"
)

type walletErc1155 struct {
	*walletNft
}

var erc1155SafeTransferDataOffsetWord = encode.ABIUint256(big.NewInt(constant.Erc1155SafeTransferFromHeadSize))

func (api *walletErc1155) BalanceOfToken(contractAddress, account string, tokenId *big.Int) (*big.Int, error) {
	erc1155 := api.w.snapshotErc1155()
	if erc1155 == nil {
		return nil, constant.ErrWalletIsNotConnected
	}

	return erc1155.BalanceOfToken(contractAddress, account, tokenId)
}

func (api *walletErc1155) BalanceOfBatch(contractAddress string, accounts []string, tokenIds []*big.Int) ([]*big.Int, error) {
	erc1155 := api.w.snapshotErc1155()
	if erc1155 == nil {
		return nil, constant.ErrWalletIsNotConnected
	}

	return erc1155.BalanceOfBatch(contractAddress, accounts, tokenIds)
}

func (api *walletErc1155) Uri(contractAddress string, tokenId *big.Int) (string, error) {
	erc1155 := api.w.snapshotErc1155()
	if erc1155 == nil {
		return "", constant.ErrWalletIsNotConnected
	}

	return erc1155.Uri(contractAddress, tokenId)
}

func (api *walletErc1155) SafeTransferFrom(ctx context.Context, contractAddress, fromAddress, toAddress string, id, amount *big.Int, data []byte, gasLimit *uint64) (*gethTypes.Receipt, error) {
	return api.waitMined(ctx, func() (common.Hash, error) {
		return api.SafeTransferFromNoWait(contractAddress, fromAddress, toAddress, id, amount, data, gasLimit)
	})
}

func (api *walletErc1155) SafeTransferFromNoWait(contractAddress, fromAddress, toAddress string, id, amount *big.Int, data []byte, gasLimit *uint64) (common.Hash, error) {
	if err := validateTransferArgs(fromAddress, toAddress, id); err != nil {
		return common.Hash{}, err
	}
	if err := validateUint256(amount); err != nil {
		return common.Hash{}, err
	}
	return api.sendNftTx(contractAddress, gasLimit, constant.Erc1155SafeTransferFromFnSignature,
		encode.ABIAddress(fromAddress),
		encode.ABIAddress(toAddress),
		encode.ABIUint256(id),
		encode.ABIUint256(amount),
		erc1155SafeTransferDataOffsetWord,
		encode.ABIBytes(data),
	)
}

func (api *walletErc1155) SafeBatchTransferFrom(ctx context.Context, contractAddress, fromAddress, toAddress string, ids, amounts []*big.Int, data []byte, gasLimit *uint64) (*gethTypes.Receipt, error) {
	return api.waitMined(ctx, func() (common.Hash, error) {
		return api.SafeBatchTransferFromNoWait(contractAddress, fromAddress, toAddress, ids, amounts, data, gasLimit)
	})
}

func (api *walletErc1155) SafeBatchTransferFromNoWait(contractAddress, fromAddress, toAddress string, ids, amounts []*big.Int, data []byte, gasLimit *uint64) (common.Hash, error) {
	if err := validateAddressPair(fromAddress, toAddress); err != nil {
		return common.Hash{}, err
	}
	if len(ids) != len(amounts) {
		return common.Hash{}, constant.ErrMismatchedArrayLength
	}
	for i, id := range ids {
		if err := validateUint256(id); err != nil {
			return common.Hash{}, err
		}
		if err := validateUint256(amounts[i]); err != nil {
			return common.Hash{}, err
		}
	}

	// safeBatchTransferFrom(address,address,uint256[],uint256[],bytes):
	// head: from(static), to(static), offsetIds, offsetAmounts, offsetData
	// tails follow in order.
	idsTail := encode.ABIUint256Array(ids)
	amountsTail := encode.ABIUint256Array(amounts)
	dataTail := encode.ABIBytes(data)
	headSize := constant.Erc1155SafeTransferFromHeadSize
	args := make([]byte, 0, headSize+len(idsTail)+len(amountsTail)+len(dataTail))
	args = append(args, encode.ABIAddress(fromAddress)...)
	args = append(args, encode.ABIAddress(toAddress)...)
	args = append(args, encode.ABIUint256(big.NewInt(int64(headSize)))...)
	args = append(args, encode.ABIUint256(big.NewInt(int64(headSize+len(idsTail))))...)
	args = append(args, encode.ABIUint256(big.NewInt(int64(headSize+len(idsTail)+len(amountsTail))))...)
	args = append(args, idsTail...)
	args = append(args, amountsTail...)
	args = append(args, dataTail...)

	return api.sendNftTx(contractAddress, gasLimit, constant.SafeBatchTransferFromFnSignature, args)
}
