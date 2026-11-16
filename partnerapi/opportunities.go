package partnerapi

import (
	"context"
	"net/http"
	"time"
)

// Opportunity is a sales opportunity (/partner/v1/opportunities). Stage is one
// of qualify, proposal, negotiation, won, lost. Amount is a decimal string.
type Opportunity struct {
	ID              string    `json:"id"`
	TenantID        string    `json:"tenant_id"`
	Name            string    `json:"name"`
	Stage           string    `json:"stage"`
	Amount          string    `json:"amount"`
	CloseDate       *string   `json:"close_date,omitempty"`
	LeadID          *string   `json:"lead_id,omitempty"`
	QuoteID         *string   `json:"quote_id,omitempty"`
	StoreLocationID *string   `json:"store_location_id,omitempty"`
	OwnerID         *string   `json:"owner_id,omitempty"`
	CreatedBy       *string   `json:"created_by,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// OpportunityListParams filters GET /partner/v1/opportunities.
type OpportunityListParams struct {
	ListParams
	Status string `json:"status,omitempty"`
}

// OpportunityCreateRequest is POST /partner/v1/opportunities. Amount is a
// decimal string; CloseDate is YYYY-MM-DD.
type OpportunityCreateRequest struct {
	Name            string `json:"name,omitempty"`
	Stage           string `json:"stage,omitempty"`
	Amount          string `json:"amount,omitempty"`
	CloseDate       string `json:"close_date,omitempty"`
	LeadID          string `json:"lead_id,omitempty"`
	QuoteID         string `json:"quote_id,omitempty"`
	StoreLocationID string `json:"store_location_id,omitempty"`
	OwnerID         string `json:"owner_id,omitempty"`
}

// OpportunitiesService calls /partner/v1/opportunities.
type OpportunitiesService struct{ c *Client }

// List returns opportunities. Scope partner-api.opportunities.read.
func (s *OpportunitiesService) List(ctx context.Context, params OpportunityListParams) (*Page[Opportunity], error) {
	return listPage[Opportunity](ctx, s.c, "/partner/v1/opportunities", "opportunities", params)
}

// Get returns one opportunity. Scope partner-api.opportunities.read.
func (s *OpportunitiesService) Get(ctx context.Context, id string) (*Opportunity, error) {
	return call[Opportunity](ctx, s.c, http.MethodGet, "/partner/v1/opportunities/"+seg(id), nil, nil)
}

// Create makes an opportunity. Scope partner-api.opportunities.write. Pass
// WithIdempotencyKey to make retries safe.
func (s *OpportunitiesService) Create(ctx context.Context, body OpportunityCreateRequest, opts ...RequestOption) (*Opportunity, error) {
	return call[Opportunity](ctx, s.c, http.MethodPost, "/partner/v1/opportunities", nil, body, opts...)
}
