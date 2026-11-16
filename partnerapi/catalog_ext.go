package partnerapi

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

// CatalogSizeRange is the width/length bounds of a catalog template.
type CatalogSizeRange struct {
	Unit      string   `json:"unit,omitempty"`
	MinWidth  *float64 `json:"min_width,omitempty"`
	MaxWidth  *float64 `json:"max_width,omitempty"`
	MinLength *float64 `json:"min_length,omitempty"`
	MaxLength *float64 `json:"max_length,omitempty"`
}

// CatalogTemplate is a sellable product template of the global Shedcloud
// catalog (/partner/v1/catalog/templates). Prices are decimal strings;
// default config and dimensions stay raw JSON.
type CatalogTemplate struct {
	ID                        string            `json:"id"`
	Slug                      string            `json:"slug"`
	Name                      string            `json:"name"`
	Subtitle                  string            `json:"subtitle,omitempty"`
	Description               string            `json:"description,omitempty"`
	Family                    string            `json:"family"`
	ProductType               string            `json:"product_type,omitempty"`
	RuntimeKey                string            `json:"runtime_key,omitempty"`
	Status                    string            `json:"status"`
	Available                 bool              `json:"available"`
	DefaultConfig             json.RawMessage   `json:"default_config,omitempty"`
	DefaultDimensions         json.RawMessage   `json:"default_dimensions,omitempty"`
	QuantitiesArePlaceholders bool              `json:"quantities_are_placeholders"`
	SortOrder                 int               `json:"sort_order"`
	SizeRange                 *CatalogSizeRange `json:"size_range,omitempty"`
	PrimaryImageURL           string            `json:"primary_image_url,omitempty"`
	Added                     bool              `json:"added"`
	AdoptedProductID          string            `json:"adopted_product_id,omitempty"`
	Private                   bool              `json:"private"`
	AcquisitionPrice          string            `json:"acquisition_price"`
	Blurb                     string            `json:"blurb,omitempty"`
	Highlights                []string          `json:"highlights"`
	MarketplaceGroup          string            `json:"marketplace_group,omitempty"`
	RatingAvg                 float64           `json:"rating_avg"`
	RatingCount               int               `json:"rating_count"`
	MyRating                  *int              `json:"my_rating,omitempty"`
	OptionCount               int               `json:"option_count"`
	ComponentCount            int               `json:"component_count"`
	CreatedAt                 time.Time         `json:"created_at"`
	UpdatedAt                 time.Time         `json:"updated_at"`
}

// CatalogImage is one gallery image (URL is a short-lived presigned GET).
type CatalogImage struct {
	ID        string `json:"id"`
	URL       string `json:"url,omitempty"`
	FileName  string `json:"file_name,omitempty"`
	AltText   string `json:"alt_text,omitempty"`
	SortOrder int    `json:"sort_order"`
	IsPrimary bool   `json:"is_primary"`
}

// CatalogOption is one configurable option of a template. Values, choices and
// visibility rules stay raw JSON.
type CatalogOption struct {
	ID           string          `json:"id"`
	GroupKey     string          `json:"group_key"`
	Key          string          `json:"key"`
	Label        string          `json:"label"`
	ValueType    string          `json:"value_type"`
	Unit         string          `json:"unit,omitempty"`
	DefaultValue json.RawMessage `json:"default_value,omitempty"`
	Choices      json.RawMessage `json:"choices,omitempty"`
	ChoiceMeta   json.RawMessage `json:"choice_meta,omitempty"`
	MinValue     *string         `json:"min_value,omitempty"`
	MaxValue     *string         `json:"max_value,omitempty"`
	Step         *string         `json:"step,omitempty"`
	Required     bool            `json:"required"`
	VisibleWhen  json.RawMessage `json:"visible_when,omitempty"`
	SortOrder    int             `json:"sort_order"`
	Notes        string          `json:"notes,omitempty"`
}

// CatalogMaterialRef is the material referenced by a catalog BOM line.
type CatalogMaterialRef struct {
	ID          string `json:"id"`
	SKU         string `json:"sku"`
	Name        string `json:"name"`
	Kind        string `json:"kind"`
	Category    string `json:"category"`
	UomCode     string `json:"uom_code"`
	DefaultCost string `json:"default_cost"`
}

// CatalogBomLine is one component of a template BOM. Quantities are decimal
// strings.
type CatalogBomLine struct {
	Component CatalogMaterialRef `json:"component"`
	Quantity  string             `json:"quantity"`
	UomCode   string             `json:"uom_code"`
	QtyBasis  string             `json:"qty_basis"`
	ScrapPct  string             `json:"scrap_pct"`
	Sequence  int                `json:"sequence"`
	Notes     string             `json:"notes,omitempty"`
}

// CatalogTemplateDetail is a template plus its gallery, options and BOM.
type CatalogTemplateDetail struct {
	CatalogTemplate
	Images  []CatalogImage   `json:"images"`
	Options []CatalogOption  `json:"options"`
	BOM     []CatalogBomLine `json:"bom"`
}

// CatalogKit is one kit (sub-assembly) of the global catalog.
type CatalogKit struct {
	ID              string          `json:"id"`
	SKU             string          `json:"sku"`
	Name            string          `json:"name"`
	Description     string          `json:"description,omitempty"`
	Category        string          `json:"category"`
	Line            string          `json:"line,omitempty"`
	Family          string          `json:"family,omitempty"`
	UomCode         string          `json:"uom_code"`
	DefaultCost     string          `json:"default_cost"`
	ThreedBinding   string          `json:"threed_binding,omitempty"`
	OptionAxes      json.RawMessage `json:"option_axes,omitempty"`
	PrimaryImageURL string          `json:"primary_image_url,omitempty"`
	ComponentCount  int             `json:"component_count"`
	UsedByCount     int             `json:"used_by_count"`
	Added           bool            `json:"added"`
}

// CatalogTemplateListParams filters GET /partner/v1/catalog/templates. The
// list is not paginated.
type CatalogTemplateListParams struct {
	Family string `json:"family,omitempty"`
	// Status is "available" or "coming_soon".
	Status string `json:"status,omitempty"`
	// Q is free text over name, subtitle and description.
	Q string `json:"q,omitempty"`
}

// CatalogKitListParams filters GET /partner/v1/catalog/kits. The list is not
// paginated.
type CatalogKitListParams struct {
	Category string `json:"category,omitempty"`
	Line     string `json:"line,omitempty"`
	// Q is free text over sku, name and description.
	Q string `json:"q,omitempty"`
}

// CatalogService calls /partner/v1/catalog (global product catalog).
type CatalogService struct{ c *Client }

// Templates lists catalog templates. Scope partner-api.catalog.read.
func (s *CatalogService) Templates(ctx context.Context, params CatalogTemplateListParams) ([]CatalogTemplate, error) {
	return listKey[CatalogTemplate](ctx, s.c, http.MethodGet, "/partner/v1/catalog/templates", "templates", params)
}

// Template returns one template with images, options and BOM, by slug. Scope
// partner-api.catalog.read.
func (s *CatalogService) Template(ctx context.Context, slug string) (*CatalogTemplateDetail, error) {
	return call[CatalogTemplateDetail](ctx, s.c, http.MethodGet, "/partner/v1/catalog/templates/"+seg(slug), nil, nil)
}

// Kits lists catalog kits. Scope partner-api.catalog.read.
func (s *CatalogService) Kits(ctx context.Context, params CatalogKitListParams) ([]CatalogKit, error) {
	return listKey[CatalogKit](ctx, s.c, http.MethodGet, "/partner/v1/catalog/kits", "kits", params)
}
