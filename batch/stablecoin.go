package batch

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/poteto-go/go-alchemy-sdk/constant"
	"github.com/poteto-go/go-alchemy-sdk/utils"
	"github.com/poteto-go/go-alchemy-sdk/validate"
)

// StableCoinBatch queues stablecoin (FiatToken-style) contract reads onto its
// Batcher. It embeds *ERC20Batch so the inherited ERC-20 reads are callable
// directly (mirroring the StableCoin namespace), and validates the same inputs.
type StableCoinBatch struct {
	*ERC20Batch
}

func (s *StableCoinBatch) IsBlacklisted(contractAddress, address string) *Result[bool] {
	if err := validate.Addresses(contractAddress, address); err != nil {
		return failed[bool](err)
	}
	return AddCall(s.b, contractAddress, constant.IsBlacklistedFnSignature, utils.DecodeBool, utils.EncodeABIAddress(address))
}

func (s *StableCoinBatch) Paused(contractAddress string) *Result[bool] {
	if err := validate.Address(contractAddress); err != nil {
		return failed[bool](err)
	}
	return AddCall(s.b, contractAddress, constant.PausedFnSignature, utils.DecodeBool)
}

func (s *StableCoinBatch) Owner(contractAddress string) *Result[common.Address] {
	if err := validate.Address(contractAddress); err != nil {
		return failed[common.Address](err)
	}
	return AddCall(s.b, contractAddress, constant.OwnerFnSignature, utils.DecodeABIAddress)
}

func (s *StableCoinBatch) MasterMinter(contractAddress string) *Result[common.Address] {
	if err := validate.Address(contractAddress); err != nil {
		return failed[common.Address](err)
	}
	return AddCall(s.b, contractAddress, constant.MasterMinterFnSignature, utils.DecodeABIAddress)
}

func (s *StableCoinBatch) Pauser(contractAddress string) *Result[common.Address] {
	if err := validate.Address(contractAddress); err != nil {
		return failed[common.Address](err)
	}
	return AddCall(s.b, contractAddress, constant.PauserFnSignature, utils.DecodeABIAddress)
}

func (s *StableCoinBatch) Blacklister(contractAddress string) *Result[common.Address] {
	if err := validate.Address(contractAddress); err != nil {
		return failed[common.Address](err)
	}
	return AddCall(s.b, contractAddress, constant.BlacklisterFnSignature, utils.DecodeABIAddress)
}

func (s *StableCoinBatch) Currency(contractAddress string) *Result[string] {
	if err := validate.Address(contractAddress); err != nil {
		return failed[string](err)
	}
	return AddCall(s.b, contractAddress, constant.CurrencyFnSignature, utils.DecodeABIString)
}

func (s *StableCoinBatch) Version(contractAddress string) *Result[string] {
	if err := validate.Address(contractAddress); err != nil {
		return failed[string](err)
	}
	return AddCall(s.b, contractAddress, constant.VersionFnSignature, utils.DecodeABIString)
}

func (s *StableCoinBatch) IsMinter(contractAddress, address string) *Result[bool] {
	if err := validate.Addresses(contractAddress, address); err != nil {
		return failed[bool](err)
	}
	return AddCall(s.b, contractAddress, constant.IsMinterFnSignature, utils.DecodeBool, utils.EncodeABIAddress(address))
}

func (s *StableCoinBatch) MinterAllowance(contractAddress, address string) *Result[*big.Int] {
	if err := validate.Addresses(contractAddress, address); err != nil {
		return failed[*big.Int](err)
	}
	return AddCall(s.b, contractAddress, constant.MinterAllowanceFnSignature, utils.DecodeUint256, utils.EncodeABIAddress(address))
}

func (s *StableCoinBatch) Nonces(contractAddress, ownerAddress string) *Result[*big.Int] {
	if err := validate.Addresses(contractAddress, ownerAddress); err != nil {
		return failed[*big.Int](err)
	}
	return AddCall(s.b, contractAddress, constant.NoncesFnSignature, utils.DecodeUint256, utils.EncodeABIAddress(ownerAddress))
}

func (s *StableCoinBatch) DomainSeparator(contractAddress string) *Result[[32]byte] {
	if err := validate.Address(contractAddress); err != nil {
		return failed[[32]byte](err)
	}
	return AddCall(s.b, contractAddress, constant.DomainSeparatorFnSignature, utils.DecodeBytes32)
}

func (s *StableCoinBatch) AuthorizationState(contractAddress, authorizer string, nonce [32]byte) *Result[bool] {
	if err := validate.Addresses(contractAddress, authorizer); err != nil {
		return failed[bool](err)
	}
	return AddCall(s.b, contractAddress, constant.AuthorizationStateFnSignature, utils.DecodeBool,
		utils.EncodeABIAddress(authorizer), nonce[:])
}
