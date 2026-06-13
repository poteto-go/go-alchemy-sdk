package wallet

import (
	"math/big"

	"github.com/poteto-go/go-alchemy-sdk/constant"
)

type walletNft struct {
	w *wallet
}

func (api *walletNft) OwnerOf(contractAddress string, tokenId *big.Int) (string, error) {
	nft := api.w.snapshotNft()
	if nft == nil {
		return "", constant.ErrWalletIsNotConnected
	}

	return nft.OwnerOf(contractAddress, tokenId)
}

func (api *walletNft) TokenURI(contractAddress string, tokenId *big.Int) (string, error) {
	nft := api.w.snapshotNft()
	if nft == nil {
		return "", constant.ErrWalletIsNotConnected
	}

	return nft.TokenURI(contractAddress, tokenId)
}

func (api *walletNft) Name(contractAddress string) (string, error) {
	nft := api.w.snapshotNft()
	if nft == nil {
		return "", constant.ErrWalletIsNotConnected
	}

	return nft.Name(contractAddress)
}

func (api *walletNft) Symbol(contractAddress string) (string, error) {
	nft := api.w.snapshotNft()
	if nft == nil {
		return "", constant.ErrWalletIsNotConnected
	}

	return nft.Symbol(contractAddress)
}

func (api *walletNft) GetApproved(contractAddress string, tokenId *big.Int) (string, error) {
	nft := api.w.snapshotNft()
	if nft == nil {
		return "", constant.ErrWalletIsNotConnected
	}

	return nft.GetApproved(contractAddress, tokenId)
}

func (api *walletNft) IsApprovedForAll(contractAddress, owner, operator string) (bool, error) {
	nft := api.w.snapshotNft()
	if nft == nil {
		return false, constant.ErrWalletIsNotConnected
	}

	return nft.IsApprovedForAll(contractAddress, owner, operator)
}
