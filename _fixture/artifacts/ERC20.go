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

// ERC20MetaData contains all meta data concerning the ERC20 contract.
var ERC20MetaData = bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_initialSupply\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	ID:  "ERC20",
	Bin: "0x60806040526040518060400160405280600d81526020017f4d696e696d616c20546f6b656e000000000000000000000000000000000000008152506003908161004891906103fc565b506040518060400160405280600381526020017f4d544b00000000000000000000000000000000000000000000000000000000008152506004908161008d91906103fc565b50601260055f6101000a81548160ff021916908360ff1602179055503480156100b4575f5ffd5b5060405161108f38038061108f83398181016040528101906100d691906104f9565b6100e633826100ec60201b60201c565b506105ac565b805f5f8282546100fc9190610551565b925050819055508060015f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f82825461014f9190610551565b925050819055508173ffffffffffffffffffffffffffffffffffffffff165f73ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef836040516101b39190610593565b60405180910390a35050565b5f81519050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5f600282049050600182168061023a57607f821691505b60208210810361024d5761024c6101f6565b5b50919050565b5f819050815f5260205f209050919050565b5f6020601f8301049050919050565b5f82821b905092915050565b5f600883026102af7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82610274565b6102b98683610274565b95508019841693508086168417925050509392505050565b5f819050919050565b5f819050919050565b5f6102fd6102f86102f3846102d1565b6102da565b6102d1565b9050919050565b5f819050919050565b610316836102e3565b61032a61032282610304565b848454610280565b825550505050565b5f5f905090565b610341610332565b61034c81848461030d565b505050565b5b8181101561036f576103645f82610339565b600181019050610352565b5050565b601f8211156103b45761038581610253565b61038e84610265565b8101602085101561039d578190505b6103b16103a985610265565b830182610351565b50505b505050565b5f82821c905092915050565b5f6103d45f19846008026103b9565b1980831691505092915050565b5f6103ec83836103c5565b9150826002028217905092915050565b610405826101bf565b67ffffffffffffffff81111561041e5761041d6101c9565b5b6104288254610223565b610433828285610373565b5f60209050601f831160018114610464575f8415610452578287015190505b61045c85826103e1565b8655506104c3565b601f19841661047286610253565b5f5b8281101561049957848901518255600182019150602085019450602081019050610474565b868310156104b657848901516104b2601f8916826103c5565b8355505b6001600288020188555050505b505050505050565b5f5ffd5b6104d8816102d1565b81146104e2575f5ffd5b50565b5f815190506104f3816104cf565b92915050565b5f6020828403121561050e5761050d6104cb565b5b5f61051b848285016104e5565b91505092915050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f61055b826102d1565b9150610566836102d1565b925082820190508082111561057e5761057d610524565b5b92915050565b61058d816102d1565b82525050565b5f6020820190506105a65f830184610584565b92915050565b610ad6806105b95f395ff3fe608060405234801561000f575f5ffd5b5060043610610091575f3560e01c8063313ce56711610064578063313ce5671461013157806370a082311461014f57806395d89b411461017f578063a9059cbb1461019d578063dd62ed3e146101cd57610091565b806306fdde0314610095578063095ea7b3146100b357806318160ddd146100e357806323b872dd14610101575b5f5ffd5b61009d6101fd565b6040516100aa9190610779565b60405180910390f35b6100cd60048036038101906100c8919061082a565b610289565b6040516100da9190610882565b60405180910390f35b6100eb610376565b6040516100f891906108aa565b60405180910390f35b61011b600480360381019061011691906108c3565b61037b565b6040516101289190610882565b60405180910390f35b610139610520565b604051610146919061092e565b60405180910390f35b61016960048036038101906101649190610947565b610532565b60405161017691906108aa565b60405180910390f35b610187610547565b6040516101949190610779565b60405180910390f35b6101b760048036038101906101b2919061082a565b6105d3565b6040516101c49190610882565b60405180910390f35b6101e760048036038101906101e29190610972565b6106e9565b6040516101f491906108aa565b60405180910390f35b6003805461020a906109dd565b80601f0160208091040260200160405190810160405280929190818152602001828054610236906109dd565b80156102815780601f1061025857610100808354040283529160200191610281565b820191905f5260205f20905b81548152906001019060200180831161026457829003601f168201915b505050505081565b5f8160025f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20819055508273ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9258460405161036491906108aa565b60405180910390a36001905092915050565b5f5481565b5f8160025f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8282546104039190610a3a565b925050819055508160015f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8282546104569190610a3a565b925050819055508160015f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8282546104a99190610a6d565b925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef8460405161050d91906108aa565b60405180910390a3600190509392505050565b60055f9054906101000a900460ff1681565b6001602052805f5260405f205f915090505481565b60048054610554906109dd565b80601f0160208091040260200160405190810160405280929190818152602001828054610580906109dd565b80156105cb5780601f106105a2576101008083540402835291602001916105cb565b820191905f5260205f20905b8154815290600101906020018083116105ae57829003601f168201915b505050505081565b5f8160015f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8282546106209190610a3a565b925050819055508160015f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8282546106739190610a6d565b925050819055508273ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef846040516106d791906108aa565b60405180910390a36001905092915050565b6002602052815f5260405f20602052805f5260405f205f91509150505481565b5f81519050919050565b5f82825260208201905092915050565b8281835e5f83830152505050565b5f601f19601f8301169050919050565b5f61074b82610709565b6107558185610713565b9350610765818560208601610723565b61076e81610731565b840191505092915050565b5f6020820190508181035f8301526107918184610741565b905092915050565b5f5ffd5b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f6107c68261079d565b9050919050565b6107d6816107bc565b81146107e0575f5ffd5b50565b5f813590506107f1816107cd565b92915050565b5f819050919050565b610809816107f7565b8114610813575f5ffd5b50565b5f8135905061082481610800565b92915050565b5f5f604083850312156108405761083f610799565b5b5f61084d858286016107e3565b925050602061085e85828601610816565b9150509250929050565b5f8115159050919050565b61087c81610868565b82525050565b5f6020820190506108955f830184610873565b92915050565b6108a4816107f7565b82525050565b5f6020820190506108bd5f83018461089b565b92915050565b5f5f5f606084860312156108da576108d9610799565b5b5f6108e7868287016107e3565b93505060206108f8868287016107e3565b925050604061090986828701610816565b9150509250925092565b5f60ff82169050919050565b61092881610913565b82525050565b5f6020820190506109415f83018461091f565b92915050565b5f6020828403121561095c5761095b610799565b5b5f610969848285016107e3565b91505092915050565b5f5f6040838503121561098857610987610799565b5b5f610995858286016107e3565b92505060206109a6858286016107e3565b9150509250929050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5f60028204905060018216806109f457607f821691505b602082108103610a0757610a066109b0565b5b50919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f610a44826107f7565b9150610a4f836107f7565b9250828203905081811115610a6757610a66610a0d565b5b92915050565b5f610a77826107f7565b9150610a82836107f7565b9250828201905080821115610a9a57610a99610a0d565b5b9291505056fea26469706673582212200162a48ac590fbd60bf0b8dac3487d4badd62d11804f4a310b5d05fa2fb2263e64736f6c634300081e0033",
}

// ERC20 is an auto generated Go binding around an Ethereum contract.
type ERC20 struct {
	abi abi.ABI
}

// NewERC20 creates a new instance of ERC20.
func NewERC20() *ERC20 {
	parsed, err := ERC20MetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &ERC20{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *ERC20) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackConstructor is the Go binding used to pack the parameters required for
// contract deployment.
//
// Solidity: constructor(uint256 _initialSupply) returns()
func (eRC20 *ERC20) PackConstructor(_initialSupply *big.Int) []byte {
	enc, err := eRC20.abi.Pack("", _initialSupply)
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
func (eRC20 *ERC20) PackAllowance(arg0 common.Address, arg1 common.Address) []byte {
	enc, err := eRC20.abi.Pack("allowance", arg0, arg1)
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
func (eRC20 *ERC20) TryPackAllowance(arg0 common.Address, arg1 common.Address) ([]byte, error) {
	return eRC20.abi.Pack("allowance", arg0, arg1)
}

// UnpackAllowance is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xdd62ed3e.
//
// Solidity: function allowance(address , address ) view returns(uint256)
func (eRC20 *ERC20) UnpackAllowance(data []byte) (*big.Int, error) {
	out, err := eRC20.abi.Unpack("allowance", data)
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
func (eRC20 *ERC20) PackApprove(spender common.Address, amount *big.Int) []byte {
	enc, err := eRC20.abi.Pack("approve", spender, amount)
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
func (eRC20 *ERC20) TryPackApprove(spender common.Address, amount *big.Int) ([]byte, error) {
	return eRC20.abi.Pack("approve", spender, amount)
}

// UnpackApprove is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (eRC20 *ERC20) UnpackApprove(data []byte) (bool, error) {
	out, err := eRC20.abi.Unpack("approve", data)
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
func (eRC20 *ERC20) PackBalanceOf(arg0 common.Address) []byte {
	enc, err := eRC20.abi.Pack("balanceOf", arg0)
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
func (eRC20 *ERC20) TryPackBalanceOf(arg0 common.Address) ([]byte, error) {
	return eRC20.abi.Pack("balanceOf", arg0)
}

// UnpackBalanceOf is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x70a08231.
//
// Solidity: function balanceOf(address ) view returns(uint256)
func (eRC20 *ERC20) UnpackBalanceOf(data []byte) (*big.Int, error) {
	out, err := eRC20.abi.Unpack("balanceOf", data)
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
func (eRC20 *ERC20) PackDecimals() []byte {
	enc, err := eRC20.abi.Pack("decimals")
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
func (eRC20 *ERC20) TryPackDecimals() ([]byte, error) {
	return eRC20.abi.Pack("decimals")
}

// UnpackDecimals is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (eRC20 *ERC20) UnpackDecimals(data []byte) (uint8, error) {
	out, err := eRC20.abi.Unpack("decimals", data)
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
func (eRC20 *ERC20) PackName() []byte {
	enc, err := eRC20.abi.Pack("name")
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
func (eRC20 *ERC20) TryPackName() ([]byte, error) {
	return eRC20.abi.Pack("name")
}

// UnpackName is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (eRC20 *ERC20) UnpackName(data []byte) (string, error) {
	out, err := eRC20.abi.Unpack("name", data)
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
func (eRC20 *ERC20) PackSymbol() []byte {
	enc, err := eRC20.abi.Pack("symbol")
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
func (eRC20 *ERC20) TryPackSymbol() ([]byte, error) {
	return eRC20.abi.Pack("symbol")
}

// UnpackSymbol is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (eRC20 *ERC20) UnpackSymbol(data []byte) (string, error) {
	out, err := eRC20.abi.Unpack("symbol", data)
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
func (eRC20 *ERC20) PackTotalSupply() []byte {
	enc, err := eRC20.abi.Pack("totalSupply")
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
func (eRC20 *ERC20) TryPackTotalSupply() ([]byte, error) {
	return eRC20.abi.Pack("totalSupply")
}

// UnpackTotalSupply is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (eRC20 *ERC20) UnpackTotalSupply(data []byte) (*big.Int, error) {
	out, err := eRC20.abi.Unpack("totalSupply", data)
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
func (eRC20 *ERC20) PackTransfer(to common.Address, amount *big.Int) []byte {
	enc, err := eRC20.abi.Pack("transfer", to, amount)
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
func (eRC20 *ERC20) TryPackTransfer(to common.Address, amount *big.Int) ([]byte, error) {
	return eRC20.abi.Pack("transfer", to, amount)
}

// UnpackTransfer is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 amount) returns(bool)
func (eRC20 *ERC20) UnpackTransfer(data []byte) (bool, error) {
	out, err := eRC20.abi.Unpack("transfer", data)
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
func (eRC20 *ERC20) PackTransferFrom(from common.Address, to common.Address, amount *big.Int) []byte {
	enc, err := eRC20.abi.Pack("transferFrom", from, to, amount)
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
func (eRC20 *ERC20) TryPackTransferFrom(from common.Address, to common.Address, amount *big.Int) ([]byte, error) {
	return eRC20.abi.Pack("transferFrom", from, to, amount)
}

// UnpackTransferFrom is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 amount) returns(bool)
func (eRC20 *ERC20) UnpackTransferFrom(data []byte) (bool, error) {
	out, err := eRC20.abi.Unpack("transferFrom", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// ERC20Approval represents a Approval event raised by the ERC20 contract.
type ERC20Approval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     *types.Log // Blockchain specific contextual infos
}

const ERC20ApprovalEventName = "Approval"

// ContractEventName returns the user-defined event name.
func (ERC20Approval) ContractEventName() string {
	return ERC20ApprovalEventName
}

// UnpackApprovalEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (eRC20 *ERC20) UnpackApprovalEvent(log *types.Log) (*ERC20Approval, error) {
	event := "Approval"
	if len(log.Topics) == 0 || log.Topics[0] != eRC20.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ERC20Approval)
	if len(log.Data) > 0 {
		if err := eRC20.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range eRC20.abi.Events[event].Inputs {
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

// ERC20Transfer represents a Transfer event raised by the ERC20 contract.
type ERC20Transfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   *types.Log // Blockchain specific contextual infos
}

const ERC20TransferEventName = "Transfer"

// ContractEventName returns the user-defined event name.
func (ERC20Transfer) ContractEventName() string {
	return ERC20TransferEventName
}

// UnpackTransferEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (eRC20 *ERC20) UnpackTransferEvent(log *types.Log) (*ERC20Transfer, error) {
	event := "Transfer"
	if len(log.Topics) == 0 || log.Topics[0] != eRC20.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ERC20Transfer)
	if len(log.Data) > 0 {
		if err := eRC20.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range eRC20.abi.Events[event].Inputs {
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
