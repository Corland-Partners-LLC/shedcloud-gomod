package partnerapi

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// QuotesService calls /partner/v1/quotes.
type QuotesService struct {
	c *Client
}

// List returns quotes for the authenticated company.
func (s *QuotesService) List(ctx context.Context, params QuoteListParams) (*PaginatedResponse[QuoteItem], error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var out PaginatedResponse[QuoteItem]
	if err := s.c.http.request(ctx, http.MethodGet, "/partner/v1/quotes", params, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Create makes a quote from an in-stock (on-lot) unit, referenced by serial
// number or work order id. The unit's work order is linked to the new quote
// (serial number, product lines, and pricing are copied), the sales location
// is assigned, and lead routing auto-assigns a salesperson when the location
// has routing configured. Requires the partner-api.quotes.write scope; the
// API responds 409 when the unit is already taken.
func (s *QuotesService) Create(ctx context.Context, body QuoteCreateRequest, opts ...RequestOption) (*QuoteItem, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var out QuoteItem
	if err := s.c.http.request(ctx, http.MethodPost, "/partner/v1/quotes", nil, body, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// Convert places a quote as a sales order — the same conversion the portal's
// "Place Order" button runs. The quote's product lines and configurator are
// cloned onto the new order, the quote is marked Sold with a back-reference,
// and a linked work order moves with the conversion. The new order starts in
// the Unsubmitted status. Requires the partner-api.orders.write scope (the
// endpoint creates a sales order); the API responds 409 when the quote is
// already Sold/Cancelled/Deleted or an active order already exists for it.
func (s *QuotesService) Convert(ctx context.Context, id string, body QuoteConvertRequest, opts ...RequestOption) (*OrderItem, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var out OrderItem
	path := fmt.Sprintf("/partner/v1/quotes/%s/convert", url.PathEscape(id))
	if err := s.c.http.request(ctx, http.MethodPost, path, nil, body, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// Get returns one quote by id.
func (s *QuotesService) Get(ctx context.Context, id string) (*QuoteItem, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var out QuoteItem
	path := fmt.Sprintf("/partner/v1/quotes/%s", url.PathEscape(id))
	if err := s.c.http.request(ctx, http.MethodGet, path, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Update patches allowlisted quote fields.
func (s *QuotesService) Update(ctx context.Context, id string, body QuotePatchRequest) (*QuoteItem, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var out QuoteItem
	path := fmt.Sprintf("/partner/v1/quotes/%s", url.PathEscape(id))
	if err := s.c.http.request(ctx, http.MethodPatch, path, nil, body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateStatus transitions a quote's status (allowlisted transitions only).
func (s *QuotesService) UpdateStatus(ctx context.Context, id string, body StatusUpdateRequest) (*QuoteItem, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var out QuoteItem
	path := fmt.Sprintf("/partner/v1/quotes/%s/status", url.PathEscape(id))
	if err := s.c.http.request(ctx, http.MethodPost, path, nil, body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
