---
sidebar_position: 13
---

# StableCoin

You can execute ERC20-compatible methods on StableCoin contracts (e.g. USDC) from the specified wallet.

:::note

Since it performs calls via bytecode, it does not require the contract implementation and can be called as long as the address is available.

:::

:::warning

- It requires connected wallet.
- It does not work on non-Ethereum compatible networks.

:::

```go
func main() {
	alchemy = gas.NewAlchemy(setting)
	w, _ := wallet.New("<privateKey>")
	w.Connect(alchemy.GetProvider())

	// call each method
	w.StableCoin().MethodXXX()
}
```

## BalanceOf

Get the StableCoin token balance of the connected wallet address.

```go
func BalanceOf(contractAddress string) (balance *big.Int, err error)
```

```go
func main() {
	contractAddress := "0x1234567890123456789012345678901234567890"
	balance, err := w.StableCoin().BalanceOf(contractAddress)
}
```

## Read Methods

You can fetch token metadata and balance information:

- TotalSupply
- Allowance
- Name
- Symbol
- Decimals
- BalanceOf

## Write Methods

### Transfer & TransferNoWait

Transfer StableCoin token from wallet.

```go
receipt, err := w.StableCoin().Transfer(
	context.Background(),
	contractAddress,
	"<toAddress>",
	big.NewInt(100),
	nil,
)
```

### Approve & ApproveNoWait

Approve a spender to spend tokens on behalf of the connected wallet.

```go
receipt, err := w.StableCoin().Approve(
	context.Background(),
	contractAddress,
	"<spenderAddress>",
	big.NewInt(100),
	nil,
)
```

### TransferFrom & TransferFromNoWait

Transfer tokens from another address using a prior allowance.

```go
receipt, err := w.StableCoin().TransferFrom(
	context.Background(),
	contractAddress,
	"<fromAddress>",
	"<toAddress>",
	big.NewInt(100),
	nil,
)
```

### Mint & MintNoWait

Mint StableCoin tokens to an address. Requires the caller to have the minter role (e.g. configured via `configureMinter` on FiatToken/USDC).

```go
func Mint(ctx context.Context, contractAddress, toAddress string, amount *big.Int, gasLimit *uint64) (*types.Receipt, error)
func MintNoWait(contractAddress, toAddress string, amount *big.Int, gasLimit *uint64) (common.Hash, error)
```

```go
receipt, err := w.StableCoin().Mint(
	context.Background(),
	contractAddress,
	"<toAddress>",
	big.NewInt(100),
	nil,
)
```

### Burn & BurnNoWait

Burn StableCoin tokens from the caller's own balance. Requires the caller to have the minter role (FiatToken/USDC behavior).

```go
func Burn(ctx context.Context, contractAddress string, amount *big.Int, gasLimit *uint64) (*types.Receipt, error)
func BurnNoWait(contractAddress string, amount *big.Int, gasLimit *uint64) (common.Hash, error)
```

```go
receipt, err := w.StableCoin().Burn(
	context.Background(),
	contractAddress,
	big.NewInt(100),
	nil,
)
```
