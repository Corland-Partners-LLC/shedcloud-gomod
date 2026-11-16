package partnerapi

import (
	"context"
	"net/http"
	"time"
)

// StockSaleRequest is a Request for Sale: a salesperson asking to sell a unit
// at a store they cannot sell from.
type StockSaleRequest struct {
	ID                string `json:"id"`
	TenantID          string `json:"tenant_id"`
	StockUnitID       string `json:"stock_unit_id"`
	SourceWorkOrderID string `json:"source_work_order_id,omitempty"`
	RequesterID       string `json:"requester_id"`
	AssigneeID        string `json:"assignee_id"`
	Message           string `json:"message,omitempty"`
	// Status: pending | approved | declined.
	Status     string     `json:"status"`
	Reason     string     `json:"reason,omitempty"`
	ResolvedBy string     `json:"resolved_by,omitempty"`
	ResolvedAt *time.Time `json:"resolved_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`

	RequesterName string `json:"requester_name,omitempty"`
	AssigneeName  string `json:"assignee_name,omitempty"`
	SerialNumber  string `json:"serial_number,omitempty"`
	ProductName   string `json:"product_name,omitempty"`
	LocationName  string `json:"location_name,omitempty"`
}

// StockSaleRequestCreateRequest is POST
// /partner/v1/stock-units/{stockUnitId}/sale-requests.
type StockSaleRequestCreateRequest struct {
	AssigneeUserID string `json:"assignee_user_id,omitempty"`
	// Message is at most 500 characters; the ERP supplies a default when empty.
	Message string `json:"message,omitempty"`
}

// StockSaleRequestListParams filters GET /partner/v1/stock-sale-requests (not
// paginated).
type StockSaleRequestListParams struct {
	// Status: pending | approved | declined.
	Status      string `json:"status,omitempty"`
	StockUnitID string `json:"stock_unit_id,omitempty"`
}

// StockSaleRequestResolveRequest is POST
// /partner/v1/stock-sale-requests/{requestId}/resolve. Reason is required by
// the ERP.
type StockSaleRequestResolveRequest struct {
	// Decision: approve | decline.
	Decision string `json:"decision,omitempty"`
	Reason   string `json:"reason,omitempty"`
}

// StockSaleRequestsService calls the Request for Sale routes.
type StockSaleRequestsService struct{ c *Client }

// Create asks to sell a stock unit. Scope partner-api.stock-sale-requests.write.
// Pass WithIdempotencyKey to make retries safe.
func (s *StockSaleRequestsService) Create(ctx context.Context, stockUnitID string, body StockSaleRequestCreateRequest, opts ...RequestOption) (*StockSaleRequest, error) {
	return call[StockSaleRequest](ctx, s.c, http.MethodPost, "/partner/v1/stock-units/"+seg(stockUnitID)+"/sale-requests", nil, body, opts...)
}

// List returns the tenant's sale requests (not paginated). Scope
// partner-api.stock-sale-requests.read.
func (s *StockSaleRequestsService) List(ctx context.Context, params StockSaleRequestListParams) ([]StockSaleRequest, error) {
	return listKey[StockSaleRequest](ctx, s.c, http.MethodGet, "/partner/v1/stock-sale-requests", "sale_requests", params)
}

// Resolve approves or declines a pending request. Scope
// partner-api.stock-sale-requests.write.
func (s *StockSaleRequestsService) Resolve(ctx context.Context, requestID string, body StockSaleRequestResolveRequest) (*StockSaleRequest, error) {
	return call[StockSaleRequest](ctx, s.c, http.MethodPost, "/partner/v1/stock-sale-requests/"+seg(requestID)+"/resolve", nil, body)
}
