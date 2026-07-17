package partnerapi

// Shared Partner API types. Shapes mirror shedcloud-api-go/docs/PARTNER_API.md.

// SortOrder is asc or desc.
type SortOrder string

const (
	SortAsc  SortOrder = "asc"
	SortDesc SortOrder = "desc"
)

// PaginatedResponse is the list envelope for every Partner API list endpoint.
type PaginatedResponse[T any] struct {
	Data  []T   `json:"data"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Total int64 `json:"total"`
}

// PartnerCustomer is the nested customer object on sales entities.
type PartnerCustomer struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
	Phone string `json:"phone,omitempty"`
}

// PartnerSalesperson is the nested salesperson object on sales entities.
type PartnerSalesperson struct {
	// ID is the user id — resolvable against Client.Users.Get.
	ID    string `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}

// ExternalReferences are your own correlation ids stamped on a record (CRM
// deal id, ERP order id, …). Echoed in every DTO and event; filterable on
// lists via the ExternalRef param ("key:value").
type ExternalReferences map[string]string

// ExternalReferencesPatch is the PATCH shape for external references: keys
// are merged into the existing map; a nil value deletes that key.
type ExternalReferencesPatch map[string]*string

// PartnerRTO is the rent-to-own block on quotes/orders when the pricing
// payment type is "rto".
type PartnerRTO struct {
	TermMonths      int     `json:"termMonths,omitempty"`
	MonthlyPayment  float64 `json:"monthlyPayment,omitempty"`
	Rent            float64 `json:"rent,omitempty"`
	SecurityDeposit float64 `json:"securityDeposit,omitempty"`
	DownPayment     float64 `json:"downPayment,omitempty"`
	DamageWaiver    float64 `json:"damageWaiver,omitempty"`
	TotalDueToday   float64 `json:"totalDueToday,omitempty"`
	Balance         float64 `json:"balance,omitempty"`
	// AgreementID is the B2B agreement linked on the sales order (SOField1724).
	AgreementID string `json:"agreementId,omitempty"`
	// ProviderName is the financing provider's name — populated on
	// get-by-id only.
	ProviderName string `json:"providerName,omitempty"`
}

// PartnerDeposits is the money collected/owed on an order.
type PartnerDeposits struct {
	DownPayment     float64 `json:"downPayment,omitempty"`
	SecurityDeposit float64 `json:"securityDeposit,omitempty"`
	TotalDueToday   float64 `json:"totalDueToday,omitempty"`
	TotalPaid       float64 `json:"totalPaid,omitempty"`
	Balance         float64 `json:"balance,omitempty"`
}

// PartnerDelivery is the transportation-run block on a work order detail
// response (get-by-id only).
type PartnerDelivery struct {
	// ScheduledDate is the per-stop delivery date on the run assignment.
	ScheduledDate string `json:"scheduledDate,omitempty"`
	RunNumber     int    `json:"runNumber,omitempty"`
	// RunStatus is "Open", "Scheduled", "En Route", or "Delivered".
	RunStatus   string `json:"runStatus,omitempty"`
	DriverName  string `json:"driverName,omitempty"`
	DeliveredAt string `json:"deliveredAt,omitempty"`
}

// PartnerLocation is the nested location object.
type PartnerLocation struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Slug string `json:"slug,omitempty"`
}

// PartnerPricing is the nested pricing object on quotes and orders.
type PartnerPricing struct {
	Subtotal       float64 `json:"subtotal"`
	Total          float64 `json:"total"`
	MonthlyPayment float64 `json:"monthlyPayment,omitempty"`
	// PaymentType is "rto" or "cash".
	PaymentType    string  `json:"paymentType,omitempty"`
	ChangeOrderFee float64 `json:"changeOrderFee,omitempty"`
}

// LotStockAttributes is the exterior configuration of a lot-stock unit — the
// same values ShedCloud's work-order detail page shows. Custom colors surface
// as "Custom Color - #HEX".
type LotStockAttributes struct {
	Siding      string `json:"siding,omitempty"`
	SidingColor string `json:"sidingColor,omitempty"`
	TrimColor   string `json:"trimColor,omitempty"`
	// RoofMaterial is "Metal Roof" or "Shingle Roof".
	RoofMaterial string `json:"roofMaterial,omitempty"`
	RoofColor    string `json:"roofColor,omitempty"`
}

// LotStockItem is one on-lot inventory row.
type LotStockItem struct {
	ID           string   `json:"id"`
	WorkOrderID  string   `json:"workOrderId,omitempty"`
	SerialNumber string   `json:"serialNumber,omitempty"`
	Title        string   `json:"title,omitempty"`
	PurchaseType string   `json:"purchaseType,omitempty"`
	BasePrice    float64  `json:"basePrice,omitempty"`
	Price        float64  `json:"price,omitempty"`
	LocationID   string   `json:"locationId,omitempty"`
	LocationName string   `json:"locationName,omitempty"`
	LocationSlug string   `json:"locationSlug,omitempty"`
	Images       []string `json:"images,omitempty"`
	// HeroImages are public CDN URLs for uploaded Hero_Images / Product_Images
	// on the work order. Empty when none are registered.
	HeroImages []string `json:"heroImages,omitempty"`
	// Attributes is nil when the unit has no configurator.
	Attributes *LotStockAttributes `json:"attributes,omitempty"`
	Sold       bool                `json:"sold"`
}

// StockTemplateItem is one stock template — a buildable catalog design (not
// physical inventory; see LotStockItem for on-lot units).
type StockTemplateItem struct {
	ID              string `json:"id"`
	WorkOrderNumber string `json:"workOrderNumber,omitempty"`
	TemplateName    string `json:"templateName,omitempty"`
	ProductName     string `json:"productName,omitempty"`
	Description     string `json:"description,omitempty"`
	DescriptionHTML string `json:"descriptionHtml,omitempty"`
	// Tags are the template's public tags (internal tags are never exposed).
	Tags   []string `json:"tags"`
	Images []string `json:"images"`
	// Prices: base price, upgrades total, and the total of removed standard
	// features.
	BasePrice     *float64 `json:"basePrice,omitempty"`
	UpgradesPrice *float64 `json:"upgradesPrice,omitempty"`
	RemovedPrice  *float64 `json:"removedPrice,omitempty"`
	// ConfiguratorID is the template's 3D configuration; StructuraURL
	// deep-links the interactive 3D view of it.
	ConfiguratorID string `json:"configuratorId,omitempty"`
	StructuraURL   string `json:"structuraUrl,omitempty"`
	ProductID      string `json:"productId,omitempty"`
	SizeID         string `json:"sizeId,omitempty"`
}

// StockTemplateListParams are query params for GET /partner/v1/stock-templates.
type StockTemplateListParams struct {
	PaginationParams // Limit is capped at 60 by the server.
	// Search is a case-insensitive match on work-order number, template
	// name, or product name.
	Search string `json:"search,omitempty"`
	// Tags filters to templates carrying all of these public tags (joined
	// with commas for you).
	Tags []string `json:"-"`
}

// LeadItem is one lead.
type LeadItem struct {
	ID              string             `json:"id"`
	OrderNumber     int                `json:"orderNumber,omitempty"`
	Status          string             `json:"status,omitempty"`
	StatusChangedAt string             `json:"statusChangedAt,omitempty"`
	Customer        PartnerCustomer    `json:"customer"`
	Salesperson     PartnerSalesperson `json:"salesperson"`
	Location        PartnerLocation    `json:"location"`
	ExternalRefs    ExternalReferences `json:"externalReferences,omitempty"`
	// Version increments on every write; usable with WithIfMatch.
	Version   int64  `json:"version,omitempty"`
	CreatedAt string `json:"createdAt,omitempty"`
	UpdatedAt string `json:"updatedAt,omitempty"`
}

// QuoteItem is one quote.
type QuoteItem struct {
	ID              string             `json:"id"`
	OrderNumber     int                `json:"orderNumber,omitempty"`
	Status          string             `json:"status,omitempty"`
	StatusChangedAt string             `json:"statusChangedAt,omitempty"`
	Customer        PartnerCustomer    `json:"customer"`
	Salesperson     PartnerSalesperson `json:"salesperson"`
	Location        PartnerLocation    `json:"location"`
	Pricing         PartnerPricing     `json:"pricing"`
	// RTO is non-nil when Pricing.PaymentType is "rto".
	RTO                  *PartnerRTO `json:"rto,omitempty"`
	SerialNumber         string      `json:"serialNumber,omitempty"`
	WorkOrderID          string      `json:"workOrderId,omitempty"`
	Converted            bool        `json:"converted"`
	ConvertedOrderID     string      `json:"convertedOrderId,omitempty"`
	ConvertedOrderNumber int         `json:"convertedOrderNumber,omitempty"`
	// ValidUntil is the quote expiration timestamp (RFC 3339). When it passes
	// while the quote is still Open/Active, a quote.expired event is emitted.
	ValidUntil   string             `json:"validUntil,omitempty"`
	ExternalRefs ExternalReferences `json:"externalReferences,omitempty"`
	Version      int64              `json:"version,omitempty"`
	CreatedAt    string             `json:"createdAt,omitempty"`
	UpdatedAt    string             `json:"updatedAt,omitempty"`
}

// OrderItem is one sales order.
type OrderItem struct {
	ID              string             `json:"id"`
	OrderNumber     int                `json:"orderNumber,omitempty"`
	Status          string             `json:"status,omitempty"`
	StatusChangedAt string             `json:"statusChangedAt,omitempty"`
	Customer        PartnerCustomer    `json:"customer"`
	Salesperson     PartnerSalesperson `json:"salesperson"`
	Location        PartnerLocation    `json:"location"`
	Pricing         PartnerPricing     `json:"pricing"`
	// RTO is non-nil when Pricing.PaymentType is "rto".
	RTO *PartnerRTO `json:"rto,omitempty"`
	// Deposits is the money collected/owed on the order.
	Deposits          *PartnerDeposits `json:"deposits,omitempty"`
	SerialNumber      string           `json:"serialNumber,omitempty"`
	WorkOrderID       string           `json:"workOrderId,omitempty"`
	SourceQuoteID     string           `json:"sourceQuoteId,omitempty"`
	SourceQuoteNumber int              `json:"sourceQuoteNumber,omitempty"`
	// ExpectedDeliveryDate is the expected delivery date (RFC 3339).
	ExpectedDeliveryDate string `json:"expectedDeliveryDate,omitempty"`
	// DeliveredAt is the delivery completion timestamp (RFC 3339).
	DeliveredAt  string             `json:"deliveredAt,omitempty"`
	ExternalRefs ExternalReferences `json:"externalReferences,omitempty"`
	Version      int64              `json:"version,omitempty"`
	CreatedAt    string             `json:"createdAt,omitempty"`
	UpdatedAt    string             `json:"updatedAt,omitempty"`
}

// WorkOrderItem is one work order.
type WorkOrderItem struct {
	ID              string          `json:"id"`
	WorkOrderNumber int             `json:"workOrderNumber,omitempty"`
	SerialNumber    string          `json:"serialNumber,omitempty"`
	Title           string          `json:"title,omitempty"`
	Status          string          `json:"status,omitempty"`
	StatusChangedAt string          `json:"statusChangedAt,omitempty"`
	OrderID         string          `json:"orderId,omitempty"`
	OrderNumber     int             `json:"orderNumber,omitempty"`
	Location        PartnerLocation `json:"location"`
	BasePrice       float64         `json:"basePrice,omitempty"`
	PromisedDate    string          `json:"promisedDate,omitempty"`
	// Delivery is the transportation-run block — populated on get-by-id only.
	Delivery     *PartnerDelivery   `json:"delivery,omitempty"`
	ExternalRefs ExternalReferences `json:"externalReferences,omitempty"`
	Version      int64              `json:"version,omitempty"`
	CreatedAt    string             `json:"createdAt,omitempty"`
	UpdatedAt    string             `json:"updatedAt,omitempty"`
}

// LocationItem is one full location record (sales lot, plant, warehouse).
type LocationItem struct {
	ID            string   `json:"id"`
	Name          string   `json:"name,omitempty"`
	Slug          string   `json:"slug,omitempty"`
	Code          string   `json:"code,omitempty"`
	Address       string   `json:"address,omitempty"`
	City          string   `json:"city,omitempty"`
	State         string   `json:"state,omitempty"`
	ZipCode       string   `json:"zipCode,omitempty"`
	Phone         string   `json:"phone,omitempty"`
	ContactPerson string   `json:"contactPerson,omitempty"`
	ContactEmail  string   `json:"contactEmail,omitempty"`
	Latitude      *float64 `json:"latitude,omitempty"`
	Longitude     *float64 `json:"longitude,omitempty"`
	Active        bool     `json:"active"`
	SalesLot      bool     `json:"salesLot"`
	Plant         bool     `json:"plant"`
	// Timezone is the company-level IANA timezone. Store hours are
	// interpreted in this zone.
	Timezone string `json:"timezone,omitempty"`
	// Region is a dealer-defined grouping label (e.g. "Southeast").
	Region string `json:"region,omitempty"`
	// StoreHours is the weekly store schedule keyed by day ("mon".."sun").
	// Empty when the location has no schedule configured.
	StoreHours map[string]StoreHoursDay `json:"storeHours,omitempty"`
	CreatedAt  string                   `json:"createdAt,omitempty"`
	UpdatedAt  string                   `json:"updatedAt,omitempty"`
}

// StoreHoursDay is one day inside LocationItem.StoreHours. Times are 24-hour
// "HH:MM" strings in the location's local time. Disabled days may still carry
// their last from/to values.
type StoreHoursDay struct {
	Enabled bool   `json:"enabled"`
	From    string `json:"from,omitempty"`
	To      string `json:"to,omitempty"`
}

// CustomerItem is one full customer record.
type CustomerItem struct {
	ID            string `json:"id"`
	Name          string `json:"name,omitempty"`
	ContactName   string `json:"contactName,omitempty"`
	ContactPerson string `json:"contactPerson,omitempty"`
	Email         string `json:"email,omitempty"`
	Phone         string `json:"phone,omitempty"`
	Address       string `json:"address,omitempty"`
	City          string `json:"city,omitempty"`
	State         string `json:"state,omitempty"`
	ZipCode       string `json:"zipCode,omitempty"`
	Code          string `json:"code,omitempty"`
	Active        bool   `json:"active"`
	// Merge lineage: a merged duplicate keeps resolving by id but carries
	// Merged=true and points at the surviving record.
	Merged       bool               `json:"merged,omitempty"`
	MergedInto   string             `json:"mergedInto,omitempty"`
	MergedAt     string             `json:"mergedAt,omitempty"`
	ExternalRefs ExternalReferences `json:"externalReferences,omitempty"`
	Version      int64              `json:"version,omitempty"`
	CreatedAt    string             `json:"createdAt,omitempty"`
	UpdatedAt    string             `json:"updatedAt,omitempty"`
}

// ProductItem is one finished catalog product (read-only). Raw materials and
// kits are never returned by the Partner API.
type ProductItem struct {
	ID          string  `json:"id"`
	Name        string  `json:"name,omitempty"`
	Description string  `json:"description,omitempty"`
	SKU         string  `json:"sku,omitempty"`
	Price       float64 `json:"price"`
	Width       float64 `json:"width,omitempty"`
	Length      float64 `json:"length,omitempty"`
	Active      bool    `json:"active"`
	// Images are ready-to-use public URLs: the uploaded gallery (upload
	// order) first, then legacy image slots on the product record.
	Images    []string `json:"images,omitempty"`
	CreatedAt string   `json:"createdAt,omitempty"`
	UpdatedAt string   `json:"updatedAt,omitempty"`
}

// PaginationParams are shared by every list endpoint.
type PaginationParams struct {
	Page  int `json:"page,omitempty"`
	Limit int `json:"limit,omitempty"`
}

// LotStockListParams are query params for GET /partner/v1/lot-stock.
type LotStockListParams struct {
	PaginationParams
	// PurchaseType is ALL (the default — no category filter), Lot Stock,
	// Rental Return, or Immediate Sale. Case-insensitive; hyphens and
	// underscores are accepted (e.g. "lot-stock").
	PurchaseType string `json:"purchaseType,omitempty"`
	Location     string `json:"location,omitempty"`
	// Region filters to units at locations carrying this exact region label
	// (dealer-defined). Combined with Location, both constraints must hold.
	Region string    `json:"region,omitempty"`
	Search string    `json:"search,omitempty"`
	Sort   string    `json:"sort,omitempty"` // serialNumber | title | price | createdAt
	Order  SortOrder `json:"order,omitempty"`
}

// SalesListParams are shared list filters for leads, quotes, and orders.
type SalesListParams struct {
	PaginationParams
	Search        string    `json:"search,omitempty"`
	Sort          string    `json:"sort,omitempty"` // orderNumber | customerName | status | total | createdAt | updatedAt
	Order         SortOrder `json:"order,omitempty"`
	Status        string    `json:"status,omitempty"`
	Location      string    `json:"location,omitempty"`
	CustomerEmail string    `json:"customerEmail,omitempty"`
	CustomerPhone string    `json:"customerPhone,omitempty"`
	OrderNumber   string    `json:"orderNumber,omitempty"`
	Salesperson   string    `json:"salesperson,omitempty"`
	CreatedFrom   string    `json:"createdFrom,omitempty"`
	CreatedTo     string    `json:"createdTo,omitempty"`
	UpdatedFrom   string    `json:"updatedFrom,omitempty"`
	UpdatedTo     string    `json:"updatedTo,omitempty"`
	// ExternalRef is a "key:value" exact match on your external references.
	ExternalRef string `json:"externalRef,omitempty"`
}

// QuoteListParams extends SalesListParams with converted.
type QuoteListParams struct {
	SalesListParams
	// Converted filters converted quotes when non-nil.
	Converted *bool `json:"converted,omitempty"`
}

// OrderListParams extends SalesListParams with payment/serial filters.
type OrderListParams struct {
	SalesListParams
	PaymentType  string `json:"paymentType,omitempty"` // rto | cash
	SerialNumber string `json:"serialNumber,omitempty"`
}

// WorkOrderListParams are query params for GET /partner/v1/work-orders.
type WorkOrderListParams struct {
	PaginationParams
	Search        string    `json:"search,omitempty"`
	Sort          string    `json:"sort,omitempty"` // workOrderNumber | serialNumber | status | createdAt | updatedAt
	Order         SortOrder `json:"order,omitempty"`
	Status        string    `json:"status,omitempty"`
	SerialNumber  string    `json:"serialNumber,omitempty"`
	OrderNumber   string    `json:"orderNumber,omitempty"`
	LinkedOrderID string    `json:"linkedOrderId,omitempty"`
	Location      string    `json:"location,omitempty"`
	CreatedFrom   string    `json:"createdFrom,omitempty"`
	CreatedTo     string    `json:"createdTo,omitempty"`
	UpdatedFrom   string    `json:"updatedFrom,omitempty"`
	UpdatedTo     string    `json:"updatedTo,omitempty"`
	// ExternalRef is a "key:value" exact match on your external references.
	ExternalRef string `json:"externalRef,omitempty"`
}

// LocationListParams are query params for GET /partner/v1/locations.
type LocationListParams struct {
	PaginationParams
	Search string    `json:"search,omitempty"`
	Sort   string    `json:"sort,omitempty"` // name | code | city | createdAt | updatedAt
	Order  SortOrder `json:"order,omitempty"`
	// Boolean flag filters apply only when non-nil (so false is still sent).
	Active   *bool `json:"active,omitempty"`
	SalesLot *bool `json:"salesLot,omitempty"`
	Plant    *bool `json:"plant,omitempty"`
	// Region is an exact match on the dealer-defined region label.
	Region string `json:"region,omitempty"`
}

// CustomerListParams are query params for GET /partner/v1/customers.
type CustomerListParams struct {
	PaginationParams
	Search      string    `json:"search,omitempty"`
	Sort        string    `json:"sort,omitempty"` // name | email | createdAt | updatedAt
	Order       SortOrder `json:"order,omitempty"`
	Email       string    `json:"email,omitempty"` // substring match
	Phone       string    `json:"phone,omitempty"` // substring match
	CreatedFrom string    `json:"createdFrom,omitempty"`
	CreatedTo   string    `json:"createdTo,omitempty"`
	UpdatedFrom string    `json:"updatedFrom,omitempty"`
	UpdatedTo   string    `json:"updatedTo,omitempty"`
	// ExternalRef is a "key:value" exact match on your external references.
	ExternalRef string `json:"externalRef,omitempty"`
	// IncludeMerged includes customers that were merged into another record
	// (excluded from lists by default).
	IncludeMerged bool `json:"includeMerged,omitempty"`
}

// ProductListParams are query params for GET /partner/v1/products.
type ProductListParams struct {
	PaginationParams
	Search string    `json:"search,omitempty"`
	Sort   string    `json:"sort,omitempty"` // name | sku | price | createdAt | updatedAt
	Order  SortOrder `json:"order,omitempty"`
	SKU    string    `json:"sku,omitempty"` // substring match
	// Active filters by the active flag when non-nil.
	Active      *bool  `json:"active,omitempty"`
	CreatedFrom string `json:"createdFrom,omitempty"`
	CreatedTo   string `json:"createdTo,omitempty"`
	UpdatedFrom string `json:"updatedFrom,omitempty"`
	UpdatedTo   string `json:"updatedTo,omitempty"`
}

// ProductCreateRequest is the body for POST /partner/v1/products.
type ProductCreateRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	SKU         string   `json:"sku,omitempty"`
	Price       *float64 `json:"price,omitempty"`
	Width       *float64 `json:"width,omitempty"`
	Length      *float64 `json:"length,omitempty"`
	// Active defaults to true.
	Active *bool `json:"active,omitempty"`
	// LineID links the new product under an existing product line (parent
	// product).
	LineID string `json:"lineId,omitempty"`
}

// ProductPatchRequest is the body for PATCH /partner/v1/products/{id}.
type ProductPatchRequest struct {
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	SKU         string   `json:"sku,omitempty"`
	Price       *float64 `json:"price,omitempty"`
	Width       *float64 `json:"width,omitempty"`
	Length      *float64 `json:"length,omitempty"`
	Active      *bool    `json:"active,omitempty"`
}

// ProductSizeCreateRequest is the body for POST /partner/v1/products/{id}/sizes.
type ProductSizeCreateRequest struct {
	// Name defaults to "<width>x<length>" when omitted.
	Name   string  `json:"name,omitempty"`
	SKU    string  `json:"sku,omitempty"`
	Width  float64 `json:"width"`
	Length float64 `json:"length"`
	Price  float64 `json:"price"`
	// Active defaults to true.
	Active *bool `json:"active,omitempty"`
}

// ProductSizeItem is one size child of a catalog product. Sizes carry the
// model's real pricing and dimensions; they never appear in the products
// list themselves.
type ProductSizeItem struct {
	ID        string  `json:"id"`
	ProductID string  `json:"productId"`
	Name      string  `json:"name,omitempty"`
	SKU       string  `json:"sku,omitempty"`
	Width     float64 `json:"width"`
	Length    float64 `json:"length"`
	Price     float64 `json:"price"`
	Active    bool    `json:"active"`
}

// LeadPatchRequest is the body for PATCH /partner/v1/leads/{id}.
type LeadPatchRequest struct {
	SalesLocation    string `json:"salesLocation,omitempty"`
	SalespersonName  string `json:"salespersonName,omitempty"`
	SalespersonEmail string `json:"salespersonEmail,omitempty"`
	// ExternalRefs keys are merged into the record's map; nil values delete keys.
	ExternalRefs ExternalReferencesPatch `json:"externalReferences,omitempty"`
}

// LeadCreateCustomer is the customer block of LeadCreateRequest. At least one
// of name, email, or phone is required.
type LeadCreateCustomer struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
	Phone string `json:"phone,omitempty"`
}

// LeadCreateRequest is the body for POST /partner/v1/leads. When the
// salesperson fields are omitted, the location's lead-routing strategy
// auto-assigns one.
type LeadCreateRequest struct {
	LocationID       string             `json:"locationId"`
	Customer         LeadCreateCustomer `json:"customer"`
	SalespersonName  string             `json:"salespersonName,omitempty"`
	SalespersonEmail string             `json:"salespersonEmail,omitempty"`
	ExternalRefs     ExternalReferences `json:"externalReferences,omitempty"`
}

// QuoteConvertRequest is the optional body for
// POST /partner/v1/quotes/{id}/convert — overrides the salesperson copied
// from the quote.
type QuoteConvertRequest struct {
	SalespersonName  string `json:"salespersonName,omitempty"`
	SalespersonEmail string `json:"salespersonEmail,omitempty"`
}

// QuoteCreateCustomer is the customer block of QuoteCreateRequest. An email,
// or a name + phone, is required; the customer is matched by email when one
// already exists, otherwise created.
type QuoteCreateCustomer struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
	Phone string `json:"phone,omitempty"`
}

// QuoteCreateDeliveryAddress is the optional delivery address stamped onto
// the created quote (and cascaded to the linked work order).
type QuoteCreateDeliveryAddress struct {
	Address   string   `json:"address,omitempty"`
	City      string   `json:"city,omitempty"`
	State     string   `json:"state,omitempty"`
	ZipCode   string   `json:"zipCode,omitempty"`
	Latitude  *float64 `json:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty"`
}

// QuoteCreateRequest is the body for POST /partner/v1/quotes — create a quote
// from an in-stock (on-lot) unit. Identify the unit with SerialNumber or
// WorkOrderID (WorkOrderID wins when both are sent).
type QuoteCreateRequest struct {
	SerialNumber string `json:"serialNumber,omitempty"`
	WorkOrderID  string `json:"workOrderId,omitempty"`
	LocationID   string `json:"locationId,omitempty"`
	// PurchaseType is "Lot Stock", "Rental Return", or "Immediate Sale";
	// inferred from the unit when omitted.
	PurchaseType    string                      `json:"purchaseType,omitempty"`
	Price           *float64                    `json:"price,omitempty"`
	Note            string                      `json:"note,omitempty"`
	Customer        QuoteCreateCustomer         `json:"customer"`
	DeliveryAddress *QuoteCreateDeliveryAddress `json:"deliveryAddress,omitempty"`
	ExternalRefs    ExternalReferences          `json:"externalReferences,omitempty"`
}

// QuotePatchRequest is the body for PATCH /partner/v1/quotes/{id}.
type QuotePatchRequest struct {
	CustomerName     string `json:"customerName,omitempty"`
	CustomerEmail    string `json:"customerEmail,omitempty"`
	CustomerPhone    string `json:"customerPhone,omitempty"`
	SalespersonName  string `json:"salespersonName,omitempty"`
	SalespersonEmail string `json:"salespersonEmail,omitempty"`
	SalesLocation    string `json:"salesLocation,omitempty"`
	DeliveryAddress  string `json:"deliveryAddress,omitempty"`
	DeliveryCity     string `json:"deliveryCity,omitempty"`
	DeliveryState    string `json:"deliveryState,omitempty"`
	DeliveryZipCode  string `json:"deliveryZipCode,omitempty"`
	// ValidUntil sets the quote's expiration ("YYYY-MM-DD" or RFC 3339;
	// date-only values cover the whole day).
	ValidUntil *string `json:"validUntil,omitempty"`
	// ExternalRefs keys are merged into the record's map; nil values delete keys.
	ExternalRefs ExternalReferencesPatch `json:"externalReferences,omitempty"`
}

// OrderCreateLineItem is one upgrade line on OrderCreateRequest.
type OrderCreateLineItem struct {
	ProductID string   `json:"productId"`
	Quantity  float64  `json:"quantity,omitempty"` // default 1
	Price     *float64 `json:"price,omitempty"`
	// LineKey is an optional stable per-line identifier (generated when
	// omitted); it becomes the line's handle for delete.
	LineKey string `json:"lineKey,omitempty"`
}

// OrderCreateConfiguration is the optional configurator payload on
// OrderCreateRequest. Pass the portal-shaped SelectedCategories array
// verbatim, or the flat display attributes (converted into category entries
// server-side).
type OrderCreateConfiguration struct {
	Model              string           `json:"model,omitempty"`
	SelectedCategories []map[string]any `json:"selectedCategories,omitempty"`
	Siding             string           `json:"siding,omitempty"`
	SidingColor        string           `json:"sidingColor,omitempty"`
	TrimColor          string           `json:"trimColor,omitempty"`
	// RoofMaterial is "Metal Roof" or "Shingle Roof".
	RoofMaterial string `json:"roofMaterial,omitempty"`
	RoofColor    string `json:"roofColor,omitempty"`
}

// OrderCreateRequest is the body for POST /partner/v1/orders — full order
// creation (customer, location, base product + optional size, upgrades,
// pricing header, configurator). New orders start in status "Unsubmitted".
type OrderCreateRequest struct {
	// CustomerID references an existing customer; alternatively set Customer
	// (matched by email or created).
	CustomerID string               `json:"customerId,omitempty"`
	Customer   *QuoteCreateCustomer `json:"customer,omitempty"`

	LocationID    string `json:"locationId"`
	SalespersonID string `json:"salespersonId,omitempty"`

	ProductID string `json:"productId"`
	SizeID    string `json:"sizeId,omitempty"`

	BasePrice *float64 `json:"basePrice,omitempty"`
	Subtotal  *float64 `json:"subtotal,omitempty"`
	Total     *float64 `json:"total,omitempty"`
	// PaymentType is "cash" or "rto".
	PaymentType string `json:"paymentType,omitempty"`

	Upgrades      []OrderCreateLineItem     `json:"upgrades,omitempty"`
	Configuration *OrderCreateConfiguration `json:"configuration,omitempty"`

	DeliveryAddress *QuoteCreateDeliveryAddress `json:"deliveryAddress,omitempty"`
	Note            string                      `json:"note,omitempty"`
	ExternalRefs    ExternalReferences          `json:"externalReferences,omitempty"`
}

// PaymentCreateRequest is the body for POST /partner/v1/orders/{id}/payments.
// Manual methods only — card/ACH records are created by the Stripe pipeline
// (use payment links instead).
type PaymentCreateRequest struct {
	// Method is one of: cash, check, financed, manual.
	Method string  `json:"method"`
	Amount float64 `json:"amount"`
	// Description defaults to "<Method> payment".
	Description string `json:"description,omitempty"`
	// CheckNumber applies to method "check" only.
	CheckNumber string `json:"checkNumber,omitempty"`
	// BankName applies to method "financed" only.
	BankName string `json:"bankName,omitempty"`
}

// PaymentLinkCreateRequest is the body for POST /partner/v1/orders/{id}/payment-links.
type PaymentLinkCreateRequest struct {
	Amount float64 `json:"amount"`
	// Name labels the Stripe Checkout line item (defaults to "Order #<n>").
	Name string `json:"name,omitempty"`
	// Email overrides the order's customer email as the link recipient.
	Email string `json:"email,omitempty"`
	// SendEmail controls the customer payment-link email (default true).
	SendEmail *bool `json:"sendEmail,omitempty"`
	// Currency defaults to "usd".
	Currency string `json:"currency,omitempty"`
}

// PaymentLinkResponse is returned by POST /partner/v1/orders/{id}/payment-links.
type PaymentLinkResponse struct {
	URL       string `json:"url"`
	SessionID string `json:"sessionId"`
	ExpiresAt string `json:"expiresAt,omitempty"`
	EmailSent bool   `json:"emailSent"`
}

// WorkOrderCreateRequest is the body for POST /partner/v1/work-orders.
type WorkOrderCreateRequest struct {
	// LocationID is the building location id or slug (required).
	LocationID string `json:"locationId"`
	// PurchaseType is one of: new-build, existing-physical-inventory, general, stock.
	PurchaseType string `json:"purchaseType,omitempty"`
	// WorkOrderType is one of: made-to-order, lot-stock, rental-return, immediate-sale, template.
	WorkOrderType string `json:"workOrderType,omitempty"`
	// DeliveryType is one of: delivery-from-factory, built-on-site,
	// delivered-from-dealer, delivered-from-lot-stock.
	DeliveryType string `json:"deliveryType,omitempty"`
	// SerialNumber is optional; uniqueness is enforced per company (409).
	SerialNumber string `json:"serialNumber,omitempty"`
	Title        string `json:"title,omitempty"`
	// SizeID optionally attaches a product size to the new work order.
	SizeID string `json:"sizeId,omitempty"`
	// PromisedDate is RFC 3339 or YYYY-MM-DD.
	PromisedDate string             `json:"promisedDate,omitempty"`
	ExternalRefs ExternalReferences `json:"externalReferences,omitempty"`
}

// LineItemCreateRequest is the body for POST /partner/v1/{quotes|orders}/{id}/line-items.
type LineItemCreateRequest struct {
	ProductID string   `json:"productId"`
	Quantity  float64  `json:"quantity,omitempty"` // default 1
	Price     *float64 `json:"price,omitempty"`
	// LineKey makes the add idempotent: re-sending the same key returns the
	// existing line instead of duplicating it.
	LineKey string `json:"lineKey,omitempty"`
}

// LineItemCreateResponse echoes the created (or matched) line.
type LineItemCreateResponse struct {
	LineID    string   `json:"lineId"`
	ProductID string   `json:"productId"`
	Quantity  float64  `json:"quantity"`
	Price     *float64 `json:"price,omitempty"`
	// Created is false when LineKey matched an existing line.
	Created bool `json:"created"`
}

// OrderPatchRequest is the body for PATCH /partner/v1/orders/{id}.
type OrderPatchRequest struct {
	CustomerName     string `json:"customerName,omitempty"`
	CustomerEmail    string `json:"customerEmail,omitempty"`
	CustomerPhone    string `json:"customerPhone,omitempty"`
	SalespersonName  string `json:"salespersonName,omitempty"`
	SalespersonEmail string `json:"salespersonEmail,omitempty"`
	SalesLocation    string `json:"salesLocation,omitempty"`
	DeliveryAddress  string `json:"deliveryAddress,omitempty"`
	DeliveryCity     string `json:"deliveryCity,omitempty"`
	DeliveryState    string `json:"deliveryState,omitempty"`
	DeliveryZipCode  string `json:"deliveryZipCode,omitempty"`
	// ExternalRefs keys are merged into the record's map; nil values delete keys.
	ExternalRefs ExternalReferencesPatch `json:"externalReferences,omitempty"`
}

// WorkOrderPatchRequest is the body for PATCH /partner/v1/work-orders/{id}.
type WorkOrderPatchRequest struct {
	Title            string `json:"title,omitempty"`
	Description      string `json:"description,omitempty"`
	BuildingLocation string `json:"buildingLocation,omitempty"`
	PromisedDate     string `json:"promisedDate,omitempty"`
	// ExternalRefs keys are merged into the record's map; nil values delete keys.
	ExternalRefs ExternalReferencesPatch `json:"externalReferences,omitempty"`
}

// LocationCreateRequest is the body for POST /partner/v1/locations.
// Name is required, plus either an Address or a Latitude/Longitude pair
// (provided together).
type LocationCreateRequest struct {
	Name          string   `json:"name"`
	Slug          string   `json:"slug,omitempty"`
	Code          string   `json:"code,omitempty"`
	Address       string   `json:"address,omitempty"`
	City          string   `json:"city,omitempty"`
	State         string   `json:"state,omitempty"`
	ZipCode       string   `json:"zipCode,omitempty"`
	Phone         string   `json:"phone,omitempty"`
	ContactPerson string   `json:"contactPerson,omitempty"`
	ContactEmail  string   `json:"contactEmail,omitempty"`
	Latitude      *float64 `json:"latitude,omitempty"`
	Longitude     *float64 `json:"longitude,omitempty"`
	Active        *bool    `json:"active,omitempty"`
	SalesLot      *bool    `json:"salesLot,omitempty"`
	Plant         *bool    `json:"plant,omitempty"`
}

// LocationPatchRequest is the body for PATCH /partner/v1/locations/{id}.
// Same key set as create, all optional.
type LocationPatchRequest struct {
	Name          string   `json:"name,omitempty"`
	Slug          string   `json:"slug,omitempty"`
	Code          string   `json:"code,omitempty"`
	Address       string   `json:"address,omitempty"`
	City          string   `json:"city,omitempty"`
	State         string   `json:"state,omitempty"`
	ZipCode       string   `json:"zipCode,omitempty"`
	Phone         string   `json:"phone,omitempty"`
	ContactPerson string   `json:"contactPerson,omitempty"`
	ContactEmail  string   `json:"contactEmail,omitempty"`
	Latitude      *float64 `json:"latitude,omitempty"`
	Longitude     *float64 `json:"longitude,omitempty"`
	Active        *bool    `json:"active,omitempty"`
	SalesLot      *bool    `json:"salesLot,omitempty"`
	Plant         *bool    `json:"plant,omitempty"`
}

// DomainProduct is one product mapping inside a domain's location entry.
type DomainProduct struct {
	ProductID string `json:"productId"`
	// Sizes are the size ids enabled for this product on this domain.
	Sizes []string `json:"sizes,omitempty"`
	// ExcludedUpgrades are upgrade ids hidden for this product on this domain.
	ExcludedUpgrades []string `json:"excludedUpgrades,omitempty"`
}

// DomainLocation is one location assignment inside a domain.
type DomainLocation struct {
	LocationID string `json:"locationId"`
	Name       string `json:"name,omitempty"`
	Code       string `json:"code,omitempty"`
	// HidePrices hides pricing for this location on this domain's storefront.
	HidePrices bool `json:"hidePrices"`
	// DefaultForStore marks this domain as the location's default storefront
	// domain. A location has at most one default domain.
	DefaultForStore bool            `json:"defaultForStore"`
	Products        []DomainProduct `json:"products,omitempty"`
}

// DomainItem is a white-label storefront subdomain from GET /partner/v1/domains.
type DomainItem struct {
	// IntegrationID is the DNS integration record the subdomain belongs to.
	IntegrationID string `json:"integrationId"`
	Subdomain     string `json:"subdomain"`
	ApexDomain    string `json:"apexDomain,omitempty"`
	// Verified reports DNS verification (per-subdomain, falling back to the
	// integration flag).
	Verified  bool             `json:"verified"`
	Locations []DomainLocation `json:"locations"`
}

// DomainListParams are query params for GET /partner/v1/domains.
type DomainListParams struct {
	PaginationParams
	// LocationID narrows to domains assigned to this location.
	LocationID string `json:"locationId,omitempty"`
	// DefaultForStore true → only location entries flagged Default For Store.
	DefaultForStore *bool `json:"defaultForStore,omitempty"`
}

// LocationDomainListParams are query params for GET /partner/v1/locations/{id}/domains.
type LocationDomainListParams struct {
	PaginationParams
	// DefaultForStore true → only the entry flagged Default For Store.
	DefaultForStore *bool `json:"defaultForStore,omitempty"`
}

// AgreementCounterparty is one endpoint of a B2B agreement.
type AgreementCounterparty struct {
	CompanyID   string `json:"companyId"`
	CompanyType string `json:"companyType,omitempty"`
	Name        string `json:"name,omitempty"`
}

// AgreementTermRate is one term/monthly-rate pair in a rate schedule.
type AgreementTermRate struct {
	Months int     `json:"months"`
	Rate   float64 `json:"rate"`
}

// AgreementSecurityDepositRule is a tiered security-deposit percent rule.
type AgreementSecurityDepositRule struct {
	MinAmount float64 `json:"minAmount"`
	Percent   float64 `json:"percent"`
}

// AgreementStateRateSchedule is per-state RTO rate configuration.
type AgreementStateRateSchedule struct {
	StateCode              string                         `json:"stateCode"`
	TermRates              []AgreementTermRate            `json:"termRates,omitempty"`
	SecurityDepositOptions []float64                      `json:"securityDepositOptions,omitempty"`
	SecurityDepositRules   []AgreementSecurityDepositRule `json:"securityDepositRules,omitempty"`
}

// AgreementSecurityDepositPolicy is the agreement-level security deposit policy.
type AgreementSecurityDepositPolicy struct {
	Type         string   `json:"type,omitempty"`
	OtherPercent *float64 `json:"otherPercent,omitempty"`
}

// AgreementRateConfig is the active or pending rate configuration on an agreement.
type AgreementRateConfig struct {
	States          []AgreementStateRateSchedule   `json:"states,omitempty"`
	SecurityDeposit AgreementSecurityDepositPolicy `json:"securityDeposit,omitempty"`
}

// AgreementPendingUpdateSummary is a redacted pending rate-change summary.
type AgreementPendingUpdateSummary struct {
	Status      string `json:"status,omitempty"`
	RequestedAt string `json:"requestedAt,omitempty"`
	Reason      string `json:"reason,omitempty"`
}

// AgreementItem is a B2B partnership agreement from GET /partner/v1/agreements.
type AgreementItem struct {
	ID                       string                         `json:"id"`
	Direction                string                         `json:"direction"`
	Status                   string                         `json:"status"`
	From                     AgreementCounterparty          `json:"from"`
	To                       AgreementCounterparty          `json:"to"`
	RateConfig               *AgreementRateConfig           `json:"rateConfig,omitempty"`
	RateConfigVersion        int64                          `json:"rateConfigVersion,omitempty"`
	PendingRateConfig        *AgreementRateConfig           `json:"pendingRateConfig,omitempty"`
	PendingRateConfigVersion int64                          `json:"pendingRateConfigVersion,omitempty"`
	PendingUpdate            *AgreementPendingUpdateSummary `json:"pendingUpdate,omitempty"`
	OrderCount               int64                          `json:"orderCount,omitempty"`
	HasOrders                bool                           `json:"hasOrders,omitempty"`
	CreatedAt                string                         `json:"createdAt,omitempty"`
	UpdatedAt                string                         `json:"updatedAt,omitempty"`
}

// AgreementListParams are query params for GET /partner/v1/agreements.
type AgreementListParams struct {
	PaginationParams
	Status string `json:"status,omitempty"`
	Search string `json:"search,omitempty"`
}

// AgreementStateLegalSection is one section in a state legal appendix.
type AgreementStateLegalSection struct {
	SortOrder int    `json:"sortOrder"`
	Title     string `json:"title,omitempty"`
	Body      string `json:"body,omitempty"`
}

// AgreementStateLegalItem is per-state RTO legal configuration.
type AgreementStateLegalItem struct {
	ID                       string                       `json:"id"`
	AgreementID              string                       `json:"agreementId"`
	StateCode                string                       `json:"stateCode"`
	LessorLegalName          string                       `json:"lessorLegalName,omitempty"`
	LessorAddressLine1       string                       `json:"lessorAddressLine1,omitempty"`
	LessorAddressLine2       string                       `json:"lessorAddressLine2,omitempty"`
	Sections                 []AgreementStateLegalSection `json:"sections,omitempty"`
	ArbitrationText          string                       `json:"arbitrationText,omitempty"`
	AcknowledgementText      string                       `json:"acknowledgementText,omitempty"`
	IncludeCustomerDataSheet bool                         `json:"includeCustomerDataSheet,omitempty"`
	IncludeStateRpaPages     bool                         `json:"includeStateRpaPages,omitempty"`
	CreatedAt                string                       `json:"createdAt,omitempty"`
	UpdatedAt                string                       `json:"updatedAt,omitempty"`
}

// AgreementStateLegalListResponse is the list envelope for state legal rows.
type AgreementStateLegalListResponse struct {
	Data []AgreementStateLegalItem `json:"data"`
}

// CustomerCreateRequest is the body for POST /partner/v1/customers.
// Email is required and must be unique within the company.
type CustomerCreateRequest struct {
	Email         string             `json:"email"`
	Name          string             `json:"name,omitempty"`
	ContactName   string             `json:"contactName,omitempty"`
	ContactPerson string             `json:"contactPerson,omitempty"`
	Phone         string             `json:"phone,omitempty"`
	Address       string             `json:"address,omitempty"`
	City          string             `json:"city,omitempty"`
	State         string             `json:"state,omitempty"`
	ZipCode       string             `json:"zipCode,omitempty"`
	Code          string             `json:"code,omitempty"`
	ExternalRefs  ExternalReferences `json:"externalReferences,omitempty"`
}

// CustomerPatchRequest is the body for PATCH /partner/v1/customers/{id}.
type CustomerPatchRequest struct {
	Name          string `json:"name,omitempty"`
	ContactName   string `json:"contactName,omitempty"`
	ContactPerson string `json:"contactPerson,omitempty"`
	Email         string `json:"email,omitempty"`
	Phone         string `json:"phone,omitempty"`
	Address       string `json:"address,omitempty"`
	City          string `json:"city,omitempty"`
	State         string `json:"state,omitempty"`
	ZipCode       string `json:"zipCode,omitempty"`
	Code          string `json:"code,omitempty"`
	Active        *bool  `json:"active,omitempty"`
	// ExternalRefs keys are merged into the record's map; nil values delete keys.
	ExternalRefs ExternalReferencesPatch `json:"externalReferences,omitempty"`
}

// CustomerMergeRequest is the body for POST /partner/v1/customers/{id}/merge.
type CustomerMergeRequest struct {
	// SurvivorID is the customer that absorbs the duplicate in the path.
	SurvivorID string `json:"survivorId"`
}

// CustomerMergeResponse reports the merge outcome.
type CustomerMergeResponse struct {
	Survivor CustomerItem `json:"survivor"`
	// MergedCustomerID is the merged duplicate's id — it now resolves with
	// Merged=true and MergedInto pointing at the survivor.
	MergedCustomerID string `json:"mergedCustomerId"`
	// RelinkedSalesEntities counts the leads/quotes/orders moved to the
	// survivor.
	RelinkedSalesEntities int64 `json:"relinkedSalesEntities"`
}

// StatusUpdateRequest is the body for POST /partner/v1/{resource}/{id}/status.
type StatusUpdateRequest struct {
	Status            string `json:"status"`
	ActionDescription string `json:"actionDescription,omitempty"`
}

// MeResponse describes the authenticated credential (GET /partner/v1/me):
// who it belongs to and what it may do.
type MeResponse struct {
	CompanyID   string `json:"companyId"`
	CompanyName string `json:"companyName,omitempty"`
	// CredentialType is "api_key" or "oauth_client".
	CredentialType string      `json:"credentialType"`
	Scopes         []string    `json:"scopes"`
	RateLimit      MeRateLimit `json:"rateLimit"`
}

// MeRateLimit reports the per-credential throttle configuration.
type MeRateLimit struct {
	RequestsPerSecond int `json:"requestsPerSecond"`
	Burst             int `json:"burst"`
}

// OAuthTokenResponse is the JSON body from POST /oauth/token.
type OAuthTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope,omitempty"`
}

// UserItem is one salesperson/user record — resolve PartnerSalesperson.ID
// from sales DTOs into a full profile.
type UserItem struct {
	ID     string `json:"id"`
	Name   string `json:"name,omitempty"`
	Email  string `json:"email,omitempty"`
	Phone  string `json:"phone,omitempty"`
	Active bool   `json:"active"`
	// LocationIDs the user can access. Empty when AllLocations is true.
	LocationIDs       []string `json:"locationIds"`
	AllLocations      bool     `json:"allLocations"`
	InLeadRoutingPool bool     `json:"inLeadRoutingPool"`
	CreatedAt         string   `json:"createdAt,omitempty"`
}

// UserListParams are query params for GET /partner/v1/users.
type UserListParams struct {
	PaginationParams
	Search string `json:"search,omitempty"`
	// Active filters by the active flag when non-nil.
	Active *bool `json:"active,omitempty"`
}

// UserCreateRequest is the body for POST /partner/v1/users. The invitation
// email sends by default; the user's Cognito login is created at their
// first sign-in via the invite, not by this call.
type UserCreateRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Phone     string `json:"phone,omitempty"`
	// RoleID references a company RBAC role — discover ids via Users.Roles.
	RoleID string `json:"roleId"`
	// LocationIDs restricts the user to specific locations; AllLocations
	// grants the unrestricted assignment instead. Mutually exclusive.
	LocationIDs  []string `json:"locationIds,omitempty"`
	AllLocations bool     `json:"allLocations,omitempty"`
	// Invite controls the invitation email (default true).
	Invite *bool `json:"invite,omitempty"`
}

// UserPatchRequest is the body for PATCH /partner/v1/users/{id}.
// Empty/omitted fields are left untouched.
type UserPatchRequest struct {
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Email     string `json:"email,omitempty"`
	Phone     string `json:"phone,omitempty"`
	RoleID    string `json:"roleId,omitempty"`
	// LocationIDs, when non-nil, REPLACES the user's location assignments.
	LocationIDs  []string `json:"locationIds,omitempty"`
	AllLocations *bool    `json:"allLocations,omitempty"`
	// Active enables/disables the user's company membership. The company
	// owner cannot be disabled.
	Active *bool `json:"active,omitempty"`
}

// RoleItem is one assignable RBAC role (GET /partner/v1/roles).
type RoleItem struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Key         string `json:"key,omitempty"`
	Description string `json:"description,omitempty"`
	IsSystem    bool   `json:"isSystem"`
}

// RolesResponse is the envelope for GET /partner/v1/roles.
type RolesResponse struct {
	Data []RoleItem `json:"data"`
}

// StatusChangeActor is who made a status change.
type StatusChangeActor struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}

// StatusChangeItem is one entry from a /status-history sub-resource.
type StatusChangeItem struct {
	Status         string            `json:"status"`
	PreviousStatus string            `json:"previousStatus,omitempty"`
	Description    string            `json:"description,omitempty"`
	ChangedAt      string            `json:"changedAt,omitempty"`
	Actor          StatusChangeActor `json:"actor"`
}

// LineItem is one curated line item from a /line-items sub-resource.
type LineItem struct {
	ID        string  `json:"id"`
	ProductID string  `json:"productId,omitempty"`
	Name      string  `json:"name"`
	ColorName string  `json:"colorName,omitempty"`
	Quantity  float64 `json:"quantity"`
	Amount    float64 `json:"amount"`
	// Status is "included", "added", or "removed".
	Status            string `json:"status"`
	IsStandardFeature bool   `json:"isStandardFeature"`
	Category          string `json:"category,omitempty"`
	Side              string `json:"side,omitempty"`
}

// LineItemTotals are counts per line item status.
type LineItemTotals struct {
	Included int `json:"included"`
	Added    int `json:"added"`
	Removed  int `json:"removed"`
}

// BuildingConfiguration is the configuration block derived from the linked
// configurator design.
type BuildingConfiguration struct {
	Model        string `json:"model,omitempty"`
	Siding       string `json:"siding,omitempty"`
	SidingColor  string `json:"sidingColor,omitempty"`
	TrimColor    string `json:"trimColor,omitempty"`
	RoofMaterial string `json:"roofMaterial,omitempty"`
	RoofColor    string `json:"roofColor,omitempty"`
}

// LineItemsResponse is the envelope for the /line-items sub-resources.
type LineItemsResponse struct {
	Data   []LineItem     `json:"data"`
	Totals LineItemTotals `json:"totals"`
	// Configuration is nil when the entity has no linked configurator design.
	Configuration *BuildingConfiguration `json:"configuration,omitempty"`
}

// ContractSummary is the read-only contract signing state from
// GET /partner/v1/orders/{id}/contract.
type ContractSummary struct {
	OrderID string `json:"orderId"`
	// Status is "draft", "out_for_signature", "partially_signed", or "completed".
	Status              string `json:"status"`
	ContractVersion     string `json:"contractVersion,omitempty"`
	ContractNumber      string `json:"contractNumber,omitempty"`
	CustomerSigned      bool   `json:"customerSigned"`
	CustomerSignedAt    string `json:"customerSignedAt,omitempty"`
	SalespersonSigned   bool   `json:"salespersonSigned"`
	SalespersonSignedAt string `json:"salespersonSignedAt,omitempty"`
	// SignedPdfDocumentID is usable against Client.Documents.Download.
	SignedPdfDocumentID string `json:"signedPdfDocumentId,omitempty"`
}

// PaymentItem is one payment record.
type PaymentItem struct {
	ID         string  `json:"id"`
	OrderID    string  `json:"orderId,omitempty"`
	CustomerID string  `json:"customerId,omitempty"`
	Amount     float64 `json:"amount"`
	// Method is "card", "ach", "cash", "check", "financed", or "manual".
	Method      string `json:"method,omitempty"`
	Status      string `json:"status,omitempty"`
	Description string `json:"description,omitempty"`
	// ProviderReference is the Stripe checkout-session / payment-intent id;
	// empty for manual records.
	ProviderReference string  `json:"providerReference,omitempty"`
	RefundedAmount    float64 `json:"refundedAmount,omitempty"`
	RefundedAt        string  `json:"refundedAt,omitempty"`
	CreatedAt         string  `json:"createdAt,omitempty"`
}

// PaymentListParams are query params for GET /partner/v1/payments.
type PaymentListParams struct {
	PaginationParams
	OrderID     string `json:"orderId,omitempty"`
	Status      string `json:"status,omitempty"`
	CreatedFrom string `json:"createdFrom,omitempty"`
	CreatedTo   string `json:"createdTo,omitempty"`
}

// DocumentItem is one file metadata row from GET /partner/v1/documents.
type DocumentItem struct {
	ID       string `json:"id"`
	Name     string `json:"name,omitempty"`
	FileName string `json:"fileName,omitempty"`
	// Type is the stored category, e.g. "Contract".
	Type      string `json:"type,omitempty"`
	MimeType  string `json:"mimeType,omitempty"`
	SizeBytes int64  `json:"sizeBytes,omitempty"`
	CreatedAt string `json:"createdAt,omitempty"`
}

// DocumentListParams are query params for GET /partner/v1/documents.
// EntityType and EntityID are required.
type DocumentListParams struct {
	PaginationParams
	// EntityType is "order", "quote", or "workOrder".
	EntityType string `json:"entityType"`
	EntityID   string `json:"entityId"`
	// Type filters by stored category, e.g. "Contract".
	Type string `json:"type,omitempty"`
}

// DocumentDownload is the short-lived presigned download from
// GET /partner/v1/documents/{id}/download. Fetch it right before
// downloading — the URL expires in ~10 minutes.
type DocumentDownload struct {
	DownloadURL string `json:"downloadUrl"`
	FileName    string `json:"fileName,omitempty"`
	ExpiresAt   string `json:"expiresAt"`
}

// EventItem is one change event from the GET /partner/v1/events feed (or a
// webhook body).
type EventItem struct {
	// ID is unique and monotonic — use as a cursor and for webhook dedupe.
	ID string `json:"id"`
	// Type is e.g. "order.status_changed", "order.cancelled", "payment.paid",
	// "contract.completed", "delivery.scheduled", "delivery.dispatched",
	// "delivery.delivered", "customer.merged", "quote.sent", "quote.expired".
	Type            string             `json:"type"`
	OccurredAt      string             `json:"occurredAt"`
	ResourceType    string             `json:"resourceType"`
	ResourceID      string             `json:"resourceId"`
	ResourceVersion int64              `json:"resourceVersion,omitempty"`
	ResourceURL     string             `json:"resourceUrl,omitempty"`
	ExternalRefs    ExternalReferences `json:"externalReferences,omitempty"`
	// Data is a compact snapshot of the changed resource.
	Data map[string]any `json:"data,omitempty"`
}

// EventListParams are query params for GET /partner/v1/events.
type EventListParams struct {
	// Cursor resumes after this event id (exclusive). Empty starts from the
	// oldest retained event.
	Cursor string `json:"cursor,omitempty"`
	// Types filters to these event types (joined with commas for you).
	Types []string `json:"-"`
	Limit int      `json:"limit,omitempty"`
}

// EventListResponse is the envelope for GET /partner/v1/events.
type EventListResponse struct {
	Data []EventItem `json:"data"`
	// NextCursor is passed as Cursor on the next call. Empty when Data is empty.
	NextCursor string `json:"nextCursor,omitempty"`
	HasMore    bool   `json:"hasMore"`
}

// EventRedeliverResponse is the body from POST /partner/v1/events/{id}/redeliver.
type EventRedeliverResponse struct {
	EventID string `json:"eventId"`
	// Enqueued is the number of active webhook subscriptions the event was
	// re-enqueued for.
	Enqueued int `json:"enqueued"`
}

// WebhookDeliveryItem is one webhook delivery attempt from
// GET /partner/v1/webhook-deliveries.
type WebhookDeliveryItem struct {
	ID             string `json:"id"`
	SubscriptionID string `json:"subscriptionId"`
	EventID        string `json:"eventId"`
	EventType      string `json:"eventType"`
	URL            string `json:"url"`
	Attempt        int    `json:"attempt"`
	StatusCode     int    `json:"statusCode,omitempty"`
	OK             bool   `json:"ok"`
	Error          string `json:"error,omitempty"`
	DurationMs     int64  `json:"durationMs,omitempty"`
	DeliveredAt    string `json:"deliveredAt"`
}

// WebhookDeliveryListParams are query params for GET /partner/v1/webhook-deliveries.
type WebhookDeliveryListParams struct {
	PaginationParams
	EventID        string `json:"eventId,omitempty"`
	SubscriptionID string `json:"subscriptionId,omitempty"`
}

// ConfiguratorSessionCreateRequest is the body for
// POST /partner/v1/configurator-sessions. Identify the design to open with
// QuoteID (reopen a saved quote) or WorkOrderID/SerialNumber (start from an
// in-stock unit); omit all three for a blank session.
type ConfiguratorSessionCreateRequest struct {
	CustomerID   string `json:"customerId"`
	LocationID   string `json:"locationId"`
	QuoteID      string `json:"quoteId,omitempty"`
	WorkOrderID  string `json:"workOrderId,omitempty"`
	SerialNumber string `json:"serialNumber,omitempty"`
	// ReturnURL must match your company's launch allowlist
	// (Settings > Developer API).
	ReturnURL string `json:"returnUrl,omitempty"`
	// TTLSeconds is the launch-token lifetime (default 900, max 3600).
	TTLSeconds int `json:"ttlSeconds,omitempty"`
}

// ConfiguratorSessionCreateResponse is the created launch session.
type ConfiguratorSessionCreateResponse struct {
	SessionID string `json:"sessionId"`
	// LaunchURL is single-use — open it in the customer's browser before
	// ExpiresAt.
	LaunchURL  string `json:"launchUrl"`
	ExpiresAt  string `json:"expiresAt"`
	CustomerID string `json:"customerId,omitempty"`
	QuoteID    string `json:"quoteId,omitempty"`
}
