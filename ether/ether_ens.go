package ether

import (
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/poteto-go/go-alchemy-sdk/constant"
	"github.com/poteto-go/go-alchemy-sdk/decode"
	"github.com/poteto-go/go-alchemy-sdk/validate"
)

// ensRegistryAddress is the canonical ENS registry deployed at the same
// address on Mainnet, Goerli, and Sepolia.
const ensRegistryAddress = "0x00000000000C2E074eC69A0dFb2997BA6C7d2e1e"

// ensNamehash implements EIP-137: recursively hashes name labels with keccak256.
func ensNamehash(name string) [32]byte {
	node := [32]byte{}
	if name == "" {
		return node
	}
	var buf [64]byte
	labels := strings.Split(name, ".")
	for i := len(labels) - 1; i >= 0; i-- {
		copy(buf[:32], node[:])
		copy(buf[32:], crypto.Keccak256([]byte(labels[i])))
		node = crypto.Keccak256Hash(buf[:])
	}
	return node
}

// ensResolverFor calls the ENS registry to get the resolver address for node.
// Returns ErrENSResolverNotFound if the registry returns the zero address.
func (ether *Ether) ensResolverFor(node [32]byte) (common.Address, error) {
	out, err := ether.CallReadMethod(constant.ENSResolverFnSignature, ensRegistryAddress, node[:])
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

// ResolveName resolves an ENS name to a lowercase hex address.
// If name is already a valid hex address it is returned as-is (lowercased).
func (ether *Ether) ResolveName(name string) (string, error) {
	if common.IsHexAddress(name) {
		return strings.ToLower(name), nil
	}

	node := ensNamehash(name)

	resolver, err := ether.ensResolverFor(node)
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

// LookupAddress performs a reverse ENS lookup (address → name).
// Returns ErrENSNameNotFound when no reverse record is registered.
func (ether *Ether) LookupAddress(address string) (string, error) {
	if err := validate.Address(address); err != nil {
		return "", err
	}

	// Reverse node: namehash("<lowercase_addr_without_0x>.addr.reverse")
	lowered := strings.TrimPrefix(strings.ToLower(address), "0x")
	reverseNode := ensNamehash(lowered + ".addr.reverse")

	resolver, err := ether.ensResolverFor(reverseNode)
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
