package xbow

import (
	"errors"
	"fmt"
	"testing"

	"github.com/doordash-oss/oapi-codegen-dd/v3/pkg/runtime"
)

func TestErrorMessage(t *testing.T) {
	tests := []struct {
		name string
		err  Error
		want string
	}{
		{
			name: "with message",
			err:  Error{StatusCode: 404, Code: "ERR_NOT_FOUND", Message: "Assessment not found"},
			want: "xbow: Assessment not found (status=404, code=ERR_NOT_FOUND)",
		},
		{
			name: "without message",
			err:  Error{StatusCode: 500, ErrorType: "Internal Server Error"},
			want: "xbow: Internal Server Error (status=500)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("Error() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestErrorIs(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		target     error
		want       bool
	}{
		{"400 is ErrBadRequest", 400, ErrBadRequest, true},
		{"401 is ErrUnauthorized", 401, ErrUnauthorized, true},
		{"403 is ErrForbidden", 403, ErrForbidden, true},
		{"404 is ErrNotFound", 404, ErrNotFound, true},
		{"429 is ErrRateLimited", 429, ErrRateLimited, true},
		{"500 is ErrInternalServer", 500, ErrInternalServer, true},
		{"502 is ErrInternalServer", 502, ErrInternalServer, true},
		{"404 is not ErrBadRequest", 404, ErrBadRequest, false},
		{"200 is not ErrNotFound", 200, ErrNotFound, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &Error{StatusCode: tt.statusCode}
			if got := errors.Is(err, tt.target); got != tt.want {
				t.Errorf("errors.Is() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrorUnwrap(t *testing.T) {
	wrapped := errors.New("underlying error")
	err := &Error{StatusCode: 500, Wrapped: wrapped}

	if !errors.Is(err, wrapped) {
		t.Error("expected Unwrap to return wrapped error")
	}
}

func TestWrapError(t *testing.T) {
	t.Run("nil returns nil", func(t *testing.T) {
		if got := wrapError(nil); got != nil {
			t.Errorf("wrapError(nil) = %v, want nil", got)
		}
	})

	t.Run("non-ClientAPIError passes through", func(t *testing.T) {
		original := errors.New("some error")
		if got := wrapError(original); got != original {
			t.Errorf("wrapError() = %v, want %v", got, original)
		}
	})

	t.Run("ClientAPIError with JSON body extracts fields", func(t *testing.T) {
		jsonErr := fmt.Errorf(`{"code":"ERR_NOT_FOUND","error":"Not Found","message":"Assessment xyz not found"}`)
		clientErr := runtime.NewClientAPIError(jsonErr, runtime.WithStatusCode(404))

		got := wrapError(clientErr)
		apiErr, ok := got.(*Error)
		if !ok {
			t.Fatalf("expected *Error, got %T", got)
		}

		if apiErr.StatusCode != 404 {
			t.Errorf("StatusCode = %d, want 404", apiErr.StatusCode)
		}
		if apiErr.Code != "ERR_NOT_FOUND" {
			t.Errorf("Code = %q, want ERR_NOT_FOUND", apiErr.Code)
		}
		if apiErr.Message != "Assessment xyz not found" {
			t.Errorf("Message = %q, want 'Assessment xyz not found'", apiErr.Message)
		}
	})

	t.Run("ClientAPIError without JSON uses status defaults", func(t *testing.T) {
		plainErr := errors.New("connection failed")
		clientErr := runtime.NewClientAPIError(plainErr, runtime.WithStatusCode(404))

		got := wrapError(clientErr)
		apiErr, ok := got.(*Error)
		if !ok {
			t.Fatalf("expected *Error, got %T", got)
		}

		if apiErr.Code != ErrCodeNotFound {
			t.Errorf("Code = %q, want %q", apiErr.Code, ErrCodeNotFound)
		}
		if apiErr.ErrorType != "Not Found" {
			t.Errorf("ErrorType = %q, want 'Not Found'", apiErr.ErrorType)
		}
	})
}

func TestIsNotFound(t *testing.T) {
	notFoundErr := &Error{StatusCode: 404}
	otherErr := &Error{StatusCode: 500}

	if !IsNotFound(notFoundErr) {
		t.Error("IsNotFound() should return true for 404")
	}
	if IsNotFound(otherErr) {
		t.Error("IsNotFound() should return false for non-404")
	}
}

func TestIsRateLimited(t *testing.T) {
	rateLimitErr := &Error{StatusCode: 429}
	otherErr := &Error{StatusCode: 500}

	if !IsRateLimited(rateLimitErr) {
		t.Error("IsRateLimited() should return true for 429")
	}
	if IsRateLimited(otherErr) {
		t.Error("IsRateLimited() should return false for non-429")
	}
}
