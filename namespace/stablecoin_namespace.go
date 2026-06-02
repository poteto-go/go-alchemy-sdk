package namespace

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/poteto-go/go-alchemy-sdk/constant"
	"github.com/poteto-go/go-alchemy-sdk/types"
	"github.com/poteto-go/go-alchemy-sdk/utils"
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
}

type stableCoin struct {
	*ERC20
}

func NewStableCoinNamespace(ether types.EtherApi) IStableCoin {
	return &stableCoin{ERC20: &ERC20{ether: ether}}
}

func decodeBoolOutput(output []byte) bool {
	return len(output) > 0 && output[len(output)-1] == 1
}

func (s *stableCoin) IsBlacklisted(contractAddress, address string) (bool, error) {
	output, err := s.ether.CallReadMethod(
		constant.IsBlacklistedFnSignature,
		contractAddress,
		common.LeftPadBytes(common.HexToAddress(address).Bytes(), 32),
	)
	if err != nil {
		return false, err
	}
	return decodeBoolOutput(output), nil
}

func (s *stableCoin) Currency(contractAddress string) (string, error) {
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
	output, err := s.ether.CallReadMethod(
		constant.PausedFnSignature,
		contractAddress,
	)
	if err != nil {
		return false, err
	}
	return decodeBoolOutput(output), nil
}

func (s *stableCoin) Owner(contractAddress string) (common.Address, error) {
	output, err := s.ether.CallReadMethod(constant.OwnerFnSignature, contractAddress)
	if err != nil {
		return common.Address{}, err
	}
	return utils.DecodeABIAddress(output), nil
}

func (s *stableCoin) MasterMinter(contractAddress string) (common.Address, error) {
	output, err := s.ether.CallReadMethod(constant.MasterMinterFnSignature, contractAddress)
	if err != nil {
		return common.Address{}, err
	}
	return utils.DecodeABIAddress(output), nil
}

func (s *stableCoin) Pauser(contractAddress string) (common.Address, error) {
	output, err := s.ether.CallReadMethod(constant.PauserFnSignature, contractAddress)
	if err != nil {
		return common.Address{}, err
	}
	return utils.DecodeABIAddress(output), nil
}

func (s *stableCoin) Blacklister(contractAddress string) (common.Address, error) {
	output, err := s.ether.CallReadMethod(constant.BlacklisterFnSignature, contractAddress)
	if err != nil {
		return common.Address{}, err
	}
	return utils.DecodeABIAddress(output), nil
}
