package namespace

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/poteto-go/go-alchemy-sdk/constant"
	"github.com/poteto-go/go-alchemy-sdk/types"
	"github.com/poteto-go/go-alchemy-sdk/utils"
	"github.com/poteto-go/go-alchemy-sdk/validate"
)

type IStableCoin interface {
	IERC20

	// IsBlacklisted returns true if the address is blacklisted on the contract.
	IsBlacklisted(contractAddress, address string) (bool, error)

	// Paused returns the current pause state of the contract.
	Paused(contractAddress string) (bool, error)

	// Owner returns the current owner address of the contract.
	Owner(contractAddress string) (common.Address, error)

	// MasterMinter returns the master minter address of the contract.
	MasterMinter(contractAddress string) (common.Address, error)

	// Pauser returns the pauser address of the contract.
	Pauser(contractAddress string) (common.Address, error)

	// Blacklister returns the blacklister address of the contract.
	Blacklister(contractAddress string) (common.Address, error)

	// Currency returns the currency identifier of the token (e.g. "USD").
	Currency(contractAddress string) (string, error)

	// Version returns the contract version string.
	Version(contractAddress string) (string, error)

	// IsMinter returns true if the address is a configured minter.
	IsMinter(contractAddress, address string) (bool, error)

	// MinterAllowance returns the remaining mint allowance for the given minter.
	MinterAllowance(contractAddress, address string) (*big.Int, error)

	// Nonces returns the current EIP-2612 permit nonce for the given owner.
	Nonces(contractAddress, ownerAddress string) (*big.Int, error)

	// DomainSeparator returns the EIP-712 domain separator for the contract.
	DomainSeparator(contractAddress string) ([32]byte, error)

	// AuthorizationState returns true if the authorization identified by (authorizer, nonce) has been used or cancelled.
	AuthorizationState(contractAddress, authorizer string, nonce [32]byte) (bool, error)
}

type stableCoin struct {
	*ERC20
}

func NewStableCoinNamespace(ether types.EtherApi) IStableCoin {
	return &stableCoin{ERC20: &ERC20{ether: ether}}
}

func (s *stableCoin) IsBlacklisted(contractAddress, address string) (bool, error) {
	if err := validate.Addresses(contractAddress, address); err != nil {
		return false, err
	}
	output, err := s.ether.CallReadMethod(
		constant.IsBlacklistedFnSignature,
		contractAddress,
		utils.EncodeABIAddress(address),
	)
	if err != nil {
		return false, err
	}
	return utils.DecodeBool(output)
}

func (s *stableCoin) Currency(contractAddress string) (string, error) {
	if err := validate.Address(contractAddress); err != nil {
		return "", err
	}
	output, err := s.ether.CallReadMethod(
		constant.CurrencyFnSignature,
		contractAddress,
	)
	if err != nil {
		return "", err
	}
	return utils.DecodeABIString(output)
}

func (s *stableCoin) Version(contractAddress string) (string, error) {
	if err := validate.Address(contractAddress); err != nil {
		return "", err
	}
	output, err := s.ether.CallReadMethod(
		constant.VersionFnSignature,
		contractAddress,
	)
	if err != nil {
		return "", err
	}
	return utils.DecodeABIString(output)
}

func (s *stableCoin) Paused(contractAddress string) (bool, error) {
	if err := validate.Address(contractAddress); err != nil {
		return false, err
	}
	output, err := s.ether.CallReadMethod(
		constant.PausedFnSignature,
		contractAddress,
	)
	if err != nil {
		return false, err
	}
	return utils.DecodeBool(output)
}

func (s *stableCoin) Owner(contractAddress string) (common.Address, error) {
	if err := validate.Address(contractAddress); err != nil {
		return common.Address{}, err
	}
	output, err := s.ether.CallReadMethod(constant.OwnerFnSignature, contractAddress)
	if err != nil {
		return common.Address{}, err
	}
	return utils.DecodeABIAddress(output)
}

func (s *stableCoin) MasterMinter(contractAddress string) (common.Address, error) {
	if err := validate.Address(contractAddress); err != nil {
		return common.Address{}, err
	}
	output, err := s.ether.CallReadMethod(constant.MasterMinterFnSignature, contractAddress)
	if err != nil {
		return common.Address{}, err
	}
	return utils.DecodeABIAddress(output)
}

func (s *stableCoin) Pauser(contractAddress string) (common.Address, error) {
	if err := validate.Address(contractAddress); err != nil {
		return common.Address{}, err
	}
	output, err := s.ether.CallReadMethod(constant.PauserFnSignature, contractAddress)
	if err != nil {
		return common.Address{}, err
	}
	return utils.DecodeABIAddress(output)
}

func (s *stableCoin) Blacklister(contractAddress string) (common.Address, error) {
	if err := validate.Address(contractAddress); err != nil {
		return common.Address{}, err
	}
	output, err := s.ether.CallReadMethod(constant.BlacklisterFnSignature, contractAddress)
	if err != nil {
		return common.Address{}, err
	}
	return utils.DecodeABIAddress(output)
}

func (s *stableCoin) IsMinter(contractAddress, address string) (bool, error) {
	if err := validate.Addresses(contractAddress, address); err != nil {
		return false, err
	}
	output, err := s.ether.CallReadMethod(
		constant.IsMinterFnSignature,
		contractAddress,
		utils.EncodeABIAddress(address),
	)
	if err != nil {
		return false, err
	}
	return utils.DecodeBool(output)
}

func (s *stableCoin) MinterAllowance(contractAddress, address string) (*big.Int, error) {
	if err := validate.Addresses(contractAddress, address); err != nil {
		return nil, err
	}
	output, err := s.ether.CallReadMethod(
		constant.MinterAllowanceFnSignature,
		contractAddress,
		utils.EncodeABIAddress(address),
	)
	if err != nil {
		return nil, err
	}
	return utils.DecodeUint256(output)
}

func (s *stableCoin) Nonces(contractAddress, ownerAddress string) (*big.Int, error) {
	if err := validate.Addresses(contractAddress, ownerAddress); err != nil {
		return nil, err
	}
	output, err := s.ether.CallReadMethod(
		constant.NoncesFnSignature,
		contractAddress,
		utils.EncodeABIAddress(ownerAddress),
	)
	if err != nil {
		return nil, err
	}
	return utils.DecodeUint256(output)
}

func (s *stableCoin) AuthorizationState(contractAddress, authorizer string, nonce [32]byte) (bool, error) {
	output, err := s.ether.CallReadMethod(
		constant.AuthorizationStateFnSignature,
		contractAddress,
		utils.EncodeABIAddress(authorizer),
		nonce[:],
	)
	if err != nil {
		return false, err
	}
	return utils.DecodeBool(output)
}

func (s *stableCoin) DomainSeparator(contractAddress string) ([32]byte, error) {
	if err := validate.Address(contractAddress); err != nil {
		return [32]byte{}, err
	}
	output, err := s.ether.CallReadMethod(
		constant.DomainSeparatorFnSignature,
		contractAddress,
	)
	if err != nil {
		return [32]byte{}, err
	}
	return utils.DecodeBytes32(output)
}
