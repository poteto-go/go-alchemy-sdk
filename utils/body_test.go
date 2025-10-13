package utils

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/agiledragon/gomonkey"
	"github.com/poteto-go/go-alchemy-sdk/constant"
	"github.com/poteto-go/go-alchemy-sdk/types"
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
				method:   constant.Eth_BlockNumber,
				params:   []string{},
				wantBody: `{"jsonrpc":"2.0","method":"eth_blockNumber","id":1}`,
			},
			{
				name:     "with params",
				id:       1,
				method:   constant.Eth_GetTransactionByHash,
				params:   []string{"0x1b4"},
				wantBody: `{"jsonrpc":"2.0","method":"eth_getTransactionByHash","params":["0x1b4"],"id":1}`,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				var params []any
				if len(tt.params) > 0 {
					for _, p := range tt.params {
						params = append(params, p)
					}
				} else {
					params = nil
				}
				got, err := CreateRequestBodyToBytes(tt.id, tt.method, params)
				assert.Nil(t, err)
				assert.Equal(t, tt.wantBody, string(got))
			})
		}
	})

	t.Run("error case:", func(t *testing.T) {
		t.Run("if failed to marshal parameter -> constant.ErrFailedToMarshalParameter", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Mock
			patches.ApplyFunc(
				json.Marshal,
				func(v interface{}) ([]byte, error) {
					return nil, errors.New("error")
				},
			)

			params := types.RequestArgs{"param1", "param2"}
			var anyParams []any
			for _, p := range params {
				anyParams = append(anyParams, p)
			}

			// Act
			_, err := CreateRequestBodyToBytes(1, "method", anyParams)

			// Assert
			assert.ErrorIs(t, err, constant.ErrFailedToMarshalParameter)
		})
	})
}
