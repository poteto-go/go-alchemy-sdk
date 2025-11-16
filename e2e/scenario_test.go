package e2e

import (
	"os"
	"strconv"
	"testing"

	"github.com/poteto-go/go-alchemy-sdk/gas"
	"github.com/stretchr/testify/assert"
)

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

func TestScenario_GetBlockNumber(t *testing.T) {
	blockNumber, err := alchemy.Core.GetBlockNumber()

	assert.Nil(t, err)
	assert.NotEmpty(t, blockNumber)
}
