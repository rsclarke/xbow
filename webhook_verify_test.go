package xbow

import (
	"bytes"
	"crypto/ed25519"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

func generateTestKey(t *testing.T) (ed25519.PrivateKey, string) {
	t.Helper()
	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}

	der, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		t.Fatalf("failed to marshal public key: %v", err)
	}

	return priv, base64.StdEncoding.EncodeToString(der)
}

func signRequest(priv ed25519.PrivateKey, timestamp string, body []byte) string {
	message := append([]byte(timestamp), body...)
	sig := ed25519.Sign(priv, message)
	return hex.EncodeToString(sig)
}

func TestNewWebhookVerifier(t *testing.T) {
	t.Run("requires at least one key", func(t *testing.T) {
		_, err := NewWebhookVerifier(nil)
		if err == nil {
			t.Error("expected error for empty keys")
		}
	})

	t.Run("parses valid key", func(t *testing.T) {
		_, b64 := generateTestKey(t)
		v, err := NewWebhookVerifier([]WebhookSigningKey{{PublicKey: b64}})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(v.publicKeys) != 1 {
			t.Errorf("expected 1 public key, got %d", len(v.publicKeys))
		}
	})

	t.Run("rejects invalid base64", func(t *testing.T) {
		_, err := NewWebhookVerifier([]WebhookSigningKey{{PublicKey: "not-valid-base64!!!"}})
		if err == nil {
			t.Error("expected error for invalid base64")
		}
	})

	t.Run("applies max clock skew option", func(t *testing.T) {
		_, b64 := generateTestKey(t)
		v, err := NewWebhookVerifier(
			[]WebhookSigningKey{{PublicKey: b64}},
			WithMaxClockSkew(10*time.Minute),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if v.maxClockSkew != 10*time.Minute {
			t.Errorf("expected 10 minute clock skew, got %v", v.maxClockSkew)
		}
	})
}

func TestWebhookVerifier_Verify(t *testing.T) {
	priv, b64 := generateTestKey(t)
	v, err := NewWebhookVerifier([]WebhookSigningKey{{PublicKey: b64}})
	if err != nil {
		t.Fatalf("failed to create verifier: %v", err)
	}

	t.Run("valid signature", func(t *testing.T) {
		body := []byte(`{"event":"ping"}`)
		timestamp := strconv.FormatInt(time.Now().Unix(), 10)
		sig := signRequest(priv, timestamp, body)

		req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader(body))
		req.Header.Set(HeaderSignatureTimestamp, timestamp)
		req.Header.Set(HeaderSignatureEd25519, sig)

		if err := v.Verify(req); err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		restoredBody, _ := io.ReadAll(req.Body)
		if !bytes.Equal(restoredBody, body) {
			t.Error("body was not restored after verification")
		}
	})

	t.Run("missing timestamp header", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/webhook", nil)
		req.Header.Set(HeaderSignatureEd25519, "abc123")

		err := v.Verify(req)
		if err == nil {
			t.Error("expected error for missing timestamp")
		}
		var xerr *Error
		if !errors.As(err, &xerr) || xerr.Code != "ERR_MISSING_TIMESTAMP" {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("missing signature header", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/webhook", nil)
		req.Header.Set(HeaderSignatureTimestamp, "1234567890")

		err := v.Verify(req)
		if err == nil {
			t.Error("expected error for missing signature")
		}
		var xerr *Error
		if !errors.As(err, &xerr) || xerr.Code != "ERR_MISSING_SIGNATURE" {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("invalid timestamp format", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/webhook", nil)
		req.Header.Set(HeaderSignatureTimestamp, "not-a-number")
		req.Header.Set(HeaderSignatureEd25519, "abc123")

		err := v.Verify(req)
		if err == nil {
			t.Error("expected error for invalid timestamp")
		}
		var xerr *Error
		if !errors.As(err, &xerr) || xerr.Code != "ERR_INVALID_TIMESTAMP" {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("expired timestamp", func(t *testing.T) {
		body := []byte(`{"event":"ping"}`)
		oldTime := time.Now().Add(-10 * time.Minute).Unix()
		timestamp := strconv.FormatInt(oldTime, 10)
		sig := signRequest(priv, timestamp, body)

		req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader(body))
		req.Header.Set(HeaderSignatureTimestamp, timestamp)
		req.Header.Set(HeaderSignatureEd25519, sig)

		err := v.Verify(req)
		if err == nil {
			t.Error("expected error for expired timestamp")
		}
		var xerr *Error
		if !errors.As(err, &xerr) || xerr.Code != "ERR_TIMESTAMP_EXPIRED" {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("future timestamp", func(t *testing.T) {
		body := []byte(`{"event":"ping"}`)
		futureTime := time.Now().Add(10 * time.Minute).Unix()
		timestamp := strconv.FormatInt(futureTime, 10)
		sig := signRequest(priv, timestamp, body)

		req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader(body))
		req.Header.Set(HeaderSignatureTimestamp, timestamp)
		req.Header.Set(HeaderSignatureEd25519, sig)

		err := v.Verify(req)
		if err == nil {
			t.Error("expected error for future timestamp")
		}
	})

	t.Run("invalid signature hex", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/webhook", nil)
		req.Header.Set(HeaderSignatureTimestamp, strconv.FormatInt(time.Now().Unix(), 10))
		req.Header.Set(HeaderSignatureEd25519, "not-hex!!!")

		err := v.Verify(req)
		if err == nil {
			t.Error("expected error for invalid signature hex")
		}
		var xerr *Error
		if !errors.As(err, &xerr) || xerr.Code != "ERR_INVALID_SIGNATURE" {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("wrong signature", func(t *testing.T) {
		body := []byte(`{"event":"ping"}`)
		timestamp := strconv.FormatInt(time.Now().Unix(), 10)
		wrongSig := hex.EncodeToString(make([]byte, ed25519.SignatureSize))

		req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader(body))
		req.Header.Set(HeaderSignatureTimestamp, timestamp)
		req.Header.Set(HeaderSignatureEd25519, wrongSig)

		err := v.Verify(req)
		if err == nil {
			t.Error("expected error for wrong signature")
		}
		var xerr *Error
		if !errors.As(err, &xerr) || xerr.Code != "ERR_SIGNATURE_INVALID" {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

func TestWebhookVerifier_MultipleKeys(t *testing.T) {
	priv1, b64_1 := generateTestKey(t)
	_, b64_2 := generateTestKey(t)

	v, err := NewWebhookVerifier([]WebhookSigningKey{
		{PublicKey: b64_1},
		{PublicKey: b64_2},
	})
	if err != nil {
		t.Fatalf("failed to create verifier: %v", err)
	}

	body := []byte(`{"event":"ping"}`)
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	sig := signRequest(priv1, timestamp, body)

	req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader(body))
	req.Header.Set(HeaderSignatureTimestamp, timestamp)
	req.Header.Set(HeaderSignatureEd25519, sig)

	if err := v.Verify(req); err != nil {
		t.Errorf("should accept signature from first key: %v", err)
	}
}

func TestWebhookVerifier_Middleware(t *testing.T) {
	priv, b64 := generateTestKey(t)
	v, err := NewWebhookVerifier([]WebhookSigningKey{{PublicKey: b64}})
	if err != nil {
		t.Fatalf("failed to create verifier: %v", err)
	}

	handlerCalled := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
	})

	mw := v.Middleware(handler)

	t.Run("passes valid request", func(t *testing.T) {
		handlerCalled = false
		body := []byte(`{"event":"ping"}`)
		timestamp := strconv.FormatInt(time.Now().Unix(), 10)
		sig := signRequest(priv, timestamp, body)

		req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader(body))
		req.Header.Set(HeaderSignatureTimestamp, timestamp)
		req.Header.Set(HeaderSignatureEd25519, sig)

		rr := httptest.NewRecorder()
		mw.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rr.Code)
		}
		if !handlerCalled {
			t.Error("handler was not called")
		}
	})

	t.Run("rejects invalid request", func(t *testing.T) {
		handlerCalled = false
		req := httptest.NewRequest(http.MethodPost, "/webhook", nil)

		rr := httptest.NewRecorder()
		mw.ServeHTTP(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Errorf("expected 401, got %d", rr.Code)
		}
		if handlerCalled {
			t.Error("handler should not have been called")
		}
	})
}
