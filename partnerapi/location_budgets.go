package partnerapi

import (
	"context"
	"net/http"
	"time"
)

// LocationBudgetMonth is one month of a location's yearly sales budget.
// Money is a decimal string.
type LocationBudgetMonth struct {
	Month        int    `json:"month"`
	WeightPct    string `json:"weight_pct"`
	Amount       string `json:"amount"`
	WeeklyAmount string `json:"weekly_amount"`
	LeadTarget   int    `json:"lead_target"`
}

// LocationBudget is a location's yearly sales budget
// (/partner/v1/location-budgets).
type LocationBudget struct {
	ID              string                `json:"id"`
	TenantID        string                `json:"tenant_id,omitempty"`
	LocationID      string                `json:"location_id"`
	LocationName    string                `json:"location_name,omitempty"`
	LocationCode    string                `json:"location_code,omitempty"`
	Year            int                   `json:"year"`
	AnnualBudget    string                `json:"annual_budget"`
	BonusThreshold1 string                `json:"bonus_threshold_1"`
	BonusThreshold2 string                `json:"bonus_threshold_2"`
	InventoryCap    string                `json:"inventory_cap"`
	CloseRateTarget string                `json:"close_rate_target"`
	Notes           string                `json:"notes,omitempty"`
	Months          []LocationBudgetMonth `json:"months"`
	CreatedBy       string                `json:"created_by,omitempty"`
	CreatedAt       time.Time             `json:"created_at"`
	UpdatedAt       time.Time             `json:"updated_at"`
}

// LocationBudgetListParams filters GET /partner/v1/location-budgets.
type LocationBudgetListParams struct {
	ListParams
	Year       int    `json:"year,omitempty"`
	LocationID string `json:"location_id,omitempty"`
}

// LocationBudgetRequest is the body of POST (create) and PUT (replace)
// /partner/v1/location-budgets. Money is a decimal string.
type LocationBudgetRequest struct {
	LocationID      string                `json:"location_id,omitempty"`
	Year            int                   `json:"year,omitempty"`
	AnnualBudget    string                `json:"annual_budget,omitempty"`
	BonusThreshold1 string                `json:"bonus_threshold_1,omitempty"`
	BonusThreshold2 string                `json:"bonus_threshold_2,omitempty"`
	InventoryCap    string                `json:"inventory_cap,omitempty"`
	CloseRateTarget string                `json:"close_rate_target,omitempty"`
	Notes           string                `json:"notes,omitempty"`
	Months          []LocationBudgetMonth `json:"months,omitempty"`
}

// LocationBudgetsService calls /partner/v1/location-budgets.
type LocationBudgetsService struct{ c *Client }

// List returns budgets (filter by Year / LocationID). Scope
// partner-api.location-budgets.read.
func (s *LocationBudgetsService) List(ctx context.Context, params LocationBudgetListParams) (*Page[LocationBudget], error) {
	return listPage[LocationBudget](ctx, s.c, "/partner/v1/location-budgets", "location_budgets", params)
}

// Get returns one budget by id. Scope partner-api.location-budgets.read.
func (s *LocationBudgetsService) Get(ctx context.Context, id string) (*LocationBudget, error) {
	return call[LocationBudget](ctx, s.c, http.MethodGet, "/partner/v1/location-budgets/"+seg(id), nil, nil)
}

// Create makes a budget. Scope partner-api.location-budgets.write.
func (s *LocationBudgetsService) Create(ctx context.Context, body LocationBudgetRequest, opts ...RequestOption) (*LocationBudget, error) {
	return call[LocationBudget](ctx, s.c, http.MethodPost, "/partner/v1/location-budgets", nil, body, opts...)
}

// Update replaces a budget (PUT). Scope partner-api.location-budgets.write.
func (s *LocationBudgetsService) Update(ctx context.Context, id string, body LocationBudgetRequest) (*LocationBudget, error) {
	return call[LocationBudget](ctx, s.c, http.MethodPut, "/partner/v1/location-budgets/"+seg(id), nil, body)
}
