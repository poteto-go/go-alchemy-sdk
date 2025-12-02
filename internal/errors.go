package internal

import (
	"net/http"
	"net/url"
	"slices"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/poteto-go/go-alchemy-sdk/constant"
)

// Determine whether to back off based on the error.
// For cases where repeated transmission yields no change in result, do not back off.
//
// refs:
//   - https://ethereum-json-rpc.com/errors
//   - https://www.alchemy.com/docs/reference/error-reference
//   - https://www.quicknode.com/docs/ethereum/error-references
//   - https://docs.starkware.co/starkex/api/spot/error_codes.html
func isAlwaysReProduceError(err error) bool {
	if err == nil {
		return false
	}

	switch assertedErr := err.(type) {
	// url error always re-produce
	case *url.Error:
		return false

	// rpc json error
	case rpc.Error:
		return isAlwaysReProduceRpcError(assertedErr)

	case rpc.HTTPError:
		return isAlwaysReProduceHttpError(assertedErr)

	default:
		return true
	}
}

func isAlwaysReProduceRpcError(err rpc.Error) bool {
	errorCode := err.ErrorCode()
	errorMsg := err.Error()

	// internal error
	// not found
	// stack limit reached
	// method handler crashed
	// execution timeout
	// nonce too low
	// filter not found
	if errorCode == -32000 {
		if errorMsg == "execution timeout" {
			return true
		}

		return false
	}

	// resource not found
	if errorCode == -32001 {
		return false
	}

	// resource unavailable
	if errorCode == -32002 {
		return false
	}

	// transaction rejected
	if errorCode == -32003 {
		return false
	}

	// method not supported
	if errorCode == -32004 {
		return false
	}

	// limit exceeded
	if errorCode == -32005 {
		// tx pool limit exceeded
		// if you wait backoff, it can work on next-try
		return true
	}

	// json rpc version not supported
	if errorCode == -32006 {
		return false
	}

	// custom error
	// Unable to complete request at this time.
	if -32099 <= errorCode && errorCode <= -32007 {
		return true
	}

	// custom error
	// any application error
	if -32599 <= errorCode && errorCode <= -32100 {
		return false
	}
	if -32699 <= errorCode && errorCode <= -32604 {
		return false
	}
	if -32768 <= errorCode && errorCode <= -32701 {
		return false
	}

	// invalid request
	if errorCode == -32600 {
		return false
	}

	// method not found
	if errorCode == -32601 {
		return false
	}

	// invalid params
	if errorCode == -32602 {
		return false
	}

	// internal error
	if errorCode == -32603 {
		return false
	}

	// parse error
	if errorCode == -32700 {
		return false
	}

	// StarkNet specific error
	if 0 <= errorCode && errorCode <= 62 {
		return false
	}

	return false
}

func isAlwaysReProduceHttpError(err rpc.HTTPError) bool {
	errorCode := err.StatusCode

	// this is alchemy rate limit
	// if you wait backoff, it can work on next-try
	if errorCode == http.StatusTooManyRequests {
		return true
	}

	// client error always re-produce
	if slices.Contains(constant.HttpClientErrorCodeList, errorCode) {
		return false
	}

	// custom error
	return true
}
