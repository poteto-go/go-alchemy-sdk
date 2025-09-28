package namespace

import (
	"math/big"

	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/poteto-go/go-alchemy-sdk/constant"
	"github.com/poteto-go/go-alchemy-sdk/ether"
	"github.com/poteto-go/go-alchemy-sdk/types"
)

type ICore interface {
	/* get  the number of the most recent block. */
	GetBlockNumber() (uint64, error)

	/* Returns the best guess of the current gas price to use in a transaction. */
	GetGasPrice() (*big.Int, error)

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
	GetTransaction(hash string) (tx *gethTypes.Transaction, isPending bool, err error)

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
		Null if the tx has not been mined.
		Returns the transaction receipt for hash.
		To stall until the transaction has been mined, consider the waitForTransaction method below.
	*/
	GetTransactionReceipt(hash string) (*gethTypes.Receipt, error)

	/*
		An enhanced API that gets all transaction receipts for a given block by number or block hash.
		Returns geth's Receipt.
	*/
	GetTransactionReceipts(arg types.TransactionReceiptsArg) ([]*gethTypes.Receipt, error)

	/*
		Returns the block from the network based on the provided block number or hash.
		Transactions on the block are represented as an array of transaction hashes.
		To get the full transaction details on the block, use {@link getBlockWithTransactions} instead.

		@param BlockHashOrBlockTag The block number or hash to get the block for.
	*/
	GetBlock(blockHashOrBlockTag types.BlockHashOrBlockTag) (*gethTypes.Block, error)
}

type Core struct {
	ether ether.EtherApi
}

func NewCore(ether ether.EtherApi) ICore {
	return &Core{
		ether: ether,
	}
}

func (c *Core) GetBlockNumber() (uint64, error) {
	blockNumber, err := c.ether.BlockNumber()
	if err != nil {
		return 0, err
	}
	return blockNumber, nil
}

func (c *Core) GetGasPrice() (*big.Int, error) {
	price, err := c.ether.GasPrice()
	if err != nil {
		return big.NewInt(0), err
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

func (c *Core) GetTransaction(hash string) (*gethTypes.Transaction, bool, error) {
	transaction, isPending, err := c.ether.GetTransaction(hash)
	if err != nil {
		return nil, false, err
	}

	return transaction, isPending, nil
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

func (c *Core) GetTransactionReceipt(hash string) (*gethTypes.Receipt, error) {
	receipt, err := c.ether.GetTransactionReceipt(hash)
	if err != nil {
		return nil, err
	}

	return receipt, nil
}

func (c *Core) GetTransactionReceipts(arg types.TransactionReceiptsArg) ([]*gethTypes.Receipt, error) {
	receipts, err := c.ether.GetTransactionReceipts(arg)
	if err != nil {
		return []*gethTypes.Receipt{}, err
	}

	return receipts, nil
}

func (c *Core) GetBlock(blockHashOrBlockTag types.BlockHashOrBlockTag) (*gethTypes.Block, error) {
	if blockHashOrBlockTag.BlockHash != "" {
		block, err := c.ether.GetBlockByHash(blockHashOrBlockTag.BlockHash)
		if err != nil {
			return nil, err
		}

		return block, nil
	}

	if blockHashOrBlockTag.BlockTag != "" {
		block, err := c.ether.GetBlockByNumber(blockHashOrBlockTag.BlockTag)
		if err != nil {
			return nil, err
		}

		return block, nil
	}

	return nil, constant.ErrInvalidArgs
}
