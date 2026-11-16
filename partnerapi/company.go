package partnerapi

import (
	"context"
	"net/http"
)

// MeResponse is GET /partner/v1/me: the calling credential. Kind is
// "api_key" or "oauth". Scopes lists the granted partner-api.* scopes.
type MeResponse struct {
	TenantID     string   `json:"tenant_id"`
	CredentialID string   `json:"credential_id"`
	Kind         string   `json:"kind"`
	Scopes       []string `json:"scopes"`
}

// Me introspects the calling credential (GET /partner/v1/me). Any valid
// credential works; no scope is needed. Use it to verify a configuration and
// discover granted scopes without probing for 403s.
//
// v1.0.0 note: the ERP /me has no company block and no rate-limit block (the
// 0.x response did); read the company with Company.Get.
func (c *Client) Me(ctx context.Context) (*MeResponse, error) {
	return call[MeResponse](ctx, c, http.MethodGet, "/partner/v1/me", nil, nil)
}

// CompanyItem is GET /partner/v1/company.
type CompanyItem struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	Type     string `json:"type"`
	Status   string `json:"status,omitempty"`
	Locale   string `json:"locale,omitempty"`
	Timezone string `json:"timezone,omitempty"`
	LogoURL  string `json:"logo_url,omitempty"`
	Email    string `json:"email,omitempty"`
	Phone    string `json:"phone,omitempty"`
	Website  string `json:"website,omitempty"`
	Country  string `json:"country,omitempty"`
	Currency string `json:"currency,omitempty"`
}

// CompanyService calls /partner/v1/company.
type CompanyService struct{ c *Client }

// Get returns the authenticated company. Scope partner-api.company.read.
func (s *CompanyService) Get(ctx context.Context) (*CompanyItem, error) {
	return call[CompanyItem](ctx, s.c, http.MethodGet, "/partner/v1/company", nil, nil)
}
