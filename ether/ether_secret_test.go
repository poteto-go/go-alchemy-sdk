package ether_test

import (
	"errors"
	"net/http"
	"reflect"
	"testing"
	"testing/synctest"
	"time"

	"github.com/agiledragon/gomonkey"
	"github.com/poteto-go/go-alchemy-sdk/ether"
	eth "github.com/poteto-go/go-alchemy-sdk/ether"
	"github.com/poteto-go/go-alchemy-sdk/gas"
	"github.com/poteto-go/go-alchemy-sdk/internal"
	"github.com/stretchr/testify/assert"
)

func newEtherApiWSecretForTest() *eth.Ether {
	provider := newProviderForTest()
	config, err := gas.NewAlchemyConfig(utAlchemySetting)
	if err != nil {
		panic(err)
	}

	decoded, err := internal.DecodeHex("0xbcdf20249abf0ed6d944c0288fad489e33f66b3960d9e6229c1cd214ed3bbe31")
	if err != nil {
		panic(err)
	}

	return ether.NewEtherApi(
		provider,
		eth.NewEtherApiConfig(
			config.GetUrl(),
			0,
			time.Duration(1*time.Second),
			nil,
			[]http.Header{},
			decoded,
			0,
		),
	).(*eth.Ether)
}

func newEtherApiWIvalidSecretForTest() *eth.Ether {
	provider := newProviderForTest()
	config, err := gas.NewAlchemyConfig(utAlchemySetting)
	if err != nil {
		panic(err)
	}

	return ether.NewEtherApi(
		provider,
		eth.NewEtherApiConfig(
			config.GetUrl(),
			0,
			time.Duration(1*time.Second),
			nil,
			[]http.Header{},
			[]byte("invalid"),
			0,
		),
	).(*eth.Ether)
}

func TestEther_BlockNumberWJwt(t *testing.T) {
	t.Run("normal case", func(t *testing.T) {
		t.Run("success request", func(t *testing.T) {
			// Arrange
			ether := newEtherApiWSecretForTest()
			alchemyMock := newAlchemyMockOnEtherTest(t)
			defer alchemyMock.DeactivateAndReset()

			// Mock
			alchemyMock.RegisterResponderOnce("eth_blockNumber", `{"jsonrpc":"2.0","id":1,"result":"0x1234"}`)

			// Act
			result, err := ether.BlockNumber()

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, result, uint64(0x1234))
		})
	})

	t.Run("error case", func(t *testing.T) {
		t.Run("if cannot create ethClient, return err", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Arrange
			ether := newEtherApiWSecretForTest()

			// Mock
			patches.ApplyMethod(
				reflect.TypeOf(ether),
				"SetEthClient",
				func(_ *eth.Ether) error {
					return errors.New("error")
				},
			)

			// Act
			_, err := ether.BlockNumber()

			// Assert
			assert.Error(t, err)
		})

		t.Run("if failed estimate gas, return error", func(t *testing.T) {
			// Arrange
			ether := newEtherApiWSecretForTest()

			// Act
			_, err := ether.BlockNumber()

			// Assert
			assert.Error(t, err)
		})
	})
}

func TestEther_RecreateExpiredJws(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		// Arrange
		ether := newEtherApiWSecretForTest()
		alchemyMock := newAlchemyMockOnEtherTest(t)
		defer alchemyMock.DeactivateAndReset()
		ether.SetEthClient()

		// Mock
		alchemyMock.RegisterResponderOnce("eth_blockNumber", `{"jsonrpc":"2.0","id":1,"result":"0x1234"}`)

		// wait 2min
		time.Sleep(2 * time.Minute)

		// Re-Act
		_, err := ether.BlockNumber()

		// Assert
		assert.Nil(t, err)
	})
}

func TestEther_JwtExpiry_SafetyMargin(t *testing.T) {
	t.Run("within 56 seconds: same client is reused", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			e := newEtherApiWSecretForTest()
			err := e.SetEthClient()
			assert.NoError(t, err)
			firstClient := e.Client()

			time.Sleep(56 * time.Second)

			err = e.SetEthClient()
			assert.NoError(t, err)

			assert.Same(t, firstClient, e.Client())
		})
	})

	t.Run("at 57 seconds: client is recreated before geth 60s window", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			e := newEtherApiWSecretForTest()
			err := e.SetEthClient()
			assert.NoError(t, err)
			firstClient := e.Client()

			time.Sleep(57 * time.Second)

			err = e.SetEthClient()
			assert.NoError(t, err)

			assert.NotSame(t, firstClient, e.Client())
		})
	})

	t.Run("after 60 seconds: client is recreated with fresh JWT", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			e := newEtherApiWSecretForTest()
			err := e.SetEthClient()
			assert.NoError(t, err)
			firstClient := e.Client()

			time.Sleep(60 * time.Second)

			err = e.SetEthClient()
			assert.NoError(t, err)

			assert.NotSame(t, firstClient, e.Client())
		})
	})
}

func TestEther_InvalidJwtSecret(t *testing.T) {
	t.Run("success request", func(t *testing.T) {
		// Arrange
		ether := newEtherApiWIvalidSecretForTest()

		// Act
		_, err := ether.BlockNumber()

		// Assert
		assert.Error(t, err)
	})
}
