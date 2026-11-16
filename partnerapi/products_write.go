package partnerapi

import (
	"context"
	"net/http"
)

// ProductManufacturerUpc is one manufacturer-supplied UPC linked to a product.
type ProductManufacturerUpc struct {
	Manufacturer string `json:"manufacturer"`
	Upc          string `json:"upc"`
}

// ProductCreateRequest is POST /partner/v1/products. Unknown fields are
// rejected by the ERP (400), so only the fields below exist. Money, weights
// and quantities are decimal strings. The booleans are pointers so false can
// be sent explicitly.
type ProductCreateRequest struct {
	SKU                      string                   `json:"sku,omitempty"`
	Name                     string                   `json:"name,omitempty"`
	Description              string                   `json:"description,omitempty"`
	ProductType              string                   `json:"product_type,omitempty"`
	CategoryID               string                   `json:"category_id,omitempty"`
	CategoryLineID           string                   `json:"category_line_id,omitempty"`
	UomCode                  string                   `json:"uom_code,omitempty"`
	SalesUomCode             string                   `json:"sales_uom_code,omitempty"`
	PurchaseUomCode          string                   `json:"purchase_uom_code,omitempty"`
	Tracking                 string                   `json:"tracking,omitempty"`
	IsStocked                *bool                    `json:"is_stocked,omitempty"`
	HideSKU                  *bool                    `json:"hide_sku,omitempty"`
	HideOnConfigurator       *bool                    `json:"hide_on_configurator,omitempty"`
	IsActive                 *bool                    `json:"is_active,omitempty"`
	Barcode                  string                   `json:"barcode,omitempty"`
	MPN                      string                   `json:"mpn,omitempty"`
	Code                     string                   `json:"code,omitempty"`
	StandardCost             string                   `json:"standard_cost,omitempty"`
	ListPrice                string                   `json:"list_price,omitempty"`
	Weight                   string                   `json:"weight,omitempty"`
	WeightUom                string                   `json:"weight_uom,omitempty"`
	DefaultPlantID           string                   `json:"default_plant_id,omitempty"`
	DefaultTicketLocationID  string                   `json:"default_ticket_location_id,omitempty"`
	DefaultReceiptLocationID string                   `json:"default_receipt_location_id,omitempty"`
	InventoryAccountID       string                   `json:"inventory_account_id,omitempty"`
	COGSAccountID            string                   `json:"cogs_account_id,omitempty"`
	RevenueAccountID         string                   `json:"revenue_account_id,omitempty"`
	ReorderPoint             string                   `json:"reorder_point,omitempty"`
	ReorderQty               string                   `json:"reorder_qty,omitempty"`
	ReorderMax               string                   `json:"reorder_max,omitempty"`
	Attributes               map[string]any           `json:"attributes,omitempty"`
	ManufacturerUpcs         []ProductManufacturerUpc `json:"manufacturer_upcs,omitempty"`
}

// ProductUpdateRequest is PATCH /partner/v1/products/{productId}. Only non-nil
// fields are sent.
type ProductUpdateRequest struct {
	SKU                *string `json:"sku,omitempty"`
	Name               *string `json:"name,omitempty"`
	Description        *string `json:"description,omitempty"`
	ProductType        *string `json:"product_type,omitempty"`
	CategoryID         *string `json:"category_id,omitempty"`
	CategoryLineID     *string `json:"category_line_id,omitempty"`
	UomCode            *string `json:"uom_code,omitempty"`
	SalesUomCode       *string `json:"sales_uom_code,omitempty"`
	PurchaseUomCode    *string `json:"purchase_uom_code,omitempty"`
	Tracking           *string `json:"tracking,omitempty"`
	IsStocked          *bool   `json:"is_stocked,omitempty"`
	HideSKU            *bool   `json:"hide_sku,omitempty"`
	HideOnConfigurator *bool   `json:"hide_on_configurator,omitempty"`
	Barcode            *string `json:"barcode,omitempty"`
	MPN                *string `json:"mpn,omitempty"`
	StandardCost       *string `json:"standard_cost,omitempty"`
	ListPrice          *string `json:"list_price,omitempty"`
	Weight             *string `json:"weight,omitempty"`
}

// ProductPricingMatrixParams are the query params of
// GET /partner/v1/products/{id}/pricing-matrix.
type ProductPricingMatrixParams struct {
	LocationID string `json:"location_id,omitempty"`
}

// ProductPricingMatrixCell is one width x length cell of the effective
// pricing matrix. Prices are decimal strings.
type ProductPricingMatrixCell struct {
	WidthFt      string `json:"width_ft"`
	LengthFt     string `json:"length_ft"`
	DimensionKey string `json:"dimension_key,omitempty"`
	SellPrice    string `json:"sell_price,omitempty"`
	HasTier      bool   `json:"has_tier"`
}

// ProductPricingMatrix is the effective price matrix of a product.
type ProductPricingMatrix struct {
	ProductID     string                     `json:"product_id"`
	PriceListID   string                     `json:"price_list_id,omitempty"`
	PriceListCode string                     `json:"price_list_code,omitempty"`
	ListPrice     string                     `json:"list_price,omitempty"`
	Cells         []ProductPricingMatrixCell `json:"cells"`
}

// Create makes a product. Scope partner-api.products.write. Pass
// WithIdempotencyKey to make retries safe.
func (s *ProductsService) Create(ctx context.Context, body ProductCreateRequest, opts ...RequestOption) (*Product, error) {
	return call[Product](ctx, s.c, http.MethodPost, "/partner/v1/products", nil, body, opts...)
}

// Update patches a product. Scope partner-api.products.write.
func (s *ProductsService) Update(ctx context.Context, id string, body ProductUpdateRequest) (*Product, error) {
	return call[Product](ctx, s.c, http.MethodPatch, "/partner/v1/products/"+seg(id), nil, body)
}

// PricingMatrix returns the effective pricing matrix of a product, optionally
// for one location. Scope partner-api.products.read.
func (s *ProductsService) PricingMatrix(ctx context.Context, id string, params ProductPricingMatrixParams) (*ProductPricingMatrix, error) {
	return call[ProductPricingMatrix](ctx, s.c, http.MethodGet, "/partner/v1/products/"+seg(id)+"/pricing-matrix", params, nil)
}
