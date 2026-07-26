package partnerapi

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// LocationBudgetsService calls /partner/v1/location-budgets.
type LocationBudgetsService struct {
	c *Client
}

// List returns paginated per-location sales budgets for the authenticated company.
func (s *LocationBudgetsService) List(ctx context.Context, params LocationBudgetListParams) (*PaginatedResponse[LocationBudgetItem], error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var out PaginatedResponse[LocationBudgetItem]
	if err := s.c.http.request(ctx, http.MethodGet, "/partner/v1/location-budgets", params, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Get returns one location sales budget by id.
func (s *LocationBudgetsService) Get(ctx context.Context, id string) (*LocationBudgetItem, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var out LocationBudgetItem
	path := fmt.Sprintf("/partner/v1/location-budgets/%s", url.PathEscape(id))
	if err := s.c.http.request(ctx, http.MethodGet, path, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Create creates a location sales budget. One budget per company, location, and year.
func (s *LocationBudgetsService) Create(ctx context.Context, body LocationBudgetUpsertRequest, opts ...RequestOption) (*LocationBudgetItem, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var out LocationBudgetItem
	if err := s.c.http.request(ctx, http.MethodPost, "/partner/v1/location-budgets", nil, body, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// Update fully replaces a location sales budget.
func (s *LocationBudgetsService) Update(ctx context.Context, id string, body LocationBudgetUpsertRequest) (*LocationBudgetItem, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var out LocationBudgetItem
	path := fmt.Sprintf("/partner/v1/location-budgets/%s", url.PathEscape(id))
	if err := s.c.http.request(ctx, http.MethodPut, path, nil, body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
