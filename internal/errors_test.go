package internal

import (
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/assert"
)

type testRPCError struct {
	code    int
	message string
}

func (e testRPCError) ErrorCode() int {
	return e.code
}

func (e testRPCError) Error() string {
	return e.message
}

func Test_isAlwaysReproduceError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
		{
			name:     "url error",
			err:      &url.Error{},
			expected: false,
		},
		{
			name:     "http.StatusTooManyRequests",
			err:      rpc.HTTPError{StatusCode: http.StatusTooManyRequests, Status: "429 Too Many Requests"},
			expected: true,
		},
		{
			name:     "client error code",
			err:      rpc.HTTPError{StatusCode: http.StatusBadRequest, Status: "400 Bad Request"},
			expected: false,
		},
		{
			name:     "other http error",
			err:      rpc.HTTPError{StatusCode: http.StatusInternalServerError, Status: "500 Internal Server Error"},
			expected: true,
		},
		{
			name:     "other error",
			err:      errors.New("other error"),
			expected: true,
		},
		{
			name:     "rpc.Error -32000 execution timeout",
			err:      testRPCError{code: -32000, message: "execution timeout"},
			expected: true,
		},
		{
			name:     "rpc.Error -32000 other",
			err:      testRPCError{code: -32000, message: "other error"},
			expected: false,
		},
		{
			name:     "rpc.Error -32001",
			err:      testRPCError{code: -32001, message: "resource not found"},
			expected: false,
		},
		{
			name:     "rpc.Error -32002",
			err:      testRPCError{code: -32002, message: "resource unavailable"},
			expected: false,
		},
		{
			name:     "rpc.Error -32003",
			err:      testRPCError{code: -32003, message: "transaction rejected"},
			expected: false,
		},
		{
			name:     "rpc.Error -32004",
			err:      testRPCError{code: -32004, message: "method not supported"},
			expected: false,
		},
		{
			name:     "rpc.Error -32005",
			err:      testRPCError{code: -32005, message: "limit exceeded"},
			expected: true,
		},
		{
			name:     "rpc.Error -32006",
			err:      testRPCError{code: -32006, message: "json rpc version not supported"},
			expected: false,
		},
		{
			name:     "rpc.Error -32050",
			err:      testRPCError{code: -32050, message: "custom error"},
			expected: true,
		},
		{
			name:     "rpc.Error -32300",
			err:      testRPCError{code: -32300, message: "custom error"},
			expected: false,
		},
		{
			name:     "rpc.Error -32600",
			err:      testRPCError{code: -32600, message: "invalid request"},
			expected: false,
		},
		{
			name:     "rpc.Error -32601",
			err:      testRPCError{code: -32601, message: "method not found"},
			expected: false,
		},
		{
			name:     "rpc.Error -32602",
			err:      testRPCError{code: -32602, message: "invalid params"},
			expected: false,
		},
		{
			name:     "rpc.Error -32603",
			err:      testRPCError{code: -32603, message: "internal error"},
			expected: false,
		},
		{
			name:     "rpc.Error -32650",
			err:      testRPCError{code: -32650, message: "custom error"},
			expected: false,
		},
		{
			name:     "rpc.Error -32700",
			err:      testRPCError{code: -32700, message: "parse error"},
			expected: false,
		},
		{
			name:     "rpc.Error -32701",
			err:      testRPCError{code: -32701, message: "custom error"},
			expected: false,
		},
		{
			name:     "rpc.Error 30",
			err:      testRPCError{code: 30, message: "starknet error"},
			expected: false,
		},
		{
			name:     "rpc.Error 100",
			err:      testRPCError{code: 100, message: "unknown error"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, isAlwaysReProduceError(tt.err))
		})
	}
}
