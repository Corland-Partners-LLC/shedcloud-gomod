package partnerapi

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

// DocumentHeader holds the header fields quotes and orders share. Money is a
// decimal string; address / setup / snapshot blobs stay json.RawMessage.
type DocumentHeader struct {
	ID                   string          `json:"id"`
	TenantID             string          `json:"tenant_id,omitempty"`
	Number               string          `json:"number,omitempty"`
	Status               string          `json:"status"`
	PriorStatus          string          `json:"prior_status,omitempty"`
	CustomerID           string          `json:"customer_id,omitempty"`
	BillToName           string          `json:"bill_to_name,omitempty"`
	BillToEmail          string          `json:"bill_to_email,omitempty"`
	BillToPhone          string          `json:"bill_to_phone,omitempty"`
	StoreLocationID      string          `json:"store_location_id,omitempty"`
	DealType             string          `json:"deal_type"`
	FulfillmentMethod    string          `json:"fulfillment_method"`
	SourceLocationID     string          `json:"source_location_id,omitempty"`
	SourceSlug           string          `json:"source_slug,omitempty"`
	BillingAddress       json.RawMessage `json:"billing_address,omitempty"`
	DeliveryAddress      json.RawMessage `json:"delivery_address,omitempty"`
	DeliverySetup        json.RawMessage `json:"delivery_setup,omitempty"`
	Currency             string          `json:"currency"`
	BasePrice            string          `json:"base_price"`
	TotalUpgrades        string          `json:"total_upgrades"`
	Subtotal             string          `json:"subtotal"`
	DiscountAmount       string          `json:"discount_amount"`
	DiscountReason       string          `json:"discount_reason,omitempty"`
	PromoCode            string          `json:"promo_code,omitempty"`
	DeliveryFee          string          `json:"delivery_fee"`
	BuildingAnchoredFee  string          `json:"building_anchored_fee"`
	MaterialSurchargeFee string          `json:"material_surcharge_fee"`
	FuelSurchargeFee     string          `json:"fuel_surcharge_fee"`
	TotalFees            string          `json:"total_fees"`
	TotalPretax          string          `json:"total_pretax"`
	TaxRate              string          `json:"tax_rate"`
	TaxAmount            string          `json:"tax_amount"`
	TaxSource            string          `json:"tax_source,omitempty"`
	TaxNote              string          `json:"tax_note,omitempty"`
	TaxExempt            bool            `json:"tax_exempt"`
	Total                string          `json:"total"`
	NearestLocationID    string          `json:"nearest_location_id,omitempty"`
	MileageTotal         string          `json:"mileage_total"`
	MileageFree          string          `json:"mileage_free"`
	MileageExtra         string          `json:"mileage_extra"`
	MileageRate          string          `json:"mileage_rate"`
	MileageMethod        string          `json:"mileage_method,omitempty"`
	MileageCharge        string          `json:"mileage_charge"`
	ValidUntil           string          `json:"valid_until,omitempty"`
	ConvertedWorkOrderID string          `json:"converted_work_order_id,omitempty"`
	PricingSnapshot      json.RawMessage `json:"pricing_snapshot,omitempty"`
	ShareTokenCreatedAt  *time.Time      `json:"share_token_created_at,omitempty"`
	// HasShareToken says the customer link exists; the raw token is only ever
	// returned once, by Quotes.Create.
	HasShareToken            bool            `json:"has_share_token,omitempty"`
	CustomerDecision         string          `json:"customer_decision,omitempty"`
	CustomerDecisionAt       *time.Time      `json:"customer_decision_at,omitempty"`
	CustomerDecisionNote     string          `json:"customer_decision_note,omitempty"`
	SentAt                   *time.Time      `json:"sent_at,omitempty"`
	CreatedBy                string          `json:"created_by,omitempty"`
	AssignedToUserID         string          `json:"assigned_to_user_id,omitempty"`
	AssignedToName           string          `json:"assigned_to_name,omitempty"`
	CreatedAt                time.Time       `json:"created_at"`
	UpdatedAt                time.Time       `json:"updated_at"`
	Lines                    []Line          `json:"lines,omitempty"`
	CustomerName             string          `json:"customer_name,omitempty"`
	StoreLocationName        string          `json:"store_location_name,omitempty"`
	PaymentOffer             json.RawMessage `json:"payment_offer,omitempty"`
	CreatedByName            string          `json:"created_by_name,omitempty"`
	SourceName               string          `json:"source_name,omitempty"`
	StatusEnteredAt          *time.Time      `json:"status_entered_at,omitempty"`
	ConvertedWorkOrderNumber string          `json:"converted_work_order_number,omitempty"`
	SerialNumber             string          `json:"serial_number,omitempty"`
	ProductSKU               string          `json:"product_sku,omitempty"`
	ProductName              string          `json:"product_name,omitempty"`
	// PreviewLine is the main line (product, serial, configuration ids).
	PreviewLine *PreviewLine `json:"preview_line,omitempty"`
	CashTerms   *CashTerms   `json:"cash_terms,omitempty"`
	RTOTerms    *RTOTerms    `json:"rto_terms,omitempty"`
}

// Quote is a quote (/partner/v1/quotes).
type Quote struct {
	DocumentHeader
	LeadID                        string          `json:"lead_id,omitempty"`
	LeadName                      string          `json:"lead_name,omitempty"`
	LostReason                    string          `json:"lost_reason,omitempty"`
	ConvertedOrderID              string          `json:"converted_order_id,omitempty"`
	ConvertedOrderNumber          string          `json:"converted_order_number,omitempty"`
	SoldAt                        *time.Time      `json:"sold_at,omitempty"`
	SentRevision                  int             `json:"sent_revision,omitempty"`
	ChangedSinceSent              bool            `json:"changed_since_sent,omitempty"`
	PreviousSentAt                *time.Time      `json:"previous_sent_at,omitempty"`
	CustomerDecisionPaymentOption string          `json:"customer_decision_payment_option,omitempty"`
	RTOSummary                    json.RawMessage `json:"rto_summary,omitempty"`
	LocationNotices               []string        `json:"location_notices,omitempty"`
}

// Order is an order (/partner/v1/orders).
type Order struct {
	DocumentHeader
	SourceQuoteID               string     `json:"source_quote_id,omitempty"`
	SourceQuoteNumber           string     `json:"source_quote_number,omitempty"`
	LegacyID                    string     `json:"legacy_id,omitempty"`
	SignatureProcessStatus      string     `json:"signature_process_status,omitempty"`
	ContractLocked              bool       `json:"contract_locked,omitempty"`
	ContractPDFStatus           string     `json:"contract_pdf_status,omitempty"`
	ContractPDFError            string     `json:"contract_pdf_error,omitempty"`
	SignatureProcessStartedAt   *time.Time `json:"signature_process_started_at,omitempty"`
	SignatureProcessCompletedAt *time.Time `json:"signature_process_completed_at,omitempty"`
	HasCorenterShareToken       bool       `json:"has_corenter_share_token,omitempty"`
	SignatureProvider           string     `json:"signature_provider,omitempty"`
	ConvertedWorkOrderStatus    string     `json:"converted_work_order_status,omitempty"`
	PlacedAt                    time.Time  `json:"placed_at"`
	SoldAt                      *time.Time `json:"sold_at,omitempty"`
	// SaleDate is the calendar sale day (read-only in v1).
	SaleDate *time.Time `json:"sale_date,omitempty"`
	// LedgerBooked / LedgerCancelled mirror the sales ledger (sticky cancel).
	LedgerBooked                bool   `json:"ledger_booked"`
	LedgerAmount                string `json:"ledger_amount,omitempty"`
	LedgerCancelled             bool   `json:"ledger_cancelled"`
	LedgerBasePrice             string `json:"ledger_base_price,omitempty"`
	LedgerUpgradesAmount        string `json:"ledger_upgrades_amount,omitempty"`
	LedgerRemovedUpgradesAmount string `json:"ledger_removed_upgrades_amount,omitempty"`
}

// QuoteCreated is the POST /partner/v1/quotes reply: the quote plus the
// customer Share link. ShareToken is returned only here, once; later reads
// expose HasShareToken only.
type QuoteCreated struct {
	Quote
	// ViewURL is the full customer page ({PORTAL_BASE_URL}/quotes/view/{token}).
	ViewURL    string `json:"view_url"`
	ShareToken string `json:"share_token"`
}

// Line is one quote or order line (QuoteID or OrderID is set, by parent).
type Line struct {
	ID                      string          `json:"id"`
	TenantID                string          `json:"tenant_id,omitempty"`
	QuoteID                 string          `json:"quote_id,omitempty"`
	OrderID                 string          `json:"order_id,omitempty"`
	LineType                string          `json:"line_type"`
	LineSource              string          `json:"line_source"`
	ProductID               string          `json:"product_id,omitempty"`
	StockUnitID             string          `json:"stock_unit_id,omitempty"`
	ConfigurationID         string          `json:"configuration_id,omitempty"`
	ConfigurationRevisionID string          `json:"configuration_revision_id,omitempty"`
	Description             string          `json:"description,omitempty"`
	Quantity                string          `json:"quantity"`
	BasePrice               string          `json:"base_price"`
	UpgradesTotal           string          `json:"upgrades_total"`
	RemovedUpgradesTotal    string          `json:"removed_upgrades_total"`
	UpgradesSnapshot        json.RawMessage `json:"upgrades_snapshot,omitempty"`
	Discount                string          `json:"discount"`
	LineTotal               string          `json:"line_total"`
	Unavailable             bool            `json:"unavailable"`
	SortOrder               int             `json:"sort_order"`
	IsCustom                bool            `json:"is_custom"`
	WidthFt                 string          `json:"width_ft,omitempty"`
	LengthFt                string          `json:"length_ft,omitempty"`
	ConvertedWorkOrderID    string          `json:"converted_work_order_id,omitempty"`
	SourceQuoteLineID       string          `json:"source_quote_line_id,omitempty"`
	CreatedAt               time.Time       `json:"created_at"`
	UpdatedAt               time.Time       `json:"updated_at"`
	ProductSKU              string          `json:"product_sku,omitempty"`
	ProductName             string          `json:"product_name,omitempty"`
	SerialNumber            string          `json:"serial_number,omitempty"`
	ReservedSerialNumber    string          `json:"reserved_serial_number,omitempty"`
	ActiveQuoteCount        int             `json:"active_quote_count,omitempty"`
	// RulesStatus / RulesEvaluation carry the configuration-rules verdict.
	RulesStatus      string          `json:"rules_status,omitempty"`
	RulesEvaluation  json.RawMessage `json:"rules_evaluation,omitempty"`
	BuildSaved       bool            `json:"build_saved,omitempty"`
	InteriorFinishes json.RawMessage `json:"interior_finishes,omitempty"`
}

// PreviewLine is the main line summary of a quote or order.
type PreviewLine struct {
	ID                      string `json:"id,omitempty"`
	ProductID               string `json:"product_id,omitempty"`
	ProductName             string `json:"product_name,omitempty"`
	ProductSKU              string `json:"product_sku,omitempty"`
	Description             string `json:"description,omitempty"`
	ConfigurationID         string `json:"configuration_id,omitempty"`
	ConfigurationRevisionID string `json:"configuration_revision_id,omitempty"`
	StockUnitID             string `json:"stock_unit_id,omitempty"`
	SerialNumber            string `json:"serial_number,omitempty"`
}

// CashTerms are the cash terms of a quote or order (QuoteID or OrderID set).
type CashTerms struct {
	QuoteID           string          `json:"quote_id,omitempty"`
	OrderID           string          `json:"order_id,omitempty"`
	TenantID          string          `json:"tenant_id,omitempty"`
	CashPurchaseTotal string          `json:"cash_purchase_total"`
	DownPayment       string          `json:"down_payment"`
	DownPaymentSource string          `json:"down_payment_source,omitempty"`
	PromoCode         string          `json:"promo_code,omitempty"`
	PromoSnapshot     json.RawMessage `json:"promo_snapshot,omitempty"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
	TotalPaid         string          `json:"total_paid,omitempty"`
	BalanceDue        string          `json:"balance_due,omitempty"`
}

// RTOTerms are the rent-to-own terms of a quote or order.
type RTOTerms struct {
	QuoteID               string          `json:"quote_id,omitempty"`
	OrderID               string          `json:"order_id,omitempty"`
	TenantID              string          `json:"tenant_id,omitempty"`
	AgreementID           string          `json:"agreement_id,omitempty"`
	ProviderName          string          `json:"provider_name,omitempty"`
	TermMonths            int             `json:"term_months"`
	SecurityDepositPct    string          `json:"security_deposit_pct"`
	SecurityDepositAmount string          `json:"security_deposit_amount"`
	Rent                  string          `json:"rent"`
	MonthlyPayment        string          `json:"monthly_payment"`
	DamageWaiver          bool            `json:"damage_waiver"`
	DamageWaiverFee       string          `json:"damage_waiver_fee"`
	MinimumThreshold      string          `json:"minimum_threshold"`
	CustomerReserve       string          `json:"customer_reserve"`
	AdditionalDeposit     string          `json:"additional_deposit"`
	RTOBalanceDue         string          `json:"rto_balance_due"`
	RateFactor            string          `json:"rate_factor"`
	PromoCode             string          `json:"promo_code,omitempty"`
	PromoSnapshot         json.RawMessage `json:"promo_snapshot,omitempty"`
	CreatedAt             time.Time       `json:"created_at"`
	UpdatedAt             time.Time       `json:"updated_at"`
	TotalPaid             string          `json:"total_paid,omitempty"`
	BalanceDue            string          `json:"balance_due,omitempty"`
}

// TermsPreview is a cash/RTO preview reply (the calculator output). Its keys
// belong to the ERP pricing engine, so it is returned as a generic map.
type TermsPreview map[string]any

// BookRepricePreview is POST .../reprice/preview.
type BookRepricePreview struct {
	Applicable        bool   `json:"applicable"`
	Stale             bool   `json:"stale"`
	CanApply          bool   `json:"can_apply"`
	Reason            string `json:"reason"`
	LineID            string `json:"line_id,omitempty"`
	Currency          string `json:"currency"`
	DocumentBasePrice string `json:"document_base_price"`
	BookBasePrice     string `json:"book_base_price,omitempty"`
	DocumentTotal     string `json:"document_total"`
	BookTotal         string `json:"book_total,omitempty"`
	PriceListCode     string `json:"price_list_code,omitempty"`
	PricingSource     string `json:"pricing_source,omitempty"`
}

// DocumentEvent is one audit event of a quote or order.
type DocumentEvent struct {
	ID          string          `json:"id"`
	TenantID    string          `json:"tenant_id,omitempty"`
	QuoteID     string          `json:"quote_id,omitempty"`
	OrderID     string          `json:"order_id,omitempty"`
	ActorUserID string          `json:"actor_user_id,omitempty"`
	ActorName   string          `json:"actor_name,omitempty"`
	EventType   string          `json:"event_type"`
	Detail      json.RawMessage `json:"detail,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
}

// DocumentListParams filters quote and order lists.
type DocumentListParams struct {
	ListParams
	Status     string `json:"status,omitempty"`
	CustomerID string `json:"customer_id,omitempty"`
	StoreID    string `json:"store_id,omitempty"`
}

// QuoteListParams / OrderListParams are DocumentListParams.
type (
	QuoteListParams = DocumentListParams
	OrderListParams = DocumentListParams
)

// DocumentCreateRequest is POST /partner/v1/quotes and /orders. Unknown fields
// are rejected by the ERP (400). StoreLocationID is the selling location.
type DocumentCreateRequest struct {
	CustomerID        string          `json:"customer_id,omitempty"`
	BillToName        string          `json:"bill_to_name,omitempty"`
	BillToEmail       string          `json:"bill_to_email,omitempty"`
	BillToPhone       string          `json:"bill_to_phone,omitempty"`
	StoreLocationID   string          `json:"store_location_id,omitempty"`
	DealType          string          `json:"deal_type,omitempty"`
	FulfillmentMethod string          `json:"fulfillment_method,omitempty"`
	SourceLocationID  string          `json:"source_location_id,omitempty"`
	SourceSlug        string          `json:"source_slug,omitempty"`
	BillingAddress    json.RawMessage `json:"billing_address,omitempty"`
	DeliveryAddress   json.RawMessage `json:"delivery_address,omitempty"`
	Currency          string          `json:"currency,omitempty"`
	TaxExempt         bool            `json:"tax_exempt,omitempty"`
	ValidUntil        string          `json:"valid_until,omitempty"`
}

// QuoteCreateRequest / OrderCreateRequest are DocumentCreateRequest.
type (
	QuoteCreateRequest = DocumentCreateRequest
	OrderCreateRequest = DocumentCreateRequest
)

// DocumentUpdateRequest is PATCH /partner/v1/quotes/{id} and /orders/{id}.
// Only non-nil fields are sent. Money fields are decimal strings.
type DocumentUpdateRequest struct {
	CustomerID          *string         `json:"customer_id,omitempty"`
	BillToName          *string         `json:"bill_to_name,omitempty"`
	BillToEmail         *string         `json:"bill_to_email,omitempty"`
	BillToPhone         *string         `json:"bill_to_phone,omitempty"`
	StoreLocationID     *string         `json:"store_location_id,omitempty"`
	DealType            *string         `json:"deal_type,omitempty"`
	FulfillmentMethod   *string         `json:"fulfillment_method,omitempty"`
	SourceLocationID    *string         `json:"source_location_id,omitempty"`
	SourceSlug          *string         `json:"source_slug,omitempty"`
	BillingAddress      json.RawMessage `json:"billing_address,omitempty"`
	DeliveryAddress     json.RawMessage `json:"delivery_address,omitempty"`
	Currency            *string         `json:"currency,omitempty"`
	DiscountAmount      *string         `json:"discount_amount,omitempty"`
	DiscountReason      *string         `json:"discount_reason,omitempty"`
	PromoCode           *string         `json:"promo_code,omitempty"`
	DeliveryFee         *string         `json:"delivery_fee,omitempty"`
	BuildingAnchoredFee *string         `json:"building_anchored_fee,omitempty"`
	TaxRate             *string         `json:"tax_rate,omitempty"`
	TaxExempt           *bool           `json:"tax_exempt,omitempty"`
	ValidUntil          *string         `json:"valid_until,omitempty"`
}

// QuoteUpdateRequest / OrderUpdateRequest are DocumentUpdateRequest.
type (
	QuoteUpdateRequest = DocumentUpdateRequest
	OrderUpdateRequest = DocumentUpdateRequest
)

// LineInput is one line of PUT .../lines (v1 sells one item per document).
// Set StockUnitID for lot stock, or ProductID (+ configuration ids) for a
// configured building. Money fields are decimal strings.
type LineInput struct {
	LineType                string          `json:"line_type,omitempty"`
	LineSource              string          `json:"line_source,omitempty"`
	ProductID               string          `json:"product_id,omitempty"`
	StockUnitID             string          `json:"stock_unit_id,omitempty"`
	ConfigurationID         string          `json:"configuration_id,omitempty"`
	ConfigurationRevisionID string          `json:"configuration_revision_id,omitempty"`
	Description             string          `json:"description,omitempty"`
	Quantity                string          `json:"quantity,omitempty"`
	BasePrice               string          `json:"base_price,omitempty"`
	UpgradesTotal           string          `json:"upgrades_total,omitempty"`
	RemovedUpgradesTotal    string          `json:"removed_upgrades_total,omitempty"`
	UpgradesSnapshot        json.RawMessage `json:"upgrades_snapshot,omitempty"`
	Discount                string          `json:"discount,omitempty"`
	SortOrder               int             `json:"sort_order,omitempty"`
	WidthFt                 string          `json:"width_ft,omitempty"`
	LengthFt                string          `json:"length_ft,omitempty"`
	IsCustom                bool            `json:"is_custom,omitempty"`
	ReservedSerialNumber    string          `json:"reserved_serial_number,omitempty"`
}

// LineBuildRequest is PATCH .../lines/{lineId}/build.
type LineBuildRequest struct {
	ConfigurationID         string          `json:"configuration_id,omitempty"`
	ConfigurationRevisionID string          `json:"configuration_revision_id,omitempty"`
	UpgradesSnapshot        json.RawMessage `json:"upgrades_snapshot,omitempty"`
	// DeriveUpgrades asks the ERP to derive the upgrades from the configuration.
	DeriveUpgrades bool `json:"derive_upgrades,omitempty"`
}

// CustomerLinkRequest is PUT .../customer: find-or-create the customer by
// contact data and link it.
type CustomerLinkRequest struct {
	Name              string `json:"name,omitempty"`
	Email             string `json:"email,omitempty"`
	Phone             string `json:"phone,omitempty"`
	SourceSlug        string `json:"source_slug,omitempty"`
	PreferredLanguage string `json:"preferred_language,omitempty"`
}

// CashTermsRequest is PUT .../cash-terms.
type CashTermsRequest struct {
	CashPurchaseTotal string          `json:"cash_purchase_total,omitempty"`
	DownPayment       string          `json:"down_payment,omitempty"`
	PromoCode         string          `json:"promo_code,omitempty"`
	PromoSnapshot     json.RawMessage `json:"promo_snapshot,omitempty"`
}

// CashTermsComputeRequest is POST .../cash-terms/preview and /compute.
type CashTermsComputeRequest struct {
	DownPayment *string `json:"down_payment,omitempty"`
	PromoCode   string  `json:"promo_code,omitempty"`
}

// RTOTermsRequest is PUT .../rto-terms.
type RTOTermsRequest struct {
	AgreementID           string `json:"agreement_id,omitempty"`
	ProviderName          string `json:"provider_name,omitempty"`
	TermMonths            int    `json:"term_months,omitempty"`
	SecurityDepositPct    string `json:"security_deposit_pct,omitempty"`
	SecurityDepositAmount string `json:"security_deposit_amount,omitempty"`
	Rent                  string `json:"rent,omitempty"`
	MonthlyPayment        string `json:"monthly_payment,omitempty"`
	DamageWaiver          bool   `json:"damage_waiver,omitempty"`
	DamageWaiverFee       string `json:"damage_waiver_fee,omitempty"`
	MinimumThreshold      string `json:"minimum_threshold,omitempty"`
	CustomerReserve       string `json:"customer_reserve,omitempty"`
	AdditionalDeposit     string `json:"additional_deposit,omitempty"`
	RTOBalanceDue         string `json:"rto_balance_due,omitempty"`
	RateFactor            string `json:"rate_factor,omitempty"`
}

// RTOTermsComputeRequest is POST .../rto-terms/preview and /compute (overrides
// are optional).
type RTOTermsComputeRequest struct {
	AgreementID                     string `json:"agreement_id,omitempty"`
	StateCode                       string `json:"state_code,omitempty"`
	TermMonths                      int    `json:"term_months,omitempty"`
	SecurityDepositPctOverride      string `json:"security_deposit_pct_override,omitempty"`
	MinimumThresholdPctOverride     string `json:"minimum_threshold_pct_override,omitempty"`
	CustomerReserveMonthsOverride   string `json:"customer_reserve_months_override,omitempty"`
	AdditionalDepositFlatOverride   string `json:"additional_deposit_flat_override,omitempty"`
	MinimumThresholdAmountOverride  string `json:"minimum_threshold_amount_override,omitempty"`
	CustomerReserveAmountOverride   string `json:"customer_reserve_amount_override,omitempty"`
	AdditionalDepositAmountOverride string `json:"additional_deposit_amount_override,omitempty"`
	TaxRateOverride                 string `json:"tax_rate_override,omitempty"`
	DamageWaiverOverride            *bool  `json:"damage_waiver_override,omitempty"`
	PromoCode                       string `json:"promo_code,omitempty"`
}

// commercialService is the call set quotes and orders share (same sub-routes
// under /partner/v1/quotes/{id} and /partner/v1/orders/{id}). T is Quote or
// Order. Scopes: reads need partner-api.<quotes|orders>.read, writes
// partner-api.<quotes|orders>.write.
type commercialService[T any] struct {
	c    *Client
	base string
	key  string
}

func (s *commercialService[T]) path(id string, rest string) string {
	return s.base + "/" + seg(id) + rest
}

// List returns one page of documents.
func (s *commercialService[T]) List(ctx context.Context, params DocumentListParams) (*Page[T], error) {
	return listPage[T](ctx, s.c, s.base, s.key, params)
}

// Get returns one document with lines and terms.
func (s *commercialService[T]) Get(ctx context.Context, id string) (*T, error) {
	return call[T](ctx, s.c, http.MethodGet, s.path(id, ""), nil, nil)
}

// Update patches header fields (PATCH).
func (s *commercialService[T]) Update(ctx context.Context, id string, body DocumentUpdateRequest) (*T, error) {
	return call[T](ctx, s.c, http.MethodPatch, s.path(id, ""), nil, body)
}

// SetCustomer links (or creates) the customer from contact data.
func (s *commercialService[T]) SetCustomer(ctx context.Context, id string, body CustomerLinkRequest) (*T, error) {
	return call[T](ctx, s.c, http.MethodPut, s.path(id, "/customer"), nil, body)
}

// ReplaceLines replaces all lines (v1 allows one item). A rental return still
// "received" is rejected, and a line whose configuration fails the rules gets
// 422 configuration_rule_violation.
func (s *commercialService[T]) ReplaceLines(ctx context.Context, id string, lines []LineInput) (*T, error) {
	body := struct {
		Lines []LineInput `json:"lines"`
	}{Lines: lines}
	return call[T](ctx, s.c, http.MethodPut, s.path(id, "/lines"), nil, body)
}

// UpdateLineBuild saves the configuration/upgrades of one line.
func (s *commercialService[T]) UpdateLineBuild(ctx context.Context, id, lineID string, body LineBuildRequest) (*T, error) {
	return call[T](ctx, s.c, http.MethodPatch, s.path(id, "/lines/"+seg(lineID)+"/build"), nil, body)
}

// Events returns the audit events of the document (not paginated).
func (s *commercialService[T]) Events(ctx context.Context, id string) ([]DocumentEvent, error) {
	return listKey[DocumentEvent](ctx, s.c, http.MethodGet, s.path(id, "/events"), "events", nil)
}

// PreviewCashTerms calculates cash terms without saving (read scope).
func (s *commercialService[T]) PreviewCashTerms(ctx context.Context, id string, body CashTermsComputeRequest) (TermsPreview, error) {
	out, err := call[TermsPreview](ctx, s.c, http.MethodPost, s.path(id, "/cash-terms/preview"), nil, body)
	if err != nil {
		return nil, err
	}
	return *out, nil
}

// ComputeCashTerms calculates and saves cash terms.
func (s *commercialService[T]) ComputeCashTerms(ctx context.Context, id string, body CashTermsComputeRequest) (*T, error) {
	return call[T](ctx, s.c, http.MethodPost, s.path(id, "/cash-terms/compute"), nil, body)
}

// SetCashTerms saves explicit cash terms.
func (s *commercialService[T]) SetCashTerms(ctx context.Context, id string, body CashTermsRequest) (*T, error) {
	return call[T](ctx, s.c, http.MethodPut, s.path(id, "/cash-terms"), nil, body)
}

// PreviewRTOTerms calculates RTO terms without saving (read scope).
func (s *commercialService[T]) PreviewRTOTerms(ctx context.Context, id string, body RTOTermsComputeRequest) (TermsPreview, error) {
	out, err := call[TermsPreview](ctx, s.c, http.MethodPost, s.path(id, "/rto-terms/preview"), nil, body)
	if err != nil {
		return nil, err
	}
	return *out, nil
}

// ComputeRTOTerms calculates and saves RTO terms.
func (s *commercialService[T]) ComputeRTOTerms(ctx context.Context, id string, body RTOTermsComputeRequest) (*T, error) {
	return call[T](ctx, s.c, http.MethodPost, s.path(id, "/rto-terms/compute"), nil, body)
}

// SetRTOTerms saves explicit RTO terms.
func (s *commercialService[T]) SetRTOTerms(ctx context.Context, id string, body RTOTermsRequest) (*T, error) {
	return call[T](ctx, s.c, http.MethodPut, s.path(id, "/rto-terms"), nil, body)
}

// ComputeMileage recomputes the delivery mileage charge.
func (s *commercialService[T]) ComputeMileage(ctx context.Context, id string) (*T, error) {
	return call[T](ctx, s.c, http.MethodPost, s.path(id, "/mileage"), nil, nil)
}

// ResolveTax recomputes tax (TaxJar) for the delivery address.
func (s *commercialService[T]) ResolveTax(ctx context.Context, id string) (*T, error) {
	return call[T](ctx, s.c, http.MethodPost, s.path(id, "/tax/resolve"), nil, nil)
}

// PreviewReprice compares the document price with the current price book
// (read scope).
func (s *commercialService[T]) PreviewReprice(ctx context.Context, id string) (*BookRepricePreview, error) {
	return call[BookRepricePreview](ctx, s.c, http.MethodPost, s.path(id, "/reprice/preview"), nil, nil)
}

// Reprice re-applies the current price book to the document.
func (s *commercialService[T]) Reprice(ctx context.Context, id string) (*T, error) {
	return call[T](ctx, s.c, http.MethodPost, s.path(id, "/reprice"), nil, nil)
}

// ---- Quotes ----------------------------------------------------------------

// QuotesService calls /partner/v1/quotes. Shared calls (List, Get, Update,
// SetCustomer, ReplaceLines, terms, mileage, tax, reprice, Events) come from
// the embedded service.
type QuotesService struct {
	*commercialService[Quote]
}

func newQuotesService(c *Client) *QuotesService {
	return &QuotesService{&commercialService[Quote]{c: c, base: "/partner/v1/quotes", key: "quotes"}}
}

// Create makes a quote (scope partner-api.quotes.write). The reply carries the
// customer ViewURL and the one-time ShareToken. Use WithIdempotencyKey for
// retry-safe creates.
func (s *QuotesService) Create(ctx context.Context, body QuoteCreateRequest, opts ...RequestOption) (*QuoteCreated, error) {
	return call[QuoteCreated](ctx, s.c, http.MethodPost, "/partner/v1/quotes", nil, body, opts...)
}

// UpdateStatus moves the quote to another status (422 on an invalid
// transition). Reason is the lost/cancel reason.
func (s *QuotesService) UpdateStatus(ctx context.Context, id, status, reason string) (*Quote, error) {
	body := struct {
		Status string `json:"status"`
		Reason string `json:"reason,omitempty"`
	}{status, reason}
	return call[Quote](ctx, s.c, http.MethodPatch, s.path(id, "/status"), nil, body)
}

// Sell is Place to Order: it commits the stock and creates the commercial
// order (Quote.ConvertedOrderID). It replaces Convert from 0.x. Scope
// partner-api.orders.write. May answer 409 configuration_rules_review_required.
func (s *QuotesService) Sell(ctx context.Context, id string) (*Quote, error) {
	return call[Quote](ctx, s.c, http.MethodPost, s.path(id, "/sell"), nil, nil)
}

// ---- Orders ----------------------------------------------------------------

// OrdersService calls /partner/v1/orders. Shared calls come from the embedded
// service.
type OrdersService struct {
	*commercialService[Order]
}

func newOrdersService(c *Client) *OrdersService {
	return &OrdersService{&commercialService[Order]{c: c, base: "/partner/v1/orders", key: "orders"}}
}

// Create makes an order directly (scope partner-api.orders.write). The usual
// path is Quotes.Sell.
func (s *OrdersService) Create(ctx context.Context, body OrderCreateRequest, opts ...RequestOption) (*Order, error) {
	return call[Order](ctx, s.c, http.MethodPost, "/partner/v1/orders", nil, body, opts...)
}

// UpdateStatus moves the order to another status. Moving to "unprocessed"
// from open / on hold needs the signed contract in Files (400 otherwise).
func (s *OrdersService) UpdateStatus(ctx context.Context, id, status string) (*Order, error) {
	body := struct {
		Status string `json:"status"`
	}{status}
	return call[Order](ctx, s.c, http.MethodPatch, s.path(id, "/status"), nil, body)
}

// OrderContract is GET /partner/v1/orders/{id}/contract: signature and
// contract-PDF progress (the PDF bytes are staff-only).
type OrderContract struct {
	OrderID                     string     `json:"order_id"`
	Number                      string     `json:"number"`
	Status                      string     `json:"status"`
	SignatureProcessStatus      string     `json:"signature_process_status"`
	ContractPDFStatus           string     `json:"contract_pdf_status"`
	ContractPDFError            string     `json:"contract_pdf_error"`
	SignatureProcessStartedAt   *time.Time `json:"signature_process_started_at"`
	SignatureProcessCompletedAt *time.Time `json:"signature_process_completed_at"`
}

// Contract returns the contract/signature status. Scope
// partner-api.contracts.read.
func (s *OrdersService) Contract(ctx context.Context, id string) (*OrderContract, error) {
	return call[OrderContract](ctx, s.c, http.MethodGet, s.path(id, "/contract"), nil, nil)
}
