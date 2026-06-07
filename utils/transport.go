package utils

import (
	"io"
	"net/http"
	"time"

	"github.com/poteto-go/go-alchemy-sdk/types"
)

// defaultRoundTripper delegates to http.DefaultTransport at call time, so that
// test frameworks (e.g. httpmock) that swap http.DefaultTransport can intercept
// requests even when the http.Client was created before the swap.
type defaultRoundTripper struct{}

func (d *defaultRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return http.DefaultTransport.RoundTrip(req)
}

type limitedReadCloser struct {
	io.Reader
	io.Closer
}

// limitedTransport wraps http.RoundTripper to cap response bodies at maxBytes.
type limitedTransport struct {
	underlying http.RoundTripper
	maxBytes   int64
}

// RoundTrip intercepts every HTTP response and caps the body at maxBytes:
//
//	caller (e.g. AlchemyFetch / geth rpc.Client)
//	  └─ http.Client.Do(req)
//	       └─ Transport.RoundTrip(req)   ← called automatically by Go's http.Client
//	            └─ limitedTransport.RoundTrip()  [utils/transport.go]
//	                 ├─ t.underlying.RoundTrip(req)  ← actual HTTP communication
//	                 └─ resp.Body = LimitReader(resp.Body, maxBytes)  ← wrapped here
func (t *limitedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := t.underlying.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	resp.Body = limitedReadCloser{
		Reader: io.LimitReader(resp.Body, t.maxBytes),
		Closer: resp.Body,
	}
	return resp, nil
}

// NewSharedHTTPClient returns an *http.Client with a limitedTransport that caps
// response bodies at maxBytes and a Timeout set to timeout. If maxBytes is 0,
// DefaultMaxResponseBytes is used. The client is intended to be constructed once
// and reused across requests so that the underlying connection pool is shared
// and no TCP/TLS handshake overhead is paid on repeated calls to the same host.
//
// transport is the caller-supplied http.RoundTripper used for the actual HTTP
// communication (connection pooling, retry/backoff, tracing, metrics, or
// provider benchmarking). If transport is nil, requests delegate to
// http.DefaultTransport at call time. The response-size cap is always applied
// on top of the supplied transport.
func NewSharedHTTPClient(maxBytes int64, timeout time.Duration, transport http.RoundTripper) *http.Client {
	if maxBytes == 0 {
		maxBytes = types.DefaultMaxResponseBytes
	}
	if transport == nil {
		transport = &defaultRoundTripper{}
	}
	return &http.Client{
		Timeout: timeout,
		Transport: &limitedTransport{
			underlying: transport,
			maxBytes:   maxBytes,
		},
	}
}
