// Package partnerapi is the official Go client for the ShedCloud ERP Partner
// API (`/partner/v1/*`), v1.x. Authenticate with a company-scoped API key
// (`sc_live_…`) or OAuth2 client credentials (`POST /oauth/token`).
//
// JSON is snake_case and money is a decimal string ("1250.00"). Shapes and
// routes mirror shedcloud-erp-api docs/partner/openapi.yaml and
// PARTNER_API.md.
//
// SDK v0.x talks to the Legacy Partner API (camelCase); do not mix them. See
// the migration guide in README.md.
package partnerapi

import "errors"

// ErrBaseURLRequired is returned by New when Options.BaseURL is empty.
//
// The ERP Partner API publishes no fixed production or sandbox host in its
// docs, so this SDK does not guess one: pass the host your ShedCloud contact
// gave you (for example "https://api.example.com").
var ErrBaseURLRequired = errors.New("partnerapi: Options.BaseURL is required (the ERP Partner API host is not built in)")

// ResolveBaseURL normalizes a host: it trims whitespace and trailing slashes.
func ResolveBaseURL(baseURL string) string {
	for len(baseURL) > 0 && (baseURL[0] == ' ' || baseURL[0] == '\t') {
		baseURL = baseURL[1:]
	}
	for len(baseURL) > 0 {
		c := baseURL[len(baseURL)-1]
		if c != '/' && c != ' ' && c != '\t' {
			break
		}
		baseURL = baseURL[:len(baseURL)-1]
	}
	return baseURL
}
