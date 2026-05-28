package utils

import (
	"io"
	"net/http"

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
// response bodies at maxBytes. If maxBytes is 0, DefaultMaxResponseBytes is used.
// The client is intended to be constructed once and reused across requests so
// that the underlying connection pool (http.DefaultTransport) is shared.
func NewSharedHTTPClient(maxBytes int64) *http.Client {
	if maxBytes == 0 {
		maxBytes = types.DefaultMaxResponseBytes
	}
	return &http.Client{
		Transport: &limitedTransport{
			underlying: &defaultRoundTripper{},
			maxBytes:   maxBytes,
		},
	}
}
