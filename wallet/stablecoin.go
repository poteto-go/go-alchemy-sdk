package wallet

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/poteto-go/go-alchemy-sdk/constant"
)

type walletStableCoin struct {
	walletERC20
}

func (api *walletStableCoin) MintNoWait(contractAddress, toAddress string, amount *big.Int, gasLimit *uint64) (common.Hash, error) {
	if err := validateAddress(toAddress); err != nil {
		return common.Hash{}, err
	}
	if err := validateUint256(amount); err != nil {
		return common.Hash{}, err
	}
	return api.sendERC20Tx(contractAddress, gasLimit, constant.MintFnSignature,
		common.LeftPadBytes(common.HexToAddress(toAddress).Bytes(), constant.ABIWordSize),
		common.LeftPadBytes(amount.Bytes(), constant.ABIWordSize),
	)
}

func (api *walletStableCoin) Mint(ctx context.Context, contractAddress, toAddress string, amount *big.Int, gasLimit *uint64) (*gethTypes.Receipt, error) {
	return api.waitMined(ctx, func() (common.Hash, error) {
		return api.MintNoWait(contractAddress, toAddress, amount, gasLimit)
	})
}

func (api *walletStableCoin) BurnNoWait(contractAddress string, amount *big.Int, gasLimit *uint64) (common.Hash, error) {
	if err := validateUint256(amount); err != nil {
		return common.Hash{}, err
	}
	return api.sendERC20Tx(contractAddress, gasLimit, constant.BurnFnSignature,
		common.LeftPadBytes(amount.Bytes(), constant.ABIWordSize),
	)
}

func (api *walletStableCoin) Burn(ctx context.Context, contractAddress string, amount *big.Int, gasLimit *uint64) (*gethTypes.Receipt, error) {
	return api.waitMined(ctx, func() (common.Hash, error) {
		return api.BurnNoWait(contractAddress, amount, gasLimit)
	})
}

func (api *walletStableCoin) BlacklistNoWait(contractAddress, address string, gasLimit *uint64) (common.Hash, error) {
	if err := validateAddress(address); err != nil {
		return common.Hash{}, err
	}
	return api.sendERC20Tx(contractAddress, gasLimit, constant.BlacklistFnSignature,
		common.LeftPadBytes(common.HexToAddress(address).Bytes(), constant.ABIWordSize),
	)
}

func (api *walletStableCoin) Blacklist(ctx context.Context, contractAddress, address string, gasLimit *uint64) (*gethTypes.Receipt, error) {
	return api.waitMined(ctx, func() (common.Hash, error) {
		return api.BlacklistNoWait(contractAddress, address, gasLimit)
	})
}

func (api *walletStableCoin) UnBlacklistNoWait(contractAddress, address string, gasLimit *uint64) (common.Hash, error) {
	if err := validateAddress(address); err != nil {
		return common.Hash{}, err
	}
	return api.sendERC20Tx(contractAddress, gasLimit, constant.UnBlacklistFnSignature,
		common.LeftPadBytes(common.HexToAddress(address).Bytes(), constant.ABIWordSize),
	)
}

func (api *walletStableCoin) UnBlacklist(ctx context.Context, contractAddress, address string, gasLimit *uint64) (*gethTypes.Receipt, error) {
	return api.waitMined(ctx, func() (common.Hash, error) {
		return api.UnBlacklistNoWait(contractAddress, address, gasLimit)
	})
}

func (api *walletStableCoin) IsBlacklisted(contractAddress, address string) (bool, error) {
	sc := api.w.snapshotStableCoin()
	if sc == nil {
		return false, constant.ErrWalletIsNotConnected
	}
	return sc.IsBlacklisted(contractAddress, address)
}

func (api *walletStableCoin) MasterMinter(contractAddress string) (common.Address, error) {
	sc := api.w.snapshotStableCoin()
	if sc == nil {
		return common.Address{}, constant.ErrWalletIsNotConnected
	}
	return sc.MasterMinter(contractAddress)
}

func (api *walletStableCoin) Pauser(contractAddress string) (common.Address, error) {
	sc := api.w.snapshotStableCoin()
	if sc == nil {
		return common.Address{}, constant.ErrWalletIsNotConnected
	}
	return sc.Pauser(contractAddress)
}

func (api *walletStableCoin) Blacklister(contractAddress string) (common.Address, error) {
	sc := api.w.snapshotStableCoin()
	if sc == nil {
		return common.Address{}, constant.ErrWalletIsNotConnected
	}
	return sc.Blacklister(contractAddress)
}

func (api *walletStableCoin) Owner(contractAddress string) (common.Address, error) {
	sc := api.w.snapshotStableCoin()
	if sc == nil {
		return common.Address{}, constant.ErrWalletIsNotConnected
	}
	return sc.Owner(contractAddress)
}

func (api *walletStableCoin) PauseNoWait(contractAddress string, gasLimit *uint64) (common.Hash, error) {
	return api.sendERC20Tx(contractAddress, gasLimit, constant.PauseFnSignature)
}

func (api *walletStableCoin) Pause(ctx context.Context, contractAddress string, gasLimit *uint64) (*gethTypes.Receipt, error) {
	return api.waitMined(ctx, func() (common.Hash, error) {
		return api.PauseNoWait(contractAddress, gasLimit)
	})
}

func (api *walletStableCoin) UnpauseNoWait(contractAddress string, gasLimit *uint64) (common.Hash, error) {
	return api.sendERC20Tx(contractAddress, gasLimit, constant.UnpauseFnSignature)
}

func (api *walletStableCoin) Unpause(ctx context.Context, contractAddress string, gasLimit *uint64) (*gethTypes.Receipt, error) {
	return api.waitMined(ctx, func() (common.Hash, error) {
		return api.UnpauseNoWait(contractAddress, gasLimit)
	})
}

func (api *walletStableCoin) Paused(contractAddress string) (bool, error) {
	sc := api.w.snapshotStableCoin()
	if sc == nil {
		return false, constant.ErrWalletIsNotConnected
	}
	return sc.Paused(contractAddress)
}

func (api *walletStableCoin) TransferOwnershipNoWait(contractAddress, newOwner string, gasLimit *uint64) (common.Hash, error) {
	if err := validateAddress(newOwner); err != nil {
		return common.Hash{}, err
	}
	return api.sendERC20Tx(contractAddress, gasLimit, constant.TransferOwnershipFnSignature,
		common.LeftPadBytes(common.HexToAddress(newOwner).Bytes(), constant.ABIWordSize),
	)
}

func (api *walletStableCoin) TransferOwnership(ctx context.Context, contractAddress, newOwner string, gasLimit *uint64) (*gethTypes.Receipt, error) {
	return api.waitMined(ctx, func() (common.Hash, error) {
		return api.TransferOwnershipNoWait(contractAddress, newOwner, gasLimit)
	})
}

func (api *walletStableCoin) Currency(contractAddress string) (string, error) {
	sc := api.w.snapshotStableCoin()
	if sc == nil {
		return "", constant.ErrWalletIsNotConnected
	}
	return sc.Currency(contractAddress)
}

func (api *walletStableCoin) Version(contractAddress string) (string, error) {
	sc := api.w.snapshotStableCoin()
	if sc == nil {
		return "", constant.ErrWalletIsNotConnected
	}
	return sc.Version(contractAddress)
}

func (api *walletStableCoin) ConfigureMinterNoWait(contractAddress, minter string, allowance *big.Int, gasLimit *uint64) (common.Hash, error) {
	if err := validateAddress(minter); err != nil {
		return common.Hash{}, err
	}
	if err := validateUint256(allowance); err != nil {
		return common.Hash{}, err
	}
	return api.sendERC20Tx(contractAddress, gasLimit, constant.ConfigureMinterFnSignature,
		common.LeftPadBytes(common.HexToAddress(minter).Bytes(), constant.ABIWordSize),
		common.LeftPadBytes(allowance.Bytes(), constant.ABIWordSize),
	)
}

func (api *walletStableCoin) ConfigureMinter(ctx context.Context, contractAddress, minter string, allowance *big.Int, gasLimit *uint64) (*gethTypes.Receipt, error) {
	return api.waitMined(ctx, func() (common.Hash, error) {
		return api.ConfigureMinterNoWait(contractAddress, minter, allowance, gasLimit)
	})
}

func (api *walletStableCoin) RemoveMinterNoWait(contractAddress, minter string, gasLimit *uint64) (common.Hash, error) {
	if err := validateAddress(minter); err != nil {
		return common.Hash{}, err
	}
	return api.sendERC20Tx(contractAddress, gasLimit, constant.RemoveMinterFnSignature,
		common.LeftPadBytes(common.HexToAddress(minter).Bytes(), constant.ABIWordSize),
	)
}

func (api *walletStableCoin) RemoveMinter(ctx context.Context, contractAddress, minter string, gasLimit *uint64) (*gethTypes.Receipt, error) {
	return api.waitMined(ctx, func() (common.Hash, error) {
		return api.RemoveMinterNoWait(contractAddress, minter, gasLimit)
	})
}

func (api *walletStableCoin) IsMinter(contractAddress, address string) (bool, error) {
	sc := api.w.snapshotStableCoin()
	if sc == nil {
		return false, constant.ErrWalletIsNotConnected
	}
	return sc.IsMinter(contractAddress, address)
}

func (api *walletStableCoin) MinterAllowance(contractAddress, address string) (*big.Int, error) {
	sc := api.w.snapshotStableCoin()
	if sc == nil {
		return nil, constant.ErrWalletIsNotConnected
	}
	return sc.MinterAllowance(contractAddress, address)
}

func (api *walletStableCoin) UpdateMasterMinterNoWait(contractAddress, newMasterMinter string, gasLimit *uint64) (common.Hash, error) {
	if err := validateAddress(newMasterMinter); err != nil {
		return common.Hash{}, err
	}
	return api.sendERC20Tx(contractAddress, gasLimit, constant.UpdateMasterMinterFnSignature,
		common.LeftPadBytes(common.HexToAddress(newMasterMinter).Bytes(), constant.ABIWordSize),
	)
}

func (api *walletStableCoin) UpdateMasterMinter(ctx context.Context, contractAddress, newMasterMinter string, gasLimit *uint64) (*gethTypes.Receipt, error) {
	return api.waitMined(ctx, func() (common.Hash, error) {
		return api.UpdateMasterMinterNoWait(contractAddress, newMasterMinter, gasLimit)
	})
}

func (api *walletStableCoin) UpdateBlacklisterNoWait(contractAddress, newBlacklister string, gasLimit *uint64) (common.Hash, error) {
	if err := validateAddress(newBlacklister); err != nil {
		return common.Hash{}, err
	}
	return api.sendERC20Tx(contractAddress, gasLimit, constant.UpdateBlacklisterFnSignature,
		common.LeftPadBytes(common.HexToAddress(newBlacklister).Bytes(), constant.ABIWordSize),
	)
}

func (api *walletStableCoin) UpdateBlacklister(ctx context.Context, contractAddress, newBlacklister string, gasLimit *uint64) (*gethTypes.Receipt, error) {
	return api.waitMined(ctx, func() (common.Hash, error) {
		return api.UpdateBlacklisterNoWait(contractAddress, newBlacklister, gasLimit)
	})
}

func (api *walletStableCoin) UpdatePauserNoWait(contractAddress, newPauser string, gasLimit *uint64) (common.Hash, error) {
	if err := validateAddress(newPauser); err != nil {
		return common.Hash{}, err
	}
	return api.sendERC20Tx(contractAddress, gasLimit, constant.UpdatePauserFnSignature,
		common.LeftPadBytes(common.HexToAddress(newPauser).Bytes(), constant.ABIWordSize),
	)
}

func (api *walletStableCoin) UpdatePauser(ctx context.Context, contractAddress, newPauser string, gasLimit *uint64) (*gethTypes.Receipt, error) {
	return api.waitMined(ctx, func() (common.Hash, error) {
		return api.UpdatePauserNoWait(contractAddress, newPauser, gasLimit)
	})
}

var permitTypeHash = crypto.Keccak256([]byte("Permit(address owner,address spender,uint256 value,uint256 nonce,uint256 deadline)"))

func (api *walletStableCoin) signPermit(contractAddress, ownerAddress, spenderAddress string, value, deadline *big.Int) (uint8, [32]byte, [32]byte, error) {
	sc := api.w.snapshotStableCoin()
	if sc == nil {
		return 0, [32]byte{}, [32]byte{}, constant.ErrWalletIsNotConnected
	}

	nonce, err := sc.Nonces(contractAddress, ownerAddress)
	if err != nil {
		return 0, [32]byte{}, [32]byte{}, err
	}

	domainSeparator, err := sc.DomainSeparator(contractAddress)
	if err != nil {
		return 0, [32]byte{}, [32]byte{}, err
	}

	// ABI-encode: permitTypeHash || owner || spender || value || nonce || deadline
	encoded := make([]byte, 0, constant.ABIWordSize*6)
	encoded = append(encoded, permitTypeHash...)
	encoded = append(encoded, common.LeftPadBytes(common.HexToAddress(ownerAddress).Bytes(), constant.ABIWordSize)...)
	encoded = append(encoded, common.LeftPadBytes(common.HexToAddress(spenderAddress).Bytes(), constant.ABIWordSize)...)
	encoded = append(encoded, common.LeftPadBytes(value.Bytes(), constant.ABIWordSize)...)
	encoded = append(encoded, common.LeftPadBytes(nonce.Bytes(), constant.ABIWordSize)...)
	encoded = append(encoded, common.LeftPadBytes(deadline.Bytes(), constant.ABIWordSize)...)

	return api.signEIP712(domainSeparator, encoded)
}

func (api *walletStableCoin) PermitNoWait(contractAddress, ownerAddress, spenderAddress string, value, deadline *big.Int, gasLimit *uint64) (common.Hash, error) {
	if err := validateAddress(ownerAddress); err != nil {
		return common.Hash{}, err
	}
	if err := validateAddress(spenderAddress); err != nil {
		return common.Hash{}, err
	}
	if err := validateUint256(value); err != nil {
		return common.Hash{}, err
	}
	if err := validateUint256(deadline); err != nil {
		return common.Hash{}, err
	}
	v, r, s, err := api.signPermit(contractAddress, ownerAddress, spenderAddress, value, deadline)
	if err != nil {
		return common.Hash{}, err
	}
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

func (api *walletStableCoin) Permit(ctx context.Context, contractAddress, ownerAddress, spenderAddress string, value, deadline *big.Int, gasLimit *uint64) (*gethTypes.Receipt, error) {
	return api.waitMined(ctx, func() (common.Hash, error) {
		return api.PermitNoWait(contractAddress, ownerAddress, spenderAddress, value, deadline, gasLimit)
	})
}

var transferWithAuthorizationTypeHash = crypto.Keccak256([]byte("TransferWithAuthorization(address from,address to,uint256 value,uint256 validAfter,uint256 validBefore,bytes32 nonce)"))
var receiveWithAuthorizationTypeHash = crypto.Keccak256([]byte("ReceiveWithAuthorization(address from,address to,uint256 value,uint256 validAfter,uint256 validBefore,bytes32 nonce)"))
var cancelAuthorizationTypeHash = crypto.Keccak256([]byte("CancelAuthorization(address authorizer,bytes32 nonce)"))

func (api *walletStableCoin) domainSeparatorFor(contractAddress string) ([32]byte, error) {
	sc := api.w.snapshotStableCoin()
	if sc == nil {
		return [32]byte{}, constant.ErrWalletIsNotConnected
	}
	return sc.DomainSeparator(contractAddress)
}

func (api *walletStableCoin) signTransferAuthorization(typeHash []byte, contractAddress, from, to string, value, validAfter, validBefore *big.Int, nonce [32]byte) (uint8, [32]byte, [32]byte, error) {
	domainSeparator, err := api.domainSeparatorFor(contractAddress)
	if err != nil {
		return 0, [32]byte{}, [32]byte{}, err
	}

	encoded := make([]byte, 0, constant.ABIWordSize*7)
	encoded = append(encoded, typeHash...)
	encoded = append(encoded, common.LeftPadBytes(common.HexToAddress(from).Bytes(), constant.ABIWordSize)...)
	encoded = append(encoded, common.LeftPadBytes(common.HexToAddress(to).Bytes(), constant.ABIWordSize)...)
	encoded = append(encoded, common.LeftPadBytes(value.Bytes(), constant.ABIWordSize)...)
	encoded = append(encoded, common.LeftPadBytes(validAfter.Bytes(), constant.ABIWordSize)...)
	encoded = append(encoded, common.LeftPadBytes(validBefore.Bytes(), constant.ABIWordSize)...)
	encoded = append(encoded, nonce[:]...)

	return api.signEIP712(domainSeparator, encoded)
}

func (api *walletStableCoin) signCancelAuthorization(contractAddress, authorizer string, nonce [32]byte) (uint8, [32]byte, [32]byte, error) {
	domainSeparator, err := api.domainSeparatorFor(contractAddress)
	if err != nil {
		return 0, [32]byte{}, [32]byte{}, err
	}

	encoded := make([]byte, 0, constant.ABIWordSize*3)
	encoded = append(encoded, cancelAuthorizationTypeHash...)
	encoded = append(encoded, common.LeftPadBytes(common.HexToAddress(authorizer).Bytes(), constant.ABIWordSize)...)
	encoded = append(encoded, nonce[:]...)

	return api.signEIP712(domainSeparator, encoded)
}

func (api *walletStableCoin) signEIP712(domainSeparator [32]byte, encoded []byte) (uint8, [32]byte, [32]byte, error) {
	structHash := crypto.Keccak256(encoded)

	msg := make([]byte, 0, 2+constant.ABIWordSize*2)
	msg = append(msg, constant.EIP191DataPrefix, constant.EIP712StructuredDataVersion)
	msg = append(msg, domainSeparator[:]...)
	msg = append(msg, structHash...)
	hash := crypto.Keccak256(msg)

	sig, err := crypto.Sign(hash, api.w.privateKey)
	if err != nil {
		return 0, [32]byte{}, [32]byte{}, err
	}

	var r, s [32]byte
	copy(r[:], sig[:constant.ABIWordSize])
	copy(s[:], sig[constant.ABIWordSize:constant.ABIWordSize*2])
	v := sig[constant.ABIWordSize*2] + constant.ECDSALegacyVOffset

	return v, r, s, nil
}

func (api *walletStableCoin) transferOrReceiveAuthorizationNoWait(typeHash, fnSig []byte, contractAddress, from, to string, value, validAfter, validBefore *big.Int, nonce [32]byte, gasLimit *uint64) (common.Hash, error) {
	if err := validateAddress(from); err != nil {
		return common.Hash{}, err
	}
	if err := validateAddress(to); err != nil {
		return common.Hash{}, err
	}
	if err := validateUint256(value); err != nil {
		return common.Hash{}, err
	}
	if err := validateUint256(validAfter); err != nil {
		return common.Hash{}, err
	}
	if err := validateUint256(validBefore); err != nil {
		return common.Hash{}, err
	}
	v, r, s, err := api.signTransferAuthorization(typeHash, contractAddress, from, to, value, validAfter, validBefore, nonce)
	if err != nil {
		return common.Hash{}, err
	}
	return api.sendERC20Tx(contractAddress, gasLimit, fnSig,
		common.LeftPadBytes(common.HexToAddress(from).Bytes(), constant.ABIWordSize),
		common.LeftPadBytes(common.HexToAddress(to).Bytes(), constant.ABIWordSize),
		common.LeftPadBytes(value.Bytes(), constant.ABIWordSize),
		common.LeftPadBytes(validAfter.Bytes(), constant.ABIWordSize),
		common.LeftPadBytes(validBefore.Bytes(), constant.ABIWordSize),
		nonce[:],
		common.LeftPadBytes([]byte{v}, constant.ABIWordSize),
		r[:],
		s[:],
	)
}

func (api *walletStableCoin) TransferWithAuthorizationNoWait(contractAddress, from, to string, value, validAfter, validBefore *big.Int, nonce [32]byte, gasLimit *uint64) (common.Hash, error) {
	return api.transferOrReceiveAuthorizationNoWait(transferWithAuthorizationTypeHash, constant.TransferWithAuthorizationFnSignature, contractAddress, from, to, value, validAfter, validBefore, nonce, gasLimit)
}

func (api *walletStableCoin) TransferWithAuthorization(ctx context.Context, contractAddress, from, to string, value, validAfter, validBefore *big.Int, nonce [32]byte, gasLimit *uint64) (*gethTypes.Receipt, error) {
	return api.waitMined(ctx, func() (common.Hash, error) {
		return api.TransferWithAuthorizationNoWait(contractAddress, from, to, value, validAfter, validBefore, nonce, gasLimit)
	})
}

func (api *walletStableCoin) ReceiveWithAuthorizationNoWait(contractAddress, from, to string, value, validAfter, validBefore *big.Int, nonce [32]byte, gasLimit *uint64) (common.Hash, error) {
	return api.transferOrReceiveAuthorizationNoWait(receiveWithAuthorizationTypeHash, constant.ReceiveWithAuthorizationFnSignature, contractAddress, from, to, value, validAfter, validBefore, nonce, gasLimit)
}

func (api *walletStableCoin) ReceiveWithAuthorization(ctx context.Context, contractAddress, from, to string, value, validAfter, validBefore *big.Int, nonce [32]byte, gasLimit *uint64) (*gethTypes.Receipt, error) {
	return api.waitMined(ctx, func() (common.Hash, error) {
		return api.ReceiveWithAuthorizationNoWait(contractAddress, from, to, value, validAfter, validBefore, nonce, gasLimit)
	})
}

func (api *walletStableCoin) CancelAuthorizationNoWait(contractAddress, authorizer string, nonce [32]byte, gasLimit *uint64) (common.Hash, error) {
	if err := validateAddress(authorizer); err != nil {
		return common.Hash{}, err
	}
	v, r, s, err := api.signCancelAuthorization(contractAddress, authorizer, nonce)
	if err != nil {
		return common.Hash{}, err
	}
	return api.sendERC20Tx(contractAddress, gasLimit, constant.CancelAuthorizationFnSignature,
		common.LeftPadBytes(common.HexToAddress(authorizer).Bytes(), constant.ABIWordSize),
		nonce[:],
		common.LeftPadBytes([]byte{v}, constant.ABIWordSize),
		r[:],
		s[:],
	)
}

func (api *walletStableCoin) CancelAuthorization(ctx context.Context, contractAddress, authorizer string, nonce [32]byte, gasLimit *uint64) (*gethTypes.Receipt, error) {
	return api.waitMined(ctx, func() (common.Hash, error) {
		return api.CancelAuthorizationNoWait(contractAddress, authorizer, nonce, gasLimit)
	})
}
