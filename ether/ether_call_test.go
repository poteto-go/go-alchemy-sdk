package ether_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
)

func TestEther_CallReadMethod(t *testing.T) {
	contractAddress := "0xE25583099BA105D9ec0A67f5Ae86D90e50036425"

	t.Run("normal case", func(t *testing.T) {
		t.Run("success request", func(t *testing.T) {
			// Arrange
			e := newEtherApiForTest()
			alchemyMock := newAlchemyMockOnEtherTest(t)
			defer alchemyMock.DeactivateAndReset()

			method := []byte("totalSupply()")
			expectedRes := "0x0000000000000000000000000000000000000000000000000000000000000001"
			expected, _ := hexutil.Decode(expectedRes)

			// Mock
			alchemyMock.RegisterResponderOnce(
				"eth_call",
				`{"jsonrpc":"2.0","id":1,"result":"`+expectedRes+`"}`,
			)

			// Act
			result, err := e.CallReadMethod(method, contractAddress)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, expected, result)
		})

		t.Run("success with multiple args", func(t *testing.T) {
			// Arrange
			e := newEtherApiForTest()
			alchemyMock := newAlchemyMockOnEtherTest(t)
			defer alchemyMock.DeactivateAndReset()

			method := []byte("balanceOf(address)")
			arg := make([]byte, 32)
			arg[31] = 0x01
			expectedRes := "0x0000000000000000000000000000000000000000000000000000000000000005"
			expected, _ := hexutil.Decode(expectedRes)

			// Mock
			alchemyMock.RegisterResponderOnce(
				"eth_call",
				`{"jsonrpc":"2.0","id":1,"result":"`+expectedRes+`"}`,
			)

			// Act
			result, err := e.CallReadMethod(method, contractAddress, arg)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, expected, result)
		})
	})

	t.Run("error case", func(t *testing.T) {
		t.Run("if alchemy call fails, return err", func(t *testing.T) {
			// Arrange
			e := newEtherApiForTest()
			alchemyMock := newAlchemyMockOnEtherTest(t)
			defer alchemyMock.DeactivateAndReset()

			method := []byte("totalSupply()")

			// Mock
			alchemyMock.RegisterResponderOnce(
				"eth_call",
				`{"jsonrpc":"2.0","id":1,"error":{"code":-32000,"message":"execution reverted"}}`,
			)

			// Act
			_, err := e.CallReadMethod(method, contractAddress)

			// Assert
			assert.Error(t, err)
		})
	})
}
