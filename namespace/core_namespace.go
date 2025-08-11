package namespace

import (
	"math/big"

	"github.com/poteto-go/go-alchemy-sdk/ether"
	"github.com/poteto-go/go-alchemy-sdk/types"
)

type ICore interface {
	/* get  the number of the most recent block. */
	GetBlockNumber() (int, error)

	/* Returns the best guess of the current gas price to use in a transaction. */
	GetGasPrice() (int, error)

	/* Returns the balance of a given address as of the provided block. */
	GetBalance(address string, blockTag string) (*big.Int, error)

	/*
		Returns the contract code of the provided address at the block.
		If there is no contract deployed, the result is 0x.
	*/
	GetCode(address, blockTag string) (string, error)

	/* Checks if the provided address is a smart contract. */
	IsContractAddress(address string) bool

	/*
		Returns the transaction with hash or null if the transaction is unknown.

		If a transaction has not been mined, this method will search the
		transaction pool. Various backends may have more restrictive transaction
		pool access (e.g. if the gas price is too low or the transaction was only
		recently sent and not yet indexed) in which case this method may also return null.

		NOTE: This is an alias for {@link TransactNamespace.getTransaction}.
	*/
	GetTransaction(hash string) (types.TransactionResponse, error)

	/*
		Return the value of the provided position at the provided address, at the provided block in `Bytes32` format.
		For inspecting solidity code.
	*/
	GetStorageAt(address, position, blockTag string) (string, error)

	/*
		Returns the ERC-20 token balances for a specific owner address
	*/
	GetTokenBalances(address string, option *types.TokenBalanceOption) (types.TokenBalanceResponse, error)

	/* Returns metadata for a given token contract address. */
	GetTokenMetadata(address string) (types.TokenMetadataResponse, error)

	/*
		Returns an array of logs that match the provided filter.
	*/
	GetLogs(filter types.Filter) ([]types.LogResponse, error)

	/*
		Returns an estimate of the amount of gas that would be required to submit transaction to the network.

		An estimate may not be accurate since there could be another transaction on the network that was not accounted for,
		but after being mined affects the relevant state.
		This is an alias for {@link TransactNamespace.estimateGas}.
	*/
	EstimateGas(tx types.TransactionRequest) (*big.Int, error)

	/*
		Returns the result of executing the transaction, using call.
		A call does not require any ether, but cannot change any state.
		This is useful for calling getters on Contracts.
	*/
	Call(tx types.TransactionRequest, blockTag string) (string, error)

	/*
		TODO: null if the tx has not been mined.
		Returns the transaction receipt for hash.
		To stall until the transaction has been mined, consider the waitForTransaction method below.
	*/
	GetTransactionReceipt(hash string) (types.TransactionReceipt, error)

	/*
		An enhanced API that gets all transaction receipts for a given block by number or block hash.
	*/
	GetTransactionReceipts(arg types.TransactionReceiptsArg) ([]types.TransactionReceipt, error)

	/*
		Simple wrapper around eth_getBlockByNumber.
		This returns the complete block information for the provided block number.
	*/
	GetBlockByBlockNumber(blockNumber string) (types.Block, error)
}

type Core struct {
	ether ether.EtherApi
}

func NewCore(ether ether.EtherApi) ICore {
	return &Core{
		ether: ether,
	}
}

func (c *Core) GetBlockNumber() (int, error) {
	blockNumber, err := c.ether.GetBlockNumber()
	if err != nil {
		return 0, err
	}
	return blockNumber, nil
}

func (c *Core) GetGasPrice() (int, error) {
	price, err := c.ether.GetGasPrice()
	if err != nil {
		return 0, err
	}
	return price, nil
}

func (c *Core) GetBalance(address string, blockTag string) (*big.Int, error) {
	balance, err := c.ether.GetBalance(address, blockTag)
	if err != nil {
		return big.NewInt(0), err
	}
	return balance, nil
}

func (c *Core) GetCode(address, blockTag string) (string, error) {
	hexCode, err := c.ether.GetCode(address, blockTag)
	if err != nil {
		return "", err
	}
	return hexCode, nil
}

/* Checks if the provided address is a smart contract. */
func (c *Core) IsContractAddress(address string) bool {
	hexCode, err := c.GetCode(address, "latest")
	if err != nil {
		return false
	}

	return hexCode != "0x"
}

func (c *Core) GetTransaction(hash string) (types.TransactionResponse, error) {
	transaction, err := c.ether.GetTransaction(hash)
	if err != nil {
		return types.TransactionResponse{}, err
	}

	return transaction, nil
}

func (c *Core) GetStorageAt(address, position, blockTag string) (string, error) {
	value, err := c.ether.GetStorageAt(address, position, blockTag)
	if err != nil {
		return "", err
	}

	return value, nil
}

func (c *Core) GetTokenBalances(address string, option *types.TokenBalanceOption) (types.TokenBalanceResponse, error) {
	params := []string{}
	if option != nil {
		params = option.ContractAddresses
	}

	result, err := c.ether.GetTokenBalances(address, params...)
	if err != nil {
		return types.TokenBalanceResponse{}, err
	}

	return result, nil
}

func (c *Core) GetTokenMetadata(address string) (types.TokenMetadataResponse, error) {
	result, err := c.ether.GetTokenMetadata(address)
	if err != nil {
		return types.TokenMetadataResponse{}, err
	}

	return result, nil
}

func (c *Core) GetLogs(filter types.Filter) ([]types.LogResponse, error) {
	logs, err := c.ether.GetLogs(filter)
	if err != nil {
		return []types.LogResponse{}, err
	}

	return logs, nil
}

func (c *Core) EstimateGas(tx types.TransactionRequest) (*big.Int, error) {
	estimatedGas, err := c.ether.EstimateGas(tx)
	if err != nil {
		return big.NewInt(0), err
	}

	return estimatedGas, nil
}

func (c *Core) Call(tx types.TransactionRequest, blockTag string) (string, error) {
	result, err := c.ether.Call(tx, blockTag)
	if err != nil {
		return "", err
	}

	return result, nil
}

func (c *Core) GetTransactionReceipt(hash string) (types.TransactionReceipt, error) {
	receipt, err := c.ether.GetTransactionReceipt(hash)
	if err != nil {
		return types.TransactionReceipt{}, err
	}

	return receipt, nil
}

func (c *Core) GetTransactionReceipts(arg types.TransactionReceiptsArg) ([]types.TransactionReceipt, error) {
	receipts, err := c.ether.GetTransactionReceipts(arg)
	if err != nil {
		return []types.TransactionReceipt{}, err
	}

	return receipts, nil
}

func (c *Core) GetBlockByBlockNumber(blockNumber string) (types.Block, error) {
	block, err := c.ether.GetBlockByNumber(blockNumber)
	if err != nil {
		return types.Block{}, err
	}

	return block, nil
}
