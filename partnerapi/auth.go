package partnerapi

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// Auth configures how the client authenticates Partner API requests.
// Exactly one of APIKey or OAuth must be set.
type Auth struct {
	// APIKey is a long-lived sc_live_… secret sent as Bearer.
	APIKey string
	// OAuth client credentials are exchanged at POST /oauth/token.
	ClientID     string
	ClientSecret string
}

func (a Auth) validate() error {
	hasKey := strings.TrimSpace(a.APIKey) != ""
	hasOAuth := strings.TrimSpace(a.ClientID) != "" && strings.TrimSpace(a.ClientSecret) != ""
	if hasKey == hasOAuth {
		return newError("auth: set exactly one of APIKey or ClientID+ClientSecret", http.StatusUnauthorized, nil, "")
	}
	return nil
}

func (a Auth) isAPIKey() bool { return strings.TrimSpace(a.APIKey) != "" }

type cachedToken struct {
	accessToken string
	expiresAt   time.Time
}

// authProvider resolves a bearer token for Partner API requests.
type authProvider struct {
	baseURL string
	auth    Auth
	httpDo  func(*http.Request) (*http.Response, error)
	skew    time.Duration
	mu      sync.Mutex
	cached  *cachedToken
}

func newAuthProvider(baseURL string, auth Auth, httpDo func(*http.Request) (*http.Response, error), skew time.Duration) *authProvider {
	if skew <= 0 {
		skew = 60 * time.Second
	}
	return &authProvider{
		baseURL: trimTrailingSlashes(baseURL),
		auth:    auth,
		httpDo:  httpDo,
		skew:    skew,
	}
}

func (p *authProvider) getAccessToken(ctx context.Context) (string, error) {
	if p.auth.isAPIKey() {
		return p.auth.APIKey, nil
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if p.cached != nil && time.Now().Before(p.cached.expiresAt.Add(-p.skew)) {
		return p.cached.accessToken, nil
	}
	return p.refreshOAuthTokenLocked(ctx)
}

// refresh forces a new OAuth token exchange (no-op for API key auth).
func (p *authProvider) refresh(ctx context.Context) (string, error) {
	if p.auth.isAPIKey() {
		return p.auth.APIKey, nil
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.cached = nil
	return p.refreshOAuthTokenLocked(ctx)
}

func (p *authProvider) refreshOAuthTokenLocked(ctx context.Context) (string, error) {
	form := url.Values{}
	form.Set("grant_type", "client_credentials")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.baseURL+"/oauth/token", strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(p.auth.ClientID, p.auth.ClientSecret)

	resp, err := p.httpDo(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	var parsed any
	if len(raw) > 0 {
		if json.Unmarshal(raw, &parsed) != nil {
			parsed = string(raw)
		}
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		msg, code := extractError(parsed, resp.Status)
		return "", newAuthError(msg, resp.StatusCode, parsed, code)
	}

	var token OAuthTokenResponse
	if err := json.Unmarshal(raw, &token); err != nil {
		return "", newAuthError("invalid OAuth token response", resp.StatusCode, parsed, "")
	}
	if token.AccessToken == "" {
		return "", newAuthError("OAuth token response missing access_token", resp.StatusCode, parsed, "")
	}

	expiresIn := token.ExpiresIn
	if expiresIn <= 0 {
		expiresIn = 3600
	}
	p.cached = &cachedToken{
		accessToken: token.AccessToken,
		expiresAt:   time.Now().Add(time.Duration(expiresIn) * time.Second),
	}
	return token.AccessToken, nil
}

func extractError(body any, fallback string) (message, code string) {
	message = fallback
	m, ok := body.(map[string]any)
	if !ok {
		return message, ""
	}
	if s, ok := m["error_description"].(string); ok && s != "" {
		message = s
	} else if s, ok := m["error"].(string); ok && s != "" {
		message = s
	} else if s, ok := m["message"].(string); ok && s != "" {
		message = s
	}
	if s, ok := m["error"].(string); ok {
		code = s
	}
	return message, code
}
