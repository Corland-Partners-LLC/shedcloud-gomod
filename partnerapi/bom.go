package partnerapi

import (
	"context"
	"net/http"
)

// BOMService calls /partner/v1/bom (read-only, partner-api.products.read).
//
// Read-only on purpose: authoring a dimension variable is a modelling
// decision made in the portal. A partner adding one would silently widen the
// vocabulary every formula author in the company then sees.
type BOMService struct {
	c *Client
}

// DimensionVariables returns the dimension-variable keys accepted in a
// DimensionMap on order and work-order creation: the ten built-in variables
// plus any this company has defined. Resolve this list rather than
// hard-coding keys — the company-defined part varies per company.
func (s *BOMService) DimensionVariables(ctx context.Context) (*DimensionVariablesResponse, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var out DimensionVariablesResponse
	if err := s.c.http.request(ctx, http.MethodGet, "/partner/v1/bom/dimension-variables", nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
