package xbow

import (
	"bytes"
	"crypto/ed25519"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"io"
	"net/http"
	"strconv"
	"time"
)

const (
	// HeaderSignatureTimestamp is the header containing the Unix timestamp.
	HeaderSignatureTimestamp = "X-Signature-Timestamp"
	// HeaderSignatureEd25519 is the header containing the hex-encoded Ed25519 signature.
	HeaderSignatureEd25519 = "X-Signature-Ed25519"
)

const defaultMaxBodyBytes = 5 * 1024 * 1024 // 5 MB

// WebhookVerifier verifies webhook signatures from XBOW.
type WebhookVerifier struct {
	publicKeys   []ed25519.PublicKey
	maxClockSkew time.Duration
	maxBodyBytes int64
}

// WebhookVerifierOption configures the WebhookVerifier.
type WebhookVerifierOption func(*WebhookVerifier)

// WithMaxClockSkew sets the maximum allowed clock skew for timestamp validation.
// Default is 5 minutes.
func WithMaxClockSkew(d time.Duration) WebhookVerifierOption {
	return func(v *WebhookVerifier) {
		v.maxClockSkew = d
	}
}

// WithMaxBodyBytes sets the maximum allowed request body size in bytes.
// Default is 5 MB. Requests with bodies exceeding this limit will be rejected.
func WithMaxBodyBytes(n int64) WebhookVerifierOption {
	return func(v *WebhookVerifier) {
		v.maxBodyBytes = n
	}
}

// NewWebhookVerifier creates a new WebhookVerifier from the signing keys
// returned by MetaService.GetWebhookSigningKeys.
//
// Example:
//
//	keys, err := client.Meta.GetWebhookSigningKeys(ctx)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	verifier, err := xbow.NewWebhookVerifier(keys)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	http.Handle("/webhook", verifier.Middleware(myHandler))
func NewWebhookVerifier(keys []WebhookSigningKey, opts ...WebhookVerifierOption) (*WebhookVerifier, error) {
	if len(keys) == 0 {
		return nil, &Error{Code: "ERR_NO_KEYS", Message: "at least one signing key is required"}
	}

	v := &WebhookVerifier{
		publicKeys:   make([]ed25519.PublicKey, 0, len(keys)),
		maxClockSkew: 5 * time.Minute,
		maxBodyBytes: defaultMaxBodyBytes,
	}

	for _, opt := range opts {
		opt(v)
	}

	for _, k := range keys {
		pub, err := parsePublicKey(k.PublicKey)
		if err != nil {
			return nil, err
		}
		v.publicKeys = append(v.publicKeys, pub)
	}

	return v, nil
}

// parsePublicKey decodes a base64-encoded SPKI public key.
func parsePublicKey(b64 string) (ed25519.PublicKey, error) {
	der, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, &Error{Code: "ERR_INVALID_KEY", Message: "failed to decode base64 public key: " + err.Error()}
	}

	pub, err := x509.ParsePKIXPublicKey(der)
	if err != nil {
		return nil, &Error{Code: "ERR_INVALID_KEY", Message: "failed to parse SPKI public key: " + err.Error()}
	}

	edPub, ok := pub.(ed25519.PublicKey)
	if !ok {
		return nil, &Error{Code: "ERR_INVALID_KEY", Message: "public key is not Ed25519"}
	}

	return edPub, nil
}

// Middleware returns an http.Handler that verifies webhook signatures.
// Requests with valid signatures are passed to the next handler.
// Invalid requests receive a 401 Unauthorized response.
func (v *WebhookVerifier) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := v.Verify(r); err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Verify checks the signature and timestamp of a webhook request.
// Returns nil if valid, or an error describing the failure.
func (v *WebhookVerifier) Verify(r *http.Request) error {
	timestamp := r.Header.Get(HeaderSignatureTimestamp)
	if timestamp == "" {
		return &Error{Code: "ERR_MISSING_TIMESTAMP", Message: "missing X-Signature-Timestamp header"}
	}

	signature := r.Header.Get(HeaderSignatureEd25519)
	if signature == "" {
		return &Error{Code: "ERR_MISSING_SIGNATURE", Message: "missing X-Signature-Ed25519 header"}
	}

	ts, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return &Error{Code: "ERR_INVALID_TIMESTAMP", Message: "invalid timestamp format"}
	}

	now := time.Now().Unix()
	diff := now - ts
	if diff < 0 {
		diff = -diff
	}
	if diff > int64(v.maxClockSkew.Seconds()) {
		return &Error{Code: "ERR_TIMESTAMP_EXPIRED", Message: "timestamp outside valid range"}
	}

	sig, err := hex.DecodeString(signature)
	if err != nil {
		return &Error{Code: "ERR_INVALID_SIGNATURE", Message: "invalid signature hex encoding"}
	}
	if len(sig) != ed25519.SignatureSize {
		return &Error{Code: "ERR_INVALID_SIGNATURE", Message: "invalid signature length"}
	}

	lr := io.LimitReader(r.Body, v.maxBodyBytes+1)
	body, err := io.ReadAll(lr)
	if err != nil {
		return &Error{Code: "ERR_READ_BODY", Message: "failed to read request body"}
	}
	if int64(len(body)) > v.maxBodyBytes {
		return &Error{Code: "ERR_BODY_TOO_LARGE", Message: "request body exceeds maximum allowed size"}
	}
	r.Body = io.NopCloser(bytes.NewReader(body))

	message := append([]byte(timestamp), body...)

	for _, pub := range v.publicKeys {
		if ed25519.Verify(pub, message, sig) {
			return nil
		}
	}

	return &Error{Code: "ERR_SIGNATURE_INVALID", Message: "signature verification failed"}
}
