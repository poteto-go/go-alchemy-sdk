package ether_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey"
	eth "github.com/poteto-go/go-alchemy-sdk/ether"
	"github.com/stretchr/testify/assert"
)

func TestEther_ResolveName(t *testing.T) {
	t.Run("returns lowercase address as-is when input is already a valid hex address", func(t *testing.T) {
		e := newEtherApiForTest()
		result, err := e.ResolveName("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045")
		assert.NoError(t, err)
		assert.Equal(t, "0xd8da6bf26964af9d7eed9e03e53415d37aa96045", result)
	})

	t.Run("resolves ENS name to address via registry and resolver", func(t *testing.T) {
		e := newEtherApiForTest()
		alchemyMock := newAlchemyMockOnEtherTest(t)
		defer alchemyMock.DeactivateAndReset()

		// resolver(namehash("vitalik.eth")) → resolver address
		alchemyMock.RegisterResponderOnce(
			"eth_call",
			`{"jsonrpc":"2.0","id":1,"result":"0x0000000000000000000000004976fb03c32e5b8cfe2b6ccb31c09ba78ebaba41"}`,
		)
		// addr(namehash("vitalik.eth")) → resolved address
		alchemyMock.RegisterResponderOnce(
			"eth_call",
			`{"jsonrpc":"2.0","id":1,"result":"0x000000000000000000000000d8da6bf26964af9d7eed9e03e53415d37aa96045"}`,
		)

		result, err := e.ResolveName("vitalik.eth")
		assert.NoError(t, err)
		assert.Equal(t, "0xd8da6bf26964af9d7eed9e03e53415d37aa96045", result)
	})

	t.Run("returns error when resolver is zero address (ENS not available on chain)", func(t *testing.T) {
		e := newEtherApiForTest()
		alchemyMock := newAlchemyMockOnEtherTest(t)
		defer alchemyMock.DeactivateAndReset()

		// resolver returns zero address
		alchemyMock.RegisterResponderOnce(
			"eth_call",
			`{"jsonrpc":"2.0","id":1,"result":"0x0000000000000000000000000000000000000000000000000000000000000000"}`,
		)

		_, err := e.ResolveName("vitalik.eth")
		assert.Error(t, err)
	})

	t.Run("returns error when SetEthClient fails", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		e := newEtherApiForTest()
		patches.ApplyMethod(
			reflect.TypeOf(e),
			"SetEthClient",
			func(_ *eth.Ether) error { return errors.New("connection failed") },
		)

		_, err := e.ResolveName("vitalik.eth")
		assert.Error(t, err)
	})
}

func TestEther_LookupAddress(t *testing.T) {
	t.Run("returns ENS name for a valid address", func(t *testing.T) {
		e := newEtherApiForTest()
		alchemyMock := newAlchemyMockOnEtherTest(t)
		defer alchemyMock.DeactivateAndReset()

		// resolver(reverseNode) → resolver address
		alchemyMock.RegisterResponderOnce(
			"eth_call",
			`{"jsonrpc":"2.0","id":1,"result":"0x0000000000000000000000004976fb03c32e5b8cfe2b6ccb31c09ba78ebaba41"}`,
		)
		// name(reverseNode) → ABI-encoded "vitalik.eth"
		// offset=0x20(32), length=0x0b(11), data="vitalik.eth" zero-padded to 32 bytes
		alchemyMock.RegisterResponderOnce(
			"eth_call",
			`{"jsonrpc":"2.0","id":1,"result":"0x0000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000b766974616c696b2e657468000000000000000000000000000000000000000000"}`,
		)

		result, err := e.LookupAddress("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045")
		assert.NoError(t, err)
		assert.Equal(t, "vitalik.eth", result)
	})

	t.Run("returns error for invalid hex address", func(t *testing.T) {
		e := newEtherApiForTest()
		_, err := e.LookupAddress("not-an-address")
		assert.Error(t, err)
	})

	t.Run("returns error when resolver is zero address", func(t *testing.T) {
		e := newEtherApiForTest()
		alchemyMock := newAlchemyMockOnEtherTest(t)
		defer alchemyMock.DeactivateAndReset()

		alchemyMock.RegisterResponderOnce(
			"eth_call",
			`{"jsonrpc":"2.0","id":1,"result":"0x0000000000000000000000000000000000000000000000000000000000000000"}`,
		)

		_, err := e.LookupAddress("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045")
		assert.Error(t, err)
	})

	t.Run("returns error when SetEthClient fails", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		e := newEtherApiForTest()
		patches.ApplyMethod(
			reflect.TypeOf(e),
			"SetEthClient",
			func(_ *eth.Ether) error { return errors.New("connection failed") },
		)

		_, err := e.LookupAddress("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045")
		assert.Error(t, err)
	})
}
