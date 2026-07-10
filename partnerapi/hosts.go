// Package partnerapi is the official Go client for the ShedCloud Partner API
// (`/partner/v1/*`). Authenticate with a company-scoped API key (`sc_live_…`)
// or OAuth2 client credentials (`POST /oauth/token`).
//
// Shapes and routes mirror shedcloud-api-go/docs/PARTNER_API.md and the
// TypeScript SDK at github.com/Corland-Partners-LLC/shedcloud-npm.
package partnerapi

// Built-in Partner API hosts. Prefer Environment over hardcoding URLs.
const (
	HostProduction = "https://go.shedcloud.com"
	HostSandbox    = "https://api.shedcloudtest.com"
)

// Environment selects a built-in Partner API host.
type Environment string

const (
	EnvironmentProduction Environment = "production"
	EnvironmentSandbox    Environment = "sandbox"
)

// DefaultBaseURL is used when neither Environment nor BaseURL is set.
const DefaultBaseURL = HostProduction

// ResolveBaseURL picks the API host from client options.
// Precedence: explicit BaseURL > Environment > production default.
func ResolveBaseURL(baseURL string, env Environment) string {
	if baseURL = trimTrailingSlashes(baseURL); baseURL != "" {
		return baseURL
	}
	switch env {
	case EnvironmentSandbox:
		return HostSandbox
	case EnvironmentProduction, "":
		return HostProduction
	default:
		return HostProduction
	}
}

func trimTrailingSlashes(s string) string {
	for len(s) > 0 && s[len(s)-1] == '/' {
		s = s[:len(s)-1]
	}
	return s
}
