package partnerapi

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

// WorkOrder is a production work order (/partner/v1/work-orders). Dates are
// "YYYY-MM-DD" strings; quantities are decimal strings. Nested collections
// (Components, Operations, Assignments, History, Display, Transport) are kept
// as raw JSON.
type WorkOrder struct {
	ID                       string          `json:"id"`
	TenantID                 string          `json:"tenant_id,omitempty"`
	WONumber                 string          `json:"wo_number,omitempty"`
	SerialNumber             string          `json:"serial_number,omitempty"`
	ProductID                string          `json:"product_id,omitempty"`
	ConfigCategory           string          `json:"config_category"`
	SalesIntent              string          `json:"sales_intent,omitempty"`
	DeliveryType             string          `json:"delivery_type,omitempty"`
	Priority                 int             `json:"priority"`
	Rush                     bool            `json:"rush"`
	PlantID                  string          `json:"plant_id,omitempty"`
	BuildAtID                string          `json:"build_at_id,omitempty"`
	BuildingLocationID       string          `json:"building_location_id,omitempty"`
	DestinationBuildID       string          `json:"destination_build_id,omitempty"`
	DestinationIsCustomer    bool            `json:"destination_is_customer"`
	SourceOrderID            string          `json:"source_order_id,omitempty"`
	SourceOrderNumber        string          `json:"source_order_number,omitempty"`
	DeliveryAddress          json.RawMessage `json:"delivery_address,omitempty"`
	Contact                  json.RawMessage `json:"contact,omitempty"`
	LegacyID                 string          `json:"legacy_id,omitempty"`
	SourceType               string          `json:"source_type,omitempty"`
	SourceID                 string          `json:"source_id,omitempty"`
	AssignedTo               string          `json:"assigned_to,omitempty"`
	ConfigurationID          string          `json:"configuration_id,omitempty"`
	ConfigurationRevisionID  string          `json:"configuration_revision_id,omitempty"`
	Dimensions               json.RawMessage `json:"dimensions,omitempty"`
	IntakeBasePrice          string          `json:"intake_base_price,omitempty"`
	IntakeUpgradesTotal      string          `json:"intake_upgrades_total,omitempty"`
	IntakePriceReason        string          `json:"intake_price_reason,omitempty"`
	QtyOrdered               string          `json:"qty_ordered"`
	QtyCompleted             string          `json:"qty_completed"`
	QtyScrapped              string          `json:"qty_scrapped"`
	Status                   string          `json:"status"`
	PriorStatus              string          `json:"prior_status,omitempty"`
	StructureFrozen          bool            `json:"structure_frozen"`
	EstimatedDeliveryDate    string          `json:"estimated_delivery_date,omitempty"`
	ScheduledFinishBuildDate string          `json:"scheduled_finish_build_date,omitempty"`
	ActualFinishBuildDate    string          `json:"actual_finish_build_date,omitempty"`
	ActualDeliveryDate       string          `json:"actual_delivery_date,omitempty"`
	ScheduledPickupAt        string          `json:"scheduled_pickup_at,omitempty"`
	ScheduledDeliveryAt      string          `json:"scheduled_delivery_at,omitempty"`
	TransportationCompanyID  string          `json:"transportation_company_id,omitempty"`
	ReviewedBy               string          `json:"reviewed_by,omitempty"`
	ReviewedAt               string          `json:"reviewed_at,omitempty"`
	Comments                 string          `json:"comments,omitempty"`
	HasShortage              bool            `json:"has_shortage"`
	ItemsSetAt               *time.Time      `json:"items_set_at"`
	BOMMode                  string          `json:"bom_mode"`
	IsCustom                 bool            `json:"is_custom"`
	CreatedBy                string          `json:"created_by,omitempty"`
	CreatedAt                time.Time       `json:"created_at"`
	UpdatedAt                time.Time       `json:"updated_at"`
	Components               json.RawMessage `json:"components,omitempty"`
	Operations               json.RawMessage `json:"operations,omitempty"`
	Assignments              json.RawMessage `json:"assignments,omitempty"`
	History                  json.RawMessage `json:"history,omitempty"`
	Display                  json.RawMessage `json:"display,omitempty"`
	Transport                json.RawMessage `json:"transport,omitempty"`
}

// WorkOrderListParams filters GET /partner/v1/work-orders.
type WorkOrderListParams struct {
	ListParams
	Status     string `json:"status,omitempty"`
	ProductID  string `json:"product_id,omitempty"`
	CustomerID string `json:"customer_id,omitempty"`
}

// WorkOrderCreateRequest is POST /partner/v1/work-orders.
type WorkOrderCreateRequest struct {
	ProductID                string          `json:"product_id,omitempty"`
	ConfigCategory           string          `json:"config_category,omitempty"`
	SalesIntent              string          `json:"sales_intent,omitempty"`
	DeliveryType             string          `json:"delivery_type,omitempty"`
	Priority                 int             `json:"priority,omitempty"`
	Rush                     bool            `json:"rush,omitempty"`
	PlantID                  string          `json:"plant_id,omitempty"`
	BuildAtID                string          `json:"build_at_id,omitempty"`
	BuildingLocationID       string          `json:"building_location_id,omitempty"`
	DestinationBuildID       string          `json:"destination_build_id,omitempty"`
	DestinationIsCustomer    bool            `json:"destination_is_customer,omitempty"`
	DeliveryAddress          json.RawMessage `json:"delivery_address,omitempty"`
	Contact                  json.RawMessage `json:"contact,omitempty"`
	SourceType               string          `json:"source_type,omitempty"`
	SourceID                 string          `json:"source_id,omitempty"`
	QtyOrdered               string          `json:"qty_ordered,omitempty"`
	EstimatedDeliveryDate    string          `json:"estimated_delivery_date,omitempty"`
	ScheduledFinishBuildDate string          `json:"scheduled_finish_build_date,omitempty"`
	ActualDeliveryDate       string          `json:"actual_delivery_date,omitempty"`
	TransportationCompanyID  string          `json:"transportation_company_id,omitempty"`
	Comments                 string          `json:"comments,omitempty"`
	SnapshotBOM              bool            `json:"snapshot_bom,omitempty"`
	ConfigurationID          string          `json:"configuration_id,omitempty"`
}

// WorkOrderUpdateRequest is PATCH /partner/v1/work-orders/{id}. Only non-nil
// fields are sent. Status changes go through UpdateStatus.
type WorkOrderUpdateRequest struct {
	ProductID                *string         `json:"product_id,omitempty"`
	ConfigCategory           *string         `json:"config_category,omitempty"`
	SalesIntent              *string         `json:"sales_intent,omitempty"`
	DeliveryType             *string         `json:"delivery_type,omitempty"`
	Priority                 *int            `json:"priority,omitempty"`
	Rush                     *bool           `json:"rush,omitempty"`
	PlantID                  *string         `json:"plant_id,omitempty"`
	BuildAtID                *string         `json:"build_at_id,omitempty"`
	BuildingLocationID       *string         `json:"building_location_id,omitempty"`
	DestinationBuildID       *string         `json:"destination_build_id,omitempty"`
	DestinationIsCustomer    *bool           `json:"destination_is_customer,omitempty"`
	DeliveryAddress          json.RawMessage `json:"delivery_address,omitempty"`
	Contact                  json.RawMessage `json:"contact,omitempty"`
	QtyOrdered               *string         `json:"qty_ordered,omitempty"`
	EstimatedDeliveryDate    *string         `json:"estimated_delivery_date,omitempty"`
	ScheduledFinishBuildDate *string         `json:"scheduled_finish_build_date,omitempty"`
	ActualFinishBuildDate    *string         `json:"actual_finish_build_date,omitempty"`
	ActualDeliveryDate       *string         `json:"actual_delivery_date,omitempty"`
	TransportationCompanyID  *string         `json:"transportation_company_id,omitempty"`
	Comments                 *string         `json:"comments,omitempty"`
	SourceType               *string         `json:"source_type,omitempty"`
	SourceID                 *string         `json:"source_id,omitempty"`
	AssignedTo               *string         `json:"assigned_to,omitempty"`
}

// WorkOrderStatusRequest is PATCH /partner/v1/work-orders/{id}/status.
type WorkOrderStatusRequest struct {
	Status string `json:"status"`
	Note   string `json:"note,omitempty"`
	// CapturedFields are status-specific values the transition asks for.
	CapturedFields map[string]any `json:"captured_fields,omitempty"`
}

// WorkOrderStatusChange is one entry of a work order's status history.
type WorkOrderStatusChange struct {
	ID          string    `json:"id"`
	TenantID    string    `json:"tenant_id,omitempty"`
	WorkOrderID string    `json:"work_order_id"`
	FromStatus  string    `json:"from_status,omitempty"`
	ToStatus    string    `json:"to_status"`
	Note        string    `json:"note,omitempty"`
	ChangedBy   string    `json:"changed_by,omitempty"`
	ChangedAt   time.Time `json:"changed_at"`
}

// WorkOrdersService calls /partner/v1/work-orders.
type WorkOrdersService struct{ c *Client }

// List returns work orders. Scope partner-api.work-orders.read.
func (s *WorkOrdersService) List(ctx context.Context, params WorkOrderListParams) (*Page[WorkOrder], error) {
	return listPage[WorkOrder](ctx, s.c, "/partner/v1/work-orders", "work_orders", params)
}

// Get returns one work order. Scope partner-api.work-orders.read.
func (s *WorkOrdersService) Get(ctx context.Context, id string) (*WorkOrder, error) {
	return call[WorkOrder](ctx, s.c, http.MethodGet, "/partner/v1/work-orders/"+seg(id), nil, nil)
}

// Create makes a work order. Scope partner-api.work-orders.write.
func (s *WorkOrdersService) Create(ctx context.Context, body WorkOrderCreateRequest, opts ...RequestOption) (*WorkOrder, error) {
	return call[WorkOrder](ctx, s.c, http.MethodPost, "/partner/v1/work-orders", nil, body, opts...)
}

// Update patches a work order. Scope partner-api.work-orders.write.
func (s *WorkOrdersService) Update(ctx context.Context, id string, body WorkOrderUpdateRequest) (*WorkOrder, error) {
	return call[WorkOrder](ctx, s.c, http.MethodPatch, "/partner/v1/work-orders/"+seg(id), nil, body)
}

// UpdateStatus transitions a work order's status (422 on an invalid
// transition). PATCH, not POST as in 0.x.
func (s *WorkOrdersService) UpdateStatus(ctx context.Context, id string, body WorkOrderStatusRequest) (*WorkOrder, error) {
	return call[WorkOrder](ctx, s.c, http.MethodPatch, "/partner/v1/work-orders/"+seg(id)+"/status", nil, body)
}

// StatusHistory returns the status changes of a work order (not paginated).
func (s *WorkOrdersService) StatusHistory(ctx context.Context, id string) ([]WorkOrderStatusChange, error) {
	return listKey[WorkOrderStatusChange](ctx, s.c, http.MethodGet, "/partner/v1/work-orders/"+seg(id)+"/status-history", "history", nil)
}
