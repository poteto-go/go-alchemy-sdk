package e2e

import (
	"math/big"
	"os"
	"strconv"
	"testing"

	"github.com/poteto-go/go-alchemy-sdk/gas"
	"github.com/stretchr/testify/assert"
)

var initAddress = "0x8943545177806ED17B9F23F0a21ee5948eCaa776"
var alchemy gas.Alchemy

func TestMain(m *testing.M) {
	setup()

	m.Run()

	teardown()
}

func setup() {
	port, err := strconv.Atoi(os.Getenv("RPC_PORT"))
	if err != nil {
		panic(err)
	}

	setting := gas.AlchemySetting{
		PrivateNetworkConfig: gas.PrivateNetworkConfig{
			Host: "127.0.0.1",
			Port: port,
		},
	}
	alchemy = gas.NewAlchemy(setting)
}

func teardown() {
}

func TestScenario_GetBalance(t *testing.T) {
	balance, err := alchemy.Core.GetBalance(initAddress, "latest")

	assert.Nil(t, err)
	assert.Equal(t, balance.Cmp(big.NewInt(0)), 1)
}
