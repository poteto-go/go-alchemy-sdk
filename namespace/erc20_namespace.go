package namespace

import (
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/poteto-go/go-alchemy-sdk/constant"
	"github.com/poteto-go/go-alchemy-sdk/types"
	"golang.org/x/crypto/sha3"
)

type IERC20 interface {
	/*
		BalanceOf returns the balance of the specified address.
	*/
	BalanceOf(
		contractAddress,
		walletAddress string,
	) (*big.Int, error)
}

type ERC20 struct {
	ether types.EtherApi
}

func NewERC20Namespace(ether types.EtherApi) IERC20 {
	return &ERC20{
		ether: ether,
	}
}

func (e *ERC20) BalanceOf(
	contractAddress,
	walletAddress string,
) (*big.Int, error) {
	hash := sha3.NewLegacyKeccak256()
	if _, err := hash.Write(constant.BalanceOfFnSignature); err != nil {
		return nil, err
	}
	methodID := hash.Sum(nil)[:4]

	contractAddr := common.HexToAddress(contractAddress)
	paddedAddress := common.LeftPadBytes(common.HexToAddress(walletAddress).Bytes(), 32)

	data := make([]byte, 0, 36)
	data = append(data, methodID...)
	data = append(data, paddedAddress...)

	msg := ethereum.CallMsg{
		To:   &contractAddr,
		Data: data,
	}

	output, err := e.ether.CallContract(
		msg,
		"latest",
	)
	if err != nil {
		return nil, err
	}

	outputInt := new(big.Int)
	return outputInt.SetBytes(output), nil
}
