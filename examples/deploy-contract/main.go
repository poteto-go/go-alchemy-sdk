package main

import (
	"fmt"

	"github.com/poteto-go/go-alchemy-sdk/gas"
	"github.com/poteto-go/go-alchemy-sdk/types"
	"github.com/poteto-go/go-alchemy-sdk/wallet"
	"github.com/poteto-go/tutorial-deploy-contract/abi"
)

func main() {
	setting := gas.AlchemySetting{
		ApiKey:  "<apiKey>",
		Network: types.EthSepolia,
	}

	alchemy := gas.NewAlchemy(setting)

	w, err := wallet.New("<privateKey>")
	if err != nil {
		panic(err)
	}
	w.Connect(alchemy.GetProvider())

	address, err := w.DeployContract(&abi.StorageMetaData)
	if err != nil {
		panic(err)
	}
	fmt.Println(address)
}
