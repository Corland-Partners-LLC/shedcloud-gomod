package partnerapi

import (
	"context"
	"net/http"
	"time"
)

// ---- Sales ledger ------------------------------------------------------------

// SalesLedgerRequest is the body of POST /partner/v1/sales-ledger (the ERP
// report is a POST with filters; it was a GET in 0.x). Dates are YYYY-MM-DD or
// RFC3339.
//
// LedgerFrom/LedgerTo filter the ledger (status/event) date; SaleFrom/SaleTo
// filter the calendar sale date. MovementMode "non_zero" keeps only movements
// whose delta is not zero.
type SalesLedgerRequest struct {
	LedgerFrom   string `json:"ledger_from,omitempty"`
	LedgerTo     string `json:"ledger_to,omitempty"`
	SaleFrom     string `json:"sale_from,omitempty"`
	SaleTo       string `json:"sale_to,omitempty"`
	StoreID      string `json:"store_id,omitempty"`
	MovementMode string `json:"movement_mode,omitempty"`
	Q            string `json:"q,omitempty"`
	Page         int    `json:"page,omitempty"`
	PageSize     int    `json:"page_size,omitempty"`
	// IncludeSums adds the SumDeltaAmount / SumCancelledDeltaAmount totals
	// (over the whole filter, not the page).
	IncludeSums bool `json:"include_sums,omitempty"`
}

// SalesLedgerRow is one ledger movement. Money is a decimal string; Delta*
// are the change this movement made, New* the values after it. Cancelled is
// sticky after a cancel/delete of a previously sold order.
type SalesLedgerRow struct {
	ID                         string     `json:"id"`
	TenantID                   string     `json:"tenant_id,omitempty"`
	SummaryID                  string     `json:"summary_id"`
	OrderID                    string     `json:"order_id"`
	OrderNumber                string     `json:"order_number,omitempty"`
	CustomerID                 string     `json:"customer_id,omitempty"`
	CustomerName               string     `json:"customer_name,omitempty"`
	CustomerPhone              string     `json:"customer_phone,omitempty"`
	StoreLocationID            string     `json:"store_location_id,omitempty"`
	StoreLocationName          string     `json:"store_location_name,omitempty"`
	SalesPersonName            string     `json:"sales_person_name,omitempty"`
	SerialNumber               string     `json:"serial_number,omitempty"`
	SourceSlug                 string     `json:"source_slug,omitempty"`
	DealType                   string     `json:"deal_type,omitempty"`
	SaleType                   string     `json:"sale_type,omitempty"`
	MovementType               string     `json:"movement_type"`
	SourceType                 string     `json:"source_type,omitempty"`
	Reason                     string     `json:"reason,omitempty"`
	PreviousStatus             string     `json:"previous_status,omitempty"`
	NewStatus                  string     `json:"new_status,omitempty"`
	CurrentStatus              string     `json:"current_status,omitempty"`
	DeltaBasePrice             string     `json:"delta_base_price"`
	DeltaUpgradesAmount        string     `json:"delta_upgrades_amount"`
	DeltaRemovedUpgradesAmount string     `json:"delta_removed_upgrades_amount"`
	DeltaAmount                string     `json:"delta_amount"`
	NewBasePrice               string     `json:"new_base_price"`
	NewUpgradesAmount          string     `json:"new_upgrades_amount"`
	NewRemovedUpgradesAmount   string     `json:"new_removed_upgrades_amount"`
	NewOrderAmount             string     `json:"new_order_amount"`
	NewLedgerAmount            string     `json:"new_ledger_amount"`
	Cancelled                  bool       `json:"cancelled"`
	SaleDate                   *time.Time `json:"sale_date,omitempty"`
	LedgerDate                 time.Time  `json:"ledger_date"`
	CreatedAt                  time.Time  `json:"created_at"`
	TotalDeltaAmount           string     `json:"total_delta_amount,omitempty"`
	ProductName                string     `json:"product_name,omitempty"`
	// ProductSize is the main line size in feet, for example "10x12".
	ProductSize string `json:"product_size,omitempty"`
}

// ReportPagination is the pagination block of the ledger report (note
// total_items instead of total).
type ReportPagination struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalItems int64 `json:"total_items"`
}

// SalesLedgerReport is the POST /partner/v1/sales-ledger reply.
type SalesLedgerReport struct {
	Rows       []SalesLedgerRow `json:"rows"`
	Pagination ReportPagination `json:"pagination"`
	// SumDeltaAmount / SumCancelledDeltaAmount are set with IncludeSums.
	SumDeltaAmount          string `json:"sum_delta_amount,omitempty"`
	SumCancelledDeltaAmount string `json:"sum_cancelled_delta_amount,omitempty"`
}

// SalesLedgerService calls POST /partner/v1/sales-ledger. Scope
// partner-api.sales-ledger.read.
type SalesLedgerService struct{ c *Client }

// Report returns ledger movements for the authenticated company.
func (s *SalesLedgerService) Report(ctx context.Context, req SalesLedgerRequest) (*SalesLedgerReport, error) {
	return call[SalesLedgerReport](ctx, s.c, http.MethodPost, "/partner/v1/sales-ledger", nil, req)
}

// ---- Revenue forecast --------------------------------------------------------

// RevenueForecastParams are the query params of GET /partner/v1/revenue-forecast.
// The ERP reads them in camelCase (the only camelCase query params of the
// API). Week bounds are YYYY-MM-DD. Slices are sent comma-joined.
type RevenueForecastParams struct {
	Year                int      `json:"year,omitempty"`
	LocationIDs         []string `json:"locationIds,omitempty"`
	ExcludedLocationIDs []string `json:"excludedLocationIds,omitempty"`
	WeekStart           string   `json:"weekStart,omitempty"`
	WeekEnd             string   `json:"weekEnd,omitempty"`
	Lang                string   `json:"lang,omitempty"`
	Page                int      `json:"page,omitempty"`
	Limit               int      `json:"limit,omitempty"`
}

// RevenueForecastItem is one tile of the forecast strip. Amount is a JSON
// number here (the ERP sends this report as a float, unlike other money).
type RevenueForecastItem struct {
	Label  string  `json:"label"`
	Amount float64 `json:"amount"`
}

// RevenueForecastResponse is the reply of GET /partner/v1/revenue-forecast.
type RevenueForecastResponse struct {
	Forecast []RevenueForecastItem `json:"forecast"`
}

// RevenueForecastService calls GET /partner/v1/revenue-forecast. Scope
// partner-api.revenue-forecast.read.
type RevenueForecastService struct{ c *Client }

// Get returns the forecast strip (label + amount rows). v1.0.0 note: this is
// the ERP's KPI forecast strip; the 0.x KPI blocks, week buckets and detail
// rows do not exist in the ERP.
func (s *RevenueForecastService) Get(ctx context.Context, params RevenueForecastParams) (*RevenueForecastResponse, error) {
	return call[RevenueForecastResponse](ctx, s.c, http.MethodGet, "/partner/v1/revenue-forecast", params, nil)
}
