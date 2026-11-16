package partnerapi

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

// Location is a company location: store, plant, yard (/partner/v1/locations).
type Location struct {
	ID              string          `json:"id"`
	TenantID        string          `json:"tenant_id,omitempty"`
	Code            string          `json:"code"`
	Name            string          `json:"name"`
	LocationType    string          `json:"location_type"`
	Phone           string          `json:"phone"`
	LogoURL         string          `json:"logo_url,omitempty"`
	IsActive        bool            `json:"is_active"`
	ConsiderLedger  bool            `json:"consider_ledger"`
	IsStockLocation bool            `json:"is_stock_location"`
	ContactName     string          `json:"contact_name,omitempty"`
	ContactEmail    string          `json:"contact_email,omitempty"`
	ContactPhone    string          `json:"contact_phone,omitempty"`
	Address         string          `json:"address,omitempty"`
	Latitude        string          `json:"latitude,omitempty"`
	Longitude       string          `json:"longitude,omitempty"`
	PlantID         string          `json:"plant_id,omitempty"`
	HostLocationID  string          `json:"host_location_id,omitempty"`
	IsVisible       bool            `json:"is_visible"`
	Slug            string          `json:"slug,omitempty"`
	LeadRouting     string          `json:"lead_routing,omitempty"`
	StoreHours      json.RawMessage `json:"store_hours,omitempty"`
	InventoryTerms  *int            `json:"inventory_terms,omitempty"`
	Notes           string          `json:"notes,omitempty"`
	PlantCode       string          `json:"plant_code,omitempty"`
	PlantName       string          `json:"plant_name,omitempty"`
	HostCode        string          `json:"host_code,omitempty"`
	HostName        string          `json:"host_name,omitempty"`
	RelatedCount    int             `json:"related_count"`
	Related         json.RawMessage `json:"related,omitempty"`
	Metadata        map[string]any  `json:"metadata,omitempty"`
	CreatedBy       string          `json:"created_by,omitempty"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
	// LotStockLimit / RentalReturnLimit are decimal strings (nil = no limit).
	LotStockLimit     *string `json:"lot_stock_limit"`
	RentalReturnLimit *string `json:"rental_return_limit"`
}

// LocationListParams filters GET /partner/v1/locations (name/code search via Q).
type LocationListParams struct {
	ListParams
}

// LocationCreateRequest is POST /partner/v1/locations.
type LocationCreateRequest struct {
	Code              string          `json:"code,omitempty"`
	Name              string          `json:"name,omitempty"`
	LocationType      string          `json:"location_type,omitempty"`
	Phone             string          `json:"phone,omitempty"`
	LogoURL           string          `json:"logo_url,omitempty"`
	IsActive          *bool           `json:"is_active,omitempty"`
	ConsiderLedger    *bool           `json:"consider_ledger,omitempty"`
	IsStockLocation   *bool           `json:"is_stock_location,omitempty"`
	ContactName       string          `json:"contact_name,omitempty"`
	ContactEmail      string          `json:"contact_email,omitempty"`
	ContactPhone      string          `json:"contact_phone,omitempty"`
	Address           string          `json:"address,omitempty"`
	Latitude          string          `json:"latitude,omitempty"`
	Longitude         string          `json:"longitude,omitempty"`
	PlantID           string          `json:"plant_id,omitempty"`
	HostLocationID    string          `json:"host_location_id,omitempty"`
	IsVisible         *bool           `json:"is_visible,omitempty"`
	Slug              string          `json:"slug,omitempty"`
	LeadRouting       string          `json:"lead_routing,omitempty"`
	StoreHours        json.RawMessage `json:"store_hours,omitempty"`
	InventoryTerms    *int            `json:"inventory_terms,omitempty"`
	Notes             string          `json:"notes,omitempty"`
	Metadata          map[string]any  `json:"metadata,omitempty"`
	LotStockLimit     string          `json:"lot_stock_limit,omitempty"`
	RentalReturnLimit string          `json:"rental_return_limit,omitempty"`
}

// LocationUpdateRequest is PATCH /partner/v1/locations/{id}. Only non-nil
// fields are sent.
type LocationUpdateRequest struct {
	Code              *string         `json:"code,omitempty"`
	Name              *string         `json:"name,omitempty"`
	LocationType      *string         `json:"location_type,omitempty"`
	Phone             *string         `json:"phone,omitempty"`
	LogoURL           *string         `json:"logo_url,omitempty"`
	IsActive          *bool           `json:"is_active,omitempty"`
	ConsiderLedger    *bool           `json:"consider_ledger,omitempty"`
	IsStockLocation   *bool           `json:"is_stock_location,omitempty"`
	ContactName       *string         `json:"contact_name,omitempty"`
	ContactEmail      *string         `json:"contact_email,omitempty"`
	ContactPhone      *string         `json:"contact_phone,omitempty"`
	Address           *string         `json:"address,omitempty"`
	Latitude          *string         `json:"latitude,omitempty"`
	Longitude         *string         `json:"longitude,omitempty"`
	PlantID           *string         `json:"plant_id,omitempty"`
	HostLocationID    *string         `json:"host_location_id,omitempty"`
	IsVisible         *bool           `json:"is_visible,omitempty"`
	Slug              *string         `json:"slug,omitempty"`
	LeadRouting       *string         `json:"lead_routing,omitempty"`
	StoreHours        json.RawMessage `json:"store_hours,omitempty"`
	InventoryTerms    *int            `json:"inventory_terms,omitempty"`
	Notes             *string         `json:"notes,omitempty"`
	LotStockLimit     *string         `json:"lot_stock_limit,omitempty"`
	RentalReturnLimit *string         `json:"rental_return_limit,omitempty"`
	Metadata          map[string]any  `json:"metadata,omitempty"`
}

// LocationsService calls /partner/v1/locations.
type LocationsService struct{ c *Client }

// List returns locations. Scope partner-api.locations.read.
func (s *LocationsService) List(ctx context.Context, params LocationListParams) (*Page[Location], error) {
	return listPage[Location](ctx, s.c, "/partner/v1/locations", "locations", params)
}

// Get returns one location. Scope partner-api.locations.read.
func (s *LocationsService) Get(ctx context.Context, id string) (*Location, error) {
	return call[Location](ctx, s.c, http.MethodGet, "/partner/v1/locations/"+seg(id), nil, nil)
}

// Create makes a location. Scope partner-api.locations.write.
func (s *LocationsService) Create(ctx context.Context, body LocationCreateRequest, opts ...RequestOption) (*Location, error) {
	return call[Location](ctx, s.c, http.MethodPost, "/partner/v1/locations", nil, body, opts...)
}

// Update patches a location. Scope partner-api.locations.write.
func (s *LocationsService) Update(ctx context.Context, id string, body LocationUpdateRequest) (*Location, error) {
	return call[Location](ctx, s.c, http.MethodPatch, "/partner/v1/locations/"+seg(id), nil, body)
}

// Domains lists the storefront domains attached to a location. Scope
// partner-api.domains.read.
func (s *LocationsService) Domains(ctx context.Context, id string) ([]Storefront, error) {
	return listKey[Storefront](ctx, s.c, http.MethodGet, "/partner/v1/locations/"+seg(id)+"/domains", "storefronts", nil)
}
