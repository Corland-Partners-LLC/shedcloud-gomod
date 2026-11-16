package partnerapi

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

// Customer is a CRM customer (/partner/v1/customers).
type Customer struct {
	ID                string         `json:"id"`
	TenantID          string         `json:"tenant_id,omitempty"`
	Name              string         `json:"name"`
	Email             string         `json:"email,omitempty"`
	Phone             string         `json:"phone,omitempty"`
	CompanyName       string         `json:"company_name,omitempty"`
	SourceSlug        string         `json:"source_slug,omitempty"`
	PreferredLanguage string         `json:"preferred_language"`
	Status            string         `json:"status"`
	AddressLine1      string         `json:"address_line1,omitempty"`
	AddressLine2      string         `json:"address_line2,omitempty"`
	City              string         `json:"city,omitempty"`
	Region            string         `json:"region,omitempty"`
	PostalCode        string         `json:"postal_code,omitempty"`
	Country           string         `json:"country,omitempty"`
	Notes             string         `json:"notes,omitempty"`
	Metadata          map[string]any `json:"metadata,omitempty"`
	CreatedBy         string         `json:"created_by,omitempty"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	ContactID         *string        `json:"contact_id,omitempty"`
}

// CustomerListParams filters GET /partner/v1/customers (free-text Q only).
type CustomerListParams struct {
	ListParams
}

// CustomerCreateRequest is POST /partner/v1/customers. A duplicate email is a
// 409.
type CustomerCreateRequest struct {
	Name              string         `json:"name,omitempty"`
	Email             string         `json:"email,omitempty"`
	Phone             string         `json:"phone,omitempty"`
	CompanyName       string         `json:"company_name,omitempty"`
	SourceSlug        string         `json:"source_slug,omitempty"`
	PreferredLanguage string         `json:"preferred_language,omitempty"`
	Status            string         `json:"status,omitempty"`
	AddressLine1      string         `json:"address_line1,omitempty"`
	AddressLine2      string         `json:"address_line2,omitempty"`
	City              string         `json:"city,omitempty"`
	Region            string         `json:"region,omitempty"`
	PostalCode        string         `json:"postal_code,omitempty"`
	Country           string         `json:"country,omitempty"`
	Notes             string         `json:"notes,omitempty"`
	Metadata          map[string]any `json:"metadata,omitempty"`
}

// CustomerUpdateRequest is PATCH /partner/v1/customers/{id}. Only non-nil
// fields are sent.
type CustomerUpdateRequest struct {
	Name              *string        `json:"name,omitempty"`
	Email             *string        `json:"email,omitempty"`
	Phone             *string        `json:"phone,omitempty"`
	CompanyName       *string        `json:"company_name,omitempty"`
	SourceSlug        *string        `json:"source_slug,omitempty"`
	PreferredLanguage *string        `json:"preferred_language,omitempty"`
	Status            *string        `json:"status,omitempty"`
	AddressLine1      *string        `json:"address_line1,omitempty"`
	AddressLine2      *string        `json:"address_line2,omitempty"`
	City              *string        `json:"city,omitempty"`
	Region            *string        `json:"region,omitempty"`
	PostalCode        *string        `json:"postal_code,omitempty"`
	Country           *string        `json:"country,omitempty"`
	Notes             *string        `json:"notes,omitempty"`
	Metadata          map[string]any `json:"metadata,omitempty"`
}

// CustomerLedger is GET /partner/v1/customers/{id}/ledger. Money is a decimal
// string; Entries stay raw JSON.
type CustomerLedger struct {
	CustomerID     string          `json:"customer_id"`
	CustomerName   string          `json:"customer_name,omitempty"`
	Entries        json.RawMessage `json:"entries,omitempty"`
	TotalCharged   string          `json:"total_charged"`
	TotalPaid      string          `json:"total_paid"`
	TotalAllocated string          `json:"total_allocated"`
	CreditBalance  string          `json:"credit_balance"`
	PendingRefunds string          `json:"pending_refunds"`
	OpenBalance    string          `json:"open_balance"`
	NetBalance     string          `json:"net_balance"`
}

// CustomersService calls /partner/v1/customers.
type CustomersService struct{ c *Client }

// List returns customers. Scope partner-api.customers.read.
func (s *CustomersService) List(ctx context.Context, params CustomerListParams) (*Page[Customer], error) {
	return listPage[Customer](ctx, s.c, "/partner/v1/customers", "customers", params)
}

// Get returns one customer. Scope partner-api.customers.read.
func (s *CustomersService) Get(ctx context.Context, id string) (*Customer, error) {
	return call[Customer](ctx, s.c, http.MethodGet, "/partner/v1/customers/"+seg(id), nil, nil)
}

// Create makes a customer. Scope partner-api.customers.write. Pass
// WithIdempotencyKey to make retries safe.
func (s *CustomersService) Create(ctx context.Context, body CustomerCreateRequest, opts ...RequestOption) (*Customer, error) {
	return call[Customer](ctx, s.c, http.MethodPost, "/partner/v1/customers", nil, body, opts...)
}

// Update patches a customer. Scope partner-api.customers.write.
func (s *CustomersService) Update(ctx context.Context, id string, body CustomerUpdateRequest) (*Customer, error) {
	return call[Customer](ctx, s.c, http.MethodPatch, "/partner/v1/customers/"+seg(id), nil, body)
}

// Payments lists the payments of one customer. Scope partner-api.payments.read.
func (s *CustomersService) Payments(ctx context.Context, id string, params ListParams) (*Page[Payment], error) {
	return listPage[Payment](ctx, s.c, "/partner/v1/customers/"+seg(id)+"/payments", "payments", params)
}

// CreatePayment records a payment for a customer. Use ReferenceType /
// ReferenceID to tie it to an order. Scope partner-api.payments.write.
// Payments are never blocked by business rules; pass WithIdempotencyKey so a
// retry cannot post twice.
func (s *CustomersService) CreatePayment(ctx context.Context, id string, body PaymentCreateRequest, opts ...RequestOption) (*Payment, error) {
	return call[Payment](ctx, s.c, http.MethodPost, "/partner/v1/customers/"+seg(id)+"/payments", nil, body, opts...)
}

// Ledger returns the customer's balance summary. Scope partner-api.payments.read.
func (s *CustomersService) Ledger(ctx context.Context, id string) (*CustomerLedger, error) {
	return call[CustomerLedger](ctx, s.c, http.MethodGet, "/partner/v1/customers/"+seg(id)+"/ledger", nil, nil)
}
