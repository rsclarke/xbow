package xbow

import (
	"context"
	"net/http"
)

// RateLimiter defines the interface for rate limiting API requests.
// Implementations should block until the request is allowed to proceed,
// or return an error (e.g., context cancellation).
type RateLimiter interface {
	// Wait blocks until the rate limiter allows a request to proceed.
	// Returns an error if the context is cancelled or the limiter fails.
	Wait(ctx context.Context) error
}

// rateLimitTransport wraps an http.RoundTripper with rate limiting.
type rateLimitTransport struct {
	base    http.RoundTripper
	limiter RateLimiter
}

func (t *rateLimitTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if err := t.limiter.Wait(req.Context()); err != nil {
		return nil, err
	}
	return t.base.RoundTrip(req)
}
