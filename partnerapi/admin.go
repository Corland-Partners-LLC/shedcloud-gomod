package partnerapi

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

// ---- Users and roles ---------------------------------------------------------

// User is the identity block of a company member.
type User struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	AvatarURL   string `json:"avatar_url"`
	Active      bool   `json:"active"`
}

// Membership links a user to the company.
type Membership struct {
	ID        string     `json:"id"`
	TenantID  string     `json:"tenant_id"`
	UserID    string     `json:"user_id"`
	Status    string     `json:"status"`
	IsDefault bool       `json:"is_default"`
	InvitedAt *time.Time `json:"invited_at,omitempty"`
	JoinedAt  *time.Time `json:"joined_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// Role is a company role.
type Role struct {
	ID          string    `json:"id"`
	Slug        string    `json:"slug"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	IsSystem    bool      `json:"is_system"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// RoleDetail is a Role plus its permissions (GET /partner/v1/roles).
type RoleDetail struct {
	Role
	Permissions json.RawMessage `json:"permissions,omitempty"`
}

// Member is one company member (/partner/v1/users).
type Member struct {
	User                User       `json:"user"`
	Membership          Membership `json:"membership"`
	Roles               []Role     `json:"roles"`
	JoinedAt            *time.Time `json:"joined_at,omitempty"`
	Status              string     `json:"status"`
	AccessBlocked       bool       `json:"access_blocked"`
	HasPermissionExtras bool       `json:"has_permission_extras"`
}

// UserListParams filters GET /partner/v1/users.
type UserListParams struct {
	ListParams
	Status string `json:"status,omitempty"`
	Role   string `json:"role,omitempty"`
}

// UserInviteRequest is POST /partner/v1/users (invites by email).
type UserInviteRequest struct {
	Email string `json:"email"`
	// Role is the role slug.
	Role string `json:"role,omitempty"`
	// SendEmail nil keeps the ERP default.
	SendEmail *bool `json:"send_email,omitempty"`
}

// Invitation is the reply of POST /partner/v1/users.
type Invitation struct {
	ID        string    `json:"id"`
	TenantID  string    `json:"tenant_id"`
	Email     string    `json:"email"`
	RoleID    string    `json:"role_id"`
	RoleName  string    `json:"role_name,omitempty"`
	Status    string    `json:"status"`
	InvitedBy *string   `json:"invited_by,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// UserUpdateRequest is PATCH /partner/v1/users/{id}. Only non-nil fields are sent.
type UserUpdateRequest struct {
	FullName  *string `json:"full_name,omitempty"`
	AvatarURL *string `json:"avatar_url,omitempty"`
	Status    *string `json:"status,omitempty"`
	// Role is the role slug.
	Role *string `json:"role,omitempty"`
}

// UsersService calls /partner/v1/users and /partner/v1/roles.
type UsersService struct{ c *Client }

// List returns company members. Scope partner-api.users.read.
func (s *UsersService) List(ctx context.Context, params UserListParams) (*Page[Member], error) {
	return listPage[Member](ctx, s.c, "/partner/v1/users", "members", params)
}

// Get returns one member. Scope partner-api.users.read.
func (s *UsersService) Get(ctx context.Context, id string) (*Member, error) {
	return call[Member](ctx, s.c, http.MethodGet, "/partner/v1/users/"+seg(id), nil, nil)
}

// Invite invites a user by email. Scope partner-api.users.write.
func (s *UsersService) Invite(ctx context.Context, body UserInviteRequest, opts ...RequestOption) (*Invitation, error) {
	return call[Invitation](ctx, s.c, http.MethodPost, "/partner/v1/users", nil, body, opts...)
}

// Update patches a member. Scope partner-api.users.write.
func (s *UsersService) Update(ctx context.Context, id string, body UserUpdateRequest) (*Member, error) {
	return call[Member](ctx, s.c, http.MethodPatch, "/partner/v1/users/"+seg(id), nil, body)
}

// Roles lists the company roles. Scope partner-api.users.read.
func (s *UsersService) Roles(ctx context.Context) ([]RoleDetail, error) {
	return listKey[RoleDetail](ctx, s.c, http.MethodGet, "/partner/v1/roles", "roles", nil)
}

// ---- Domains -----------------------------------------------------------------

// Storefront is a storefront domain with the locations it serves.
type Storefront struct {
	ID                string          `json:"id"`
	TenantID          string          `json:"tenant_id,omitempty"`
	IntegrationID     string          `json:"integration_id,omitempty"`
	ApexDomain        string          `json:"apex_domain"`
	Subdomain         string          `json:"subdomain"`
	Draft             bool            `json:"draft"`
	CnameData         map[string]any  `json:"cname_data,omitempty"`
	Locations         json.RawMessage `json:"locations,omitempty"`
	SortOrder         int             `json:"sort_order"`
	HostingProvider   string          `json:"hosting_provider,omitempty"`
	HostingStatus     string          `json:"hosting_status,omitempty"`
	DNSVerifiedAt     *time.Time      `json:"dns_verified_at,omitempty"`
	HostingVerifiedAt *time.Time      `json:"hosting_verified_at,omitempty"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
}

// DomainsService calls /partner/v1/domains. Scope partner-api.domains.read.
type DomainsService struct{ c *Client }

// List returns the company storefront domains (not paginated).
func (s *DomainsService) List(ctx context.Context) ([]Storefront, error) {
	return listKey[Storefront](ctx, s.c, http.MethodGet, "/partner/v1/domains", "storefronts", nil)
}

// ---- Documents ---------------------------------------------------------------

// Document is a file attached to an entity (/partner/v1/documents).
type Document struct {
	ID          string         `json:"id"`
	LinkID      string         `json:"link_id"`
	EntityType  string         `json:"entity_type"`
	EntityID    string         `json:"entity_id"`
	Role        string         `json:"role"`
	Position    int            `json:"position"`
	IsPrimary   bool           `json:"is_primary"`
	Caption     string         `json:"caption,omitempty"`
	FileName    string         `json:"file_name"`
	MimeType    string         `json:"mime_type,omitempty"`
	SizeBytes   int64          `json:"size_bytes,omitempty"`
	UploadedBy  string         `json:"uploaded_by,omitempty"`
	Metadata    map[string]any `json:"metadata,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	SourceLabel string         `json:"source_label,omitempty"`
}

// DocumentFilesParams selects the entity whose files are listed. EntityType
// and EntityID are required.
type DocumentFilesParams struct {
	EntityType string `json:"entity_type"`
	EntityID   string `json:"entity_id"`
	Role       string `json:"role,omitempty"`
}

// PresignedURL is a short-lived download URL.
type PresignedURL struct {
	URL       string            `json:"url"`
	Method    string            `json:"method"`
	ExpiresAt time.Time         `json:"expires_at"`
	Headers   map[string]string `json:"headers,omitempty"`
}

// DocumentsService calls /partner/v1/documents. Scope partner-api.documents.read.
type DocumentsService struct{ c *Client }

// List returns the files of one entity (not paginated).
func (s *DocumentsService) List(ctx context.Context, params DocumentFilesParams) ([]Document, error) {
	return listKey[Document](ctx, s.c, http.MethodGet, "/partner/v1/documents", "documents", params)
}

// Download returns a presigned URL for a file (nodeId = Document.ID).
func (s *DocumentsService) Download(ctx context.Context, nodeID string) (*PresignedURL, error) {
	return call[PresignedURL](ctx, s.c, http.MethodGet, "/partner/v1/documents/"+seg(nodeID)+"/download", nil, nil)
}

// ---- Events and webhook deliveries -------------------------------------------

// Event is a partner event from the outbox (/partner/v1/events).
type Event struct {
	ID         string         `json:"id"`
	TenantID   string         `json:"tenant_id,omitempty"`
	EntityType string         `json:"entity_type"`
	EntityID   string         `json:"entity_id"`
	Action     string         `json:"action"`
	Payload    map[string]any `json:"payload,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
}

// EventListParams are the query params of GET /partner/v1/events.
type EventListParams struct {
	// After is an RFC3339 lower bound on the event time.
	After string `json:"after,omitempty"`
	Limit int    `json:"limit,omitempty"`
}

// WebhookDelivery is one webhook delivery attempt row.
type WebhookDelivery struct {
	ID             string     `json:"id"`
	SubscriptionID string     `json:"subscription_id"`
	EventID        string     `json:"event_id"`
	Status         string     `json:"status"`
	AttemptCount   int        `json:"attempt_count"`
	ResponseCode   *int       `json:"response_code,omitempty"`
	DeliveredAt    *time.Time `json:"delivered_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
}

// WebhookDeliveryListParams filters GET /partner/v1/webhook-deliveries.
type WebhookDeliveryListParams struct {
	// Status: pending | delivering | delivered | failed | dead.
	Status string `json:"status,omitempty"`
	// Limit defaults to 50, max 200.
	Limit int `json:"limit,omitempty"`
}

// WebhookDeliveryList is the reply of GET /partner/v1/webhook-deliveries.
// Catalog is the list of event keys the ERP can emit.
type WebhookDeliveryList struct {
	Deliveries []WebhookDelivery `json:"deliveries"`
	Catalog    json.RawMessage   `json:"catalog,omitempty"`
}

// EventRedeliverResponse is the reply of POST /partner/v1/events/{id}/redeliver.
type EventRedeliverResponse struct {
	EventID  string `json:"event_id"`
	Enqueued int    `json:"enqueued"`
	// EventKey is "entity_type.action".
	EventKey string `json:"event_key"`
}

// EventsService calls /partner/v1/events and /partner/v1/webhook-deliveries.
// Scope partner-api.events.read (also for Redeliver).
type EventsService struct{ c *Client }

// List returns events (not paginated; use After + Limit to walk).
func (s *EventsService) List(ctx context.Context, params EventListParams) ([]Event, error) {
	return listKey[Event](ctx, s.c, http.MethodGet, "/partner/v1/events", "events", params)
}

// Deliveries returns webhook delivery attempts.
func (s *EventsService) Deliveries(ctx context.Context, params WebhookDeliveryListParams) (*WebhookDeliveryList, error) {
	return call[WebhookDeliveryList](ctx, s.c, http.MethodGet, "/partner/v1/webhook-deliveries", params, nil)
}

// Redeliver re-enqueues an event for the active webhook subscriptions.
func (s *EventsService) Redeliver(ctx context.Context, eventID string) (*EventRedeliverResponse, error) {
	return call[EventRedeliverResponse](ctx, s.c, http.MethodPost, "/partner/v1/events/"+seg(eventID)+"/redeliver", nil, nil)
}

// ---- Configurator sessions ---------------------------------------------------

// ConfiguratorSessionCreateRequest is POST /partner/v1/configurator-sessions.
// ReturnURL must match the company's allowed launch return URLs.
type ConfiguratorSessionCreateRequest struct {
	ProductID  string `json:"product_id"`
	LocationID string `json:"location_id,omitempty"`
	ReturnURL  string `json:"return_url,omitempty"`
}

// ConfiguratorSession is the launch URL for the 3D configurator.
type ConfiguratorSession struct {
	URL       string `json:"url"`
	ProductID string `json:"product_id"`
	CompanyID string `json:"company_id"`
}

// ConfiguratorSessionsService calls /partner/v1/configurator-sessions. Scope
// partner-api.configurator-sessions.write.
type ConfiguratorSessionsService struct{ c *Client }

// Create returns a configurator launch URL for a product.
func (s *ConfiguratorSessionsService) Create(ctx context.Context, body ConfiguratorSessionCreateRequest, opts ...RequestOption) (*ConfiguratorSession, error) {
	return call[ConfiguratorSession](ctx, s.c, http.MethodPost, "/partner/v1/configurator-sessions", nil, body, opts...)
}
