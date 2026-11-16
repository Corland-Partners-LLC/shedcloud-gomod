package partnerapi

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

// StockUnit is one lot-stock unit (/partner/v1/stock-units). It replaces the
// Legacy LotStockItem. Money is a decimal string.
type StockUnit struct {
	ID                      string `json:"id"`
	TenantID                string `json:"tenant_id,omitempty"`
	ProductID               string `json:"product_id"`
	SerialNumber            string `json:"serial_number,omitempty"`
	LotNumber               string `json:"lot_number,omitempty"`
	SourceWorkOrderID       string `json:"source_work_order_id,omitempty"`
	ConfigurationID         string `json:"configuration_id,omitempty"`
	ConfigurationRevisionID string `json:"configuration_revision_id,omitempty"`
	LocationID              string `json:"location_id,omitempty"`
	// Status: received | available | reserved | sold (partner reads hide others).
	Status string `json:"status"`
	// Origin is how the unit entered stock (for example built or return).
	Origin      string          `json:"origin"`
	UnitCost    string          `json:"unit_cost"`
	AskingPrice string          `json:"asking_price,omitempty"`
	Specs       json.RawMessage `json:"specs,omitempty"`
	SoldQuoteID string          `json:"sold_quote_id,omitempty"`
	SoldAt      *time.Time      `json:"sold_at,omitempty"`
	// PartnerSoldVisibleDays overrides, per unit, how long a sold unit stays listed.
	PartnerSoldVisibleDays *int   `json:"partner_sold_visible_days,omitempty"`
	PartnerKeepVisible     bool   `json:"partner_keep_visible"`
	HideFromEmbed          bool   `json:"hide_from_embed"`
	DiscountReason         string `json:"discount_reason,omitempty"`
	ConditionGrade         string `json:"condition_grade,omitempty"`

	ReturnedCustomerID     string `json:"returned_customer_id,omitempty"`
	ReturnedOrderID        string `json:"returned_order_id,omitempty"`
	ReturnedRTOAgreementID string `json:"returned_rto_agreement_id,omitempty"`
	CreatedBy              string `json:"created_by,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	ProductSKU      string   `json:"product_sku,omitempty"`
	ProductName     string   `json:"product_name,omitempty"`
	LocationCode    string   `json:"location_code,omitempty"`
	StoreLocationID string   `json:"store_location_id,omitempty"`
	LocationName    string   `json:"location_name,omitempty"`
	LocationState   string   `json:"location_state,omitempty"`
	SourceWONumber  string   `json:"source_wo_number,omitempty"`
	WidthFt         *float64 `json:"width_ft,omitempty"`
	LengthFt        *float64 `json:"length_ft,omitempty"`
	QuoteCount      int      `json:"quote_count,omitempty"`
	OrderCount      int      `json:"order_count,omitempty"`
	// SaleRequest is the open Request-for-Sale summary, when any.
	SaleRequest                json.RawMessage `json:"sale_request,omitempty"`
	ReturnedCustomerName       string          `json:"returned_customer_name,omitempty"`
	ReturnedOrderNumber        string          `json:"returned_order_number,omitempty"`
	ReturnedRTOAgreementNumber string          `json:"returned_rto_agreement_number,omitempty"`
}

// StockUnitListParams filters GET /partner/v1/stock-units.
//
// By default the list hides sold units past their visibility window and units
// still "received" (Receiving). Pass IncludeExpired / IncludeReceived to see
// them.
type StockUnitListParams struct {
	ListParams
	Origin            string `json:"origin,omitempty"`
	SourceWorkOrderID string `json:"source_work_order_id,omitempty"`
	// Statuses is sent comma-joined (statuses=available,reserved).
	Statuses        []string `json:"statuses,omitempty"`
	IncludeExpired  bool     `json:"include_expired,omitempty"`
	IncludeReceived bool     `json:"include_received,omitempty"`
}

// StockUnitGetParams are optional flags for Get / PDF.
type StockUnitGetParams struct {
	// IncludeReceived lets Get/PDF see a unit still in Receiving (else 404).
	IncludeReceived bool `json:"include_received,omitempty"`
}

// StockUnitCreateRequest is POST /partner/v1/stock-units.
type StockUnitCreateRequest struct {
	ProductID               string          `json:"product_id"`
	SerialNumber            string          `json:"serial_number,omitempty"`
	LotNumber               string          `json:"lot_number,omitempty"`
	SourceWorkOrderID       string          `json:"source_work_order_id,omitempty"`
	ConfigurationID         string          `json:"configuration_id,omitempty"`
	ConfigurationRevisionID string          `json:"configuration_revision_id,omitempty"`
	LocationID              string          `json:"location_id,omitempty"`
	Status                  string          `json:"status,omitempty"`
	Origin                  string          `json:"origin,omitempty"`
	UnitCost                string          `json:"unit_cost,omitempty"`
	AskingPrice             string          `json:"asking_price,omitempty"`
	Specs                   json.RawMessage `json:"specs,omitempty"`
	ConditionGrade          string          `json:"condition_grade,omitempty"`
	ReturnedCustomerID      string          `json:"returned_customer_id,omitempty"`
	ReturnedOrderID         string          `json:"returned_order_id,omitempty"`
	ReturnedRTOAgreementID  string          `json:"returned_rto_agreement_id,omitempty"`
}

// StockUnitRegisterReturnRequest is POST /partner/v1/stock-units/register-return
// (a rental return coming back to lot, held in "received" until released).
type StockUnitRegisterReturnRequest struct {
	ProductID              string `json:"product_id"`
	ConfigCategory         string `json:"config_category,omitempty"`
	LocationID             string `json:"location_id,omitempty"`
	SerialNumber           string `json:"serial_number,omitempty"`
	ConditionGrade         string `json:"condition_grade,omitempty"`
	UnitCost               string `json:"unit_cost,omitempty"`
	AskingPrice            string `json:"asking_price,omitempty"`
	ReturnedCustomerID     string `json:"returned_customer_id,omitempty"`
	ReturnedOrderID        string `json:"returned_order_id,omitempty"`
	ReturnedRTOAgreementID string `json:"returned_rto_agreement_id,omitempty"`
}

// StockUnitUpdateRequest is PATCH /partner/v1/stock-units/{id}. Only non-nil
// fields are sent.
type StockUnitUpdateRequest struct {
	SerialNumber           *string         `json:"serial_number,omitempty"`
	LotNumber              *string         `json:"lot_number,omitempty"`
	LocationID             *string         `json:"location_id,omitempty"`
	Status                 *string         `json:"status,omitempty"`
	UnitCost               *string         `json:"unit_cost,omitempty"`
	AskingPrice            *string         `json:"asking_price,omitempty"`
	Specs                  json.RawMessage `json:"specs,omitempty"`
	PartnerSoldVisibleDays *int            `json:"partner_sold_visible_days,omitempty"`
	PartnerKeepVisible     *bool           `json:"partner_keep_visible,omitempty"`
	HideFromEmbed          *bool           `json:"hide_from_embed,omitempty"`
	DiscountReason         *string         `json:"discount_reason,omitempty"`
	ConditionGrade         *string         `json:"condition_grade,omitempty"`
	ReturnedCustomerID     *string         `json:"returned_customer_id,omitempty"`
	ReturnedOrderID        *string         `json:"returned_order_id,omitempty"`
	ReturnedRTOAgreementID *string         `json:"returned_rto_agreement_id,omitempty"`
}

// StockUnitsService calls /partner/v1/stock-units (ERP replacement of
// /partner/v1/lot-stock).
type StockUnitsService struct{ c *Client }

// List returns lot stock. Scope partner-api.lot-stock.read.
func (s *StockUnitsService) List(ctx context.Context, params StockUnitListParams) (*Page[StockUnit], error) {
	return listPage[StockUnit](ctx, s.c, "/partner/v1/stock-units", "stock_units", params)
}

// Get returns one unit. Scope partner-api.lot-stock.read.
func (s *StockUnitsService) Get(ctx context.Context, id string) (*StockUnit, error) {
	return s.GetWithParams(ctx, id, StockUnitGetParams{})
}

// GetWithParams is Get with include_received.
func (s *StockUnitsService) GetWithParams(ctx context.Context, id string, params StockUnitGetParams) (*StockUnit, error) {
	return call[StockUnit](ctx, s.c, http.MethodGet, "/partner/v1/stock-units/"+seg(id), params, nil)
}

// Create adds a unit. Scope partner-api.lot-stock.write.
func (s *StockUnitsService) Create(ctx context.Context, body StockUnitCreateRequest, opts ...RequestOption) (*StockUnit, error) {
	return call[StockUnit](ctx, s.c, http.MethodPost, "/partner/v1/stock-units", nil, body, opts...)
}

// RegisterReturn registers a rental return as a "received" unit. Scope
// partner-api.lot-stock.write.
func (s *StockUnitsService) RegisterReturn(ctx context.Context, body StockUnitRegisterReturnRequest, opts ...RequestOption) (*StockUnit, error) {
	return call[StockUnit](ctx, s.c, http.MethodPost, "/partner/v1/stock-units/register-return", nil, body, opts...)
}

// Update patches a unit. Scope partner-api.lot-stock.write.
func (s *StockUnitsService) Update(ctx context.Context, id string, body StockUnitUpdateRequest) (*StockUnit, error) {
	return call[StockUnit](ctx, s.c, http.MethodPatch, "/partner/v1/stock-units/"+seg(id), nil, body)
}

// Release releases a returned unit to sellable stock. The ERP may answer 409
// until staff upload the work-order preview and floorplan files. Scope
// partner-api.lot-stock.write.
func (s *StockUnitsService) Release(ctx context.Context, id string) (*StockUnit, error) {
	return call[StockUnit](ctx, s.c, http.MethodPost, "/partner/v1/stock-units/"+seg(id)+"/release", nil, nil)
}

// PDF downloads the unit info sheet as PDF bytes. Scope
// partner-api.lot-stock.read.
func (s *StockUnitsService) PDF(ctx context.Context, id string, params StockUnitGetParams) ([]byte, error) {
	return s.c.http.do(ctx, http.MethodGet, "/partner/v1/stock-units/"+seg(id)+"/pdf", params, nil, withAccept("application/pdf"))
}
