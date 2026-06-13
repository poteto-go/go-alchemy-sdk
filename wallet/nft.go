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

type walletNft struct {
	w *wallet
}

// safeTransferDataOffsetWord is the ABI offset word for the `bytes` argument of
// safeTransferFrom(address,address,uint256,bytes). The head holds four words
// (from, to, tokenId, offset), so the tail begins at 4*32 bytes. Precomputed
// once since it is a constant for this fixed signature.
var safeTransferDataOffsetWord = encode.ABIUint256(big.NewInt(4 * constant.ABIWordSize))

func (api *walletNft) waitMined(ctx context.Context, send func() (common.Hash, error)) (*gethTypes.Receipt, error) {
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

func (api *walletNft) sendNftTx(contractAddress string, gasLimit *uint64, sig []byte, params ...[]byte) (common.Hash, error) {
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

// validateTransferArgs validates the from/to addresses and tokenId shared by
// every ERC721 transfer variant.
func validateTransferArgs(fromAddress, toAddress string, tokenId *big.Int) error {
	if err := validateAddress(fromAddress); err != nil {
		return err
	}
	if err := validateAddress(toAddress); err != nil {
		return err
	}
	return validateUint256(tokenId)
}

func (api *walletNft) TransferFrom(ctx context.Context, contractAddress, fromAddress, toAddress string, tokenId *big.Int, gasLimit *uint64) (*gethTypes.Receipt, error) {
	return api.waitMined(ctx, func() (common.Hash, error) {
		return api.TransferFromNoWait(contractAddress, fromAddress, toAddress, tokenId, gasLimit)
	})
}

func (api *walletNft) TransferFromNoWait(contractAddress, fromAddress, toAddress string, tokenId *big.Int, gasLimit *uint64) (common.Hash, error) {
	if err := validateTransferArgs(fromAddress, toAddress, tokenId); err != nil {
		return common.Hash{}, err
	}
	return api.sendNftTx(contractAddress, gasLimit, constant.TransferFromFnSignature,
		encode.ABIAddress(fromAddress),
		encode.ABIAddress(toAddress),
		encode.ABIUint256(tokenId),
	)
}

func (api *walletNft) SafeTransferFrom(ctx context.Context, contractAddress, fromAddress, toAddress string, tokenId *big.Int, gasLimit *uint64) (*gethTypes.Receipt, error) {
	return api.waitMined(ctx, func() (common.Hash, error) {
		return api.SafeTransferFromNoWait(contractAddress, fromAddress, toAddress, tokenId, gasLimit)
	})
}

func (api *walletNft) SafeTransferFromNoWait(contractAddress, fromAddress, toAddress string, tokenId *big.Int, gasLimit *uint64) (common.Hash, error) {
	if err := validateTransferArgs(fromAddress, toAddress, tokenId); err != nil {
		return common.Hash{}, err
	}
	return api.sendNftTx(contractAddress, gasLimit, constant.SafeTransferFromFnSignature,
		encode.ABIAddress(fromAddress),
		encode.ABIAddress(toAddress),
		encode.ABIUint256(tokenId),
	)
}

func (api *walletNft) SafeTransferFromWithData(ctx context.Context, contractAddress, fromAddress, toAddress string, tokenId *big.Int, data []byte, gasLimit *uint64) (*gethTypes.Receipt, error) {
	return api.waitMined(ctx, func() (common.Hash, error) {
		return api.SafeTransferFromWithDataNoWait(contractAddress, fromAddress, toAddress, tokenId, data, gasLimit)
	})
}

func (api *walletNft) SafeTransferFromWithDataNoWait(contractAddress, fromAddress, toAddress string, tokenId *big.Int, data []byte, gasLimit *uint64) (common.Hash, error) {
	if err := validateTransferArgs(fromAddress, toAddress, tokenId); err != nil {
		return common.Hash{}, err
	}
	// `data` is a dynamic `bytes` argument: emit the offset word pointing past
	// the head, then the length+data tail.
	return api.sendNftTx(contractAddress, gasLimit, constant.SafeTransferFromWithDataFnSignature,
		encode.ABIAddress(fromAddress),
		encode.ABIAddress(toAddress),
		encode.ABIUint256(tokenId),
		safeTransferDataOffsetWord,
		encode.ABIBytes(data),
	)
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
