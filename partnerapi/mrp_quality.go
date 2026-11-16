package partnerapi

import (
	"context"
	"net/http"
	"time"
)

// MRPPlannedOrder is a suggested buy or make quantity from an MRP run.
// Qty is a decimal string. Action is "buy" or "make".
type MRPPlannedOrder struct {
	ID                     string    `json:"id"`
	TenantID               string    `json:"tenant_id"`
	RunID                  string    `json:"run_id"`
	ProductID              string    `json:"product_id"`
	Action                 string    `json:"action"`
	Qty                    string    `json:"qty"`
	ConvertedRequisitionID *string   `json:"converted_requisition_id,omitempty"`
	ConvertedWorkOrderID   *string   `json:"converted_work_order_id,omitempty"`
	CreatedAt              time.Time `json:"created_at"`
}

// MRPRun is one material-requirements planning execution
// (/partner/v1/mrp/runs). Status is running, completed or failed.
// PlannedOrders is filled by Get and by Run; the list may omit it.
type MRPRun struct {
	ID            string            `json:"id"`
	TenantID      string            `json:"tenant_id"`
	Status        string            `json:"status"`
	Note          string            `json:"note,omitempty"`
	CreatedBy     *string           `json:"created_by,omitempty"`
	CreatedAt     time.Time         `json:"created_at"`
	PlannedOrders []MRPPlannedOrder `json:"planned_orders,omitempty"`
}

// MRPRunListParams filters GET /partner/v1/mrp/runs (paging only).
type MRPRunListParams struct {
	ListParams
}

// MRPRunRequest is POST /partner/v1/mrp/runs.
type MRPRunRequest struct {
	Note string `json:"note,omitempty"`
}

// MRPService calls /partner/v1/mrp.
type MRPService struct{ c *Client }

// ListRuns returns MRP runs. Scope partner-api.mrp.read.
func (s *MRPService) ListRuns(ctx context.Context, params MRPRunListParams) (*Page[MRPRun], error) {
	return listPage[MRPRun](ctx, s.c, "/partner/v1/mrp/runs", "runs", params)
}

// GetRun returns one run with its planned orders. Scope partner-api.mrp.read.
func (s *MRPService) GetRun(ctx context.Context, runID string) (*MRPRun, error) {
	return call[MRPRun](ctx, s.c, http.MethodGet, "/partner/v1/mrp/runs/"+seg(runID), nil, nil)
}

// Run executes a new MRP run (POST, 201). Scope partner-api.mrp.write.
func (s *MRPService) Run(ctx context.Context, body MRPRunRequest, opts ...RequestOption) (*MRPRun, error) {
	return call[MRPRun](ctx, s.c, http.MethodPost, "/partner/v1/mrp/runs", nil, body, opts...)
}

// QualityNCR is a non-conformance report (/partner/v1/quality/ncrs).
// SourceType is work_order, service_ticket or goods_receipt; Severity is
// minor, major or critical; Disposition is open, rework, scrap, use_as_is or
// closed; Status is draft, open or closed.
type QualityNCR struct {
	ID          string    `json:"id"`
	TenantID    string    `json:"tenant_id"`
	Number      string    `json:"number,omitempty"`
	SourceType  string    `json:"source_type"`
	SourceID    string    `json:"source_id"`
	Severity    string    `json:"severity"`
	Disposition string    `json:"disposition"`
	CAPA        string    `json:"capa,omitempty"`
	Status      string    `json:"status"`
	CreatedBy   *string   `json:"created_by,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// QualityNCRListParams filters GET /partner/v1/quality/ncrs.
type QualityNCRListParams struct {
	ListParams
	Status string `json:"status,omitempty"`
}

// QualityNCRCreateRequest is POST /partner/v1/quality/ncrs.
type QualityNCRCreateRequest struct {
	SourceType  string `json:"source_type,omitempty"`
	SourceID    string `json:"source_id,omitempty"`
	Severity    string `json:"severity,omitempty"`
	Disposition string `json:"disposition,omitempty"`
	CAPA        string `json:"capa,omitempty"`
}

// QualityService calls /partner/v1/quality.
type QualityService struct{ c *Client }

// ListNCRs returns NCRs (filter by Status). Scope partner-api.quality.read.
func (s *QualityService) ListNCRs(ctx context.Context, params QualityNCRListParams) (*Page[QualityNCR], error) {
	return listPage[QualityNCR](ctx, s.c, "/partner/v1/quality/ncrs", "ncrs", params)
}

// GetNCR returns one NCR. Scope partner-api.quality.read.
func (s *QualityService) GetNCR(ctx context.Context, ncrID string) (*QualityNCR, error) {
	return call[QualityNCR](ctx, s.c, http.MethodGet, "/partner/v1/quality/ncrs/"+seg(ncrID), nil, nil)
}

// CreateNCR opens an NCR (POST, 201). Scope partner-api.quality.write. Pass
// WithIdempotencyKey to make retries safe.
func (s *QualityService) CreateNCR(ctx context.Context, body QualityNCRCreateRequest, opts ...RequestOption) (*QualityNCR, error) {
	return call[QualityNCR](ctx, s.c, http.MethodPost, "/partner/v1/quality/ncrs", nil, body, opts...)
}
