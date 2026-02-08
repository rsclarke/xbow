package xbow

import (
	"crypto/rand"
	"math"
	"math/big"
	"net/http"
	"time"
)

// RetryPolicy configures automatic retry behavior for transient failures.
type RetryPolicy struct {
	MaxAttempts          int
	InitialBackoff       time.Duration
	MaxBackoff           time.Duration
	Jitter               bool
	RetryableStatusCodes []int
	RetryPOST            bool
}

func (p *RetryPolicy) defaults() {
	if p.MaxAttempts <= 0 {
		p.MaxAttempts = 3
	}
	if p.InitialBackoff <= 0 {
		p.InitialBackoff = 500 * time.Millisecond
	}
	if p.MaxBackoff <= 0 {
		p.MaxBackoff = 30 * time.Second
	}
	if p.RetryableStatusCodes == nil {
		p.RetryableStatusCodes = []int{429, 500, 502, 503, 504}
	}
}

// WithRetryPolicy enables automatic retries with exponential backoff for
// transient failures. By default, only idempotent HTTP methods (GET, HEAD,
// PUT, DELETE) are retried. Set RetryPOST to true to also retry POST requests.
//
// Retries are performed with exponential backoff and optional jitter (enabled
// by default) to avoid thundering herd problems.
//
//	client, err := xbow.NewClient(
//	    xbow.WithOrganizationKey("key"),
//	    xbow.WithRetryPolicy(&xbow.RetryPolicy{
//	        MaxAttempts:    4,
//	        InitialBackoff: time.Second,
//	    }),
//	)
func WithRetryPolicy(p *RetryPolicy) ClientOption {
	return func(c *clientConfig) {
		c.retryPolicy = p
	}
}

// retryTransport wraps an http.RoundTripper with retry logic.
type retryTransport struct {
	base   http.RoundTripper
	policy RetryPolicy
}

func (t *retryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if !t.isRetryableMethod(req.Method) {
		return t.base.RoundTrip(req)
	}

	var resp *http.Response
	var err error

	for attempt := range t.policy.MaxAttempts {
		resp, err = t.base.RoundTrip(req)
		if err != nil {
			return nil, err
		}

		if !t.isRetryableStatus(resp.StatusCode) {
			return resp, nil
		}

		if attempt == t.policy.MaxAttempts-1 {
			return resp, nil
		}

		_ = resp.Body.Close()

		backoff := t.backoff(attempt)
		timer := time.NewTimer(backoff)
		select {
		case <-req.Context().Done():
			timer.Stop()
			return nil, req.Context().Err()
		case <-timer.C:
		}
	}

	return resp, err
}

func (t *retryTransport) isRetryableMethod(method string) bool {
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodPut, http.MethodDelete:
		return true
	case http.MethodPost:
		return t.policy.RetryPOST
	}
	return false
}

func (t *retryTransport) isRetryableStatus(status int) bool {
	for _, s := range t.policy.RetryableStatusCodes {
		if s == status {
			return true
		}
	}
	return false
}

func (t *retryTransport) backoff(attempt int) time.Duration {
	backoff := float64(t.policy.InitialBackoff) * math.Pow(2, float64(attempt))
	if backoff > float64(t.policy.MaxBackoff) {
		backoff = float64(t.policy.MaxBackoff)
	}
	if t.policy.Jitter {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(backoff)))
		backoff = float64(n.Int64())
	}
	return time.Duration(backoff)
}
