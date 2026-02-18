package internal

import (
	"errors"
	"net/http"
	"net/url"
	"slices"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/poteto-go/go-alchemy-sdk/constant"
)

// Determine whether provided error always reproduce or not.
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

	if _, ok := errors.AsType[*url.Error](err); ok {
		return true
	}

	if rpcError, ok := errors.AsType[rpc.Error](err); ok {
		return isAlwaysReProduceRpcError(rpcError)
	}

	if httpError, ok := errors.AsType[rpc.HTTPError](err); ok {
		return isAlwaysReProduceHttpError(httpError)
	}

	return false
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
		return errorMsg != "execution timeout"
	}

	// resource not found
	if errorCode == -32001 {
		return true
	}

	// resource unavailable
	if errorCode == -32002 {
		return true
	}

	// transaction rejected
	if errorCode == -32003 {
		return true
	}

	// method not supported
	if errorCode == -32004 {
		return true
	}

	// limit exceeded
	if errorCode == -32005 {
		// tx pool limit exceeded
		// if you wait backoff, it can work on next-try
		return false
	}

	// json rpc version not supported
	if errorCode == -32006 {
		return true
	}

	// custom error
	// Unable to complete request at this time.
	if -32099 <= errorCode && errorCode <= -32007 {
		return false
	}

	// custom error
	// any application error
	if -32599 <= errorCode && errorCode <= -32100 {
		return true
	}
	if -32699 <= errorCode && errorCode <= -32604 {
		return true
	}
	if -32768 <= errorCode && errorCode <= -32701 {
		return true
	}

	// invalid request
	if errorCode == -32600 {
		return true
	}

	// method not found
	if errorCode == -32601 {
		return true
	}

	// invalid params
	if errorCode == -32602 {
		return true
	}

	// internal error
	if errorCode == -32603 {
		return true
	}

	// parse error
	if errorCode == -32700 {
		return true
	}

	// StarkNet specific error
	if 0 <= errorCode && errorCode <= 62 {
		return true
	}

	return true
}

func isAlwaysReProduceHttpError(err rpc.HTTPError) bool {
	errorCode := err.StatusCode

	// this is alchemy rate limit
	// if you wait backoff, it can work on next-try
	if errorCode == http.StatusTooManyRequests {
		return false
	}

	// client error always re-produce
	if slices.Contains(constant.HttpClientErrorCodeList, errorCode) {
		return true
	}

	// custom error
	return false
}
