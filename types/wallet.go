package types

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
)

// Wallet class inherits Signer and can sign transactions and messages using
type Wallet interface {
	// get address of wallet
	GetAddress() string

	// get balance of native token
	GetBalance() (balance *big.Int, err error)

	// connect provider to wallet
	Connect(provider IAlchemyProvider)

	/*
		PendingNonceAt returns the account nonce of the given account in the pending state.
		This is the nonce that should be used for the next transaction.

		internal call geth
	*/
	PendingNonceAt() (nonce uint64, err error)

	/*
		sign Transaction by wallet's p8 key
		using latest EIP155Signer

		EIP155Signer sign w/ ChainID to protect replay-attack
	*/
	SignTx(txRequest TransactionRequest) (signedTx *gethTypes.Transaction, err error)

	// Signs tx and sends it to the pending pool for execution
	// Returns the transaction hash of the submitted transaction
	SendTransaction(txRequest TransactionRequest) (txHash common.Hash, err error)

	/*
		DeployContract creates and submits a deployment transaction based on the
		deployer bytecode.
		It returns the address and creation transaction of the pending contract,
		or an error if the creation failed.

		It does not work on non-Ethernet compatible networks.
		The operation stops waiting when ctx is canceled.
	*/
	DeployContract(ctx context.Context, metaData *bind.MetaData) (common.Address, error)

	/*
		transact of Deployment Contract to tx pool.

		You can wait deployment using deployRes.

			deployRes, err := wallet.DeployContractNoWait(&<your-metadata>)
			tx := deployRes.Txs[metaData.ID]
			addr, err := alchemy.Transact.WaitDeployed(tx.Hash().Hex())
	*/
	DeployContractNoWait(metaData *bind.MetaData) (*bind.DeploymentResult, error)

	/*
		ContractTransact executes a transaction on a deployed contract.
		It waits for the transaction to be mined and returns the transaction receipt.
		The operation stops waiting when ctx is canceled.
	*/
	ContractTransact(
		ctx context.Context,
		contractAddress string,
		data []byte,
	) (*gethTypes.Receipt, error)

	/*
		ContractTransact executes a transaction on a deployed contract.

		You can wait deployment using deployRes.

			tx, err := wallet.ContractTransactNoWait(addr, data)
			txReceipt, err := alchemy.Transact.WaitDeployed(tx.Hash().Hex())
	*/
	ContractTransactNoWait(
		contractAddress string,
		data []byte,
	) (*gethTypes.Transaction, error)

	/*
		ContractCall calls a contract method.
		It is used for read-only methods.
	*/
	ContractCall(
		contractAddress string,
		opts *bind.CallOpts,
		callData []byte,
		unpack func([]byte) (any, error),
	) (any, error)

	/*
		EIP-712 signing by private key.

		EIP-712 is a standard for hashing and signing of typed structured data, as opposed to just arbitrary bytes.
		It is used to prevent signing of unintended data and to make the signed data more human-readable.

		refs: https://eips.ethereum.org/EIPS/eip-712
	*/
	SignEIP712(
		domainSeparator [32]byte, encoded []byte,
	) (Signature, error)

	/* ERC20 support */
	ERC20() WalletERC20

	/* StableCoin support */
	StableCoin() WalletStableCoin

	/* Nft (ERC721) support */
	Nft() WalletNft

	/* Erc1155 (multi-token) support */
	Erc1155() WalletErc1155

	/*
		ResetPool clears the cached ChainID and TransactOpts.
		Call this when you need to refresh the cached values.

		If switch network, you need to call this.
	*/
	ResetPool()
}
