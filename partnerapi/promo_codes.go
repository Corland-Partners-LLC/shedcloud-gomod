package partnerapi

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

// PromoCode is a tenant promo code (/partner/v1/promo-codes). Rules is the
// versioned rule document, kept as raw JSON.
type PromoCode struct {
	ID               string          `json:"id"`
	TenantID         string          `json:"tenant_id,omitempty"`
	Code             string          `json:"code"`
	DescriptionHTML  string          `json:"description_html,omitempty"`
	Status           string          `json:"status"`
	StartsAt         *time.Time      `json:"starts_at,omitempty"`
	EndsAt           *time.Time      `json:"ends_at,omitempty"`
	Rules            json.RawMessage `json:"rules"`
	MaxRedemptions   *int            `json:"max_redemptions,omitempty"`
	PerCustomerLimit *int            `json:"per_customer_limit,omitempty"`
	RedemptionCount  int             `json:"redemption_count"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
}

// PromoCodeListParams filters GET /partner/v1/promo-codes. The list is not
// paginated.
type PromoCodeListParams struct {
	// Status is draft | scheduled | active | paused | expired | archived.
	Status string `json:"status,omitempty"`
	Q      string `json:"q,omitempty"`
}

// PromoCodeCreateRequest is POST /partner/v1/promo-codes. StartsAt and EndsAt
// are RFC3339. Caps above zero are applied; zero means no cap.
type PromoCodeCreateRequest struct {
	Code             string          `json:"code,omitempty"`
	DescriptionHTML  string          `json:"description_html,omitempty"`
	Status           string          `json:"status,omitempty"`
	Rules            json.RawMessage `json:"rules,omitempty"`
	MaxRedemptions   int             `json:"max_redemptions,omitempty"`
	PerCustomerLimit int             `json:"per_customer_limit,omitempty"`
	StartsAt         string          `json:"starts_at,omitempty"`
	EndsAt           string          `json:"ends_at,omitempty"`
}

// PromoCodeUpdateRequest is PATCH /partner/v1/promo-codes/{promoId}. Only
// non-nil fields are sent. An empty StartsAt/EndsAt clears the schedule bound
// and a MaxRedemptions/PerCustomerLimit of 0 clears the cap.
type PromoCodeUpdateRequest struct {
	Code             *string         `json:"code,omitempty"`
	DescriptionHTML  *string         `json:"description_html,omitempty"`
	Status           *string         `json:"status,omitempty"`
	StartsAt         *string         `json:"starts_at,omitempty"`
	EndsAt           *string         `json:"ends_at,omitempty"`
	Rules            json.RawMessage `json:"rules,omitempty"`
	MaxRedemptions   *int            `json:"max_redemptions,omitempty"`
	PerCustomerLimit *int            `json:"per_customer_limit,omitempty"`
}

// PromoCodeValidateRequest is POST /partner/v1/promo-codes/validate: the code
// plus the quote context it is evaluated against. Subtotal is a decimal
// string.
type PromoCodeValidateRequest struct {
	Code        string   `json:"code,omitempty"`
	QuoteID     string   `json:"quote_id,omitempty"`
	CustomerID  string   `json:"customer_id,omitempty"`
	Subtotal    string   `json:"subtotal,omitempty"`
	DealType    string   `json:"deal_type,omitempty"`
	TermMonths  int      `json:"term_months,omitempty"`
	StateCode   string   `json:"state_code,omitempty"`
	LineSources []string `json:"line_sources,omitempty"`
}

// PromoCodeAdjustment is one evaluated change a promo makes to deal terms.
// Operator is "set" or "subtract"; Value is a decimal string.
type PromoCodeAdjustment struct {
	Field     string `json:"field"`
	Operator  string `json:"operator"`
	ValueUnit string `json:"value_unit,omitempty"`
	Value     string `json:"value"`
}

// PromoCodeEffects is the outcome of evaluating a promo against a quote.
type PromoCodeEffects struct {
	Matched     bool                  `json:"matched"`
	Adjustments []PromoCodeAdjustment `json:"adjustments,omitempty"`
}

// PromoCodeValidation is the validate reply. When Valid is false, Reason says
// why and Promo/Effects are absent.
type PromoCodeValidation struct {
	Valid   bool              `json:"valid"`
	Reason  string            `json:"reason,omitempty"`
	Promo   *PromoCode        `json:"promo,omitempty"`
	Effects *PromoCodeEffects `json:"effects,omitempty"`
}

// PromoCodesService calls /partner/v1/promo-codes.
type PromoCodesService struct{ c *Client }

// List returns promo codes. Scope partner-api.promo-codes.read.
func (s *PromoCodesService) List(ctx context.Context, params PromoCodeListParams) ([]PromoCode, error) {
	return listKey[PromoCode](ctx, s.c, http.MethodGet, "/partner/v1/promo-codes", "promo_codes", params)
}

// Create makes a promo code. Scope partner-api.promo-codes.write. Pass
// WithIdempotencyKey to make retries safe.
func (s *PromoCodesService) Create(ctx context.Context, body PromoCodeCreateRequest, opts ...RequestOption) (*PromoCode, error) {
	return call[PromoCode](ctx, s.c, http.MethodPost, "/partner/v1/promo-codes", nil, body, opts...)
}

// Validate checks a code against a quote context without persisting anything.
// Scope partner-api.promo-codes.read.
func (s *PromoCodesService) Validate(ctx context.Context, body PromoCodeValidateRequest) (*PromoCodeValidation, error) {
	return call[PromoCodeValidation](ctx, s.c, http.MethodPost, "/partner/v1/promo-codes/validate", nil, body)
}

// Get returns one promo code. Scope partner-api.promo-codes.read.
func (s *PromoCodesService) Get(ctx context.Context, id string) (*PromoCode, error) {
	return call[PromoCode](ctx, s.c, http.MethodGet, "/partner/v1/promo-codes/"+seg(id), nil, nil)
}

// Update patches a promo code. Scope partner-api.promo-codes.write.
func (s *PromoCodesService) Update(ctx context.Context, id string, body PromoCodeUpdateRequest) (*PromoCode, error) {
	return call[PromoCode](ctx, s.c, http.MethodPatch, "/partner/v1/promo-codes/"+seg(id), nil, body)
}
