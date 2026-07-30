package partnerapi

import (
	"context"
	"net/http"
)

// SalesLedgerService calls GET /partner/v1/sales-ledger — read-only sales-order
// ledger movements (same data as the portal Sales Ledger). Requires
// partner-api.sales-ledger.read.
//
// SaleDate is the calendar sale day for store / month reporting; LedgerDate is
// the status/event register and is independent of SaleDate.
type SalesLedgerService struct {
	c *Client
}

// List returns ledger movements for the authenticated company.
func (s *SalesLedgerService) List(ctx context.Context, params SalesLedgerListParams) (*PaginatedResponse[SalesLedgerItem], error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var out PaginatedResponse[SalesLedgerItem]
	if err := s.c.http.request(ctx, http.MethodGet, "/partner/v1/sales-ledger", params, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
