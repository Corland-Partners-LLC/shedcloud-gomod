package partnerapi

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

// EngineeringChange is a product engineering change order with BOM/options
// snapshots (/partner/v1/engineering-changes). Status is draft, pending,
// approved or rejected. The snapshots are opaque JSON documents.
type EngineeringChange struct {
	ID              string          `json:"id"`
	TenantID        string          `json:"tenant_id"`
	Number          string          `json:"number,omitempty"`
	ProductID       string          `json:"product_id"`
	Status          string          `json:"status"`
	BOMSnapshot     json.RawMessage `json:"bom_snapshot"`
	OptionsSnapshot json.RawMessage `json:"options_snapshot"`
	ApprovedBy      *string         `json:"approved_by,omitempty"`
	ApprovedAt      *time.Time      `json:"approved_at,omitempty"`
	CreatedBy       *string         `json:"created_by,omitempty"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

// EngineeringChangeListParams filters GET /partner/v1/engineering-changes.
type EngineeringChangeListParams struct {
	ListParams
	Status string `json:"status,omitempty"`
}

// EngineeringChangeCreateRequest is POST /partner/v1/engineering-changes.
// The snapshots are raw JSON.
type EngineeringChangeCreateRequest struct {
	ProductID       string          `json:"product_id,omitempty"`
	BOMSnapshot     json.RawMessage `json:"bom_snapshot,omitempty"`
	OptionsSnapshot json.RawMessage `json:"options_snapshot,omitempty"`
}

// EngineeringChangesService calls /partner/v1/engineering-changes.
type EngineeringChangesService struct{ c *Client }

// List returns change orders (filter by Status). Scope
// partner-api.engineering-changes.read.
func (s *EngineeringChangesService) List(ctx context.Context, params EngineeringChangeListParams) (*Page[EngineeringChange], error) {
	return listPage[EngineeringChange](ctx, s.c, "/partner/v1/engineering-changes", "engineering_changes", params)
}

// Get returns one change order. Scope partner-api.engineering-changes.read.
func (s *EngineeringChangesService) Get(ctx context.Context, changeID string) (*EngineeringChange, error) {
	return call[EngineeringChange](ctx, s.c, http.MethodGet, "/partner/v1/engineering-changes/"+seg(changeID), nil, nil)
}

// Create opens a change order (POST, 201). Scope
// partner-api.engineering-changes.write. Pass WithIdempotencyKey to make
// retries safe.
func (s *EngineeringChangesService) Create(ctx context.Context, body EngineeringChangeCreateRequest, opts ...RequestOption) (*EngineeringChange, error) {
	return call[EngineeringChange](ctx, s.c, http.MethodPost, "/partner/v1/engineering-changes", nil, body, opts...)
}
