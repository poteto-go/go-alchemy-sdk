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
- Currency
- Version
- MasterMinter
- Pauser
- Blacklister
- Owner
- IsMinter
- MinterAllowance

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

### Blacklist & BlacklistNoWait

Add an address to the blacklist. Requires the caller to have the blacklister role (FiatToken/USDC behavior).

```go
func Blacklist(ctx context.Context, contractAddress, address string, gasLimit *uint64) (*types.Receipt, error)
func BlacklistNoWait(contractAddress, address string, gasLimit *uint64) (common.Hash, error)
```

```go
receipt, err := w.StableCoin().Blacklist(
	context.Background(),
	contractAddress,
	"<targetAddress>",
	nil,
)
```

### UnBlacklist & UnBlacklistNoWait

Remove an address from the blacklist. Requires the caller to have the blacklister role (FiatToken/USDC behavior).

```go
func UnBlacklist(ctx context.Context, contractAddress, address string, gasLimit *uint64) (*types.Receipt, error)
func UnBlacklistNoWait(contractAddress, address string, gasLimit *uint64) (common.Hash, error)
```

```go
receipt, err := w.StableCoin().UnBlacklist(
	context.Background(),
	contractAddress,
	"<targetAddress>",
	nil,
)
```

### IsBlacklisted

Check whether an address is blacklisted on the contract.

```go
func IsBlacklisted(contractAddress, address string) (bool, error)
```

```go
isBlacklisted, err := w.StableCoin().IsBlacklisted(contractAddress, "<targetAddress>")
```

### Pause & PauseNoWait

Pause all token transfers. Requires the caller to have the pauser role (FiatToken/USDC behavior).

```go
func Pause(ctx context.Context, contractAddress string, gasLimit *uint64) (*types.Receipt, error)
func PauseNoWait(contractAddress string, gasLimit *uint64) (common.Hash, error)
```

```go
receipt, err := w.StableCoin().Pause(context.Background(), contractAddress, nil)
```

### Unpause & UnpauseNoWait

Resume token transfers. Requires the caller to have the pauser role (FiatToken/USDC behavior).

```go
func Unpause(ctx context.Context, contractAddress string, gasLimit *uint64) (*types.Receipt, error)
func UnpauseNoWait(contractAddress string, gasLimit *uint64) (common.Hash, error)
```

```go
receipt, err := w.StableCoin().Unpause(context.Background(), contractAddress, nil)
```

### Paused

Check whether the contract is currently paused.

```go
func Paused(contractAddress string) (bool, error)
```

```go
paused, err := w.StableCoin().Paused(contractAddress)
```

### UpdateMasterMinter & UpdateMasterMinterNoWait

Update the master minter address. Requires the caller to be the current owner (FiatToken/USDC behavior).

```go
func UpdateMasterMinter(ctx context.Context, contractAddress, newMasterMinter string, gasLimit *uint64) (*types.Receipt, error)
func UpdateMasterMinterNoWait(contractAddress, newMasterMinter string, gasLimit *uint64) (common.Hash, error)
```

```go
receipt, err := w.StableCoin().UpdateMasterMinter(
	context.Background(),
	contractAddress,
	"<newMasterMinterAddress>",
	nil,
)
```

### UpdateBlacklister & UpdateBlacklisterNoWait

Update the blacklister address. Requires the caller to be the current owner (FiatToken/USDC behavior).

```go
func UpdateBlacklister(ctx context.Context, contractAddress, newBlacklister string, gasLimit *uint64) (*types.Receipt, error)
func UpdateBlacklisterNoWait(contractAddress, newBlacklister string, gasLimit *uint64) (common.Hash, error)
```

```go
receipt, err := w.StableCoin().UpdateBlacklister(
	context.Background(),
	contractAddress,
	"<newBlacklisterAddress>",
	nil,
)
```

### UpdatePauser & UpdatePauserNoWait

Update the pauser address. Requires the caller to be the current owner (FiatToken/USDC behavior).

```go
func UpdatePauser(ctx context.Context, contractAddress, newPauser string, gasLimit *uint64) (*types.Receipt, error)
func UpdatePauserNoWait(contractAddress, newPauser string, gasLimit *uint64) (common.Hash, error)
```

```go
receipt, err := w.StableCoin().UpdatePauser(
	context.Background(),
	contractAddress,
	"<newPauserAddress>",
	nil,
)
```

### TransferOwnership & TransferOwnershipNoWait

Transfer contract ownership to a new address. Requires the caller to be the current owner (FiatToken/USDC behavior).

```go
func TransferOwnership(ctx context.Context, contractAddress, newOwner string, gasLimit *uint64) (*types.Receipt, error)
func TransferOwnershipNoWait(contractAddress, newOwner string, gasLimit *uint64) (common.Hash, error)
```

```go
receipt, err := w.StableCoin().TransferOwnership(
	context.Background(),
	contractAddress,
	"<newOwnerAddress>",
	nil,
)
```

### MasterMinter

Get the master minter address of the contract. FiatToken/USDC compatibility.

```go
func MasterMinter(contractAddress string) (common.Address, error)
```

```go
masterMinter, err := w.StableCoin().MasterMinter(contractAddress)
```

### Pauser

Get the pauser address of the contract. FiatToken/USDC compatibility.

```go
func Pauser(contractAddress string) (common.Address, error)
```

```go
pauser, err := w.StableCoin().Pauser(contractAddress)
```

### Blacklister

Get the blacklister address of the contract. FiatToken/USDC compatibility.

```go
func Blacklister(contractAddress string) (common.Address, error)
```

```go
blacklister, err := w.StableCoin().Blacklister(contractAddress)
```

### Owner

Read the current owner address of the contract.

```go
func Owner(contractAddress string) (common.Address, error)
```

```go
owner, err := w.StableCoin().Owner(contractAddress)
```

### Currency

Get the currency identifier of the token (e.g. `"USD"`). FiatToken/USDC compatibility.

```go
func Currency(contractAddress string) (string, error)
```

```go
currency, err := w.StableCoin().Currency(contractAddress)
```

### Version

Get the contract version string. FiatToken/USDC compatibility.

```go
func Version(contractAddress string) (string, error)
```

```go
version, err := w.StableCoin().Version(contractAddress)
```

### ConfigureMinter & ConfigureMinterNoWait

Configure a minter with a mint allowance. Requires the caller to have the masterMinter role (FiatToken/USDC behavior).

```go
func ConfigureMinter(ctx context.Context, contractAddress, minter string, allowance *big.Int, gasLimit *uint64) (*types.Receipt, error)
func ConfigureMinterNoWait(contractAddress, minter string, allowance *big.Int, gasLimit *uint64) (common.Hash, error)
```

```go
receipt, err := w.StableCoin().ConfigureMinter(
	context.Background(),
	contractAddress,
	"<minterAddress>",
	big.NewInt(1_000_000),
	nil,
)
```

### RemoveMinter & RemoveMinterNoWait

Remove a minter, revoking their ability to mint. Requires the caller to have the masterMinter role (FiatToken/USDC behavior).

```go
func RemoveMinter(ctx context.Context, contractAddress, minter string, gasLimit *uint64) (*types.Receipt, error)
func RemoveMinterNoWait(contractAddress, minter string, gasLimit *uint64) (common.Hash, error)
```

```go
receipt, err := w.StableCoin().RemoveMinter(
	context.Background(),
	contractAddress,
	"<minterAddress>",
	nil,
)
```

### IsMinter

Check whether an address is a configured minter on the contract. FiatToken/USDC compatibility.

```go
func IsMinter(contractAddress, address string) (bool, error)
```

```go
isMinter, err := w.StableCoin().IsMinter(contractAddress, "<minterAddress>")
```

### MinterAllowance

Get the remaining mint allowance for a minter. FiatToken/USDC compatibility.

```go
func MinterAllowance(contractAddress, address string) (*big.Int, error)
```

```go
allowance, err := w.StableCoin().MinterAllowance(contractAddress, "<minterAddress>")
```

## EIP-2612 Methods

### Permit & PermitNoWait

Submit an EIP-2612 permit transaction, allowing a spender to use tokens on behalf of an owner via a pre-signed signature.

```go
func Permit(ctx context.Context, contractAddress, ownerAddress, spenderAddress string, value, deadline *big.Int, gasLimit *uint64) (*types.Receipt, error)
func PermitNoWait(contractAddress, ownerAddress, spenderAddress string, value, deadline *big.Int, gasLimit *uint64) (common.Hash, error)
```

The wallet automatically fetches the on-chain nonce and domain separator, then signs the EIP-712 permit message using the wallet's private key.

```go
receipt, err := w.StableCoin().Permit(
	context.Background(),
	contractAddress,
	"<ownerAddress>",
	"<spenderAddress>",
	big.NewInt(100),
	big.NewInt(deadline),
	nil,
)
```

## EIP-3009 Methods

### TransferWithAuthorization & TransferWithAuthorizationNoWait

Submit an EIP-3009 transfer-with-authorization transaction. The wallet signs a `TransferWithAuthorization` EIP-712 message and submits it on-chain. The `nonce` is a random `[32]byte` chosen by the caller.

```go
func TransferWithAuthorization(ctx context.Context, contractAddress, from, to string, value, validAfter, validBefore *big.Int, nonce [32]byte, gasLimit *uint64) (*types.Receipt, error)
func TransferWithAuthorizationNoWait(contractAddress, from, to string, value, validAfter, validBefore *big.Int, nonce [32]byte, gasLimit *uint64) (common.Hash, error)
```

Use `utils.NewAuthorizationNonce()` to generate a cryptographically random nonce. Hold on to the returned value if you may need to cancel the authorization later.

```go
nonce := utils.NewAuthorizationNonce()

receipt, err := w.StableCoin().TransferWithAuthorization(
	context.Background(),
	contractAddress,
	"<fromAddress>",
	"<toAddress>",
	big.NewInt(100),
	big.NewInt(validAfter),
	big.NewInt(validBefore),
	nonce,
	nil,
)
```

### ReceiveWithAuthorization & ReceiveWithAuthorizationNoWait

Submit an EIP-3009 receive-with-authorization transaction. The caller must be the `to` address. The wallet signs a `ReceiveWithAuthorization` EIP-712 message.

```go
func ReceiveWithAuthorization(ctx context.Context, contractAddress, from, to string, value, validAfter, validBefore *big.Int, nonce [32]byte, gasLimit *uint64) (*types.Receipt, error)
func ReceiveWithAuthorizationNoWait(contractAddress, from, to string, value, validAfter, validBefore *big.Int, nonce [32]byte, gasLimit *uint64) (common.Hash, error)
```

```go
nonce := utils.NewAuthorizationNonce()

receipt, err := w.StableCoin().ReceiveWithAuthorization(
	context.Background(),
	contractAddress,
	"<fromAddress>",
	"<toAddress>",
	big.NewInt(100),
	big.NewInt(validAfter),
	big.NewInt(validBefore),
	nonce,
	nil,
)
```

### CancelAuthorization & CancelAuthorizationNoWait

Cancel an outstanding EIP-3009 authorization by its nonce. The wallet signs a `CancelAuthorization` EIP-712 message.

```go
func CancelAuthorization(ctx context.Context, contractAddress, authorizer string, nonce [32]byte, gasLimit *uint64) (*types.Receipt, error)
func CancelAuthorizationNoWait(contractAddress, authorizer string, nonce [32]byte, gasLimit *uint64) (common.Hash, error)
```

```go
// nonce must be the same value used when the authorization was submitted
receipt, err := w.StableCoin().CancelAuthorization(
	context.Background(),
	contractAddress,
	"<authorizerAddress>",
	nonce,
	nil,
)
```
