package constant

import (
	"errors"
	"net/http"
)

var (
	ErrInvalidHexString                 = errors.New("invalid hex string")
	ErrFailedToMarshalParameter         = errors.New("failed to marshal parameter")
	ErrFailedToCreateRequest            = errors.New("failed to create request")
	ErrFailedToConnect                  = errors.New("failed to connect")
	ErrFailedToUnmarshalResponse        = errors.New("failed to unmarshal response")
	ErrFailedToUnmarshalTransaction     = errors.New("failed to unmarshal transaction")
	ErrFailedToMapTransaction           = errors.New("failed to map transaction")
	ErrFailedToMapTransactionReceipt    = errors.New("failed to map transaction receipt")
	ErrFailedToMapTransactionReceipts   = errors.New("failed to map transaction receipts")
	ErrFailedToMapTokenResponse         = errors.New("failed to map token response")
	ErrFailedToMapBlockResponse         = errors.New("failed to map block response")
	ErrBatcherNotRunning                = errors.New("batcher not running")
	ErrRequestTimeoutTooShort           = errors.New("request timeout should be longer than batch timeout")
	ErrRequestTimeout                   = errors.New("request timeout")
	ErrOverMaxRetries                   = errors.New("over max retries")
	ErrInvalidBlockTag                  = errors.New("invalid block tag")
	ErrInvalidGetTransactionReceiptsArg = errors.New("invalid get transaction receipts arg, need blockHash or blockNumber")
	ErrFailedToTransformBlockNumber     = errors.New("failed to transform block number")
	ErrFailedToTransformType            = errors.New("failed to transform type")
	ErrFailedToTransformNonce           = errors.New("failed to transform nonce")
	ErrFailedToTransformGasPrice        = errors.New("failed to transform gas price")
	ErrFailedToTransformGasLimit        = errors.New("failed to transform gas limit")
	ErrFailedToTransformValue           = errors.New("failed to transform value")
	ErrFailedToTransformChainId         = errors.New("failed to transform chain id")
	ErrFailedToTransformV               = errors.New("failed to transform v")
	ErrFailedToTransformTimestamp       = errors.New("failed to transform timestamp")
	ErrFailedToTransformDifficulty      = errors.New("failed to transform difficulty")
	ErrResultIsNil                      = errors.New("result is nil")
	ErrFailedToCreateRequestBody        = errors.New("failed to create request body")
	ErrInvalidArgs                      = errors.New("invalid args")
	ErrOverFlow                         = errors.New("overflow")
	ErrWalletIsNotConnected             = errors.New("wallet is not connected")
	ErrContractInstanceIsNil            = errors.New("contract instance is nil")
	ErrNilAmount                        = errors.New("amount must not be nil")
	ErrNegativeAmount                   = errors.New("amount must not be negative")
	ErrAmountExceedsUint256             = errors.New("amount exceeds uint256 max")
	ErrInvalidAddress                   = errors.New("invalid hex address")
	ErrInvalidABIString                 = errors.New("invalid ABI string")
	ErrFailedToReadResponse             = errors.New("failed to read response body")
	ErrUnexpectedResponseType           = errors.New("unexpected response type")
)

var HttpClientErrorCodeList = []int{
	http.StatusBadRequest,
	http.StatusUnauthorized,
	http.StatusForbidden,
	http.StatusPaymentRequired,
	http.StatusNotFound,
	http.StatusMethodNotAllowed,
	http.StatusNotAcceptable,
	http.StatusProxyAuthRequired,
	http.StatusConflict,
	http.StatusGone,
	http.StatusLengthRequired,
	http.StatusPreconditionFailed,
	http.StatusRequestEntityTooLarge,
	http.StatusRequestURITooLong,
	http.StatusUnsupportedMediaType,
	http.StatusRequestedRangeNotSatisfiable,
	http.StatusExpectationFailed,
	http.StatusTeapot,
	http.StatusMisdirectedRequest,
	http.StatusUnprocessableEntity,
	http.StatusUpgradeRequired,
	http.StatusPreconditionRequired,
}
