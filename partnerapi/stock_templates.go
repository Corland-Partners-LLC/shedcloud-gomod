package partnerapi

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// StockTemplatesService calls /partner/v1/stock-templates — buildable
// catalog designs (not physical inventory; see LotStock for on-lot units).
// Requires the partner-api.lot-stock.read scope.
type StockTemplatesService struct {
	c *Client
}

// stockTemplateListQuery flattens StockTemplateListParams for query encoding
// (Tags joined).
type stockTemplateListQuery struct {
	Page   int    `json:"page,omitempty"`
	Limit  int    `json:"limit,omitempty"`
	Search string `json:"search,omitempty"`
	Tags   string `json:"tags,omitempty"`
}

// List returns stock templates for the authenticated company. The server
// caps Limit at 60.
func (s *StockTemplatesService) List(ctx context.Context, params StockTemplateListParams) (*PaginatedResponse[StockTemplateItem], error) {
	if ctx == nil {
		ctx = context.Background()
	}
	q := stockTemplateListQuery{
		Page:   params.Page,
		Limit:  params.Limit,
		Search: params.Search,
		Tags:   strings.Join(params.Tags, ","),
	}
	var out PaginatedResponse[StockTemplateItem]
	if err := s.c.http.request(ctx, http.MethodGet, "/partner/v1/stock-templates", q, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Get returns one stock template by its work-order id.
func (s *StockTemplatesService) Get(ctx context.Context, id string) (*StockTemplateItem, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var out StockTemplateItem
	path := fmt.Sprintf("/partner/v1/stock-templates/%s", url.PathEscape(id))
	if err := s.c.http.request(ctx, http.MethodGet, path, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
