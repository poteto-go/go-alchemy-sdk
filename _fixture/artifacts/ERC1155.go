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

// ERC1155MetaData contains all meta data concerning the ERC1155 contract.
var ERC1155MetaData = bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"ApprovalForAll\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"TransferSingle\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"accounts\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"ids\",\"type\":\"uint256[]\"}],\"name\":\"balanceOfBatch\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"isApprovedForAll\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"setApprovalForAll\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"uri\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	ID:  "ERC1155",
	Bin: "0x60e06040526025608081815290610a6e60a0395f9061001e90826100c8565b5034801561002a575f5ffd5b50610182565b634e487b7160e01b5f52604160045260245ffd5b600181811c9082168061005857607f821691505b60208210810361007657634e487b7160e01b5f52602260045260245ffd5b50919050565b601f8211156100c357805f5260205f20601f840160051c810160208510156100a15750805b601f840160051c820191505b818110156100c0575f81556001016100ad565b50505b505050565b81516001600160401b038111156100e1576100e1610030565b6100f5816100ef8454610044565b8461007c565b6020601f821160018114610127575f83156101105750848201515b5f19600385901b1c1916600184901b1784556100c0565b5f84815260208120601f198516915b828110156101565787850151825560209485019460019092019101610136565b508482101561017357868401515f19600387901b60f8161c191681555b50505050600190811b01905550565b6108df8061018f5f395ff3fe608060405234801561000f575f5ffd5b506004361061005f575f3560e01c8062fdd58e146100635780630e89341c14610089578063156e29f6146100a95780634e1273f4146100be578063a22cb465146100de578063e985e9c5146100f1575b5f5ffd5b610076610071366004610559565b61013c565b6040519081526020015b60405180910390f35b61009c610097366004610581565b6101d5565b6040516100809190610598565b6100bc6100b73660046105cd565b610266565b005b6100d16100cc3660046106cf565b610349565b6040516100809190610792565b6100bc6100ec3660046107d4565b610469565b61012c6100ff36600461080d565b6001600160a01b039182165f90815260026020908152604080832093909416825291909152205460ff1690565b6040519015158152602001610080565b5f6001600160a01b0383166101ab5760405162461bcd60e51b815260206004820152602a60248201527f455243313135353a2061646472657373207a65726f206973206e6f742061207660448201526930b634b21037bbb732b960b11b60648201526084015b60405180910390fd5b505f8181526001602090815260408083206001600160a01b03861684529091529020545b92915050565b60605f80546101e39061083e565b80601f016020809104026020016040519081016040528092919081815260200182805461020f9061083e565b801561025a5780601f106102315761010080835404028352916020019161025a565b820191905f5260205f20905b81548152906001019060200180831161023d57829003601f168201915b50505050509050919050565b6001600160a01b0383166102c65760405162461bcd60e51b815260206004820152602160248201527f455243313135353a206d696e7420746f20746865207a65726f206164647265736044820152607360f81b60648201526084016101a2565b5f8281526001602090815260408083206001600160a01b0387168452909152812080548392906102f7908490610876565b909155505060408051838152602081018390526001600160a01b038516915f9133917fc3d58168c5ae7397731d063d5bbf3d657854427343f4c083240f7aacaa2d0f62910160405180910390a4505050565b606081518351146103ae5760405162461bcd60e51b815260206004820152602960248201527f455243313135353a206163636f756e747320616e6420696473206c656e677468604482015268040dad2e6dac2e8c6d60bb1b60648201526084016101a2565b5f835167ffffffffffffffff8111156103c9576103c96105fd565b6040519080825280602002602001820160405280156103f2578160200160208202803683370190505b5090505f5b84518110156104615761043c85828151811061041557610415610895565b602002602001015185838151811061042f5761042f610895565b602002602001015161013c565b82828151811061044e5761044e610895565b60209081029190910101526001016103f7565b509392505050565b336001600160a01b038316036104d35760405162461bcd60e51b815260206004820152602960248201527f455243313135353a2073657474696e6720617070726f76616c20737461747573604482015268103337b91039b2b63360b91b60648201526084016101a2565b335f8181526002602090815260408083206001600160a01b03871680855290835292819020805460ff191686151590811790915590519081529192917f17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31910160405180910390a35050565b80356001600160a01b0381168114610554575f5ffd5b919050565b5f5f6040838503121561056a575f5ffd5b6105738361053e565b946020939093013593505050565b5f60208284031215610591575f5ffd5b5035919050565b602081525f82518060208401528060208501604085015e5f604082850101526040601f19601f83011684010191505092915050565b5f5f5f606084860312156105df575f5ffd5b6105e88461053e565b95602085013595506040909401359392505050565b634e487b7160e01b5f52604160045260245ffd5b604051601f8201601f1916810167ffffffffffffffff8111828210171561063a5761063a6105fd565b604052919050565b5f67ffffffffffffffff82111561065b5761065b6105fd565b5060051b60200190565b5f82601f830112610674575f5ffd5b813561068761068282610642565b610611565b8082825260208201915060208360051b8601019250858311156106a8575f5ffd5b602085015b838110156106c55780358352602092830192016106ad565b5095945050505050565b5f5f604083850312156106e0575f5ffd5b823567ffffffffffffffff8111156106f6575f5ffd5b8301601f81018513610706575f5ffd5b803561071461068282610642565b8082825260208201915060208360051b850101925087831115610735575f5ffd5b6020840193505b8284101561075e5761074d8461053e565b82526020938401939091019061073c565b9450505050602083013567ffffffffffffffff81111561077c575f5ffd5b61078885828601610665565b9150509250929050565b602080825282518282018190525f918401906040840190835b818110156107c95783518352602093840193909201916001016107ab565b509095945050505050565b5f5f604083850312156107e5575f5ffd5b6107ee8361053e565b915060208301358015158114610802575f5ffd5b809150509250929050565b5f5f6040838503121561081e575f5ffd5b6108278361053e565b91506108356020840161053e565b90509250929050565b600181811c9082168061085257607f821691505b60208210810361087057634e487b7160e01b5f52602260045260245ffd5b50919050565b808201808211156101cf57634e487b7160e01b5f52601160045260245ffd5b634e487b7160e01b5f52603260045260245ffdfea264697066735822122015f01b97ac18380c4cdb098c8e515d1209fb2c2e80daef91ab7feafa9ab6bdf864736f6c634300081e003368747470733a2f2f6578616d706c652e636f6d2f657263313135352f7b69647d2e6a736f6e",
}

// ERC1155 is an auto generated Go binding around an Ethereum contract.
type ERC1155 struct {
	abi abi.ABI
}

// NewERC1155 creates a new instance of ERC1155.
func NewERC1155() *ERC1155 {
	parsed, err := ERC1155MetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &ERC1155{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *ERC1155) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackBalanceOf is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x00fdd58e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function balanceOf(address account, uint256 id) view returns(uint256)
func (eRC1155 *ERC1155) PackBalanceOf(account common.Address, id *big.Int) []byte {
	enc, err := eRC1155.abi.Pack("balanceOf", account, id)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackBalanceOf is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x00fdd58e.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function balanceOf(address account, uint256 id) view returns(uint256)
func (eRC1155 *ERC1155) TryPackBalanceOf(account common.Address, id *big.Int) ([]byte, error) {
	return eRC1155.abi.Pack("balanceOf", account, id)
}

// UnpackBalanceOf is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x00fdd58e.
//
// Solidity: function balanceOf(address account, uint256 id) view returns(uint256)
func (eRC1155 *ERC1155) UnpackBalanceOf(data []byte) (*big.Int, error) {
	out, err := eRC1155.abi.Unpack("balanceOf", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackBalanceOfBatch is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4e1273f4.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function balanceOfBatch(address[] accounts, uint256[] ids) view returns(uint256[])
func (eRC1155 *ERC1155) PackBalanceOfBatch(accounts []common.Address, ids []*big.Int) []byte {
	enc, err := eRC1155.abi.Pack("balanceOfBatch", accounts, ids)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackBalanceOfBatch is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4e1273f4.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function balanceOfBatch(address[] accounts, uint256[] ids) view returns(uint256[])
func (eRC1155 *ERC1155) TryPackBalanceOfBatch(accounts []common.Address, ids []*big.Int) ([]byte, error) {
	return eRC1155.abi.Pack("balanceOfBatch", accounts, ids)
}

// UnpackBalanceOfBatch is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x4e1273f4.
//
// Solidity: function balanceOfBatch(address[] accounts, uint256[] ids) view returns(uint256[])
func (eRC1155 *ERC1155) UnpackBalanceOfBatch(data []byte) ([]*big.Int, error) {
	out, err := eRC1155.abi.Unpack("balanceOfBatch", data)
	if err != nil {
		return *new([]*big.Int), err
	}
	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)
	return out0, nil
}

// PackIsApprovedForAll is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe985e9c5.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function isApprovedForAll(address account, address operator) view returns(bool)
func (eRC1155 *ERC1155) PackIsApprovedForAll(account common.Address, operator common.Address) []byte {
	enc, err := eRC1155.abi.Pack("isApprovedForAll", account, operator)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackIsApprovedForAll is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe985e9c5.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function isApprovedForAll(address account, address operator) view returns(bool)
func (eRC1155 *ERC1155) TryPackIsApprovedForAll(account common.Address, operator common.Address) ([]byte, error) {
	return eRC1155.abi.Pack("isApprovedForAll", account, operator)
}

// UnpackIsApprovedForAll is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address account, address operator) view returns(bool)
func (eRC1155 *ERC1155) UnpackIsApprovedForAll(data []byte) (bool, error) {
	out, err := eRC1155.abi.Unpack("isApprovedForAll", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackMint is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x156e29f6.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function mint(address to, uint256 id, uint256 amount) returns()
func (eRC1155 *ERC1155) PackMint(to common.Address, id *big.Int, amount *big.Int) []byte {
	enc, err := eRC1155.abi.Pack("mint", to, id, amount)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMint is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x156e29f6.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function mint(address to, uint256 id, uint256 amount) returns()
func (eRC1155 *ERC1155) TryPackMint(to common.Address, id *big.Int, amount *big.Int) ([]byte, error) {
	return eRC1155.abi.Pack("mint", to, id, amount)
}

// PackSetApprovalForAll is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa22cb465.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (eRC1155 *ERC1155) PackSetApprovalForAll(operator common.Address, approved bool) []byte {
	enc, err := eRC1155.abi.Pack("setApprovalForAll", operator, approved)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSetApprovalForAll is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa22cb465.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (eRC1155 *ERC1155) TryPackSetApprovalForAll(operator common.Address, approved bool) ([]byte, error) {
	return eRC1155.abi.Pack("setApprovalForAll", operator, approved)
}

// PackUri is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0e89341c.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function uri(uint256 ) view returns(string)
func (eRC1155 *ERC1155) PackUri(arg0 *big.Int) []byte {
	enc, err := eRC1155.abi.Pack("uri", arg0)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackUri is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0e89341c.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function uri(uint256 ) view returns(string)
func (eRC1155 *ERC1155) TryPackUri(arg0 *big.Int) ([]byte, error) {
	return eRC1155.abi.Pack("uri", arg0)
}

// UnpackUri is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x0e89341c.
//
// Solidity: function uri(uint256 ) view returns(string)
func (eRC1155 *ERC1155) UnpackUri(data []byte) (string, error) {
	out, err := eRC1155.abi.Unpack("uri", data)
	if err != nil {
		return *new(string), err
	}
	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, nil
}

// ERC1155ApprovalForAll represents a ApprovalForAll event raised by the ERC1155 contract.
type ERC1155ApprovalForAll struct {
	Account  common.Address
	Operator common.Address
	Approved bool
	Raw      *types.Log // Blockchain specific contextual infos
}

const ERC1155ApprovalForAllEventName = "ApprovalForAll"

// ContractEventName returns the user-defined event name.
func (ERC1155ApprovalForAll) ContractEventName() string {
	return ERC1155ApprovalForAllEventName
}

// UnpackApprovalForAllEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event ApprovalForAll(address indexed account, address indexed operator, bool approved)
func (eRC1155 *ERC1155) UnpackApprovalForAllEvent(log *types.Log) (*ERC1155ApprovalForAll, error) {
	event := "ApprovalForAll"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != eRC1155.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(ERC1155ApprovalForAll)
	if len(log.Data) > 0 {
		if err := eRC1155.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range eRC1155.abi.Events[event].Inputs {
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

// ERC1155TransferSingle represents a TransferSingle event raised by the ERC1155 contract.
type ERC1155TransferSingle struct {
	Operator common.Address
	From     common.Address
	To       common.Address
	Id       *big.Int
	Value    *big.Int
	Raw      *types.Log // Blockchain specific contextual infos
}

const ERC1155TransferSingleEventName = "TransferSingle"

// ContractEventName returns the user-defined event name.
func (ERC1155TransferSingle) ContractEventName() string {
	return ERC1155TransferSingleEventName
}

// UnpackTransferSingleEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event TransferSingle(address indexed operator, address indexed from, address indexed to, uint256 id, uint256 value)
func (eRC1155 *ERC1155) UnpackTransferSingleEvent(log *types.Log) (*ERC1155TransferSingle, error) {
	event := "TransferSingle"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != eRC1155.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(ERC1155TransferSingle)
	if len(log.Data) > 0 {
		if err := eRC1155.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range eRC1155.abi.Events[event].Inputs {
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
