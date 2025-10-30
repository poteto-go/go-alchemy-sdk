package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/joho/godotenv"
	"github.com/poteto-go/go-alchemy-sdk/gas"
	"github.com/poteto-go/go-alchemy-sdk/types"
	"github.com/poteto-go/go-alchemy-sdk/wallet"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		panic(err)
	}
}

func main() {
	setting := gas.AlchemySetting{
		ApiKey:  os.Getenv("API_KEY"),
		Network: types.PolygonAmoy,
	}

	alchemy := gas.NewAlchemy(setting)

	w, err := wallet.New(os.Getenv("PRIVATE_KEY"))
	if err != nil {
		panic(err)
	}

	fmt.Println("complete: create wallet")

	w.Connect(alchemy.GetProvider())

	fmt.Println("complete: connect to wallet")

	sc, err := ioutil.ReadFile("./abi/PotetoToken.json")
	if err != nil {
		panic(err)
	}

	fmt.Println("complete: load abi")

	txRequest := types.TransactionRequest{
		From:     w.GetAddress(),
		Data:     string(sc),
		GasLimit: 10000,
	}

	if err := w.SendTransaction(txRequest); err != nil {
		panic(err)
	}

	fmt.Println("complete: deploy")
}
