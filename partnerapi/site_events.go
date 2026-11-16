package partnerapi

import (
	"context"
	"net/http"
)

// SiteEventCustomer enriches identity events (customer.profile,
// identity.resolved) so the visitor profile gains a name/email/phone instead
// of staying anonymous. snake_case.
type SiteEventCustomer struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Email     string `json:"email,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Zip       string `json:"zip,omitempty"`
}

// SiteEventInput is one event of a SiteEvents.Track batch (snake_case).
//
// v1.0.0 note: the ERP ingest has no Delivery or Payment block (0.x had
// them); put that detail in Payload.
type SiteEventInput struct {
	EventType string `json:"event_type"`
	// EventID is an optional client id for de-duplication.
	EventID string `json:"event_id,omitempty"`
	// OccurredAt is RFC3339. Defaults to receipt time when omitted.
	OccurredAt  string `json:"occurred_at,omitempty"`
	VisitorID   string `json:"visitor_id,omitempty"`
	UserID      string `json:"user_id,omitempty"`
	LocationID  string `json:"location_id,omitempty"`
	QuoteID     string `json:"quote_id,omitempty"`
	ConfigID    string `json:"config_id,omitempty"`
	ProductID   string `json:"product_id,omitempty"`
	ProductName string `json:"product_name,omitempty"`
	Step        string `json:"step,omitempty"`
	Page        string `json:"page,omitempty"`
	MenuItem    string `json:"menu_item,omitempty"`
	// Customer enriches identity events.
	Customer *SiteEventCustomer `json:"customer,omitempty"`
	// Payload is arbitrary JSON detail (keep it small — the batch is size-limited).
	Payload any `json:"payload,omitempty"`
}

// SiteEventsTrackRequest is the batched envelope for SiteEvents.Track
// (snake_case, as in 0.x).
type SiteEventsTrackRequest struct {
	// SessionID identifies one visit (rotate per browser session).
	SessionID string `json:"session_id"`
	// VisitorID is the permanent visitor identity (localStorage UUID).
	VisitorID string `json:"visitor_id,omitempty"`
	// SiteHost is the property emitting events, e.g. "example-sheds.com".
	SiteHost string `json:"site_host,omitempty"`
	// Timezone is the browser's IANA timezone; Language is navigator.language.
	Timezone      string           `json:"timezone,omitempty"`
	Language      string           `json:"language,omitempty"`
	CorrelationID string           `json:"correlation_id,omitempty"`
	Events        []SiteEventInput `json:"events"`
	// ClientIP is the end shopper's IP as seen by your backend proxy. Without
	// it every event carries your server's IP.
	ClientIP string `json:"client_ip,omitempty"`
	// ClientUserAgent is the end shopper's browser User-Agent.
	ClientUserAgent string `json:"client_user_agent,omitempty"`
	// Referrer is the visit's external document.referrer.
	Referrer string `json:"referrer,omitempty"`
	// UTM carries utm_* query params captured client-side.
	UTM map[string]string `json:"utm,omitempty"`
}

// SiteEventsTrackResponse reports how many events were accepted (HTTP 202).
type SiteEventsTrackResponse struct {
	Success  bool `json:"success"`
	Accepted int  `json:"accepted"`
}

// SiteEventItem is one tracked event of GET /partner/v1/site-events. The read
// side stays camelCase (it did in 0.x too).
type SiteEventItem struct {
	EventID    string `json:"eventId"`
	EventType  string `json:"eventType"`
	OccurredAt string `json:"occurredAt"`
	SessionID  string `json:"sessionId"`
	VisitorID  string `json:"visitorId,omitempty"`
	// Source is "marketing" or "configurator".
	Source     string `json:"source,omitempty"`
	SiteHost   string `json:"siteHost,omitempty"`
	ProductID  string `json:"productId,omitempty"`
	QuoteID    string `json:"quoteId,omitempty"`
	ConfigID   string `json:"configId,omitempty"`
	GeoCountry string `json:"geoCountry,omitempty"`
	GeoRegion  string `json:"geoRegion,omitempty"`
	GeoCity    string `json:"geoCity,omitempty"`
	Payload    any    `json:"payload,omitempty"`
}

// SiteEventListParams are the query params of GET /partner/v1/site-events.
// The ERP has no cursor: it returns one batch of up to Limit events.
type SiteEventListParams struct {
	// Limit is the max events returned (ERP default 100).
	Limit int `json:"limit,omitempty"`
	// SessionID restricts to one session.
	SessionID string `json:"sessionId,omitempty"`
	// Types filters to these event types (joined with commas for you).
	Types []string `json:"types,omitempty"`
	// From / To are RFC3339 bounds on occurredAt.
	From string `json:"from,omitempty"`
	To   string `json:"to,omitempty"`
}

// SiteEventListResponse is the GET /partner/v1/site-events envelope.
type SiteEventListResponse struct {
	Data []SiteEventItem `json:"data"`
	// HasMore is always false in the ERP today (no pagination).
	HasMore bool `json:"hasMore"`
}

// SiteEventsService calls /partner/v1/site-events — visitor behavioral
// tracking for partner marketing sites. Track needs
// partner-api.site-events.write, List partner-api.site-events.read.
//
// Partner keys are server-side secrets: proxy events through your backend,
// never call this from the browser.
type SiteEventsService struct{ c *Client }

// Track ingests a batch of behavioral events.
func (s *SiteEventsService) Track(ctx context.Context, req SiteEventsTrackRequest) (*SiteEventsTrackResponse, error) {
	return call[SiteEventsTrackResponse](ctx, s.c, http.MethodPost, "/partner/v1/site-events", nil, req)
}

// List returns tracked events.
func (s *SiteEventsService) List(ctx context.Context, params SiteEventListParams) (*SiteEventListResponse, error) {
	return call[SiteEventListResponse](ctx, s.c, http.MethodGet, "/partner/v1/site-events", params, nil)
}
