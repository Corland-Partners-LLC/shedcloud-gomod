package partnerapi

import (
	"context"
	"net/http"
	"time"
)

// Subscription is a recurring customer charge (/partner/v1/subscriptions).
// Amount is a decimal string; NextChargeOn is YYYY-MM-DD.
type Subscription struct {
	ID           string    `json:"id"`
	TenantID     string    `json:"tenant_id,omitempty"`
	CustomerID   string    `json:"customer_id"`
	Amount       string    `json:"amount"`
	Cadence      string    `json:"cadence"`
	NextChargeOn string    `json:"next_charge_on"`
	Status       string    `json:"status"`
	LastChargeID *string   `json:"last_charge_id,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// SubscriptionListParams filters GET /partner/v1/subscriptions.
type SubscriptionListParams struct {
	ListParams
	Status string `json:"status,omitempty"`
}

// SubscriptionRequest is the body of POST /partner/v1/subscriptions. Amount is
// a decimal string; NextChargeOn is YYYY-MM-DD.
type SubscriptionRequest struct {
	CustomerID   string `json:"customer_id,omitempty"`
	Amount       string `json:"amount,omitempty"`
	Cadence      string `json:"cadence,omitempty"`
	NextChargeOn string `json:"next_charge_on,omitempty"`
}

// SubscriptionsService calls /partner/v1/subscriptions.
type SubscriptionsService struct{ c *Client }

// List returns subscriptions (filter by Status). Scope
// partner-api.subscriptions.read.
func (s *SubscriptionsService) List(ctx context.Context, params SubscriptionListParams) (*Page[Subscription], error) {
	return listPage[Subscription](ctx, s.c, "/partner/v1/subscriptions", "subscriptions", params)
}

// Get returns one subscription. Scope partner-api.subscriptions.read.
func (s *SubscriptionsService) Get(ctx context.Context, subscriptionID string) (*Subscription, error) {
	return call[Subscription](ctx, s.c, http.MethodGet, "/partner/v1/subscriptions/"+seg(subscriptionID), nil, nil)
}

// Create makes a subscription. Scope partner-api.subscriptions.write.
func (s *SubscriptionsService) Create(ctx context.Context, body SubscriptionRequest, opts ...RequestOption) (*Subscription, error) {
	return call[Subscription](ctx, s.c, http.MethodPost, "/partner/v1/subscriptions", nil, body, opts...)
}
