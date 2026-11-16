package partnerapi

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

// ServiceTicket is a service (warranty/repair) ticket
// (/partner/v1/service-tickets). Money is a decimal string.
type ServiceTicket struct {
	ID                 string `json:"id"`
	TenantID           string `json:"tenant_id"`
	TicketNumber       string `json:"ticket_number,omitempty"`
	WorkOrderID        string `json:"work_order_id"`
	SerialNumber       string `json:"serial_number,omitempty"`
	ProductID          string `json:"product_id,omitempty"`
	ConfigurationID    string `json:"configuration_id,omitempty"`
	PlantID            string `json:"plant_id,omitempty"`
	BuildingLocationID string `json:"building_location_id,omitempty"`
	BuildingSize       string `json:"building_size,omitempty"`

	SourceType string `json:"source_type,omitempty"`
	SourceID   string `json:"source_id,omitempty"`
	// Contact and ServiceAddress are free-form JSON blobs.
	Contact        json.RawMessage `json:"contact,omitempty"`
	ServiceAddress json.RawMessage `json:"service_address,omitempty"`

	SubmittedType   string `json:"submitted_type,omitempty"`
	CategoryOfIssue string `json:"category_of_issue,omitempty"`
	InventoryType   string `json:"inventory_type,omitempty"`
	Priority        string `json:"priority"`
	Description     string `json:"description,omitempty"`

	AssignedTechID string `json:"assigned_tech_id,omitempty"`
	TechName       string `json:"tech_name,omitempty"`
	TechEmail      string `json:"tech_email,omitempty"`
	TechPhone      string `json:"tech_phone,omitempty"`
	SubmittedBy    string `json:"submitted_by,omitempty"`

	Manufacturer      string `json:"manufacturer,omitempty"`
	DeliveryInstaller string `json:"delivery_installer,omitempty"`
	NotesText         string `json:"notes_text,omitempty"`

	Status          string `json:"status"`
	PriorStatus     string `json:"prior_status,omitempty"`
	StatusEnteredAt string `json:"status_entered_at,omitempty"`

	ScheduledServiceDate string `json:"scheduled_service_date,omitempty"`
	SolvedDate           string `json:"solved_date,omitempty"`
	BalanceDue           string `json:"balance_due"`

	EntitlementStartsOn string `json:"entitlement_starts_on,omitempty"`
	EntitlementEndsOn   string `json:"entitlement_ends_on,omitempty"`
	WarrantyMonths      *int   `json:"warranty_months,omitempty"`
	SLADueAt            string `json:"sla_due_at,omitempty"`
	SLABreached         bool   `json:"sla_breached"`
	SuppliesConsumed    bool   `json:"supplies_consumed"`

	CyflexTicketID     string `json:"cyflex_ticket_id,omitempty"`
	CyflexTicketNumber string `json:"cyflex_ticket_number,omitempty"`
	CyflexProviderName string `json:"cyflex_provider_name,omitempty"`

	CreatedBy string    `json:"created_by,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Supplies []ServiceTicketSupply `json:"supplies,omitempty"`
	Display  *ServiceTicketDisplay `json:"display,omitempty"`
}

// ServiceTicketDisplay is the denormalized display block of a ticket.
type ServiceTicketDisplay struct {
	WorkOrderNumber string `json:"work_order_number,omitempty"`
	ProductName     string `json:"product_name,omitempty"`
	CustomerName    string `json:"customer_name,omitempty"`
	TechName        string `json:"tech_name,omitempty"`
	PlantName       string `json:"plant_name,omitempty"`
	StatusEnteredAt string `json:"status_entered_at,omitempty"`
}

// ServiceTicketSupply is one supply line consumed by a ticket. Quantities and
// money are decimal strings.
type ServiceTicketSupply struct {
	ID              string    `json:"id"`
	TenantID        string    `json:"tenant_id"`
	ServiceTicketID string    `json:"service_ticket_id"`
	LineNo          int       `json:"line_no"`
	ProductID       string    `json:"product_id,omitempty"`
	Description     string    `json:"description"`
	Quantity        string    `json:"quantity"`
	UnitOfMeasure   string    `json:"unit_of_measure"`
	UnitCost        string    `json:"unit_cost"`
	LineTotal       string    `json:"line_total"`
	InventoryType   string    `json:"inventory_type,omitempty"`
	Notes           string    `json:"notes,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}

// ServiceTicketListParams filters GET /partner/v1/service-tickets. Sort is the
// ERP sort key. The ERP does not send a pagination block for this list, so
// Page.Pagination stays zero.
type ServiceTicketListParams struct {
	ListParams
	Status      string `json:"status,omitempty"`
	WorkOrderID string `json:"work_order_id,omitempty"`
	Sort        string `json:"sort,omitempty"`
}

// ServiceTicketCreateRequest is POST /partner/v1/service-tickets. Unknown
// fields are rejected by the ERP (400).
type ServiceTicketCreateRequest struct {
	WorkOrderID     string `json:"work_order_id,omitempty"`
	SubmittedType   string `json:"submitted_type,omitempty"`
	CategoryOfIssue string `json:"category_of_issue,omitempty"`
	InventoryType   string `json:"inventory_type,omitempty"`
	Priority        string `json:"priority,omitempty"`
	Description     string `json:"description,omitempty"`
	AssignedTechID  string `json:"assigned_tech_id,omitempty"`
	TechName        string `json:"tech_name,omitempty"`
	TechEmail       string `json:"tech_email,omitempty"`
	TechPhone       string `json:"tech_phone,omitempty"`
	// ScheduledDate is the service date (YYYY-MM-DD).
	ScheduledDate string `json:"scheduled_date,omitempty"`
}

// ServiceTicketUpdateRequest is PATCH /partner/v1/service-tickets/{id}.
// Tri-state: nil keeps the value, a pointer to "" clears it, a value sets it.
// Status changes go through UpdateStatus. The contact and service_address
// blobs are not exposed.
type ServiceTicketUpdateRequest struct {
	SubmittedType        *string `json:"submitted_type,omitempty"`
	CategoryOfIssue      *string `json:"category_of_issue,omitempty"`
	InventoryType        *string `json:"inventory_type,omitempty"`
	Priority             *string `json:"priority,omitempty"`
	Description          *string `json:"description,omitempty"`
	AssignedTechID       *string `json:"assigned_tech_id,omitempty"`
	TechName             *string `json:"tech_name,omitempty"`
	TechEmail            *string `json:"tech_email,omitempty"`
	TechPhone            *string `json:"tech_phone,omitempty"`
	Manufacturer         *string `json:"manufacturer,omitempty"`
	DeliveryInstaller    *string `json:"delivery_installer,omitempty"`
	NotesText            *string `json:"notes_text,omitempty"`
	ScheduledServiceDate *string `json:"scheduled_service_date,omitempty"`
	SolvedDate           *string `json:"solved_date,omitempty"`
	BalanceDue           *string `json:"balance_due,omitempty"`
	SourceType           *string `json:"source_type,omitempty"`
	SourceID             *string `json:"source_id,omitempty"`
	EntitlementStartsOn  *string `json:"entitlement_starts_on,omitempty"`
	EntitlementEndsOn    *string `json:"entitlement_ends_on,omitempty"`
	WarrantyMonths       *int    `json:"warranty_months,omitempty"`
	SLADueAt             *string `json:"sla_due_at,omitempty"`
	SLABreached          *bool   `json:"sla_breached,omitempty"`
}

// ServiceTicketsService calls /partner/v1/service-tickets.
type ServiceTicketsService struct{ c *Client }

// List returns service tickets. Scope partner-api.service-tickets.read.
func (s *ServiceTicketsService) List(ctx context.Context, params ServiceTicketListParams) (*Page[ServiceTicket], error) {
	return listPage[ServiceTicket](ctx, s.c, "/partner/v1/service-tickets", "service_tickets", params)
}

// Get returns one ticket. Scope partner-api.service-tickets.read.
func (s *ServiceTicketsService) Get(ctx context.Context, id string) (*ServiceTicket, error) {
	return call[ServiceTicket](ctx, s.c, http.MethodGet, "/partner/v1/service-tickets/"+seg(id), nil, nil)
}

// Create opens a ticket. Scope partner-api.service-tickets.write. Pass
// WithIdempotencyKey to make retries safe.
func (s *ServiceTicketsService) Create(ctx context.Context, body ServiceTicketCreateRequest, opts ...RequestOption) (*ServiceTicket, error) {
	return call[ServiceTicket](ctx, s.c, http.MethodPost, "/partner/v1/service-tickets", nil, body, opts...)
}

// Update patches a ticket header. Scope partner-api.service-tickets.write.
func (s *ServiceTicketsService) Update(ctx context.Context, id string, body ServiceTicketUpdateRequest) (*ServiceTicket, error) {
	return call[ServiceTicket](ctx, s.c, http.MethodPatch, "/partner/v1/service-tickets/"+seg(id), nil, body)
}

// UpdateStatus moves a ticket to another status key. Scope
// partner-api.service-tickets.write.
func (s *ServiceTicketsService) UpdateStatus(ctx context.Context, id, status string) (*ServiceTicket, error) {
	body := struct {
		Status string `json:"status"`
	}{status}
	return call[ServiceTicket](ctx, s.c, http.MethodPatch, "/partner/v1/service-tickets/"+seg(id)+"/status", nil, body)
}
