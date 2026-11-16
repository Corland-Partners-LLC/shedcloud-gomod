package partnerapi

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

// AnalyticsExportJob is an async dataset dump
// (/partner/v1/analytics/exports). Dataset is orders, sales_ledger, journals
// or inventory; Status is queued, running, done or failed.
type AnalyticsExportJob struct {
	ID         string     `json:"id"`
	TenantID   string     `json:"tenant_id"`
	Dataset    string     `json:"dataset"`
	Status     string     `json:"status"`
	StorageKey *string    `json:"storage_key,omitempty"`
	Error      string     `json:"error,omitempty"`
	CreatedBy  *string    `json:"created_by,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	FinishedAt *time.Time `json:"finished_at,omitempty"`
}

// AnalyticsExportListParams filters GET /partner/v1/analytics/exports.
type AnalyticsExportListParams struct {
	ListParams
	Status string `json:"status,omitempty"`
}

// AnalyticsService calls /partner/v1/analytics (read only).
type AnalyticsService struct{ c *Client }

// ListExports returns export jobs (filter by Status). Scope
// partner-api.analytics.read.
func (s *AnalyticsService) ListExports(ctx context.Context, params AnalyticsExportListParams) (*Page[AnalyticsExportJob], error) {
	return listPage[AnalyticsExportJob](ctx, s.c, "/partner/v1/analytics/exports", "jobs", params)
}

// GetExport returns one export job. Scope partner-api.analytics.read.
func (s *AnalyticsService) GetExport(ctx context.Context, jobID string) (*AnalyticsExportJob, error) {
	return call[AnalyticsExportJob](ctx, s.c, http.MethodGet, "/partner/v1/analytics/exports/"+seg(jobID), nil, nil)
}

// EDIMap is a tenant EDI transaction-set mapping (/partner/v1/edi-maps).
// TransactionSet is 850, 855, 856 or 810. Mapping is an opaque JSON document.
type EDIMap struct {
	ID             string          `json:"id"`
	TenantID       string          `json:"tenant_id"`
	TransactionSet string          `json:"transaction_set"`
	Name           string          `json:"name"`
	Mapping        json.RawMessage `json:"mapping"`
	IsActive       bool            `json:"is_active"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

// EDIMapListParams filters GET /partner/v1/edi-maps.
type EDIMapListParams struct {
	ListParams
}

// EDIMapCreateRequest is POST /partner/v1/edi-maps. Mapping is raw JSON.
type EDIMapCreateRequest struct {
	TransactionSet string          `json:"transaction_set,omitempty"`
	Name           string          `json:"name,omitempty"`
	Mapping        json.RawMessage `json:"mapping,omitempty"`
}

// EDIMapsService calls /partner/v1/edi-maps.
type EDIMapsService struct{ c *Client }

// List returns EDI maps. Scope partner-api.edi-maps.read.
func (s *EDIMapsService) List(ctx context.Context, params EDIMapListParams) (*Page[EDIMap], error) {
	return listPage[EDIMap](ctx, s.c, "/partner/v1/edi-maps", "edi_maps", params)
}

// Get returns one EDI map. Scope partner-api.edi-maps.read.
func (s *EDIMapsService) Get(ctx context.Context, mapID string) (*EDIMap, error) {
	return call[EDIMap](ctx, s.c, http.MethodGet, "/partner/v1/edi-maps/"+seg(mapID), nil, nil)
}

// Create makes an EDI map (POST, 201). Scope partner-api.edi-maps.write.
// Pass WithIdempotencyKey to make retries safe.
func (s *EDIMapsService) Create(ctx context.Context, body EDIMapCreateRequest, opts ...RequestOption) (*EDIMap, error) {
	return call[EDIMap](ctx, s.c, http.MethodPost, "/partner/v1/edi-maps", nil, body, opts...)
}
