package partnerapi

import (
	"context"
	"net/http"
	"time"
)

// FreightBill is a carrier freight bill (/partner/v1/freight-bills). Money is
// a decimal string.
type FreightBill struct {
	ID               string    `json:"id"`
	TenantID         string    `json:"tenant_id,omitempty"`
	Number           string    `json:"number,omitempty"`
	CarrierID        *string   `json:"carrier_id,omitempty"`
	RunID            *string   `json:"run_id,omitempty"`
	TransportOrderID *string   `json:"transport_order_id,omitempty"`
	OfferID          *string   `json:"offer_id,omitempty"`
	StripePayoutID   string    `json:"stripe_payout_id,omitempty"`
	PayoutStatus     string    `json:"payout_status,omitempty"`
	InvoiceAmount    string    `json:"invoice_amount"`
	RatedAmount      string    `json:"rated_amount"`
	VarianceAmount   string    `json:"variance_amount"`
	Status           string    `json:"status"`
	CreatedBy        *string   `json:"created_by,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// FreightBillListParams filters GET /partner/v1/freight-bills.
type FreightBillListParams struct {
	ListParams
	Status string `json:"status,omitempty"`
}

// FreightBillRequest is the body of POST /partner/v1/freight-bills. Money is a
// decimal string.
type FreightBillRequest struct {
	CarrierID        string `json:"carrier_id,omitempty"`
	OfferID          string `json:"offer_id,omitempty"`
	RunID            string `json:"run_id,omitempty"`
	TransportOrderID string `json:"transport_order_id,omitempty"`
	InvoiceAmount    string `json:"invoice_amount,omitempty"`
	RatedAmount      string `json:"rated_amount,omitempty"`
	Status           string `json:"status,omitempty"`
}

// FreightBillsService calls /partner/v1/freight-bills.
type FreightBillsService struct{ c *Client }

// List returns freight bills (filter by Status). Scope
// partner-api.freight-bills.read.
func (s *FreightBillsService) List(ctx context.Context, params FreightBillListParams) (*Page[FreightBill], error) {
	return listPage[FreightBill](ctx, s.c, "/partner/v1/freight-bills", "freight_bills", params)
}

// Get returns one freight bill. Scope partner-api.freight-bills.read.
func (s *FreightBillsService) Get(ctx context.Context, billID string) (*FreightBill, error) {
	return call[FreightBill](ctx, s.c, http.MethodGet, "/partner/v1/freight-bills/"+seg(billID), nil, nil)
}

// Create records a freight bill. Scope partner-api.freight-bills.write.
func (s *FreightBillsService) Create(ctx context.Context, body FreightBillRequest, opts ...RequestOption) (*FreightBill, error) {
	return call[FreightBill](ctx, s.c, http.MethodPost, "/partner/v1/freight-bills", nil, body, opts...)
}

// FixedAsset is a depreciable fixed asset (/partner/v1/fixed-assets). Dates
// are YYYY-MM-DD; money is a decimal string.
type FixedAsset struct {
	ID                 string    `json:"id"`
	TenantID           string    `json:"tenant_id,omitempty"`
	Name               string    `json:"name"`
	AssetAccountID     *string   `json:"asset_account_id,omitempty"`
	AccumDeprAccountID *string   `json:"accum_depr_account_id,omitempty"`
	ExpenseAccountID   *string   `json:"expense_account_id,omitempty"`
	PlacedInService    *string   `json:"placed_in_service,omitempty"`
	Cost               string    `json:"cost"`
	Salvage            string    `json:"salvage"`
	LifeMonths         int       `json:"life_months"`
	Method             string    `json:"method"`
	Status             string    `json:"status"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// FixedAssetListParams filters GET /partner/v1/fixed-assets.
type FixedAssetListParams struct {
	ListParams
	Status string `json:"status,omitempty"`
}

// FixedAssetRequest is the body of POST /partner/v1/fixed-assets.
// PlacedInService is YYYY-MM-DD; Cost and Salvage are decimal strings.
type FixedAssetRequest struct {
	Name               string `json:"name,omitempty"`
	AssetAccountID     string `json:"asset_account_id,omitempty"`
	AccumDeprAccountID string `json:"accum_depr_account_id,omitempty"`
	ExpenseAccountID   string `json:"expense_account_id,omitempty"`
	PlacedInService    string `json:"placed_in_service,omitempty"`
	Cost               string `json:"cost,omitempty"`
	Salvage            string `json:"salvage,omitempty"`
	LifeMonths         int    `json:"life_months,omitempty"`
}

// FixedAssetsService calls /partner/v1/fixed-assets.
type FixedAssetsService struct{ c *Client }

// List returns fixed assets (filter by Status). Scope
// partner-api.fixed-assets.read.
func (s *FixedAssetsService) List(ctx context.Context, params FixedAssetListParams) (*Page[FixedAsset], error) {
	return listPage[FixedAsset](ctx, s.c, "/partner/v1/fixed-assets", "fixed_assets", params)
}

// Get returns one fixed asset. Scope partner-api.fixed-assets.read.
func (s *FixedAssetsService) Get(ctx context.Context, assetID string) (*FixedAsset, error) {
	return call[FixedAsset](ctx, s.c, http.MethodGet, "/partner/v1/fixed-assets/"+seg(assetID), nil, nil)
}

// Create registers a fixed asset. Scope partner-api.fixed-assets.write.
func (s *FixedAssetsService) Create(ctx context.Context, body FixedAssetRequest, opts ...RequestOption) (*FixedAsset, error) {
	return call[FixedAsset](ctx, s.c, http.MethodPost, "/partner/v1/fixed-assets", nil, body, opts...)
}
