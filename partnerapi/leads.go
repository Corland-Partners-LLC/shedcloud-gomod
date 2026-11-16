package partnerapi

import (
	"context"
	"net/http"
	"time"
)

// Lead is a CRM lead (/partner/v1/leads). Money is a decimal string.
type Lead struct {
	ID                string         `json:"id"`
	TenantID          string         `json:"tenant_id,omitempty"`
	Name              string         `json:"name"`
	FirstName         string         `json:"first_name,omitempty"`
	LastName          string         `json:"last_name,omitempty"`
	Email             string         `json:"email,omitempty"`
	Phone             string         `json:"phone,omitempty"`
	MobilePhone       string         `json:"mobile_phone,omitempty"`
	CompanyName       string         `json:"company_name,omitempty"`
	Website           string         `json:"website,omitempty"`
	JobTitle          string         `json:"job_title,omitempty"`
	Status            string         `json:"status"`
	LifecycleStage    string         `json:"lifecycle_stage"`
	SourceSlug        string         `json:"source_slug,omitempty"`
	SourceName        string         `json:"source_name,omitempty"`
	LeadScore         int            `json:"lead_score"`
	Priority          string         `json:"priority"`
	Rating            string         `json:"rating,omitempty"`
	AssignedTo        string         `json:"assigned_to,omitempty"`
	AssignedName      string         `json:"assigned_name,omitempty"`
	AssignedEmail     string         `json:"assigned_email,omitempty"`
	CustomerID        string         `json:"customer_id,omitempty"`
	CustomerName      string         `json:"customer_name,omitempty"`
	EstimatedValue    string         `json:"estimated_value,omitempty"`
	Currency          string         `json:"currency"`
	ExpectedCloseDate *time.Time     `json:"expected_close_date,omitempty"`
	AddressLine1      string         `json:"address_line1,omitempty"`
	AddressLine2      string         `json:"address_line2,omitempty"`
	City              string         `json:"city,omitempty"`
	Region            string         `json:"region,omitempty"`
	PostalCode        string         `json:"postal_code,omitempty"`
	Country           string         `json:"country,omitempty"`
	LastContactedAt   *time.Time     `json:"last_contacted_at,omitempty"`
	NextFollowUpAt    *time.Time     `json:"next_follow_up_at,omitempty"`
	ConvertedAt       *time.Time     `json:"converted_at,omitempty"`
	Notes             string         `json:"notes,omitempty"`
	Metadata          map[string]any `json:"metadata,omitempty"`
	CreatedBy         string         `json:"created_by,omitempty"`
	UpdatedBy         string         `json:"updated_by,omitempty"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	ContactID         *string        `json:"contact_id,omitempty"`
	OrganizationID    *string        `json:"organization_id,omitempty"`
	ExternalRef       *string        `json:"external_ref,omitempty"`
	ReceivedAt        *time.Time     `json:"received_at,omitempty"`
	FirstContactedAt  *time.Time     `json:"first_contacted_at,omitempty"`
	Subject           *string        `json:"subject,omitempty"`
	Interest          *string        `json:"interest,omitempty"`
	SourceDetail      *string        `json:"source_detail,omitempty"`
	ArchivedAt        *time.Time     `json:"archived_at,omitempty"`
	Revision          int64          `json:"revision,omitempty"`
	PriorityRank      int            `json:"priority_rank,omitempty"`
}

// LeadListParams filters GET /partner/v1/leads.
type LeadListParams struct {
	ListParams
	Status string `json:"status,omitempty"`
}

// LeadCreateRequest is POST /partner/v1/leads. Unknown fields are rejected by
// the ERP (400), so only the fields below exist.
type LeadCreateRequest struct {
	Name              string         `json:"name,omitempty"`
	FirstName         string         `json:"first_name,omitempty"`
	LastName          string         `json:"last_name,omitempty"`
	Email             string         `json:"email,omitempty"`
	Phone             string         `json:"phone,omitempty"`
	MobilePhone       string         `json:"mobile_phone,omitempty"`
	CompanyName       string         `json:"company_name,omitempty"`
	Website           string         `json:"website,omitempty"`
	JobTitle          string         `json:"job_title,omitempty"`
	Status            string         `json:"status,omitempty"`
	LifecycleStage    string         `json:"lifecycle_stage,omitempty"`
	SourceSlug        string         `json:"source_slug,omitempty"`
	LeadScore         int            `json:"lead_score,omitempty"`
	Priority          string         `json:"priority,omitempty"`
	Rating            string         `json:"rating,omitempty"`
	AssignedTo        string         `json:"assigned_to,omitempty"`
	EstimatedValue    string         `json:"estimated_value,omitempty"`
	Currency          string         `json:"currency,omitempty"`
	ExpectedCloseDate string         `json:"expected_close_date,omitempty"`
	AddressLine1      string         `json:"address_line1,omitempty"`
	AddressLine2      string         `json:"address_line2,omitempty"`
	City              string         `json:"city,omitempty"`
	Region            string         `json:"region,omitempty"`
	PostalCode        string         `json:"postal_code,omitempty"`
	Country           string         `json:"country,omitempty"`
	Notes             string         `json:"notes,omitempty"`
	Metadata          map[string]any `json:"metadata,omitempty"`
}

// LeadUpdateRequest is PATCH /partner/v1/leads/{id}. Only non-nil fields are
// sent. Status changes go through UpdateStatus.
type LeadUpdateRequest struct {
	Name           *string        `json:"name,omitempty"`
	FirstName      *string        `json:"first_name,omitempty"`
	LastName       *string        `json:"last_name,omitempty"`
	Email          *string        `json:"email,omitempty"`
	Phone          *string        `json:"phone,omitempty"`
	MobilePhone    *string        `json:"mobile_phone,omitempty"`
	CompanyName    *string        `json:"company_name,omitempty"`
	Website        *string        `json:"website,omitempty"`
	JobTitle       *string        `json:"job_title,omitempty"`
	LifecycleStage *string        `json:"lifecycle_stage,omitempty"`
	SourceSlug     *string        `json:"source_slug,omitempty"`
	LeadScore      *int           `json:"lead_score,omitempty"`
	Priority       *string        `json:"priority,omitempty"`
	Rating         *string        `json:"rating,omitempty"`
	AssignedTo     *string        `json:"assigned_to,omitempty"`
	EstimatedValue *string        `json:"estimated_value,omitempty"`
	Currency       *string        `json:"currency,omitempty"`
	AddressLine1   *string        `json:"address_line1,omitempty"`
	AddressLine2   *string        `json:"address_line2,omitempty"`
	City           *string        `json:"city,omitempty"`
	Region         *string        `json:"region,omitempty"`
	PostalCode     *string        `json:"postal_code,omitempty"`
	Country        *string        `json:"country,omitempty"`
	Notes          *string        `json:"notes,omitempty"`
	Metadata       map[string]any `json:"metadata,omitempty"`
}

// LeadStatusRequest is PATCH /partner/v1/leads/{id}/status.
type LeadStatusRequest struct {
	Status string `json:"status"`
	Note   string `json:"note,omitempty"`
}

// LeadConvertRequest is POST /partner/v1/leads/{id}/convert. Leave
// CreateCustomer nil for the ERP default (create a customer from the lead
// when CustomerID is empty).
type LeadConvertRequest struct {
	CustomerID     string `json:"customer_id,omitempty"`
	CreateCustomer *bool  `json:"create_customer,omitempty"`
}

// LeadConvertResult is the convert reply. The ERP serializes it with Go field
// names (Lead, Customer, Created); encoding/json matches them
// case-insensitively.
type LeadConvertResult struct {
	Lead     Lead     `json:"lead"`
	Customer Customer `json:"customer"`
	// Created is true when a new customer was created.
	Created bool `json:"created"`
}

// LeadStatusHistoryEntry is one status change of a lead.
type LeadStatusHistoryEntry struct {
	ID         string    `json:"id"`
	LeadID     string    `json:"lead_id"`
	FromStatus string    `json:"from_status"`
	ToStatus   string    `json:"to_status"`
	Note       string    `json:"note,omitempty"`
	ChangedBy  string    `json:"changed_by,omitempty"`
	ChangedAt  time.Time `json:"changed_at"`
	Source     string    `json:"source,omitempty"`
}

// LeadsService calls /partner/v1/leads.
type LeadsService struct{ c *Client }

// List returns leads. Scope partner-api.leads.read.
func (s *LeadsService) List(ctx context.Context, params LeadListParams) (*Page[Lead], error) {
	return listPage[Lead](ctx, s.c, "/partner/v1/leads", "leads", params)
}

// Create makes a lead. Scope partner-api.leads.write. Pass
// WithIdempotencyKey to make retries safe.
func (s *LeadsService) Create(ctx context.Context, body LeadCreateRequest, opts ...RequestOption) (*Lead, error) {
	return call[Lead](ctx, s.c, http.MethodPost, "/partner/v1/leads", nil, body, opts...)
}

// Get returns one lead. Scope partner-api.leads.read.
func (s *LeadsService) Get(ctx context.Context, id string) (*Lead, error) {
	return call[Lead](ctx, s.c, http.MethodGet, "/partner/v1/leads/"+seg(id), nil, nil)
}

// Update patches a lead. Scope partner-api.leads.write.
func (s *LeadsService) Update(ctx context.Context, id string, body LeadUpdateRequest) (*Lead, error) {
	return call[Lead](ctx, s.c, http.MethodPatch, "/partner/v1/leads/"+seg(id), nil, body)
}

// UpdateStatus transitions a lead's status (PATCH, not POST as in 0.x).
func (s *LeadsService) UpdateStatus(ctx context.Context, id string, body LeadStatusRequest) (*Lead, error) {
	return call[Lead](ctx, s.c, http.MethodPatch, "/partner/v1/leads/"+seg(id)+"/status", nil, body)
}

// Assign sets the lead's owner (user id). Scope partner-api.leads.write.
func (s *LeadsService) Assign(ctx context.Context, id, assignedTo string) (*Lead, error) {
	body := struct {
		AssignedTo string `json:"assigned_to"`
	}{assignedTo}
	return call[Lead](ctx, s.c, http.MethodPatch, "/partner/v1/leads/"+seg(id)+"/assign", nil, body)
}

// Convert converts a lead into a customer. Scope partner-api.leads.write.
func (s *LeadsService) Convert(ctx context.Context, id string, body LeadConvertRequest) (*LeadConvertResult, error) {
	return call[LeadConvertResult](ctx, s.c, http.MethodPost, "/partner/v1/leads/"+seg(id)+"/convert", nil, body)
}

// StatusHistory returns the lead's status changes (not paginated).
func (s *LeadsService) StatusHistory(ctx context.Context, id string) ([]LeadStatusHistoryEntry, error) {
	return listKey[LeadStatusHistoryEntry](ctx, s.c, http.MethodGet, "/partner/v1/leads/"+seg(id)+"/status-history", "history", nil)
}
