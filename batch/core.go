package batch

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/poteto-go/go-alchemy-sdk/constant"
)

// CoreBatch queues geth-backed scalar reads onto its Batcher. Mirrors the
// read-only methods of the Core namespace.
type CoreBatch struct {
	b *Batcher
}

/* Queues eth_blockNumber. */
func (c *CoreBatch) BlockNumber() *Result[uint64] {
	return Add(c.b, constant.Eth_BlockNumber, nil, new(hexutil.Uint64), fromHexUint64)
}

/* Queues eth_gasPrice. */
func (c *CoreBatch) GasPrice() *Result[*big.Int] {
	return Add(c.b, constant.Eth_GasPrice, nil, new(hexutil.Big), fromHexBig)
}

/* Queues eth_chainId. */
func (c *CoreBatch) ChainID() *Result[*big.Int] {
	return Add(c.b, constant.Eth_ChainId, nil, new(hexutil.Big), fromHexBig)
}

/* Queues net_peerCount. */
func (c *CoreBatch) PeerCount() *Result[uint64] {
	return Add(c.b, constant.Net_PeerCount, nil, new(hexutil.Uint64), fromHexUint64)
}

/* Queues eth_getBalance for the given address at the block tag. */
func (c *CoreBatch) Balance(address, blockTag string) *Result[*big.Int] {
	return Add(c.b, constant.Eth_GetBalance, []any{strings.ToLower(address), blockTag}, new(hexutil.Big), fromHexBig)
}

/* Queues eth_getCode for the given address at the block tag. */
func (c *CoreBatch) Code(address, blockTag string) *Result[string] {
	return Add(c.b, constant.Eth_GetCode, []any{address, blockTag}, new(hexutil.Bytes), fromHexBytesEncoded)
}

/* Queues eth_getStorageAt for the given address and position at the block tag. */
func (c *CoreBatch) StorageAt(address, position, blockTag string) *Result[string] {
	// matches ether.StorageAt output (unprefixed hex).
	return Add(c.b, constant.Eth_GetStorageAt, []any{address, position, blockTag}, new(hexutil.Bytes), fromHexBytesRaw)
}

// --- decoded-value converters ------------------------------------------------

func fromHexUint64(v *hexutil.Uint64) uint64 { return uint64(*v) }

func fromHexBig(v *hexutil.Big) *big.Int { return (*big.Int)(v) }

func fromHexBytesEncoded(v *hexutil.Bytes) string { return hexutil.Encode(*v) }

func fromHexBytesRaw(v *hexutil.Bytes) string { return common.Bytes2Hex(*v) }
