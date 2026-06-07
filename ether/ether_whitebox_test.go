package ether

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestEther_SetEthClient_ConnCountRollbackOnError verifies that a failed
// SetEthClient does not permanently inflate connCount (issue #324).
func TestEther_SetEthClient_ConnCountRollbackOnError(t *testing.T) {
	e := &Ether{
		config: NewEtherApiConfig("", 0, time.Duration(0), nil, nil, []byte(""), 0, nil),
		mu:     &sync.Mutex{},
	}

	// Act: SetEthClient with an empty URL must fail at createRpcClient.
	err := e.SetEthClient()

	// Assert: error surfaced, and connCount was rolled back to 0.
	assert.Error(t, err)
	assert.Equal(t, 0, e.connCount, "connCount must not leak on error")
}
