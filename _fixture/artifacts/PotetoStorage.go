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
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Stored\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"retrieve\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"number\",\"type\":\"uint256\"}],\"name\":\"store\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	ID:  "9dd2cd8c45e1186b82c402350e1b9df1e6",
	Bin: "0x6080604052348015600e575f5ffd5b506101918061001c5f395ff3fe608060405234801561000f575f5ffd5b5060043610610034575f3560e01c80632e64cec1146100385780636057361d14610056575b5f5ffd5b610040610072565b60405161004d91906100e9565b60405180910390f35b610070600480360381019061006b9190610130565b61007a565b005b5f5f54905090565b805f819055503373ffffffffffffffffffffffffffffffffffffffff167febfcf7c0a1b09f6499e519a8d8bb85ce33cd539ec6cbd964e116cd74943ead1a826040516100c691906100e9565b60405180910390a250565b5f819050919050565b6100e3816100d1565b82525050565b5f6020820190506100fc5f8301846100da565b92915050565b5f5ffd5b61010f816100d1565b8114610119575f5ffd5b50565b5f8135905061012a81610106565b92915050565b5f6020828403121561014557610144610102565b5b5f6101528482850161011c565b9150509291505056fea2646970667358221220d066542b4f92e39c3ad04b021821e2f19be6c01ce6734afd178f96a11584626864736f6c634300081e0033",
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

// PotetoStorageStored represents a Stored event raised by the PotetoStorage contract.
type PotetoStorageStored struct {
	Sender common.Address
	Value  *big.Int
	Raw    *types.Log // Blockchain specific contextual infos
}

const PotetoStorageStoredEventName = "Stored"

// ContractEventName returns the user-defined event name.
func (PotetoStorageStored) ContractEventName() string {
	return PotetoStorageStoredEventName
}

// UnpackStoredEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Stored(address indexed sender, uint256 value)
func (potetoStorage *PotetoStorage) UnpackStoredEvent(log *types.Log) (*PotetoStorageStored, error) {
	event := "Stored"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != potetoStorage.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(PotetoStorageStored)
	if len(log.Data) > 0 {
		if err := potetoStorage.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range potetoStorage.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}
