package partnerapi

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// LeadsService calls /partner/v1/leads.
type LeadsService struct {
	c *Client
}

// List returns leads for the authenticated company.
func (s *LeadsService) List(ctx context.Context, params SalesListParams) (*PaginatedResponse[LeadItem], error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var out PaginatedResponse[LeadItem]
	if err := s.c.http.request(ctx, http.MethodGet, "/partner/v1/leads", params, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Get returns one lead by id.
func (s *LeadsService) Get(ctx context.Context, id string) (*LeadItem, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var out LeadItem
	path := fmt.Sprintf("/partner/v1/leads/%s", url.PathEscape(id))
	if err := s.c.http.request(ctx, http.MethodGet, path, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Update patches allowlisted lead fields.
func (s *LeadsService) Update(ctx context.Context, id string, body LeadPatchRequest) (*LeadItem, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var out LeadItem
	path := fmt.Sprintf("/partner/v1/leads/%s", url.PathEscape(id))
	if err := s.c.http.request(ctx, http.MethodPatch, path, nil, body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateStatus transitions a lead's status (allowlisted transitions only).
func (s *LeadsService) UpdateStatus(ctx context.Context, id string, body StatusUpdateRequest) (*LeadItem, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var out LeadItem
	path := fmt.Sprintf("/partner/v1/leads/%s/status", url.PathEscape(id))
	if err := s.c.http.request(ctx, http.MethodPost, path, nil, body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
