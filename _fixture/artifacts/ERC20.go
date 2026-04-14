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

// ArtifactsMetaData contains all meta data concerning the Artifacts contract.
var ArtifactsMetaData = bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_initialSupply\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	ID:  "Artifacts",
}

// Artifacts is an auto generated Go binding around an Ethereum contract.
type Artifacts struct {
	abi abi.ABI
}

// NewArtifacts creates a new instance of Artifacts.
func NewArtifacts() *Artifacts {
	parsed, err := ArtifactsMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &Artifacts{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *Artifacts) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackConstructor is the Go binding used to pack the parameters required for
// contract deployment.
//
// Solidity: constructor(uint256 _initialSupply) returns()
func (artifacts *Artifacts) PackConstructor(_initialSupply *big.Int) []byte {
	enc, err := artifacts.abi.Pack("", _initialSupply)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackAllowance is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdd62ed3e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function allowance(address , address ) view returns(uint256)
func (artifacts *Artifacts) PackAllowance(arg0 common.Address, arg1 common.Address) []byte {
	enc, err := artifacts.abi.Pack("allowance", arg0, arg1)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackAllowance is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdd62ed3e.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function allowance(address , address ) view returns(uint256)
func (artifacts *Artifacts) TryPackAllowance(arg0 common.Address, arg1 common.Address) ([]byte, error) {
	return artifacts.abi.Pack("allowance", arg0, arg1)
}

// UnpackAllowance is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xdd62ed3e.
//
// Solidity: function allowance(address , address ) view returns(uint256)
func (artifacts *Artifacts) UnpackAllowance(data []byte) (*big.Int, error) {
	out, err := artifacts.abi.Unpack("allowance", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackApprove is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x095ea7b3.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (artifacts *Artifacts) PackApprove(spender common.Address, amount *big.Int) []byte {
	enc, err := artifacts.abi.Pack("approve", spender, amount)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackApprove is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x095ea7b3.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (artifacts *Artifacts) TryPackApprove(spender common.Address, amount *big.Int) ([]byte, error) {
	return artifacts.abi.Pack("approve", spender, amount)
}

// UnpackApprove is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (artifacts *Artifacts) UnpackApprove(data []byte) (bool, error) {
	out, err := artifacts.abi.Unpack("approve", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackBalanceOf is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x70a08231.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function balanceOf(address ) view returns(uint256)
func (artifacts *Artifacts) PackBalanceOf(arg0 common.Address) []byte {
	enc, err := artifacts.abi.Pack("balanceOf", arg0)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackBalanceOf is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x70a08231.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function balanceOf(address ) view returns(uint256)
func (artifacts *Artifacts) TryPackBalanceOf(arg0 common.Address) ([]byte, error) {
	return artifacts.abi.Pack("balanceOf", arg0)
}

// UnpackBalanceOf is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x70a08231.
//
// Solidity: function balanceOf(address ) view returns(uint256)
func (artifacts *Artifacts) UnpackBalanceOf(data []byte) (*big.Int, error) {
	out, err := artifacts.abi.Unpack("balanceOf", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackDecimals is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x313ce567.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function decimals() view returns(uint8)
func (artifacts *Artifacts) PackDecimals() []byte {
	enc, err := artifacts.abi.Pack("decimals")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackDecimals is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x313ce567.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function decimals() view returns(uint8)
func (artifacts *Artifacts) TryPackDecimals() ([]byte, error) {
	return artifacts.abi.Pack("decimals")
}

// UnpackDecimals is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (artifacts *Artifacts) UnpackDecimals(data []byte) (uint8, error) {
	out, err := artifacts.abi.Unpack("decimals", data)
	if err != nil {
		return *new(uint8), err
	}
	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)
	return out0, nil
}

// PackName is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x06fdde03.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function name() view returns(string)
func (artifacts *Artifacts) PackName() []byte {
	enc, err := artifacts.abi.Pack("name")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackName is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x06fdde03.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function name() view returns(string)
func (artifacts *Artifacts) TryPackName() ([]byte, error) {
	return artifacts.abi.Pack("name")
}

// UnpackName is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (artifacts *Artifacts) UnpackName(data []byte) (string, error) {
	out, err := artifacts.abi.Unpack("name", data)
	if err != nil {
		return *new(string), err
	}
	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, nil
}

// PackSymbol is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x95d89b41.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function symbol() view returns(string)
func (artifacts *Artifacts) PackSymbol() []byte {
	enc, err := artifacts.abi.Pack("symbol")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSymbol is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x95d89b41.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function symbol() view returns(string)
func (artifacts *Artifacts) TryPackSymbol() ([]byte, error) {
	return artifacts.abi.Pack("symbol")
}

// UnpackSymbol is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (artifacts *Artifacts) UnpackSymbol(data []byte) (string, error) {
	out, err := artifacts.abi.Unpack("symbol", data)
	if err != nil {
		return *new(string), err
	}
	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, nil
}

// PackTotalSupply is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x18160ddd.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function totalSupply() view returns(uint256)
func (artifacts *Artifacts) PackTotalSupply() []byte {
	enc, err := artifacts.abi.Pack("totalSupply")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTotalSupply is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x18160ddd.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function totalSupply() view returns(uint256)
func (artifacts *Artifacts) TryPackTotalSupply() ([]byte, error) {
	return artifacts.abi.Pack("totalSupply")
}

// UnpackTotalSupply is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (artifacts *Artifacts) UnpackTotalSupply(data []byte) (*big.Int, error) {
	out, err := artifacts.abi.Unpack("totalSupply", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackTransfer is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa9059cbb.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function transfer(address to, uint256 amount) returns(bool)
func (artifacts *Artifacts) PackTransfer(to common.Address, amount *big.Int) []byte {
	enc, err := artifacts.abi.Pack("transfer", to, amount)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTransfer is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa9059cbb.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function transfer(address to, uint256 amount) returns(bool)
func (artifacts *Artifacts) TryPackTransfer(to common.Address, amount *big.Int) ([]byte, error) {
	return artifacts.abi.Pack("transfer", to, amount)
}

// UnpackTransfer is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 amount) returns(bool)
func (artifacts *Artifacts) UnpackTransfer(data []byte) (bool, error) {
	out, err := artifacts.abi.Unpack("transfer", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackTransferFrom is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x23b872dd.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function transferFrom(address from, address to, uint256 amount) returns(bool)
func (artifacts *Artifacts) PackTransferFrom(from common.Address, to common.Address, amount *big.Int) []byte {
	enc, err := artifacts.abi.Pack("transferFrom", from, to, amount)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTransferFrom is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x23b872dd.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function transferFrom(address from, address to, uint256 amount) returns(bool)
func (artifacts *Artifacts) TryPackTransferFrom(from common.Address, to common.Address, amount *big.Int) ([]byte, error) {
	return artifacts.abi.Pack("transferFrom", from, to, amount)
}

// UnpackTransferFrom is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 amount) returns(bool)
func (artifacts *Artifacts) UnpackTransferFrom(data []byte) (bool, error) {
	out, err := artifacts.abi.Unpack("transferFrom", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// ArtifactsApproval represents a Approval event raised by the Artifacts contract.
type ArtifactsApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     *types.Log // Blockchain specific contextual infos
}

const ArtifactsApprovalEventName = "Approval"

// ContractEventName returns the user-defined event name.
func (ArtifactsApproval) ContractEventName() string {
	return ArtifactsApprovalEventName
}

// UnpackApprovalEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (artifacts *Artifacts) UnpackApprovalEvent(log *types.Log) (*ArtifactsApproval, error) {
	event := "Approval"
	if len(log.Topics) == 0 || log.Topics[0] != artifacts.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ArtifactsApproval)
	if len(log.Data) > 0 {
		if err := artifacts.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range artifacts.abi.Events[event].Inputs {
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

// ArtifactsTransfer represents a Transfer event raised by the Artifacts contract.
type ArtifactsTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   *types.Log // Blockchain specific contextual infos
}

const ArtifactsTransferEventName = "Transfer"

// ContractEventName returns the user-defined event name.
func (ArtifactsTransfer) ContractEventName() string {
	return ArtifactsTransferEventName
}

// UnpackTransferEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (artifacts *Artifacts) UnpackTransferEvent(log *types.Log) (*ArtifactsTransfer, error) {
	event := "Transfer"
	if len(log.Topics) == 0 || log.Topics[0] != artifacts.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ArtifactsTransfer)
	if len(log.Data) > 0 {
		if err := artifacts.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range artifacts.abi.Events[event].Inputs {
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
