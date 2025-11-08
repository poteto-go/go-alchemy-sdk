// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package artifacts

import (
	"bytes"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = bytes.Equal
	_ = errors.New
	_ = big.NewInt
	_ = common.Big1
	_ = types.BloomLookup
	_ = abi.ConvertType
)

// PotetoStorageMetaData contains all meta data concerning the PotetoStorage contract.
var PotetoStorageMetaData = bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"retrieve\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"number\",\"type\":\"uint256\"}],\"name\":\"store\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	ID:  "9dd2cd8c45e1186b82c402350e1b9df1e6",
	Bin: "0x6080604052348015600e575f5ffd5b506101298061001c5f395ff3fe6080604052348015600e575f5ffd5b50600436106030575f3560e01c80632e64cec11460345780636057361d14604e575b5f5ffd5b603a6066565b60405160459190608d565b60405180910390f35b606460048036038101906060919060cd565b606e565b005b5f5f54905090565b805f8190555050565b5f819050919050565b6087816077565b82525050565b5f602082019050609e5f8301846080565b92915050565b5f5ffd5b60af816077565b811460b8575f5ffd5b50565b5f8135905060c78160a8565b92915050565b5f6020828403121560df5760de60a4565b5b5f60ea8482850160bb565b9150509291505056fea2646970667358221220526ad04d169847529a898a88dd9dd46c6b4f07e916d1f44d4528760afc07be3864736f6c634300081e0033",
}

// PotetoStorage is an auto generated Go binding around an Ethereum contract.
type PotetoStorage struct {
	abi abi.ABI
}

// NewPotetoStorage creates a new instance of PotetoStorage.
func NewPotetoStorage() *PotetoStorage {
	parsed, err := PotetoStorageMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &PotetoStorage{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *PotetoStorage) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackRetrieve is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2e64cec1.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function retrieve() view returns(uint256)
func (potetoStorage *PotetoStorage) PackRetrieve() []byte {
	enc, err := potetoStorage.abi.Pack("retrieve")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRetrieve is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2e64cec1.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function retrieve() view returns(uint256)
func (potetoStorage *PotetoStorage) TryPackRetrieve() ([]byte, error) {
	return potetoStorage.abi.Pack("retrieve")
}

// UnpackRetrieve is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x2e64cec1.
//
// Solidity: function retrieve() view returns(uint256)
func (potetoStorage *PotetoStorage) UnpackRetrieve(data []byte) (*big.Int, error) {
	out, err := potetoStorage.abi.Unpack("retrieve", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackStore is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6057361d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function store(uint256 number) returns()
func (potetoStorage *PotetoStorage) PackStore(number *big.Int) []byte {
	enc, err := potetoStorage.abi.Pack("store", number)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackStore is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6057361d.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function store(uint256 number) returns()
func (potetoStorage *PotetoStorage) TryPackStore(number *big.Int) ([]byte, error) {
	return potetoStorage.abi.Pack("store", number)
}
