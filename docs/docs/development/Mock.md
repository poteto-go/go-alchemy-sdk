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

	mock.RegisterResponderOnce("eth_getBalance", `{"jsonrpc":"2.0","id":1,"result":"0x1234"}`)

	balance, err := alchemy.Core.GetBalance("0x", "latest")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "4660", balance.String())
}
```

## WebSocket subscriptions

`NewAlchemyHttpMock` mocks one request → one response over HTTP, which cannot model
a subscription (one `eth_subscribe` request followed by a server-pushed stream).
For subscription-based code, use `NewAlchemyWsMock`: it stands up a real in-process
WebSocket JSON-RPC server and lets you push canned notifications with the `Emit*`
helpers.

```go
func TestSubscription(t *testing.T) {
	mock := alchemymock.NewAlchemyWsMock(t)
	defer mock.Close()

	// A ws-scheme url selects the WebSocket provider.
	alchemy, _ := gas.NewAlchemy(gas.AlchemySetting{
		PrivateNetworkConfig: gas.PrivateNetworkConfig{Url: mock.URL()},
	})
	defer alchemy.GetProvider().Eth().Shutdown()

	sub := alchemy.GetProvider().(types.ISubscribeProvider)
	ch := make(chan *gethTypes.Header, 4)
	subscription, _ := sub.Subscribe(context.Background(), ch, "newHeads")
	defer subscription.Unsubscribe()

	// Push canned heads after the subscription is established.
	mock.EmitNewHeads(
		&gethTypes.Header{Number: big.NewInt(0x10), Difficulty: big.NewInt(0)},
	)

	head := <-ch
	assert.Equal(t, int64(0x10), head.Number.Int64())
}
```

## Detail

If you want to test your code without making changes to a public chain, you can easily do so with mocks.
If you want not to change on public chain, you can test w/ mock easily w/o change on public chain.
