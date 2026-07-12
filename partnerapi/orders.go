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

// Update patches allowlisted sales order fields. Pass WithIfMatch(version)
// for optimistic concurrency.
func (s *OrdersService) Update(ctx context.Context, id string, body OrderPatchRequest, opts ...RequestOption) (*OrderItem, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var out OrderItem
	path := fmt.Sprintf("/partner/v1/orders/%s", url.PathEscape(id))
	if err := s.c.http.request(ctx, http.MethodPatch, path, nil, body, &out, opts...); err != nil {
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

// StatusHistory returns the order's status change log (newest first).
func (s *OrdersService) StatusHistory(ctx context.Context, id string, params PaginationParams) (*PaginatedResponse[StatusChangeItem], error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var out PaginatedResponse[StatusChangeItem]
	path := fmt.Sprintf("/partner/v1/orders/%s/status-history", url.PathEscape(id))
	if err := s.c.http.request(ctx, http.MethodGet, path, params, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// LineItems returns the order's curated line items plus the building
// configuration block.
func (s *OrdersService) LineItems(ctx context.Context, id string) (*LineItemsResponse, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var out LineItemsResponse
	path := fmt.Sprintf("/partner/v1/orders/%s/line-items", url.PathEscape(id))
	if err := s.c.http.request(ctx, http.MethodGet, path, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Contract returns the order's read-only contract signing state. Requires
// the partner-api.contracts.read scope.
func (s *OrdersService) Contract(ctx context.Context, id string) (*ContractSummary, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var out ContractSummary
	path := fmt.Sprintf("/partner/v1/orders/%s/contract", url.PathEscape(id))
	if err := s.c.http.request(ctx, http.MethodGet, path, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Payments returns the payments recorded against this order. Requires the
// partner-api.payments.read scope.
func (s *OrdersService) Payments(ctx context.Context, id string, params PaginationParams) (*PaginatedResponse[PaymentItem], error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var out PaginatedResponse[PaymentItem]
	path := fmt.Sprintf("/partner/v1/orders/%s/payments", url.PathEscape(id))
	if err := s.c.http.request(ctx, http.MethodGet, path, params, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
