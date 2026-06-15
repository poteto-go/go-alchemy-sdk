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
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"ApprovalForAll\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"ids\",\"type\":\"uint256[]\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"values\",\"type\":\"uint256[]\"}],\"name\":\"TransferBatch\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"TransferSingle\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"accounts\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"ids\",\"type\":\"uint256[]\"}],\"name\":\"balanceOfBatch\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"isApprovedForAll\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256[]\",\"name\":\"ids\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"safeBatchTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"setApprovalForAll\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"uri\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	ID:  "ERC1155",
	Bin: "0x6080604052604051806060016040528060258152602001612497602591395f908161002a9190610279565b50348015610036575f5ffd5b50610348565b5f81519050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5f60028204905060018216806100b757607f821691505b6020821081036100ca576100c9610073565b5b50919050565b5f819050815f5260205f209050919050565b5f6020601f8301049050919050565b5f82821b905092915050565b5f6008830261012c7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff826100f1565b61013686836100f1565b95508019841693508086168417925050509392505050565b5f819050919050565b5f819050919050565b5f61017a6101756101708461014e565b610157565b61014e565b9050919050565b5f819050919050565b61019383610160565b6101a761019f82610181565b8484546100fd565b825550505050565b5f5f905090565b6101be6101af565b6101c981848461018a565b505050565b5b818110156101ec576101e15f826101b6565b6001810190506101cf565b5050565b601f82111561023157610202816100d0565b61020b846100e2565b8101602085101561021a578190505b61022e610226856100e2565b8301826101ce565b50505b505050565b5f82821c905092915050565b5f6102515f1984600802610236565b1980831691505092915050565b5f6102698383610242565b9150826002028217905092915050565b6102828261003c565b67ffffffffffffffff81111561029b5761029a610046565b5b6102a582546100a0565b6102b08282856101f0565b5f60209050601f8311600181146102e1575f84156102cf578287015190505b6102d9858261025e565b865550610340565b601f1984166102ef866100d0565b5f5b82811015610316578489015182556001820191506020850194506020810190506102f1565b86831015610333578489015161032f601f891682610242565b8355505b6001600288020188555050505b505050505050565b612142806103555f395ff3fe608060405234801561000f575f5ffd5b5060043610610085575f3560e01c80634e1273f4116100595780634e1273f414610121578063a22cb46514610151578063e985e9c51461016d578063f242432a1461019d57610085565b8062fdd58e146100895780630e89341c146100b9578063156e29f6146100e95780632eb2c2d614610105575b5f5ffd5b6100a3600480360381019061009e9190611151565b6101b9565b6040516100b0919061119e565b60405180910390f35b6100d360048036038101906100ce91906111b7565b61027d565b6040516100e09190611252565b60405180910390f35b61010360048036038101906100fe9190611272565b61030e565b005b61011f600480360381019061011a91906114b2565b610461565b005b61013b6004803603810190610136919061163d565b610983565b604051610148919061176a565b60405180910390f35b61016b600480360381019061016691906117bf565b610a97565b005b610187600480360381019061018291906117fd565b610bfd565b604051610194919061184a565b60405180910390f35b6101b760048036038101906101b29190611863565b610c8b565b005b5f5f73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1603610228576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161021f90611966565b60405180910390fd5b60015f8381526020019081526020015f205f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2054905092915050565b60605f805461028b906119b1565b80601f01602080910402602001604051908101604052809291908181526020018280546102b7906119b1565b80156103025780601f106102d957610100808354040283529160200191610302565b820191905f5260205f20905b8154815290600101906020018083116102e557829003601f168201915b50505050509050919050565b5f73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff160361037c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161037390611a51565b60405180910390fd5b8060015f8481526020019081526020015f205f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8282546103d79190611a9c565b925050819055508273ffffffffffffffffffffffffffffffffffffffff165f73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167fc3d58168c5ae7397731d063d5bbf3d657854427343f4c083240f7aacaa2d0f628585604051610454929190611acf565b60405180910390a4505050565b3373ffffffffffffffffffffffffffffffffffffffff168573ffffffffffffffffffffffffffffffffffffffff1614806104a157506104a08533610bfd565b5b6104e0576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016104d790611b66565b60405180910390fd5b5f73ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff160361054e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161054590611bf4565b60405180910390fd5b8151835114610592576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161058990611c82565b60405180910390fd5b5f5f90505b835181101561079a578281815181106105b3576105b2611ca0565b5b602002602001015160015f8684815181106105d1576105d0611ca0565b5b602002602001015181526020019081526020015f205f8873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20541015610663576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161065a90611d3d565b60405180910390fd5b82818151811061067657610675611ca0565b5b602002602001015160015f86848151811061069457610693611ca0565b5b602002602001015181526020019081526020015f205f8873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8282546106f29190611d5b565b9250508190555082818151811061070c5761070b611ca0565b5b602002602001015160015f86848151811061072a57610729611ca0565b5b602002602001015181526020019081526020015f205f8773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8282546107889190611a9c565b92505081905550806001019050610597565b508373ffffffffffffffffffffffffffffffffffffffff168573ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167f4a39dc06d4c0dbc64b70af90fd698a233a518aa5d07e595d983b8c0526c8f7fb8686604051610811929190611d8e565b60405180910390a45f8473ffffffffffffffffffffffffffffffffffffffff163b111561097c578373ffffffffffffffffffffffffffffffffffffffff1663bc197c8133878686866040518663ffffffff1660e01b8152600401610879959493929190611e24565b6020604051808303815f875af19250505080156108b457506040513d601f19601f820116820180604052508101906108b19190611edf565b60015b6108f3576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016108ea90611f7a565b60405180910390fd5b63bc197c8160e01b7bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916817bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19161461097a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161097190611f7a565b60405180910390fd5b505b5050505050565b606081518351146109c9576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016109c090612008565b60405180910390fd5b5f835167ffffffffffffffff8111156109e5576109e46112c6565b5b604051908082528060200260200182016040528015610a135781602001602082028036833780820191505090505b5090505f5f90505b8451811015610a8c57610a62858281518110610a3a57610a39611ca0565b5b6020026020010151858381518110610a5557610a54611ca0565b5b60200260200101516101b9565b828281518110610a7557610a74611ca0565b5b602002602001018181525050806001019050610a1b565b508091505092915050565b3373ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603610b05576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610afc90612096565b60405180910390fd5b8060025f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f6101000a81548160ff0219169083151502179055508173ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167f17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c3183604051610bf1919061184a565b60405180910390a35050565b5f60025f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f9054906101000a900460ff16905092915050565b3373ffffffffffffffffffffffffffffffffffffffff168573ffffffffffffffffffffffffffffffffffffffff161480610ccb5750610cca8533610bfd565b5b610d0a576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610d0190611b66565b60405180910390fd5b5f73ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff1603610d78576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610d6f90611bf4565b60405180910390fd5b8160015f8581526020019081526020015f205f8773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20541015610e07576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610dfe90611d3d565b60405180910390fd5b8160015f8581526020019081526020015f205f8773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f828254610e629190611d5b565b925050819055508160015f8581526020019081526020015f205f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f828254610ec49190611a9c565b925050819055508373ffffffffffffffffffffffffffffffffffffffff168573ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167fc3d58168c5ae7397731d063d5bbf3d657854427343f4c083240f7aacaa2d0f628686604051610f41929190611acf565b60405180910390a45f8473ffffffffffffffffffffffffffffffffffffffff163b11156110ac578373ffffffffffffffffffffffffffffffffffffffff1663f23a6e6133878686866040518663ffffffff1660e01b8152600401610fa99594939291906120b4565b6020604051808303815f875af1925050508015610fe457506040513d601f19601f82011682018060405250810190610fe19190611edf565b60015b611023576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161101a90611f7a565b60405180910390fd5b63f23a6e6160e01b7bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916817bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916146110aa576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016110a190611f7a565b60405180910390fd5b505b5050505050565b5f604051905090565b5f5ffd5b5f5ffd5b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f6110ed826110c4565b9050919050565b6110fd816110e3565b8114611107575f5ffd5b50565b5f81359050611118816110f4565b92915050565b5f819050919050565b6111308161111e565b811461113a575f5ffd5b50565b5f8135905061114b81611127565b92915050565b5f5f60408385031215611167576111666110bc565b5b5f6111748582860161110a565b92505060206111858582860161113d565b9150509250929050565b6111988161111e565b82525050565b5f6020820190506111b15f83018461118f565b92915050565b5f602082840312156111cc576111cb6110bc565b5b5f6111d98482850161113d565b91505092915050565b5f81519050919050565b5f82825260208201905092915050565b8281835e5f83830152505050565b5f601f19601f8301169050919050565b5f611224826111e2565b61122e81856111ec565b935061123e8185602086016111fc565b6112478161120a565b840191505092915050565b5f6020820190508181035f83015261126a818461121a565b905092915050565b5f5f5f60608486031215611289576112886110bc565b5b5f6112968682870161110a565b93505060206112a78682870161113d565b92505060406112b88682870161113d565b9150509250925092565b5f5ffd5b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b6112fc8261120a565b810181811067ffffffffffffffff8211171561131b5761131a6112c6565b5b80604052505050565b5f61132d6110b3565b905061133982826112f3565b919050565b5f67ffffffffffffffff821115611358576113576112c6565b5b602082029050602081019050919050565b5f5ffd5b5f61137f61137a8461133e565b611324565b905080838252602082019050602084028301858111156113a2576113a1611369565b5b835b818110156113cb57806113b7888261113d565b8452602084019350506020810190506113a4565b5050509392505050565b5f82601f8301126113e9576113e86112c2565b5b81356113f984826020860161136d565b91505092915050565b5f5ffd5b5f67ffffffffffffffff8211156114205761141f6112c6565b5b6114298261120a565b9050602081019050919050565b828183375f83830152505050565b5f61145661145184611406565b611324565b90508281526020810184848401111561147257611471611402565b5b61147d848285611436565b509392505050565b5f82601f830112611499576114986112c2565b5b81356114a9848260208601611444565b91505092915050565b5f5f5f5f5f60a086880312156114cb576114ca6110bc565b5b5f6114d88882890161110a565b95505060206114e98882890161110a565b945050604086013567ffffffffffffffff81111561150a576115096110c0565b5b611516888289016113d5565b935050606086013567ffffffffffffffff811115611537576115366110c0565b5b611543888289016113d5565b925050608086013567ffffffffffffffff811115611564576115636110c0565b5b61157088828901611485565b9150509295509295909350565b5f67ffffffffffffffff821115611597576115966112c6565b5b602082029050602081019050919050565b5f6115ba6115b58461157d565b611324565b905080838252602082019050602084028301858111156115dd576115dc611369565b5b835b8181101561160657806115f2888261110a565b8452602084019350506020810190506115df565b5050509392505050565b5f82601f830112611624576116236112c2565b5b81356116348482602086016115a8565b91505092915050565b5f5f60408385031215611653576116526110bc565b5b5f83013567ffffffffffffffff8111156116705761166f6110c0565b5b61167c85828601611610565b925050602083013567ffffffffffffffff81111561169d5761169c6110c0565b5b6116a9858286016113d5565b9150509250929050565b5f81519050919050565b5f82825260208201905092915050565b5f819050602082019050919050565b6116e58161111e565b82525050565b5f6116f683836116dc565b60208301905092915050565b5f602082019050919050565b5f611718826116b3565b61172281856116bd565b935061172d836116cd565b805f5b8381101561175d57815161174488826116eb565b975061174f83611702565b925050600181019050611730565b5085935050505092915050565b5f6020820190508181035f830152611782818461170e565b905092915050565b5f8115159050919050565b61179e8161178a565b81146117a8575f5ffd5b50565b5f813590506117b981611795565b92915050565b5f5f604083850312156117d5576117d46110bc565b5b5f6117e28582860161110a565b92505060206117f3858286016117ab565b9150509250929050565b5f5f60408385031215611813576118126110bc565b5b5f6118208582860161110a565b92505060206118318582860161110a565b9150509250929050565b6118448161178a565b82525050565b5f60208201905061185d5f83018461183b565b92915050565b5f5f5f5f5f60a0868803121561187c5761187b6110bc565b5b5f6118898882890161110a565b955050602061189a8882890161110a565b94505060406118ab8882890161113d565b93505060606118bc8882890161113d565b925050608086013567ffffffffffffffff8111156118dd576118dc6110c0565b5b6118e988828901611485565b9150509295509295909350565b7f455243313135353a2061646472657373207a65726f206973206e6f74206120765f8201527f616c6964206f776e657200000000000000000000000000000000000000000000602082015250565b5f611950602a836111ec565b915061195b826118f6565b604082019050919050565b5f6020820190508181035f83015261197d81611944565b9050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5f60028204905060018216806119c857607f821691505b6020821081036119db576119da611984565b5b50919050565b7f455243313135353a206d696e7420746f20746865207a65726f206164647265735f8201527f7300000000000000000000000000000000000000000000000000000000000000602082015250565b5f611a3b6021836111ec565b9150611a46826119e1565b604082019050919050565b5f6020820190508181035f830152611a6881611a2f565b9050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f611aa68261111e565b9150611ab18361111e565b9250828201905080821115611ac957611ac8611a6f565b5b92915050565b5f604082019050611ae25f83018561118f565b611aef602083018461118f565b9392505050565b7f455243313135353a2063616c6c6572206973206e6f7420746f6b656e206f776e5f8201527f6572206f7220617070726f766564000000000000000000000000000000000000602082015250565b5f611b50602e836111ec565b9150611b5b82611af6565b604082019050919050565b5f6020820190508181035f830152611b7d81611b44565b9050919050565b7f455243313135353a207472616e7366657220746f20746865207a65726f2061645f8201527f6472657373000000000000000000000000000000000000000000000000000000602082015250565b5f611bde6025836111ec565b9150611be982611b84565b604082019050919050565b5f6020820190508181035f830152611c0b81611bd2565b9050919050565b7f455243313135353a2069647320616e6420616d6f756e7473206c656e677468205f8201527f6d69736d61746368000000000000000000000000000000000000000000000000602082015250565b5f611c6c6028836111ec565b9150611c7782611c12565b604082019050919050565b5f6020820190508181035f830152611c9981611c60565b9050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603260045260245ffd5b7f455243313135353a20696e73756666696369656e742062616c616e636520666f5f8201527f72207472616e7366657200000000000000000000000000000000000000000000602082015250565b5f611d27602a836111ec565b9150611d3282611ccd565b604082019050919050565b5f6020820190508181035f830152611d5481611d1b565b9050919050565b5f611d658261111e565b9150611d708361111e565b9250828203905081811115611d8857611d87611a6f565b5b92915050565b5f6040820190508181035f830152611da6818561170e565b90508181036020830152611dba818461170e565b90509392505050565b611dcc816110e3565b82525050565b5f81519050919050565b5f82825260208201905092915050565b5f611df682611dd2565b611e008185611ddc565b9350611e108185602086016111fc565b611e198161120a565b840191505092915050565b5f60a082019050611e375f830188611dc3565b611e446020830187611dc3565b8181036040830152611e56818661170e565b90508181036060830152611e6a818561170e565b90508181036080830152611e7e8184611dec565b90509695505050505050565b5f7fffffffff0000000000000000000000000000000000000000000000000000000082169050919050565b611ebe81611e8a565b8114611ec8575f5ffd5b50565b5f81519050611ed981611eb5565b92915050565b5f60208284031215611ef457611ef36110bc565b5b5f611f0184828501611ecb565b91505092915050565b7f455243313135353a207472616e7366657220746f206e6f6e2d455243313135355f8201527f5265636569766572000000000000000000000000000000000000000000000000602082015250565b5f611f646028836111ec565b9150611f6f82611f0a565b604082019050919050565b5f6020820190508181035f830152611f9181611f58565b9050919050565b7f455243313135353a206163636f756e747320616e6420696473206c656e6774685f8201527f206d69736d617463680000000000000000000000000000000000000000000000602082015250565b5f611ff26029836111ec565b9150611ffd82611f98565b604082019050919050565b5f6020820190508181035f83015261201f81611fe6565b9050919050565b7f455243313135353a2073657474696e6720617070726f76616c207374617475735f8201527f20666f722073656c660000000000000000000000000000000000000000000000602082015250565b5f6120806029836111ec565b915061208b82612026565b604082019050919050565b5f6020820190508181035f8301526120ad81612074565b9050919050565b5f60a0820190506120c75f830188611dc3565b6120d46020830187611dc3565b6120e1604083018661118f565b6120ee606083018561118f565b81810360808301526121008184611dec565b9050969550505050505056fea2646970667358221220f3f25ca720084bf5ee0d779258efd1cc16b45db9126d2673ac29dd112df4b37e64736f6c634300081e003368747470733a2f2f6578616d706c652e636f6d2f657263313135352f7b69647d2e6a736f6e",
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

// PackSafeBatchTransferFrom is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2eb2c2d6.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function safeBatchTransferFrom(address from, address to, uint256[] ids, uint256[] amounts, bytes data) returns()
func (eRC1155 *ERC1155) PackSafeBatchTransferFrom(from common.Address, to common.Address, ids []*big.Int, amounts []*big.Int, data []byte) []byte {
	enc, err := eRC1155.abi.Pack("safeBatchTransferFrom", from, to, ids, amounts, data)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSafeBatchTransferFrom is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2eb2c2d6.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function safeBatchTransferFrom(address from, address to, uint256[] ids, uint256[] amounts, bytes data) returns()
func (eRC1155 *ERC1155) TryPackSafeBatchTransferFrom(from common.Address, to common.Address, ids []*big.Int, amounts []*big.Int, data []byte) ([]byte, error) {
	return eRC1155.abi.Pack("safeBatchTransferFrom", from, to, ids, amounts, data)
}

// PackSafeTransferFrom is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf242432a.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 id, uint256 amount, bytes data) returns()
func (eRC1155 *ERC1155) PackSafeTransferFrom(from common.Address, to common.Address, id *big.Int, amount *big.Int, data []byte) []byte {
	enc, err := eRC1155.abi.Pack("safeTransferFrom", from, to, id, amount, data)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSafeTransferFrom is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf242432a.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 id, uint256 amount, bytes data) returns()
func (eRC1155 *ERC1155) TryPackSafeTransferFrom(from common.Address, to common.Address, id *big.Int, amount *big.Int, data []byte) ([]byte, error) {
	return eRC1155.abi.Pack("safeTransferFrom", from, to, id, amount, data)
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

// ERC1155TransferBatch represents a TransferBatch event raised by the ERC1155 contract.
type ERC1155TransferBatch struct {
	Operator common.Address
	From     common.Address
	To       common.Address
	Ids      []*big.Int
	Values   []*big.Int
	Raw      *types.Log // Blockchain specific contextual infos
}

const ERC1155TransferBatchEventName = "TransferBatch"

// ContractEventName returns the user-defined event name.
func (ERC1155TransferBatch) ContractEventName() string {
	return ERC1155TransferBatchEventName
}

// UnpackTransferBatchEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event TransferBatch(address indexed operator, address indexed from, address indexed to, uint256[] ids, uint256[] values)
func (eRC1155 *ERC1155) UnpackTransferBatchEvent(log *types.Log) (*ERC1155TransferBatch, error) {
	event := "TransferBatch"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != eRC1155.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(ERC1155TransferBatch)
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
