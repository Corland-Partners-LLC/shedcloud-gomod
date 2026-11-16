package partnerapi

import (
	"context"
	"net/http"
	"time"
)

// FXRate is a currency exchange rate (/partner/v1/fx/rates). Rate is a decimal
// string; AsOf is YYYY-MM-DD.
type FXRate struct {
	ID           string    `json:"id"`
	TenantID     string    `json:"tenant_id,omitempty"`
	FromCurrency string    `json:"from_currency"`
	ToCurrency   string    `json:"to_currency"`
	Rate         string    `json:"rate"`
	AsOf         string    `json:"as_of"`
	CreatedAt    time.Time `json:"created_at"`
}

// FXRateRequest is the body of POST /partner/v1/fx/rates (upsert). Rate is a
// decimal string; AsOf is YYYY-MM-DD.
type FXRateRequest struct {
	FromCurrency string `json:"from_currency,omitempty"`
	ToCurrency   string `json:"to_currency,omitempty"`
	Rate         string `json:"rate,omitempty"`
	AsOf         string `json:"as_of,omitempty"`
}

// FXRateListParams are the paging params of GET /partner/v1/fx/rates.
type FXRateListParams struct {
	ListParams
}

// FXConvertParams are the query params of GET /partner/v1/fx/convert. Amount is
// a decimal string; AsOf (YYYY-MM-DD) is optional.
type FXConvertParams struct {
	Amount string `json:"amount,omitempty"`
	From   string `json:"from,omitempty"`
	To     string `json:"to,omitempty"`
	AsOf   string `json:"as_of,omitempty"`
}

// FXConversion is the result of GET /partner/v1/fx/convert: the converted
// amount (decimal string) and the rate row that was applied.
type FXConversion struct {
	Amount string `json:"amount"`
	Rate   FXRate `json:"rate"`
}

// FXService calls /partner/v1/fx.
type FXService struct{ c *Client }

// Rates lists stored exchange rates. Scope partner-api.fx.read.
func (s *FXService) Rates(ctx context.Context, params FXRateListParams) (*Page[FXRate], error) {
	return listPage[FXRate](ctx, s.c, "/partner/v1/fx/rates", "currency_rates", params)
}

// UpsertRate creates or updates a rate. Scope partner-api.fx.write.
func (s *FXService) UpsertRate(ctx context.Context, body FXRateRequest, opts ...RequestOption) (*FXRate, error) {
	return call[FXRate](ctx, s.c, http.MethodPost, "/partner/v1/fx/rates", nil, body, opts...)
}

// Convert converts an amount between currencies. Scope partner-api.fx.read.
func (s *FXService) Convert(ctx context.Context, params FXConvertParams) (*FXConversion, error) {
	return call[FXConversion](ctx, s.c, http.MethodGet, "/partner/v1/fx/convert", params, nil)
}

// ConsolidationRun is one multi-entity consolidation run
// (/partner/v1/consolidations). AsOf is YYYY-MM-DD.
type ConsolidationRun struct {
	ID             string    `json:"id"`
	TenantID       string    `json:"tenant_id,omitempty"`
	AsOf           string    `json:"as_of"`
	Status         string    `json:"status"`
	JournalEntryID *string   `json:"journal_entry_id,omitempty"`
	CreatedBy      *string   `json:"created_by,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

// ConsolidationListParams filters GET /partner/v1/consolidations.
type ConsolidationListParams struct {
	ListParams
	Status string `json:"status,omitempty"`
}

// ConsolidationRequest is the body of POST /partner/v1/consolidations. AsOf is
// YYYY-MM-DD.
type ConsolidationRequest struct {
	AsOf string `json:"as_of,omitempty"`
}

// ConsolidationsService calls /partner/v1/consolidations.
type ConsolidationsService struct{ c *Client }

// List returns consolidation runs (filter by Status). Scope
// partner-api.consolidations.read.
func (s *ConsolidationsService) List(ctx context.Context, params ConsolidationListParams) (*Page[ConsolidationRun], error) {
	return listPage[ConsolidationRun](ctx, s.c, "/partner/v1/consolidations", "consolidations", params)
}

// Get returns one consolidation run. Scope partner-api.consolidations.read.
func (s *ConsolidationsService) Get(ctx context.Context, runID string) (*ConsolidationRun, error) {
	return call[ConsolidationRun](ctx, s.c, http.MethodGet, "/partner/v1/consolidations/"+seg(runID), nil, nil)
}

// Create starts a consolidation run as of a date. Scope
// partner-api.consolidations.write.
func (s *ConsolidationsService) Create(ctx context.Context, body ConsolidationRequest, opts ...RequestOption) (*ConsolidationRun, error) {
	return call[ConsolidationRun](ctx, s.c, http.MethodPost, "/partner/v1/consolidations", nil, body, opts...)
}
