package partnerapi

import (
	"context"
	"net/http"
	"time"
)

// InventorySettings is the per-tenant lot-stock sold-visibility setting
// (/partner/v1/inventory-settings).
type InventorySettings struct {
	TenantID string `json:"tenant_id"`
	// PartnerSoldVisibleDays is how long a sold unit stays listed (default 30).
	PartnerSoldVisibleDays int `json:"partner_sold_visible_days"`
	// PartnerSoldVisibleMode: hide_after_days | always_show_sold.
	PartnerSoldVisibleMode string `json:"partner_sold_visible_mode"`
	// InventoryLimitMode says how location inventory limits apply.
	InventoryLimitMode string    `json:"inventory_limit_mode"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// InventorySettingsUpdateRequest is PUT /partner/v1/inventory-settings. Only
// non-nil fields are sent.
type InventorySettingsUpdateRequest struct {
	PartnerSoldVisibleDays *int    `json:"partner_sold_visible_days,omitempty"`
	PartnerSoldVisibleMode *string `json:"partner_sold_visible_mode,omitempty"`
}

// InventoryAdjustment is a batch manual stock correction at one location.
type InventoryAdjustment struct {
	ID           string    `json:"id"`
	TenantID     string    `json:"tenant_id"`
	AdjustmentNo string    `json:"adjustment_no"`
	LocationID   string    `json:"location_id"`
	Notes        string    `json:"notes,omitempty"`
	CreatedBy    string    `json:"created_by,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	LocationCode string    `json:"location_code,omitempty"`
	LocationName string    `json:"location_name,omitempty"`
	ActorLabel   string    `json:"actor_label,omitempty"`
	LineCount    int       `json:"line_count,omitempty"`
}

// InventoryAdjustmentLine is one movement row of an adjustment. Quantities and
// costs are decimal strings; Quantity is signed.
type InventoryAdjustmentLine struct {
	ID           string    `json:"id"`
	ProductID    string    `json:"product_id"`
	ProductSKU   string    `json:"product_sku,omitempty"`
	ProductName  string    `json:"product_name,omitempty"`
	Quantity     string    `json:"quantity"`
	QtyBefore    string    `json:"qty_before"`
	QtyAfter     string    `json:"qty_after"`
	UnitCost     string    `json:"unit_cost"`
	Notes        string    `json:"notes,omitempty"`
	OccurredAt   time.Time `json:"occurred_at"`
	MovementType string    `json:"movement_type"`
}

// InventoryAdjustmentDetail is an adjustment with its lines (create and get
// replies are {"adjustment": {...}, "lines": [...]}).
type InventoryAdjustmentDetail struct {
	Adjustment InventoryAdjustment       `json:"adjustment"`
	Lines      []InventoryAdjustmentLine `json:"lines"`
}

// InventoryAdjustmentReportRow is one denormalized line of the adjustments
// report.
type InventoryAdjustmentReportRow struct {
	AdjustmentID  string    `json:"adjustment_id"`
	AdjustmentNo  string    `json:"adjustment_no"`
	MovementID    string    `json:"movement_id"`
	ProductID     string    `json:"product_id"`
	ProductSKU    string    `json:"product_sku"`
	ProductName   string    `json:"product_name"`
	LocationID    string    `json:"location_id"`
	LocationCode  string    `json:"location_code"`
	LocationName  string    `json:"location_name"`
	ActorLabel    string    `json:"actor_label"`
	QtyBefore     string    `json:"qty_before"`
	Quantity      string    `json:"quantity"`
	QtyAfter      string    `json:"qty_after"`
	UnitCost      string    `json:"unit_cost"`
	ExtendedCost  string    `json:"extended_cost"`
	OccurredAt    time.Time `json:"occurred_at"`
	LineNotes     string    `json:"line_notes,omitempty"`
	DocumentNotes string    `json:"document_notes,omitempty"`
}

// InventoryAdjustmentListParams filters GET /partner/v1/inventory/adjustments
// (default page size 25). From/To are RFC3339 or YYYY-MM-DD.
type InventoryAdjustmentListParams struct {
	ListParams
	LocationID string `json:"location_id,omitempty"`
	From       string `json:"from,omitempty"`
	To         string `json:"to,omitempty"`
}

// InventoryAdjustmentReportParams filters
// GET /partner/v1/inventory/adjustments/report (default page size 25).
type InventoryAdjustmentReportParams struct {
	ListParams
	// LocationIDs is sent comma-joined (location_id=a,b).
	LocationIDs []string `json:"location_id,omitempty"`
	From        string   `json:"from,omitempty"`
	To          string   `json:"to,omitempty"`
}

// InventoryAdjustmentLineRequest is one line of an adjustment create.
type InventoryAdjustmentLineRequest struct {
	ProductID string `json:"product_id,omitempty"`
	// Quantity is a signed decimal string.
	Quantity string `json:"quantity,omitempty"`
	UnitCost string `json:"unit_cost,omitempty"`
	Notes    string `json:"notes,omitempty"`
}

// InventoryAdjustmentCreateRequest is POST /partner/v1/inventory/adjustments.
type InventoryAdjustmentCreateRequest struct {
	LocationID string                           `json:"location_id,omitempty"`
	Notes      string                           `json:"notes,omitempty"`
	Lines      []InventoryAdjustmentLineRequest `json:"lines,omitempty"`
}

// InventoryTransfer moves stock from one location to another.
type InventoryTransfer struct {
	ID               string    `json:"id"`
	TenantID         string    `json:"tenant_id"`
	TransferNo       string    `json:"transfer_no"`
	FromLocationID   string    `json:"from_location_id"`
	ToLocationID     string    `json:"to_location_id"`
	Notes            string    `json:"notes,omitempty"`
	CreatedBy        string    `json:"created_by,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	FromLocationCode string    `json:"from_location_code,omitempty"`
	FromLocationName string    `json:"from_location_name,omitempty"`
	ToLocationCode   string    `json:"to_location_code,omitempty"`
	ToLocationName   string    `json:"to_location_name,omitempty"`
	ActorLabel       string    `json:"actor_label,omitempty"`
	LineCount        int       `json:"line_count,omitempty"`
}

// InventoryTransferLine is one transferred product (source and destination
// legs paired). Quantity is the positive amount moved.
type InventoryTransferLine struct {
	ID            string    `json:"id"`
	ProductID     string    `json:"product_id"`
	ProductSKU    string    `json:"product_sku,omitempty"`
	ProductName   string    `json:"product_name,omitempty"`
	Quantity      string    `json:"quantity"`
	UnitCost      string    `json:"unit_cost"`
	ExtendedCost  string    `json:"extended_cost"`
	FromQtyBefore string    `json:"from_qty_before"`
	FromQtyAfter  string    `json:"from_qty_after"`
	ToQtyBefore   string    `json:"to_qty_before"`
	ToQtyAfter    string    `json:"to_qty_after"`
	Notes         string    `json:"notes,omitempty"`
	OccurredAt    time.Time `json:"occurred_at"`
}

// InventoryTransferDetail is a transfer with its lines (create and get replies
// are {"transfer": {...}, "lines": [...]}).
type InventoryTransferDetail struct {
	Transfer InventoryTransfer       `json:"transfer"`
	Lines    []InventoryTransferLine `json:"lines"`
}

// InventoryTransferReportRow is one denormalized line of the transfers report.
type InventoryTransferReportRow struct {
	TransferID       string    `json:"transfer_id"`
	TransferNo       string    `json:"transfer_no"`
	MovementID       string    `json:"movement_id"`
	ProductID        string    `json:"product_id"`
	ProductSKU       string    `json:"product_sku"`
	ProductName      string    `json:"product_name"`
	FromLocationID   string    `json:"from_location_id"`
	FromLocationCode string    `json:"from_location_code"`
	FromLocationName string    `json:"from_location_name"`
	ToLocationID     string    `json:"to_location_id"`
	ToLocationCode   string    `json:"to_location_code"`
	ToLocationName   string    `json:"to_location_name"`
	ActorLabel       string    `json:"actor_label"`
	Quantity         string    `json:"quantity"`
	UnitCost         string    `json:"unit_cost"`
	ExtendedCost     string    `json:"extended_cost"`
	OccurredAt       time.Time `json:"occurred_at"`
	LineNotes        string    `json:"line_notes,omitempty"`
	DocumentNotes    string    `json:"document_notes,omitempty"`
}

// InventoryTransferListParams filters GET /partner/v1/inventory/transfers
// (default page size 25). From/To are RFC3339 or YYYY-MM-DD.
type InventoryTransferListParams struct {
	ListParams
	FromLocationID string `json:"from_location_id,omitempty"`
	ToLocationID   string `json:"to_location_id,omitempty"`
	From           string `json:"from,omitempty"`
	To             string `json:"to,omitempty"`
}

// InventoryTransferReportParams filters
// GET /partner/v1/inventory/transfers/report (default page size 25).
type InventoryTransferReportParams struct {
	ListParams
	// LocationIDs is sent comma-joined (location_id=a,b).
	LocationIDs []string `json:"location_id,omitempty"`
	From        string   `json:"from,omitempty"`
	To          string   `json:"to,omitempty"`
}

// InventoryTransferLineRequest is one line of a transfer create.
type InventoryTransferLineRequest struct {
	ProductID string `json:"product_id,omitempty"`
	Quantity  string `json:"quantity,omitempty"`
	Notes     string `json:"notes,omitempty"`
}

// InventoryTransferCreateRequest is POST /partner/v1/inventory/transfers.
type InventoryTransferCreateRequest struct {
	FromLocationID string                         `json:"from_location_id,omitempty"`
	ToLocationID   string                         `json:"to_location_id,omitempty"`
	Notes          string                         `json:"notes,omitempty"`
	Lines          []InventoryTransferLineRequest `json:"lines,omitempty"`
}

// InventoryService calls /partner/v1/inventory-settings and
// /partner/v1/inventory/{adjustments,transfers}.
type InventoryService struct{ c *Client }

// GetSettings returns the inventory settings. Scope
// partner-api.inventory-settings.read.
func (s *InventoryService) GetSettings(ctx context.Context) (*InventorySettings, error) {
	return call[InventorySettings](ctx, s.c, http.MethodGet, "/partner/v1/inventory-settings", nil, nil)
}

// UpdateSettings updates the inventory settings (PUT, only non-nil fields are
// sent). Scope partner-api.inventory-settings.write.
func (s *InventoryService) UpdateSettings(ctx context.Context, body InventorySettingsUpdateRequest) (*InventorySettings, error) {
	return call[InventorySettings](ctx, s.c, http.MethodPut, "/partner/v1/inventory-settings", nil, body)
}

// ListAdjustments returns inventory adjustments. Scope
// partner-api.inventory-adjustments.read.
func (s *InventoryService) ListAdjustments(ctx context.Context, params InventoryAdjustmentListParams) (*Page[InventoryAdjustment], error) {
	return listPage[InventoryAdjustment](ctx, s.c, "/partner/v1/inventory/adjustments", "adjustments", params)
}

// CreateAdjustment posts an adjustment document. Scope
// partner-api.inventory-adjustments.write. Pass WithIdempotencyKey to make
// retries safe.
func (s *InventoryService) CreateAdjustment(ctx context.Context, body InventoryAdjustmentCreateRequest, opts ...RequestOption) (*InventoryAdjustmentDetail, error) {
	return call[InventoryAdjustmentDetail](ctx, s.c, http.MethodPost, "/partner/v1/inventory/adjustments", nil, body, opts...)
}

// AdjustmentsReport returns the denormalized adjustment lines report. Scope
// partner-api.inventory-adjustments.read.
func (s *InventoryService) AdjustmentsReport(ctx context.Context, params InventoryAdjustmentReportParams) (*Page[InventoryAdjustmentReportRow], error) {
	return listPage[InventoryAdjustmentReportRow](ctx, s.c, "/partner/v1/inventory/adjustments/report", "rows", params)
}

// GetAdjustment returns one adjustment with its lines. Scope
// partner-api.inventory-adjustments.read.
func (s *InventoryService) GetAdjustment(ctx context.Context, id string) (*InventoryAdjustmentDetail, error) {
	return call[InventoryAdjustmentDetail](ctx, s.c, http.MethodGet, "/partner/v1/inventory/adjustments/"+seg(id), nil, nil)
}

// ListTransfers returns inventory transfers. Scope
// partner-api.inventory-transfers.read.
func (s *InventoryService) ListTransfers(ctx context.Context, params InventoryTransferListParams) (*Page[InventoryTransfer], error) {
	return listPage[InventoryTransfer](ctx, s.c, "/partner/v1/inventory/transfers", "transfers", params)
}

// CreateTransfer posts a transfer document. Scope
// partner-api.inventory-transfers.write. Pass WithIdempotencyKey to make
// retries safe.
func (s *InventoryService) CreateTransfer(ctx context.Context, body InventoryTransferCreateRequest, opts ...RequestOption) (*InventoryTransferDetail, error) {
	return call[InventoryTransferDetail](ctx, s.c, http.MethodPost, "/partner/v1/inventory/transfers", nil, body, opts...)
}

// TransfersReport returns the denormalized transfer lines report. Scope
// partner-api.inventory-transfers.read.
func (s *InventoryService) TransfersReport(ctx context.Context, params InventoryTransferReportParams) (*Page[InventoryTransferReportRow], error) {
	return listPage[InventoryTransferReportRow](ctx, s.c, "/partner/v1/inventory/transfers/report", "rows", params)
}

// GetTransfer returns one transfer with its lines. Scope
// partner-api.inventory-transfers.read.
func (s *InventoryService) GetTransfer(ctx context.Context, id string) (*InventoryTransferDetail, error) {
	return call[InventoryTransferDetail](ctx, s.c, http.MethodGet, "/partner/v1/inventory/transfers/"+seg(id), nil, nil)
}
