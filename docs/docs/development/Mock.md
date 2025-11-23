---
sidebar_position: 1
---

# Mock

## Overview

You can mock go-alchemy-sdk's response in UT.

## Usage

```go
func TestSomething(t *testing.T) {
	setting := gas.AlchemySetting{
		ApiKey:  "<alchemy-api-key>",
		Network: types.EthSepolia,
	}
	alchemy := gas.NewAlchemy(setting)

	mock := alchemymock.NewAlchemyHttpMock(setting, t)
	defer mock.DeactivateAndReset()

	mock.RegisterResponder("eth_getBalance", `{"jsonrpc":"2.0","id":1,"result":"0x1234"}`)

	balance, err := alchemy.Core.GetBalance("0x", "latest")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "4660", balance.String())
}
```

## Detail

`alchemymock` mock jsonrpc response from blockchain.
If you want not to change on public chain, you can test w/ mock easily w/o change on public chain.
