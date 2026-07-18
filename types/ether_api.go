package types

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
)

// EthClient is the go-ethereum client surface that Ether depends on, narrowed to
// the subset that BOTH *ethclient.Client (URL/RPC backed) and simulated.Client
// (in-process simulated.Backend) implement. It lets Ether hold a single client
// value without caring which transport it talks to.
//
// Methods exposed only by *ethclient.Client are intentionally excluded because
// simulated.Client (an interface that deliberately hides the raw client) does
// not provide them:
//   - Client() *rpc.Client (raw rpc access used by BatchCall)
//   - CodeAtHash
//   - PeerCount
//
// bind.ContractBackend / bind.DeployBackend are embedded so an EthClient value
// can be passed straight to bind.LinkAndDeploy / WaitMined / WaitDeployed.
type EthClient interface {
	bind.ContractBackend
	bind.DeployBackend

	BlockNumber(ctx context.Context) (uint64, error)
	ChainID(ctx context.Context) (*big.Int, error)
	StorageAt(ctx context.Context, account common.Address, key common.Hash, blockNumber *big.Int) ([]byte, error)
	BlockByNumber(ctx context.Context, number *big.Int) (*gethTypes.Block, error)
	BlockByHash(ctx context.Context, hash common.Hash) (*gethTypes.Block, error)
	BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error)
	TransactionByHash(ctx context.Context, txHash common.Hash) (tx *gethTypes.Transaction, isPending bool, err error)
}

type EtherApi interface {
	/*
		set ether client if client is nil,
		if client exists add connCount to re-use.

		In WS, we persist the client.
	*/
	SetEthClient() error

	/*
		decrement connCount.
		if connCount <= 0, close client & set nil

		The design of WS prevents it from being closed using `defer Close`.
	*/
	Close()

	/*
		shutdown ws client or kill http client
	*/
	Shutdown()

	/*
		get raw ethclient
	*/
	Client() EthClient

	/*
		Commit backend;
		mined transaction
		! this only works on simulated backend
	*/
	Commit() (common.Hash, error)

	/*
		Fork backend;
		you can revert to the commit hash point
		! this only works on simulated backend
	*/
	Fork(snapShotHash common.Hash) error

	/*
		BatchCall sends multiple JSON-RPC requests in a single HTTP round-trip
		using geth's underlying rpc.Client.

		Each element's Result/Error is populated in place (geth semantics): a
		per-request RPC error is stored on the element's Error field, while the
		returned error is only set for I/O level failures.
	*/
	BatchCall(elems []rpc.BatchElem) error

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
		SuggestGasTipCap returns the suggested maxPriorityFeePerGas (eth_maxPriorityFeePerGas)
		for EIP-1559 transactions.
	*/
	SuggestGasTipCap() (*big.Int, error)

	/*
		SuggestEIP1559Fees returns (maxPriorityFeePerGas, maxFeePerGas) ready to use in a
		TransactionRequest. maxFeePerGas is derived as baseFee*2 + maxPriorityFeePerGas.
		Returns an error on chains that do not support EIP-1559.
	*/
	SuggestEIP1559Fees() (maxPriorityFeePerGas *big.Int, maxFeePerGas *big.Int, err error)

	/*
		Read method call for Any Smart Contract
	*/
	CallReadMethod(
		method []byte,
		contractAddress string,
		args ...[]byte,
	) ([]byte, error)

	/*
		Returns the result of executing the transaction, using call.
		A call does not require any ether, but cannot change any state.
		This is useful for calling getters on Contracts.
	*/
	Call(tx TransactionRequest, blockTag string) (string, error)

	/*
		Return the result of eth_call of smart contract by provided call message.
	*/
	CallContract(msg ethereum.CallMsg, blockTag string) ([]byte, error)

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
		Deploy Contract to tx pool.

		internal call geth
	*/
	DeployContract(
		auth *bind.TransactOpts,
		metaData *bind.MetaData,
	) (*bind.DeploymentResult, error)

	/*
		ContractTransact transacts with a contract.

		internal call geth
	*/
	ContractTransact(
		auth *bind.TransactOpts,
		contractAddress string,
		data []byte,
	) (txReceipt *gethTypes.Transaction, err error)

	/*
		WaitMined waits for a transaction with the provided hash and
		returns the transaction receipt when it is mined.
		It stops waiting when ctx is canceled.
	*/
	WaitMined(ctx context.Context, hash common.Hash) (*gethTypes.Receipt, error)

	/*
		WaitDeployed waits for a contract deployment transaction with the provided hash and
		returns the contract address.
		It stops waiting when ctx is canceled.
	*/
	WaitDeployed(ctx context.Context, hash common.Hash) (common.Address, error)

	/*
		Snapshot takes a snapshot of the current blockchain state with evm_snapshot
		and returns the snapshot id.

		Only supported on development chains (hardhat, anvil, ganache, ...).
	*/
	Snapshot() (*big.Int, error)

	/*
		RevertTo reverts the blockchain state to the provided snapshot id with evm_revert.
		It returns true if the state was reverted.

		Only supported on development chains (hardhat, anvil, ganache, ...).
	*/
	RevertTo(snapshotId *big.Int) (bool, error)

	/*
		Network returns the Alchemy network this client is connected to.
		Returns an empty string for simulated backends.
	*/
	Network() Network

	/*
		ResolveNameBy resolves an ENS name to a lowercase hex address using the
		provided ENS registry contract address.
		If name is already a valid hex address it is returned as-is (lowercased).
	*/
	ResolveNameBy(registryAddress string, name string) (string, error)

	/*
		LookupAddressBy performs a reverse ENS lookup (address → name) using the
		provided ENS registry contract address.
		Returns an error when no reverse record is registered.
	*/
	LookupAddressBy(registryAddress string, address string) (string, error)

	/*
		ContractCall calls a contract.

		internal call geth
	*/
	ContractCall(
		contractAddress common.Address,
		ops *bind.CallOpts,
		callData []byte,
		unpack func([]byte) (any, error),
	) (any, error)

	/*
		GetAssetTransfers fetches asset transfer history matching the given params.

		NOTE: This is an Alchemy-specific API (alchemy_getAssetTransfers).
		It is not available on non-Alchemy endpoints such as simulated backends.
	*/
	GetAssetTransfers(params AssetTransfersParams) (AssetTransfersResponse, error)

	// ! WsEtherApi is the interface for Ether's websocket provider.
	// to use this, you need set UseWebsocket: true
	//
	//  setting := gas.AlchemySetting{
	//    ApiKey:       "<alchemy-api-key>",
	//    Network:      types.EthSepolia,
	//    UseWebsocket: true,
	//  }
	WsEtherApi
}

type WsEtherApi interface {
	Subscribe(ctx context.Context, channel any, params ...any) (ethereum.Subscription, error)
	SubscribeNewHead(ctx context.Context, headerChan chan<- *gethTypes.Header) (ethereum.Subscription, error)
	SubscribeFilterLogs(ctx context.Context, query ethereum.FilterQuery, logChan chan<- gethTypes.Log) (ethereum.Subscription, error)
	SubscribeTxReceipts(ctx context.Context, q *ethereum.TransactionReceiptsQuery, receiptsChan chan<- []*gethTypes.Receipt) (ethereum.Subscription, error)
}
