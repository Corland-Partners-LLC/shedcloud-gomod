package partnerapi

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

// Product is a catalog product (/partner/v1/products). Money, quantities and
// weights are decimal strings. Only the commonly used fields are modelled;
// accounting, vendor and reorder fields are not decoded.
type Product struct {
	ID                 string         `json:"id"`
	TenantID           string         `json:"tenant_id,omitempty"`
	SKU                string         `json:"sku"`
	DisplaySKU         string         `json:"display_sku,omitempty"`
	Name               string         `json:"name"`
	Description        string         `json:"description,omitempty"`
	ProductType        string         `json:"product_type"`
	CategoryID         *string        `json:"category_id,omitempty"`
	CategoryLineID     *string        `json:"category_line_id,omitempty"`
	UomCode            string         `json:"uom_code"`
	SalesUomCode       *string        `json:"sales_uom_code,omitempty"`
	PurchaseUomCode    *string        `json:"purchase_uom_code,omitempty"`
	Tracking           string         `json:"tracking"`
	IsStocked          bool           `json:"is_stocked"`
	HideSKU            bool           `json:"hide_sku"`
	HideOnConfigurator bool           `json:"hide_on_configurator"`
	Barcode            string         `json:"barcode,omitempty"`
	MPN                string         `json:"mpn,omitempty"`
	Code               string         `json:"code,omitempty"`
	StandardCost       string         `json:"standard_cost"`
	ListPrice          string         `json:"list_price"`
	Weight             *string        `json:"weight,omitempty"`
	WeightUom          *string        `json:"weight_uom,omitempty"`
	DefaultPlantID     *string        `json:"default_plant_id,omitempty"`
	Attributes         map[string]any `json:"attributes,omitempty"`
	IsActive           bool           `json:"is_active"`
	IsSystem           bool           `json:"is_system"`
	BOMMode            string         `json:"bom_mode"`
	SortOrder          *int           `json:"sort_order,omitempty"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
}

// ProductListParams filters GET /partner/v1/products (free-text Q).
type ProductListParams struct {
	ListParams
}

// PricingResolveParams are the query params of GET /partner/v1/pricing/resolve.
type PricingResolveParams struct {
	ProductID  string `json:"product_id"`
	LocationID string `json:"location_id,omitempty"`
	Quantity   string `json:"quantity,omitempty"`
	WidthFt    string `json:"width_ft,omitempty"`
	LengthFt   string `json:"length_ft,omitempty"`
	// Date is YYYY-MM-DD (price effective date).
	Date string `json:"date,omitempty"`
}

// ProductsService calls /partner/v1/products (read) and /partner/v1/pricing.
type ProductsService struct{ c *Client }

// List returns products. Scope partner-api.products.read.
func (s *ProductsService) List(ctx context.Context, params ProductListParams) (*Page[Product], error) {
	return listPage[Product](ctx, s.c, "/partner/v1/products", "products", params)
}

// Get returns one product. Scope partner-api.products.read.
func (s *ProductsService) Get(ctx context.Context, id string) (*Product, error) {
	return call[Product](ctx, s.c, http.MethodGet, "/partner/v1/products/"+seg(id), nil, nil)
}

// ResolvePrice returns the effective sell price (a decimal string) and its
// source ("item" | "override") as a generic map. Scope partner-api.pricing.read.
func (s *ProductsService) ResolvePrice(ctx context.Context, params PricingResolveParams) (map[string]any, error) {
	out, err := call[map[string]any](ctx, s.c, http.MethodGet, "/partner/v1/pricing/resolve", params, nil)
	if err != nil {
		return nil, err
	}
	return *out, nil
}

// ---- Agreements --------------------------------------------------------------

// Agreement is a capital/RTO program agreement (/partner/v1/agreements).
// The wet-signature acceptance evidence is never exposed to partners. Nested
// rate schedules are kept as raw JSON.
type Agreement struct {
	ID                        string          `json:"id"`
	TenantID                  string          `json:"tenant_id,omitempty"`
	CounterpartyTenantID      string          `json:"counterparty_tenant_id"`
	Type                      string          `json:"type"`
	Title                     string          `json:"title"`
	Status                    string          `json:"status"`
	Terms                     string          `json:"terms,omitempty"`
	EffectiveDate             *time.Time      `json:"effective_date,omitempty"`
	ExpiryDate                *time.Time      `json:"expiry_date,omitempty"`
	Metadata                  map[string]any  `json:"metadata,omitempty"`
	AcceptedAt                *time.Time      `json:"accepted_at,omitempty"`
	RateConfigVersion         int64           `json:"rate_config_version"`
	PendingRateConfigVersion  int64           `json:"pending_rate_config_version"`
	PendingUpdate             json.RawMessage `json:"pending_update,omitempty"`
	RevokeReason              string          `json:"revoke_reason,omitempty"`
	RevokedAt                 *time.Time      `json:"revoked_at,omitempty"`
	CreatedAt                 time.Time       `json:"created_at"`
	UpdatedAt                 time.Time       `json:"updated_at"`
	DealerID                  string          `json:"dealer_id,omitempty"`
	DealerName                string          `json:"dealer_name,omitempty"`
	IsInhouse                 bool            `json:"is_inhouse"`
	OwnerName                 string          `json:"owner_name,omitempty"`
	OwnerSlug                 string          `json:"owner_slug,omitempty"`
	OwnerType                 string          `json:"owner_type,omitempty"`
	OwnerAgreementCode        string          `json:"owner_agreement_code,omitempty"`
	CounterpartyName          string          `json:"counterparty_name,omitempty"`
	CounterpartySlug          string          `json:"counterparty_slug,omitempty"`
	CounterpartyType          string          `json:"counterparty_type,omitempty"`
	CounterpartyAgreementCode string          `json:"counterparty_agreement_code,omitempty"`
	ActiveSchedule            json.RawMessage `json:"active_schedule,omitempty"`
	PendingSchedule           json.RawMessage `json:"pending_schedule,omitempty"`
	UsageCount                int64           `json:"usage_count,omitempty"`
	HasUsage                  bool            `json:"has_usage,omitempty"`
}

// AgreementListParams filters GET /partner/v1/agreements.
type AgreementListParams struct {
	ListParams
	Status string `json:"status,omitempty"`
	Type   string `json:"type,omitempty"`
}

// AgreementStateLegal is the state-specific legal text of an agreement.
// Sections and Pages stay raw JSON.
type AgreementStateLegal struct {
	ID                       string          `json:"id"`
	AgreementID              string          `json:"agreement_id"`
	StateCode                string          `json:"state_code"`
	LessorLegalName          string          `json:"lessor_legal_name,omitempty"`
	LessorAddressLine1       string          `json:"lessor_address_line1,omitempty"`
	LessorAddressLine2       string          `json:"lessor_address_line2,omitempty"`
	LessorCity               string          `json:"lessor_city,omitempty"`
	LessorRegion             string          `json:"lessor_region,omitempty"`
	LessorPostalCode         string          `json:"lessor_postal_code,omitempty"`
	Sections                 json.RawMessage `json:"sections,omitempty"`
	ArbitrationText          string          `json:"arbitration_text,omitempty"`
	AcknowledgementText      string          `json:"acknowledgement_text,omitempty"`
	IncludeCustomerDataSheet bool            `json:"include_customer_data_sheet"`
	IncludeStateRPAPages     bool            `json:"include_state_rpa_pages"`
	Pages                    json.RawMessage `json:"pages,omitempty"`
	CreatedAt                time.Time       `json:"created_at"`
	UpdatedAt                time.Time       `json:"updated_at"`
}

// AgreementsService calls /partner/v1/agreements. Scope
// partner-api.agreements.read.
type AgreementsService struct{ c *Client }

// List returns agreements.
func (s *AgreementsService) List(ctx context.Context, params AgreementListParams) (*Page[Agreement], error) {
	return listPage[Agreement](ctx, s.c, "/partner/v1/agreements", "agreements", params)
}

// Active returns only active agreements (optionally filtered by Type / Q).
func (s *AgreementsService) Active(ctx context.Context, params AgreementListParams) (*Page[Agreement], error) {
	return listPage[Agreement](ctx, s.c, "/partner/v1/agreements/active", "agreements", params)
}

// Get returns one agreement.
func (s *AgreementsService) Get(ctx context.Context, id string) (*Agreement, error) {
	return call[Agreement](ctx, s.c, http.MethodGet, "/partner/v1/agreements/"+seg(id), nil, nil)
}

// ListStateLegal lists the state legal rows of an agreement.
func (s *AgreementsService) ListStateLegal(ctx context.Context, id string) ([]AgreementStateLegal, error) {
	return listKey[AgreementStateLegal](ctx, s.c, http.MethodGet, "/partner/v1/agreements/"+seg(id)+"/state-legal", "state_legal", nil)
}

// GetStateLegal returns one state's legal row (two-letter state code).
func (s *AgreementsService) GetStateLegal(ctx context.Context, id, state string) (*AgreementStateLegal, error) {
	return call[AgreementStateLegal](ctx, s.c, http.MethodGet, "/partner/v1/agreements/"+seg(id)+"/state-legal/"+seg(state), nil, nil)
}
