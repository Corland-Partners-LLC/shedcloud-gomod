package partnerapi

import (
	"context"
	"net/http"
	"time"
)

// RTOInstallment is one row of an RTO agreement's payment schedule. Money is
// a decimal string.
type RTOInstallment struct {
	ID             string     `json:"id"`
	TenantID       string     `json:"tenant_id"`
	RTOAgreementID string     `json:"rto_agreement_id"`
	SequenceNo     int        `json:"sequence_no"`
	LineKind       string     `json:"line_kind"`
	DueDate        string     `json:"due_date"`
	Amount         string     `json:"amount"`
	ChargeID       string     `json:"charge_id,omitempty"`
	ChargeNumber   string     `json:"charge_number,omitempty"`
	Status         string     `json:"status"`
	PaidAmount     string     `json:"paid_amount,omitempty"`
	Balance        string     `json:"balance,omitempty"`
	Trigger        string     `json:"trigger,omitempty"`
	TriggeredAt    *time.Time `json:"triggered_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// RTOAgreement is a customer rent-to-own contract with its payment schedule
// (/partner/v1/rto-agreements). Status is draft, active, completed,
// cancelled, defaulted or returned. Money and rates are decimal strings.
type RTOAgreement struct {
	ID                 string     `json:"id"`
	TenantID           string     `json:"tenant_id"`
	CustomerID         string     `json:"customer_id"`
	CustomerName       string     `json:"customer_name,omitempty"`
	QuoteID            string     `json:"quote_id,omitempty"`
	QuoteNumber        string     `json:"quote_number,omitempty"`
	OrderID            string     `json:"order_id,omitempty"`
	OrderNumber        string     `json:"order_number,omitempty"`
	ProgramAgreementID string     `json:"program_agreement_id,omitempty"`
	AgreementNumber    string     `json:"agreement_number"`
	Status             string     `json:"status"`
	StartDate          string     `json:"start_date,omitempty"`
	FirstPaymentDate   string     `json:"first_payment_date,omitempty"`
	TermMonths         int        `json:"term_months"`
	MonthlyPayment     string     `json:"monthly_payment"`
	SecurityDeposit    string     `json:"security_deposit"`
	AdditionalDeposit  string     `json:"additional_deposit"`
	CustomerReserve    string     `json:"customer_reserve"`
	DamageWaiver       bool       `json:"damage_waiver"`
	DamageWaiverFee    string     `json:"damage_waiver_fee"`
	Rent               string     `json:"rent"`
	RateFactor         string     `json:"rate_factor"`
	ProviderName       string     `json:"provider_name,omitempty"`
	StateCode          string     `json:"state_code,omitempty"`
	Notes              string     `json:"notes,omitempty"`
	ActivatedAt        *time.Time `json:"activated_at,omitempty"`
	CreatedBy          string     `json:"created_by,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`

	Origin           string         `json:"origin,omitempty"`
	DealerID         string         `json:"dealer_id,omitempty"`
	DealerName       string         `json:"dealer_name,omitempty"`
	CapitalProductID string         `json:"capital_product_id,omitempty"`
	Product          map[string]any `json:"product,omitempty"`
	CashPrice        string         `json:"cash_price,omitempty"`
	TaxRate          string         `json:"tax_rate,omitempty"`
	TaxAmount        string         `json:"tax_amount,omitempty"`
	OwnershipTotal   string         `json:"ownership_total,omitempty"`
	DeliveryAddress  map[string]any `json:"delivery_address,omitempty"`
	ProgramTitle     string         `json:"program_title,omitempty"`

	SecurityDepositPct            string     `json:"security_deposit_pct,omitempty"`
	SecurityDepositAmountOverride string     `json:"security_deposit_amount_override,omitempty"`
	CustomerReserveMonths         int        `json:"customer_reserve_months"`
	DamageWaiverPct               string     `json:"damage_waiver_pct,omitempty"`
	DamageWaiverFlat              string     `json:"damage_waiver_flat,omitempty"`
	UseDamageWaiverPct            bool       `json:"use_damage_waiver_pct"`
	ReturnedAt                    *time.Time `json:"returned_at,omitempty"`

	// ContractKind is rto or cash; ScheduleKind is installments or milestones.
	ContractKind        string     `json:"contract_kind"`
	ScheduleKind        string     `json:"schedule_kind"`
	LateFeePolicyID     string     `json:"late_fee_policy_id,omitempty"`
	EpoPolicyID         string     `json:"epo_policy_id,omitempty"`
	UnitID              string     `json:"unit_id,omitempty"`
	ClosedAt            *time.Time `json:"closed_at,omitempty"`
	ClosedReason        string     `json:"closed_reason,omitempty"`
	AssignedCollectorID string     `json:"assigned_collector_id,omitempty"`
	LastContactAt       *time.Time `json:"last_contact_at,omitempty"`

	Installments []RTOInstallment `json:"installments,omitempty"`

	TotalScheduled string `json:"total_scheduled,omitempty"`
	TotalPaid      string `json:"total_paid,omitempty"`
	BalanceDue     string `json:"balance_due,omitempty"`
}

// RTOAgreementListParams filters GET /partner/v1/rto-agreements.
type RTOAgreementListParams struct {
	ListParams
	CustomerID string `json:"customer_id,omitempty"`
	QuoteID    string `json:"quote_id,omitempty"`
	Status     string `json:"status,omitempty"`
}

// RTOAgreementCreateRequest is POST /partner/v1/rto-agreements. Dates are
// YYYY-MM-DD; money and RateFactor are decimal strings. Unknown fields are
// rejected by the ERP (400).
type RTOAgreementCreateRequest struct {
	CustomerID         string `json:"customer_id,omitempty"`
	QuoteID            string `json:"quote_id,omitempty"`
	ProgramAgreementID string `json:"program_agreement_id,omitempty"`
	StartDate          string `json:"start_date,omitempty"`
	FirstPaymentDate   string `json:"first_payment_date,omitempty"`
	TermMonths         int    `json:"term_months,omitempty"`
	MonthlyPayment     string `json:"monthly_payment,omitempty"`
	SecurityDeposit    string `json:"security_deposit,omitempty"`
	AdditionalDeposit  string `json:"additional_deposit,omitempty"`
	CustomerReserve    string `json:"customer_reserve,omitempty"`
	DamageWaiver       bool   `json:"damage_waiver,omitempty"`
	DamageWaiverFee    string `json:"damage_waiver_fee,omitempty"`
	Rent               string `json:"rent,omitempty"`
	RateFactor         string `json:"rate_factor,omitempty"`
	ProviderName       string `json:"provider_name,omitempty"`
	StateCode          string `json:"state_code,omitempty"`
	Notes              string `json:"notes,omitempty"`
}

// RTOAgreementsService calls /partner/v1/rto-agreements.
type RTOAgreementsService struct{ c *Client }

// List returns customer RTO agreements. Scope partner-api.rto-agreements.read.
func (s *RTOAgreementsService) List(ctx context.Context, params RTOAgreementListParams) (*Page[RTOAgreement], error) {
	return listPage[RTOAgreement](ctx, s.c, "/partner/v1/rto-agreements", "agreements", params)
}

// Create makes a draft RTO agreement. Scope partner-api.rto-agreements.write.
// Pass WithIdempotencyKey to make retries safe.
func (s *RTOAgreementsService) Create(ctx context.Context, body RTOAgreementCreateRequest, opts ...RequestOption) (*RTOAgreement, error) {
	return call[RTOAgreement](ctx, s.c, http.MethodPost, "/partner/v1/rto-agreements", nil, body, opts...)
}

// Get returns one agreement with its installments. Scope
// partner-api.rto-agreements.read.
func (s *RTOAgreementsService) Get(ctx context.Context, id string) (*RTOAgreement, error) {
	return call[RTOAgreement](ctx, s.c, http.MethodGet, "/partner/v1/rto-agreements/"+seg(id), nil, nil)
}

// Activate activates a draft agreement. Scope partner-api.rto-agreements.write.
func (s *RTOAgreementsService) Activate(ctx context.Context, id string) (*RTOAgreement, error) {
	return call[RTOAgreement](ctx, s.c, http.MethodPost, "/partner/v1/rto-agreements/"+seg(id)+"/activate", nil, nil)
}
