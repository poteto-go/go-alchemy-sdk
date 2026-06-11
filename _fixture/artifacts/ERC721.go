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

// ERC721MetaData contains all meta data concerning the ERC721 contract.
var ERC721MetaData = bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"approved\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"ApprovalForAll\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getApproved\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"isApprovedForAll\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ownerOf\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"setApprovalForAll\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"tokenURI\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	ID:  "ERC721",
	Bin: "0x60c0604052600b60809081526a135a5b9a5b585b0813919560aa1b60a0525f9061002990826100fc565b506040805180820190915260048152631353919560e21b602082015260019061005290826100fc565b5034801561005e575f5ffd5b506101b6565b634e487b7160e01b5f52604160045260245ffd5b600181811c9082168061008c57607f821691505b6020821081036100aa57634e487b7160e01b5f52602260045260245ffd5b50919050565b601f8211156100f757805f5260205f20601f840160051c810160208510156100d55750805b601f840160051c820191505b818110156100f4575f81556001016100e1565b50505b505050565b81516001600160401b0381111561011557610115610064565b610129816101238454610078565b846100b0565b6020601f82116001811461015b575f83156101445750848201515b5f19600385901b1c1916600184901b1784556100f4565b5f84815260208120601f198516915b8281101561018a578785015182556020948501946001909201910161016a565b50848210156101a757868401515f19600387901b60f8161c191681555b50505050600190811b01905550565b610db3806101c35f395ff3fe608060405234801561000f575f5ffd5b50600436106100a6575f3560e01c80636352211e1161006e5780636352211e1461012e57806370a082311461014157806395d89b4114610162578063a22cb4651461016a578063c87b56dd1461017d578063e985e9c514610190575f5ffd5b806306fdde03146100aa578063081812fc146100c8578063095ea7b3146100f357806323b872dd1461010857806340c10f191461011b575b5f5ffd5b6100b26101b3565b6040516100bf9190610ac8565b60405180910390f35b6100db6100d6366004610afd565b61023e565b6040516001600160a01b0390911681526020016100bf565b610106610101366004610b2f565b610295565b005b610106610116366004610b57565b6103f5565b610106610129366004610b2f565b61060e565b6100db61013c366004610afd565b61074d565b61015461014f366004610b91565b610787565b6040519081526020016100bf565b6100b261080b565b610106610178366004610bb1565b610818565b6100b261018b366004610afd565b6108db565b6101a361019e366004610bea565b610940565b60405190151581526020016100bf565b5f80546101bf90610c1b565b80601f01602080910402602001604051908101604052809291908181526020018280546101eb90610c1b565b80156102365780601f1061020d57610100808354040283529160200191610236565b820191905f5260205f20905b81548152906001019060200180831161021957829003601f168201915b505050505081565b5f818152600260205260408120546001600160a01b031661027a5760405162461bcd60e51b815260040161027190610c53565b60405180910390fd5b505f908152600460205260409020546001600160a01b031690565b5f61029f8261074d565b9050806001600160a01b0316836001600160a01b03160361030c5760405162461bcd60e51b815260206004820152602160248201527f4552433732313a20617070726f76616c20746f2063757272656e74206f776e656044820152603960f91b6064820152608401610271565b336001600160a01b038216148061032857506103288133610940565b61039a5760405162461bcd60e51b815260206004820152603d60248201527f4552433732313a20617070726f76652063616c6c6572206973206e6f7420746f60448201527f6b656e206f776e6572206f7220617070726f76656420666f7220616c6c0000006064820152608401610271565b5f8281526004602052604080822080546001600160a01b0319166001600160a01b0387811691821790925591518593918516917f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92591a4505050565b6103ff338261096d565b6104615760405162461bcd60e51b815260206004820152602d60248201527f4552433732313a2063616c6c6572206973206e6f7420746f6b656e206f776e6560448201526c1c881bdc88185c1c1c9bdd9959609a1b6064820152608401610271565b6001600160a01b0382166104c35760405162461bcd60e51b8152602060048201526024808201527f4552433732313a207472616e7366657220746f20746865207a65726f206164646044820152637265737360e01b6064820152608401610271565b5f6104cd8261074d565b9050836001600160a01b0316816001600160a01b03161461053e5760405162461bcd60e51b815260206004820152602560248201527f4552433732313a207472616e736665722066726f6d20696e636f72726563742060448201526437bbb732b960d91b6064820152608401610271565b5f82815260046020908152604080832080546001600160a01b03191690556001600160a01b038716835260039091528120805460019290610580908490610c9e565b90915550506001600160a01b0383165f9081526003602052604081208054600192906105ad908490610cb1565b90915550505f8281526002602052604080822080546001600160a01b0319166001600160a01b0387811691821790925591518593918816917fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef91a450505050565b6001600160a01b0382166106645760405162461bcd60e51b815260206004820181905260248201527f4552433732313a206d696e7420746f20746865207a65726f20616464726573736044820152606401610271565b5f818152600260205260409020546001600160a01b0316156106c85760405162461bcd60e51b815260206004820152601c60248201527f4552433732313a20746f6b656e20616c7265616479206d696e746564000000006044820152606401610271565b6001600160a01b0382165f9081526003602052604081208054600192906106f0908490610cb1565b90915550505f8181526002602052604080822080546001600160a01b0319166001600160a01b03861690811790915590518392907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef908290a45050565b5f818152600260205260408120546001600160a01b0316806107815760405162461bcd60e51b815260040161027190610c53565b92915050565b5f6001600160a01b0382166107f05760405162461bcd60e51b815260206004820152602960248201527f4552433732313a2061646472657373207a65726f206973206e6f7420612076616044820152683634b21037bbb732b960b91b6064820152608401610271565b506001600160a01b03165f9081526003602052604090205490565b600180546101bf90610c1b565b336001600160a01b038316036108705760405162461bcd60e51b815260206004820152601960248201527f4552433732313a20617070726f766520746f2063616c6c6572000000000000006044820152606401610271565b335f8181526005602090815260408083206001600160a01b03871680855290835292819020805460ff191686151590811790915590519081529192917f17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31910160405180910390a35050565b5f818152600260205260409020546060906001600160a01b03166109115760405162461bcd60e51b815260040161027190610c53565b61091a826109cb565b60405160200161092a9190610cc4565b6040516020818303038152906040529050919050565b6001600160a01b039182165f90815260056020908152604080832093909416825291909152205460ff1690565b5f5f6109788361074d565b9050806001600160a01b0316846001600160a01b0316148061099f575061099f8185610940565b806109c35750836001600160a01b03166109b88461023e565b6001600160a01b0316145b949350505050565b6060815f036109f15750506040805180820190915260018152600360fc1b602082015290565b815f5b8115610a1a5780610a0481610d03565b9150610a139050600a83610d2f565b91506109f4565b5f8167ffffffffffffffff811115610a3457610a34610d42565b6040519080825280601f01601f191660200182016040528015610a5e576020820181803683370190505b5090505b84156109c357610a73600183610c9e565b9150610a80600a86610d56565b610a8b906030610cb1565b60f81b818381518110610aa057610aa0610d69565b60200101906001600160f81b03191690815f1a905350610ac1600a86610d2f565b9450610a62565b602081525f82518060208401528060208501604085015e5f604082850101526040601f19601f83011684010191505092915050565b5f60208284031215610b0d575f5ffd5b5035919050565b80356001600160a01b0381168114610b2a575f5ffd5b919050565b5f5f60408385031215610b40575f5ffd5b610b4983610b14565b946020939093013593505050565b5f5f5f60608486031215610b69575f5ffd5b610b7284610b14565b9250610b8060208501610b14565b929592945050506040919091013590565b5f60208284031215610ba1575f5ffd5b610baa82610b14565b9392505050565b5f5f60408385031215610bc2575f5ffd5b610bcb83610b14565b915060208301358015158114610bdf575f5ffd5b809150509250929050565b5f5f60408385031215610bfb575f5ffd5b610c0483610b14565b9150610c1260208401610b14565b90509250929050565b600181811c90821680610c2f57607f821691505b602082108103610c4d57634e487b7160e01b5f52602260045260245ffd5b50919050565b60208082526018908201527f4552433732313a20696e76616c696420746f6b656e2049440000000000000000604082015260600190565b634e487b7160e01b5f52601160045260245ffd5b8181038181111561078157610781610c8a565b8082018082111561078157610781610c8a565b7f68747470733a2f2f6578616d706c652e636f6d2f6e66742f000000000000000081525f82518060208501601885015e5f920160180191825250919050565b5f60018201610d1457610d14610c8a565b5060010190565b634e487b7160e01b5f52601260045260245ffd5b5f82610d3d57610d3d610d1b565b500490565b634e487b7160e01b5f52604160045260245ffd5b5f82610d6457610d64610d1b565b500690565b634e487b7160e01b5f52603260045260245ffdfea2646970667358221220d2a4e4678f35fa32056390063c2d2bdc05f0759fee2b2e57ac28b730b141acb964736f6c634300081e0033",
}

// ERC721 is an auto generated Go binding around an Ethereum contract.
type ERC721 struct {
	abi abi.ABI
}

// NewERC721 creates a new instance of ERC721.
func NewERC721() *ERC721 {
	parsed, err := ERC721MetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &ERC721{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *ERC721) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackApprove is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x095ea7b3.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (eRC721 *ERC721) PackApprove(to common.Address, tokenId *big.Int) []byte {
	enc, err := eRC721.abi.Pack("approve", to, tokenId)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackApprove is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x095ea7b3.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (eRC721 *ERC721) TryPackApprove(to common.Address, tokenId *big.Int) ([]byte, error) {
	return eRC721.abi.Pack("approve", to, tokenId)
}

// PackBalanceOf is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x70a08231.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (eRC721 *ERC721) PackBalanceOf(owner common.Address) []byte {
	enc, err := eRC721.abi.Pack("balanceOf", owner)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackBalanceOf is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x70a08231.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (eRC721 *ERC721) TryPackBalanceOf(owner common.Address) ([]byte, error) {
	return eRC721.abi.Pack("balanceOf", owner)
}

// UnpackBalanceOf is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (eRC721 *ERC721) UnpackBalanceOf(data []byte) (*big.Int, error) {
	out, err := eRC721.abi.Unpack("balanceOf", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackGetApproved is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x081812fc.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (eRC721 *ERC721) PackGetApproved(tokenId *big.Int) []byte {
	enc, err := eRC721.abi.Pack("getApproved", tokenId)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetApproved is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x081812fc.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (eRC721 *ERC721) TryPackGetApproved(tokenId *big.Int) ([]byte, error) {
	return eRC721.abi.Pack("getApproved", tokenId)
}

// UnpackGetApproved is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (eRC721 *ERC721) UnpackGetApproved(data []byte) (common.Address, error) {
	out, err := eRC721.abi.Unpack("getApproved", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackIsApprovedForAll is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe985e9c5.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (eRC721 *ERC721) PackIsApprovedForAll(owner common.Address, operator common.Address) []byte {
	enc, err := eRC721.abi.Pack("isApprovedForAll", owner, operator)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackIsApprovedForAll is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe985e9c5.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (eRC721 *ERC721) TryPackIsApprovedForAll(owner common.Address, operator common.Address) ([]byte, error) {
	return eRC721.abi.Pack("isApprovedForAll", owner, operator)
}

// UnpackIsApprovedForAll is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (eRC721 *ERC721) UnpackIsApprovedForAll(data []byte) (bool, error) {
	out, err := eRC721.abi.Unpack("isApprovedForAll", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackMint is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x40c10f19.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function mint(address to, uint256 tokenId) returns()
func (eRC721 *ERC721) PackMint(to common.Address, tokenId *big.Int) []byte {
	enc, err := eRC721.abi.Pack("mint", to, tokenId)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMint is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x40c10f19.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function mint(address to, uint256 tokenId) returns()
func (eRC721 *ERC721) TryPackMint(to common.Address, tokenId *big.Int) ([]byte, error) {
	return eRC721.abi.Pack("mint", to, tokenId)
}

// PackName is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x06fdde03.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function name() view returns(string)
func (eRC721 *ERC721) PackName() []byte {
	enc, err := eRC721.abi.Pack("name")
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
func (eRC721 *ERC721) TryPackName() ([]byte, error) {
	return eRC721.abi.Pack("name")
}

// UnpackName is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (eRC721 *ERC721) UnpackName(data []byte) (string, error) {
	out, err := eRC721.abi.Unpack("name", data)
	if err != nil {
		return *new(string), err
	}
	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, nil
}

// PackOwnerOf is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6352211e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (eRC721 *ERC721) PackOwnerOf(tokenId *big.Int) []byte {
	enc, err := eRC721.abi.Pack("ownerOf", tokenId)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackOwnerOf is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6352211e.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (eRC721 *ERC721) TryPackOwnerOf(tokenId *big.Int) ([]byte, error) {
	return eRC721.abi.Pack("ownerOf", tokenId)
}

// UnpackOwnerOf is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (eRC721 *ERC721) UnpackOwnerOf(data []byte) (common.Address, error) {
	out, err := eRC721.abi.Unpack("ownerOf", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackSetApprovalForAll is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa22cb465.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (eRC721 *ERC721) PackSetApprovalForAll(operator common.Address, approved bool) []byte {
	enc, err := eRC721.abi.Pack("setApprovalForAll", operator, approved)
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
func (eRC721 *ERC721) TryPackSetApprovalForAll(operator common.Address, approved bool) ([]byte, error) {
	return eRC721.abi.Pack("setApprovalForAll", operator, approved)
}

// PackSymbol is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x95d89b41.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function symbol() view returns(string)
func (eRC721 *ERC721) PackSymbol() []byte {
	enc, err := eRC721.abi.Pack("symbol")
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
func (eRC721 *ERC721) TryPackSymbol() ([]byte, error) {
	return eRC721.abi.Pack("symbol")
}

// UnpackSymbol is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (eRC721 *ERC721) UnpackSymbol(data []byte) (string, error) {
	out, err := eRC721.abi.Unpack("symbol", data)
	if err != nil {
		return *new(string), err
	}
	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, nil
}

// PackTokenURI is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc87b56dd.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (eRC721 *ERC721) PackTokenURI(tokenId *big.Int) []byte {
	enc, err := eRC721.abi.Pack("tokenURI", tokenId)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTokenURI is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc87b56dd.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (eRC721 *ERC721) TryPackTokenURI(tokenId *big.Int) ([]byte, error) {
	return eRC721.abi.Pack("tokenURI", tokenId)
}

// UnpackTokenURI is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (eRC721 *ERC721) UnpackTokenURI(data []byte) (string, error) {
	out, err := eRC721.abi.Unpack("tokenURI", data)
	if err != nil {
		return *new(string), err
	}
	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, nil
}

// PackTransferFrom is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x23b872dd.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (eRC721 *ERC721) PackTransferFrom(from common.Address, to common.Address, tokenId *big.Int) []byte {
	enc, err := eRC721.abi.Pack("transferFrom", from, to, tokenId)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTransferFrom is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x23b872dd.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (eRC721 *ERC721) TryPackTransferFrom(from common.Address, to common.Address, tokenId *big.Int) ([]byte, error) {
	return eRC721.abi.Pack("transferFrom", from, to, tokenId)
}

// ERC721Approval represents a Approval event raised by the ERC721 contract.
type ERC721Approval struct {
	Owner    common.Address
	Approved common.Address
	TokenId  *big.Int
	Raw      *types.Log // Blockchain specific contextual infos
}

const ERC721ApprovalEventName = "Approval"

// ContractEventName returns the user-defined event name.
func (ERC721Approval) ContractEventName() string {
	return ERC721ApprovalEventName
}

// UnpackApprovalEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (eRC721 *ERC721) UnpackApprovalEvent(log *types.Log) (*ERC721Approval, error) {
	event := "Approval"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != eRC721.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(ERC721Approval)
	if len(log.Data) > 0 {
		if err := eRC721.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range eRC721.abi.Events[event].Inputs {
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

// ERC721ApprovalForAll represents a ApprovalForAll event raised by the ERC721 contract.
type ERC721ApprovalForAll struct {
	Owner    common.Address
	Operator common.Address
	Approved bool
	Raw      *types.Log // Blockchain specific contextual infos
}

const ERC721ApprovalForAllEventName = "ApprovalForAll"

// ContractEventName returns the user-defined event name.
func (ERC721ApprovalForAll) ContractEventName() string {
	return ERC721ApprovalForAllEventName
}

// UnpackApprovalForAllEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (eRC721 *ERC721) UnpackApprovalForAllEvent(log *types.Log) (*ERC721ApprovalForAll, error) {
	event := "ApprovalForAll"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != eRC721.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(ERC721ApprovalForAll)
	if len(log.Data) > 0 {
		if err := eRC721.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range eRC721.abi.Events[event].Inputs {
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

// ERC721Transfer represents a Transfer event raised by the ERC721 contract.
type ERC721Transfer struct {
	From    common.Address
	To      common.Address
	TokenId *big.Int
	Raw     *types.Log // Blockchain specific contextual infos
}

const ERC721TransferEventName = "Transfer"

// ContractEventName returns the user-defined event name.
func (ERC721Transfer) ContractEventName() string {
	return ERC721TransferEventName
}

// UnpackTransferEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (eRC721 *ERC721) UnpackTransferEvent(log *types.Log) (*ERC721Transfer, error) {
	event := "Transfer"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != eRC721.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(ERC721Transfer)
	if len(log.Data) > 0 {
		if err := eRC721.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range eRC721.abi.Events[event].Inputs {
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
