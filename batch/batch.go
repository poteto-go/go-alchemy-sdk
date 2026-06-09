package batch

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/poteto-go/go-alchemy-sdk/constant"
	"github.com/poteto-go/go-alchemy-sdk/encode"
	"github.com/poteto-go/go-alchemy-sdk/types"
)

/*
Result is a typed handle to a single request queued on a Batcher.

The value is not available until the batch is sent with Batcher.Send. Calling
Unwrap before then returns constant.ErrBatchNotSent.
*/
type Result[T any] struct {
	value   T
	err     error
	settled bool
}

/*
Unwrap returns the decoded value and the per-request error.

It returns constant.ErrBatchNotSent if the owning batch has not been sent yet.
Once sent, err is the RPC error for this specific request (nil on success); the
value is the zero value of T when err is non-nil.
*/
func (r *Result[T]) Unwrap() (T, error) {
	if !r.settled {
		var zero T
		return zero, constant.ErrBatchNotSent
	}
	return r.value, r.err
}

// ethCallObject is the call object marshaled for an eth_call read.
type ethCallObject struct {
	To   string `json:"to"`
	Data string `json:"data"`
}

// failed returns an already-settled Result carrying err. It is used for inputs
// rejected before the call is queued (e.g. address validation), so the failure
// surfaces through Unwrap like any other per-request error.
func failed[T any](err error) *Result[T] {
	return &Result[T]{err: err, settled: true}
}

/*
Batcher collects multiple read-only calls and sends them in a single HTTP
round-trip while keeping the SDK's typed, decoded ergonomics.

Queue calls through the typed sub-namespaces, which mirror the Alchemy
namespaces, then call Send once and read each Result with Unwrap:

	b := batch.NewBatcher(alchemy.GetProvider().Eth())
	bn := b.Core.BlockNumber()
	name := b.ERC20.Name(contract)
	owner := b.StableCoin.Owner(contract)
	if err := b.Send(); err != nil { // I/O failure only
		return err
	}
	n, err := bn.Unwrap()
*/
type Batcher struct {
	ether      types.EtherApi
	elems      []rpc.BatchElem
	finalizers []func(rpc.BatchElem)
	sent       bool

	Core       *CoreBatch
	ERC20      *ERC20Batch
	StableCoin *StableCoinBatch
}

// NewBatcher creates a Batcher bound to the given EtherApi (e.g.
// alchemy.GetProvider().Eth()).
func NewBatcher(ether types.EtherApi) *Batcher {
	b := &Batcher{ether: ether}
	b.Core = &CoreBatch{b: b}
	b.ERC20 = &ERC20Batch{b: b}
	b.StableCoin = &StableCoinBatch{ERC20Batch: b.ERC20}
	return b
}

/*
Send dispatches every queued call in one HTTP round-trip and decodes each
result into its Result.

The returned error is only set for I/O level failures (the whole batch failed);
per-request RPC errors are surfaced through each Result.Unwrap. A batch can only
be sent once; a second call returns constant.ErrBatchAlreadySent.
*/
func (b *Batcher) Send() error {
	if b.sent {
		return constant.ErrBatchAlreadySent
	}
	b.sent = true

	if len(b.elems) == 0 {
		return nil
	}

	if err := b.ether.BatchCall(b.elems); err != nil {
		return err
	}

	for i, finalize := range b.finalizers {
		finalize(b.elems[i])
	}
	return nil
}

// addConv queues one call: geth unmarshals its raw result into target during
// Send, then convert turns that decoded value into T. It is the single place
// the elem + finalizer are wired, shared by Add and AddCall.
func addConv[D, T any](
	b *Batcher,
	method string,
	args []any,
	target *D,
	convert func(*D) (T, error),
) *Result[T] {
	res := &Result[T]{}

	b.elems = append(b.elems, rpc.BatchElem{
		Method: method,
		Args:   args,
		Result: target,
	})
	b.finalizers = append(b.finalizers, func(elem rpc.BatchElem) {
		res.settled = true
		if elem.Error != nil {
			res.err = elem.Error
			return
		}
		res.value, res.err = convert(target)
	})

	return res
}

/*
Add is the low-level typed escape hatch: it queues an arbitrary JSON-RPC method
whose raw result geth unmarshals into target during Send, then convert turns the
decoded value into T. Use it to batch a method the typed sub-namespaces do not
expose.
*/
func Add[D, T any](
	b *Batcher,
	method string,
	args []any,
	target *D,
	convert func(*D) T,
) *Result[T] {
	return addConv(b, method, args, target, func(d *D) (T, error) {
		return convert(d), nil
	})
}

/*
AddCall is the typed escape hatch for contract reads: it queues an eth_call (at
the latest block) and decodes the returned bytes with decode. The typed token
sub-namespaces (ERC20, StableCoin) are built on it.

signature is the function signature (e.g. []byte("balanceOf(address)")) and args
are the already ABI-encoded 32-byte words (see encode.ABIAddress).
*/
func AddCall[T any](
	b *Batcher,
	contractAddress string,
	signature []byte,
	decode func([]byte) (T, error),
	args ...[]byte,
) *Result[T] {
	calldata := encode.ReadCalldata(signature, args...)
	call := ethCallObject{
		To:   common.HexToAddress(contractAddress).Hex(),
		Data: hexutil.Encode(calldata),
	}
	target := new(hexutil.Bytes)

	return addConv(b, constant.Eth_Call, []any{call, "latest"}, target,
		func(t *hexutil.Bytes) (T, error) { return decode([]byte(*t)) })
}
