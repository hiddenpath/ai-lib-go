// Package errors provides structured error handling.
// 错误处理，遵循 AI-Protocol 标准错误码体系。
package errors

import (
	"fmt"
	"net/http"
)

// StandardErrorCode represents a standard error code.
// 标准错误码 (E1001-E9999)。
type StandardErrorCode string

const (
	// E1xxx: Client errors
	ErrInvalidRequest    StandardErrorCode = "E1001" // 400
	ErrAuthentication    StandardErrorCode = "E1002" // 401
	ErrPermissionDenied  StandardErrorCode = "E1003" // 403
	ErrNotFound          StandardErrorCode = "E1004" // 404
	ErrMethodNotAllowed  StandardErrorCode = "E1005" // 405
	ErrConflict          StandardErrorCode = "E1006" // 409
	ErrPayloadTooLarge   StandardErrorCode = "E1007" // 413
	ErrUnsupportedMedia  StandardErrorCode = "E1015" // 415
	ErrUnprocessable     StandardErrorCode = "E1022" // 422

	// E2xxx: Rate limiting errors
	ErrRateLimited      StandardErrorCode = "E2001" // 429
	ErrQuotaExhausted   StandardErrorCode = "E2002" // 429 (quota specific)

	// E3xxx: Server errors
	ErrInternalError    StandardErrorCode = "E3001" // 500
	ErrNotImplemented   StandardErrorCode = "E3002" // 501
	ErrBadGateway       Standard ErrorCode = "E3003" // 502
	ErrServiceUnavailable StandardErrorCode = "E3004" // 503
	ErrGatewayTimeout   StandardErrorCode = "E3005" // 504
	ErrTimeout          StandardErrorCode = "E3006" // timeout

	// E4xxx: Content errors
	ErrContentFilter    StandardErrorCode = "E4001"
	ErrSensitiveContent StandardErrorCode = "E4002"

	// E5xxx: Model errors
	ErrModelNotFound    StandardErrorCode = "E5001"
	ErrModelOverloaded  StandardErrorCode = "E5002"
	ErrContextLengthExceeded StandardErrorCode = "E5003"

	// E9xxx: Transport errors
	ErrNetworkError     StandardErrorCode = "E9001"
	ErrConnectionFailed StandardErrorCode = "E9002"
	ErrSSLError         StandardErrorCode = "E9003"
)

// Error represents a structured error.
type Error struct {
	Code       StandardErrorCode `json:"code"`
	StatusCode int               `json:"status_code"`
	Message    string            `json:"message"`
	Provider   string            `json:"provider,omitempty"`
	RequestID  string            `json:"request_id,omitempty"`
	Details    map[string]any    `json:"details,omitempty"`
	Retryable  bool              `json:"retryable"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// NewError creates a new error.
func NewError(code StandardErrorCode, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

// WithStatusCode sets the HTTP status code.
func (e *Error) WithStatusCode(status int) *Error {
	e.StatusCode = status
	return e
}

// WithProvider sets the provider name.
func (e *Error) WithProvider(provider string) *Error {
	e.Provider = provider
	return e
}

// WithRequestID sets the request ID.
func (e *Error) WithRequestID(id string) *Error {
	e.RequestID = id
	return e
}

// WithDetails sets additional details.
func (e *Error) WithDetails(details map[string]any) *Error {
	e.Details = details
	return e
}

// WithRetryable sets whether the error is retryable.
func (e *Error) WithRetryable(retryable bool) *Error {
	e.Retryable = retryable
	return e
}

// IsRetryable returns whether the error is retryable.
func (e *Error) IsRetryable() bool {
	return e.Retryable
}

// Classification classifies errors by HTTP status code.
type Classification struct {
	ByHTTPStatus map[int]StandardErrorCode
	ByProviderCode map[string]StandardErrorCode
}

// DefaultClassification returns the default error classification.
func DefaultClassification() *Classification {
	return &Classification{
		ByHTTPStatus: map[int]StandardErrorCode{
			http.StatusBadRequest:          ErrInvalidRequest,
			http.StatusUnauthorized:        ErrAuthentication,
			http.StatusForbidden:           ErrPermissionDenied,
			http.StatusNotFound:            ErrNotFound,
			http.StatusMethodNotAllowed:    ErrMethodNotAllowed,
			http.StatusConflict:            ErrConflict,
			http.StatusRequestEntityTooLarge: ErrPayloadTooLarge,
			http.StatusUnsupportedMediaType: ErrUnsupportedMedia,
			http.StatusUnprocessableEntity: ErrUnprocessable,
			http.StatusTooManyRequests:     ErrRateLimited,
			http.StatusInternalServerError: ErrInternalError,
			http.StatusNotImplemented:      ErrNotImplemented,
			http.StatusBadGateway:          ErrBadGateway,
			http.StatusServiceUnavailable:  ErrServiceUnavailable,
			http.StatusGatewayTimeout:      ErrGatewayTimeout,
		},
		ByProviderCode: make(map[string]StandardErrorCode),
	}
}

// Classify classifies an HTTP status code to a standard error code.
func (c *Classification) Classify(statusCode int) StandardErrorCode {
	if code, ok := c.ByHTTPStatus[statusCode]; ok {
		return code
	}
	if statusCode >= 500 {
		return ErrInternalError
	}
	if statusCode >= 400 {
		return ErrInvalidRequest
	}
	return ErrInternalError
}

// IsRetryable returns whether the status code indicates a retryable error.
func IsRetryable(statusCode int) bool {
	return statusCode == http.StatusTooManyRequests ||
		statusCode == http.StatusInternalServerError ||
		statusCode == http.StatusBadGateway ||
		statusCode == http.StatusServiceUnavailable ||
		statusCode == http.StatusGatewayTimeout
}

// IsAuthError returns whether the error is an authentication error.
func IsAuthError(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.Code == ErrAuthentication || e.Code == ErrPermissionDenied
	}
	return false
}

// IsRateLimited returns whether the error is a rate limit error.
func IsRateLimited(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.Code == ErrRateLimited || e.Code == ErrQuotaExhausted
	}
	return false
}

// IsTimeout returns whether the error is a timeout error.
func IsTimeout(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.Code == ErrTimeout || e.Code == ErrGatewayTimeout
	}
	return false
}

// IsServerError returns whether the error is a server error.
func IsServerError(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.StatusCode >= 500
	}
	return false
}

// Wrap wraps an error with context.
func Wrap(err error, code StandardErrorCode, message string) *Error {
	if e, ok := err.(*Error); ok {
		return e
	}
	return NewError(code, message)
}

// WrapWithCode wraps an error with a code.
func WrapWithCode(err error, code StandardErrorCode) *Error {
	if e, ok := err.(*Error); ok {
		return e
	}
	return NewError(code, err.Error())
}
