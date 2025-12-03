package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type ContractInstance interface {
	Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract
}

type EtherApi interface {
	GetEthClient() (*ethclient.Client, error)

	/* get  the number of the most recent block. */
	BlockNumber() (uint64, error)

	/* Returns the best guess of the current gas price to use in a transaction. */
	GasPrice() (*big.Int, error)

	/* Returns the balance of a given address as of the provided block. */
	GetBalance(address string, blockTag string) (*big.Int, error)

	/*
		StorageAt returns the value of key in the contract storage of the given account.
		The block number can be nil, in which case the value is taken from the latest known block.
	*/
	CodeAt(address string, blockTag string) (string, error)

	/*
		CodeAtHash returns the contract code of the given account.
	*/
	CodeAtHash(address string, blockHash string) (string, error)

	/*
		Returns the transaction with hash or null if the transaction is unknown.

		If a transaction has not been mined, this method will search the
		transaction pool. Various backends may have more restrictive transaction
		pool access (e.g. if the gas price is too low or the transaction was only
		recently sent and not yet indexed) in which case this method may also return null.

		internal call geth.TransactionByHash

		NOTE: This is an alias for {@link TransactNamespace.getTransaction}.
	*/
	GetTransaction(hash string) (tx *gethTypes.Transaction, isPending bool, err error)

	/*
		internal call geth ethclient.StorageAt
	*/
	StorageAt(address, position, blockTag string) (string, error)

	/*
		Returns the ERC-20 token balances for a specific owner address w, w/o params
	*/
	GetTokenBalances(address string, params ...string) (TokenBalanceResponse, error)

	/* Returns metadata for a given token contract address. */
	GetTokenMetadata(address string) (TokenMetadataResponse, error)

	/*
		Returns an array of logs that match the provided filter.
	*/
	GetLogs(filter Filter) ([]LogResponse, error)

	/*
		Returns an estimate of the amount of gas that would be required to submit transaction to the network.

		An estimate may not be accurate since there could be another transaction on the network that was not accounted for,
		but after being mined affects the relevant state.
		This is an alias for {@link TransactNamespace.estimateGas}.
	*/
	EstimateGas(tx TransactionRequest) (*big.Int, error)

	/*
		SuggestGasPrice retrieves the currently suggested gas price to allow a timely
		execution of a transaction.
	*/
	SuggestGasPrice() (*big.Int, error)

	/*
		Returns the result of executing the transaction, using call.
		A call does not require any ether, but cannot change any state.
		This is useful for calling getters on Contracts.
	*/
	Call(tx TransactionRequest, blockTag string) (string, error)

	/*
		Null if the tx has not been mined.
		Returns the transaction receipt for hash.
		To stall until the transaction has been mined, consider the waitForTransaction method below.
	*/
	GetTransactionReceipt(hash string) (*gethTypes.Receipt, error)

	/*
		An enhanced API that gets all transaction receipts for a given block by number or block hash.
		Returns geth's Receipt.
	*/
	GetTransactionReceipts(arg BlockNumberOrHash) ([]*gethTypes.Receipt, error)

	/*
		Simple wrapper around eth_getBlockByNumber.
		This returns the complete block information for the provided block number.
	*/
	GetBlockByNumber(blockNumber string) (*gethTypes.Block, error)

	/*
		Simple wrapper around eth_getBlockByHash.
		This returns the complete block information for the provided block hash.
	*/
	GetBlockByHash(blockHash string) (*gethTypes.Block, error)

	/*
		PendingNonceAt returns the account nonce of the given account in the pending state.
		This is the nonce that should be used for the next transaction.

		internal call geth
	*/
	PendingNonceAt(address string) (uint64, error)

	// send signed tx into the pending pool for execution w/geth
	SendRawTransaction(signedTx *gethTypes.Transaction) error

	/*
		ChainID retrieves the current chain ID for transaction replay protection.

		internal call geth
	*/
	ChainID() (*big.Int, error)

	/*
		PeerCount returns the number of p2p peers as reported by the net_peerCount method.

		internal call geth
	*/
	PeerCount() (uint64, error)

	/*
		WaitDeployed waits for a contract deployment transaction with the provided hash and
		returns the on-chain contract address when it is mined.
		It stops waiting when ctx is canceled.

		internal call geth
	*/
	DeployContract(
		auth *bind.TransactOpts,
		metaData *bind.MetaData,
	) (common.Address, error)

	/*
		ContractTransact transacts with a contract.
		Wait for mined
	*/
	ContractTransact(
		auth *bind.TransactOpts,
		contract ContractInstance,
		contractAddress string,
		data []byte,
	) (txReceipt *gethTypes.Receipt, err error)

	/*
		WaitMined waits for a transaction with the provided hash and
		returns the transaction receipt when it is mined.
		It stops waiting when ctx is canceled.
	*/
	WaitMined(hash common.Hash) (*gethTypes.Receipt, error)

	/*
		WaitDeployed waits for a contract deployment transaction with the provided hash and
		returns the contract address
		It stops waiting when ctx is canceled.
	*/
	WaitDeployed(hash common.Hash) (common.Address, error)
}
