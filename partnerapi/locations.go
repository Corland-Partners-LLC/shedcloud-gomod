package partnerapi

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// LocationsService calls /partner/v1/locations.
type LocationsService struct {
	c *Client
}

// List returns locations (sales lots, plants, warehouses) for the
// authenticated company.
func (s *LocationsService) List(ctx context.Context, params LocationListParams) (*PaginatedResponse[LocationItem], error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var out PaginatedResponse[LocationItem]
	if err := s.c.http.request(ctx, http.MethodGet, "/partner/v1/locations", params, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Get returns one location by id.
func (s *LocationsService) Get(ctx context.Context, id string) (*LocationItem, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var out LocationItem
	path := fmt.Sprintf("/partner/v1/locations/%s", url.PathEscape(id))
	if err := s.c.http.request(ctx, http.MethodGet, path, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Create creates a location. Name is required, plus either an address or a
// latitude/longitude pair. The location code must be unique in the company
// (409 on duplicates).
func (s *LocationsService) Create(ctx context.Context, body LocationCreateRequest) (*LocationItem, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var out LocationItem
	if err := s.c.http.request(ctx, http.MethodPost, "/partner/v1/locations", nil, body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Update patches allowlisted location fields.
func (s *LocationsService) Update(ctx context.Context, id string, body LocationPatchRequest) (*LocationItem, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var out LocationItem
	path := fmt.Sprintf("/partner/v1/locations/%s", url.PathEscape(id))
	if err := s.c.http.request(ctx, http.MethodPatch, path, nil, body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
