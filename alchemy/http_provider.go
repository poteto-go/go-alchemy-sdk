package alchemy

import (
	"net/http"

	"github.com/poteto-go/go-alchemy-sdk/core"
	"github.com/poteto-go/go-alchemy-sdk/internal"
	"github.com/poteto-go/go-alchemy-sdk/types"
	"github.com/poteto-go/go-alchemy-sdk/utils"
)

type AlchemyProvider struct {
	config AlchemyConfig
	id     int
}

func NewAlchemyProvider(config AlchemyConfig) types.IAlchemyProvider {
	provider := &AlchemyProvider{
		config: config,
		id:     1,
	}

	return provider
}

/* get  the number of the most recent block. */
func (provider *AlchemyProvider) GetBlockNumber() (int, error) {
	blockNumberHex, err := provider.Send("eth_blockNumber")
	if err != nil {
		return 0, err
	}
	blockNumber, err := utils.FromHex(blockNumberHex)
	if err != nil {
		return 0, err
	}
	return blockNumber, nil
}

func (provider *AlchemyProvider) Send(method string, params ...string) (string, error) {
	return provider.send(method, params...)
}

func (provider *AlchemyProvider) send(method string, params ...string) (string, error) {
	if len(params) == 0 {
		params = []string{}
	}

	body := types.AlchemyRequestBody{
		Jsonrpc: "2.0",
		Method:  method,
		Params:  params,
		Id:      provider.id,
	}

	req, err := http.NewRequest("POST", provider.config.GetUrl(), nil)
	if err != nil {
		return "", core.ErrFailedToCreateRequest
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Alchemy-Ethers-Sdk-Method", "send")

	request := types.AlchemyRequest{
		Body:    body,
		Request: req,
	}

	result, err := internal.RequestHttpWithBackoff(
		provider.config.backoffConfig,
		utils.AlchemyFetch,
		request,
	)
	if err != nil {
		return "", err
	}

	provider.id++

	return result.Result, nil
}
