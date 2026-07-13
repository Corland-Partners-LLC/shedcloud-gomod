package partnerapi

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// DomainsService calls /partner/v1/domains (read-only). Each domain is a
// white-label storefront subdomain carrying per-location assignments:
// hidePrices, defaultForStore, and product/size mappings.
type DomainsService struct {
	c *Client
}

// List returns the company's storefront domains.
func (s *DomainsService) List(ctx context.Context, params DomainListParams) (*PaginatedResponse[DomainItem], error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var out PaginatedResponse[DomainItem]
	if err := s.c.http.request(ctx, http.MethodGet, "/partner/v1/domains", params, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ForLocation returns the domains a location is assigned to, with each
// domain's locations narrowed to that location's entry.
func (s *DomainsService) ForLocation(ctx context.Context, locationID string, params LocationDomainListParams) (*PaginatedResponse[DomainItem], error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var out PaginatedResponse[DomainItem]
	path := fmt.Sprintf("/partner/v1/locations/%s/domains", url.PathEscape(locationID))
	if err := s.c.http.request(ctx, http.MethodGet, path, params, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
