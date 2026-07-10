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
