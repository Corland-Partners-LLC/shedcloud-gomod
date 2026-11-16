package partnerapi

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

// Payment is a customer payment (/partner/v1/payments). Money is a decimal
// string; PaymentDate is "YYYY-MM-DD".
type Payment struct {
	ID              string          `json:"id"`
	TenantID        string          `json:"tenant_id,omitempty"`
	CustomerID      string          `json:"customer_id"`
	CustomerName    string          `json:"customer_name,omitempty"`
	PaymentNumber   string          `json:"payment_number"`
	PaymentDate     string          `json:"payment_date"`
	Amount          string          `json:"amount"`
	AllocatedAmount string          `json:"allocated_amount"`
	UnappliedAmount string          `json:"unapplied_amount"`
	RefundedAmount  string          `json:"refunded_amount"`
	Method          string          `json:"method"`
	Reference       string          `json:"reference,omitempty"`
	Status          string          `json:"status"`
	Notes           string          `json:"notes,omitempty"`
	ExternalRef     string          `json:"external_ref,omitempty"`
	Gateway         string          `json:"gateway,omitempty"`
	CardBrand       string          `json:"card_brand,omitempty"`
	CardLast4       string          `json:"card_last4,omitempty"`
	CardWallet      string          `json:"card_wallet,omitempty"`
	ReferenceType   string          `json:"reference_type,omitempty"`
	ReferenceID     string          `json:"reference_id,omitempty"`
	JournalEntryID  string          `json:"journal_entry_id,omitempty"`
	Allocations     json.RawMessage `json:"allocations,omitempty"`
	CreatedBy       string          `json:"created_by,omitempty"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

// PaymentListParams filters GET /partner/v1/payments. To list the payments of
// one order use ReferenceType "order" and ReferenceID = the order id.
type PaymentListParams struct {
	ListParams
	CustomerID    string `json:"customer_id,omitempty"`
	ReferenceType string `json:"reference_type,omitempty"`
	ReferenceID   string `json:"reference_id,omitempty"`
}

// PaymentCreateRequest is POST /partner/v1/customers/{id}/payments. Amount is
// a decimal string.
type PaymentCreateRequest struct {
	PaymentNumber string `json:"payment_number,omitempty"`
	PaymentDate   string `json:"payment_date,omitempty"`
	Amount        string `json:"amount"`
	Method        string `json:"method,omitempty"`
	Reference     string `json:"reference,omitempty"`
	Notes         string `json:"notes,omitempty"`
	ExternalRef   string `json:"external_ref,omitempty"`
	Gateway       string `json:"gateway,omitempty"`
	ReferenceType string `json:"reference_type,omitempty"`
	ReferenceID   string `json:"reference_id,omitempty"`
}

// PaymentsService calls /partner/v1/payments (read). Create payments with
// Customers.CreatePayment.
type PaymentsService struct{ c *Client }

// List returns payments. Scope partner-api.payments.read.
func (s *PaymentsService) List(ctx context.Context, params PaymentListParams) (*Page[Payment], error) {
	return listPage[Payment](ctx, s.c, "/partner/v1/payments", "payments", params)
}

// Get returns one payment. Scope partner-api.payments.read.
func (s *PaymentsService) Get(ctx context.Context, id string) (*Payment, error) {
	return call[Payment](ctx, s.c, http.MethodGet, "/partner/v1/payments/"+seg(id), nil, nil)
}
