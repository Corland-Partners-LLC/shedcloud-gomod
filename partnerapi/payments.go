package partnerapi

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// PaymentsService calls /partner/v1/payments — read-only payment records
// (card/ACH via Stripe plus manual entries). Raw provider payloads are never
// exposed; PaymentItem.ProviderReference carries the Stripe id for
// correlation. Requires the partner-api.payments.read scope.
type PaymentsService struct {
	c *Client
}

// List returns payments for the authenticated company.
func (s *PaymentsService) List(ctx context.Context, params PaymentListParams) (*PaginatedResponse[PaymentItem], error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var out PaginatedResponse[PaymentItem]
	if err := s.c.http.request(ctx, http.MethodGet, "/partner/v1/payments", params, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Get returns one payment by id.
func (s *PaymentsService) Get(ctx context.Context, id string) (*PaymentItem, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var out PaymentItem
	path := fmt.Sprintf("/partner/v1/payments/%s", url.PathEscape(id))
	if err := s.c.http.request(ctx, http.MethodGet, path, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
