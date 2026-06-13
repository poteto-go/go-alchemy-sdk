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
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"approved\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"ApprovalForAll\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getApproved\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"isApprovedForAll\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ownerOf\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"setApprovalForAll\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"tokenURI\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	ID:  "ERC721",
	Bin: "0x60c0604052600b60809081526a135a5b9a5b585b0813919560aa1b60a0525f9061002990826100fc565b506040805180820190915260048152631353919560e21b602082015260019061005290826100fc565b5034801561005e575f5ffd5b506101b6565b634e487b7160e01b5f52604160045260245ffd5b600181811c9082168061008c57607f821691505b6020821081036100aa57634e487b7160e01b5f52602260045260245ffd5b50919050565b601f8211156100f757805f5260205f20601f840160051c810160208510156100d55750805b601f840160051c820191505b818110156100f4575f81556001016100e1565b50505b505050565b81516001600160401b0381111561011557610115610064565b610129816101238454610078565b846100b0565b6020601f82116001811461015b575f83156101445750848201515b5f19600385901b1c1916600184901b1784556100f4565b5f84815260208120601f198516915b8281101561018a578785015182556020948501946001909201910161016a565b50848210156101a757868401515f19600387901b60f8161c191681555b50505050600190811b01905550565b61109a806101c35f395ff3fe608060405234801561000f575f5ffd5b50600436106100cb575f3560e01c80636352211e11610088578063a22cb46511610063578063a22cb465146101a2578063b88d4fde146101b5578063c87b56dd146101c8578063e985e9c5146101db575f5ffd5b80636352211e1461016657806370a082311461017957806395d89b411461019a575f5ffd5b806306fdde03146100cf578063081812fc146100ed578063095ea7b31461011857806323b872dd1461012d57806340c10f191461014057806342842e0e14610153575b5f5ffd5b6100d76101fe565b6040516100e49190610c92565b60405180910390f35b6101006100fb366004610cab565b610289565b6040516001600160a01b0390911681526020016100e4565b61012b610126366004610cdd565b6102e0565b005b61012b61013b366004610d05565b610440565b61012b61014e366004610cdd565b610659565b61012b610161366004610d05565b610798565b610100610174366004610cab565b6107b7565b61018c610187366004610d3f565b6107f1565b6040519081526020016100e4565b6100d7610875565b61012b6101b0366004610d58565b610882565b61012b6101c3366004610da5565b610945565b6100d76101d6366004610cab565b6109c9565b6101ee6101e9366004610e82565b610a2e565b60405190151581526020016100e4565b5f805461020a90610eb3565b80601f016020809104026020016040519081016040528092919081815260200182805461023690610eb3565b80156102815780601f1061025857610100808354040283529160200191610281565b820191905f5260205f20905b81548152906001019060200180831161026457829003601f168201915b505050505081565b5f818152600260205260408120546001600160a01b03166102c55760405162461bcd60e51b81526004016102bc90610eeb565b60405180910390fd5b505f908152600460205260409020546001600160a01b031690565b5f6102ea826107b7565b9050806001600160a01b0316836001600160a01b0316036103575760405162461bcd60e51b815260206004820152602160248201527f4552433732313a20617070726f76616c20746f2063757272656e74206f776e656044820152603960f91b60648201526084016102bc565b336001600160a01b038216148061037357506103738133610a2e565b6103e55760405162461bcd60e51b815260206004820152603d60248201527f4552433732313a20617070726f76652063616c6c6572206973206e6f7420746f60448201527f6b656e206f776e6572206f7220617070726f76656420666f7220616c6c00000060648201526084016102bc565b5f8281526004602052604080822080546001600160a01b0319166001600160a01b0387811691821790925591518593918516917f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92591a4505050565b61044a3382610a5b565b6104ac5760405162461bcd60e51b815260206004820152602d60248201527f4552433732313a2063616c6c6572206973206e6f7420746f6b656e206f776e6560448201526c1c881bdc88185c1c1c9bdd9959609a1b60648201526084016102bc565b6001600160a01b03821661050e5760405162461bcd60e51b8152602060048201526024808201527f4552433732313a207472616e7366657220746f20746865207a65726f206164646044820152637265737360e01b60648201526084016102bc565b5f610518826107b7565b9050836001600160a01b0316816001600160a01b0316146105895760405162461bcd60e51b815260206004820152602560248201527f4552433732313a207472616e736665722066726f6d20696e636f72726563742060448201526437bbb732b960d91b60648201526084016102bc565b5f82815260046020908152604080832080546001600160a01b03191690556001600160a01b0387168352600390915281208054600192906105cb908490610f36565b90915550506001600160a01b0383165f9081526003602052604081208054600192906105f8908490610f49565b90915550505f8281526002602052604080822080546001600160a01b0319166001600160a01b0387811691821790925591518593918816917fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef91a450505050565b6001600160a01b0382166106af5760405162461bcd60e51b815260206004820181905260248201527f4552433732313a206d696e7420746f20746865207a65726f206164647265737360448201526064016102bc565b5f818152600260205260409020546001600160a01b0316156107135760405162461bcd60e51b815260206004820152601c60248201527f4552433732313a20746f6b656e20616c7265616479206d696e7465640000000060448201526064016102bc565b6001600160a01b0382165f90815260036020526040812080546001929061073b908490610f49565b90915550505f8181526002602052604080822080546001600160a01b0319166001600160a01b03861690811790915590518392907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef908290a45050565b6107b283838360405180602001604052805f815250610945565b505050565b5f818152600260205260408120546001600160a01b0316806107eb5760405162461bcd60e51b81526004016102bc90610eeb565b92915050565b5f6001600160a01b03821661085a5760405162461bcd60e51b815260206004820152602960248201527f4552433732313a2061646472657373207a65726f206973206e6f7420612076616044820152683634b21037bbb732b960b91b60648201526084016102bc565b506001600160a01b03165f9081526003602052604090205490565b6001805461020a90610eb3565b336001600160a01b038316036108da5760405162461bcd60e51b815260206004820152601960248201527f4552433732313a20617070726f766520746f2063616c6c65720000000000000060448201526064016102bc565b335f8181526005602090815260408083206001600160a01b03871680855290835292819020805460ff191686151590811790915590519081529192917f17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31910160405180910390a35050565b610950848484610440565b61095c84848484610ab9565b6109c35760405162461bcd60e51b815260206004820152603260248201527f4552433732313a207472616e7366657220746f206e6f6e20455243373231526560448201527131b2b4bb32b91034b6b83632b6b2b73a32b960711b60648201526084016102bc565b50505050565b5f818152600260205260409020546060906001600160a01b03166109ff5760405162461bcd60e51b81526004016102bc90610eeb565b610a0882610b67565b604051602001610a189190610f5c565b6040516020818303038152906040529050919050565b6001600160a01b039182165f90815260056020908152604080832093909416825291909152205460ff1690565b5f5f610a66836107b7565b9050806001600160a01b0316846001600160a01b03161480610a8d5750610a8d8185610a2e565b80610ab15750836001600160a01b0316610aa684610289565b6001600160a01b0316145b949350505050565b5f836001600160a01b03163b5f03610ad357506001610ab1565b604051630a85bd0160e11b81526001600160a01b0385169063150b7a0290610b05903390899088908890600401610f9b565b6020604051808303815f875af1925050508015610b3f575060408051601f3d908101601f19168201909252610b3c91810190610fd7565b60015b610b4a57505f610ab1565b6001600160e01b031916630a85bd0160e11b149050949350505050565b6060815f03610b8d5750506040805180820190915260018152600360fc1b602082015290565b815f5b8115610bb65780610ba081610ffe565b9150610baf9050600a8361102a565b9150610b90565b5f8167ffffffffffffffff811115610bd057610bd0610d91565b6040519080825280601f01601f191660200182016040528015610bfa576020820181803683370190505b5090505b8415610ab157610c0f600183610f36565b9150610c1c600a8661103d565b610c27906030610f49565b60f81b818381518110610c3c57610c3c611050565b60200101906001600160f81b03191690815f1a905350610c5d600a8661102a565b9450610bfe565b5f81518084528060208401602086015e5f602082860101526020601f19601f83011685010191505092915050565b602081525f610ca46020830184610c64565b9392505050565b5f60208284031215610cbb575f5ffd5b5035919050565b80356001600160a01b0381168114610cd8575f5ffd5b919050565b5f5f60408385031215610cee575f5ffd5b610cf783610cc2565b946020939093013593505050565b5f5f5f60608486031215610d17575f5ffd5b610d2084610cc2565b9250610d2e60208501610cc2565b929592945050506040919091013590565b5f60208284031215610d4f575f5ffd5b610ca482610cc2565b5f5f60408385031215610d69575f5ffd5b610d7283610cc2565b915060208301358015158114610d86575f5ffd5b809150509250929050565b634e487b7160e01b5f52604160045260245ffd5b5f5f5f5f60808587031215610db8575f5ffd5b610dc185610cc2565b9350610dcf60208601610cc2565b925060408501359150606085013567ffffffffffffffff811115610df1575f5ffd5b8501601f81018713610e01575f5ffd5b803567ffffffffffffffff811115610e1b57610e1b610d91565b604051601f8201601f19908116603f0116810167ffffffffffffffff81118282101715610e4a57610e4a610d91565b604052818152828201602001891015610e61575f5ffd5b816020840160208301375f6020838301015280935050505092959194509250565b5f5f60408385031215610e93575f5ffd5b610e9c83610cc2565b9150610eaa60208401610cc2565b90509250929050565b600181811c90821680610ec757607f821691505b602082108103610ee557634e487b7160e01b5f52602260045260245ffd5b50919050565b60208082526018908201527f4552433732313a20696e76616c696420746f6b656e2049440000000000000000604082015260600190565b634e487b7160e01b5f52601160045260245ffd5b818103818111156107eb576107eb610f22565b808201808211156107eb576107eb610f22565b7f68747470733a2f2f6578616d706c652e636f6d2f6e66742f000000000000000081525f82518060208501601885015e5f920160180191825250919050565b6001600160a01b03858116825284166020820152604081018390526080606082018190525f90610fcd90830184610c64565b9695505050505050565b5f60208284031215610fe7575f5ffd5b81516001600160e01b031981168114610ca4575f5ffd5b5f6001820161100f5761100f610f22565b5060010190565b634e487b7160e01b5f52601260045260245ffd5b5f8261103857611038611016565b500490565b5f8261104b5761104b611016565b500690565b634e487b7160e01b5f52603260045260245ffdfea2646970667358221220d3e56928ede6c41df13eca0dbf8d34ab29fe50ee5a694cde294ce1aafa69ea5a64736f6c634300081e0033",
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

// PackSafeTransferFrom is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x42842e0e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (eRC721 *ERC721) PackSafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) []byte {
	enc, err := eRC721.abi.Pack("safeTransferFrom", from, to, tokenId)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSafeTransferFrom is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x42842e0e.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (eRC721 *ERC721) TryPackSafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) ([]byte, error) {
	return eRC721.abi.Pack("safeTransferFrom", from, to, tokenId)
}

// PackSafeTransferFrom0 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb88d4fde.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (eRC721 *ERC721) PackSafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) []byte {
	enc, err := eRC721.abi.Pack("safeTransferFrom0", from, to, tokenId, data)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSafeTransferFrom0 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb88d4fde.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (eRC721 *ERC721) TryPackSafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) ([]byte, error) {
	return eRC721.abi.Pack("safeTransferFrom0", from, to, tokenId, data)
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
