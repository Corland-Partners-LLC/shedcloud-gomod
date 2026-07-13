package partnerapi

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// ProductsService calls /partner/v1/products. Reads return finished catalog
// products (raw materials and kits are excluded); serialized on-lot units
// are served by LotStockService instead. Writes require the
// partner-api.products.write scope. Pricing/dimensions for configurable
// models live on size children — see CreateSize.
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

// Create creates a finished catalog product. Name is required; optionally
// link the product under an existing line with LineID.
func (s *ProductsService) Create(ctx context.Context, body ProductCreateRequest, opts ...RequestOption) (*ProductItem, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var out ProductItem
	if err := s.c.http.request(ctx, http.MethodPost, "/partner/v1/products", nil, body, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// Update patches allowlisted product fields.
func (s *ProductsService) Update(ctx context.Context, id string, body ProductPatchRequest) (*ProductItem, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var out ProductItem
	path := fmt.Sprintf("/partner/v1/products/%s", url.PathEscape(id))
	if err := s.c.http.request(ctx, http.MethodPatch, path, nil, body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateSize adds a size child (width/length/price) under a catalog product.
func (s *ProductsService) CreateSize(ctx context.Context, id string, body ProductSizeCreateRequest, opts ...RequestOption) (*ProductSizeItem, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var out ProductSizeItem
	path := fmt.Sprintf("/partner/v1/products/%s/sizes", url.PathEscape(id))
	if err := s.c.http.request(ctx, http.MethodPost, path, nil, body, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}
