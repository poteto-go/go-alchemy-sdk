package batch

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/poteto-go/go-alchemy-sdk/constant"
	"github.com/poteto-go/go-alchemy-sdk/decode"
	"github.com/poteto-go/go-alchemy-sdk/encode"
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
	return AddCall(s.b, contractAddress, constant.IsBlacklistedFnSignature, decode.Bool, encode.ABIAddress(address))
}

func (s *StableCoinBatch) Paused(contractAddress string) *Result[bool] {
	if err := validate.Address(contractAddress); err != nil {
		return failed[bool](err)
	}
	return AddCall(s.b, contractAddress, constant.PausedFnSignature, decode.Bool)
}

func (s *StableCoinBatch) Owner(contractAddress string) *Result[common.Address] {
	if err := validate.Address(contractAddress); err != nil {
		return failed[common.Address](err)
	}
	return AddCall(s.b, contractAddress, constant.OwnerFnSignature, decode.ABIAddress)
}

func (s *StableCoinBatch) MasterMinter(contractAddress string) *Result[common.Address] {
	if err := validate.Address(contractAddress); err != nil {
		return failed[common.Address](err)
	}
	return AddCall(s.b, contractAddress, constant.MasterMinterFnSignature, decode.ABIAddress)
}

func (s *StableCoinBatch) Pauser(contractAddress string) *Result[common.Address] {
	if err := validate.Address(contractAddress); err != nil {
		return failed[common.Address](err)
	}
	return AddCall(s.b, contractAddress, constant.PauserFnSignature, decode.ABIAddress)
}

func (s *StableCoinBatch) Blacklister(contractAddress string) *Result[common.Address] {
	if err := validate.Address(contractAddress); err != nil {
		return failed[common.Address](err)
	}
	return AddCall(s.b, contractAddress, constant.BlacklisterFnSignature, decode.ABIAddress)
}

func (s *StableCoinBatch) Currency(contractAddress string) *Result[string] {
	if err := validate.Address(contractAddress); err != nil {
		return failed[string](err)
	}
	return AddCall(s.b, contractAddress, constant.CurrencyFnSignature, decode.ABIString)
}

func (s *StableCoinBatch) Version(contractAddress string) *Result[string] {
	if err := validate.Address(contractAddress); err != nil {
		return failed[string](err)
	}
	return AddCall(s.b, contractAddress, constant.VersionFnSignature, decode.ABIString)
}

func (s *StableCoinBatch) IsMinter(contractAddress, address string) *Result[bool] {
	if err := validate.Addresses(contractAddress, address); err != nil {
		return failed[bool](err)
	}
	return AddCall(s.b, contractAddress, constant.IsMinterFnSignature, decode.Bool, encode.ABIAddress(address))
}

func (s *StableCoinBatch) MinterAllowance(contractAddress, address string) *Result[*big.Int] {
	if err := validate.Addresses(contractAddress, address); err != nil {
		return failed[*big.Int](err)
	}
	return AddCall(s.b, contractAddress, constant.MinterAllowanceFnSignature, decode.Uint256, encode.ABIAddress(address))
}

func (s *StableCoinBatch) Nonces(contractAddress, ownerAddress string) *Result[*big.Int] {
	if err := validate.Addresses(contractAddress, ownerAddress); err != nil {
		return failed[*big.Int](err)
	}
	return AddCall(s.b, contractAddress, constant.NoncesFnSignature, decode.Uint256, encode.ABIAddress(ownerAddress))
}

func (s *StableCoinBatch) DomainSeparator(contractAddress string) *Result[[32]byte] {
	if err := validate.Address(contractAddress); err != nil {
		return failed[[32]byte](err)
	}
	return AddCall(s.b, contractAddress, constant.DomainSeparatorFnSignature, decode.Bytes32)
}

func (s *StableCoinBatch) AuthorizationState(contractAddress, authorizer string, nonce [32]byte) *Result[bool] {
	if err := validate.Addresses(contractAddress, authorizer); err != nil {
		return failed[bool](err)
	}
	return AddCall(s.b, contractAddress, constant.AuthorizationStateFnSignature, decode.Bool,
		encode.ABIAddress(authorizer), nonce[:])
}
