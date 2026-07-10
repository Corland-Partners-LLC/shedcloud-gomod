package partnerapi

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// CustomersService calls /partner/v1/customers.
type CustomersService struct {
	c *Client
}

// List returns customers for the authenticated company.
func (s *CustomersService) List(ctx context.Context, params CustomerListParams) (*PaginatedResponse[CustomerItem], error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var out PaginatedResponse[CustomerItem]
	if err := s.c.http.request(ctx, http.MethodGet, "/partner/v1/customers", params, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Get returns one customer by id.
func (s *CustomersService) Get(ctx context.Context, id string) (*CustomerItem, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var out CustomerItem
	path := fmt.Sprintf("/partner/v1/customers/%s", url.PathEscape(id))
	if err := s.c.http.request(ctx, http.MethodGet, path, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Create creates a customer. Email is required and must be unique within the
// company (409 on duplicates).
func (s *CustomersService) Create(ctx context.Context, body CustomerCreateRequest) (*CustomerItem, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var out CustomerItem
	if err := s.c.http.request(ctx, http.MethodPost, "/partner/v1/customers", nil, body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Update patches allowlisted customer fields (contact info, address, code,
// active flag).
func (s *CustomersService) Update(ctx context.Context, id string, body CustomerPatchRequest) (*CustomerItem, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var out CustomerItem
	path := fmt.Sprintf("/partner/v1/customers/%s", url.PathEscape(id))
	if err := s.c.http.request(ctx, http.MethodPatch, path, nil, body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
