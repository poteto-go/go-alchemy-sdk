package types

import (
	"testing"

	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
)

func TestTransactionRequestWithBlockTag(t *testing.T) {
	t.Run("normal case:", func(t *testing.T) {
		// Arrange
		txs := RequestArgs{
			TransactionRequest{
				To:    "0x2345",
				Value: "0x1",
			},
			"latest",
		}
		expected := `[{"to":"0x2345","value":"0x1"},"latest"]`

		// Act
		actual, _ := json.Marshal(txs)

		// Assert
		assert.Equal(t, expected, string(actual))
	})
}
