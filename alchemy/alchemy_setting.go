package alchemy

import "github.com/poteto-go/go-alchemy-sdk/types"

type AlchemySetting struct {
	ApiKey  string        `yaml:"api_key"`
	Network types.Network `yaml:"network"`
}
