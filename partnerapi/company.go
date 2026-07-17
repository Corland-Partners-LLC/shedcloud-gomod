package partnerapi

import (
	"context"
	"net/http"
)

// CompanyService calls GET /partner/v1/company (read-only).
type CompanyService struct {
	c *Client
}

// Get returns the authenticated company's profile.
func (s *CompanyService) Get(ctx context.Context) (*CompanyItem, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var out CompanyItem
	if err := s.c.http.request(ctx, http.MethodGet, "/partner/v1/company", nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
