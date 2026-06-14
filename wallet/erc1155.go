package wallet

import (
	"math/big"

	"github.com/poteto-go/go-alchemy-sdk/constant"
)

type walletErc1155 struct {
	*walletNft
}

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
