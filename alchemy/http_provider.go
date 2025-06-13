package alchemy

import (
	"bytes"
	"io"
	"net/http"

	"github.com/goccy/go-json"
	"github.com/poteto-go/go-alchemy-sdk/core"
	"github.com/poteto-go/go-alchemy-sdk/types"
	"github.com/poteto-go/go-alchemy-sdk/utils"
)

type AlchemyProvider struct {
	config AlchemyConfig
	id     int
}

func NewAlchemyProvider(config AlchemyConfig) types.IAlchemyProvider {
	return &AlchemyProvider{
		config: config,
		id:     1,
	}
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
	return int(blockNumber), nil
}

func (provider *AlchemyProvider) Send(method string, params ...string) (string, error) {
	return provider.send(method, params...)
}

func (provider *AlchemyProvider) send(method string, params ...string) (string, error) {
	if len(params) == 0 {
		params = []string{}
	}

	reqParam := types.AlchemyRequest{
		Jsonrpc: "2.0",
		Method:  method,
		Params:  params,
		Id:      provider.id,
	}

	paramJson, err := json.Marshal(reqParam)
	if err != nil {
		return "", core.ErrFailedToMarshalParameter
	}

	req, err := http.NewRequest("POST", provider.config.GetUrl(), bytes.NewBuffer(paramJson))
	if err != nil {
		return "", core.ErrFailedToCreateRequest
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", core.ErrFailedToConnect
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	result := types.AlchemyResponse{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", core.ErrFailedToUnmarshalResponse
	}

	provider.id++

	return result.Result, nil
}
