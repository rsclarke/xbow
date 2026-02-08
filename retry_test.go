package xbow

import (
	"context"
	"errors"
	"net/http"
	"sync/atomic"
	"testing"
	"time"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func newRetryTransport(base http.RoundTripper, policy *RetryPolicy) *retryTransport {
	if policy == nil {
		policy = &RetryPolicy{}
	}
	policy.defaults()
	return &retryTransport{base: base, policy: *policy}
}

func TestRetryTransport_NoRetryOnSuccess(t *testing.T) {
	var calls atomic.Int32
	rt := newRetryTransport(roundTripFunc(func(req *http.Request) (*http.Response, error) {
		calls.Add(1)
		return &http.Response{StatusCode: 200, Body: http.NoBody}, nil
	}), nil)

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://example.com", nil)
	resp, err := rt.RoundTrip(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Errorf("StatusCode = %d, want 200", resp.StatusCode)
	}
	if got := calls.Load(); got != 1 {
		t.Errorf("calls = %d, want 1", got)
	}
}

func TestRetryTransport_RetriesOnRetryableStatus(t *testing.T) {
	for _, status := range []int{429, 500, 502, 503, 504} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			var calls atomic.Int32
			rt := newRetryTransport(roundTripFunc(func(req *http.Request) (*http.Response, error) {
				calls.Add(1)
				return &http.Response{StatusCode: status, Body: http.NoBody}, nil
			}), &RetryPolicy{
				MaxAttempts:    3,
				InitialBackoff: time.Millisecond,
				MaxBackoff:     10 * time.Millisecond,
				Jitter:         true,
			})

			req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://example.com", nil)
			resp, err := rt.RoundTrip(req)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode != status {
				t.Errorf("StatusCode = %d, want %d", resp.StatusCode, status)
			}
			if got := calls.Load(); got != 3 {
				t.Errorf("calls = %d, want 3", got)
			}
		})
	}
}

func TestRetryTransport_ReturnsOnEventualSuccess(t *testing.T) {
	var calls atomic.Int32
	rt := newRetryTransport(roundTripFunc(func(req *http.Request) (*http.Response, error) {
		n := calls.Add(1)
		if n < 3 {
			return &http.Response{StatusCode: 503, Body: http.NoBody}, nil
		}
		return &http.Response{StatusCode: 200, Body: http.NoBody}, nil
	}), &RetryPolicy{
		MaxAttempts:    4,
		InitialBackoff: time.Millisecond,
		MaxBackoff:     10 * time.Millisecond,
		Jitter:         true,
	})

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://example.com", nil)
	resp, err := rt.RoundTrip(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Errorf("StatusCode = %d, want 200", resp.StatusCode)
	}
	if got := calls.Load(); got != 3 {
		t.Errorf("calls = %d, want 3", got)
	}
}

func TestRetryTransport_SkipsPOSTByDefault(t *testing.T) {
	var calls atomic.Int32
	rt := newRetryTransport(roundTripFunc(func(req *http.Request) (*http.Response, error) {
		calls.Add(1)
		return &http.Response{StatusCode: 503, Body: http.NoBody}, nil
	}), &RetryPolicy{
		InitialBackoff: time.Millisecond,
	})

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, "https://example.com", nil)
	resp, err := rt.RoundTrip(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 503 {
		t.Errorf("StatusCode = %d, want 503", resp.StatusCode)
	}
	if got := calls.Load(); got != 1 {
		t.Errorf("calls = %d, want 1 (no retry for POST)", got)
	}
}

func TestRetryTransport_RetriesPOSTWhenEnabled(t *testing.T) {
	var calls atomic.Int32
	rt := newRetryTransport(roundTripFunc(func(req *http.Request) (*http.Response, error) {
		calls.Add(1)
		return &http.Response{StatusCode: 503, Body: http.NoBody}, nil
	}), &RetryPolicy{
		MaxAttempts:    2,
		InitialBackoff: time.Millisecond,
		MaxBackoff:     10 * time.Millisecond,
		Jitter:         true,
		RetryPOST:      true,
	})

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, "https://example.com", nil)
	resp, err := rt.RoundTrip(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()
	if got := calls.Load(); got != 2 {
		t.Errorf("calls = %d, want 2", got)
	}
}

func TestRetryTransport_RespectsContextCancellation(t *testing.T) {
	var calls atomic.Int32
	rt := newRetryTransport(roundTripFunc(func(req *http.Request) (*http.Response, error) {
		calls.Add(1)
		return &http.Response{StatusCode: 503, Body: http.NoBody}, nil
	}), &RetryPolicy{
		MaxAttempts:    5,
		InitialBackoff: 100 * time.Millisecond,
		MaxBackoff:     time.Second,
	})

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "https://example.com", nil)
	resp, err := rt.RoundTrip(req)
	if resp != nil {
		_ = resp.Body.Close()
	}
	if !errors.Is(err, context.Canceled) {
		t.Errorf("err = %v, want context.Canceled", err)
	}
	if got := calls.Load(); got >= 5 {
		t.Errorf("calls = %d, expected fewer than MaxAttempts (5)", got)
	}
}

func TestRetryTransport_NoRetryOnTransportError(t *testing.T) {
	var calls atomic.Int32
	transportErr := &netError{msg: "connection refused"}
	rt := newRetryTransport(roundTripFunc(func(req *http.Request) (*http.Response, error) {
		calls.Add(1)
		return nil, transportErr
	}), &RetryPolicy{
		MaxAttempts:    3,
		InitialBackoff: time.Millisecond,
	})

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://example.com", nil)
	resp, err := rt.RoundTrip(req)
	if resp != nil {
		_ = resp.Body.Close()
	}
	if !errors.Is(err, transportErr) {
		t.Errorf("err = %v, want %v", err, transportErr)
	}
	if got := calls.Load(); got != 1 {
		t.Errorf("calls = %d, want 1 (no retry on transport error)", got)
	}
}

type netError struct {
	msg string
}

func (e *netError) Error() string { return e.msg }

func TestRetryTransport_IdempotentMethods(t *testing.T) {
	for _, method := range []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodDelete} {
		t.Run(method, func(t *testing.T) {
			var calls atomic.Int32
			rt := newRetryTransport(roundTripFunc(func(req *http.Request) (*http.Response, error) {
				calls.Add(1)
				return &http.Response{StatusCode: 503, Body: http.NoBody}, nil
			}), &RetryPolicy{
				MaxAttempts:    2,
				InitialBackoff: time.Millisecond,
				MaxBackoff:     10 * time.Millisecond,
				Jitter:         true,
			})

			req, _ := http.NewRequestWithContext(context.Background(), method, "https://example.com", nil)
			resp, err := rt.RoundTrip(req)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			defer resp.Body.Close()
			if got := calls.Load(); got != 2 {
				t.Errorf("calls = %d, want 2", got)
			}
		})
	}
}

func TestRetryTransport_NonRetryableStatusNotRetried(t *testing.T) {
	for _, status := range []int{400, 401, 403, 404} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			var calls atomic.Int32
			rt := newRetryTransport(roundTripFunc(func(req *http.Request) (*http.Response, error) {
				calls.Add(1)
				return &http.Response{StatusCode: status, Body: http.NoBody}, nil
			}), &RetryPolicy{
				InitialBackoff: time.Millisecond,
			})

			req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://example.com", nil)
			resp, err := rt.RoundTrip(req)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode != status {
				t.Errorf("StatusCode = %d, want %d", resp.StatusCode, status)
			}
			if got := calls.Load(); got != 1 {
				t.Errorf("calls = %d, want 1", got)
			}
		})
	}
}

func TestRetryPolicyDefaults(t *testing.T) {
	p := &RetryPolicy{}
	p.defaults()

	if p.MaxAttempts != 3 {
		t.Errorf("MaxAttempts = %d, want 3", p.MaxAttempts)
	}
	if p.InitialBackoff != 500*time.Millisecond {
		t.Errorf("InitialBackoff = %v, want 500ms", p.InitialBackoff)
	}
	if p.MaxBackoff != 30*time.Second {
		t.Errorf("MaxBackoff = %v, want 30s", p.MaxBackoff)
	}
	if p.RetryableStatusCodes == nil {
		t.Fatal("RetryableStatusCodes should not be nil")
	}
	want := []int{429, 500, 502, 503, 504}
	if len(p.RetryableStatusCodes) != len(want) {
		t.Fatalf("RetryableStatusCodes length = %d, want %d", len(p.RetryableStatusCodes), len(want))
	}
	for i, v := range want {
		if p.RetryableStatusCodes[i] != v {
			t.Errorf("RetryableStatusCodes[%d] = %d, want %d", i, p.RetryableStatusCodes[i], v)
		}
	}
}

func TestRetryPolicyDefaultsPreservesExplicit(t *testing.T) {
	p := &RetryPolicy{
		MaxAttempts:          5,
		InitialBackoff:       time.Second,
		MaxBackoff:           time.Minute,
		RetryableStatusCodes: []int{429},
	}
	p.defaults()

	if p.MaxAttempts != 5 {
		t.Errorf("MaxAttempts = %d, want 5", p.MaxAttempts)
	}
	if p.InitialBackoff != time.Second {
		t.Errorf("InitialBackoff = %v, want 1s", p.InitialBackoff)
	}
	if p.MaxBackoff != time.Minute {
		t.Errorf("MaxBackoff = %v, want 1m", p.MaxBackoff)
	}
	if len(p.RetryableStatusCodes) != 1 || p.RetryableStatusCodes[0] != 429 {
		t.Errorf("RetryableStatusCodes = %v, want [429]", p.RetryableStatusCodes)
	}
}

func TestBackoff(t *testing.T) {
	rt := &retryTransport{
		policy: RetryPolicy{
			InitialBackoff: 100 * time.Millisecond,
			MaxBackoff:     time.Second,
			Jitter:         false,
		},
	}

	tests := []struct {
		attempt int
		want    time.Duration
	}{
		{0, 100 * time.Millisecond},
		{1, 200 * time.Millisecond},
		{2, 400 * time.Millisecond},
		{3, 800 * time.Millisecond},
		{4, time.Second},
		{5, time.Second},
	}

	for _, tt := range tests {
		got := rt.backoff(tt.attempt)
		if got != tt.want {
			t.Errorf("backoff(%d) = %v, want %v", tt.attempt, got, tt.want)
		}
	}
}

func TestBackoffWithJitter(t *testing.T) {
	rt := &retryTransport{
		policy: RetryPolicy{
			InitialBackoff: 100 * time.Millisecond,
			MaxBackoff:     time.Second,
			Jitter:         true,
		},
	}

	for range 100 {
		got := rt.backoff(0)
		if got < 0 || got > 100*time.Millisecond {
			t.Errorf("backoff(0) with jitter = %v, want [0, 100ms]", got)
		}
	}
}
