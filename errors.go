package xbow

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/doordash-oss/oapi-codegen-dd/v3/pkg/runtime"
)

// apiErrorEnvelope is used to extract structured error info from API responses.
type apiErrorEnvelope struct {
	Code    string `json:"code"`
	Error   string `json:"error"`
	Message string `json:"message"`
}

// Error codes returned by the API.
const (
	ErrCodeValidation     = "FST_ERR_VALIDATION"
	ErrCodeNotFound       = "ERR_NOT_FOUND"
	ErrCodeQuotaExhausted = "ERR_QUOTA_EXHAUSTED"
)

// Sentinel errors for use with errors.Is.
var (
	ErrNotFound       = errors.New("resource not found")
	ErrBadRequest     = errors.New("bad request")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrForbidden      = errors.New("forbidden")
	ErrRateLimited    = errors.New("rate limited")
	ErrInternalServer = errors.New("internal server error")

	// Client-side configuration errors.
	ErrMissingOrgKey         = errors.New("xbow: organization key is required")
	ErrMissingIntegrationKey = errors.New("xbow: integration key is required")
	ErrMissingAnyKey         = errors.New("xbow: organization key or integration key is required")
)

// Error represents an API error response.
type Error struct {
	StatusCode int    `json:"-"`
	Code       string `json:"code"`
	ErrorType  string `json:"error"`
	Message    string `json:"message"`
	Wrapped    error  `json:"-"`
}

func (e *Error) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("xbow: %s (status=%d, code=%s)", e.Message, e.StatusCode, e.Code)
	}
	return fmt.Sprintf("xbow: %s (status=%d)", e.ErrorType, e.StatusCode)
}

// Unwrap returns the wrapped error.
func (e *Error) Unwrap() error {
	return e.Wrapped
}

// Is implements errors.Is for API errors.
func (e *Error) Is(target error) bool {
	switch {
	case errors.Is(target, ErrBadRequest) && e.StatusCode == 400:
		return true
	case errors.Is(target, ErrUnauthorized) && e.StatusCode == 401:
		return true
	case errors.Is(target, ErrForbidden) && e.StatusCode == 403:
		return true
	case errors.Is(target, ErrNotFound) && e.StatusCode == 404:
		return true
	case errors.Is(target, ErrRateLimited) && e.StatusCode == 429:
		return true
	case errors.Is(target, ErrInternalServer) && e.StatusCode >= 500:
		return true
	}
	return false
}

// IsNotFound returns true if the error is a 404 Not Found error.
func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

// IsRateLimited returns true if the error is a 429 Rate Limited error.
func IsRateLimited(err error) bool {
	return errors.Is(err, ErrRateLimited)
}

// wrapError converts a generated client error to our Error type.
func wrapError(err error) error {
	if err == nil {
		return nil
	}

	var clientErr *runtime.ClientAPIError
	if errors.As(err, &clientErr) {
		apiErr := &Error{
			StatusCode: clientErr.StatusCode(),
			Wrapped:    err,
		}

		// Try to extract structured error info from the wrapped error's message.
		// The generated client wraps parsed error responses that contain code/error/message.
		if parsed := parseAPIError(clientErr.Unwrap()); parsed != nil {
			apiErr.Code = parsed.Code
			apiErr.ErrorType = parsed.Error
			apiErr.Message = parsed.Message
		} else {
			// Fall back to status-based defaults
			switch apiErr.StatusCode {
			case 400:
				apiErr.ErrorType = "Bad Request"
				apiErr.Code = ErrCodeValidation
			case 401:
				apiErr.ErrorType = "Unauthorized"
			case 403:
				apiErr.ErrorType = "Forbidden"
			case 404:
				apiErr.ErrorType = "Not Found"
				apiErr.Code = ErrCodeNotFound
			case 429:
				apiErr.ErrorType = "Too Many Requests"
			default:
				if apiErr.StatusCode >= 500 {
					apiErr.ErrorType = "Internal Server Error"
				}
			}
			apiErr.Message = err.Error()
		}

		return apiErr
	}

	return err
}

// wrapRawError creates a structured *Error from a raw HTTP response status and body.
// It mirrors the logic in wrapError but works without a runtime.ClientAPIError.
func wrapRawError(statusCode int, body []byte) *Error {
	apiErr := &Error{
		StatusCode: statusCode,
	}

	var envelope apiErrorEnvelope
	if json.Unmarshal(body, &envelope) == nil && envelope.Code != "" {
		apiErr.Code = envelope.Code
		apiErr.ErrorType = envelope.Error
		apiErr.Message = envelope.Message
	} else {
		switch statusCode {
		case 400:
			apiErr.ErrorType = "Bad Request"
			apiErr.Code = ErrCodeValidation
		case 401:
			apiErr.ErrorType = "Unauthorized"
		case 403:
			apiErr.ErrorType = "Forbidden"
		case 404:
			apiErr.ErrorType = "Not Found"
			apiErr.Code = ErrCodeNotFound
		case 429:
			apiErr.ErrorType = "Too Many Requests"
		default:
			if statusCode >= 500 {
				apiErr.ErrorType = "Internal Server Error"
			}
		}
		apiErr.Message = string(body)
	}

	return apiErr
}

// parseAPIError attempts to extract structured error info from an error.
// It handles both JSON-formatted error messages and typed error responses.
func parseAPIError(err error) *apiErrorEnvelope {
	if err == nil {
		return nil
	}

	// Try to unmarshal the error message as JSON (for fmt.Errorf wrapped errors)
	msg := err.Error()
	var envelope apiErrorEnvelope
	if json.Unmarshal([]byte(msg), &envelope) == nil && envelope.Code != "" {
		return &envelope
	}

	return nil
}
