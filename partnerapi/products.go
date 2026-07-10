package partnerapi

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// ProductsService calls /partner/v1/products. Catalog products (models,
// components, upgrades) are read-only — serialized on-lot units are served by
// LotStockService instead.
type ProductsService struct {
	c *Client
}

// List returns catalog products for the authenticated company.
func (s *ProductsService) List(ctx context.Context, params ProductListParams) (*PaginatedResponse[ProductItem], error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var out PaginatedResponse[ProductItem]
	if err := s.c.http.request(ctx, http.MethodGet, "/partner/v1/products", params, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Get returns one catalog product by id.
func (s *ProductsService) Get(ctx context.Context, id string) (*ProductItem, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var out ProductItem
	path := fmt.Sprintf("/partner/v1/products/%s", url.PathEscape(id))
	if err := s.c.http.request(ctx, http.MethodGet, path, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
