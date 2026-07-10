package partnerapi

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// OrdersService calls /partner/v1/orders.
type OrdersService struct {
	c *Client
}

// List returns sales orders for the authenticated company.
func (s *OrdersService) List(ctx context.Context, params OrderListParams) (*PaginatedResponse[OrderItem], error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var out PaginatedResponse[OrderItem]
	if err := s.c.http.request(ctx, http.MethodGet, "/partner/v1/orders", params, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Get returns one sales order by id.
func (s *OrdersService) Get(ctx context.Context, id string) (*OrderItem, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var out OrderItem
	path := fmt.Sprintf("/partner/v1/orders/%s", url.PathEscape(id))
	if err := s.c.http.request(ctx, http.MethodGet, path, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Update patches allowlisted sales order fields.
func (s *OrdersService) Update(ctx context.Context, id string, body OrderPatchRequest) (*OrderItem, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var out OrderItem
	path := fmt.Sprintf("/partner/v1/orders/%s", url.PathEscape(id))
	if err := s.c.http.request(ctx, http.MethodPatch, path, nil, body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateStatus transitions a sales order's status (allowlisted transitions only;
// Processed is blocked via the Partner API).
func (s *OrdersService) UpdateStatus(ctx context.Context, id string, body StatusUpdateRequest) (*OrderItem, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var out OrderItem
	path := fmt.Sprintf("/partner/v1/orders/%s/status", url.PathEscape(id))
	if err := s.c.http.request(ctx, http.MethodPost, path, nil, body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
