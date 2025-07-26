package utils

import (
	"errors"
	"testing"

	"github.com/agiledragon/gomonkey"
	"github.com/goccy/go-json"
	"github.com/poteto-go/go-alchemy-sdk/core"
	"github.com/stretchr/testify/assert"
)

func TestCreateRequestBodyToBytes(t *testing.T) {
	t.Run("normal case:", func(t *testing.T) {
		tests := []struct {
			name     string
			id       int
			method   string
			params   []string
			wantBody string
		}{
			{
				name:     "no params",
				id:       1,
				method:   core.Eth_BlockNumber,
				params:   []string{},
				wantBody: `{"jsonrpc":"2.0","method":"eth_blockNumber","id":1}`,
			},
			{
				name:     "with params",
				id:       1,
				method:   core.Eth_GetTransactionByHash,
				params:   []string{"0x1b4"},
				wantBody: `{"jsonrpc":"2.0","method":"eth_getTransactionByHash","params":["0x1b4"],"id":1}`,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := CreateRequestBodyToBytes(tt.id, tt.method, tt.params)
				if err != nil {
					t.Fatal(err)
				}
				assert.JSONEq(t, tt.wantBody, string(got))
			})
		}
	})

	t.Run("error case:", func(t *testing.T) {
		t.Run("if failed to marshal parameter -> core.ErrFailedToMarshalParameter", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Mock
			patches.ApplyFunc(
				json.Marshal,
				func(v interface{}) ([]byte, error) {
					return nil, errors.New("error")
				},
			)

			// Act
			_, err := CreateRequestBodyToBytes(1, "method", []string{"param1", "param2"})

			// Assert
			assert.ErrorIs(t, err, core.ErrFailedToMarshalParameter)
		})
	})
}
