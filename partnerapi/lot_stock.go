package partnerapi

import (
	"context"
	"net/http"
)

// LotStockService calls GET /partner/v1/lot-stock.
type LotStockService struct {
	c *Client
}

// List returns on-lot inventory for the authenticated company.
func (s *LotStockService) List(ctx context.Context, params LotStockListParams) (*PaginatedResponse[LotStockItem], error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var out PaginatedResponse[LotStockItem]
	if err := s.c.http.request(ctx, http.MethodGet, "/partner/v1/lot-stock", params, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
