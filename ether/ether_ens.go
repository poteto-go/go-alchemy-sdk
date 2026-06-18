package ether

import (
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/poteto-go/go-alchemy-sdk/constant"
	"github.com/poteto-go/go-alchemy-sdk/decode"
	"github.com/poteto-go/go-alchemy-sdk/validate"
)

// ensResolverFor calls registryAddr to get the resolver address for node.
// Returns ErrENSResolverNotFound if the registry returns the zero address.
func (ether *Ether) ensResolverFor(registryAddr common.Address, node [32]byte) (common.Address, error) {
	out, err := ether.CallReadMethod(constant.ENSResolverFnSignature, registryAddr.Hex(), node[:])
	if err != nil {
		return common.Address{}, err
	}
	resolver, err := decode.ABIAddress(out)
	if err != nil {
		return common.Address{}, err
	}
	if resolver == (common.Address{}) {
		return common.Address{}, constant.ErrENSResolverNotFound
	}
	return resolver, nil
}

// ResolveNameBy resolves an ENS name to a lowercase hex address using the
// provided ENS registry contract address.
// If name is already a valid hex address it is returned as-is (lowercased).
func (ether *Ether) ResolveNameBy(registryAddress string, name string) (string, error) {
	if err := validate.Address(registryAddress); err != nil {
		return "", err
	}
	if common.IsHexAddress(name) {
		return strings.ToLower(name), nil
	}

	node := decode.ENSNamehash(name)
	registry := common.HexToAddress(registryAddress)

	resolver, err := ether.ensResolverFor(registry, node)
	if err != nil {
		return "", err
	}

	out, err := ether.CallReadMethod(constant.ENSAddrFnSignature, resolver.Hex(), node[:])
	if err != nil {
		return "", err
	}
	addr, err := decode.ABIAddress(out)
	if err != nil {
		return "", err
	}
	return strings.ToLower(addr.Hex()), nil
}

// LookupAddressBy performs a reverse ENS lookup (address → name) using the
// provided ENS registry contract address.
// Returns ErrENSNameNotFound when no reverse record is registered.
func (ether *Ether) LookupAddressBy(registryAddress string, address string) (string, error) {
	if err := validate.Address(registryAddress); err != nil {
		return "", err
	}
	if err := validate.Address(address); err != nil {
		return "", err
	}

	lowered := strings.TrimPrefix(strings.ToLower(address), "0x")
	reverseNode := decode.ENSNamehash(lowered + ".addr.reverse")
	registry := common.HexToAddress(registryAddress)

	resolver, err := ether.ensResolverFor(registry, reverseNode)
	if err != nil {
		return "", err
	}

	out, err := ether.CallReadMethod(constant.ENSNameFnSignature, resolver.Hex(), reverseNode[:])
	if err != nil {
		return "", err
	}
	name, err := decode.ABIString(out)
	if err != nil {
		return "", err
	}
	if name == "" {
		return "", constant.ErrENSNameNotFound
	}
	return name, nil
}
