package namespace

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/poteto-go/go-alchemy-sdk/constant"
	"github.com/poteto-go/go-alchemy-sdk/types"
)

type IEIP2612 interface {
	// Nonces returns the current nonce for the given owner address.
	Nonces(contractAddress, ownerAddress string) (*big.Int, error)

	// DomainSeparator returns the EIP-712 domain separator for the contract.
	DomainSeparator(contractAddress string) ([32]byte, error)
}

type eip2612 struct {
	ether types.EtherApi
}

func NewEIP2612Namespace(ether types.EtherApi) IEIP2612 {
	return &eip2612{ether: ether}
}

func (e *eip2612) Nonces(contractAddress, ownerAddress string) (*big.Int, error) {
	output, err := e.ether.CallReadMethod(
		constant.NoncesFnSignature,
		contractAddress,
		common.LeftPadBytes(common.HexToAddress(ownerAddress).Bytes(), constant.ABIWordSize),
	)
	if err != nil {
		return nil, err
	}
	return new(big.Int).SetBytes(output), nil
}

func (e *eip2612) DomainSeparator(contractAddress string) ([32]byte, error) {
	output, err := e.ether.CallReadMethod(
		constant.DomainSeparatorFnSignature,
		contractAddress,
	)
	if err != nil {
		return [32]byte{}, err
	}
	if len(output) < 32 {
		return [32]byte{}, fmt.Errorf("unexpected output length: %d", len(output))
	}
	var result [32]byte
	copy(result[:], output[:32])
	return result, nil
}
