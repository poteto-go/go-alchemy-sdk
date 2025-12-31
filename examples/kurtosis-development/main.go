package main

import (
	"os"
	"strconv"

	"github.com/poteto-go/go-alchemy-sdk/gas"
)

var alchemy gas.Alchemy

func init() {
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

func main() {
	blockNumber, err := alchemy.Core.GetBlockNumber()
	if err != nil {
		panic(err)
	}
	println(blockNumber)
}
