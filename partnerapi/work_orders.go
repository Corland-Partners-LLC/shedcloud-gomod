package partnerapi

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// WorkOrdersService calls /partner/v1/work-orders.
type WorkOrdersService struct {
	c *Client
}

// List returns work orders for the authenticated company.
func (s *WorkOrdersService) List(ctx context.Context, params WorkOrderListParams) (*PaginatedResponse[WorkOrderItem], error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var out PaginatedResponse[WorkOrderItem]
	if err := s.c.http.request(ctx, http.MethodGet, "/partner/v1/work-orders", params, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Get returns one work order by id.
func (s *WorkOrdersService) Get(ctx context.Context, id string) (*WorkOrderItem, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var out WorkOrderItem
	path := fmt.Sprintf("/partner/v1/work-orders/%s", url.PathEscape(id))
	if err := s.c.http.request(ctx, http.MethodGet, path, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Create creates a manufacturing work order (portal building-wizard parity:
// status starts in "Customer Care", the number is allocated automatically,
// and an optional SizeID attaches the product). Pass WithIdempotencyKey to
// make retries safe.
func (s *WorkOrdersService) Create(ctx context.Context, body WorkOrderCreateRequest, opts ...RequestOption) (*WorkOrderItem, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var out WorkOrderItem
	if err := s.c.http.request(ctx, http.MethodPost, "/partner/v1/work-orders", nil, body, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// Update patches allowlisted work order fields. Pass WithIfMatch(version)
// for optimistic concurrency.
func (s *WorkOrdersService) Update(ctx context.Context, id string, body WorkOrderPatchRequest, opts ...RequestOption) (*WorkOrderItem, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var out WorkOrderItem
	path := fmt.Sprintf("/partner/v1/work-orders/%s", url.PathEscape(id))
	if err := s.c.http.request(ctx, http.MethodPatch, path, nil, body, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateStatus transitions a work order's status (allowlisted transitions only).
func (s *WorkOrdersService) UpdateStatus(ctx context.Context, id string, body StatusUpdateRequest) (*WorkOrderItem, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var out WorkOrderItem
	path := fmt.Sprintf("/partner/v1/work-orders/%s/status", url.PathEscape(id))
	if err := s.c.http.request(ctx, http.MethodPost, path, nil, body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// StatusHistory returns the work order's status change log (newest first).
func (s *WorkOrdersService) StatusHistory(ctx context.Context, id string, params PaginationParams) (*PaginatedResponse[StatusChangeItem], error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var out PaginatedResponse[StatusChangeItem]
	path := fmt.Sprintf("/partner/v1/work-orders/%s/status-history", url.PathEscape(id))
	if err := s.c.http.request(ctx, http.MethodGet, path, params, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
