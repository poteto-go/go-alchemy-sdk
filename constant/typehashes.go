package constant

import "github.com/ethereum/go-ethereum/crypto"

var (
	// EIP-2612 type hashes
	PermitTypeHash = crypto.Keccak256([]byte("Permit(address owner,address spender,uint256 value,uint256 nonce,uint256 deadline)"))

	// EIP-3009 type hashes
	TransferWithAuthorizationTypeHash = crypto.Keccak256([]byte("TransferWithAuthorization(address from,address to,uint256 value,uint256 validAfter,uint256 validBefore,bytes32 nonce)"))
	ReceiveWithAuthorizationTypeHash  = crypto.Keccak256([]byte("ReceiveWithAuthorization(address from,address to,uint256 value,uint256 validAfter,uint256 validBefore,bytes32 nonce)"))
	CancelAuthorizationTypeHash       = crypto.Keccak256([]byte("CancelAuthorization(address authorizer,bytes32 nonce)"))
)
