---
sidebar_position: 4
---

# JPYC (Japanese Yen Stablecoin)

## Overview

[JPYC](https://jpyc.co.jp/) is a Japanese-yen–pegged stablecoin issued by JPYC Inc. as a
regulated *electronic payment instrument* (資金移動業 / funds-transfer service). It is built on
the same `FiatToken` base as USDC, so the SDK's `StableCoin` namespace works against it directly,
and it supports **EIP-3009 (`transferWithAuthorization`)** — the gasless-transfer primitive used by
x402-style payment flows.

Two things to know before you start:

- **JPYC has `18` decimals** — *not* `6` like USDC / USDT. This is the most common mistake when
  formatting balances.
- **Use the new (資金移動業) JPYC contract.** The current token deploys to the **same address on
  every chain**: `0xE7C3D8C9a439feDe00D2600032D5dB0Be71C3c29` (Ethereum, Polygon, Avalanche).
  This differs from the older JPYC v1 (prepaid) token. The SDK's `famous` registry returns it via
  `famous.ContractAddress(network, famous.JPYC)` — see [Notes](#notes).

This tutorial uses **Polygon mainnet via Alchemy** (read paths need no private key). A gasless
`transferWithAuthorization` example using the `SimulatedBackend` is sketched at the end.

## Setup

```go
package main

import (
	"fmt"
	"log"
	"math/big"

	"github.com/poteto-go/go-alchemy-sdk/famous"
	"github.com/poteto-go/go-alchemy-sdk/gas"
	"github.com/poteto-go/go-alchemy-sdk/types"
)

func main() {
	setting := gas.AlchemySetting{
		ApiKey:  "<alchemy-api-key>",
		Network: types.PolygonMainnet,
	}
	alchemy, err := gas.NewAlchemy(setting)
	if err != nil {
		log.Fatal(err)
	}

	// Resolve the JPYC address from the famous registry (新 資金移動業 JPYC).
	jpycAddr, err := famous.ContractAddress(types.PolygonMainnet, famous.JPYC)
	if err != nil {
		log.Fatal(err)
	}
	jpyc := jpycAddr.Hex()

	// ... see sections below
	_ = alchemy
	_ = jpyc
}
```

## Reading JPYC state

The `StableCoin` namespace exposes ERC20 + FiatToken read methods. No wallet/private key is
required for reads.

```go
holder := "0x0000000000000000000000000000000000000000" // any JPYC holder

// Balance — remember JPYC is 18 decimals.
raw, err := alchemy.StableCoin.BalanceOf(jpyc, holder)
if err != nil {
	log.Fatal(err)
}
fmt.Printf("balance: %s JPYC\n", formatUnits(raw, 18))

// FiatToken metadata
currency, _ := alchemy.StableCoin.Currency(jpyc) // "JPY"
version, _ := alchemy.StableCoin.Version(jpyc)
fmt.Printf("currency=%s version=%s\n", currency, version)
```

`formatUnits` converts the raw `*big.Int` (wei-style integer) to a human-readable amount:

```go
func formatUnits(raw *big.Int, decimals int) string {
	denom := new(big.Float).SetInt(new(big.Int).Exp(
		big.NewInt(10), big.NewInt(int64(decimals)), nil,
	))
	return new(big.Float).Quo(new(big.Float).SetInt(raw), denom).Text('f', 2)
}
```

## Fetching Transfer events

JPYC `Transfer` logs are retrievable with `Core.GetLogs`. The `Transfer(address,address,uint256)`
event topic is a constant:

```go
// keccak256("Transfer(address,address,uint256)")
const transferTopic = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"

logs, err := alchemy.Core.GetLogs(types.Filter{
	FromBlock: "0x...",  // e.g. a recent block in hex
	ToBlock:   "latest",
	Address:   jpyc,
	Topics:    []string{transferTopic},
})
if err != nil {
	log.Fatal(err)
}

for _, l := range logs {
	// l.Topics[1] = from (32-byte padded), l.Topics[2] = to, l.Data = value (uint256)
	value, _ := new(big.Int).SetString(strings.TrimPrefix(l.Data, "0x"), 16)
	fmt.Printf("from=%s to=%s value=%s JPYC\n",
		topicToAddress(l.Topics[1]),
		topicToAddress(l.Topics[2]),
		formatUnits(value, 18),
	)
}
```

```go
// a 32-byte topic holds a left-padded 20-byte address
func topicToAddress(topic string) string {
	h := strings.TrimPrefix(topic, "0x")
	return "0x" + h[len(h)-40:]
}
```

:::warning Alchemy block-range limit

`eth_getLogs` on Alchemy caps how wide a `FromBlock`→`ToBlock` range you can request in one call
(and how many results it returns). For a busy token like JPYC, page through bounded windows
(e.g. a few thousand blocks per call) instead of scanning from genesis, or the call will be
rejected / truncated.

:::

## Gasless transfers with EIP-3009

Because JPYC is a `FiatToken`, a holder can authorize a transfer with an **off-chain EIP-712
signature** and let a relayer pay the gas. This is the primitive x402 / agentic payments rely on.
The SDK ships helpers for the whole flow:

```go
package main

import (
	"context"
	"math/big"
	"time"

	"github.com/poteto-go/go-alchemy-sdk/constant"
	"github.com/poteto-go/go-alchemy-sdk/famous"
	"github.com/poteto-go/go-alchemy-sdk/gas"
	"github.com/poteto-go/go-alchemy-sdk/typeddata"
	"github.com/poteto-go/go-alchemy-sdk/types"
	"github.com/poteto-go/go-alchemy-sdk/utils"
	"github.com/poteto-go/go-alchemy-sdk/wallet"
)

func main() {
	alchemy, _ := gas.NewAlchemy(gas.AlchemySetting{
		ApiKey:  "<alchemy-api-key>",
		Network: types.PolygonMainnet,
	})
	jpycAddr, _ := famous.ContractAddress(types.PolygonMainnet, famous.JPYC)
	jpyc := jpycAddr.Hex()

	// the relayer wallet submits the tx and pays gas; `from` only signs
	w, _ := wallet.New("<relayerPrivateKey>")
	w.Connect(alchemy.GetProvider())

	now := time.Now().Unix()
	validAfter := big.NewInt(0)
	validBefore := big.NewInt(now + int64(10*time.Minute/time.Second))
	value := big.NewInt(100)               // base units — JPYC has 18 decimals
	nonce := utils.NewAuthorizationNonce() // cryptographically random [32]byte

	domainSeparator, _ := alchemy.StableCoin.DomainSeparator(jpyc)

	// the holder signs the authorization off-chain (no gas, no tx)
	sig, _ := typeddata.SignEIP712Str(
		"<fromPrivateKey>",
		domainSeparator,
		typeddata.EncodeWords(
			constant.TransferWithAuthorizationTypeHash,
			"<fromAddress>", "<toAddress>", value, validAfter, validBefore, nonce,
		),
	)

	// the relayer broadcasts it and pays the gas
	receipt, _ := w.StableCoin().TransferWithAuthorization(
		context.Background(), jpyc,
		"<fromAddress>", "<toAddress>", value, validAfter, validBefore, nonce, sig, nil,
	)
	_ = receipt
}
```

### Trying it without real funds (SimulatedBackend)

You can run the whole flow offline against the in-process `SimulatedBackend`. JPYC is a
`FiatToken`, so deploying the `FiatToken` fixture reproduces the exact gasless behavior. The
example below deploys the token, mints to a **holder**, then has the holder sign off-chain while a
separate **relayer** submits the transaction and pays the gas:

```go
package main

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient/simulated"
	"github.com/poteto-go/go-alchemy-sdk/_fixture/artifacts"
	"github.com/poteto-go/go-alchemy-sdk/constant"
	"github.com/poteto-go/go-alchemy-sdk/deployer"
	"github.com/poteto-go/go-alchemy-sdk/gas"
	"github.com/poteto-go/go-alchemy-sdk/typeddata"
	"github.com/poteto-go/go-alchemy-sdk/types"
	"github.com/poteto-go/go-alchemy-sdk/utils"
	"github.com/poteto-go/go-alchemy-sdk/wallet"
)

// anvil dev accounts: relayer pays gas, holder only signs
const (
	relayerPK   = "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	relayerAddr = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
	holderPK    = "0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d"
	holderAddr  = "0x70997970C51812dc3A010C7d01b50e0d17dc79C8"
)

func main() {
	oneEth := big.NewInt(1_000_000_000_000_000_000)
	backend := simulated.NewBackend(gethTypes.GenesisAlloc{
		common.HexToAddress(relayerAddr): {Balance: new(big.Int).Mul(big.NewInt(1000), oneEth)},
		common.HexToAddress(holderAddr):  {Balance: new(big.Int).Mul(big.NewInt(1000), oneEth)},
	})
	defer backend.Close()

	alchemy, _ := gas.NewSimulatedAlchemy(backend)

	w, _ := wallet.New(relayerPK) // relayer wallet (pays gas)
	w.Connect(alchemy.GetProvider())

	// deploy the FiatToken fixture (JPYC is a FiatToken)
	owner := common.HexToAddress(relayerAddr)
	meta := &artifacts.FiatTokenMetaData
	deployer.BindDeploymentMetadata(meta,
		"JPY Coin", "JPYC", "JPY", uint8(18),
		owner, owner, owner, owner, // masterMinter, pauser, blacklister, owner
	)
	contractAddr, _ := w.DeployContract(context.Background(), meta)
	jpyc := contractAddr.Hex()

	// configure the relayer as a minter, then mint to the holder
	fiatToken := artifacts.NewFiatToken()
	txHash, _ := w.SendTransaction(types.TransactionRequest{
		From: relayerAddr, To: jpyc, Value: "0x0", GasLimit: 300000,
		Data: fiatToken.PackConfigureMinter(owner, big.NewInt(1_000_000_000)),
	})
	alchemy.Transact.WaitMined(context.Background(), txHash.Hex())
	w.StableCoin().Mint(context.Background(), jpyc, holderAddr, big.NewInt(1000), nil)

	// --- gasless flow ---
	value := big.NewInt(100)
	validAfter := big.NewInt(0)
	validBefore := big.NewInt(9999999999)
	nonce := utils.NewAuthorizationNonce()

	domainSeparator, _ := alchemy.StableCoin.DomainSeparator(jpyc)

	// the holder signs off-chain (no gas, no tx)
	sig, _ := typeddata.SignEIP712Str(
		holderPK,
		domainSeparator,
		typeddata.EncodeWords(
			constant.TransferWithAuthorizationTypeHash,
			holderAddr,  // from
			relayerAddr, // to
			value, validAfter, validBefore, nonce,
		),
	)

	// the relayer broadcasts it and pays the gas
	receipt, _ := w.StableCoin().TransferWithAuthorization(
		context.Background(), jpyc,
		holderAddr, relayerAddr, value, validAfter, validBefore, nonce, sig, nil,
	)
	_ = receipt

	// holder balance: 1000 → 900, recipient: 0 → 100; the nonce is now spent
	used, _ := alchemy.StableCoin.AuthorizationState(jpyc, holderAddr, nonce)
	_ = used // true
}
```

`Nonces` and `AuthorizationState` let you check whether a given authorization has already been
used before submitting. For a real network, swap `simulated.NewBackend` for `gas.NewAlchemy` with
your Alchemy key and use the existing on-chain JPYC address instead of deploying the fixture.

## Notes

- **Decimals: 18** (not 6). Always format with the token's own decimals.
- **Contract version.** The new 資金移動業 JPYC is `0xE7C3D8C9a439feDe00D2600032D5dB0Be71C3c29`
  (same on Ethereum / Polygon / Avalanche). The older JPYC v1 (prepaid) token uses different
  addresses and is a separate product.
- **`famous` registry.** `famous.ContractAddress(network, famous.JPYC)` returns the current
  資金移動業 JPYC address (`0xE7C3D8C9a439feDe00D2600032D5dB0Be71C3c29`) on Ethereum and Polygon —
  which is why the examples above resolve the address from the registry instead of hard-coding it.
