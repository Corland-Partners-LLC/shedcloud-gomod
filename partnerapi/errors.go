package partnerapi

import (
	"fmt"
	"net/http"
)

// Error is a typed Partner API failure (non-2xx HTTP response or client-side
// auth/timeout problem).
//
// The ERP answers errors as {"error": "...", "code": "...", "message": "..."}.
// Message prefers the plain-language "message" and falls back to "error";
// Code is the machine code ("forbidden", "configuration_rule_violation", ...).
// Body keeps the full decoded response (for example violations[]).
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
func (e *Error) IsConflict() bool     { return e != nil && e.Status == http.StatusConflict }
func (e *Error) IsRateLimited() bool  { return e != nil && e.Status == http.StatusTooManyRequests }

// IsUnprocessable reports a 422 (invalid transition or validation, including
// "configuration_rule_violation").
func (e *Error) IsUnprocessable() bool { return e != nil && e.Status == http.StatusUnprocessableEntity }

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

// extractError reads the ERP envelope {error, code, message}. The plain-language
// "message" wins for Message; "code" is the machine code. For OAuth-style
// bodies without "code", "error" doubles as the code.
func extractError(body any, fallback string) (message, code string) {
	message = fallback
	m, ok := body.(map[string]any)
	if !ok {
		return message, ""
	}
	errStr, _ := m["error"].(string)
	switch {
	case str(m["message"]) != "":
		message = str(m["message"])
	case str(m["error_description"]) != "":
		message = str(m["error_description"])
	case errStr != "":
		message = errStr
	}
	code = str(m["code"])
	if code == "" && str(m["error_description"]) != "" {
		code = errStr
	}
	return message, code
}

func str(v any) string {
	s, _ := v.(string)
	return s
}
