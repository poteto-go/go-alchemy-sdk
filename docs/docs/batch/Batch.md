![](https://img.shields.io/badge/go-geth-lightblue)

Sends multiple read-only calls in a **single HTTP round-trip** while keeping the
SDK's typed, decoded ergonomics — callers never touch geth's raw `rpc.BatchElem`
or hex strings.

`batch.NewBatcher(ether)` returns a `*Batcher` whose typed sub-namespaces mirror
the Alchemy namespaces (`Core`, `ERC20`, `StableCoin`). Queue calls (each returns
a `*Result[T]`), call `Send` once, then read each result with `Unwrap`, which
returns the decoded native Go value.

- a `*Result[T]` value is available only after `Send`; `Unwrap` before then
  returns `constant.ErrBatchNotSent`.
- `Send`'s `error` is set only for I/O level failures (the whole batch). A
  per-request RPC error (or a rejected input such as an invalid address) is
  surfaced through that request's `Unwrap`.
- a batch sends once; a second `Send` returns `constant.ErrBatchAlreadySent`.

```go
func batch.NewBatcher(ether types.EtherApi) *batch.Batcher
func (b *Batcher) Send() error
func (r *Result[T]) Unwrap() (T, error)

// b.Core   — geth scalar reads
func (c *CoreBatch) BlockNumber() *Result[uint64]
func (c *CoreBatch) GasPrice() *Result[*big.Int]
func (c *CoreBatch) ChainID() *Result[*big.Int]
func (c *CoreBatch) PeerCount() *Result[uint64]
func (c *CoreBatch) Balance(address, blockTag string) *Result[*big.Int]
func (c *CoreBatch) Code(address, blockTag string) *Result[string]
func (c *CoreBatch) StorageAt(address, position, blockTag string) *Result[string]

// b.ERC20  — ERC-20 reads (BalanceOf/TotalSupply/Allowance/Name/Symbol/Decimals)
// b.StableCoin — FiatToken reads (Owner/Paused/IsMinter/Nonces/DomainSeparator/...)
```

```go
func main() {
	...
	alchemy := gas.NewAlchemy(setting)

	b := batch.NewBatcher(alchemy.GetProvider().Eth())

	blockNumber := b.Core.BlockNumber()                 // *Result[uint64]
	balance := b.Core.Balance("0xabc...", "latest")     // *Result[*big.Int]
	name := b.ERC20.Name(contract)                       // *Result[string]
	decimals := b.ERC20.Decimals(contract)               // *Result[uint8]
	owner := b.StableCoin.Owner(contract)                // *Result[common.Address]

	if err := b.Send(); err != nil {
		// I/O level failure (whole batch)
	}

	bn, _ := blockNumber.Unwrap()
	bal, _ := balance.Unwrap()
	nm, _ := name.Unwrap()
	o, _ := owner.Unwrap()
}
```

## Escape hatches

To batch a method the sub-namespaces don't expose, use the generic primitives:

- `batch.Add(b, method, args, target, convert)` — any JSON-RPC method.
- `batch.AddCall(b, contractAddress, signature, decode, args...)` — any contract
  `eth_call` read.
