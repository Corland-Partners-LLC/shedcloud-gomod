package partnerapi

import (
	"context"
	"time"
)

// CommissionPlanRule is one rule of a commission plan. RatePct is a decimal
// string.
type CommissionPlanRule struct {
	ID              string    `json:"id"`
	TenantID        string    `json:"tenant_id,omitempty"`
	PlanID          string    `json:"plan_id"`
	StoreLocationID *string   `json:"store_location_id,omitempty"`
	RoleSlug        string    `json:"role_slug,omitempty"`
	ProductID       *string   `json:"product_id,omitempty"`
	AppliesTo       string    `json:"applies_to"`
	RatePct         string    `json:"rate_pct"`
	CreatedAt       time.Time `json:"created_at"`
}

// CommissionPlan is a commission plan (/partner/v1/commissions/plans).
type CommissionPlan struct {
	ID        string               `json:"id"`
	TenantID  string               `json:"tenant_id,omitempty"`
	Name      string               `json:"name"`
	IsActive  bool                 `json:"is_active"`
	CreatedAt time.Time            `json:"created_at"`
	UpdatedAt time.Time            `json:"updated_at"`
	Rules     []CommissionPlanRule `json:"rules,omitempty"`
}

// CommissionAccrual is a commission earned on an order
// (/partner/v1/commissions/accruals). Amount is a decimal string.
type CommissionAccrual struct {
	ID            string    `json:"id"`
	TenantID      string    `json:"tenant_id,omitempty"`
	PlanID        *string   `json:"plan_id,omitempty"`
	OrderID       string    `json:"order_id"`
	SalespersonID string    `json:"salesperson_id"`
	Status        string    `json:"status"`
	Amount        string    `json:"amount"`
	SalesReturnID *string   `json:"sales_return_id,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// CommissionListParams are the paging params of the commission lists.
type CommissionListParams struct {
	ListParams
}

// CommissionsService calls /partner/v1/commissions.
type CommissionsService struct{ c *Client }

// Plans lists commission plans. Scope partner-api.commissions.read.
func (s *CommissionsService) Plans(ctx context.Context, params CommissionListParams) (*Page[CommissionPlan], error) {
	return listPage[CommissionPlan](ctx, s.c, "/partner/v1/commissions/plans", "plans", params)
}

// Accruals lists commission accruals. Scope partner-api.commissions.read.
func (s *CommissionsService) Accruals(ctx context.Context, params CommissionListParams) (*Page[CommissionAccrual], error) {
	return listPage[CommissionAccrual](ctx, s.c, "/partner/v1/commissions/accruals", "accruals", params)
}
