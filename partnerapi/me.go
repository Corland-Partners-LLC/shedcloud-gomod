package partnerapi

import (
	"context"
	"net/http"
)

// Me introspects the calling credential (GET /partner/v1/me): company,
// credential type, granted scopes, and rate-limit configuration. Requires
// any valid Partner API credential — no specific scope. Useful to verify a
// configuration (and discover granted scopes) without probing endpoints
// for 403s.
func (c *Client) Me(ctx context.Context) (*MeResponse, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var out MeResponse
	if err := c.http.request(ctx, http.MethodGet, "/partner/v1/me", nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
