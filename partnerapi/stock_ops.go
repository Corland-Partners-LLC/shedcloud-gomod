package partnerapi

import (
	"context"
	"net/http"
	"time"
)

// ---- Sales returns (RMA) ----

// SalesReturn is a customer RMA. Money is a decimal string.
type SalesReturn struct {
	ID       string `json:"id"`
	TenantID string `json:"tenant_id"`
	Number   string `json:"number,omitempty"`
	// Status: draft | approved | completed | cancelled.
	Status          string     `json:"status"`
	OrderID         string     `json:"order_id"`
	CustomerID      string     `json:"customer_id,omitempty"`
	StoreLocationID string     `json:"store_location_id,omitempty"`
	Reason          string     `json:"reason,omitempty"`
	Notes           string     `json:"notes,omitempty"`
	RestockTotal    string     `json:"restock_total"`
	CreditTotal     string     `json:"credit_total"`
	CreatedBy       string     `json:"created_by,omitempty"`
	ApprovedBy      string     `json:"approved_by,omitempty"`
	ApprovedAt      *time.Time `json:"approved_at,omitempty"`
	CompletedAt     *time.Time `json:"completed_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	OrderNumber     string     `json:"order_number,omitempty"`
	CustomerName    string     `json:"customer_name,omitempty"`
	LineCount       int        `json:"line_count,omitempty"`
}

// SalesReturnLine is one returned item of an RMA.
type SalesReturnLine struct {
	ID          string `json:"id"`
	TenantID    string `json:"tenant_id"`
	ReturnID    string `json:"return_id"`
	OrderLineID string `json:"order_line_id,omitempty"`
	StockUnitID string `json:"stock_unit_id,omitempty"`
	ProductID   string `json:"product_id,omitempty"`
	// Disposition: restock | scrap | credit_only.
	Disposition string    `json:"disposition"`
	Qty         string    `json:"qty"`
	Amount      string    `json:"amount"`
	Reason      string    `json:"reason,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// SalesReturnDetail is the GET /partner/v1/sales-returns/{returnId} reply. The
// ERP repeats a few header fields at the top level next to the full header
// (SalesReturn) and the lines.
type SalesReturnDetail struct {
	ID          string            `json:"id"`
	Number      string            `json:"number"`
	Status      string            `json:"status"`
	OrderID     string            `json:"order_id"`
	SalesReturn SalesReturn       `json:"sales_return"`
	Lines       []SalesReturnLine `json:"lines"`
}

// SalesReturnCreateResult is the POST /partner/v1/sales-returns reply.
type SalesReturnCreateResult struct {
	SalesReturn SalesReturn       `json:"sales_return"`
	Lines       []SalesReturnLine `json:"lines"`
}

// SalesReturnListParams filters GET /partner/v1/sales-returns.
type SalesReturnListParams struct {
	ListParams
	Status  string `json:"status,omitempty"`
	OrderID string `json:"order_id,omitempty"`
}

// SalesReturnLineRequest is one line of a sales return create.
type SalesReturnLineRequest struct {
	OrderLineID string `json:"order_line_id,omitempty"`
	StockUnitID string `json:"stock_unit_id,omitempty"`
	ProductID   string `json:"product_id,omitempty"`
	// Disposition: restock | scrap | credit_only.
	Disposition string `json:"disposition,omitempty"`
	Qty         string `json:"qty,omitempty"`
	Amount      string `json:"amount,omitempty"`
	Reason      string `json:"reason,omitempty"`
}

// SalesReturnCreateRequest is POST /partner/v1/sales-returns.
type SalesReturnCreateRequest struct {
	OrderID         string                   `json:"order_id,omitempty"`
	CustomerID      string                   `json:"customer_id,omitempty"`
	StoreLocationID string                   `json:"store_location_id,omitempty"`
	Reason          string                   `json:"reason,omitempty"`
	Notes           string                   `json:"notes,omitempty"`
	Lines           []SalesReturnLineRequest `json:"lines,omitempty"`
}

// SalesReturnsService calls /partner/v1/sales-returns.
type SalesReturnsService struct{ c *Client }

// List returns sales returns. Scope partner-api.sales-returns.read.
func (s *SalesReturnsService) List(ctx context.Context, params SalesReturnListParams) (*Page[SalesReturn], error) {
	return listPage[SalesReturn](ctx, s.c, "/partner/v1/sales-returns", "sales_returns", params)
}

// Create opens a sales return (RMA). Scope partner-api.sales-returns.write.
// Pass WithIdempotencyKey to make retries safe.
func (s *SalesReturnsService) Create(ctx context.Context, body SalesReturnCreateRequest, opts ...RequestOption) (*SalesReturnCreateResult, error) {
	return call[SalesReturnCreateResult](ctx, s.c, http.MethodPost, "/partner/v1/sales-returns", nil, body, opts...)
}

// Get returns one sales return with its lines. Scope
// partner-api.sales-returns.read.
func (s *SalesReturnsService) Get(ctx context.Context, returnID string) (*SalesReturnDetail, error) {
	return call[SalesReturnDetail](ctx, s.c, http.MethodGet, "/partner/v1/sales-returns/"+seg(returnID), nil, nil)
}

// ---- Stock reconciliations ----

// StockReconciliationSession is one walk-the-lot reconciliation session.
type StockReconciliationSession struct {
	ID           string `json:"id"`
	TenantID     string `json:"tenant_id"`
	LocationID   string `json:"location_id,omitempty"`
	AllLocations bool   `json:"all_locations"`
	// Status: active | finished | cancelled.
	Status        string     `json:"status"`
	StartedBy     string     `json:"started_by,omitempty"`
	StartedByName string     `json:"started_by_name,omitempty"`
	LocationName  string     `json:"location_name,omitempty"`
	ReviewerIDs   []string   `json:"reviewer_ids"`
	ItemCount     int        `json:"item_count"`
	DecidedCount  int        `json:"decided_count"`
	FinishedAt    *time.Time `json:"finished_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// StockReconciliationListParams filters GET /partner/v1/stock-reconciliations.
// The ERP does not read q on this route.
type StockReconciliationListParams struct {
	ListParams
	// Status: active | finished | cancelled.
	Status     string `json:"status,omitempty"`
	LocationID string `json:"location_id,omitempty"`
}

// StockReconciliationStartRequest is POST /partner/v1/stock-reconciliations.
type StockReconciliationStartRequest struct {
	LocationID   string   `json:"location_id,omitempty"`
	AllLocations bool     `json:"all_locations,omitempty"`
	ReviewerIDs  []string `json:"reviewer_ids,omitempty"`
}

// StockReconciliationsService calls /partner/v1/stock-reconciliations.
type StockReconciliationsService struct{ c *Client }

// List returns reconciliation sessions. Scope
// partner-api.stock-reconciliation.read.
func (s *StockReconciliationsService) List(ctx context.Context, params StockReconciliationListParams) (*Page[StockReconciliationSession], error) {
	return listPage[StockReconciliationSession](ctx, s.c, "/partner/v1/stock-reconciliations", "sessions", params)
}

// Start opens a reconciliation session. Scope
// partner-api.stock-reconciliation.write. Pass WithIdempotencyKey to make
// retries safe.
func (s *StockReconciliationsService) Start(ctx context.Context, body StockReconciliationStartRequest, opts ...RequestOption) (*StockReconciliationSession, error) {
	return call[StockReconciliationSession](ctx, s.c, http.MethodPost, "/partner/v1/stock-reconciliations", nil, body, opts...)
}

// Get returns one reconciliation session. Scope
// partner-api.stock-reconciliation.read.
func (s *StockReconciliationsService) Get(ctx context.Context, sessionID string) (*StockReconciliationSession, error) {
	return call[StockReconciliationSession](ctx, s.c, http.MethodGet, "/partner/v1/stock-reconciliations/"+seg(sessionID), nil, nil)
}

// ---- Cycle counts ----

// CycleCountSession is a quantity count at one business location.
type CycleCountSession struct {
	ID         string `json:"id"`
	TenantID   string `json:"tenant_id"`
	LocationID string `json:"location_id"`
	// Status: open | completed | cancelled.
	Status       string     `json:"status"`
	AdjustmentID string     `json:"adjustment_id,omitempty"`
	StartedBy    string     `json:"started_by,omitempty"`
	CompletedAt  *time.Time `json:"completed_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	LineCount    int        `json:"line_count,omitempty"`
}

// CycleCountLine is one product's system vs counted quantity (decimal strings).
type CycleCountLine struct {
	ID          string    `json:"id"`
	TenantID    string    `json:"tenant_id"`
	SessionID   string    `json:"session_id"`
	ProductID   string    `json:"product_id"`
	SystemQty   string    `json:"system_qty"`
	CountedQty  string    `json:"counted_qty,omitempty"`
	VarianceQty string    `json:"variance_qty"`
	CreatedAt   time.Time `json:"created_at"`
}

// CycleCountDetail is a cycle count session with its lines (start and get
// replies are {"session": {...}, "lines": [...]}).
type CycleCountDetail struct {
	Session CycleCountSession `json:"session"`
	Lines   []CycleCountLine  `json:"lines"`
}

// CycleCountListParams filters GET /partner/v1/cycle-counts. The ERP does not
// read q on this route.
type CycleCountListParams struct {
	ListParams
	LocationID string `json:"location_id,omitempty"`
	// Status: open | completed | cancelled.
	Status string `json:"status,omitempty"`
}

// CycleCountStartRequest is POST /partner/v1/cycle-counts.
type CycleCountStartRequest struct {
	LocationID string `json:"location_id,omitempty"`
}

// CycleCountsService calls /partner/v1/cycle-counts.
type CycleCountsService struct{ c *Client }

// List returns cycle count sessions. Scope partner-api.cycle-counts.read.
func (s *CycleCountsService) List(ctx context.Context, params CycleCountListParams) (*Page[CycleCountSession], error) {
	return listPage[CycleCountSession](ctx, s.c, "/partner/v1/cycle-counts", "cycle_counts", params)
}

// Start opens a cycle count at a location. Scope
// partner-api.cycle-counts.write. Pass WithIdempotencyKey to make retries safe.
func (s *CycleCountsService) Start(ctx context.Context, body CycleCountStartRequest, opts ...RequestOption) (*CycleCountDetail, error) {
	return call[CycleCountDetail](ctx, s.c, http.MethodPost, "/partner/v1/cycle-counts", nil, body, opts...)
}

// Get returns one cycle count with its lines. Scope
// partner-api.cycle-counts.read.
func (s *CycleCountsService) Get(ctx context.Context, sessionID string) (*CycleCountDetail, error) {
	return call[CycleCountDetail](ctx, s.c, http.MethodGet, "/partner/v1/cycle-counts/"+seg(sessionID), nil, nil)
}

// ---- Warehouse tasks ----

// WarehouseTask is a WMS putaway / pick / count task. Qty is a decimal string.
type WarehouseTask struct {
	ID       string `json:"id"`
	TenantID string `json:"tenant_id"`
	// Kind: putaway | pick | count.
	Kind string `json:"kind"`
	// Status: open | in_progress | done | cancelled.
	Status     string    `json:"status"`
	LocationID string    `json:"location_id,omitempty"`
	ProductID  string    `json:"product_id,omitempty"`
	Qty        string    `json:"qty"`
	SourceType string    `json:"source_type,omitempty"`
	SourceID   string    `json:"source_id,omitempty"`
	AssignedTo string    `json:"assigned_to,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// WarehouseTaskListParams filters GET /partner/v1/warehouse-tasks. The ERP does
// not read q on this route.
type WarehouseTaskListParams struct {
	ListParams
	// Status: open | in_progress | done | cancelled.
	Status string `json:"status,omitempty"`
}

// WarehouseTaskCreateRequest is POST /partner/v1/warehouse-tasks.
type WarehouseTaskCreateRequest struct {
	// Kind: putaway | pick | count.
	Kind       string `json:"kind,omitempty"`
	LocationID string `json:"location_id,omitempty"`
	ProductID  string `json:"product_id,omitempty"`
	Qty        string `json:"qty,omitempty"`
	SourceType string `json:"source_type,omitempty"`
	SourceID   string `json:"source_id,omitempty"`
	AssignedTo string `json:"assigned_to,omitempty"`
}

// WarehouseTasksService calls /partner/v1/warehouse-tasks.
type WarehouseTasksService struct{ c *Client }

// List returns warehouse tasks. Scope partner-api.warehouse-tasks.read.
func (s *WarehouseTasksService) List(ctx context.Context, params WarehouseTaskListParams) (*Page[WarehouseTask], error) {
	return listPage[WarehouseTask](ctx, s.c, "/partner/v1/warehouse-tasks", "warehouse_tasks", params)
}

// Create makes a warehouse task. Scope partner-api.warehouse-tasks.write. Pass
// WithIdempotencyKey to make retries safe.
func (s *WarehouseTasksService) Create(ctx context.Context, body WarehouseTaskCreateRequest, opts ...RequestOption) (*WarehouseTask, error) {
	return call[WarehouseTask](ctx, s.c, http.MethodPost, "/partner/v1/warehouse-tasks", nil, body, opts...)
}

// Get returns one warehouse task. Scope partner-api.warehouse-tasks.read.
func (s *WarehouseTasksService) Get(ctx context.Context, taskID string) (*WarehouseTask, error) {
	return call[WarehouseTask](ctx, s.c, http.MethodGet, "/partner/v1/warehouse-tasks/"+seg(taskID), nil, nil)
}
