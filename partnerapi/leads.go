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

// Create makes a lead at a sales location with the customer's contact
// information. When no salesperson is given, the location's lead-routing
// strategy (round-robin, availability, skill-based) auto-assigns one — the
// same routing in-stock quote creation uses. Requires the
// partner-api.leads.write scope.
func (s *LeadsService) Create(ctx context.Context, body LeadCreateRequest, opts ...RequestOption) (*LeadItem, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var out LeadItem
	if err := s.c.http.request(ctx, http.MethodPost, "/partner/v1/leads", nil, body, &out, opts...); err != nil {
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

// Update patches allowlisted lead fields. Pass WithIfMatch(version) for
// optimistic concurrency.
func (s *LeadsService) Update(ctx context.Context, id string, body LeadPatchRequest, opts ...RequestOption) (*LeadItem, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var out LeadItem
	path := fmt.Sprintf("/partner/v1/leads/%s", url.PathEscape(id))
	if err := s.c.http.request(ctx, http.MethodPatch, path, nil, body, &out, opts...); err != nil {
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

// StatusHistory returns the lead's status change log (newest first).
func (s *LeadsService) StatusHistory(ctx context.Context, id string, params PaginationParams) (*PaginatedResponse[StatusChangeItem], error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var out PaginatedResponse[StatusChangeItem]
	path := fmt.Sprintf("/partner/v1/leads/%s/status-history", url.PathEscape(id))
	if err := s.c.http.request(ctx, http.MethodGet, path, params, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
