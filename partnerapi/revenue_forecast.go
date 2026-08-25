package partnerapi

import (
	"context"
	"net/http"
)

// RevenueForecastService calls GET /partner/v1/revenue-forecast — read-only
// Revenue/Forecast KPIs and Authorized/Delivery week buckets. Requires
// partner-api.revenue-forecast.read.
//
// Qualification is delivery / promised dates. Money is the Sales Ledger order
// amount (base + upgrades − removed) so a related order matches
// GET /partner/v1/sales-ledger newOrderAmount. Period totals do not equal the
// Sales Ledger report for the same calendar window.
type RevenueForecastService struct {
	c *Client
}

// Get returns KPIs for the selected Period (or from/to range).
func (s *RevenueForecastService) Get(ctx context.Context, params RevenueForecastParams) (*RevenueForecastResponse, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var out RevenueForecastResponse
	if err := s.c.http.request(ctx, http.MethodGet, "/partner/v1/revenue-forecast", params, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
