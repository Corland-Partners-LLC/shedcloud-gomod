package partnerapi

import (
	"fmt"
	"net/http"
)

// Error is a typed Partner API failure (non-2xx HTTP response or client-side
// auth/timeout problem).
type Error struct {
	Message string
	Status  int
	Body    any
	Code    string
}

func (e *Error) Error() string {
	if e == nil {
		return "partnerapi: <nil>"
	}
	if e.Code != "" {
		return fmt.Sprintf("partnerapi: %s (status=%d code=%s)", e.Message, e.Status, e.Code)
	}
	return fmt.Sprintf("partnerapi: %s (status=%d)", e.Message, e.Status)
}

func (e *Error) IsUnauthorized() bool { return e != nil && e.Status == http.StatusUnauthorized }
func (e *Error) IsForbidden() bool    { return e != nil && e.Status == http.StatusForbidden }
func (e *Error) IsNotFound() bool     { return e != nil && e.Status == http.StatusNotFound }
func (e *Error) IsRateLimited() bool  { return e != nil && e.Status == http.StatusTooManyRequests }

// AuthError is returned when OAuth token exchange fails.
type AuthError struct {
	Message string
	Status  int
	Body    any
	Code    string
}

func (e *AuthError) Error() string {
	if e == nil {
		return "partnerapi auth: <nil>"
	}
	if e.Code != "" {
		return fmt.Sprintf("partnerapi auth: %s (status=%d code=%s)", e.Message, e.Status, e.Code)
	}
	return fmt.Sprintf("partnerapi auth: %s (status=%d)", e.Message, e.Status)
}

func (e *AuthError) IsUnauthorized() bool { return e != nil && e.Status == http.StatusUnauthorized }

func newError(message string, status int, body any, code string) *Error {
	return &Error{Message: message, Status: status, Body: body, Code: code}
}

func newAuthError(message string, status int, body any, code string) *AuthError {
	return &AuthError{Message: message, Status: status, Body: body, Code: code}
}
